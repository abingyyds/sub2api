package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// AgentHandler handles user-facing agent requests
type AgentHandler struct {
	agentService *service.AgentService
}

// NewAgentHandler creates a new AgentHandler
func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

// GetStatus returns the agent status for the current user
// GET /api/v1/agent/status
func (h *AgentHandler) GetStatus(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	status, err := h.agentService.GetAgentStatus(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, status)
}

// Apply submits an agent application
// POST /api/v1/agent/apply
func (h *AgentHandler) Apply(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// note is optional
	}

	if err := h.agentService.ApplyForAgent(c.Request.Context(), subject.UserID, req.Note); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "application submitted"})
}

// Dashboard returns the agent dashboard stats
// GET /api/v1/agent/dashboard
func (h *AgentHandler) Dashboard(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	stats, err := h.agentService.GetDashboard(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, stats)
}

// GetLink returns the agent's invite link
// GET /api/v1/agent/link
func (h *AgentHandler) GetLink(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	code, err := h.agentService.GetInviteLink(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"invite_code": code,
	})
}

// ListSubUsers returns the agent's sub-users
// GET /api/v1/agent/sub-users
func (h *AgentHandler) ListSubUsers(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	search := c.Query("search")

	users, pag, err := h.agentService.ListSubUsers(c.Request.Context(), subject.UserID, params, search)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Paginated(c, users, pag.Total, pag.Page, pag.PageSize)
}

// ListFinancialLogs returns financial logs of the agent's sub-users
// GET /api/v1/agent/financial-logs
func (h *AgentHandler) ListFinancialLogs(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	search := c.Query("search")

	logs, pag, err := h.agentService.ListSubUserFinancialLogs(c.Request.Context(), subject.UserID, params, search)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Paginated(c, logs, pag.Total, pag.Page, pag.PageSize)
}

// ListCommissions returns the agent's commission records
// GET /api/v1/agent/commissions
func (h *AgentHandler) ListCommissions(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	status := c.Query("status")

	commissions, pag, err := h.agentService.ListCommissions(c.Request.Context(), subject.UserID, params, status)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Paginated(c, commissions, pag.Total, pag.Page, pag.PageSize)
}
