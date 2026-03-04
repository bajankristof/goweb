//go:build !dev

package main

import (
	"log/slog"
	"os"

	"github.com/bajankristof/goweb/slogctx"
)

func init() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(slogctx.NewHandler(handler))
	slog.SetDefault(logger)
}
