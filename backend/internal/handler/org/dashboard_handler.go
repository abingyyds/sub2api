package org

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles org dashboard
type DashboardHandler struct {
	orgService *service.OrganizationService
}

// NewDashboardHandler creates a new org dashboard handler
func NewDashboardHandler(orgService *service.OrganizationService) *DashboardHandler {
	return &DashboardHandler{orgService: orgService}
}

// GetDashboard returns the organization dashboard overview
// GET /api/v1/org/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	dashboard, err := h.orgService.GetDashboard(c.Request.Context(), org.ID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrgDashboardFromService(dashboard))
}
