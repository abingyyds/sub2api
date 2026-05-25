package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newPaymentOrderEntRepo(t *testing.T) (*paymentOrderRepo, *dbent.Client) {
	t.Helper()

	db, err := sql.Open("sqlite", "file:payment_orders?mode=memory&cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() })

	return &paymentOrderRepo{client: client, db: db}, client
}

func TestPaymentOrderRepoListAllFiltersPendingInvoices(t *testing.T) {
	repo, client := newPaymentOrderEntRepo(t)
	ctx := context.Background()

	user := client.User.Create().
		SetEmail("invoice-admin-test@example.com").
		SetPasswordHash("hash").
		SaveX(ctx)

	now := time.Date(2026, 5, 25, 8, 0, 0, 0, time.UTC)
	oldRequest := now.Add(-2 * time.Hour)
	newRequest := now.Add(-30 * time.Minute)

	createOrder := func(orderNo string) *dbent.PaymentOrderCreate {
		return client.PaymentOrder.Create().
			SetOrderNo(orderNo).
			SetUserID(user.ID).
			SetPlanKey("basic").
			SetGroupID(1).
			SetAmountFen(1000).
			SetValidityDays(30).
			SetStatus(service.PaymentOrderStatusPaid).
			SetPayMethod("wechat_native").
			SetExpiredAt(now.Add(24 * time.Hour)).
			SetCreatedAt(now.Add(-4 * time.Hour))
	}

	createOrder("no-invoice").SaveX(ctx)
	createOrder("processed-invoice").
		SetInvoiceRequestedAt(oldRequest).
		SetInvoiceProcessedAt(now.Add(-time.Hour)).
		SaveX(ctx)
	createOrder("pending-old").
		SetInvoiceRequestedAt(oldRequest).
		SaveX(ctx)
	createOrder("pending-new").
		SetInvoiceRequestedAt(newRequest).
		SaveX(ctx)

	orders, pag, err := repo.ListAll(ctx, pagination.PaginationParams{Page: 1, PageSize: 10}, "", "", "pending")
	require.NoError(t, err)
	require.Equal(t, int64(2), pag.Total)
	require.Len(t, orders, 2)
	require.Equal(t, "pending-new", orders[0].OrderNo)
	require.Equal(t, "pending-old", orders[1].OrderNo)
	require.NotNil(t, orders[0].InvoiceRequestedAt)
	require.Nil(t, orders[0].InvoiceProcessedAt)
}
