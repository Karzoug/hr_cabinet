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
        IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'work_type') THEN
            create type work_type AS ENUM ('Удаленная', 'Смешанная', 'Офис');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS "departments"
(
    "id"         bigserial PRIMARY KEY,
    "name"       varchar NOT NULL,
    "created_at" timestamptz DEFAULT (now()),
    "updated_at" timestamptz
);

CREATE TABLE IF NOT EXISTS "positions"
(
    "id"            bigserial PRIMARY KEY,
    "name"          varchar NOT NULL,
    "department_id" bigint  NOT NULL,
    "position_id"   bigint  NOT NULL,
    "created_at"    timestamptz DEFAULT (now()),
    "updated_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "organization_structure"
(
    "id"                        bigserial PRIMARY KEY,
    "head_department_id"        bigint,
    "head_position_id"          bigint,
    "subordinate_department_id" bigint NOT NULL,
    "created_at"                timestamptz DEFAULT (now()),
    "updated_at"                timestamptz
);

CREATE TABLE IF NOT EXISTS "users"
(
    "id"                   bigserial PRIMARY KEY,
    "lastname"             varchar NOT NULL,
    "firstname"            varchar NOT NULL,
    "middlename"           varchar NOT NULL,
    "gender"               gender  NOT NULL,
    "date_of_birth"        date    NOT NULL,
    "place_of_birth"       date    NOT NULL,
    "position_id"          bigint  NOT NULL,
    "department_id"        bigint  NOT NULL,
    "grade"                varchar,
    "phone_numbers"        jsonb,
    "work_email"           varchar NOT NULL,
    "registration_address" varchar NOT NULL,
    "residential_address"  varchar NOT NULL,
    "nationality"          varchar NOT NULL,
    "insurance_number"     varchar NOT NULL,
    "taxpayer_number"      varchar NOT NULL,
    "created_at"           timestamptz DEFAULT (now()),
    "updated_at"           timestamptz
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
    "password_hash" varchar NOT NULL,
    "role"          bigint  NOT NULL,
    "created_at"    timestamptz DEFAULT (now()),
    "updated_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "educations"
(
    "id"                  bigserial PRIMARY KEY,
    "user_id"             bigint  NOT NULL,
    "name_of_institution" varchar NOT NULL,
    "document_number"     varchar NOT NULL,
    "year_of_begin"       date    NOT NULL,
    "year_of_end"         date    NOT NULL,
    "created_at"          timestamptz DEFAULT (now()),
    "updated_at"          timestamptz
);

CREATE TABLE IF NOT EXISTS "trainings"
(
    "id"                  bigserial PRIMARY KEY,
    "user_id"             bigint  NOT NULL,
    "name_of_institution" varchar NOT NULL,
    "name_of_program"     varchar NOT NULL,
    "cost"                integer NOT NULL,
    "date_begin"          date    NOT NULL,
    "date_end"            date    NOT NULL,
    "created_at"          timestamptz DEFAULT (now()),
    "updated_at"          timestamptz
);

CREATE TABLE IF NOT EXISTS "passports"
(
    "id"          bigserial PRIMARY KEY,
    "user_id"     bigint  NOT NULL,
    "series"      integer NOT NULL,
    "number"      integer NOT NULL,
    "issued_date" date    NOT NULL,
    "issued_by"   varchar NOT NULL,
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz
);

CREATE TABLE IF NOT EXISTS "militaries"
(
    "id"                   bigserial PRIMARY KEY,
    "user_id"              bigint  NOT NULL,
    "rank"                 varchar NOT NULL,
    "name_of_commissariat" varchar NOT NULL,
    "specialty"            varchar NOT NULL,
    "category_of_validity" varchar NOT NULL,
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
    "cost"       integer NOT NULL,
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

CREATE TABLE IF NOT EXISTS "contracts"
(
    "id"            bigserial PRIMARY KEY,
    "user_id"       bigint        NOT NULL,
    "number"        integer       NOT NULL,
    "contract_type" contract_type NOT NULL,
    "work_type"     work_type     NOT NULL,
    "date_begin"    date          NOT NULL,
    "date_end"      date,
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

CREATE TABLE IF NOT EXISTS "indexations"
(
    "id"            bigserial PRIMARY KEY,
    "user_id"       bigint  NOT NULL,
    "date_begin"    date    NOT NULL,
    "percents"      integer NOT NULL,
    "currency"      integer NOT NULL,
    "salary_before" integer NOT NULL,
    "created_at"    timestamptz DEFAULT (now()),
    "updated_at"    timestamptz
);

CREATE TABLE IF NOT EXISTS "finances"
(
    "id"                  bigserial PRIMARY KEY,
    "user_id"             bigint  NOT NULL,
    "salary"              integer NOT NULL,
    "fee_to_pension_fund" integer NOT NULL,
    "fee_to_tax"          integer NOT NULL,
    "created_at"          timestamptz DEFAULT (now()),
    "updated_at"          timestamptz
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

ALTER TABLE "positions"
    ADD FOREIGN KEY ("position_id") REFERENCES "positions" ("id");

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
    ADD FOREIGN KEY ("role") REFERENCES "roles" ("id");

ALTER TABLE "educations"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "trainings"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "passports"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "militaries"
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

ALTER TABLE "experiences"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "indexations"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "finances"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
