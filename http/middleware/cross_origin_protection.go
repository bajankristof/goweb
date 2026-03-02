package middleware

import (
	"log/slog"
	"net/http"
)

func CrossOriginProtection(origins []string) func(next http.Handler) http.Handler {
	cop := &http.CrossOriginProtection{}

	for _, o := range origins {
		err := cop.AddTrustedOrigin(o)
		if err != nil {
			slog.Warn("failed to add trusted origin to cross-origin protection", "origin", o, "err", err)
		}
	}

	return cop.Handler
}
