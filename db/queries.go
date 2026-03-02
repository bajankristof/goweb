package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) tx(ctx context.Context, fn func(*Queries) error) error {
	if c, ok := q.db.(*DB); ok {
		return c.WithTx(ctx, func(tx pgx.Tx) error {
			qtx := q.WithTx(tx)
			return fn(qtx)
		})
	}

	return fn(q)
}
