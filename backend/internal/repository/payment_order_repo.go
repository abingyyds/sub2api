package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentorder"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type paymentOrderRepo struct {
	client *dbent.Client
	db     *sql.DB
}

func NewPaymentOrderRepo(client *dbent.Client, db *sql.DB) service.PaymentOrderRepository {
	return &paymentOrderRepo{client: client, db: db}
}

func (r *paymentOrderRepo) Create(ctx context.Context, order *service.PaymentOrder) error {
	builder := r.client.PaymentOrder.Create().
		SetOrderNo(order.OrderNo).
		SetUserID(order.UserID).
		SetPlanKey(order.PlanKey).
		SetGroupID(order.GroupID).
		SetAmountFen(order.AmountFen).
		SetValidityDays(order.ValidityDays).
		SetOrderType(order.OrderType).
		SetBalanceAmount(order.BalanceAmount).
		SetPromoCode(order.PromoCode).
		SetDiscountAmount(order.DiscountAmount).
		SetStatus(order.Status).
		SetPayMethod(order.PayMethod).
		SetExpiredAt(order.ExpiredAt)

	if order.CodeURL != nil {
		builder.SetCodeURL(*order.CodeURL)
	}
	if order.SubSiteID != nil {
		builder.SetSubSiteID(*order.SubSiteID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return fmt.Errorf("create payment order: %w", err)
	}
	order.ID = created.ID
	order.CreatedAt = created.CreatedAt
	order.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *paymentOrderRepo) GetByID(ctx context.Context, id int64) (*service.PaymentOrder, error) {
	order, err := r.client.PaymentOrder.Get(ctx, id)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrPaymentOrderNotFound
		}
		return nil, fmt.Errorf("get payment order: %w", err)
	}
	return toServicePaymentOrder(order), nil
}

func (r *paymentOrderRepo) GetByOrderNo(ctx context.Context, orderNo string) (*service.PaymentOrder, error) {
	order, err := r.client.PaymentOrder.Query().
		Where(paymentorder.OrderNoEQ(orderNo)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrPaymentOrderNotFound
		}
		return nil, fmt.Errorf("get payment order by order_no: %w", err)
	}
	return toServicePaymentOrder(order), nil
}

func (r *paymentOrderRepo) UpdateStatus(ctx context.Context, orderNo string, status string, transactionID *string, paidAt *time.Time) error {
	// 先查询订单以确定支付方式
	order, err := r.client.PaymentOrder.Query().
		Where(paymentorder.OrderNoEQ(orderNo)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("get order for update: %w", err)
	}

	builder := r.client.PaymentOrder.Update().
		Where(paymentorder.OrderNoEQ(orderNo)).
		SetStatus(status)

	if transactionID != nil {
		// 根据支付方式设置不同的交易号字段
		switch order.PayMethod {
		case "alipay_native":
			builder.SetAlipayTradeNo(*transactionID)
		case "epay_alipay", "epay_wxpay":
			builder.SetEpayTradeNo(*transactionID)
		default:
			builder.SetWechatTransactionID(*transactionID)
		}
	}
	if paidAt != nil {
		builder.SetPaidAt(*paidAt)
	}

	_, err = builder.Save(ctx)
	return err
}

func (r *paymentOrderRepo) CompareAndUpdateStatus(ctx context.Context, orderNo string, expectedStatus string, newStatus string, transactionID *string, paidAt *time.Time) (bool, error) {
	// 先查询订单以确定支付方式
	order, err := r.client.PaymentOrder.Query().
		Where(paymentorder.OrderNoEQ(orderNo)).
		Only(ctx)
	if err != nil {
		return false, fmt.Errorf("get order for CAS update: %w", err)
	}

	// CAS: 只在当前状态匹配时更新
	builder := r.client.PaymentOrder.Update().
		Where(
			paymentorder.OrderNoEQ(orderNo),
			paymentorder.StatusEQ(expectedStatus),
		).
		SetStatus(newStatus)

	if transactionID != nil {
		switch order.PayMethod {
		case "alipay_native":
			builder.SetAlipayTradeNo(*transactionID)
		case "epay_alipay", "epay_wxpay":
			builder.SetEpayTradeNo(*transactionID)
		default:
			builder.SetWechatTransactionID(*transactionID)
		}
	}
	if paidAt != nil {
		builder.SetPaidAt(*paidAt)
	}

	n, err := builder.Save(ctx)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *paymentOrderRepo) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.PaymentOrder, *pagination.PaginationResult, error) {
	query := r.client.PaymentOrder.Query().
		Where(paymentorder.UserIDEQ(userID))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("count payment orders: %w", err)
	}

	orders, err := query.
		Order(dbent.Desc(paymentorder.FieldCreatedAt)).
		Limit(params.PageSize).
		Offset((params.Page - 1) * params.PageSize).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list payment orders: %w", err)
	}

	result := make([]service.PaymentOrder, len(orders))
	for i, o := range orders {
		result[i] = *toServicePaymentOrder(o)
	}

	return result, &pagination.PaginationResult{
		Total:    int64(total),
		Page:     params.Page,
		PageSize: params.PageSize,
		Pages:    (total + params.PageSize - 1) / params.PageSize,
	}, nil
}

