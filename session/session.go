package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Session struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"-"`
	RefreshTokenHash string    `json:"-"`
	UserAgent        string    `json:"userAgent"`
	ExpiresAt        time.Time `json:"expiresAt"`
	CreatedAt        time.Time `json:"createdAt"`
	RefreshedAt      time.Time `json:"refreshedAt"`
	RevokedAt        null.Time `json:"-"`
}

type CreateParams struct {
	UserID           uuid.UUID
	UserAgent        string
	RefreshTokenHash string
	ExpiresAt        time.Time
}

type RefreshParams struct {
	ID               uuid.UUID
	UserAgent        string
	RefreshTokenHash string
	ExpiresAt        time.Time
}

type Store interface {
	Create(ctx context.Context, params CreateParams) (Session, error)
	Get(ctx context.Context, id uuid.UUID) (Session, error)
	GetByRefreshTokenHash(ctx context.Context, hash string) (Session, error)
	List(ctx context.Context, userID uuid.UUID) ([]Session, error)
	Refresh(ctx context.Context, params RefreshParams) (Session, error)
	Revoke(ctx context.Context, id uuid.UUID) error
}
