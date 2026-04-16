package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// GetPlans returns available payment plans
// GET /api/v1/payment/plans
func (h *PaymentHandler) GetPlans(c *gin.Context) {
	plans, err := h.paymentService.GetPlans(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, plans)
}

// GetRechargeInfo returns recharge plans and minimum amount
// GET /api/v1/payment/recharge-info
func (h *PaymentHandler) GetRechargeInfo(c *gin.Context) {
	info, err := h.paymentService.GetRechargeInfo(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, info)
}

// GetNewcomerStatus checks if the current user is eligible for newcomer plans
// GET /api/v1/payment/newcomer-status
func (h *PaymentHandler) GetNewcomerStatus(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	eligible, err := h.paymentService.CheckNewcomerEligibility(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{
		"eligible": eligible,
	})
}

// GetPayMethods returns currently enabled payment methods
// GET /api/v1/payment/methods
func (h *PaymentHandler) GetPayMethods(c *gin.Context) {
	methods := h.paymentService.GetAvailablePayMethods(c.Request.Context())
	response.Success(c, gin.H{
		"methods": methods,
	})
}

// CreateOrderRequest represents the request body for creating a payment order
type CreateOrderRequest struct {
	PlanKey   string `json:"plan_key" binding:"required"`
	PromoCode string `json:"promo_code"` // 优惠码
	PayMethod string `json:"pay_method"` // "wechat" | "alipay" | "epay_alipay" | "epay_wxpay"
}

// CreateRechargeRequest represents the request body for creating a recharge order
type CreateRechargeRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	PromoCode string  `json:"promo_code"` // 优惠码
	PayMethod string  `json:"pay_method"` // "wechat" | "alipay" | "epay_alipay" | "epay_wxpay"
	PlanKey   string  `json:"plan_key"`   // 充值套餐key（可选，用于验证和限购）
}

type CreateAgentActivationOrderRequest struct {
	PayMethod string `json:"pay_method"` // "wechat" | "alipay" | "epay_alipay" | "epay_wxpay"
}

type CreateSubSiteActivationOrderRequest struct {
	PayMethod       string                               `json:"pay_method"`
	ActivationInput service.CreateSubSiteActivationInput `json:"activation_input"`
}

type CreateSubSiteTopupOrderRequest struct {
	SiteID    int64  `json:"site_id" binding:"required"`
	AmountFen int    `json:"amount_fen" binding:"required,gt=0"`
	PayMethod string `json:"pay_method"`
}