func (r *paymentOrderRepo) ListAll(ctx context.Context, params pagination.PaginationParams, status string, orderType string) ([]service.PaymentOrder, *pagination.PaginationResult, error) {
	query := r.client.PaymentOrder.Query()

	if status != "" {
		query = query.Where(paymentorder.StatusEQ(status))
	}
	if orderType != "" {
		query = query.Where(paymentorder.OrderTypeEQ(orderType))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("count payment orders: %w", err)
	}

	orders, err := query.
		Order(dbent.Desc(paymentorder.FieldCreatedAt)).
		Limit(params.PageSize).
		Offset((params.Page - 1) * params.PageSize).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list all payment orders: %w", err)
	}

	result := make([]service.PaymentOrder, len(orders))
	for i, o := range orders {
		result[i] = *toServicePaymentOrder(o)
	}

	return result, &pagination.PaginationResult{
		Total:    int64(total),
		Page:     params.Page,
		PageSize: params.PageSize,
		Pages:    (total + params.PageSize - 1) / params.PageSize,
	}, nil
}

func (r *paymentOrderRepo) SubmitInvoiceRequest(ctx context.Context, userID int64, orderNos []string, invoice service.InvoiceRequest) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("start invoice request transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	orders, err := tx.PaymentOrder.Query().
		Where(
			paymentorder.UserIDEQ(userID),
			paymentorder.OrderNoIn(orderNos...),
		).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query invoice orders: %w", err)
	}
	if len(orders) != len(orderNos) {
		return service.ErrPaymentOrderNotFound
	}

	for _, order := range orders {
		if order.Status != service.PaymentOrderStatusPaid {
			return service.ErrInvoiceOrderNotPaid
		}
		if order.InvoiceRequestedAt != nil {
			return service.ErrInvoiceAlreadyFiled
		}
	}

	now := time.Now()
	updated, err := tx.PaymentOrder.Update().
		Where(
			paymentorder.UserIDEQ(userID),
			paymentorder.OrderNoIn(orderNos...),
			paymentorder.InvoiceRequestedAtIsNil(),
		).
		SetInvoiceCompanyName(invoice.CompanyName).
		SetInvoiceTaxID(invoice.TaxID).
		SetInvoiceEmail(invoice.Email).
		SetInvoiceRemark(invoice.Remark).
		SetInvoiceRequestedAt(now).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("update invoice request: %w", err)
	}
	if updated != len(orderNos) {
		return service.ErrInvoiceAlreadyFiled
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit invoice request transaction: %w", err)
	}
	return nil
}

func (r *paymentOrderRepo) MarkInvoiceProcessed(ctx context.Context, orderID int64) error {
	order, err := r.client.PaymentOrder.Get(ctx, orderID)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrPaymentOrderNotFound
		}
		return fmt.Errorf("get order before mark invoice processed: %w", err)
	}
	if order.InvoiceRequestedAt == nil {
		return service.ErrInvoiceNotRequested
	}
	if order.InvoiceProcessedAt != nil {
		return service.ErrInvoiceAlreadyHandled
	}

	updated, err := r.client.PaymentOrder.Update().
		Where(
			paymentorder.IDEQ(orderID),
			paymentorder.InvoiceRequestedAtNotNil(),
			paymentorder.InvoiceProcessedAtIsNil(),
		).
		SetInvoiceProcessedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("mark invoice processed: %w", err)
	}
	if updated == 0 {
		return service.ErrInvoiceAlreadyHandled
	}
	return nil
}

func (r *paymentOrderRepo) CloseExpiredOrders(ctx context.Context) (int64, error) {
	affected, err := r.client.PaymentOrder.Update().
		Where(
			paymentorder.StatusEQ(service.PaymentOrderStatusPending),
			paymentorder.ExpiredAtLT(time.Now()),
		).
		SetStatus(service.PaymentOrderStatusClosed).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("close expired orders: %w", err)
	}
	return int64(affected), nil
}

func (r *paymentOrderRepo) CountPaidByUserAndPlanKey(ctx context.Context, userID int64, planKey string) (int, error) {
	count, err := r.client.PaymentOrder.Query().
		Where(
			paymentorder.UserIDEQ(userID),
			paymentorder.PlanKeyEQ(planKey),
			paymentorder.StatusEQ(service.PaymentOrderStatusPaid),
		).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count paid orders by user and plan_key: %w", err)
	}
	return count, nil
}

func toServicePaymentOrder(e *dbent.PaymentOrder) *service.PaymentOrder {
	return &service.PaymentOrder{
		ID:                  e.ID,
		OrderNo:             e.OrderNo,
		UserID:              e.UserID,
		PlanKey:             e.PlanKey,
		GroupID:             e.GroupID,
		AmountFen:           e.AmountFen,
		ValidityDays:        e.ValidityDays,
		OrderType:           e.OrderType,
		BalanceAmount:       e.BalanceAmount,
		SubSiteID:           e.SubSiteID,
		PromoCode:           e.PromoCode,
		DiscountAmount:      e.DiscountAmount,
		Status:              e.Status,
		PayMethod:           e.PayMethod,
		WechatTransactionID: e.WechatTransactionID,
		AlipayTradeNo:       e.AlipayTradeNo,
		EpayTradeNo:         e.EpayTradeNo,
		InvoiceCompanyName:  e.InvoiceCompanyName,
		InvoiceTaxID:        e.InvoiceTaxID,
		InvoiceEmail:        e.InvoiceEmail,
		InvoiceRemark:       e.InvoiceRemark,
		InvoiceRequestedAt:  e.InvoiceRequestedAt,
		InvoiceProcessedAt:  e.InvoiceProcessedAt,
		CodeURL:             e.CodeURL,
		PaidAt:              e.PaidAt,
		ExpiredAt:           e.ExpiredAt,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}
}
