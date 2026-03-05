package org

import (
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuditLogHandler handles org audit log management
type AuditLogHandler struct {
	auditService *service.OrgAuditService
}

// NewAuditLogHandler creates a new org audit log handler
func NewAuditLogHandler(auditService *service.OrgAuditService) *AuditLogHandler {
	return &AuditLogHandler{auditService: auditService}
}

// UpdateAuditConfigRequest represents update audit config request
type UpdateAuditConfigRequest struct {
	AuditMode string `json:"audit_mode" binding:"required,oneof=metadata summary full"`
}

// List handles listing audit logs
// GET /api/v1/org/audit-logs
func (h *AuditLogHandler) List(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	page, pageSize := response.ParsePagination(c)

	// Parse filters from query params
	var filters service.AuditLogFilters
	if memberIDStr := c.Query("member_id"); memberIDStr != "" {
		if memberID, err := strconv.ParseInt(memberIDStr, 10, 64); err == nil {
			filters.MemberID = &memberID
		}
	}
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		if projectID, err := strconv.ParseInt(projectIDStr, 10, 64); err == nil {
			filters.ProjectID = &projectID
		}
	}
	if action := c.Query("action"); action != "" {
		filters.Action = action
	}
	if model := c.Query("model"); model != "" {
		filters.Model = model
	}
	if flaggedStr := c.Query("flagged"); flaggedStr != "" {
		flagged := flaggedStr == "true"
		filters.Flagged = &flagged
	}
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filters.StartDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filters.EndDate = &t
		}
	}

	logs, pagination, err := h.auditService.ListAuditLogs(c.Request.Context(), org.ID, page, pageSize, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.OrgAuditLog, 0, len(logs))
	for i := range logs {
		out = append(out, *dto.OrgAuditLogFromService(&logs[i]))
	}
	response.PaginatedWithResult(c, out, toOrgResponsePagination(pagination))
}

// GetByID handles getting an audit log by ID
// GET /api/v1/org/audit-logs/:id
func (h *AuditLogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid audit log ID")
		return
	}

	auditLog, err := h.auditService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if auditLog == nil {
		response.Error(c, 404, "Audit log not found")
		return
	}

	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || auditLog.OrgID != org.ID {
		response.Error(c, 403, "Audit log does not belong to this organization")
		return
	}

	response.Success(c, dto.OrgAuditLogFromService(auditLog))
}

// CountFlagged handles getting the count of flagged audit logs
// GET /api/v1/org/audit-logs/flagged-count
func (h *AuditLogHandler) CountFlagged(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	count, err := h.auditService.CountFlagged(c.Request.Context(), org.ID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"count": count})
}

// GetAuditConfig handles getting the audit config
// GET /api/v1/org/audit-config
func (h *AuditLogHandler) GetAuditConfig(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	auditMode, err := h.auditService.GetAuditConfig(c.Request.Context(), org.ID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"audit_mode": auditMode})
}

// UpdateAuditConfig handles updating the audit config
// PUT /api/v1/org/audit-config
func (h *AuditLogHandler) UpdateAuditConfig(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	var req UpdateAuditConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.auditService.UpdateAuditConfig(c.Request.Context(), org.ID, req.AuditMode); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"audit_mode": req.AuditMode})
}
