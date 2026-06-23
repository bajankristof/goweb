package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bajankristof/trustedproxy"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Auth     AuthConfig            `toml:"auth"`
	Database DatabaseConfig        `toml:"database"`
	OIDC     map[string]OIDCConfig `toml:"oidc" env:"-"`
	JWT      JWTConfig             `toml:"jwt"`
	Server   ServerConfig          `toml:"server"`
}

func Default() *Config {
	return &Config{
		Auth: AuthConfig{
			AccessTokenTTL:  Duration(15 * time.Minute),
			RefreshTokenTTL: Duration(7 * 24 * time.Hour),
		},
		Database: DatabaseConfig{
			AutoMigrate:         true,
			MigrationAttempts:   3,
			MigrationRetryDelay: Duration(5 * time.Second),
		},
		OIDC: make(map[string]OIDCConfig),
		JWT:  JWTConfig{},
		Server: ServerConfig{
			Host:            "",
			Port:            8080,
			ReadTimeout:     Duration(15 * time.Second),
			WriteTimeout:    Duration(15 * time.Second),
			IdleTimeout:     Duration(60 * time.Second),
			ShutdownTimeout: Duration(15 * time.Second),
			AllowedOrigins:  []string{},
			TrustedProxies:  trustedproxy.DefaultPrefixes(),
		},
	}
}

func (c *Config) Load(ctx context.Context, paths ...string) error {
	for _, path := range paths {
		_, err := toml.DecodeFile(path, c)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config: load error: %w", err)
		}
	}

	err := env.Parse(c)
	if err != nil {
		return fmt.Errorf("config: env parse error: %w", err)
	}

	idpc := OIDCConfig{}
	err = env.Parse(&idpc)
	if err != nil {
		return fmt.Errorf("config: env parse error: %w", err)
	}

	if idpc.IssuerURL.String() != "" {
		c.OIDC["default"] = idpc
	}

	return nil
}
