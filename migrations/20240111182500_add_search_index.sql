-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS users_lastname_idx_gin ON users USING gin (lastname gin_trgm_ops);
CREATE INDEX IF NOT EXISTS departments_title_idx_gin ON departments USING gin (title gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_lastname_idx_gin;
DROP INDEX IF EXISTS departments_title_idx_gin;
-- +goose StatementEnd