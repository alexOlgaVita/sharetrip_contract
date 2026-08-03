-- +goose Up
-- +goose StatementBegin

-- Область хранения: договоры (см. README.md → «Схема хранения»).
-- Поля: id, company_id, status, starts_at, ends_at.

CREATE TYPE contract_status AS ENUM (
    'draft',
    'active',
    'suspended',
    'terminated'
);

CREATE TABLE contracts (
                           id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           company_id  UUID NOT NULL,
                           status      contract_status NOT NULL,
                           starts_at   TIMESTAMPTZ NOT NULL,
                           ends_at     TIMESTAMPTZ NULL,
                           created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
                           updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- ends_at IS NULL — бессрочный договор; иначе окончание строго после начала.
                           CONSTRAINT contracts_period_check
                               CHECK (ends_at IS NULL OR ends_at > starts_at)
);

-- Индексы из README: company_id, status, (company_id, status).
CREATE INDEX idx_contracts_company_id
    ON contracts (company_id);

CREATE INDEX idx_contracts_status
    ON contracts (status);

CREATE INDEX idx_contracts_company_id_status
    ON contracts (company_id, status);

-- У компании только один договор со статусом active (partial unique index).
CREATE UNIQUE INDEX uq_contracts_one_active_per_company
    ON contracts (company_id)
    WHERE status = 'active';

COMMENT ON TABLE contracts IS 'Договоры компаний.';
COMMENT ON COLUMN contracts.id IS 'Идентификатор договора.';
COMMENT ON COLUMN contracts.company_id IS 'Идентификатор компании.';
COMMENT ON COLUMN contracts.status IS 'Статус: draft, active, suspended, terminated (ENUM contract_status).';
COMMENT ON COLUMN contracts.starts_at IS 'Дата и время начала действия договора.';
COMMENT ON COLUMN contracts.ends_at IS 'Дата и время окончания; NULL — бессрочный.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS uq_contracts_one_active_per_company;
DROP INDEX IF EXISTS idx_contracts_company_id_status;
DROP INDEX IF EXISTS idx_contracts_status;
DROP INDEX IF EXISTS idx_contracts_company_id;
DROP TABLE IF EXISTS contracts;
DROP TYPE IF EXISTS contract_status;
-- +goose StatementEnd