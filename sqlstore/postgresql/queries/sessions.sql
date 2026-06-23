-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token_hash,
  user_agent,
  expires_at,
  created_at,
  refreshed_at
)
VALUES (
  @id,
  @user_id,
  @refresh_token_hash,
  @user_agent,
  @expires_at,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 AND revoked_at IS NULL;

-- name: GetSessionByRefreshTokenHash :one
SELECT * FROM sessions
WHERE refresh_token_hash = $1 AND revoked_at IS NULL;

-- name: ListSessions :many
SELECT * FROM sessions
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: RefreshSession :one
UPDATE sessions
SET
  refresh_token_hash = @refresh_token_hash,
  user_agent = @user_agent,
  expires_at = @expires_at,
  refreshed_at = CURRENT_TIMESTAMP
WHERE id = @id AND revoked_at IS NULL
RETURNING *;

-- name: RevokeSession :exec
UPDATE sessions
SET
  revoked_at = CURRENT_TIMESTAMP
WHERE id = $1 AND revoked_at IS NULL;
