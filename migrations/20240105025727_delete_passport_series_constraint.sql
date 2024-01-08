-- +goose Up
-- +goose StatementBegin
ALTER TABLE passports
    ALTER COLUMN series DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
