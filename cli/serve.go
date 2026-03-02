package cli

import (
	"context"
	"time"

	"github.com/bajankristof/watchbowl/config"
	"github.com/bajankristof/watchbowl/db"
	"github.com/bajankristof/watchbowl/http/handler"
	"github.com/bajankristof/watchbowl/http/server"
	"github.com/urfave/cli/v3"
)

func newServeCommand() *cli.Command {
	return &cli.Command{
		Name:   "serve",
		Usage:  "Start the WatchBowl server",
		Action: serve,
		Flags: []cli.Flag{
			&cli.BoolWithInverseFlag{
				Name:  "migrate",
				Usage: "If set, the server will automatically update the database schema to the latest version on startup",
				Value: true,
			},
		},
	}
}

func serve(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load(cmd.String("config"))
	if err != nil {
		return err
	}

	// if err := cfg.Validate(); err != nil {
	// 	return err
	// }

	dbc, err := connectDB(ctx, cfg)
	if err != nil {
		return err
	}
	defer dbc.Close()

	if cmd.Bool("migrate") {
		if err := dbc.Migrate(ctx); err != nil {
			return err
		}
	}

	jwts, err := newJWTS(cfg)
	if err != nil {
		return err
	}

	hdlr := handler.New(
		db.New(dbc),
		jwts,
		newOIDR(cfg),
		cfg.Server.AllowedOrigins,
		cfg.Server.TrustedProxies,
	)

	srv := server.New(
		hdlr,
		server.WithAddr(cfg.Server.Addr()),
		server.WithReadTimeout(time.Duration(cfg.Server.ReadTimeout)),
		server.WithWriteTimeout(time.Duration(cfg.Server.WriteTimeout)),
		server.WithIdleTimeout(time.Duration(cfg.Server.IdleTimeout)),
		server.WithShutdownTimeout(time.Duration(cfg.Server.ShutdownTimeout)),
	)

	return srv.ListenAndServe(ctx)
}
