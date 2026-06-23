package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bajankristof/goweb/auth"
	"github.com/bajankristof/goweb/config"
	"github.com/bajankristof/goweb/session"
	"github.com/bajankristof/goweb/sqlstore"
	"github.com/bajankristof/goweb/sqlstore/postgresql"
	"github.com/bajankristof/goweb/sqlstore/sqlite"
	"github.com/bajankristof/goweb/user"
)

type db struct {
	sqld
	userStore    user.Store
	sessionStore session.Store
	authUoW      auth.UnitOfWork
}

type sqld interface {
	Migrate(ctx context.Context) error
	Rollback(ctx context.Context) error
	Close() error
}

func openDB(ctx context.Context, cfg *config.Config) (db, error) {
	switch cfg.Database.URL.Scheme {
	case "postgres", "postgresql":
		d, err := postgresql.Connect(
			ctx,
			cfg.Database.URL.Unwrap(),
			postgresql.WithMigrationAttempts(cfg.Database.MigrationAttempts),
			postgresql.WithMigrationRetryDelay(time.Duration(cfg.Database.MigrationRetryDelay)),
		)
		if err != nil {
			return db{}, fmt.Errorf("connect postgresql: %w", err)
		}

		if cfg.Database.AutoMigrate {
			err = d.Migrate(ctx)
			if err != nil {
				return db{}, errors.Join(fmt.Errorf("migrate postgresql: %w", err), d.Close())
			}
		}

		return db{
			sqld:         d,
			userStore:    postgresql.NewUserStore(d),
			sessionStore: postgresql.NewSessionStore(d),
			authUoW: sqlstore.NewUnitOfWork(d, func(tx postgresql.Tx) auth.Stores {
				return auth.Stores{
					UserStore:    postgresql.NewUserStore(tx),
					SessionStore: postgresql.NewSessionStore(tx),
				}
			}),
		}, nil
	case "sqlite":
		d, err := sqlite.Open(
			ctx,
			cfg.Database.URL.Unwrap(),
			sqlite.WithMigrationAttempts(cfg.Database.MigrationAttempts),
			sqlite.WithMigrationRetryDelay(time.Duration(cfg.Database.MigrationRetryDelay)),
		)
		if err != nil {
			return db{}, fmt.Errorf("open sqlite: %w", err)
		}

		if cfg.Database.AutoMigrate {
			err = d.Migrate(ctx)
			if err != nil {
				return db{}, errors.Join(fmt.Errorf("migrate sqlite: %w", err), d.Close())
			}
		}

		return db{
			sqld:         d,
			userStore:    sqlite.NewUserStore(d),
			sessionStore: sqlite.NewSessionStore(d),
			authUoW: sqlstore.NewUnitOfWork(d, func(tx sqlite.Tx) auth.Stores {
				return auth.Stores{
					UserStore:    sqlite.NewUserStore(tx),
					SessionStore: sqlite.NewSessionStore(tx),
				}
			}),
		}, nil
	default:
		return db{}, fmt.Errorf("unsupported database scheme %q", cfg.Database.URL.Scheme)
	}
}
