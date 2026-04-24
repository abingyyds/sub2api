ALTER TABLE payment_orders
  ADD COLUMN IF NOT EXISTS sub_site_id BIGINT REFERENCES sub_sites(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_payment_orders_sub_site_id ON payment_orders(sub_site_id);

UPDATE payment_orders po
SET sub_site_id = (regexp_match(po.plan_key, ':subsite:([0-9]+)$'))[1]::BIGINT
WHERE po.sub_site_id IS NULL
  AND po.order_type = 'balance'
  AND po.plan_key ~ ':subsite:[0-9]+$';
