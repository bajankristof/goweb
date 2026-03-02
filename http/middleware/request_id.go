package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/bajankristof/watchbowl/slogctx"
	"github.com/go-chi/chi/v5/middleware"
)

func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if id := GetRequestID(ctx); id != "" {
			ctx = slogctx.WithAttrs(ctx, slog.String("context.trace_id", id))
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	}))
}

func GetRequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}
