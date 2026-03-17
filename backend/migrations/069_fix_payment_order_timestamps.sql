-- Fix payment_orders time columns: timestamp -> timestamptz
-- Existing data is assumed to be in Asia/Shanghai timezone (server default)
ALTER TABLE payment_orders
    ALTER COLUMN paid_at TYPE TIMESTAMPTZ USING paid_at AT TIME ZONE 'Asia/Shanghai',
    ALTER COLUMN expired_at TYPE TIMESTAMPTZ USING expired_at AT TIME ZONE 'Asia/Shanghai',
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'Asia/Shanghai',
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'Asia/Shanghai';
