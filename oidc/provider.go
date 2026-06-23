package oidc

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	issuerURL    string
	clientID     string
	clientSecret string
	scopes       []string
}

type ProviderOption func(*Provider)

func NewProvider(
	issuerURL,
	clientID,
	clientSecret string,
	opts ...ProviderOption,
) *Provider {
	i := &Provider{
		issuerURL:    issuerURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (p *Provider) AuthCodeURL(
	ctx context.Context,
	callbackURL,
	state,
	nonce,
	codeVerifier string,
) (string, error) {
	op, err := oidc.NewProvider(ctx, p.issuerURL)
	if err != nil {
		return "", newSentinelError(ErrDiscovery, err)
	}

	c := oauth2.Config{
		ClientID:     p.clientID,
		ClientSecret: p.clientSecret,
		Endpoint:     op.Endpoint(),
		Scopes:       p.scopes,
		RedirectURL:  callbackURL,
	}

	return c.AuthCodeURL(
		state,
		oidc.Nonce(nonce),
		oauth2.S256ChallengeOption(codeVerifier),
	), nil
}

func (p *Provider) Exchange(
	ctx context.Context,
	callbackURL,
	code,
	nonce,
	codeVerifier string,
) (User, error) {
	u := User{}
	op, err := oidc.NewProvider(ctx, p.issuerURL)
	if err != nil {
		return u, newSentinelError(ErrDiscovery, err)
	}

	c := oauth2.Config{
		ClientID:     p.clientID,
		ClientSecret: p.clientSecret,
		Endpoint:     op.Endpoint(),
		Scopes:       p.scopes,
		RedirectURL:  callbackURL,
	}

	t, err := c.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return u, newSentinelError(ErrExchange, err)
	}

	v := op.VerifierContext(ctx, &oidc.Config{ClientID: p.clientID})
	id, err := v.Verify(ctx, t.Extra("id_token").(string))
	if err != nil {
		return u, newSentinelError(ErrInvalidToken, err)
	}

	if id.Nonce != nonce {
		return u, ErrNonceMismatch
	}

	err = id.Claims(&u)
	if err != nil {
		return u, newSentinelError(ErrInvalidClaims, err)
	}

	return u, nil
}

func WithScopes(scopes ...string) ProviderOption {
	return func(p *Provider) {
		p.scopes = scopes
	}
}

func newSentinelError(sentinel error, err error) error {
	m := err.Error()
	if m, ok := strings.CutPrefix(m, "oidc: "); ok {
		return fmt.Errorf("%w: %s", sentinel, m)
	}

	return fmt.Errorf("%w: %s", sentinel, m)
}
