package requestutil

import (
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"path"
	"strings"

	"github.com/bajankristof/trustedproxy"
	"github.com/go-chi/chi/v5"
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

// NetIPAddr extracts the remote IP address from the request's RemoteAddr field.
func NetIPAddr(r *http.Request) (netip.Addr, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return netip.Addr{}, err
	}

	return netip.ParseAddr(host)
}

// RewriteURL checks if the request's route pattern ends with any of the given suffixes,
// and if so, rewrites the URL to the given path and returns the new URL and true.
// If no suffix matches, it returns nil and false.
func RewriteURL(r *http.Request, to string, from ...string) (*url.URL, bool) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return nil, false
	}

	pattern := rctx.RoutePattern()
	for _, f := range from {
		if prefix, ok := strings.CutSuffix(pattern, f); ok {
			return AbsoluteURL(r, path.Join(prefix, to)), true
		}
	}

	return nil, false
}
