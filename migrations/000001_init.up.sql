-- create pgcrypto extension if not exists
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- user_role enum type
CREATE TYPE user_role AS ENUM ('admin', 'operator');

-- users table
CREATE TABLE IF NOT EXISTS users (
    name VARCHAR(255),
    hashed_password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,

    CONSTRAINT pk_users
        PRIMARY KEY (name)
);

INSERT INTO users (name, hashed_password, role)
VALUES
    ('admin', crypt('admin', gen_salt('bf')), 'admin'),
    ('operator', crypt('operator', gen_salt('bf')), 'operator');


-- currencies table
CREATE TABLE IF NOT EXISTS currencies (
    code CHAR(3),
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(8) NOT NULL,
    minor_unit SMALLINT NOT NULL,

    CONSTRAINT pk_currencies
        PRIMARY KEY (code),

    CONSTRAINT chk_currencies_minor_unit_range
        CHECK (minor_unit BETWEEN 0 AND 4)
);

INSERT INTO currencies (code, name, symbol, minor_unit)
VALUES
    ('USD', 'Доллар США', '$', 2),
    ('BYN', 'Белорусский рубль', 'Br', 2),
    ('RUB', 'Российский рубль', '₽', 2),
    ('EUR', 'Евро', '€', 2),
    ('GBP', 'Фунт стерлингов', '£', 2),
    ('CHF', 'Швейцарский франк', 'Fr', 2),
    ('PLN', 'Польский злотый', 'zł', 2),
    ('CNY', 'Китайский юань', '¥', 2),
    ('JPY', 'Японская иена', '¥', 0),
    ('KZT', 'Казахстанский тенге', '₸', 2),
    ('TRY', 'Турецкая лира', '₺', 2),
    ('AED', 'Дирхам ОАЭ', 'د.إ', 2),
    ('GEL', 'Грузинский лари', '₾', 2),
    ('UAH', 'Украинская гривна', '₴', 2),
    ('CZK', 'Чешская крона', 'Kč', 2),
    ('CAD', 'Канадский доллар', '$', 2),
    ('AUD', 'Австралийский доллар', '$', 2)
;

-- exchange_rates table
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID DEFAULT gen_random_uuid(),
    base_currency_code CHAR(3) NOT NULL,
    quote_currency_code CHAR(3) NOT NULL,
    buy_rate NUMERIC(20, 8) NOT NULL,
    sell_rate NUMERIC(20, 8) NOT NULL,
    is_buy_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_sell_active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT pk_exchange_rates
        PRIMARY KEY (id),

    CONSTRAINT fk_exchange_rates_base_currency
        FOREIGN KEY (base_currency_code)
        REFERENCES currencies(code)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_exchange_rates_quote_currency
        FOREIGN KEY (quote_currency_code)
        REFERENCES currencies(code)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT chk_exchange_rates_base_currency_not_equal_to_quote_currency
        CHECK (base_currency_code <> quote_currency_code),

    CONSTRAINT chk_exchange_rates_buy_rate_positive
        CHECK (buy_rate > 0),

    CONSTRAINT chk_exchange_rates_sell_rate_positive
        CHECK (sell_rate > 0),

    CONSTRAINT chk_exchange_rates_sell_rate_greater_than_or_equal_to_buy_rate
        CHECK (sell_rate >= buy_rate),

    CONSTRAINT uq_exchange_rates_currency_pair
        UNIQUE (base_currency_code, quote_currency_code)
);

