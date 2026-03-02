package schema

import (
	"fmt"
	"log/slog"
	"strings"
)

type Logger interface {
	Printf(format string, v ...any)
	Fatalf(format string, v ...any)
}

type logger struct{}

func (l logger) Printf(format string, v ...any) {
	slog.Info(l.format(format, v...))
}

func (l logger) Fatalf(format string, v ...any) {
	slog.Error(l.format(format, v...))
}

func (l logger) format(format string, v ...any) string {
	if len(format) >= 7 && format[:7] == "goose: " {
		format = format[7:]
	}

	format = strings.ReplaceAll(format, ".", ":")

	return fmt.Sprintf(format, v...)
}
