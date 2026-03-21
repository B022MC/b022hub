-- 077_add_payment_orders.sql
-- LinuxDo Credit payment orders

CREATE TABLE IF NOT EXISTS payment_orders (
    id                    BIGSERIAL PRIMARY KEY,
    provider              VARCHAR(32) NOT NULL DEFAULT 'linuxdo_credit',
    out_trade_no          VARCHAR(64) NOT NULL UNIQUE,
    provider_trade_no     VARCHAR(128),
    user_id               BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title                 VARCHAR(64) NOT NULL,
    amount                DECIMAL(20,8) NOT NULL DEFAULT 0,
    credited_amount       DECIMAL(20,8) NOT NULL DEFAULT 0,
    status                VARCHAR(20) NOT NULL DEFAULT 'pending',
    raw_provider_payload  TEXT,
    paid_at               TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_orders_user_created_at
    ON payment_orders(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_payment_orders_provider_status
    ON payment_orders(provider, status);

CREATE INDEX IF NOT EXISTS idx_payment_orders_provider_trade_no
    ON payment_orders(provider_trade_no);
