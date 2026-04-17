-- 分站双模式 + 分站主自有收款 + 提现来源扩展
--
-- sub_sites.mode:
--   pool  — 分站主先向主站充值到独立资金池，线下卖用户余额赚差价
--   rate  — 无独立资金池；用户充值仍到主站，消费时按分站复合倍率扣用户，
--           每级分站按 (compound_i - compound_{i-1}) × base 的差额入账到 balance_fen，
--           分站主走提现出库。
-- owner_payment_config: pool 模式下分站主可配置自有微信/支付宝/易支付凭据，
--   用户通过分站域名给主站余额充值时将使用此凭据收款并自动从分站池扣等额（自动进货）。
-- total_withdrawn_fen: 该分站已提现累计额（rate 模式结算用）。

ALTER TABLE sub_sites
    ADD COLUMN IF NOT EXISTS mode VARCHAR(16) NOT NULL DEFAULT 'pool',
    ADD COLUMN IF NOT EXISTS owner_payment_config JSONB,
    ADD COLUMN IF NOT EXISTS total_withdrawn_fen BIGINT NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_sub_sites_mode ON sub_sites(mode);

-- 提现请求复用给分站主 rate 模式利润提现。
--   source_type = 'agent_commission' 走 agent_profiles.withdrawable_balance
--   source_type = 'sub_site_profit'  走 sub_sites.balance_fen（rate 模式）
ALTER TABLE agent_withdraw_requests
    ADD COLUMN IF NOT EXISTS source_type VARCHAR(32) NOT NULL DEFAULT 'agent_commission',
    ADD COLUMN IF NOT EXISTS source_sub_site_id BIGINT REFERENCES sub_sites(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_agent_withdraw_source_sub_site ON agent_withdraw_requests(source_sub_site_id);
CREATE INDEX IF NOT EXISTS idx_agent_withdraw_source_type ON agent_withdraw_requests(source_type, status);
