package middleware

import (
	"context"
	"net/http"

	"github.com/bajankristof/watchbowl/http/cookieutil"
	"github.com/bajankristof/watchbowl/jwt"
	"github.com/google/uuid"
)

type authContextKey struct{}

var authKey = authContextKey{}

func Auth(jwts *jwt.Signer, errorHandler http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := cookieutil.Get(r, "access_token")
			if accessToken == "" {
				if errorHandler != nil {
					errorHandler.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r)
				}
				return
			}

			claims := &jwt.RegisteredClaims{}
			err := jwts.Verify(accessToken, claims)
			if err != nil {
				if errorHandler != nil {
					errorHandler.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r)
				}
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				if errorHandler != nil {
					errorHandler.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r)
				}
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, authKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetCurrentUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(authKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}
