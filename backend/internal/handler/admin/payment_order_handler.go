package admin

import (
	"log/slog"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// PaymentOrderHandler handles admin payment order management
type PaymentOrderHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentOrderHandler creates a new admin payment order handler
func NewPaymentOrderHandler(paymentService *service.PaymentService) *PaymentOrderHandler {
	return &PaymentOrderHandler{paymentService: paymentService}
}

// List handles listing all payment orders with pagination and filters
// GET /api/v1/admin/orders
func (h *PaymentOrderHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	status := c.Query("status")
	orderType := c.Query("order_type")

	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	orders, pag, err := h.paymentService.ListAllOrders(c.Request.Context(), params, status, orderType)
	if err != nil {
		slog.Error("[AdminOrders] ListAllOrders failed", "error", err)
		response.ErrorFrom(c, err)
		return
	}

	items := make([]gin.H, len(orders))
	for i, order := range orders {
		items[i] = gin.H{
			"id":                    order.ID,
			"order_no":              order.OrderNo,
			"user_id":               order.UserID,
			"plan_key":              order.PlanKey,
			"group_id":              order.GroupID,
			"amount_fen":            order.AmountFen,
			"validity_days":         order.ValidityDays,
			"order_type":            order.OrderType,
			"balance_amount":        order.BalanceAmount,
			"sub_site_id":           order.SubSiteID,
			"status":                order.Status,
			"pay_method":            order.PayMethod,
			"wechat_transaction_id": order.WechatTransactionID,
			"alipay_trade_no":       order.AlipayTradeNo,
			"epay_trade_no":         order.EpayTradeNo,
			"invoice_company_name":  order.InvoiceCompanyName,
			"invoice_tax_id":        order.InvoiceTaxID,
			"invoice_email":         order.InvoiceEmail,
			"invoice_remark":        order.InvoiceRemark,
			"invoice_requested_at":  order.InvoiceRequestedAt,
			"invoice_processed_at":  order.InvoiceProcessedAt,
			"paid_at":               order.PaidAt,
			"expired_at":            order.ExpiredAt,
			"created_at":            order.CreatedAt,
		}
	}

	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}

// MarkInvoiceProcessed marks an order invoice request as processed.
// POST /api/v1/admin/orders/:id/invoice/processed
func (h *PaymentOrderHandler) MarkInvoiceProcessed(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID <= 0 {
		response.BadRequest(c, "invalid order id")
		return
	}

	if err := h.paymentService.MarkInvoiceProcessed(c.Request.Context(), orderID); err != nil {
		slog.Error("[AdminOrders] MarkInvoiceProcessed failed", "order_id", orderID, "error", err)
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
}

// Repair manually fulfills an order when the payment callback was missed.
// POST /api/v1/admin/orders/:id/repair
func (h *PaymentOrderHandler) Repair(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID <= 0 {
		response.BadRequest(c, "invalid order id")
		return
	}

	if err := h.paymentService.RepairOrder(c.Request.Context(), orderID); err != nil {
		slog.Error("[AdminOrders] RepairOrder failed", "order_id", orderID, "error", err)
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
}
