//go:build unit

package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

type paymentPlanGroupRepoStub struct {
	groups []Group
}

func (s *paymentPlanGroupRepoStub) Create(ctx context.Context, group *Group) error {
	panic("unexpected Create call")
}

func (s *paymentPlanGroupRepoStub) GetByID(ctx context.Context, id int64) (*Group, error) {
	panic("unexpected GetByID call")
}

func (s *paymentPlanGroupRepoStub) GetByIDLite(ctx context.Context, id int64) (*Group, error) {
	panic("unexpected GetByIDLite call")
}

func (s *paymentPlanGroupRepoStub) Update(ctx context.Context, group *Group) error {
	panic("unexpected Update call")
}

func (s *paymentPlanGroupRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *paymentPlanGroupRepoStub) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
}

func (s *paymentPlanGroupRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *paymentPlanGroupRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *paymentPlanGroupRepoStub) ListActive(ctx context.Context) ([]Group, error) {
	return s.groups, nil
}

func (s *paymentPlanGroupRepoStub) ListActiveByPlatform(ctx context.Context, platform string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
}

func (s *paymentPlanGroupRepoStub) ExistsByName(ctx context.Context, name string) (bool, error) {
	panic("unexpected ExistsByName call")
}

func (s *paymentPlanGroupRepoStub) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected GetAccountCount call")
}

func (s *paymentPlanGroupRepoStub) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
}

type paymentActiveSubscriptionRepoStub struct {
	active *UserSubscription
}

func (s *paymentActiveSubscriptionRepoStub) Create(ctx context.Context, sub *UserSubscription) error {
	panic("unexpected Create call")
}

func (s *paymentActiveSubscriptionRepoStub) GetByID(ctx context.Context, id int64) (*UserSubscription, error) {
	panic("unexpected GetByID call")
}

func (s *paymentActiveSubscriptionRepoStub) GetByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error) {
	panic("unexpected GetByUserIDAndGroupID call")
}

func (s *paymentActiveSubscriptionRepoStub) GetActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error) {
	if s.active == nil || s.active.UserID != userID || s.active.GroupID != groupID {
		return nil, ErrSubscriptionNotFound
	}
	clone := *s.active
	return &clone, nil
}

func (s *paymentActiveSubscriptionRepoStub) Update(ctx context.Context, sub *UserSubscription) error {
	panic("unexpected Update call")
}

func (s *paymentActiveSubscriptionRepoStub) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
}

func (s *paymentActiveSubscriptionRepoStub) ListByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	panic("unexpected ListByUserID call")
}

func (s *paymentActiveSubscriptionRepoStub) ListActiveByUserID(ctx context.Context, userID int64) ([]UserSubscription, error) {
	panic("unexpected ListActiveByUserID call")
}

func (s *paymentActiveSubscriptionRepoStub) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]UserSubscription, *pagination.PaginationResult, error) {
	panic("unexpected ListByGroupID call")
}

func (s *paymentActiveSubscriptionRepoStub) List(ctx context.Context, params pagination.PaginationParams, userID, groupID *int64, status, sortBy, sortOrder string) ([]UserSubscription, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *paymentActiveSubscriptionRepoStub) ExistsByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error) {
	panic("unexpected ExistsByUserIDAndGroupID call")
}

func (s *paymentActiveSubscriptionRepoStub) ExtendExpiry(ctx context.Context, subscriptionID int64, newExpiresAt time.Time) error {
	panic("unexpected ExtendExpiry call")
}

func (s *paymentActiveSubscriptionRepoStub) UpdateStatus(ctx context.Context, subscriptionID int64, status string) error {
	panic("unexpected UpdateStatus call")
}

func (s *paymentActiveSubscriptionRepoStub) UpdateNotes(ctx context.Context, subscriptionID int64, notes string) error {
	panic("unexpected UpdateNotes call")
}

func (s *paymentActiveSubscriptionRepoStub) ActivateWindows(ctx context.Context, id int64, start time.Time) error {
	panic("unexpected ActivateWindows call")
}

func (s *paymentActiveSubscriptionRepoStub) ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetDailyUsage call")
}

func (s *paymentActiveSubscriptionRepoStub) ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetWeeklyUsage call")
}

func (s *paymentActiveSubscriptionRepoStub) ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error {
	panic("unexpected ResetMonthlyUsage call")
}

func (s *paymentActiveSubscriptionRepoStub) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	panic("unexpected IncrementUsage call")
}

