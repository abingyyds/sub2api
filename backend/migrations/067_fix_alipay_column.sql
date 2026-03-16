-- Fix: re-add alipay_trade_no column that was accidentally dropped
-- by migration 066 (which contained goose Down section that the
-- custom migration runner executed as part of the full file)
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS alipay_trade_no VARCHAR(64);
