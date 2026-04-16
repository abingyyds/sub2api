ALTER TABLE sub_sites
    ADD COLUMN IF NOT EXISTS consume_rate_multiplier DOUBLE PRECISION NOT NULL DEFAULT 1.0;

UPDATE sub_sites
SET consume_rate_multiplier = 1.0
WHERE consume_rate_multiplier <= 0;
