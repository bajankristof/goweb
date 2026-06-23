-- name: CreateUser :one
INSERT INTO users (
  open_id,
  idp,
  email,
  display_name,
  created_at,
  updated_at
)
VALUES (
  @open_id,
  @idp,
  @email,
  @display_name,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
)
ON CONFLICT (open_id, idp) DO
UPDATE SET
  email = EXCLUDED.email,
  display_name = EXCLUDED.display_name,
  updated_at = CURRENT_TIMESTAMP,
  deleted_at = NULL
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;
