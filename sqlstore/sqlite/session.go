package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bajankristof/goweb/session"
	"github.com/google/uuid"
)

var _ session.Store = &SessionStore{}

type SessionStore struct {
	queries *Queries
}

func NewSessionStore(dbtx DBTX) *SessionStore {
	return &SessionStore{queries: New(dbtx)}
}

func (st *SessionStore) Create(ctx context.Context, params session.CreateParams) (session.Session, error) {
	ds, err := st.queries.CreateSession(ctx, CreateSessionParams{
		ID:               uuid.New(),
		UserID:           params.UserID,
		RefreshTokenHash: params.RefreshTokenHash,
		UserAgent:        params.UserAgent,
		ExpiresAt:        params.ExpiresAt,
	})
	if err != nil {
		return session.Session{}, fmt.Errorf("session: create error: %w", err)
	}

	return newSessionFromDB(ds), nil
}

func (st *SessionStore) Get(ctx context.Context, id uuid.UUID) (session.Session, error) {
	ds, err := st.queries.GetSession(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return session.Session{}, session.ErrNotFound
	} else if err != nil {
		return session.Session{}, fmt.Errorf("session: get error: %w", err)
	}

	return newSessionFromDB(ds), nil
}

func (st *SessionStore) GetByRefreshTokenHash(ctx context.Context, hash string) (session.Session, error) {
	ds, err := st.queries.GetSessionByRefreshTokenHash(ctx, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return session.Session{}, session.ErrNotFound
	} else if err != nil {
		return session.Session{}, fmt.Errorf("session: get by refresh token hash error: %w", err)
	}

	return newSessionFromDB(ds), nil
}

func (st *SessionStore) List(ctx context.Context, userID uuid.UUID) ([]session.Session, error) {
	dslc, err := st.queries.ListSessions(ctx, userID)
	if err != nil {
		return []session.Session{}, fmt.Errorf("session: list error: %w", err)
	}

	slc := make([]session.Session, len(dslc))
	for i, ds := range dslc {
		slc[i] = newSessionFromDB(ds)
	}

	return slc, nil
}

func (st *SessionStore) Refresh(ctx context.Context, params session.RefreshParams) (session.Session, error) {
	ds, err := st.queries.RefreshSession(ctx, RefreshSessionParams{
		ID:               params.ID,
		RefreshTokenHash: params.RefreshTokenHash,
		UserAgent:        params.UserAgent,
		ExpiresAt:        params.ExpiresAt,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return session.Session{}, session.ErrNotFound
	} else if err != nil {
		return session.Session{}, fmt.Errorf("session: refresh error: %w", err)
	}

	return newSessionFromDB(ds), nil
}

func (st *SessionStore) Revoke(ctx context.Context, id uuid.UUID) error {
	err := st.queries.RevokeSession(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return session.ErrNotFound
	} else if err != nil {
		return fmt.Errorf("session: revoke error: %w", err)
	}

	return nil
}

func newSessionFromDB(ds Session) session.Session {
	return session.Session{
		ID:               ds.ID,
		UserID:           ds.UserID,
		RefreshTokenHash: ds.RefreshTokenHash,
		UserAgent:        ds.UserAgent,
		ExpiresAt:        ds.ExpiresAt,
		CreatedAt:        ds.CreatedAt,
		RefreshedAt:      ds.RefreshedAt,
		RevokedAt:        ds.RevokedAt,
	}
}
