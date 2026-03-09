-- 059_add_payment_orders.sql
-- Add payment_orders table for WeChat Pay Native integration

CREATE TABLE IF NOT EXISTS payment_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_key VARCHAR(50) NOT NULL,
    group_id BIGINT NOT NULL,
    amount_fen INT NOT NULL,
    validity_days INT NOT NULL DEFAULT 30,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    pay_method VARCHAR(20) NOT NULL DEFAULT 'wechat_native',
    wechat_transaction_id VARCHAR(64),
    code_url TEXT,
    paid_at TIMESTAMP,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payment_orders_user_id ON payment_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_order_no ON payment_orders(order_no);
CREATE INDEX IF NOT EXISTS idx_payment_orders_status ON payment_orders(status);
