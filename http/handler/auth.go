package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/bajankristof/watchbowl/db"
	"github.com/bajankristof/watchbowl/http/cookieutil"
	"github.com/bajankristof/watchbowl/http/dto"
	"github.com/bajankristof/watchbowl/http/requestutil"
	"github.com/bajankristof/watchbowl/jwt"
	"github.com/bajankristof/watchbowl/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"golang.org/x/oauth2"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type authState struct {
	jwt.RegisteredClaims
	Provider     string `json:"idp,omitempty"`
	CallbackURL  string `json:"cbk,omitempty"`
	CodeVerifier string `json:"cv,omitempty"`
}

type authURLGenerator interface {
	AuthURL(ctx context.Context, idp, callbackURL, state, nonce, verifier string) (string, error)
}

type authCodeExchanger interface {
	Exchange(ctx context.Context, idp, callbackURL, code, nonce, verifier string) (*oidc.UserInfo, error)
}

type authSignOutQueries interface {
	RevokeSession(ctx context.Context, refreshTokenHash string) (db.Session, error)
}

type authCallbackQueries interface {
	CreateWebUser(ctx context.Context, arg db.CreateWebUserParams) (db.User, error)
}

type authRefreshQueries interface {
	RotateSession(ctx context.Context, arg db.RotateSessionParams) (db.Session, error)
}

func authHandler(dbq *db.Queries, jwts *jwt.Signer, oidr *oidc.Registry) http.Handler {
	r := chi.NewRouter()
	r.Get("/well-known", authWellKnownHandler(oidr))
	r.Get("/signin", authSignInHandler(jwts, oidr))
	r.Get("/signin/{idp}", authSignInHandler(jwts, oidr))
	r.Get("/signout", authSignOutHandler(dbq))
	r.Get("/callback", authCallbackHandler(dbq, jwts, oidr))
	r.Get("/refresh", authRefreshHandler(dbq, jwts))
	return r
}

func authWellKnownHandler(oidr *oidc.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, &dto.AuthWellKnownResponse{
			Providers: oidr.All(),
		})
	}
}

