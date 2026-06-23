package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/bajankristof/goweb/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	queries *Queries
}

func NewUserStore(dbtx DBTX) *UserStore {
	return &UserStore{queries: New(dbtx)}
}

func (s *UserStore) Create(ctx context.Context, params user.CreateParams) (user.User, error) {
	du, err := s.queries.CreateUser(ctx, CreateUserParams{
		OpenID:      params.OpenID,
		IDP:         params.IDP,
		Email:       params.Email,
		DisplayName: params.DisplayName,
	})
	if err != nil {
		return user.User{}, fmt.Errorf("user: create error: %w", err)
	}

	return newUserFromDB(du), nil
}

func (s *UserStore) Get(ctx context.Context, id uuid.UUID) (user.User, error) {
	du, err := s.queries.GetUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, user.ErrNotFound
	} else if err != nil {
		return user.User{}, fmt.Errorf("user: get error: %w", err)
	}

	return newUserFromDB(du), nil
}

func newUserFromDB(du User) user.User {
	return user.User{
		ID:          du.ID,
		OpenID:      du.OpenID,
		IDP:         du.IDP,
		Email:       du.Email,
		DisplayName: du.DisplayName,
		CreatedAt:   du.CreatedAt,
		UpdatedAt:   du.UpdatedAt,
		DeletedAt:   du.DeletedAt,
	}
}
