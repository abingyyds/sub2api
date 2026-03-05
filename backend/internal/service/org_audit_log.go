package service

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// Audit log action constants
const (
	AuditActionAPIRequest   = "api_request"
	AuditActionMemberCreate = "member.create"
	AuditActionMemberUpdate = "member.update"
	AuditActionMemberDelete = "member.delete"
	AuditActionPolicyUpdate = "policy.update"
)

type OrgAuditLog struct {
	ID              int64
	OrgID           int64
	UserID          int64
	MemberID        *int64
	ProjectID       *int64
	UsageLogID      *int64
	Action          string
	Model           *string
	AuditMode       string
	RequestSummary  *string
	RequestContent  *string
	ResponseSummary *string
	Keywords        []string
	Flagged         bool
	FlagReason      *string
	InputTokens     *int
	OutputTokens    *int
	CostUSD         *float64
	IPAddress       *string
	UserAgent       *string
	Detail          map[string]interface{}
	CreatedAt       time.Time
}

type OrgAuditLogRepository interface {
	Create(ctx context.Context, log *OrgAuditLog) error
	GetByID(ctx context.Context, id int64) (*OrgAuditLog, error)
	List(ctx context.Context, orgID int64, params pagination.PaginationParams, filters AuditLogFilters) ([]OrgAuditLog, *pagination.PaginationResult, error)
	CountFlagged(ctx context.Context, orgID int64) (int, error)
}

// AuditLogFilters contains optional filters for listing audit logs.
type AuditLogFilters struct {
	MemberID  *int64
	ProjectID *int64
	Action    string
	Model     string
	Flagged   *bool
	StartDate *time.Time
	EndDate   *time.Time
}
