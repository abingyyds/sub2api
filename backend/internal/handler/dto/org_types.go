package dto

import "time"

// Organization represents an organization for API responses
type Organization struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Description    *string  `json:"description"`
	OwnerUserID    int64    `json:"owner_user_id"`
	BillingMode    string   `json:"billing_mode"`
	Balance        float64  `json:"balance"`
	MonthlyBudgetUSD *float64 `json:"monthly_budget_usd"`
	MaxMembers     int      `json:"max_members"`
	MaxAPIKeys     int      `json:"max_api_keys"`
	Status         string   `json:"status"`
	AuditMode      string   `json:"audit_mode"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Owner       *User `json:"owner,omitempty"`
	MemberCount int   `json:"member_count,omitempty"`
}

// OrgMember represents an organization member for API responses
type OrgMember struct {
	ID               int64    `json:"id"`
	OrgID            int64    `json:"org_id"`
	UserID           int64    `json:"user_id"`
	Role             string   `json:"role"`
	MonthlyQuotaUSD  *float64 `json:"monthly_quota_usd"`
	DailyQuotaUSD    *float64 `json:"daily_quota_usd"`
	MonthlyUsageUSD  float64  `json:"monthly_usage_usd"`
	DailyUsageUSD    float64  `json:"daily_usage_usd"`
	Status           string   `json:"status"`
	Notes            *string  `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	User *User `json:"user,omitempty"`
}

// OrgSubscription represents an org subscription for API responses
type OrgSubscription struct {
	ID        int64     `json:"id"`
	OrgID     int64     `json:"org_id"`
	GroupID   int64     `json:"group_id"`
	StartsAt  time.Time `json:"starts_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"`

	DailyUsageUSD   float64 `json:"daily_usage_usd"`
	WeeklyUsageUSD  float64 `json:"weekly_usage_usd"`
	MonthlyUsageUSD float64 `json:"monthly_usage_usd"`

	AssignedBy *int64    `json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at"`
	Notes      *string   `json:"notes,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Group *Group `json:"group,omitempty"`
}

// OrgDashboard represents the org dashboard overview
type OrgDashboard struct {
	Organization *Organization `json:"organization"`
	MemberCount  int           `json:"member_count"`
	APIKeyCount  int           `json:"api_key_count"`
}

// OrgProject represents an org project for API responses
type OrgProject struct {
	ID                 int64      `json:"id"`
	OrgID              int64      `json:"org_id"`
	Name               string     `json:"name"`
	Description        *string    `json:"description"`
	GroupID            *int64     `json:"group_id"`
	AllowedModels      []string   `json:"allowed_models"`
	MonthlyBudgetUSD   *float64   `json:"monthly_budget_usd"`
	MonthlyUsageUSD    float64    `json:"monthly_usage_usd"`
	MonthlyWindowStart *time.Time `json:"monthly_window_start"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// OrgAuditLog represents an org audit log for API responses
type OrgAuditLog struct {
	ID              int64                  `json:"id"`
	OrgID           int64                  `json:"org_id"`
	UserID          int64                  `json:"user_id"`
	MemberID        *int64                 `json:"member_id"`
	ProjectID       *int64                 `json:"project_id"`
	UsageLogID      *int64                 `json:"usage_log_id"`
	Action          string                 `json:"action"`
	Model           *string                `json:"model"`
	AuditMode       string                 `json:"audit_mode"`
	RequestSummary  *string                `json:"request_summary"`
	RequestContent  *string                `json:"request_content,omitempty"`
	ResponseSummary *string                `json:"response_summary"`
	Keywords        []string               `json:"keywords"`
	Flagged         bool                   `json:"flagged"`
	FlagReason      *string                `json:"flag_reason"`
	InputTokens     *int                   `json:"input_tokens"`
	OutputTokens    *int                   `json:"output_tokens"`
	CostUSD         *float64               `json:"cost_usd"`
	IPAddress       *string                `json:"ip_address"`
	UserAgent       *string                `json:"user_agent"`
	Detail          map[string]interface{} `json:"detail,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}
