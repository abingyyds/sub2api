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

// CreateOrderRequest represents the request body for creating a payment order
type CreateOrderRequest struct {
	PlanKey string `json:"plan_key" binding:"required"`
}

// CreateRechargeRequest represents the request body for creating a recharge order
type CreateRechargeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
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

	order, err := h.paymentService.CreateOrder(c.Request.Context(), subject.UserID, req.PlanKey)
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

	order, err := h.paymentService.CreateRechargeOrder(c.Request.Context(), subject.UserID, req.Amount)
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
		"order_no":   order.OrderNo,
		"plan_key":   order.PlanKey,
		"amount_fen": order.AmountFen,
		"status":     order.Status,
		"pay_method": order.PayMethod,
		"code_url":   order.CodeURL,
		"paid_at":    order.PaidAt,
		"expired_at": order.ExpiredAt,
		"created_at": order.CreatedAt,
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
			"order_no":   order.OrderNo,
			"plan_key":   order.PlanKey,
			"amount_fen": order.AmountFen,
			"status":     order.Status,
			"pay_method": order.PayMethod,
			"paid_at":    order.PaidAt,
			"expired_at": order.ExpiredAt,
			"created_at": order.CreatedAt,
		}
	}

	response.Paginated(c, items, paginationResult.Total, paginationResult.Page, paginationResult.PageSize)
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
