package sqlstore

import (
	"context"
	"database/sql"
	"io/fs"
	"log/slog"
	"time"

	"github.com/pressly/goose/v3"
)

const TableSchemaMigrations = "schema_migrations"

type Schema struct {
	Dialect             goose.Dialect
	FS                  fs.FS
	MigrationAttempts   int
	MigrationRetryDelay time.Duration
}

func (s *Schema) Migrate(ctx context.Context, db *sql.DB) error {
	return s.modify(ctx, db, func(p *goose.Provider) error {
		_, err := p.Up(ctx)
		return err
	})
}

func (s *Schema) Rollback(ctx context.Context, db *sql.DB) error {
	return s.modify(ctx, db, func(p *goose.Provider) error {
		_, err := p.Down(ctx)
		return err
	})
}

func (s *Schema) modify(ctx context.Context, db *sql.DB, f func(p *goose.Provider) error) error {
	var err error
	p, err := goose.NewProvider(
		s.Dialect,
		db,
		s.FS,
		goose.WithSlog(slog.Default()),
		goose.WithTableName(TableSchemaMigrations),
	)
	if err != nil {
		return err
	}

	for range s.MigrationAttempts {
		err = f(p)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.MigrationRetryDelay):
		}
	}

	return err
}
