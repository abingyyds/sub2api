//go:build unit

package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type paymentOrderRepoStub struct {
	order             *PaymentOrder
	compareCalls      int
	updateStatusCalls int
}

func (s *paymentOrderRepoStub) Create(ctx context.Context, order *PaymentOrder) error {
	panic("unexpected Create call")
}

func (s *paymentOrderRepoStub) GetByID(ctx context.Context, id int64) (*PaymentOrder, error) {
	if s.order == nil || s.order.ID != id {
		return nil, ErrPaymentOrderNotFound
	}
	clone := *s.order
	return &clone, nil
}

func (s *paymentOrderRepoStub) GetByOrderNo(ctx context.Context, orderNo string) (*PaymentOrder, error) {
	if s.order == nil || s.order.OrderNo != orderNo {
		return nil, ErrPaymentOrderNotFound
	}
	clone := *s.order
	return &clone, nil
}

func (s *paymentOrderRepoStub) UpdateStatus(ctx context.Context, orderNo string, status string, transactionID *string, paidAt *time.Time) error {
	s.updateStatusCalls++
	if s.order != nil && s.order.OrderNo == orderNo {
		s.order.Status = status
		s.order.PaidAt = paidAt
	}
	return nil
}

func (s *paymentOrderRepoStub) CompareAndUpdateStatus(ctx context.Context, orderNo string, expectedStatus string, newStatus string, transactionID *string, paidAt *time.Time) (bool, error) {
	s.compareCalls++
	if s.order == nil || s.order.OrderNo != orderNo {
		return false, ErrPaymentOrderNotFound
	}
	if s.order.Status != expectedStatus {
		return false, nil
	}
	s.order.Status = newStatus
	s.order.PaidAt = paidAt
	return true, nil
}

func (s *paymentOrderRepoStub) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]PaymentOrder, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserID call")
}

func (s *paymentOrderRepoStub) ListAll(ctx context.Context, params pagination.PaginationParams, status string, orderType string) ([]PaymentOrder, *pagination.PaginationResult, error) {
	panic("unexpected ListAll call")
}

func (s *paymentOrderRepoStub) SubmitInvoiceRequest(ctx context.Context, userID int64, orderNos []string, invoice InvoiceRequest) error {
	panic("unexpected SubmitInvoiceRequest call")
}

func (s *paymentOrderRepoStub) MarkInvoiceProcessed(ctx context.Context, orderID int64) error {
	panic("unexpected MarkInvoiceProcessed call")
}

func (s *paymentOrderRepoStub) CloseExpiredOrders(ctx context.Context) (int64, error) {
	panic("unexpected CloseExpiredOrders call")
}

func (s *paymentOrderRepoStub) CountPaidByUserAndPlanKey(ctx context.Context, userID int64, planKey string) (int, error) {
	panic("unexpected CountPaidByUserAndPlanKey call")
}

type paymentBalanceUserRepoStub struct {
	*userRepoStub
	balanceUpdates []float64
}

func (s *paymentBalanceUserRepoStub) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	s.balanceUpdates = append(s.balanceUpdates, amount)
	if s.userRepoStub != nil && s.userRepoStub.user != nil {
		s.userRepoStub.user.Balance += amount
	}
	return nil
}

func newPaymentServiceForTest(settings map[string]string) *PaymentService {
	cfg := &config.Config{}
	settingService := NewSettingService(&settingRepoStub{values: settings}, cfg)
	return &PaymentService{settingService: settingService}
}

func TestPaymentServiceGetAvailablePayMethodsFiltersIncompleteConfigs(t *testing.T) {
	t.Parallel()

	svc := newPaymentServiceForTest(map[string]string{
		SettingKeyPaymentEnabled:       "true",
		SettingKeyWechatPayAppID:       "wx-app",
		SettingKeyWechatPayMchID:       "wx-mch",
		SettingKeyWechatPayNotifyURL:   "https://example.com/wx-notify",
		SettingKeyWechatPayPrivateKey:  "private-key",
		SettingKeyWechatPayMchSerialNo: "serial",
		SettingKeyAlipayEnabled:        "true",
		SettingKeyAlipayAppID:          "ali-app",
		SettingKeyAlipayPrivateKey:     "ali-private",
		SettingKeyAlipayPublicKey:      "ali-public",
		SettingKeyEpayEnabled:          "true",
		SettingKeyEpayGateway:          "https://epay.example.com",
		SettingKeyEpayPID:              "10001",
		SettingKeyEpayPKey:             "secret",
		SettingKeyEpayNotifyURL:        "https://example.com/epay-notify",
	})

	require.Equal(t, []string{"wechat", "epay_alipay", "epay_wxpay"}, svc.GetAvailablePayMethods(context.Background()))
}

