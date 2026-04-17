package service

import "time"

// 分站池流水类型
const (
	// 入账
	SubSiteLedgerTopupOnline    = "topup_online"    // 分站主走主站支付通道充值
	SubSiteLedgerTopupAdmin     = "topup_admin"     // 平台管理员人工加余额
	SubSiteLedgerManualCredit   = "manual_credit"   // 其他人工入账
	SubSiteLedgerProfit         = "profit"          // rate 模式下用户消费产生的分站主利润入账
	SubSiteLedgerWithdrawReject = "withdraw_reject" // 提现被拒，利润退回
	SubSiteLedgerWithdrawCancel = "withdraw_cancel" // 用户撤销提现，利润退回
	// 出账
	SubSiteLedgerConsume          = "consume"            // 分站用户消费时扣 1× 成本（pool 模式）
	SubSiteLedgerOfflineUserTopup = "offline_user_topup" // 分站主线下收款后给用户加余额，从池扣等额
	SubSiteLedgerAutoRestock      = "auto_restock"       // pool 模式用户线上充值到分站主账号时，自动扣池（自动进货）
	SubSiteLedgerWithdrawApply    = "withdraw_apply"     // 提现申请锁定利润
	SubSiteLedgerWithdrawPaid     = "withdraw_paid"      // 提现打款完成（仅记账，不再扣余额）
	SubSiteLedgerRefund           = "refund"             // 退款
	SubSiteLedgerManualDebit      = "manual_debit"       // 其他人工出账
)

// SubSiteLedgerEntry 分站余额池流水记录
type SubSiteLedgerEntry struct {
	ID                int64     `json:"id"`
	SubSiteID         int64     `json:"sub_site_id"`
	TxType            string    `json:"tx_type"`
	DeltaFen          int64     `json:"delta_fen"`
	BalanceAfterFen   int64     `json:"balance_after_fen"`
	RelatedUserID     *int64    `json:"related_user_id,omitempty"`
	RelatedUsageLogID *int64    `json:"related_usage_log_id,omitempty"`
	RelatedOrderID    *int64    `json:"related_order_id,omitempty"`
	OperatorID        *int64    `json:"operator_id,omitempty"`
	Note              string    `json:"note,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
