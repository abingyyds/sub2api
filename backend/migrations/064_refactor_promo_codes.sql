-- 064: Refactor promo codes from registration bonus to order discount
-- Change promo codes from giving bonus balance at registration to providing discounts at payment

-- 1. promo_codes: rename bonus_amount → discount_amount, add discount_type and min_order_amount
ALTER TABLE promo_codes RENAME COLUMN bonus_amount TO discount_amount;
ALTER TABLE promo_codes ADD COLUMN IF NOT EXISTS discount_type VARCHAR(20) DEFAULT 'fixed';
ALTER TABLE promo_codes ADD COLUMN IF NOT EXISTS min_order_amount INT DEFAULT 0;

-- 2. promo_code_usages: rename bonus_amount → discount_amount, add order_no
ALTER TABLE promo_code_usages RENAME COLUMN bonus_amount TO discount_amount;
ALTER TABLE promo_code_usages ADD COLUMN IF NOT EXISTS order_no VARCHAR(64);

-- 3. payment_orders: add promo_code and discount_amount
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS promo_code VARCHAR(32) DEFAULT '';
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS discount_amount INT DEFAULT 0;
