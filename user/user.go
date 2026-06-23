package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type User struct {
	ID          uuid.UUID   `json:"id"`
	OpenID      string      `json:"-"`
	IDP         string      `json:"-"`
	Email       string      `json:"email"`
	DisplayName null.String `json:"displayName"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	DeletedAt   null.Time   `json:"-"`
}

type CreateParams struct {
	OpenID      string
	IDP         string
	Email       string
	DisplayName null.String
}

type Store interface {
	Create(ctx context.Context, params CreateParams) (User, error)
	Get(ctx context.Context, id uuid.UUID) (User, error)
}
