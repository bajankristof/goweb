package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist
var FS embed.FS

type Handler struct {
	fs   fs.FS
	next http.Handler
}

func NewHandler() http.Handler {
	sub, _ := fs.Sub(FS, "dist")
	next := http.FileServerFS(sub)

	return &Handler{
		fs:   sub,
		next: next,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	name = strings.TrimSuffix(name, "/")
	if _, err := fs.Stat(h.fs, name); err == nil {
		h.next.ServeHTTP(w, r)
	} else {
		h.Index(w, r)
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	r = r.Clone(r.Context())
	r.URL.Path = "/"
	h.next.ServeHTTP(w, r)
}
