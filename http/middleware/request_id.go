package middleware

import (
	"context"
	"net/http"

	"github.com/bajankristof/goweb/slogcontext"
	"github.com/go-chi/chi/v5/middleware"
)

func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if i := GetRequestID(ctx); i != "" {
			r = slogcontext.Inject(r, "http.request.id", i)
		}
		next.ServeHTTP(w, r)
	}))
}

func GetRequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}
