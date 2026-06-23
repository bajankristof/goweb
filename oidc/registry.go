package oidc

import (
	"maps"
	"slices"
)

type Registry map[string]*Provider

func NewRegistry() Registry {
	return make(map[string]*Provider)
}

func (r Registry) Add(name string, p *Provider) {
	r[name] = p
}

func (r Registry) Get(name string) (*Provider, error) {
	i, ok := r[name]
	if !ok {
		return nil, ErrInvalidProvider
	}
	return i, nil
}

func (r Registry) Remove(name string) {
	delete(r, name)
}

func (r Registry) Names() []string {
	return slices.Collect(maps.Keys(r))
}
