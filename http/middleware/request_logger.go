package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v3"
)

func RequestLogger() func(next http.Handler) http.Handler {
	return httplog.RequestLogger(
		slog.Default(),
		&httplog.Options{
			Level:  slog.LevelDebug,
			Schema: httplog.SchemaOTEL,
		},
	)
}
