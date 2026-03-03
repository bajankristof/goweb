package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	inner        *oidc.Provider
	id           string
	name         string
	issuer       string
	clientID     string
	clientSecret string
	scopes       []string
	ttl          time.Duration
	updatedAt    time.Time
	mu           sync.RWMutex
}

type ProviderOption func(*Provider)

func NewProvider(
	id string,
	issuer string,
	clientID string,
	clientSecret string,
	opts ...ProviderOption,
) *Provider {
	p := &Provider{
		id:           id,
		issuer:       issuer,
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Provider) AuthURL(
	ctx context.Context,
	callbackURL,
	state,
	nonce,
	verifier string,
) (string, error) {
	cfg, err := p.config(ctx, callbackURL)
	if err != nil {
		return "", err
	}

	return cfg.AuthCodeURL(
		state,
		oidc.Nonce(nonce),
		oauth2.S256ChallengeOption(verifier),
	), nil
}

func (p *Provider) Exchange(
	ctx context.Context,
	callbackURL,
	code,
	nonce,
	verifier string,
) (*UserInfo, error) {
	cfg, err := p.config(ctx, callbackURL)
	if err != nil {
		return nil, err
	}

	token, err := cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		if _, ok := errors.AsType[*oauth2.RetrieveError](err); ok {
			return nil, &ExchangeError{err}
		}
		return nil, err
	}

	idToken, err := p.verify(ctx, token, nonce)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{}
	if err := idToken.Claims(userInfo); err != nil {
		return nil, fmt.Errorf("unmarshal user info: %w", err)
	}

	return userInfo, nil
}

// MarshalJSON implements json.Marshaler.
func (p *Provider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"id":     p.id,
		"name":   p.name,
		"issuer": p.issuer,
	})
}

func (p *Provider) config(ctx context.Context, callbackURL string) (*oauth2.Config, error) {
	err := p.sync(ctx)
	if err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	return &oauth2.Config{
		ClientID:     p.clientID,
		ClientSecret: p.clientSecret,
		Endpoint:     p.inner.Endpoint(),
		RedirectURL:  callbackURL,
		Scopes:       p.scopes,
	}, nil
}

func (p *Provider) sync(ctx context.Context) error {
	p.mu.RLock()
	stale := p.inner == nil || time.Since(p.updatedAt) > p.ttl
	p.mu.RUnlock()

	if !stale {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.inner != nil && time.Since(p.updatedAt) <= p.ttl {
		return nil
	}

	inner, err := oidc.NewProvider(ctx, p.issuer)
	if err != nil {
		return &SyncError{fmt.Errorf("sync provider %q: %w", p.id, err)}
	}

	p.inner = inner
	p.updatedAt = time.Now()

	return nil
}

func (p *Provider) verify(ctx context.Context, token *oauth2.Token, nonce string) (*oidc.IDToken, error) {
	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, &ExchangeError{errors.New("verify token: ID token not found")}
	}

	if err := p.sync(ctx); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	idToken, err := p.inner.Verifier(&oidc.Config{ClientID: p.clientID}).Verify(ctx, idTokenStr)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	if idToken.Nonce != nonce {
		return nil, &ExchangeError{errors.New("verify token: nonce mismatch")}
	}

	return idToken, nil
}

func WithTTL(ttl time.Duration) ProviderOption {
	return func(p *Provider) {
		p.ttl = ttl
	}
}

func WithName(name string) ProviderOption {
	return func(p *Provider) {
		p.name = name
	}
}

type UserInfo struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	DisplayName   string `json:"name"`
	Picture       string `json:"picture"`
	Timezone      string `json:"zoneinfo"`
	Locale        string `json:"locale"`
}