func (s *paymentActiveSubscriptionRepoStub) BatchUpdateExpiredStatus(ctx context.Context) (int64, error) {
	panic("unexpected BatchUpdateExpiredStatus call")
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

func TestPaymentServiceCreateOrderBlocksRepurchaseWhileSubscriptionActive(t *testing.T) {
	t.Parallel()

	groupRepo := &paymentPlanGroupRepoStub{
		groups: []Group{
			{
				ID:                  88,
				Name:                "Pro",
				Description:         "pro plan",
				PriceFen:            29900,
				Listed:              true,
				DefaultValidityDays: 30,
				SubscriptionType:    SubscriptionTypeSubscription,
				Status:              StatusActive,
			},
		},
	}
	subRepo := &paymentActiveSubscriptionRepoStub{
		active: &UserSubscription{
			ID:        10,
			UserID:    7,
			GroupID:   88,
			Status:    SubscriptionStatusActive,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	}

	svc := &PaymentService{
		settingService: NewSettingService(&settingRepoStub{values: map[string]string{
			SettingKeyPaymentEnabled: "true",
		}}, &config.Config{}),
		groupRepo: groupRepo,
		subscriptionService: &SubscriptionService{
			userSubRepo: subRepo,
		},
	}

	_, err := svc.CreateOrder(context.Background(), 7, "group_88", "", "")
	require.ErrorIs(t, err, ErrPaymentSubscriptionRepurchaseBlocked)
	require.Equal(t, http.StatusConflict, infraerrors.Code(err))
	require.Equal(t, "SUBSCRIPTION_REPURCHASE_BLOCKED", infraerrors.Reason(err))
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

func TestPaymentServiceQueryOrderRepairsPendingWechatOrderViaUpstreamQuery(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Contains(t, r.URL.Path, "/v3/pay/transactions/out-trade-no/PO202604080005")
		require.Equal(t, "10001", r.URL.Query().Get("mchid"))
		require.NotEmpty(t, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"trade_state":"SUCCESS","transaction_id":"wx-trade-100","success_time":"2026-04-08T12:00:00+08:00"}`))
	}))
	defer server.Close()

	originalBaseURL := wechatPayAPIBaseURL
	wechatPayAPIBaseURL = server.URL
	t.Cleanup(func() {
		wechatPayAPIBaseURL = originalBaseURL
	})

	orderRepo := &paymentOrderRepoStub{
		order: &PaymentOrder{
			ID:            5,
			OrderNo:       "PO202604080005",
			UserID:        7,
			OrderType:     PaymentOrderTypeBalance,
			BalanceAmount: 25,
			Status:        PaymentOrderStatusPending,
			PayMethod:     PaymentMethodWechatNative,
			ExpiredAt:     time.Now().Add(30 * time.Minute),
		},
	}
	userRepo := &paymentBalanceUserRepoStub{userRepoStub: &userRepoStub{user: &User{ID: 7, Balance: 10}}}
	svc := &PaymentService{
		orderRepo: orderRepo,
		settingService: NewSettingService(&settingRepoStub{values: map[string]string{
			SettingKeyWechatPayMchID:       "10001",
			SettingKeyWechatPayPrivateKey:  string(privateKeyPEM),
			SettingKeyWechatPayMchSerialNo: "serial-001",
		}}, &config.Config{}),
		userRepo: userRepo,
	}

	order, err := svc.QueryOrder(context.Background(), 7, "PO202604080005")
	require.NoError(t, err)
	require.Equal(t, PaymentOrderStatusPaid, order.Status)
	require.NotNil(t, order.PaidAt)
	require.Equal(t, []float64{25}, userRepo.balanceUpdates)
	require.Equal(t, 35.0, userRepo.user.Balance)
	require.Equal(t, 1, orderRepo.compareCalls)
}

func TestLegacySubscriptionOrderStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		subscriptionStatus string
		expected           string
	}{
		{
			name:               "active subscription stays paid",
			subscriptionStatus: SubscriptionStatusActive,
			expected:           PaymentOrderStatusPaid,
		},
		{
			name:               "expired subscription stays paid",
			subscriptionStatus: SubscriptionStatusExpired,
			expected:           PaymentOrderStatusPaid,
		},
		{
			name:               "suspended subscription becomes closed",
			subscriptionStatus: SubscriptionStatusSuspended,
			expected:           PaymentOrderStatusClosed,
		},
		{
			name:               "unknown status is treated as closed",
			subscriptionStatus: "revoked",
			expected:           PaymentOrderStatusClosed,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, legacySubscriptionOrderStatus(tc.subscriptionStatus))
		})
	}
}
