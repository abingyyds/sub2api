package org

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
)

// OrgHandlers aggregates all org-related handlers
type OrgHandlers struct {
	Dashboard *DashboardHandler
	Member    *MemberHandler
	Project   *ProjectHandler
	AuditLog  *AuditLogHandler
}

// toOrgResponsePagination converts pagination result for response
func toOrgResponsePagination(p *pagination.PaginationResult) *response.PaginationResult {
	if p == nil {
		return nil
	}
	return &response.PaginationResult{
		Total:    p.Total,
		Page:     p.Page,
		PageSize: p.PageSize,
		Pages:    p.Pages,
	}
}