INSERT INTO exchange_rates (base_currency_code, quote_currency_code, buy_rate, sell_rate)
VALUES
    ('USD', 'BYN', 3.20000000, 3.30000000),
    ('EUR', 'BYN', 3.48000000, 3.58000000),
    ('RUB', 'BYN', 0.03500000, 0.03700000),
    ('GBP', 'BYN', 4.05000000, 4.20000000),
    ('CHF', 'BYN', 3.62000000, 3.76000000),
    ('PLN', 'BYN', 0.80000000, 0.85000000),
    ('CNY', 'BYN', 0.44000000, 0.47000000),
    ('JPY', 'BYN', 0.02100000, 0.02300000),
    ('KZT', 'BYN', 0.00680000, 0.00740000),
    ('TRY', 'BYN', 0.09500000, 0.11000000),
    ('AED', 'BYN', 0.87000000, 0.92000000),
    ('GEL', 'BYN', 1.18000000, 1.28000000),
    ('UAH', 'BYN', 0.07800000, 0.08600000),
    ('CZK', 'BYN', 0.14000000, 0.15500000),
    ('CAD', 'BYN', 2.32000000, 2.45000000),
    ('AUD', 'BYN', 2.09000000, 2.22000000);


-- exchange_rates_history table
CREATE TABLE IF NOT EXISTS exchange_rates_history (
    exchange_rate_id UUID NOT NULL,
    buy_rate NUMERIC(20, 8) NOT NULL,
    sell_rate NUMERIC(20, 8) NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    valid_to TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT pk_exchange_rates_history
        PRIMARY KEY (exchange_rate_id, valid_from),

    CONSTRAINT fk_exchange_rates_history_exchange_rate
        FOREIGN KEY (exchange_rate_id)
        REFERENCES exchange_rates(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT chk_exchange_rates_history_buy_rate_positive
        CHECK (buy_rate > 0),

    CONSTRAINT chk_exchange_rates_history_sell_rate_positive
        CHECK (sell_rate > 0),

    CONSTRAINT chk_exchange_rates_history_sell_rate_greater_than_or_equal_to_buy_rate
        CHECK (sell_rate >= buy_rate)
);

CREATE INDEX IF NOT EXISTS idx_exchange_rates_history_valid_from
    ON exchange_rates_history (valid_from DESC);

CREATE INDEX IF NOT EXISTS idx_exchange_rates_history_rate_valid_from
    ON exchange_rates_history (exchange_rate_id, valid_from DESC);

INSERT INTO exchange_rates_history (exchange_rate_id, buy_rate, sell_rate, valid_from, valid_to)
SELECT id, buy_rate, sell_rate, updated_at, updated_at
FROM exchange_rates
WHERE NOT EXISTS (
    SELECT 1
    FROM exchange_rates_history
    WHERE exchange_rates_history.exchange_rate_id = exchange_rates.id
);

-- user_actions table
CREATE TABLE IF NOT EXISTS user_actions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_name  VARCHAR(255) NOT NULL,
    session_id  VARCHAR(128) NOT NULL,
    request     JSONB NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_actions_actor
        FOREIGN KEY (actor_name)
        REFERENCES users(name)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_actions_session
    ON user_actions (session_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_user_actions_actor
    ON user_actions (actor_name, created_at DESC);

-- exchanges table
CREATE TABLE IF NOT EXISTS exchanges (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_name       VARCHAR(255) NOT NULL,
    session_id          VARCHAR(128) NOT NULL,
    base_currency_code  CHAR(3) NOT NULL,
    quote_currency_code CHAR(3) NOT NULL,
    side                VARCHAR(8) NOT NULL,
    amount              NUMERIC(20, 8) NOT NULL,
    rate                NUMERIC(20, 8) NOT NULL,
    result_amount       NUMERIC(20, 8) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_exchanges_operator
        FOREIGN KEY (operator_name)
        REFERENCES users(name)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_exchanges_base_currency
        FOREIGN KEY (base_currency_code)
        REFERENCES currencies(code)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_exchanges_quote_currency
        FOREIGN KEY (quote_currency_code)
        REFERENCES currencies(code)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT chk_exchanges_side
        CHECK (side IN ('buy', 'sell')),

    CONSTRAINT chk_exchanges_amount_positive
        CHECK (amount > 0),

    CONSTRAINT chk_exchanges_rate_positive
        CHECK (rate > 0),

    CONSTRAINT chk_exchanges_result_amount_positive
        CHECK (result_amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_exchanges_created_at
    ON exchanges (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_exchanges_session
    ON exchanges (session_id, created_at DESC);
