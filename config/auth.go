package config

import (
	"time"
)

type AuthConfig struct {
	Issuer       string               `toml:"issuer" env:"OIDC_ISSUER"`
	ClientID     string               `toml:"client_id" env:"OIDC_CLIENT_ID"`
	ClientSecret string               `toml:"client_secret" env:"OIDC_CLIENT_SECRET"`
	TTLSeconds   int                  `toml:"ttl"`
	Providers    []AuthProviderConfig `toml:"providers"`
}

type AuthProviderConfig struct {
	ID           string `toml:"id"`
	Issuer       string `toml:"issuer"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	TTLSeconds   int    `toml:"ttl"`
}

func (c AuthConfig) TTL() time.Duration {
	return time.Duration(c.TTLSeconds) * time.Second
}

func (c AuthProviderConfig) TTL() time.Duration {
	return time.Duration(c.TTLSeconds) * time.Second
}