func authSignInHandler(jwts *jwt.Signer, oidr authURLGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbackURL, ok := requestutil.RewriteURL(r, "/callback", "/signin", "/signin/{idp}")
		if !ok {
			slog.ErrorContext(r.Context(), "failed to determine callback URL")
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		state := authState{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "watchbowl",
				Subject:   rand.Text(),
				ExpiresAt: jwt.After(5 * time.Minute),
				IssuedAt:  jwt.Now(),
			},
			Provider:     chi.URLParam(r, "idp"),
			CallbackURL:  callbackURL.String(),
			CodeVerifier: oauth2.GenerateVerifier(),
		}

		stateStr, err := jwts.Sign(state)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to sign auth state", "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		authURL, err := oidr.AuthURL(
			r.Context(),
			state.Provider,
			state.CallbackURL,
			stateStr,
			state.Subject,
			state.CodeVerifier,
		)
		if _, ok := errors.AsType[*oidc.RegistryError](err); ok {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, http.StatusText(http.StatusNotFound))
			return
		} else if _, ok := errors.AsType[*oidc.SyncError](err); ok {
			slog.ErrorContext(r.Context(), "failed to sync OIDC provider", "err", err)
			render.Status(r, http.StatusBadGateway)
			render.PlainText(w, r, http.StatusText(http.StatusBadGateway))
			return
		} else if err != nil {
			slog.ErrorContext(r.Context(), "failed to generate auth URL", "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

func authSignOutHandler(dbq authSignOutQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL, ok := requestutil.RewriteURL(r, "/", "/signout")
		if !ok {
			slog.ErrorContext(r.Context(), "failed to determine base URL")
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		refreshToken := cookieutil.Get(r, "refresh_token")
		if refreshToken == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		refreshTokenHash := newRefreshTokenHash(refreshToken)
		_, err := dbq.RevokeSession(r.Context(), refreshTokenHash)
		if errors.Is(err, db.ErrNoRows) {
			slog.WarnContext(r.Context(), "session does not exist or already revoked", "refresh_token_hash", refreshTokenHash)
		} else if err != nil {
			slog.ErrorContext(r.Context(), "failed to revoke session", "refresh_token_hash", refreshTokenHash, "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		cookieutil.Set(w, r, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})
		cookieutil.Set(w, r, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     baseURL.Path,
			HttpOnly: true,
			MaxAge:   -1,
		})

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func authCallbackHandler(dbq authCallbackQueries, jwts *jwt.Signer, oidr authCodeExchanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL, ok := requestutil.RewriteURL(r, "/", "/callback")
		if !ok {
			slog.ErrorContext(r.Context(), "failed to determine base URL")
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		}

		stateStr := r.URL.Query().Get("state")
		if stateStr == "" {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		}

		state := authState{}
		err := jwts.Verify(stateStr, &state)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		}

		userInfo, err := oidr.Exchange(
			r.Context(),
			state.Provider,
			state.CallbackURL,
			code,
			state.Subject,
			state.CodeVerifier,
		)
		if _, ok := errors.AsType[*oidc.SyncError](err); ok {
			slog.ErrorContext(r.Context(), "failed to sync OIDC provider", "err", err)
			render.Status(r, http.StatusBadGateway)
			render.PlainText(w, r, http.StatusText(http.StatusBadGateway))
			return
		} else if _, ok := errors.AsType[*oidc.ExchangeError](err); ok {
			slog.DebugContext(r.Context(), "rejected auth code exchange", "err", err)
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		} else if err != nil {
			slog.ErrorContext(r.Context(), "failed to exchange auth code", "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		refreshToken := rand.Text()
		user, err := dbq.CreateWebUser(r.Context(), db.CreateWebUserParams{
			OpenID:           userInfo.ID,
			Provider:         state.Provider,
			Email:            userInfo.Email,
			DisplayName:      null.StringFrom(userInfo.DisplayName),
			RefreshTokenHash: newRefreshTokenHash(refreshToken),
			UserAgent:        r.UserAgent(),
		})
		if errors.Is(err, db.ErrNoRows) {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		} else if err != nil {
			slog.ErrorContext(r.Context(), "failed to create user", "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		accessToken, err := newAccessToken(jwts, user.UserID)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to generate access token", "user_id", user.UserID, "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		setAccessToken(w, r, accessToken)
		setRefreshToken(w, r, refreshToken, baseURL.Path)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func authRefreshHandler(dbq authRefreshQueries, jwts *jwt.Signer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL, ok := requestutil.RewriteURL(r, "/", "/refresh")
		if !ok {
			slog.ErrorContext(r.Context(), "failed to determine base URL")
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		refreshToken := cookieutil.Get(r, "refresh_token")
		if refreshToken == "" {
			unsetAccessToken(w, r)
			unsetRefreshToken(w, r, baseURL.Path)
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		}

		newRefreshToken := rand.Text()
		refreshTokenHash := newRefreshTokenHash(refreshToken)
		sess, err := dbq.RotateSession(r.Context(), db.RotateSessionParams{
			RefreshTokenHash:    refreshTokenHash,
			NewRefreshTokenHash: newRefreshTokenHash(newRefreshToken),
		})
		if errors.Is(err, db.ErrNoRows) {
			slog.WarnContext(r.Context(), "session does not exist or already revoked", "refresh_token_hash", refreshTokenHash)
			unsetAccessToken(w, r)
			unsetRefreshToken(w, r, baseURL.Path)
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, http.StatusText(http.StatusUnauthorized))
			return
		} else if err != nil {
			slog.ErrorContext(r.Context(), "failed to rotate session", "refresh_token_hash", refreshTokenHash, "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		accessToken, err := newAccessToken(jwts, sess.UserID)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to generate access token", "user_id", sess.UserID, "err", err)
			unsetAccessToken(w, r)
			unsetRefreshToken(w, r, baseURL.Path)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
			return
		}

		setAccessToken(w, r, accessToken)
		setRefreshToken(w, r, newRefreshToken, baseURL.Path)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func newAccessToken(
	jwts *jwt.Signer,
	userID uuid.UUID,
) (string, error) {
	return jwts.Sign(jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.After(accessTokenTTL),
		IssuedAt:  jwt.Now(),
	})
}

func newRefreshTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func setAccessToken(w http.ResponseWriter, r *http.Request, token string) {
	cookieutil.Set(w, r, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(accessTokenTTL),
	})
}

func unsetAccessToken(w http.ResponseWriter, r *http.Request) {
	cookieutil.Set(w, r, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func setRefreshToken(w http.ResponseWriter, r *http.Request, token, path string) {
	cookieutil.Set(w, r, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     path,
		HttpOnly: true,
		Expires:  time.Now().Add(refreshTokenTTL),
	})
}

func unsetRefreshToken(w http.ResponseWriter, r *http.Request, path string) {
	cookieutil.Set(w, r, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     path,
		MaxAge:   -1,
	})
}
