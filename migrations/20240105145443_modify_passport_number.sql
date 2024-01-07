-- +goose Up
-- +goose StatementBegin
ALTER TABLE passports
    ALTER COLUMN "number" TYPE varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
