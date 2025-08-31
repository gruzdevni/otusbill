CREATE TABLE users(
    guid                UUID PRIMARY KEY        NOT NULL,
    created_at          TIMESTAMPTZ             NOT NULL DEFAULT now()
);

COMMENT ON COLUMN users.guid          IS 'GUID пользователя';
COMMENT ON COLUMN users.created_at    IS 'Дата создания';

CREATE TABLE balance_register(
    id              SERIAL PRIMARY KEY,
    user_guid       UUID            NOT NULL,
    operation_ref   TEXT            NOT NULL DEFAULT '',
    amount          NUMERIC(14, 2)  NOT NULL DEFAULT 0.00,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

COMMENT ON COLUMN balance_register.id               IS 'id строки';
COMMENT ON COLUMN balance_register.user_guid        IS 'GUID пользователя';
COMMENT ON COLUMN balance_register.operation_ref    IS 'Описание операции';
COMMENT ON COLUMN balance_register.amount           IS 'Сумма';
COMMENT ON COLUMN balance_register.created_at       IS 'Дата создания';

CREATE INDEX idx_balance_register_user_guid ON balance_register (user_guid);
