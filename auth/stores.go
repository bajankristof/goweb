package auth

import (
	"context"

	"github.com/bajankristof/goweb/session"
	"github.com/bajankristof/goweb/user"
)

type Stores struct {
	UserStore    user.Store
	SessionStore session.Store
}

type UnitOfWork interface {
	Do(ctx context.Context, f func(s Stores) error) error
}
