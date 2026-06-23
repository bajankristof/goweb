-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  open_id TEXT NOT NULL,
  idp TEXT NOT NULL,
  email TEXT NOT NULL,
  display_name TEXT DEFAULT NULL,
  created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  deleted_at DATETIME DEFAULT NULL
);

CREATE UNIQUE INDEX idx_users_open_id_connect ON users (open_id, idp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
