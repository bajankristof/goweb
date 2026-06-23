package authhttp

import (
	"net/http"

	"github.com/bajankristof/goweb/auth"
	"github.com/bajankristof/goweb/slogcontext"
)

func Authn(s *auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(CookieAccessToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := s.Verify(r.Context(), c.Value)
			// TODO: log the error (?)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := auth.WithAuthUserID(r.Context(), userID)
			ctx = slogcontext.With(ctx, "user.id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
