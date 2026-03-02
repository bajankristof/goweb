package cli

import (
	"context"

	"github.com/bajankristof/watchbowl/config"
	"github.com/urfave/cli/v3"
)

func newDBCommand() *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "Manage the database schema",
		Commands: []*cli.Command{
			{
				Name:   "migrate",
				Usage:  "Update the schema to the latest version",
				Action: modifyDBSchema,
			},
			{
				Name:   "rollback",
				Usage:  "Roll the schema back to the previous version",
				Action: modifyDBSchema,
			},
		},
	}
}

func modifyDBSchema(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load(cmd.String("config"))
	if err != nil {
		return err
	}

	// if err := cfg.Validate(); err != nil {
	// 	return err
	// }

	d, err := connectDB(ctx, cfg)
	if err != nil {
		return err
	}
	defer d.Close()

	if cmd.Name == "rollback" {
		return d.Rollback(ctx)
	}
	return d.Migrate(ctx)
}
