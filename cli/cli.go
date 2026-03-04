package cli

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/bajankristof/goweb/config"
	"github.com/bajankristof/goweb/db"
	"github.com/bajankristof/goweb/jwt"
	"github.com/bajankristof/goweb/oidc"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:           "goweb",
		Usage:          "GoWeb",
		ErrWriter:      io.Discard,
		ExitErrHandler: func(context.Context, *cli.Command, error) {},
		Commands: []*cli.Command{
			newServeCommand(),
			newDBCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the configuration file",
				Value:   filepath.Join("/etc", "goweb", "config.toml"),
			},
		},
	}
}

func connectDB(ctx context.Context, cfg *config.Config) (*db.DB, error) {
	return db.Connect(
		ctx,
		cfg.Database.URL,
		db.WithMigrationAttempts(cfg.Database.MigrationAttempts),
		db.WithMigrationRetryDelay(time.Duration(cfg.Database.MigrationRetryDelay)),
	)
}

func newJWTS(cfg *config.Config) (*jwt.Signer, error) {
	return jwt.NewSigner(
		cfg.Encryption.PublicKey,
		cfg.Encryption.PrivateKey,
	)
}

func newOIDR(cfg *config.Config) *oidc.Registry {
	r := oidc.NewRegistry()

	if cfg.Auth.Issuer != "" && cfg.Auth.ClientID != "" && cfg.Auth.ClientSecret != "" {
		r.Register(
			"",
			cfg.Auth.Issuer,
			cfg.Auth.ClientID,
			cfg.Auth.ClientSecret,
			oidc.WithTTL(cfg.Auth.TTL()),
		)
	}

	for _, p := range cfg.Auth.Providers {
		r.Register(
			p.ID,
			p.Issuer,
			p.ClientID,
			p.ClientSecret,
			oidc.WithName(p.Name),
			oidc.WithTTL(p.TTL()),
		)
	}

	return r
}
