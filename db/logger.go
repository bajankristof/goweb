package db

import (
	"context"
	"log/slog"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct{}

func (t *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	msg = strings.ToLower(msg)
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		k = strcase.ToSnake(k)
		attrs = append(attrs, slog.Any(k, v))
	}

	var l slog.Level
	switch level {
	case tracelog.LogLevelTrace:
		l = slog.LevelDebug - 1
	case tracelog.LogLevelDebug:
		l = slog.LevelDebug - 1
	case tracelog.LogLevelInfo:
		l = slog.LevelDebug
	case tracelog.LogLevelWarn:
		l = slog.LevelDebug
	case tracelog.LogLevelError:
		l = slog.LevelWarn
	default:
		l = slog.LevelDebug
	}

	slog.LogAttrs(ctx, l, msg, attrs...)
}
