-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuidv4(),
  open_id VARCHAR(256) NOT NULL,
  idp VARCHAR(128) NOT NULL,
  email VARCHAR(256) NOT NULL,
  display_name VARCHAR(256) DEFAULT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE UNIQUE INDEX idx_users_open_id_connect ON users (open_id, idp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
