package sqlite

import (
	"context"
	"database/sql"
	"net/url"
	"path/filepath"
	"time"

	"github.com/bajankristof/goweb/sqlstore"
	"github.com/bajankristof/goweb/sqlstore/sqlite/schema"
	_ "modernc.org/sqlite"
)

var _ DBTX = (*DB)(nil)

type DB struct {
	inner  *sql.DB
	schema *sqlstore.Schema
}

type Tx = *sql.Tx

type OpenOption func(*DB)

func Open(ctx context.Context, url *url.URL, opts ...OpenOption) (*DB, error) {
	db, err := sql.Open("sqlite", buildDSN(url))
	if err != nil {
		return nil, err
	}

	d := &DB{
		inner: db,
		schema: &sqlstore.Schema{
			Dialect:             "sqlite3",
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
	return d.inner.Close()
}

func (d *DB) Migrate(ctx context.Context) error {
	return d.schema.Migrate(ctx, d.inner)
}

func (d *DB) Rollback(ctx context.Context) error {
	return d.schema.Rollback(ctx, d.inner)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.inner.ExecContext(ctx, query, args...)
}

func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.inner.PrepareContext(ctx, query)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.inner.QueryContext(ctx, query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.inner.QueryRowContext(ctx, query, args...)
}

func (d *DB) RunInTx(ctx context.Context, f func(Tx) error) error {
	tx, err := d.inner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func WithMigrationAttempts(attempts int) OpenOption {
	return func(d *DB) {
		d.schema.MigrationAttempts = attempts
	}
}

func WithMigrationRetryDelay(delay time.Duration) OpenOption {
	return func(d *DB) {
		d.schema.MigrationRetryDelay = delay
	}
}

func buildDSN(u *url.URL) string {
	var base string
	if u.Host == ":memory:" {
		base = ":memory:"
	} else {
		base = "file:" + filepath.Join(u.Host, u.Path)
	}

	query := u.Query()
	query.Del("_pragma")
	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "busy_timeout(5000)")
	query.Add("_pragma", "journal_mode(WAL)")

	rawQuery := query.Encode()
	if rawQuery == "" {
		return base
	}

	return base + "?" + rawQuery
}
