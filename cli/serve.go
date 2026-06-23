package cli

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/bajankristof/goweb/auth/authhttp"
	"github.com/bajankristof/goweb/frontend"
	"github.com/bajankristof/goweb/http/middleware"
	"github.com/bajankristof/goweb/http/server"
	"github.com/bajankristof/goweb/user/userhttp"
	"github.com/bajankristof/goweb/wellknown/wellknownhttp"
	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v3"
)

func newServeCommand() *cli.Command {
	return &cli.Command{
		Name:   "serve",
		Usage:  "Start the server",
		Action: serve,
		Flags: []cli.Flag{
			&cli.BoolWithInverseFlag{
				Name:  flagMigrate,
				Usage: "If set, the server will automatically update the database schema to the latest version on startup",
				Value: true,
			},
		},
	}
}

func serve(ctx context.Context, cmd *cli.Command) error {
	cfg, err := loadConfig(ctx, cmd.StringSlice(flagConfig)...)
	if err != nil {
		return err
	}

	if cmd.IsSet(flagMigrate) {
		cfg.Database.AutoMigrate = cmd.Bool(flagMigrate)
	}

	app, err := newApp(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		err := app.close()
		if err != nil {
			slog.ErrorContext(ctx, "failed to shutdown cleanly", "err", err)
		}
	}()

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.TrustedProxy(cfg.Server.TrustedProxies))
	mux.Use(middleware.RequestLogger())
	mux.Use(middleware.CrossOriginProtection(cfg.Server.AllowedOrigins))
	mux.Use(middleware.CORS(cfg.Server.AllowedOrigins))
	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.Mount("/auth", authhttp.NewHandler(app.auth))
	mux.Route("/api/v1", func(r chi.Router) {
		r.Mount("/info", wellknownhttp.NewHandler(app.info()))
		r.Group(func(r chi.Router) {
			r.Use(authhttp.Authn(app.auth))
			r.Mount("/users", userhttp.NewHandler(app.user))
		})
	})
	mux.Mount("/", frontend.NewHandler())

	srv := server.New(
		mux,
		server.WithAddr(cfg.Server.Addr()),
		server.WithReadTimeout(time.Duration(cfg.Server.ReadTimeout)),
		server.WithWriteTimeout(time.Duration(cfg.Server.WriteTimeout)),
		server.WithIdleTimeout(time.Duration(cfg.Server.IdleTimeout)),
		server.WithShutdownTimeout(time.Duration(cfg.Server.ShutdownTimeout)),
	)
	err = srv.ListenAndServe(ctx)
	if err != nil {
		return err
	}

	return nil
}
