//go:build dev

package main

import (
	"log/slog"
	"os"

	"github.com/bajankristof/goweb/slogcontext"
	"github.com/lmittmann/tint"
)

func init() {
	handler := tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug})
	logger := slog.New(slogcontext.NewHandler(handler))
	slog.SetDefault(logger)
}
