package oidc

import (
	"context"
	"fmt"
	"sync"
)

type Registry struct {
	mu   sync.RWMutex
	idps map[string]*Provider
}

type RegistryOption func(*Registry)

func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{idps: make(map[string]*Provider)}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Registry) Register(id, issuer, clientID, clientSecret string, opts ...ProviderOption) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.idps[id] = NewProvider(id, issuer, clientID, clientSecret, opts...)
}

func (r *Registry) All() []*Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idps := make([]*Provider, 0, len(r.idps))

	for _, idp := range r.idps {
		idps = append(idps, idp)
	}

	return idps
}

func (r *Registry) AuthURL(ctx context.Context, id, callbackURL, state, nonce, verifier string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idp, ok := r.idps[id]
	if !ok {
		return "", &RegistryError{Message: fmt.Sprintf("auth URL: unknown provider %q", id)}
	}

	return idp.AuthURL(ctx, callbackURL, state, nonce, verifier)
}

func (r *Registry) Exchange(ctx context.Context, id, callbackURL, code, nonce, verifier string) (*UserInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idp, ok := r.idps[id]
	if !ok {
		return nil, &RegistryError{Message: fmt.Sprintf("exchange: unknown provider %q", id)}
	}

	return idp.Exchange(ctx, callbackURL, code, nonce, verifier)
}

func WithProvider(
	id string,
	issuer string,
	clientID string,
	clientSecret string,
	opts... ProviderOption,
) RegistryOption {
	return func(r *Registry) {
		r.Register(id, issuer, clientID, clientSecret, opts...)
	}
}
