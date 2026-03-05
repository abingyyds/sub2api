package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// OrganizationHandler handles admin organization management
type OrganizationHandler struct {
	orgService *service.OrganizationService
}

// NewOrganizationHandler creates a new admin organization handler
func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

// CreateOrganizationRequest represents create org request
type CreateOrganizationRequest struct {
	Name           string   `json:"name" binding:"required"`
	Slug           string   `json:"slug"`
	Description    *string  `json:"description"`
	OwnerUserID    int64    `json:"owner_user_id" binding:"required"`
	BillingMode    string   `json:"billing_mode"`
	Balance        float64  `json:"balance"`
	MonthlyBudgetUSD *float64 `json:"monthly_budget_usd"`
	MaxMembers     int      `json:"max_members"`
	MaxAPIKeys     int      `json:"max_api_keys"`
}

// UpdateOrganizationRequest represents update org request
type UpdateOrganizationRequest struct {
	Name           string   `json:"name" binding:"required"`
	Slug           string   `json:"slug" binding:"required"`
	Description    *string  `json:"description"`
	BillingMode    string   `json:"billing_mode"`
	MonthlyBudgetUSD *float64 `json:"monthly_budget_usd"`
	MaxMembers     int      `json:"max_members"`
	MaxAPIKeys     int      `json:"max_api_keys"`
	Status         string   `json:"status"`
	AuditMode      string   `json:"audit_mode"`
}

// UpdateBalanceRequest represents balance update request
type UpdateOrgBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	Action string  `json:"action" binding:"required,oneof=add set"`
}

// List handles listing all organizations
// GET /api/v1/admin/organizations
func (h *OrganizationHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	status := c.Query("status")
	search := c.Query("search")

	orgs, pagination, err := h.orgService.List(c.Request.Context(), page, pageSize, status, search)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.Organization, 0, len(orgs))
	for i := range orgs {
		out = append(out, *dto.OrganizationFromService(&orgs[i]))
	}
	response.PaginatedWithResult(c, out, toResponsePagination(pagination))
}

// GetByID handles getting an organization by ID
// GET /api/v1/admin/organizations/:id
func (h *OrganizationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	org, err := h.orgService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrganizationFromService(org))
}

// Create handles creating a new organization
// POST /api/v1/admin/organizations
func (h *OrganizationHandler) Create(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	org, err := h.orgService.Create(c.Request.Context(), &service.CreateOrganizationInput{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		OwnerUserID:      req.OwnerUserID,
		BillingMode:      req.BillingMode,
		Balance:          req.Balance,
		MonthlyBudgetUSD: req.MonthlyBudgetUSD,
		MaxMembers:       req.MaxMembers,
		MaxAPIKeys:       req.MaxAPIKeys,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Created(c, dto.OrganizationFromService(org))
}

// Update handles updating an organization
// PUT /api/v1/admin/organizations/:id
func (h *OrganizationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	org, err := h.orgService.Update(c.Request.Context(), id, &service.UpdateOrganizationInput{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		BillingMode:      req.BillingMode,
		MonthlyBudgetUSD: req.MonthlyBudgetUSD,
		MaxMembers:       req.MaxMembers,
		MaxAPIKeys:       req.MaxAPIKeys,
		Status:           req.Status,
		AuditMode:        req.AuditMode,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrganizationFromService(org))
}

// Delete handles deleting an organization
// DELETE /api/v1/admin/organizations/:id
func (h *OrganizationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	if err := h.orgService.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Organization deleted successfully"})
}

// UpdateBalance handles updating organization balance
// POST /api/v1/admin/organizations/:id/balance
func (h *OrganizationHandler) UpdateBalance(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	var req UpdateOrgBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	org, err := h.orgService.UpdateBalance(c.Request.Context(), id, &service.UpdateOrgBalanceInput{
		Amount: req.Amount,
		Action: req.Action,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrganizationFromService(org))
}
