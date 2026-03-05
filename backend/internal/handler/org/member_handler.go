package org

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// MemberHandler handles org member management
type MemberHandler struct {
	memberService *service.OrgMemberService
}

// NewMemberHandler creates a new org member handler
func NewMemberHandler(memberService *service.OrgMemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

// CreateEmployeeRequest represents create employee request
type CreateEmployeeRequest struct {
	Email           string   `json:"email" binding:"required,email"`
	Password        string   `json:"password" binding:"required,min=6"`
	Username        string   `json:"username"`
	DailyQuotaUSD   *float64 `json:"daily_quota_usd"`
	MonthlyQuotaUSD *float64 `json:"monthly_quota_usd"`
	Notes           *string  `json:"notes"`
}

// UpdateMemberRequest represents update member request
type UpdateMemberRequest struct {
	Role            string   `json:"role"`
	DailyQuotaUSD   *float64 `json:"daily_quota_usd"`
	MonthlyQuotaUSD *float64 `json:"monthly_quota_usd"`
	Notes           *string  `json:"notes"`
}

// List handles listing all members of the organization
// GET /api/v1/org/members
func (h *MemberHandler) List(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	page, pageSize := response.ParsePagination(c)
	members, pagination, err := h.memberService.ListByOrg(c.Request.Context(), org.ID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.OrgMember, 0, len(members))
	for i := range members {
		out = append(out, *dto.OrgMemberFromService(&members[i]))
	}
	response.PaginatedWithResult(c, out, toOrgResponsePagination(pagination))
}

// Create handles creating a new employee
// POST /api/v1/org/members
func (h *MemberHandler) Create(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	member, err := h.memberService.CreateEmployee(c.Request.Context(), org.ID, &service.CreateEmployeeInput{
		Email:           req.Email,
		Password:        req.Password,
		Username:        req.Username,
		DailyQuotaUSD:   req.DailyQuotaUSD,
		MonthlyQuotaUSD: req.MonthlyQuotaUSD,
		Notes:           req.Notes,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Created(c, dto.OrgMemberFromService(member))
}

// GetByID handles getting a member by ID
// GET /api/v1/org/members/:id
func (h *MemberHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	member, err := h.memberService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Verify member belongs to the org
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || member.OrgID != org.ID {
		response.Error(c, 403, "Member does not belong to this organization")
		return
	}

	response.Success(c, dto.OrgMemberFromService(member))
}

// Update handles updating a member
// PUT /api/v1/org/members/:id
func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Verify member belongs to the org
	existing, err := h.memberService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || existing.OrgID != org.ID {
		response.Error(c, 403, "Member does not belong to this organization")
		return
	}

	member, err := h.memberService.Update(c.Request.Context(), id, &service.UpdateMemberInput{
		Role:            req.Role,
		DailyQuotaUSD:   req.DailyQuotaUSD,
		MonthlyQuotaUSD: req.MonthlyQuotaUSD,
		Notes:           req.Notes,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrgMemberFromService(member))
}

// Delete handles removing a member
// DELETE /api/v1/org/members/:id
func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	// Verify member belongs to the org
	existing, err := h.memberService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || existing.OrgID != org.ID {
		response.Error(c, 403, "Member does not belong to this organization")
		return
	}

	if err := h.memberService.Remove(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Member removed successfully"})
}

// Suspend handles suspending a member
// POST /api/v1/org/members/:id/suspend
func (h *MemberHandler) Suspend(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	// Verify member belongs to the org
	existing, err := h.memberService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || existing.OrgID != org.ID {
		response.Error(c, 403, "Member does not belong to this organization")
		return
	}

	member, err := h.memberService.Suspend(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrgMemberFromService(member))
}
