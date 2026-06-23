package user

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	return s.store.Get(ctx, id)
}
