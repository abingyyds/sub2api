package service

import (
	"context"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type OrgAuditService struct {
	auditRepo   OrgAuditLogRepository
	orgRepo     OrganizationRepository
}

func NewOrgAuditService(
	auditRepo OrgAuditLogRepository,
	orgRepo OrganizationRepository,
) *OrgAuditService {
	return &OrgAuditService{
		auditRepo: auditRepo,
		orgRepo:   orgRepo,
	}
}

// WriteAPIRequestAudit records an API request audit log entry.
// The depth of content recorded depends on the org's audit_mode.
type WriteAuditInput struct {
	OrgID       int64
	UserID      int64
	MemberID    *int64
	ProjectID   *int64
	UsageLogID  *int64
	Action      string
	Model       *string
	AuditMode   string
	InputTokens *int
	OutputTokens *int
	CostUSD     *float64
	IPAddress   *string
	UserAgent   *string
	// Content fields (filled based on audit_mode)
	RequestBody  string
	ResponseBody string
	// Operation audit detail
	Detail map[string]interface{}
}

func (s *OrgAuditService) WriteAuditLog(ctx context.Context, input *WriteAuditInput) error {
	auditLog := &OrgAuditLog{
		OrgID:       input.OrgID,
		UserID:      input.UserID,
		MemberID:    input.MemberID,
		ProjectID:   input.ProjectID,
		UsageLogID:  input.UsageLogID,
		Action:      input.Action,
		Model:       input.Model,
		AuditMode:   input.AuditMode,
		Flagged:     false,
		InputTokens: input.InputTokens,
		OutputTokens: input.OutputTokens,
		CostUSD:     input.CostUSD,
		IPAddress:   input.IPAddress,
		UserAgent:   input.UserAgent,
		Detail:      input.Detail,
	}

	switch input.AuditMode {
	case OrgAuditModeSummary:
		if input.RequestBody != "" {
			summary := truncate(input.RequestBody, 200)
			auditLog.RequestSummary = &summary
		}
		if input.ResponseBody != "" {
			summary := truncate(input.ResponseBody, 200)
			auditLog.ResponseSummary = &summary
		}
		keywords := extractKeywords(input.RequestBody)
		if len(keywords) > 0 {
			auditLog.Keywords = keywords
		}
	case OrgAuditModeFull:
		if input.RequestBody != "" {
			auditLog.RequestContent = &input.RequestBody
		}
		if input.ResponseBody != "" {
			summary := truncate(input.ResponseBody, 500)
			auditLog.ResponseSummary = &summary
		}
		if input.RequestBody != "" {
			summary := truncate(input.RequestBody, 200)
			auditLog.RequestSummary = &summary
		}
		keywords := extractKeywords(input.RequestBody)
		if len(keywords) > 0 {
			auditLog.Keywords = keywords
		}
	// metadata mode: only metadata fields are filled (default)
	}

	return s.auditRepo.Create(ctx, auditLog)
}

// WriteOperationAudit records a non-API operation audit log (member changes, policy updates, etc.)
func (s *OrgAuditService) WriteOperationAudit(ctx context.Context, orgID, userID int64, action string, detail map[string]interface{}) error {
	auditLog := &OrgAuditLog{
		OrgID:     orgID,
		UserID:    userID,
		Action:    action,
		AuditMode: OrgAuditModeMetadata,
		Flagged:   false,
		Detail:    detail,
	}
	return s.auditRepo.Create(ctx, auditLog)
}

func (s *OrgAuditService) ListAuditLogs(ctx context.Context, orgID int64, page, pageSize int, filters AuditLogFilters) ([]OrgAuditLog, *pagination.PaginationResult, error) {
	params := pagination.NewPaginationParams(page, pageSize)
	return s.auditRepo.List(ctx, orgID, params, filters)
}

func (s *OrgAuditService) GetByID(ctx context.Context, id int64) (*OrgAuditLog, error) {
	return s.auditRepo.GetByID(ctx, id)
}

func (s *OrgAuditService) CountFlagged(ctx context.Context, orgID int64) (int, error) {
	return s.auditRepo.CountFlagged(ctx, orgID)
}

// GetAuditConfig returns the organization's audit mode.
func (s *OrgAuditService) GetAuditConfig(ctx context.Context, orgID int64) (string, error) {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return "", err
	}
	return org.AuditMode, nil
}

// UpdateAuditConfig updates the organization's audit mode.
func (s *OrgAuditService) UpdateAuditConfig(ctx context.Context, orgID int64, auditMode string) error {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}
	org.AuditMode = auditMode
	return s.orgRepo.Update(ctx, org)
}

// truncate returns the first n characters of s, or s if shorter.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

// extractKeywords does simple keyword extraction from request content.
// It extracts unique words longer than 3 characters, up to 20 keywords.
func extractKeywords(content string) []string {
	if content == "" {
		return nil
	}

	// Simple word extraction
	words := strings.Fields(content)
	seen := make(map[string]bool)
	var keywords []string

	for _, word := range words {
		w := strings.ToLower(strings.Trim(word, ".,;:!?\"'()[]{}"))
		if len(w) > 3 && !seen[w] && !isCommonWord(w) {
			seen[w] = true
			keywords = append(keywords, w)
			if len(keywords) >= 20 {
				break
			}
		}
	}
	return keywords
}

func isCommonWord(w string) bool {
	common := map[string]bool{
		"this": true, "that": true, "with": true, "from": true,
		"have": true, "been": true, "will": true, "would": true,
		"could": true, "should": true, "about": true, "which": true,
		"their": true, "there": true, "these": true, "those": true,
		"what": true, "when": true, "where": true, "your": true,
		"they": true, "them": true, "then": true, "than": true,
		"some": true, "into": true, "also": true, "just": true,
		"like": true, "more": true, "very": true, "most": true,
		"only": true, "does": true, "each": true, "make": true,
	}
	return common[w]
}
