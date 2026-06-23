package requestutil

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/bajankristof/trustedproxy"
	"github.com/go-chi/chi/v5"
)

var (
	ErrNoRouteContext = errors.New("no route context in request")
	ErrNoRouteMatch   = errors.New("no route pattern matches the request")
)

// AbsoluteURL returns a new URL based on the request's scheme and host, with the given path.
func AbsoluteURL(r *http.Request, path string) *url.URL {
	scheme := "http"
	if IsSecure(r) {
		scheme = "https"
	}
	return &url.URL{Scheme: scheme, Host: r.Host, Path: path}
}

// IsSecure checks if the request is made over a secure connection (HTTPS).
func IsSecure(r *http.Request) bool {
	return trustedproxy.IsSecure(r)
}

// RewriteURL checks if the request's route pattern ends with any of the given suffixes,
// and if so, rewrites the URL to the given path and returns the new URL and true.
// If no suffix matches, it returns nil and false.
func RewriteURL(r *http.Request, to string, from ...string) (*url.URL, error) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return nil, ErrNoRouteContext
	}

	pattern := rctx.RoutePattern()
	for _, f := range from {
		if prefix, ok := strings.CutSuffix(pattern, f); ok {
			return AbsoluteURL(r, path.Join(prefix, to)), nil
		}
	}

	return nil, fmt.Errorf("%w: %s does not end with any of %v", ErrNoRouteMatch, pattern, from)
}
