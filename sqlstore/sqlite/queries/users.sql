-- name: CreateUser :one
INSERT INTO users (
  id,
  open_id,
  idp,
  email,
  display_name,
  created_at,
  updated_at
)
VALUES (
  @id,
  @open_id,
  @idp,
  @email,
  @display_name,
  strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
  strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
)
ON CONFLICT (open_id, idp) DO
UPDATE SET
  email = EXCLUDED.email,
  display_name = EXCLUDED.display_name,
  updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
  deleted_at = NULL
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? AND deleted_at IS NULL;
