package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bajankristof/trustedproxy"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Database   DatabaseConfig   `toml:"database"`
	Server     ServerConfig     `toml:"server"`
	Encryption EncryptionConfig `toml:"encryption"`
	Auth       AuthConfig       `toml:"auth"`
}

func New() *Config {
	return &Config{
		Database: DatabaseConfig{
			MigrationAttempts:   3,
			MigrationRetryDelay: Duration(5 * time.Second),
		},
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

func Load(path string) (*Config, error) {
	c := New()

	_, err := toml.DecodeFile(path, c)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}

	err = env.Parse(c)
	if err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	err = c.Encryption.LoadKeys()
	if err != nil {
		return nil, err
	}

	return c, nil
}
