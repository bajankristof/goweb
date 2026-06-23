package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/bajankristof/goweb/jwt"
	"github.com/bajankristof/goweb/oidc"
	"github.com/bajankristof/goweb/session"
	"github.com/bajankristof/goweb/user"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

const (
	audOAuthState  = "goweb-oauth-state"
	audAccessToken = "goweb-access"
)

type Service struct {
	stores Stores
	uow    UnitOfWork
	signer *jwt.Signer
	idps   oidc.Registry

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(
	stores Stores,
	uow UnitOfWork,
	signer *jwt.Signer,
	idps oidc.Registry,
	opts ...ServiceOption,
) *Service {
	s := &Service{
		stores: stores,
		uow:    uow,
		signer: signer,
		idps:   idps,

		accessTokenTTL:  time.Minute * 15,
		refreshTokenTTL: time.Hour * 24 * 7,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type StateClaims struct {
	jwt.RegisteredClaims
	IDP          string `json:"idp"`
	CallbackURL  string `json:"cbk"`
	CodeVerifier string `json:"pkce"`
}

func (s *Service) BeginSignIn(ctx context.Context, via, nonce, callbackURL string) (string, error) {
	idp, err := s.idps.Get(via)
	if err != nil {
		return "", err
	}

	codeVerifier := rand.Text()
	state, err := s.signer.Sign(StateClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "goweb",
			Subject:   nonce,
			Audience:  []string{audOAuthState},
			ExpiresAt: jwt.After(time.Minute * 5),
			IssuedAt:  jwt.Now(),
		},
		IDP:          via,
		CodeVerifier: codeVerifier,
		CallbackURL:  callbackURL,
	})
	if err != nil {
		return "", fmt.Errorf("auth: state signing error: %w", err)
	}

	authURL, err := idp.AuthCodeURL(ctx, callbackURL, state, nonce, codeVerifier)
	if err != nil {
		return "", fmt.Errorf("auth: auth code URL generation error: %w", err)
	}

	return authURL, nil
}

func (s *Service) CompleteSignIn(ctx context.Context, nonce, state, code, userAgent string) (User, error) {
	var u User
	var claims StateClaims
	err := s.signer.Verify(state, &claims, jwt.WithAudience(audOAuthState))
	if err != nil {
		return u, fmt.Errorf("%w: %w", ErrInvalidState, err)
	}

	if claims.Subject != nonce {
		return u, fmt.Errorf("%w: nonce mismatch", ErrInvalidState)
	}

	idp, err := s.idps.Get(claims.IDP)
	if err != nil {
		return u, err
	}

	ou, err := idp.Exchange(ctx, claims.CallbackURL, code, claims.Subject, claims.CodeVerifier)
	if err != nil {
		return u, err
	}

	return s.SignUp(ctx, SignUpParams{
		OpenID:      ou.ID,
		IDP:         claims.IDP,
		Email:       ou.Email,
		DisplayName: null.StringFrom(ou.DisplayName),
		UserAgent:   userAgent,
	})
}

type SignUpParams struct {
	OpenID      string
	IDP         string
	Email       string
	DisplayName null.String
	UserAgent   string
}

func (s *Service) SignUp(ctx context.Context, params SignUpParams) (User, error) {
	var u User
	err := s.uow.Do(ctx, func(stores Stores) error {
		var err error
		u.User, err = stores.UserStore.Create(ctx, user.CreateParams{
			OpenID:      params.OpenID,
			IDP:         params.IDP,
			Email:       params.Email,
			DisplayName: params.DisplayName,
		})
		if err != nil {
			return err
		}

		u.RefreshToken = rand.Text()
		u.Session, err = stores.SessionStore.Create(ctx, session.CreateParams{
			UserID:           u.ID,
			UserAgent:        params.UserAgent,
			RefreshTokenHash: newHexSum(u.RefreshToken),
			ExpiresAt:        time.Now().Add(s.refreshTokenTTL),
		})
		if err != nil {
			return err
		}

		u.AccessToken, err = s.newAccessToken(u.ID)

		return err
	})
	return u, err
}

func (s *Service) SignOut(ctx context.Context, refreshToken string) error {
	return s.uow.Do(ctx, func(stores Stores) error {
		sess, err := stores.SessionStore.GetByRefreshTokenHash(ctx, newHexSum(refreshToken))
		if err != nil {
			return err
		}

		err = stores.SessionStore.Revoke(ctx, sess.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) Verify(ctx context.Context, accessToken string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	err := s.signer.Verify(accessToken, &claims, jwt.WithAudience(audAccessToken))
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth: access token verification error: %w", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth: invalid user ID in access token: %w", err)
	}

	return userID, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken,
	userAgent string,
) (User, error) {
	var u User
	err := s.uow.Do(ctx, func(stores Stores) error {
		sess, err := stores.SessionStore.GetByRefreshTokenHash(ctx, newHexSum(refreshToken))
		if err != nil {
			return err
		}

		if sess.ExpiresAt.Before(time.Now()) {
			return session.ErrExpired
		}

		u.User, err = stores.UserStore.Get(ctx, sess.UserID)
		if err != nil {
			return err
		}

		u.RefreshToken = rand.Text()
		u.Session, err = stores.SessionStore.Refresh(ctx, session.RefreshParams{
			ID:               sess.ID,
			UserAgent:        userAgent,
			RefreshTokenHash: newHexSum(u.RefreshToken),
			ExpiresAt:        time.Now().Add(s.refreshTokenTTL),
		})
		if err != nil {
			return err
		}

		u.AccessToken, err = s.newAccessToken(u.ID)

		return err
	})
	return u, err
}

func (s *Service) newAccessToken(userID uuid.UUID) (string, error) {
	t, err := s.signer.Sign(jwt.RegisteredClaims{
		Issuer:    "goweb",
		Subject:   userID.String(),
		Audience:  []string{audAccessToken},
		ExpiresAt: jwt.After(s.accessTokenTTL),
		IssuedAt:  jwt.Now(),
	})
	if err != nil {
		return "", fmt.Errorf("auth: access token signing error: %w", err)
	}

	return t, nil
}

type ServiceOption func(*Service)

func WithAccessTokenTTL(ttl time.Duration) func(*Service) {
	return func(s *Service) {
		s.accessTokenTTL = ttl
	}
}

func WithRefreshTokenTTL(ttl time.Duration) func(*Service) {
	return func(s *Service) {
		s.refreshTokenTTL = ttl
	}
}

func newHexSum(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
