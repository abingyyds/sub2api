-- +migrate Up
ALTER TABLE payment_orders
    ADD COLUMN IF NOT EXISTS invoice_company_name VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS invoice_tax_id VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS invoice_email VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS invoice_remark TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS invoice_requested_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS invoice_processed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_payment_orders_invoice_requested_at ON payment_orders(invoice_requested_at);
CREATE INDEX IF NOT EXISTS idx_payment_orders_invoice_processed_at ON payment_orders(invoice_processed_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_payment_orders_invoice_processed_at;
DROP INDEX IF EXISTS idx_payment_orders_invoice_requested_at;

ALTER TABLE payment_orders
    DROP COLUMN IF EXISTS invoice_processed_at,
    DROP COLUMN IF EXISTS invoice_requested_at,
    DROP COLUMN IF EXISTS invoice_remark,
    DROP COLUMN IF EXISTS invoice_email,
    DROP COLUMN IF EXISTS invoice_tax_id,
    DROP COLUMN IF EXISTS invoice_company_name;
