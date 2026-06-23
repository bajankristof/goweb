package config

type OIDCConfig struct {
	IssuerURL    URL    `toml:"issuer_url" env:"OIDC_ISSUER_URL"`
	ClientID     string `toml:"client_id" env:"OIDC_CLIENT_ID"`
	ClientSecret string `toml:"client_secret" env:"OIDC_CLIENT_SECRET"`
}
