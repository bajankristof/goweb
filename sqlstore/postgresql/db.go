package postgresql

import (
	"context"
	"net/url"
	"time"

	"github.com/bajankristof/goweb/sqlstore"
	"github.com/bajankristof/goweb/sqlstore/postgresql/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var _ DBTX = (*DB)(nil)

type DB struct {
	inner  *pgxpool.Pool
	schema *sqlstore.Schema
}

type Tx = pgx.Tx

type ConnectOption func(*DB)

func Connect(ctx context.Context, url *url.URL, opts ...ConnectOption) (*DB, error) {
	pool, err := pgxpool.New(ctx, url.String())
	if err != nil {
		return nil, err
	}

	d := &DB{
		inner: pool,
		schema: &sqlstore.Schema{
			Dialect:             goose.DialectPostgres,
			FS:                  schema.FS,
			MigrationAttempts:   3,
			MigrationRetryDelay: time.Second * 2,
		},
	}

	for _, opt := range opts {
		opt(d)
	}

	return d, nil
}

func (d *DB) Close() error {
	d.inner.Close()
	return nil
}

func (d *DB) Migrate(ctx context.Context) error {
	sqld := stdlib.OpenDBFromPool(d.inner)
	defer sqld.Close()

	return d.schema.Migrate(ctx, sqld)
}

func (d *DB) Rollback(ctx context.Context) error {
	sqld := stdlib.OpenDBFromPool(d.inner)
	defer sqld.Close()

	return d.schema.Rollback(ctx, sqld)
}

func (d *DB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return d.inner.Exec(ctx, query, args...)
}

func (d *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.inner.Query(ctx, query, args...)
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.inner.QueryRow(ctx, query, args...)
}

func (d *DB) RunInTx(ctx context.Context, f func(Tx) error) error {
	tx, err := d.inner.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func WithMigrationAttempts(attempts int) ConnectOption {
	return func(d *DB) {
		d.schema.MigrationAttempts = attempts
	}
}

func WithMigrationRetryDelay(delay time.Duration) ConnectOption {
	return func(d *DB) {
		d.schema.MigrationRetryDelay = delay
	}
}
