-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token_hash, ip_address, user_agent)
VALUES (@user_id, @refresh_token_hash, @ip_address, @user_agent)
RETURNING *;

-- name: RotateSession :one
UPDATE sessions
SET refresh_token_hash = @new_refresh_token_hash, last_used_at = CURRENT_TIMESTAMP
WHERE refresh_token_hash = @refresh_token_hash
  AND revoked_at IS NULL
  AND EXISTS (SELECT 1 FROM users WHERE users.user_id = sessions.user_id AND users.deleted_at IS NULL)
RETURNING *;

-- name: RevokeSession :one
UPDATE sessions
SET revoked_at = CURRENT_TIMESTAMP
WHERE refresh_token_hash = $1 AND revoked_at IS NULL
RETURNING *;