func TestPaymentServiceCreateEpayOrderEncodesFormValues(t *testing.T) {
	t.Parallel()

	var receivedForm neturl.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/mapi.php", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		receivedForm, err = neturl.ParseQuery(string(body))
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"qrcode":"https://pay.example.com/qr"}`))
	}))
	defer server.Close()

	svc := newPaymentServiceForTest(map[string]string{
		SettingKeyEpayGateway:   server.URL,
		SettingKeyEpayPID:       "10001",
		SettingKeyEpayPKey:      "secret",
		SettingKeyEpayNotifyURL: "https://example.com/payment/notify?source=pricing",
	})

	order := &PaymentOrder{
		OrderNo:   "PO202604080001",
		AmountFen: 9900,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}

	codeURL, err := svc.createEpayOrder(context.Background(), order, "套餐充值 Pro Max", "alipay")
	require.NoError(t, err)
	require.Equal(t, "https://pay.example.com/qr", codeURL)
	require.Equal(t, "套餐充值 Pro Max", receivedForm.Get("name"))
	require.Equal(t, "https://example.com/payment/notify?source=pricing", receivedForm.Get("notify_url"))
	require.Equal(t, "99.00", receivedForm.Get("money"))
	require.NotEmpty(t, receivedForm.Get("sign"))
}

func TestPaymentServiceCreateEpayOrderUnexpectedResponseReturnsBadRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html>gateway error</html>"))
	}))
	defer server.Close()

	svc := newPaymentServiceForTest(map[string]string{
		SettingKeyEpayGateway:   server.URL,
		SettingKeyEpayPID:       "10001",
		SettingKeyEpayPKey:      "secret",
		SettingKeyEpayNotifyURL: "https://example.com/payment/notify",
	})

	order := &PaymentOrder{
		OrderNo:   "PO202604080002",
		AmountFen: 19900,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}

	_, err := svc.createEpayOrder(context.Background(), order, "测试订单", "wxpay")
	require.Error(t, err)
	require.Equal(t, http.StatusBadRequest, infraerrors.Code(err))
	require.Equal(t, "EPAY_API_ERROR", infraerrors.Reason(err))
}

func TestPaymentServiceRepairOrder_ReplaysBalanceFulfillment(t *testing.T) {
	t.Parallel()

	orderRepo := &paymentOrderRepoStub{
		order: &PaymentOrder{
			ID:            1,
			OrderNo:       "PO202604080003",
			UserID:        7,
			OrderType:     "balance",
			BalanceAmount: 25,
			Status:        PaymentOrderStatusClosed,
		},
	}
	userRepo := &paymentBalanceUserRepoStub{userRepoStub: &userRepoStub{user: &User{ID: 7, Balance: 10}}}
	svc := &PaymentService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}

	err := svc.RepairOrder(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, PaymentOrderStatusPaid, orderRepo.order.Status)
	require.NotNil(t, orderRepo.order.PaidAt)
	require.Equal(t, []float64{25}, userRepo.balanceUpdates)
	require.Equal(t, 35.0, userRepo.user.Balance)
	require.Equal(t, 1, orderRepo.compareCalls)
	require.Zero(t, orderRepo.updateStatusCalls)
}

func TestPaymentServiceRepairOrder_RejectsPaidOrder(t *testing.T) {
	t.Parallel()

	orderRepo := &paymentOrderRepoStub{
		order: &PaymentOrder{
			ID:        2,
			OrderNo:   "PO202604080004",
			UserID:    7,
			OrderType: "balance",
			Status:    PaymentOrderStatusPaid,
		},
	}
	userRepo := &paymentBalanceUserRepoStub{userRepoStub: &userRepoStub{user: &User{ID: 7, Balance: 10}}}
	svc := &PaymentService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}

	err := svc.RepairOrder(context.Background(), 2)
	require.ErrorIs(t, err, ErrPaymentOrderPaid)
	require.Empty(t, userRepo.balanceUpdates)
	require.Zero(t, orderRepo.compareCalls)
}
