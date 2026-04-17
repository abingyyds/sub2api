package admin

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type SubSiteHandler struct {
	subSiteService  *service.SubSiteService
	withdrawService *service.WithdrawService
}

func NewSubSiteHandler(subSiteService *service.SubSiteService, withdrawService *service.WithdrawService) *SubSiteHandler {
	return &SubSiteHandler{subSiteService: subSiteService, withdrawService: withdrawService}
}

func (h *SubSiteHandler) GetPlatformConfig(c *gin.Context) {
	config, err := h.subSiteService.GetPlatformConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *SubSiteHandler) UpdatePlatformConfig(c *gin.Context) {
	var req service.UpdatePlatformSubSiteConfigInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	config, err := h.subSiteService.UpdatePlatformConfig(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *SubSiteHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.subSiteService.List(c.Request.Context(), params, c.Query("search"), c.Query("status"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

func (h *SubSiteHandler) Create(c *gin.Context) {
	var req service.CreateSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	site, err := h.subSiteService.Create(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

func (h *SubSiteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req service.UpdateSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = id
	site, err := h.subSiteService.Update(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

func (h *SubSiteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	if err := h.subSiteService.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "sub-site deleted"})
}

type adminSubSiteTopupRequest struct {
	AmountFen int64  `json:"amount_fen"`
	Note      string `json:"note"`
}

// Topup 给分站池人工加余额（tx_type=topup_admin）。
func (h *SubSiteHandler) Topup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req adminSubSiteTopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if req.AmountFen <= 0 {
		response.BadRequest(c, "amount_fen must be positive")
		return
	}
	var operatorID int64
	if subject, ok := middleware2.GetAuthSubjectFromContext(c); ok {
		operatorID = subject.UserID
	}
	site, err := h.subSiteService.AdminTopupPool(c.Request.Context(), id, req.AmountFen, operatorID, strings.TrimSpace(req.Note))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

// Ledger 查询分站池流水。
func (h *SubSiteHandler) Ledger(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.subSiteService.ListPoolLedger(c.Request.Context(), id, params, c.Query("tx_type"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if items == nil {
		items = []service.SubSiteLedgerEntry{}
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

type setModeRequest struct {
	Mode string `json:"mode" binding:"required"`
}

// SetMode 管理员切换分站模式（pool ↔ rate）。
func (h *SubSiteHandler) SetMode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req setModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	hasPending, err := h.withdrawService.HasPendingWithdrawForSubSite(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	site, err := h.subSiteService.SetSubSiteMode(c.Request.Context(), id, req.Mode, hasPending)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}
