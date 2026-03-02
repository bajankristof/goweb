package cookieutil

import (
	"net/http"

	"github.com/bajankristof/watchbowl/http/requestutil"
)

// Get retrieves the value of the cookie with the given name from the request.
// If the cookie is not found, it returns an empty string.
func Get(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// Set sets the given cookie in the response.
// It ensures that the cookie is marked as secure if the request is secure,
// and that it has a default path of "/" and a default SameSite policy of Lax if not specified.
func Set(w http.ResponseWriter, r *http.Request, c *http.Cookie) {
	http.SetCookie(w, expand(r, c))
}

// expand modifies the cookie to ensure it has the appropriate attributes based on the request context.
func expand(r *http.Request, c *http.Cookie) *http.Cookie {
	if requestutil.IsSecure(r) {
		c.Secure = true
	}

	if c.Path == "" {
		c.Path = "/"
	}

	if c.SameSite == http.SameSiteDefaultMode {
		c.SameSite = http.SameSiteLaxMode
	}

	return c
}
