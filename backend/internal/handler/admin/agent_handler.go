package admin

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// AgentHandler handles admin agent management
type AgentHandler struct {
	agentService *service.AgentService
}

// NewAgentHandler creates a new admin agent handler
func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

// List handles listing all agents with pagination
// GET /api/v1/admin/agents
func (h *AgentHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	status := strings.TrimSpace(c.Query("status"))
	search := strings.TrimSpace(c.Query("search"))
	if len(search) > 100 {
		search = search[:100]
	}

	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	agents, pag, err := h.agentService.AdminListAgents(c.Request.Context(), params, status, search)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Paginated(c, agents, pag.Total, pag.Page, pag.PageSize)
}

// Approve approves an agent application
// POST /api/v1/admin/agents/:id/approve
func (h *AgentHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid agent ID")
		return
	}

	var req struct {
		CommissionRate float64 `json:"commission_rate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// commission_rate is optional, will use default
	}

	if err := h.agentService.AdminApproveAgent(c.Request.Context(), id, req.CommissionRate); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "agent approved"})
}

// Reject rejects an agent application
// POST /api/v1/admin/agents/:id/reject
func (h *AgentHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid agent ID")
		return
	}

	if err := h.agentService.AdminRejectAgent(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "agent rejected"})
}

// Update updates the agent's commission rate
// PUT /api/v1/admin/agents/:id
func (h *AgentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid agent ID")
		return
	}

	var req struct {
		CommissionRate float64 `json:"commission_rate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "commission_rate is required")
		return
	}

	if err := h.agentService.AdminUpdateCommissionRate(c.Request.Context(), id, req.CommissionRate); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "commission rate updated"})
}

// Settle settles all pending commissions for an agent
// POST /api/v1/admin/agents/:id/settle
func (h *AgentHandler) Settle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid agent ID")
		return
	}

	amount, err := h.agentService.AdminSettleCommissions(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"settled_amount": amount})
}
