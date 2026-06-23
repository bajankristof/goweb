package config

import (
	"net"
	"net/netip"
	"strconv"
)

type ServerConfig struct {
	Host            string         `toml:"host" env:"HOST"`
	Port            int            `toml:"port" env:"PORT"`
	ReadTimeout     Duration       `toml:"read_timeout" env:"-"`
	WriteTimeout    Duration       `toml:"write_timeout" env:"-"`
	IdleTimeout     Duration       `toml:"idle_timeout" env:"-"`
	ShutdownTimeout Duration       `toml:"shutdown_timeout" env:"-"`
	AllowedOrigins  []string       `toml:"allowed_origins" env:"-" envSeparator:","`
	TrustedProxies  []netip.Prefix `toml:"trusted_proxies" env:"-"`
}

func (c ServerConfig) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
