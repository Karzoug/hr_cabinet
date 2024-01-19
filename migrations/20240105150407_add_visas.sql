-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE IF NOT EXISTS "visas"
(
    "id"                   bigserial PRIMARY KEY,
    "user_id"              bigint  NOT NULL,
    "number"               varchar    NOT NULL,
    "issued_state"         varchar,
    "valid_to"             date NOT NULL,
    "valid_from"           date NOT NULL,
    "type"                 varchar NOT NULL,
    "created_at"           timestamptz DEFAULT (now()),
    "updated_at"           timestamptz
);

ALTER TABLE "visas"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

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

COMMIT;
-- +goose StatementEnd
