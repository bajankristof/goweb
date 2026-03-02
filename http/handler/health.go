package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type readyzQueries interface {
	Ping(ctx context.Context) (string, error)
}

func healthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, render.M{})
	}
}

func readyzHandler(dbq readyzQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := dbq.Ping(r.Context())
		if err != nil {
			slog.WarnContext(r.Context(), "failed to ping database", "err", err)
			render.Status(r, http.StatusServiceUnavailable)
			render.JSON(w, r, render.M{})
			return
		}

		render.JSON(w, r, render.M{})
	}
}
