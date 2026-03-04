package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bajankristof/goweb/db"
	"github.com/bajankristof/goweb/http/dto"
	"github.com/bajankristof/goweb/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type usersYouQueries interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (db.User, error)
}

func usersHandler(dbq *db.Queries) http.Handler {
	r := chi.NewRouter()
	r.Get("/u", usersYouHandler(dbq))
	return r
}

func usersYouHandler(dbq usersYouQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetCurrentUserID(r.Context())
		user, err := dbq.GetUserByID(r.Context(), userID)
		if errors.Is(err, db.ErrNoRows) {
			render.Render(w, r, &dto.ErrResponse{
				Err:    err,
				Status: http.StatusNotFound,
			})
			return
		} else if err != nil {
			slog.ErrorContext(r.Context(), "failed to get user by ID", "user_id", userID, "err", err)
			render.Render(w, r, &dto.ErrResponse{
				Err:    err,
				Status: http.StatusInternalServerError,
			})
			return
		}

		render.JSON(w, r, user)
	}
}
