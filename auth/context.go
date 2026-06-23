package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey struct{}

var (
	contextKeyAuthUserID = contextKey{}
)

func WithAuthUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyAuthUserID, userID)
}

func GetAuthUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(contextKeyAuthUserID).(uuid.UUID)
	return id, ok
}
