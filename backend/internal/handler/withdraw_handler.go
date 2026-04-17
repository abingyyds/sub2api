package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type WithdrawHandler struct {
	withdrawService *service.WithdrawService
}

func NewWithdrawHandler(withdrawService *service.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{withdrawService: withdrawService}
}

type applyWithdrawRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	AlipayName      string  `json:"alipay_name" binding:"required"`
	AlipayAccount   string  `json:"alipay_account" binding:"required"`
	AlipayQRImage   string  `json:"alipay_qr_image"`
	SourceType      string  `json:"source_type" binding:"required"`
	SourceSubSiteID *int64  `json:"source_sub_site_id"`
}

func (h *WithdrawHandler) Apply(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}
	var req applyWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	result, err := h.withdrawService.ApplyWithdraw(c.Request.Context(), subject.UserID, req.Amount,
		req.AlipayName, req.AlipayAccount, req.AlipayQRImage, req.SourceType, req.SourceSubSiteID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *WithdrawHandler) Cancel(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	result, err := h.withdrawService.CancelWithdraw(c.Request.Context(), subject.UserID, id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *WithdrawHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	sourceType := c.Query("source_type")
	status := c.Query("status")
	items, pag, err := h.withdrawService.ListMyWithdraws(c.Request.Context(), subject.UserID, params, sourceType, status)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if items == nil {
		items = []service.WithdrawRequest{}
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

func (h *WithdrawHandler) AdminList(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	sourceType := c.Query("source_type")
	status := c.Query("status")
	items, pag, err := h.withdrawService.AdminListWithdraws(c.Request.Context(), params, sourceType, status)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if items == nil {
		items = []service.WithdrawRequest{}
	}
	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

type reviewWithdrawRequest struct {
	Approve    bool   `json:"approve"`
	ReviewNote string `json:"review_note"`
}

func (h *WithdrawHandler) AdminReview(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req reviewWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	result, err := h.withdrawService.ReviewWithdraw(c.Request.Context(), id, req.Approve, req.ReviewNote)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *WithdrawHandler) AdminPay(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	result, err := h.withdrawService.PayWithdraw(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}
