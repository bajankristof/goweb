-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
  session_id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  refresh_token_hash VARCHAR(64) UNIQUE NOT NULL,
  ip_address INET NOT NULL,
  user_agent TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  revoked_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_refresh_token_hash ON sessions (refresh_token_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions CASCADE;
-- +goose StatementEnd
