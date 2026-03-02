package db

import "github.com/jackc/pgx/v5"

var (
	ErrNoRows           = pgx.ErrNoRows
	ErrTooManyRows      = pgx.ErrTooManyRows
	ErrTxClosed         = pgx.ErrTxClosed
	ErrTxCommitRollback = pgx.ErrTxCommitRollback
)
