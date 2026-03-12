package repository

import (
	"context"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentorder"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type paymentOrderRepo struct {
	client *dbent.Client
}

func NewPaymentOrderRepo(client *dbent.Client) service.PaymentOrderRepository {
	return &paymentOrderRepo{client: client}
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
	builder := r.client.PaymentOrder.Update().
		Where(paymentorder.OrderNoEQ(orderNo)).
		SetStatus(status)

	if transactionID != nil {
		builder.SetWechatTransactionID(*transactionID)
	}
	if paidAt != nil {
		builder.SetPaidAt(*paidAt)
	}

	_, err := builder.Save(ctx)
	return err
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
		PromoCode:           e.PromoCode,
		DiscountAmount:      e.DiscountAmount,
		Status:              e.Status,
		PayMethod:           e.PayMethod,
		WechatTransactionID: e.WechatTransactionID,
		CodeURL:             e.CodeURL,
		PaidAt:              e.PaidAt,
		ExpiredAt:           e.ExpiredAt,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}
}
