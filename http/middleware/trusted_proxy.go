package middleware

import (
	"net/http"
	"net/netip"

	"github.com/bajankristof/trustedproxy"
)

func TrustedProxy(prefixes []netip.Prefix) func(next http.Handler) http.Handler {
	tp := &trustedproxy.TrustedProxy{}
	for _, p := range prefixes {
		tp.AddPrefix(p)
	}

	return tp.Handler
}