type SubmitInvoiceRequest struct {
	OrderNos    []string `json:"order_nos" binding:"required,min=1"`
	CompanyName string   `json:"company_name" binding:"required"`
	TaxID       string   `json:"tax_id" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Remark      string   `json:"remark"`
}

// CreateOrder creates a new payment order
// POST /api/v1/payment/orders
func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: plan_key is required")
		return
	}

	// 默认使用微信支付
	payMethod := req.PayMethod
	if payMethod == "" {
		payMethod = "wechat"
	}

	order, err := h.paymentService.CreateOrder(c.Request.Context(), subject.UserID, req.PlanKey, req.PromoCode, payMethod)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order_no":        order.OrderNo,
		"code_url":        order.CodeURL,
		"amount_fen":      order.AmountFen,
		"discount_amount": order.DiscountAmount,
		"expired_at":      order.ExpiredAt,
	})
}

// CreateRecharge creates a new balance recharge order with custom amount
// POST /api/v1/payment/recharge
func (h *PaymentHandler) CreateRecharge(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req CreateRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: amount is required and must be positive")
		return
	}

	// 默认使用微信支付
	payMethod := req.PayMethod
	if payMethod == "" {
		payMethod = "wechat"
	}

	order, err := h.paymentService.CreateRechargeOrder(c.Request.Context(), subject.UserID, req.Amount, req.PromoCode, payMethod, req.PlanKey)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order_no":        order.OrderNo,
		"code_url":        order.CodeURL,
		"amount_fen":      order.AmountFen,
		"discount_amount": order.DiscountAmount,
		"expired_at":      order.ExpiredAt,
	})
}

// CreateAgentActivationOrder creates an order for the agent activation fee.
// POST /api/v1/payment/agent-activation
func (h *PaymentHandler) CreateAgentActivationOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req CreateAgentActivationOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		response.BadRequest(c, "invalid request")
		return
	}

	payMethod := req.PayMethod
	if payMethod == "" {
		payMethod = "wechat"
	}

	order, err := h.paymentService.CreateAgentActivationOrder(c.Request.Context(), subject.UserID, payMethod)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order_no":   order.OrderNo,
		"code_url":   order.CodeURL,
		"amount_fen": order.AmountFen,
		"expired_at": order.ExpiredAt,
	})
}

// CreateSubSiteActivationOrder creates an order for self-service sub-site activation.
// POST /api/v1/payment/subsite-activation
func (h *PaymentHandler) CreateSubSiteActivationOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req CreateSubSiteActivationOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	payMethod := req.PayMethod
	if payMethod == "" {
		payMethod = "wechat"
	}

	order, err := h.paymentService.CreateSubSiteActivationOrder(c.Request.Context(), subject.UserID, req.ActivationInput, payMethod)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order_no":   order.OrderNo,
		"code_url":   order.CodeURL,
		"amount_fen": order.AmountFen,
		"expired_at": order.ExpiredAt,
	})
}

// CreateSubSiteTopupOrder creates an order for sub-site pool online topup.
// POST /api/v1/payment/subsite-topup
func (h *PaymentHandler) CreateSubSiteTopupOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req CreateSubSiteTopupOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: site_id and amount_fen are required")
		return
	}

	payMethod := req.PayMethod
	if payMethod == "" {
		payMethod = "wechat"
	}

	order, err := h.paymentService.CreateSubSiteTopupOrder(c.Request.Context(), subject.UserID, req.SiteID, req.AmountFen, payMethod)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order_no":   order.OrderNo,
		"code_url":   order.CodeURL,
		"amount_fen": order.AmountFen,
		"expired_at": order.ExpiredAt,
	})
}

// QueryOrder queries the status of a payment order
// GET /api/v1/payment/orders/:orderNo
func (h *PaymentHandler) QueryOrder(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	orderNo := c.Param("orderNo")
	if orderNo == "" {
		response.BadRequest(c, "order_no is required")
		return
	}

	order, err := h.paymentService.QueryOrder(c.Request.Context(), subject.UserID, orderNo)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"order_no":             order.OrderNo,
		"plan_key":             order.PlanKey,
		"amount_fen":           order.AmountFen,
		"discount_amount":      order.DiscountAmount,
		"status":               order.Status,
		"pay_method":           order.PayMethod,
		"code_url":             order.CodeURL,
		"paid_at":              order.PaidAt,
		"expired_at":           order.ExpiredAt,
		"created_at":           order.CreatedAt,
		"invoice_company_name": order.InvoiceCompanyName,
		"invoice_tax_id":       order.InvoiceTaxID,
		"invoice_email":        order.InvoiceEmail,
		"invoice_remark":       order.InvoiceRemark,
		"invoice_requested_at": order.InvoiceRequestedAt,
		"invoice_processed_at": order.InvoiceProcessedAt,
	})
}

// ListOrders lists payment orders for the current user
// GET /api/v1/payment/orders
func (h *PaymentHandler) ListOrders(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	page, pageSize := response.ParsePagination(c)

	orders, paginationResult, err := h.paymentService.ListOrders(c.Request.Context(), subject.UserID, pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Convert to response DTOs
	items := make([]gin.H, len(orders))
	for i, order := range orders {
		items[i] = gin.H{
			"order_no":             order.OrderNo,
			"plan_key":             order.PlanKey,
			"amount_fen":           order.AmountFen,
			"discount_amount":      order.DiscountAmount,
			"status":               order.Status,
			"pay_method":           order.PayMethod,
			"paid_at":              order.PaidAt,
			"expired_at":           order.ExpiredAt,
			"created_at":           order.CreatedAt,
			"invoice_company_name": order.InvoiceCompanyName,
			"invoice_tax_id":       order.InvoiceTaxID,
			"invoice_email":        order.InvoiceEmail,
			"invoice_remark":       order.InvoiceRemark,
			"invoice_requested_at": order.InvoiceRequestedAt,
			"invoice_processed_at": order.InvoiceProcessedAt,
		}
	}

	response.Paginated(c, items, paginationResult.Total, paginationResult.Page, paginationResult.PageSize)
}

// SubmitInvoice submits invoice requests for one or more paid orders.
// POST /api/v1/payment/invoice-requests
func (h *PaymentHandler) SubmitInvoice(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req SubmitInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid invoice request")
		return
	}

	err := h.paymentService.SubmitInvoiceRequest(c.Request.Context(), subject.UserID, req.OrderNos, service.InvoiceRequest{
		CompanyName: req.CompanyName,
		TaxID:       req.TaxID,
		Email:       req.Email,
		Remark:      req.Remark,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
}

// WechatNotify handles WeChat Pay callback notifications
// POST /api/v1/payment/wechat/notify
func (h *PaymentHandler) WechatNotify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FAIL",
			"message": "read body failed",
		})
		return
	}

	// Extract WeChat Pay headers
	wechatpayTimestamp := c.GetHeader("Wechatpay-Timestamp")
	wechatpayNonce := c.GetHeader("Wechatpay-Nonce")
	wechatpaySignature := c.GetHeader("Wechatpay-Signature")
	wechatpaySerial := c.GetHeader("Wechatpay-Serial")

	log.Printf("[Payment] WechatNotify headers: Timestamp=%q Nonce=%q Serial=%q Signature_len=%d Body_len=%d",
		wechatpayTimestamp, wechatpayNonce, wechatpaySerial, len(wechatpaySignature), len(body))

	if err := h.paymentService.HandleWechatNotify(
		c.Request.Context(),
		body,
		wechatpayTimestamp,
		wechatpayNonce,
		wechatpaySignature,
		wechatpaySerial,
	); err != nil {
		log.Printf("[Payment] WechatNotify error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FAIL",
			"message": err.Error(),
		})
		return
	}

	// WeChat expects this exact response format for success
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "OK",
	})
}

// AlipayNotify handles Alipay callback notifications
// POST /api/v1/payment/alipay/notify
func (h *PaymentHandler) AlipayNotify(c *gin.Context) {
	if err := h.paymentService.HandleAlipayNotify(c.Request.Context(), c.Request); err != nil {
		log.Printf("[Payment] AlipayNotify error: %v", err)
		c.String(http.StatusInternalServerError, "failure")
		return
	}

	// 支付宝要求返回纯文本 "success"
	c.String(http.StatusOK, "success")
}

// EpayNotify handles Epay (易支付) callback notifications
// POST /api/v1/payment/epay/notify
func (h *PaymentHandler) EpayNotify(c *gin.Context) {
	// Epay sends params via GET query or POST form
	params := make(map[string]string)
	// Try query params first, then form params
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	if c.Request.Method == "POST" {
		if err := c.Request.ParseForm(); err == nil {
			for key, values := range c.Request.PostForm {
				if len(values) > 0 {
					params[key] = values[0]
				}
			}
		}
	}

	log.Printf("[Payment] EpayNotify %s params: %v", c.Request.Method, params)

	if err := h.paymentService.HandleEpayNotify(c.Request.Context(), params); err != nil {
		log.Printf("[Payment] EpayNotify error: %v", err)
		c.String(http.StatusInternalServerError, "failure")
		return
	}

	c.String(http.StatusOK, "success")
}
