package wellknownhttp

import (
	"net/http"

	"github.com/bajankristof/goweb/wellknown"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	mux  chi.Router
	info wellknown.Info
}

func NewHandler(info wellknown.Info) *Handler {
	h := &Handler{mux: chi.NewRouter(), info: info}

	h.mux.Get("/", h.Info)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, InfoResponse{Info: h.info})
}

type InfoResponse struct {
	Info wellknown.Info `json:"info"`
}

func (i InfoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
