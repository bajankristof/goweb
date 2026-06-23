//go:build !dev

package main

import (
	"log/slog"
	"os"

	"github.com/bajankristof/goweb/slogcontext"
)

func init() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(slogcontext.NewHandler(handler))
	slog.SetDefault(logger)
}
