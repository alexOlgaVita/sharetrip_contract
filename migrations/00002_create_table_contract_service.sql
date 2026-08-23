-- +goose Up
-- +goose StatementBegin

-- Область хранения: услуги, доступные по договору (см. README.md → «Схема хранения»).
-- Поля: service_code (список услуг), is_available (признак доступности).

CREATE TABLE contract_services (
                                   id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                   contract_id  UUID NOT NULL,
                                   service_code TEXT NOT NULL,
                                   is_available BOOLEAN NOT NULL DEFAULT TRUE,
                                   created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- FK: contract_services.contract_id → contracts.id
                                   CONSTRAINT fk_contract_services_contract
                                       FOREIGN KEY (contract_id)
                                           REFERENCES contracts (id)
                                           ON DELETE CASCADE,

    -- Одна и та же услуга не дублируется в рамках одного договора.
                                   CONSTRAINT uq_contract_services_contract_service
                                       UNIQUE (contract_id, service_code),

                                   CONSTRAINT contract_services_service_code_not_empty
                                       CHECK (char_length(trim(service_code)) > 0)
);

-- Индексы из README: contract_id, service_code.
CREATE INDEX idx_contract_services_contract_id
    ON contract_services (contract_id);

CREATE INDEX idx_contract_services_service_code
    ON contract_services (service_code);

COMMENT ON TABLE contract_services IS 'Услуги, доступные по договору.';
COMMENT ON COLUMN contract_services.contract_id IS 'Идентификатор договора (FK → contracts.id).';
COMMENT ON COLUMN contract_services.service_code IS 'Код услуги в списке доступных услуг договора.';
COMMENT ON COLUMN contract_services.is_available IS 'Признак доступности услуги по договору.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_contract_services_service_code;
DROP INDEX IF EXISTS idx_contract_services_contract_id;
DROP TABLE IF EXISTS contract_services;

-- +goose StatementEnd