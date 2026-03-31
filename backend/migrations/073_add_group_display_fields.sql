-- Add display fields for group cards in key management page
-- tags: custom tag labels like "官方 API", "逆向", "推荐", "暂不可用" etc.
-- display_price: custom price text like "6 块 / 1 美元"
-- display_discount: custom discount text like "8.3折"

ALTER TABLE groups ADD COLUMN IF NOT EXISTS tags jsonb DEFAULT '[]';
ALTER TABLE groups ADD COLUMN IF NOT EXISTS display_price text NOT NULL DEFAULT '';
ALTER TABLE groups ADD COLUMN IF NOT EXISTS display_discount text NOT NULL DEFAULT '';
