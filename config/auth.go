package config

type AuthConfig struct {
	AccessTokenTTL  Duration `toml:"access_token_ttl" env:"-"`
	RefreshTokenTTL Duration `toml:"refresh_token_ttl" env:"-"`
}
