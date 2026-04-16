package handler

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// SubSiteAdminHandler handles sub-site owner admin-panel requests.
// 路由前缀：/api/v1/subsite/admin/:siteId/*
// 中间件链：JWTAuth → SubSiteOwnerRequired（两层鉴权，service 层是第三层）。
type SubSiteAdminHandler struct {
	adminService   *service.SubSiteAdminService
	subSiteService *service.SubSiteService
}

func NewSubSiteAdminHandler(adminService *service.SubSiteAdminService, subSiteService *service.SubSiteService) *SubSiteAdminHandler {
	return &SubSiteAdminHandler{adminService: adminService, subSiteService: subSiteService}
}

func (h *SubSiteAdminHandler) siteIDOrAbort(c *gin.Context) (int64, bool) {
	siteID, err := strconv.ParseInt(c.Param("siteId"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return 0, false
	}
	return siteID, true
}

func (h *SubSiteAdminHandler) ownerOrAbort(c *gin.Context) (int64, bool) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return 0, false
	}
	return subject.UserID, true
}

// Dashboard 返回分站概览卡片数据。
func (h *SubSiteAdminHandler) Dashboard(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	stats, err := h.adminService.Dashboard(c.Request.Context(), ownerID, siteID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, stats)
}

// SiteDetail 返回当前分站的详情（复用 SubSiteService.GetOwnedSite）。
func (h *SubSiteAdminHandler) SiteDetail(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	item, err := h.subSiteService.GetOwnedSite(c.Request.Context(), ownerID, siteID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

// UpdateSite 更新当前分站的设置（复用 UpdateOwnedSite）。
func (h *SubSiteAdminHandler) UpdateSite(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	var req service.UpdateOwnedSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = siteID
	item, err := h.subSiteService.UpdateOwnedSite(c.Request.Context(), ownerID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

// ListUsers 返回绑定到本分站的用户。
func (h *SubSiteAdminHandler) ListUsers(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.adminService.ListUsers(c.Request.Context(), ownerID, siteID, params, c.Query("search"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

// ListOrders 返回本分站用户的订单。
func (h *SubSiteAdminHandler) ListOrders(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.adminService.ListOrders(c.Request.Context(), ownerID, siteID, params, c.Query("status"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

// ListUsage 返回本分站用户的用量。
func (h *SubSiteAdminHandler) ListUsage(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.adminService.ListUsage(c.Request.Context(), ownerID, siteID, params, c.Query("model"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

type offlineTopupRequest struct {
	UserID    int64  `json:"user_id"`
	AmountFen int64  `json:"amount_fen"`
	Note      string `json:"note"`
}

// OfflineTopup 给本分站用户线下加余额（从分站池等额扣除）。
func (h *SubSiteAdminHandler) OfflineTopup(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	var req offlineTopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if req.AmountFen <= 0 {
		response.BadRequest(c, "amount_fen must be positive")
		return
	}
	site, err := h.subSiteService.OfflineTopupUser(c.Request.Context(), ownerID, siteID, req.UserID, req.AmountFen, strings.TrimSpace(req.Note))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

// Ledger 返回本分站池流水。
func (h *SubSiteAdminHandler) Ledger(c *gin.Context) {
	ownerID, ok := h.ownerOrAbort(c)
	if !ok {
		return
	}
	siteID, ok := h.siteIDOrAbort(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.subSiteService.ListOwnedPoolLedger(c.Request.Context(), ownerID, siteID, params, c.Query("tx_type"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if items == nil {
		items = []service.SubSiteLedgerEntry{}
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}
