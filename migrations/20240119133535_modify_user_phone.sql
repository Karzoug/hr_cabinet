-- +goose Up
-- +goose StatementBegin
BEGIN;
ALTER TABLE users
    ADD COLUMN "mobile_phone_number" varchar,
    ADD COLUMN "office_phone_number" varchar NOT NULL DEFAULT ''; 

UPDATE users SET mobile_phone_number = phone_numbers ->> 'mobile';

ALTER TABLE users
    DROP COLUMN "phone_numbers",
    ALTER COLUMN "mobile_phone_number" SET NOT NULL;
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd