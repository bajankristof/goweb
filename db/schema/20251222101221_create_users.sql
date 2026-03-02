-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  user_id UUID PRIMARY KEY DEFAULT uuidv7(),
  open_id VARCHAR(256) NOT NULL,
  provider VARCHAR(128) NOT NULL,
  email VARCHAR(256) NOT NULL,
  display_name VARCHAR(256) DEFAULT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE UNIQUE INDEX idx_users_open_id_connect ON users (open_id, provider);

CREATE TRIGGER set_updated_at_on_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE bump_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
