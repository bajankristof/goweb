package dto

import (
	"net/http"

	"github.com/bajankristof/watchbowl/oidc"
)

type AuthWellKnownResponse struct {
	Providers []*oidc.Provider `json:"providers"`
}

func (resp *AuthWellKnownResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
