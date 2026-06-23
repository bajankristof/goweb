-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT uuidv4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  refresh_token_hash VARCHAR(64) UNIQUE NOT NULL,
  user_agent TEXT NOT NULL,
  refreshed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions CASCADE;
-- +goose StatementEnd
