package service

import "time"

// AgentProfile stores lightweight onboarding and risk-control state for an agent.
type AgentProfile struct {
	UserID              int64      `json:"user_id"`
	RealName            string     `json:"real_name"`
	IDCardNo            string     `json:"id_card_no"`
	Phone               string     `json:"phone"`
	IdentityStatus      string     `json:"identity_status"`
	IdentitySubmittedAt *time.Time `json:"identity_submitted_at,omitempty"`
	ContractStatus      string     `json:"contract_status"`
	ContractVersion     string     `json:"contract_version"`
	ContractSignedAt    *time.Time `json:"contract_signed_at,omitempty"`
	ContractIP          string     `json:"contract_ip,omitempty"`
	ContractFileData    string     `json:"contract_file_data,omitempty"`
	ActivationOrderID   *int64     `json:"activation_order_id,omitempty"`
	ActivationFeePaidAt *time.Time `json:"activation_fee_paid_at,omitempty"`
	FrozenBalance       float64    `json:"frozen_balance"`
	WithdrawableBalance float64    `json:"withdrawable_balance"`
	TotalWithdrawn      float64    `json:"total_withdrawn"`
	IsFrozen            bool       `json:"is_frozen"`
	FrozenReason        string     `json:"frozen_reason,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// AgentWalletSummary provides a simplified wallet view for the current agent.
type AgentWalletSummary struct {
	SiteBalance         float64 `json:"site_balance"`
	FrozenBalance       float64 `json:"frozen_balance"`
	WithdrawableBalance float64 `json:"withdrawable_balance"`
	TotalWithdrawn      float64 `json:"total_withdrawn"`
}

// AgentWithdrawWindow describes when withdrawals can be submitted.
type AgentWithdrawWindow struct {
	Weekday   int    `json:"weekday"`
	StartHour int    `json:"start_hour"`
	EndHour   int    `json:"end_hour"`
	Label     string `json:"label"`
}

// AgentStatusView is the aggregated status payload returned to the frontend.
type AgentStatusView struct {
	Enabled            bool                `json:"enabled"`
	IsAgent            bool                `json:"is_agent"`
	AgentStatus        string              `json:"agent_status"`
	CommissionRate     float64             `json:"commission_rate"`
	InviteCode         string              `json:"invite_code,omitempty"`
	CanApply           bool                `json:"can_apply"`
	ActivationFee      float64             `json:"activation_fee"`
	ContractTemplate   string              `json:"contract_template,omitempty"`
	Profile            *AgentProfile       `json:"profile,omitempty"`
	Wallet             AgentWalletSummary  `json:"wallet"`
	WithdrawFreezeDays int                 `json:"withdraw_freeze_days"`
	WithdrawWindow     AgentWithdrawWindow `json:"withdraw_window"`
}

// AgentCommission represents a commission record for an agent.
type AgentCommission struct {
	ID               int64      `json:"id"`
	AgentID          int64      `json:"agent_id"`
	UserID           int64      `json:"user_id"`
	OrderID          *int64     `json:"order_id,omitempty"`
	SourceType       string     `json:"source_type"`
	SourceAmount     float64    `json:"source_amount"`
	CommissionRate   float64    `json:"commission_rate"`
	CommissionAmount float64    `json:"commission_amount"`
	Status           string     `json:"status"`
	SettledAt        *time.Time `json:"settled_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Joined fields for display
	UserEmail  string `json:"user_email,omitempty"`
	AgentEmail string `json:"agent_email,omitempty"`
	OrderNo    string `json:"order_no,omitempty"`
}

