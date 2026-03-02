package db

import (
	"context"
	"net/netip"

	"github.com/guregu/null/v6"
)

type CreateWebUserParams struct {
	OpenID           string
	Provider         string
	Email            string
	DisplayName      null.String
	RefreshTokenHash string
	IpAddress        netip.Addr
	UserAgent        string
}

func (q *Queries) CreateWebUser(ctx context.Context, arg CreateWebUserParams) (User, error) {
	var user User
	return user, q.tx(ctx, func(qtx *Queries) error {
		var err error
		user, err = qtx.CreateUser(ctx, CreateUserParams{
			OpenID:      arg.OpenID,
			Provider:    arg.Provider,
			Email:       arg.Email,
			DisplayName: arg.DisplayName,
		})
		if err != nil {
			return err
		}

		_, err = qtx.CreateSession(ctx, CreateSessionParams{
			UserID:           user.UserID,
			RefreshTokenHash: arg.RefreshTokenHash,
			IpAddress:        arg.IpAddress,
			UserAgent:        arg.UserAgent,
		})
		return err
	})
}
