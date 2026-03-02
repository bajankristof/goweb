package schema

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.SetBaseFS(FS)
	goose.SetLogger(logger{})
	goose.SetTableName("schema_migrations")
	err := goose.SetDialect("postgres")
	if err != nil {
		panic(err)
	}
}

func Migrate(ctx context.Context, db *sql.DB) error {
	return goose.UpContext(ctx, db, ".")
}

func Rollback(ctx context.Context, db *sql.DB) error {
	return goose.DownContext(ctx, db, ".")
}
