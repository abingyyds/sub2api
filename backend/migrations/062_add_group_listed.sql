-- Add listed field to groups table
-- Controls whether a group appears as a purchasable plan on the pricing page
ALTER TABLE groups ADD COLUMN IF NOT EXISTS listed BOOLEAN NOT NULL DEFAULT false;
