package frontend

import (
	"bytes"
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist
var FS embed.FS

const noncePlaceholder = "NONCEPLACEHOLDER"

type Handler struct {
	fs         fs.FS
	next       http.Handler
	htmlChunks [][]byte
}

func NewHandler() (http.Handler, error) {
	sub, err := fs.Sub(FS, "dist")
	if err != nil {
		return nil, fmt.Errorf("frontend: filesystem error: %w", err)
	}

	html, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		return nil, fmt.Errorf("frontend: index.html read error: %w", err)
	}

	return &Handler{
		fs:         sub,
		next:       http.FileServerFS(sub),
		htmlChunks: bytes.Split(html, []byte(noncePlaceholder)),
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	name = strings.TrimSuffix(name, "/")
	if name == "index.html" {
		http.Redirect(w, r, "/", http.StatusFound)
	} else if _, err := fs.Stat(h.fs, name); err == nil {
		h.next.ServeHTTP(w, r)
	} else {
		h.Index(w, r)
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	nonce := rand.Text()

	w.Header().Set("Content-Security-Policy", CSP(nonce))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for i, chunk := range h.htmlChunks {
		if i > 0 {
			w.Write([]byte(nonce))
		}

		w.Write(chunk)
	}
}

func CSP(nonce string) string {
	return "default-src 'self'; " +
		"script-src 'self' 'nonce-" + nonce + "'; " +
		"style-src 'self' 'nonce-" + nonce + "'; " +
		"img-src 'self' data:; " +
		"connect-src 'self'; " +
		"manifest-src 'self'; " +
		"worker-src 'self'; " +
		"frame-src 'none'; " +
		"object-src 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self';"
}
