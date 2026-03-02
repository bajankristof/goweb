-- name: CreateUser :one
INSERT INTO users (open_id, provider, email, display_name)
VALUES (@open_id, @provider, @email, @display_name)
ON CONFLICT (open_id, provider) DO
UPDATE SET email = EXCLUDED.email, display_name = EXCLUDED.display_name
WHERE users.deleted_at IS NULL
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1 AND deleted_at IS NULL;
