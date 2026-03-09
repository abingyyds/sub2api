package admin

import (
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
		response.ErrorFrom(c, err)
		return
	}

	items := make([]gin.H, len(orders))
	for i, order := range orders {
		items[i] = gin.H{
			"id":                     order.ID,
			"order_no":               order.OrderNo,
			"user_id":                order.UserID,
			"plan_key":               order.PlanKey,
			"group_id":               order.GroupID,
			"amount_fen":             order.AmountFen,
			"validity_days":          order.ValidityDays,
			"order_type":             order.OrderType,
			"balance_amount":         order.BalanceAmount,
			"status":                 order.Status,
			"pay_method":             order.PayMethod,
			"wechat_transaction_id":  order.WechatTransactionID,
			"paid_at":                order.PaidAt,
			"expired_at":             order.ExpiredAt,
			"created_at":             order.CreatedAt,
		}
	}

	response.Paginated(c, items, pag.Total, pag.Page, pag.PageSize)
}
