-- +goose Up
-- +goose StatementBegin
BEGIN;

DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'passport_type') THEN
            create type passport_type AS ENUM ('Внутренний', 'Заграничный', 'Иностранного гражданина');
        END IF;
    END
$$;

ALTER TABLE passports
    ADD COLUMN "type" passport_type;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

ALTER TABLE passports
    DROP COLUMN "type";
DROP TYPE IF EXISTS passport_type;

COMMIT;
-- +goose StatementEnd
