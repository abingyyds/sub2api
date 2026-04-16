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

// SubSiteHandler handles user-facing sub-site requests.
type SubSiteHandler struct {
	subSiteService *service.SubSiteService
}

func NewSubSiteHandler(subSiteService *service.SubSiteService) *SubSiteHandler {
	return &SubSiteHandler{subSiteService: subSiteService}
}

// GetOpenInfo returns self-service opening information for the current site scope.
func (h *SubSiteHandler) GetOpenInfo(c *gin.Context) {
	info, err := h.subSiteService.GetOpenInfo(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, info)
}

// ListOwned returns all sub-sites owned by the current user.
func (h *SubSiteHandler) ListOwned(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	items, err := h.subSiteService.ListOwnedSites(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

// GetOwned returns one owned sub-site detail.
func (h *SubSiteHandler) GetOwned(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	siteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	item, err := h.subSiteService.GetOwnedSite(c.Request.Context(), subject.UserID, siteID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

// UpdateOwned updates a sub-site owned by the current user.
func (h *SubSiteHandler) UpdateOwned(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	siteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req service.UpdateOwnedSubSiteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = siteID
	item, err := h.subSiteService.UpdateOwnedSite(c.Request.Context(), subject.UserID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

type ownedSubSiteOfflineTopupRequest struct {
	UserID    int64  `json:"user_id"`
	AmountFen int64  `json:"amount_fen"`
	Note      string `json:"note"`
}

// OfflineTopupOwnedUser 分站主给自己分站的用户线下加余额；从分站池同额扣除。
func (h *SubSiteHandler) OfflineTopupOwnedUser(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	siteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	var req ownedSubSiteOfflineTopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if req.AmountFen <= 0 {
		response.BadRequest(c, "amount_fen must be positive")
		return
	}
	site, err := h.subSiteService.OfflineTopupUser(c.Request.Context(), subject.UserID, siteID, req.UserID, req.AmountFen, strings.TrimSpace(req.Note))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, site)
}

// OwnedLedger 分站主查询自己分站的流水。
func (h *SubSiteHandler) OwnedLedger(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	siteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || siteID <= 0 {
		response.BadRequest(c, "invalid sub-site ID")
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	items, pag, err := h.subSiteService.ListOwnedPoolLedger(c.Request.Context(), subject.UserID, siteID, params, c.Query("tx_type"))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if items == nil {
		items = []service.SubSiteLedgerEntry{}
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}
