-- +goose Up
-- +goose StatementBegin
BEGIN;

DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'visa_number_entries') THEN
            create type visa_number_entries AS ENUM ('1', '2', 'mult');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS "visas"
(
    "id"                   bigserial PRIMARY KEY,
    "user_id"              bigint  NOT NULL,
    "passport_id"          bigint  NOT NULL,
    "number"               varchar    NOT NULL,
    "issued_state"         varchar    NOT NULL,
    "valid_to"             date NOT NULL,
    "valid_from"           date NOT NULL,
    "number_entries" visa_number_entries NOT NULL,
    "created_at"           timestamptz DEFAULT (now()),
    "updated_at"           timestamptz
);

ALTER TABLE "visas"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    ADD FOREIGN KEY ("passport_id") REFERENCES "passports" ("id");

CREATE OR REPLACE TRIGGER trigger_visas_set_updated_at
    BEFORE UPDATE
    ON visas
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS visas;
DROP TYPE IF EXISTS visa_number_entries;

COMMIT;
-- +goose StatementEnd
