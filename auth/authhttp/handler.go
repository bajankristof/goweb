package authhttp

import (
	"crypto/rand"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/bajankristof/goweb/auth"
	"github.com/bajankristof/goweb/http/requestutil"
	"github.com/bajankristof/goweb/oidc"
	"github.com/bajankristof/goweb/session"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	mux chi.Router
	svc *auth.Service
}

func NewHandler(svc *auth.Service) *Handler {
	h := &Handler{mux: chi.NewRouter(), svc: svc}

	h.mux.Get("/signin/{via}", h.SignIn)
	h.mux.Get("/callback", h.Callback)
	h.mux.Post("/refresh", h.Refresh)
	h.mux.Post("/signout", h.SignOut)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	callbackURL, err := requestutil.RewriteURL(r, "/callback", "/signin/{via}")
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to construct callback URL", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	via := chi.URLParam(r, "via")
	nonce := rand.Text()
	authURL, err := h.svc.BeginSignIn(r.Context(), via, nonce, callbackURL.String())
	if errors.Is(err, oidc.ErrInvalidProvider) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.ErrorContext(r.Context(), "failed to begin sign-in", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieOAuthState,
		Value:    nonce,
		Path:     callbackURL.Path,
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Minute * 5),
	})

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	refreshURL, err := requestutil.RewriteURL(r, "/refresh", "/callback")
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to construct refresh URL", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie(CookieOAuthState)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieOAuthState,
		Value:    "",
		Path:     r.URL.Path,
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	userAgent := r.UserAgent()
	user, err := h.svc.CompleteSignIn(r.Context(), c.Value, state, code, userAgent)
	if errors.Is(err, auth.ErrInvalidState) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		slog.ErrorContext(r.Context(), "failed to complete sign-in", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieAccessToken,
		Value:    user.AccessToken,
		Path:     "/",
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  user.Session.ExpiresAt,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieRefreshToken,
		Value:    user.RefreshToken,
		Path:     refreshURL.Path,
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  user.Session.ExpiresAt,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(CookieRefreshToken)
	if c == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u, err := h.svc.Refresh(r.Context(), c.Value, r.UserAgent())
	if errors.Is(err, session.ErrNotFound) ||
		errors.Is(err, session.ErrExpired) ||
		errors.Is(err, session.ErrRevoked) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		slog.ErrorContext(r.Context(), "failed to refresh token", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieAccessToken,
		Value:    u.AccessToken,
		Path:     "/",
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  u.Session.ExpiresAt,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieRefreshToken,
		Value:    u.RefreshToken,
		Path:     r.URL.Path,
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  u.Session.ExpiresAt,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	refreshURL, err := requestutil.RewriteURL(r, "/refresh", "/signout")
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to construct refresh URL", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c, _ := r.Cookie(CookieAccessToken)
	if c == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = h.svc.SignOut(r.Context(), c.Value)
	if err != nil && !errors.Is(err, session.ErrNotFound) {
		slog.ErrorContext(r.Context(), "failed to sign out", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieAccessToken,
		Value:    "",
		Path:     "/",
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieRefreshToken,
		Value:    "",
		Path:     refreshURL.Path,
		Secure:   requestutil.IsSecure(r),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})

	w.WriteHeader(http.StatusNoContent)
}
