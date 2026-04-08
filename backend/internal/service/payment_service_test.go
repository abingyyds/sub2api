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
	"github.com/stretchr/testify/require"
)

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
