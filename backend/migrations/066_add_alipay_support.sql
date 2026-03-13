-- Add Alipay support to payment_orders table
-- +goose Up
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS alipay_trade_no VARCHAR(64);

COMMENT ON COLUMN payment_orders.pay_method IS 'wechat_native | alipay_native';
COMMENT ON COLUMN payment_orders.wechat_transaction_id IS '微信交易号';
COMMENT ON COLUMN payment_orders.alipay_trade_no IS '支付宝交易号';

-- +goose Down
ALTER TABLE payment_orders DROP COLUMN IF EXISTS alipay_trade_no;