// AgentSubUser represents a sub-user under an agent.
type AgentSubUser struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username,omitempty"`
	Balance        float64   `json:"balance"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	TotalRecharge  float64   `json:"total_recharge"`
	TotalConsumed  float64   `json:"total_consumed"`
	OrderCount     int       `json:"order_count"`
	CommissionRate *float64  `json:"commission_rate"`
	IsAgent        bool      `json:"is_agent"`
}

// AgentDashboardStats holds aggregated stats for an agent's dashboard.
type AgentDashboardStats struct {
	TotalSubUsers       int64   `json:"total_sub_users"`
	TotalRecharge       float64 `json:"total_recharge"`
	TotalConsumed       float64 `json:"total_consumed"`
	TotalCommission     float64 `json:"total_commission"`
	PendingCommission   float64 `json:"pending_commission"`
	SettledCommission   float64 `json:"settled_commission"`
	TodayNewUsers       int64   `json:"today_new_users"`
	TodayRecharge       float64 `json:"today_recharge"`
	SiteBalance         float64 `json:"site_balance"`
	FrozenBalance       float64 `json:"frozen_balance"`
	WithdrawableBalance float64 `json:"withdrawable_balance"`
	TotalWithdrawn      float64 `json:"total_withdrawn"`
}

// AgentFinancialLog represents a financial event for a sub-user (payment or consumption).
type AgentFinancialLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	UserEmail string    `json:"user_email"`
	Type      string    `json:"type"` // "payment" or "usage"
	Amount    float64   `json:"amount"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

// AgentInfo holds agent profile information (for admin listing).
type AgentInfo struct {
	ID                  int64      `json:"id"`
	Email               string     `json:"email"`
	Username            string     `json:"username,omitempty"`
	IsAgent             bool       `json:"is_agent"`
	AgentStatus         string     `json:"agent_status"`
	CommissionRate      float64    `json:"agent_commission_rate"`
	AgentNote           string     `json:"agent_note,omitempty"`
	ApprovedAt          *time.Time `json:"agent_approved_at,omitempty"`
	InviteCode          string     `json:"invite_code,omitempty"`
	SubUserCount        int64      `json:"sub_user_count"`
	TotalCommission     float64    `json:"total_commission"`
	PendingCommission   float64    `json:"pending_commission"`
	RealName            string     `json:"real_name,omitempty"`
	IDCardNo            string     `json:"id_card_no,omitempty"`
	Phone               string     `json:"phone,omitempty"`
	IdentityStatus      string     `json:"identity_status,omitempty"`
	IdentitySubmittedAt *time.Time `json:"identity_submitted_at,omitempty"`
	ContractStatus      string     `json:"contract_status,omitempty"`
	ContractVersion     string     `json:"contract_version,omitempty"`
	ContractSignedAt    *time.Time `json:"contract_signed_at,omitempty"`
	ContractIP          string     `json:"contract_ip,omitempty"`
	ContractFileData    string     `json:"contract_file_data,omitempty"`
	ActivationFeePaidAt *time.Time `json:"activation_fee_paid_at,omitempty"`
	IsFrozen            bool       `json:"is_frozen"`
	FrozenReason        string     `json:"frozen_reason,omitempty"`
	FrozenBalance       float64    `json:"frozen_balance"`
	WithdrawableBalance float64    `json:"withdrawable_balance"`
	TotalWithdrawn      float64    `json:"total_withdrawn"`
	CreatedAt           time.Time  `json:"created_at"`
}

// 提现请求状态
const (
	WithdrawStatusPending  = "pending"
	WithdrawStatusApproved = "approved"
	WithdrawStatusRejected = "rejected"
	WithdrawStatusPaid     = "paid"
	WithdrawStatusCancelled = "cancelled"
)

// 提现来源
const (
	WithdrawSourceAgentCommission = "agent_commission"
	WithdrawSourceSubSiteProfit   = "sub_site_profit"
)

// WithdrawRequest 提现请求（复用 agent_withdraw_requests 表）。
type WithdrawRequest struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	Amount          float64    `json:"amount"`
	AlipayName      string     `json:"alipay_name"`
	AlipayAccount   string     `json:"alipay_account"`
	AlipayQRImage   string     `json:"alipay_qr_image,omitempty"`
	Status          string     `json:"status"`
	ReviewNote      string     `json:"review_note,omitempty"`
	SourceType      string     `json:"source_type"`
	SourceSubSiteID *int64     `json:"source_sub_site_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	// joined
	UserEmail       string `json:"user_email,omitempty"`
	SubSiteName     string `json:"sub_site_name,omitempty"`
}
