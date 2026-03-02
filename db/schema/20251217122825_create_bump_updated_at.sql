-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION bump_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS bump_updated_at() CASCADE;
-- +goose StatementEnd
