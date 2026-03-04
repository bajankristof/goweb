package handler

import (
	"net/http"
	"net/netip"

	"github.com/bajankristof/goweb/db"
	"github.com/bajankristof/goweb/frontend"
	"github.com/bajankristof/goweb/http/dto"
	"github.com/bajankristof/goweb/http/middleware"
	"github.com/bajankristof/goweb/jwt"
	"github.com/bajankristof/goweb/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var jsonAuthErrorHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, &dto.ErrResponse{
		Status:  http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
		Code:    "unauthorized",
	})
}

func New(
	dbq *db.Queries,
	jwts *jwt.Signer,
	oidr *oidc.Registry,
	allowedOrigins []string,
	trustedProxies []netip.Prefix,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.TrustedProxy(trustedProxies))
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS(allowedOrigins))
	r.Use(middleware.CrossOriginProtection(allowedOrigins))
	r.Use(middleware.Recoverer)

	r.Get("/readyz", readyzHandler(dbq))
	r.Get("/healthz", healthzHandler())
	r.Mount("/auth", authHandler(dbq, jwts, oidr))

	r.Group(func(g chi.Router) {
		g.Use(middleware.Auth(jwts, jsonAuthErrorHandler))
		g.Mount("/api/v1/users", usersHandler(dbq))
	})

	r.Mount("/", frontend.FileServer())

	return r
}
