-- Add price_fen column to groups for direct pricing on groups
ALTER TABLE groups ADD COLUMN IF NOT EXISTS price_fen INTEGER NOT NULL DEFAULT 0;
