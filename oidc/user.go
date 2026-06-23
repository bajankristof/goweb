package oidc

type User struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	DisplayName   string `json:"name"`
	Picture       string `json:"picture"`
	Timezone      string `json:"zoneinfo"`
	Locale        string `json:"locale"`
}
