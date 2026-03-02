package config

import (
	"fmt"
	"net/netip"
)

type ServerConfig struct {
	Host            string         `toml:"host" env:"HOST"`
	Port            int            `toml:"port" env:"PORT"`
	ReadTimeout     Duration       `toml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout    Duration       `toml:"write_timeout" env:"WRITE_TIMEOUT"`
	IdleTimeout     Duration       `toml:"idle_timeout" env:"IDLE_TIMEOUT"`
	ShutdownTimeout Duration       `toml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
	AllowedOrigins  []string       `toml:"allowed_origins" env:"ALLOWED_ORIGINS" envSeparator:","`
	TrustedProxies  []netip.Prefix `toml:"trusted_proxies" env:"TRUSTED_PROXIES" envSeparator:","`
}

func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
