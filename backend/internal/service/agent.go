package service

import "time"

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
}

// AgentDashboardStats holds aggregated stats for an agent's dashboard.
type AgentDashboardStats struct {
	TotalSubUsers      int64   `json:"total_sub_users"`
	TotalRecharge      float64 `json:"total_recharge"`
	TotalConsumed      float64 `json:"total_consumed"`
	TotalCommission    float64 `json:"total_commission"`
	PendingCommission  float64 `json:"pending_commission"`
	SettledCommission  float64 `json:"settled_commission"`
	TodayNewUsers      int64   `json:"today_new_users"`
	TodayRecharge      float64 `json:"today_recharge"`
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
	ID               int64      `json:"id"`
	Email            string     `json:"email"`
	Username         string     `json:"username,omitempty"`
	IsAgent          bool       `json:"is_agent"`
	AgentStatus      string     `json:"agent_status"`
	CommissionRate   float64    `json:"commission_rate"`
	AgentNote        string     `json:"agent_note,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	InviteCode       string     `json:"invite_code,omitempty"`
	SubUserCount     int64      `json:"sub_user_count"`
	TotalCommission  float64    `json:"total_commission"`
	CreatedAt        time.Time  `json:"created_at"`
}
