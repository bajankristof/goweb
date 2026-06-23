package config

type JWTConfig struct {
	SigningKey PrivateKey `toml:"signing_key" env:"JWT_SIGNING_KEY"`
}
