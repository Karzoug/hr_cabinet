-- +goose Up
-- +goose StatementBegin
BEGIN;

DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'gender') THEN
            create type gender AS ENUM ('Мужской', 'Женский');
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'contract_type') THEN
            create type contract_type AS ENUM ('Срочный', 'Бессрочный');
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'passport_type') THEN
            create type passport_type AS ENUM ('Внутренний', 'Заграничный');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS "departments"
(
    "id"          bigserial PRIMARY KEY,
    "title"       varchar NOT NULL,
    "description" varchar,
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz
);

CREATE TABLE IF NOT EXISTS "positions"
(
    "id"            bigserial PRIMARY KEY,
    "title"         varchar NOT NULL,
    "description"   varchar,
    "department_id" bigint,
    "created_at"    timestamptz DEFAULT (now()),
    "updated_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "organization_structure"
(
    "id"                        bigserial PRIMARY KEY,
    "head_department_id"        bigint,
    "head_position_id"          bigint,
    "subordinate_department_id" bigint,
    "created_at"                timestamptz DEFAULT (now()),
    "updated_at"                timestamptz
);

CREATE TABLE IF NOT EXISTS "users"
(
    "id"               bigserial PRIMARY KEY,
    "lastname"         varchar NOT NULL,
    "firstname"        varchar NOT NULL,
    "middlename"       varchar NOT NULL,
    "gender"           gender  NOT NULL,
    "date_of_birth"    date    NOT NULL,
    "place_of_birth"   varchar NOT NULL,
    "position_id"      bigint  NOT NULL,
    "department_id"    bigint  NOT NULL,
    "grade"            varchar NOT NULL,
    "phone_numbers"    jsonb   NOT NULL,
    "work_email"       varchar NOT NULL,
    "registration_address" varchar NOT NULL,
    "residential_address"  varchar NOT NULL,
    "insurance_number" varchar NOT NULL,
    "taxpayer_number"  varchar NOT NULL,
    "created_at"       timestamptz DEFAULT (now()),
    "updated_at"       timestamptz,
    UNIQUE (work_email)
);

CREATE TABLE IF NOT EXISTS "roles"
(
    "id"          bigserial PRIMARY KEY,
    "title"       varchar NOT NULL,
    "description" varchar,
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz
);

CREATE TABLE IF NOT EXISTS "authorizations"
(
    "id"            bigserial PRIMARY KEY,
    "user_id"       bigint  NOT NULL,
    "role_id"       bigint  NOT NULL,
    "password_hash" varchar NOT NULL,
    "created_at"    timestamptz DEFAULT (now()),
    "updated_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "passports"
(
    "id"                   bigserial PRIMARY KEY,
    "user_id"              bigint  NOT NULL,
    "number"               varchar NOT NULL,
    "citizenship"         varchar NOT NULL,
    "type"                 passport_type,
    "issued_date"          date    NOT NULL,
    "issued_by"            varchar,
    "issued_by_code"       varchar,
    "created_at"           timestamptz DEFAULT (now()),
    "updated_at"           timestamptz
);

CREATE TABLE IF NOT EXISTS "militaries"
(
    "id"                    bigserial PRIMARY KEY,
    "user_id"               bigint  NOT NULL,
    "rank"                  varchar NOT NULL,
    "title_of_commissariat" varchar NOT NULL,
    "specialty"             varchar NOT NULL,
    "category_of_validity"  varchar NOT NULL,
    "created_at"            timestamptz DEFAULT (now()),
    "updated_at"            timestamptz
);

CREATE TABLE IF NOT EXISTS "educations"
(
    "id"                   bigserial PRIMARY KEY,
    "user_id"              bigint  NOT NULL,
    "title_of_institution" varchar NOT NULL,
    "title_of_program"     varchar NOT NULL,
    "document_number"      varchar NOT NULL,
    "year_of_begin"        date    NOT NULL,
    "year_of_end"          date    NOT NULL,
    "created_at"           timestamptz DEFAULT (now()),
    "updated_at"           timestamptz
);

CREATE TABLE IF NOT EXISTS "trainings"
(
    "id"                   bigserial PRIMARY KEY,
    "user_id"              bigint  NOT NULL,
    "title_of_institution" varchar NOT NULL,
    "title_of_program"     varchar NOT NULL,
    "cost"                 bigint  NOT NULL,
    "date_begin"           date    NOT NULL,
    "date_end"             date    NOT NULL,
    "created_at"           timestamptz DEFAULT (now()),
    "updated_at"           timestamptz
);

CREATE TABLE IF NOT EXISTS "vacations"
(
    "id"         bigserial PRIMARY KEY,
    "user_id"    bigint NOT NULL,
    "date_begin" date   NOT NULL,
    "date_end"   date   NOT NULL,
    "created_at" timestamptz DEFAULT (now()),
    "updated_at" timestamptz
);

CREATE TABLE IF NOT EXISTS "benefits"
(
    "id"         bigserial PRIMARY KEY,
    "title"      varchar NOT NULL,
    "cost"       bigint  NOT NULL,
    "created_at" timestamptz DEFAULT (now()),
    "updated_at" timestamptz
);

CREATE TABLE IF NOT EXISTS "benefit_uses"
(
    "id"         bigserial PRIMARY KEY,
    "user_id"    bigint NOT NULL,
    "benefit_id" bigint NOT NULL,
    "date_begin" date   NOT NULL,
    "date_end"   date   NOT NULL,
    "created_at" timestamptz DEFAULT (now()),
    "updated_at" timestamptz,
    UNIQUE (benefit_id)
);

CREATE TABLE IF NOT EXISTS "work_types"
(
    "id"          bigserial PRIMARY KEY,
    "title"       varchar NOT NULL,
    "description" varchar,
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz
);

CREATE TABLE IF NOT EXISTS "contracts"
(
    "id"               bigserial PRIMARY KEY,
    "user_id"          bigint        NOT NULL,
    "number"           varchar       NOT NULL,
    "contract_type"    contract_type NOT NULL,
    "work_type_id"     bigint        NOT NULL,
    "probation_period" integer,
    "date_begin"       date          NOT NULL,
    "date_end"         date,
    "created_at"       timestamptz DEFAULT (now()),
    "updated_at"       timestamptz
);

CREATE TABLE IF NOT EXISTS "finances"
(
    "id"                  bigserial PRIMARY KEY,
    "user_id"             bigint NOT NULL,
    "contract_id"         bigint NOT NULL,
    "salary"              bigint NOT NULL,
    "salary_rate"         double precision,
    "social_security_tax" bigint NOT NULL,
    "income_tax"          bigint NOT NULL,
    "created_at"          timestamptz DEFAULT (now()),
    "updated_at"          timestamptz
);

CREATE TABLE IF NOT EXISTS "indexations"
(
    "id"            bigserial PRIMARY KEY,
    "user_id"       bigint  NOT NULL,
    "date_begin"    date    NOT NULL,
    "percents"      integer NOT NULL,
    "currency"      bigint  NOT NULL,
    "salary_before" bigint  NOT NULL,
    "created_at"    timestamptz DEFAULT (now()),
    "updated_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "experiences"
(
    "id"           bigserial PRIMARY KEY,
    "user_id"      bigint  NOT NULL,
    "company_name" varchar NOT NULL,
    "date_begin"   date    NOT NULL,
    "date_end"     date    NOT NULL,
    "position"     varchar NOT NULL,
    "functional"   varchar NOT NULL,
    "awards"       varchar NOT NULL,
    "created_at"   timestamptz DEFAULT (now()),
    "updated_at"   timestamptz
);

CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS
$update_updated_at$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$update_updated_at$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_departments_set_updated_at
    BEFORE UPDATE
    ON departments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_positions_set_updated_at
    BEFORE UPDATE
    ON positions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_organization_structure_set_updated_at
    BEFORE UPDATE
    ON organization_structure
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_users_set_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_roles_set_updated_at
    BEFORE UPDATE
    ON roles
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_authorizations_set_updated_at
    BEFORE UPDATE
    ON authorizations
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_educations_set_updated_at
    BEFORE UPDATE
    ON educations
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_trainings_set_updated_at
    BEFORE UPDATE
    ON trainings
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_passports_set_updated_at
    BEFORE UPDATE
    ON passports
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_militaries_set_updated_at
    BEFORE UPDATE
    ON militaries
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_vacations_set_updated_at
    BEFORE UPDATE
    ON vacations
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_benefits_set_updated_at
    BEFORE UPDATE
    ON benefits
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_benefit_uses_set_updated_at
    BEFORE UPDATE
    ON benefit_uses
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_contracts_set_updated_at
    BEFORE UPDATE
    ON contracts
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_experiences_set_updated_at
    BEFORE UPDATE
    ON experiences
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_indexations_set_updated_at
    BEFORE UPDATE
    ON indexations
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE TRIGGER trigger_finances_set_updated_at
    BEFORE UPDATE
    ON finances
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

ALTER TABLE "positions"
    ADD FOREIGN KEY ("department_id") REFERENCES "departments" ("id");

ALTER TABLE "organization_structure"
    ADD FOREIGN KEY ("head_department_id") REFERENCES "departments" ("id");

ALTER TABLE "organization_structure"
    ADD FOREIGN KEY ("head_position_id") REFERENCES "positions" ("id");

ALTER TABLE "organization_structure"
    ADD FOREIGN KEY ("subordinate_department_id") REFERENCES "departments" ("id");

ALTER TABLE "users"
    ADD FOREIGN KEY ("position_id") REFERENCES "positions" ("id");

ALTER TABLE "users"
    ADD FOREIGN KEY ("department_id") REFERENCES "departments" ("id");

ALTER TABLE "authorizations"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "authorizations"
    ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "passports"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "militaries"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "educations"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "trainings"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "vacations"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "benefit_uses"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE TABLE "benefits_benefit_uses"
(
    "benefits_id"             bigint,
    "benefit_uses_benefit_id" bigint,
    PRIMARY KEY ("benefits_id", "benefit_uses_benefit_id")
);

ALTER TABLE "benefits_benefit_uses"
    ADD FOREIGN KEY ("benefits_id") REFERENCES "benefits" ("id");

ALTER TABLE "benefits_benefit_uses"
    ADD FOREIGN KEY ("benefit_uses_benefit_id") REFERENCES "benefit_uses" ("benefit_id");

ALTER TABLE "contracts"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "contracts"
    ADD FOREIGN KEY ("work_type_id") REFERENCES "work_types" ("id");

ALTER TABLE "finances"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "finances"
    ADD FOREIGN KEY ("contract_id") REFERENCES "contracts" ("id");

ALTER TABLE "indexations"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "experiences"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS organization_structure CASCADE;
DROP TABLE IF EXISTS positions CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP TABLE IF EXISTS benefits CASCADE;
DROP TABLE IF EXISTS benefit_uses CASCADE;
DROP TABLE IF EXISTS benefits_benefit_uses;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS authorizations;
DROP TABLE IF EXISTS educations;
DROP TABLE IF EXISTS trainings;
DROP TABLE IF EXISTS passports;
DROP TABLE IF EXISTS militaries;
DROP TABLE IF EXISTS vacations;
DROP TABLE IF EXISTS finances CASCADE;
DROP TABLE IF EXISTS work_types CASCADE;
DROP TABLE IF EXISTS contracts;
DROP TABLE IF EXISTS experiences;
DROP TABLE IF EXISTS indexations;

DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS contract_type;
DROP TYPE IF EXISTS passport_type;

DROP FUNCTION IF EXISTS update_updated_at();

COMMIT;
-- +goose StatementEnd
