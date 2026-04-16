-- 分站余额池 + 统一流水表
-- 支撑的业务模型：
--   1. 分站主从主站"采购"余额到分站池（balance_fen）
--   2. 分站用户消费时，用户扣 multiplier×、分站池扣 1×，差额即分站主利润（保留在池内）
--   3. 线上充值（走主站支付）+ 线下充值（分站主后台手工）由 allow_online_topup / allow_offline_topup 开关控制

ALTER TABLE sub_sites
    ADD COLUMN IF NOT EXISTS balance_fen BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_topup_fen BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_consumed_fen BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS allow_online_topup BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS allow_offline_topup BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE IF NOT EXISTS sub_site_ledger (
    id BIGSERIAL PRIMARY KEY,
    sub_site_id BIGINT NOT NULL REFERENCES sub_sites(id) ON DELETE CASCADE,
    tx_type VARCHAR(32) NOT NULL,
    delta_fen BIGINT NOT NULL,
    balance_after_fen BIGINT NOT NULL,
    related_user_id BIGINT,
    related_usage_log_id BIGINT,
    related_order_id BIGINT,
    operator_id BIGINT,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sub_site_ledger_site_created ON sub_site_ledger(sub_site_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sub_site_ledger_related_user ON sub_site_ledger(related_user_id);
CREATE INDEX IF NOT EXISTS idx_sub_site_ledger_tx_type ON sub_site_ledger(sub_site_id, tx_type);

-- 注意：sub_site_group_prices / sub_site_recharge_prices（由 079 引入）已在业务层废弃，
-- 一期保留表结构仅为兼容，二期随同相关代码一并 DROP。
