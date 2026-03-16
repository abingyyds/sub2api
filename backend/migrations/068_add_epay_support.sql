-- 068: Add Epay (易支付) support
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS epay_trade_no VARCHAR(64);
COMMENT ON COLUMN payment_orders.pay_method IS 'wechat_native | alipay_native | epay_alipay | epay_wxpay';
