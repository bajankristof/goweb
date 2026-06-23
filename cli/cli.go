package cli

import (
	"context"

	"github.com/bajankristof/goweb/config"
	"github.com/urfave/cli/v3"
)

const (
	flagConfig      = "config"
	flagAliasConfig = "c"
	flagMigrate     = "migrate"
)

func New() *cli.Command {
	return &cli.Command{
		Name:    "goweb",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    flagConfig,
				Aliases: []string{flagAliasConfig},
				Usage:   "Path to configuration file(s)",
				Value:   []string{"goweb.toml"},
			},
		},
		Commands: []*cli.Command{
			newServeCommand(),
		},
	}
}

func loadConfig(ctx context.Context, paths ...string) (*config.Config, error) {
	cfg := config.Default()
	if err := cfg.Load(ctx, paths...); err != nil {
		return nil, err
	}
	return cfg, nil
}
