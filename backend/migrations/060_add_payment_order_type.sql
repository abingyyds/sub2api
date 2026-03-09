-- Add order_type and balance_amount columns to payment_orders
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS order_type VARCHAR(20) NOT NULL DEFAULT 'subscription';
ALTER TABLE payment_orders ADD COLUMN IF NOT EXISTS balance_amount DOUBLE PRECISION NOT NULL DEFAULT 0;
