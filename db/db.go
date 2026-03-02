package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bajankristof/watchbowl/db/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
)

type DB struct {
	pool                *pgxpool.Pool
	migrationAttempts   int
	migrationRetryDelay time.Duration
}

type ConnectOption func(*DB)

func Connect(ctx context.Context, url string, opts ...ConnectOption) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	cfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &Logger{},
		LogLevel: tracelog.LogLevelTrace,
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	d := &DB{
		pool:                pool,
		migrationAttempts:   3,
		migrationRetryDelay: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(d)
	}

	return d, nil
}

func (d *DB) Close() {
	d.pool.Close()
}

func (d *DB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}

func (d *DB) Migrate(ctx context.Context) error {
	return d.modifySchema(ctx, schema.Migrate)
}

func (d *DB) Rollback(ctx context.Context) error {
	return d.modifySchema(ctx, schema.Rollback)
}

func (d *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

func (d *DB) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func WithMigrationAttempts(attempts int) ConnectOption {
	return func(d *DB) {
		d.migrationAttempts = attempts
	}
}

func WithMigrationRetryDelay(delay time.Duration) ConnectOption {
	return func(d *DB) {
		d.migrationRetryDelay = delay
	}
}

func (d *DB) modifySchema(ctx context.Context, f func(context.Context, *sql.DB) error) error {
	var err error
	for i := 0; i < d.migrationAttempts; i++ {
		if i > 0 {
			time.Sleep(d.migrationRetryDelay)
		}

		sdb := stdlib.OpenDBFromPool(d.pool)
		err = f(ctx, sdb)
		sdb.Close()
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("modify schema after %d attempts: %w", d.migrationAttempts, err)
}
