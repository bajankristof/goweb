package frontend

import (
	"embed"
	"net/http"
)

//go:embed dist/*
var FS embed.FS

func FileServer() http.Handler {
	return http.FileServer(http.FS(FS))
}
