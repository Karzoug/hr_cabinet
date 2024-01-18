-- +goose Up
-- +goose StatementBegin
BEGIN;

DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'scan_type') THEN
            create type scan_type AS ENUM
                ('Паспорт', 'ИНН', 'СНИЛС', 'Трудовой договор', 'Согласие на обработку данных', 'Военный билет',
                    'Документ об образовании', 'Сертификат', 'Инструктаж', 'Разрешение на работу', 'Свидетельство о браке',
                    'Свидетельство о рождении', 'Другое');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS "scans"
(
    "id"          bigserial PRIMARY KEY,
    "user_id"     bigint    NOT NULL,
    "document_id" bigint    NOT NULL,
    "type"        scan_type NOT NULL,
    "description" varchar,
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz
);

ALTER TABLE "scans"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE OR REPLACE TRIGGER trigger_scans_set_updated_at
    BEFORE UPDATE
    ON scans
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS scans;
DROP TYPE IF EXISTS scan_type;

COMMIT;
-- +goose StatementEnd
