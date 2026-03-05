package org

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ProjectHandler handles org project management
type ProjectHandler struct {
	projectService *service.OrgProjectService
}

// NewProjectHandler creates a new org project handler
func NewProjectHandler(projectService *service.OrgProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// CreateProjectRequest represents create project request
type CreateProjectRequest struct {
	Name             string   `json:"name" binding:"required"`
	Description      *string  `json:"description"`
	GroupID          *int64   `json:"group_id"`
	AllowedModels    []string `json:"allowed_models"`
	MonthlyBudgetUSD *float64 `json:"monthly_budget_usd"`
}

// UpdateProjectRequest represents update project request
type UpdateProjectRequest struct {
	Name             string   `json:"name" binding:"required"`
	Description      *string  `json:"description"`
	GroupID          *int64   `json:"group_id"`
	AllowedModels    []string `json:"allowed_models"`
	MonthlyBudgetUSD *float64 `json:"monthly_budget_usd"`
	Status           string   `json:"status"`
}

// List handles listing all projects of the organization
// GET /api/v1/org/projects
func (h *ProjectHandler) List(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	page, pageSize := response.ParsePagination(c)
	projects, pagination, err := h.projectService.ListByOrg(c.Request.Context(), org.ID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.OrgProject, 0, len(projects))
	for i := range projects {
		out = append(out, *dto.OrgProjectFromService(&projects[i]))
	}
	response.PaginatedWithResult(c, out, toOrgResponsePagination(pagination))
}

// Create handles creating a new project
// POST /api/v1/org/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok {
		response.Error(c, 403, "Organization not found in context")
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	project, err := h.projectService.Create(c.Request.Context(), org.ID, &service.CreateOrgProjectInput{
		Name:             req.Name,
		Description:      req.Description,
		GroupID:          req.GroupID,
		AllowedModels:    req.AllowedModels,
		MonthlyBudgetUSD: req.MonthlyBudgetUSD,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Created(c, dto.OrgProjectFromService(project))
}

// GetByID handles getting a project by ID
// GET /api/v1/org/projects/:id
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || project.OrgID != org.ID {
		response.Error(c, 403, "Project does not belong to this organization")
		return
	}

	response.Success(c, dto.OrgProjectFromService(project))
}

// Update handles updating a project
// PUT /api/v1/org/projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Verify project belongs to the org
	existing, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || existing.OrgID != org.ID {
		response.Error(c, 403, "Project does not belong to this organization")
		return
	}

	project, err := h.projectService.Update(c.Request.Context(), id, &service.UpdateOrgProjectInput{
		Name:             req.Name,
		Description:      req.Description,
		GroupID:          req.GroupID,
		AllowedModels:    req.AllowedModels,
		MonthlyBudgetUSD: req.MonthlyBudgetUSD,
		Status:           req.Status,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.OrgProjectFromService(project))
}

// Delete handles deleting a project
// DELETE /api/v1/org/projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	// Verify project belongs to the org
	existing, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	org, ok := middleware.GetOrganizationFromContext(c)
	if !ok || existing.OrgID != org.ID {
		response.Error(c, 403, "Project does not belong to this organization")
		return
	}

	if err := h.projectService.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Project deleted successfully"})
}
