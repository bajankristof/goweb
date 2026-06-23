package auth

import (
	"github.com/bajankristof/goweb/session"
	"github.com/bajankristof/goweb/user"
)

type User struct {
	user.User
	Session      session.Session `json:"session"`
	AccessToken  string          `json:"-"`
	RefreshToken string          `json:"-"`
}
