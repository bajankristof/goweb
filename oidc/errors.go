package oidc

import "errors"

var (
	ErrInvalidProvider = errors.New("oidc: invalid provider")
	ErrDiscovery       = errors.New("oidc: discovery error")
	ErrExchange        = errors.New("oidc: exchange error")
	ErrInvalidToken    = errors.New("oidc: invalid token")
	ErrNonceMismatch   = errors.New("oidc: nonce mismatch")
	ErrInvalidClaims   = errors.New("oidc: invalid claims")
)
