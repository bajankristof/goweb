package userhttp

import (
	"errors"
	"net/http"

	"github.com/bajankristof/goweb/auth"
	"github.com/bajankristof/goweb/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	mux chi.Router
	svc *user.Service
}

func NewHandler(svc *user.Service) *Handler {
	h := &Handler{mux: chi.NewRouter(), svc: svc}

	h.mux.Get("/me", h.GetCurrentUser)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetAuthUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u, err := h.svc.GetUser(r.Context(), userID)
	if errors.Is(err, user.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Render(w, r, GetCurrentUserResponse{User: u})
}

type GetCurrentUserResponse struct {
	User user.User `json:"user"`
}

func (t GetCurrentUserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
