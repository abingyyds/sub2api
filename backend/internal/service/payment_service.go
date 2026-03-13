package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/smartwalle/alipay/v3"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrPaymentDisabled      = infraerrors.Forbidden("PAYMENT_DISABLED", "online payment is disabled")
	ErrPaymentPlanNotFound  = infraerrors.NotFound("PAYMENT_PLAN_NOT_FOUND", "payment plan not found")
	ErrPaymentOrderNotFound = infraerrors.NotFound("PAYMENT_ORDER_NOT_FOUND", "payment order not found")
	ErrPaymentConfigMissing = infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "payment configuration is incomplete")
	ErrPaymentCreateFailed  = infraerrors.BadRequest("PAYMENT_CREATE_FAILED", "failed to create payment order")
	ErrPaymentSignature     = infraerrors.BadRequest("PAYMENT_SIGNATURE_INVALID", "payment notification signature invalid")
)

// PaymentPlan 套餐配置
type PaymentPlan struct {
	Key           string   `json:"key"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"` // 套餐描述
	Features      []string `json:"features,omitempty"`    // 特性列表
	AmountFen     int      `json:"amount_fen"`
	GroupID       int64    `json:"group_id"`
	ValidityDays  int      `json:"validity_days"`
	Type          string   `json:"type,omitempty"`           // "subscription" (default) or "balance"
	BalanceAmount float64  `json:"balance_amount,omitempty"` // 充值金额（元），仅 type=balance 时使用
}

// PaymentOrder 支付订单
type PaymentOrder struct {
	ID                  int64
	OrderNo             string
	UserID              int64
	PlanKey             string
	GroupID             int64
	AmountFen           int
	ValidityDays        int
	OrderType           string  // "subscription" or "balance"
	BalanceAmount       float64 // 充值金额（元），仅 balance 类型
	PromoCode           string  // 使用的优惠码
	DiscountAmount      int     // 优惠码折扣金额（分）
	Status              string
	PayMethod           string
	WechatTransactionID *string
	AlipayTradeNo       *string
	CodeURL             *string
	PaidAt              *time.Time
	ExpiredAt           time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type PaymentOrderRepository interface {
	Create(ctx context.Context, order *PaymentOrder) error
	GetByID(ctx context.Context, id int64) (*PaymentOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*PaymentOrder, error)
	UpdateStatus(ctx context.Context, orderNo string, status string, transactionID *string, paidAt *time.Time) error
	ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]PaymentOrder, *pagination.PaginationResult, error)
	ListAll(ctx context.Context, params pagination.PaginationParams, status string, orderType string) ([]PaymentOrder, *pagination.PaginationResult, error)
	CloseExpiredOrders(ctx context.Context) (int64, error)
}

// PaymentService 支付服务
type PaymentService struct {
	orderRepo           PaymentOrderRepository
	settingService      *SettingService
	subscriptionService *SubscriptionService
	billingCache        *BillingCacheService
	userRepo            UserRepository
	groupRepo           GroupRepository
	promoService        *PromoService
}

// NewPaymentService 创建支付服务
func NewPaymentService(
	orderRepo PaymentOrderRepository,
	settingService *SettingService,
	subscriptionService *SubscriptionService,
	billingCache *BillingCacheService,
	userRepo UserRepository,
	groupRepo GroupRepository,
	promoService *PromoService,
) *PaymentService {
	return &PaymentService{
		orderRepo:           orderRepo,
		settingService:      settingService,
		subscriptionService: subscriptionService,
		billingCache:        billingCache,
		userRepo:            userRepo,
		groupRepo:           groupRepo,
		promoService:        promoService,
	}
}

// GetPlans 获取所有套餐（从分组的 price_fen 自动生成 + 兼容旧 JSON 配置）
func (s *PaymentService) GetPlans(ctx context.Context) ([]PaymentPlan, error) {
	var plans []PaymentPlan

	// 1. 从分组自动生成套餐：已上架（listed=true）且 price_fen > 0 的活跃分组
	groups, err := s.groupRepo.ListActive(ctx)
	if err != nil {
		log.Printf("[Payment] Failed to list active groups for plans: %v", err)
	} else {
		for _, g := range groups {
			if !g.Listed || g.PriceFen <= 0 {
				continue
			}
			plan := PaymentPlan{
				Key:          fmt.Sprintf("group_%d", g.ID),
				Name:         g.Name,
				Description:  g.Description,
				AmountFen:    g.PriceFen,
				GroupID:      g.ID,
				ValidityDays: g.DefaultValidityDays,
				Type:         "subscription",
			}
			// 特性列表：优先使用管理员自定义，否则自动生成
			if len(g.PlanFeatures) > 0 {
				plan.Features = g.PlanFeatures
			} else {
				features := make([]string, 0, 4)
				if g.DailyLimitUSD != nil && *g.DailyLimitUSD > 0 {
					features = append(features, fmt.Sprintf("每日额度 $%.0f", *g.DailyLimitUSD))
				}
				if g.MonthlyLimitUSD != nil && *g.MonthlyLimitUSD > 0 {
					features = append(features, fmt.Sprintf("每月额度 $%.0f", *g.MonthlyLimitUSD))
				}
				if g.RateMultiplier != 1.0 {
					features = append(features, fmt.Sprintf("费率倍率 %.1fx", g.RateMultiplier))
				}
				plan.Features = features
			}
			plans = append(plans, plan)
		}
	}

	// 2. 兼容旧的 JSON 配置（如有）
	plansJSON, err := s.settingService.GetSettingValue(ctx, SettingKeyPaymentPlans)
	if err == nil && plansJSON != "" {
		var jsonPlans []PaymentPlan
		if err := json.Unmarshal([]byte(plansJSON), &jsonPlans); err != nil {
			log.Printf("[Payment] Failed to parse payment_plans JSON: %v, raw=%q", err, plansJSON)
		} else {
			// 用 map 去重：如果分组已通过 price_fen 生成了套餐，跳过 JSON 中同 group_id 的配置
			existingGroupIDs := make(map[int64]bool, len(plans))
			for _, p := range plans {
				existingGroupIDs[p.GroupID] = true
			}
			for i := range jsonPlans {
				if jsonPlans[i].Type == "" {
					jsonPlans[i].Type = "subscription"
				}
				if jsonPlans[i].Type == "subscription" && existingGroupIDs[jsonPlans[i].GroupID] {
					continue // 分组已通过 price_fen 生成
				}
				plans = append(plans, jsonPlans[i])
			}
			// 对 JSON 中的套餐补充分组信息
			s.enrichPlansFromGroups(ctx, plans)
		}
	}

	if plans == nil {
		plans = []PaymentPlan{}
	}
	return plans, nil
}

// RechargePlan 充值优惠套餐
type RechargePlan struct {
	Key           string  `json:"key"`
	Name          string  `json:"name"`
	Description   string  `json:"description,omitempty"`
	PayAmountFen  int     `json:"pay_amount_fen"` // 实付金额（分）
	BalanceAmount float64 `json:"balance_amount"` // 到账金额（美元）
	Popular       bool    `json:"popular,omitempty"`
}

// RechargeInfo 充值信息（含优惠套餐和最低金额）
type RechargeInfo struct {
	MinAmount float64        `json:"min_amount"` // 最低充值金额（元）
	Plans     []RechargePlan `json:"plans"`
}

// GetRechargeInfo 获取充值信息
func (s *PaymentService) GetRechargeInfo(ctx context.Context) (*RechargeInfo, error) {
	info := &RechargeInfo{
		MinAmount: 0,
		Plans:     []RechargePlan{},
	}

	// 读取最低充值金额
	if minStr, _ := s.settingService.GetSettingValue(ctx, SettingKeyRechargeMinAmount); minStr != "" {
		if v, err := strconv.ParseFloat(minStr, 64); err == nil && v > 0 {
			info.MinAmount = v
		}
	}

	// 读取充值优惠套餐
	if plansJSON, _ := s.settingService.GetSettingValue(ctx, SettingKeyRechargePlans); plansJSON != "" {
		if err := json.Unmarshal([]byte(plansJSON), &info.Plans); err != nil {
			log.Printf("[Payment] Failed to parse recharge_plans JSON: %v", err)
		}
	}

	return info, nil
}

// enrichPlansFromGroups 从分组信息自动填充套餐的描述和特性
func (s *PaymentService) enrichPlansFromGroups(ctx context.Context, plans []PaymentPlan) {
	for i := range plans {
		plan := &plans[i]
		if plan.Type != "subscription" || plan.GroupID == 0 {
			continue
		}
		// 如果已手动配置了 features，跳过
		if len(plan.Features) > 0 {
			continue
		}
		group, err := s.groupRepo.GetByID(ctx, plan.GroupID)
		if err != nil {
			log.Printf("[Payment] enrichPlans: failed to get group %d for plan %s: %v", plan.GroupID, plan.Key, err)
			continue
		}
		log.Printf("[Payment] enrichPlans: plan=%s group=%d name=%q desc=%q daily=%.2f monthly=%.2f rate=%.2f",
			plan.Key, group.ID, group.Name, group.Description,
			floatVal(group.DailyLimitUSD), floatVal(group.MonthlyLimitUSD), group.RateMultiplier)
		// 填充描述
		if plan.Description == "" && group.Description != "" {
			plan.Description = group.Description
		}
		// 自动生成特性列表
		features := make([]string, 0, 4)
		if group.DailyLimitUSD != nil && *group.DailyLimitUSD > 0 {
			features = append(features, fmt.Sprintf("每日额度 $%.0f", *group.DailyLimitUSD))
		}
		if group.MonthlyLimitUSD != nil && *group.MonthlyLimitUSD > 0 {
			features = append(features, fmt.Sprintf("每月额度 $%.0f", *group.MonthlyLimitUSD))
		}
		if group.RateMultiplier != 1.0 {
			features = append(features, fmt.Sprintf("费率倍率 %.1fx", group.RateMultiplier))
		}
		plan.Features = features
	}
}

func floatVal(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// CreateOrder 创建支付订单
func (s *PaymentService) CreateOrder(ctx context.Context, userID int64, planKey string, promoCode string, payMethod string) (*PaymentOrder, error) {
	// 检查支付是否启用
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return nil, ErrPaymentDisabled
	}

	// 查找套餐
	plans, err := s.GetPlans(ctx)
	if err != nil {
		return nil, err
	}
	var plan *PaymentPlan
	for i := range plans {
		if plans[i].Key == planKey {
			plan = &plans[i]
			break
		}
	}
	if plan == nil {
		return nil, ErrPaymentPlanNotFound
	}

	// 计算优惠码折扣
	var discountFen int
	promoCode = strings.TrimSpace(promoCode)
	if promoCode != "" && s.promoService != nil && s.settingService.IsPromoCodeEnabled(ctx) {
		discount, _, err := s.promoService.ValidateAndCalculateDiscount(ctx, userID, promoCode, plan.AmountFen)
		if err != nil {
			return nil, err
		}
		discountFen = discount
	}

	// 计算最终金额
	finalAmount := plan.AmountFen - discountFen
	if finalAmount < 1 {
		finalAmount = 1 // 最低 1 分
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 确定订单类型
	orderType := plan.Type
	if orderType == "" {
		orderType = "subscription"
	}

	// 确定支付方式
	if payMethod == "" {
		payMethod = "wechat"
	}
	var payMethodStr string
	if payMethod == "alipay" {
		payMethodStr = PaymentMethodAlipayNative
	} else {
		payMethodStr = PaymentMethodWechatNative
	}

	// 创建订单记录
	order := &PaymentOrder{
		OrderNo:        orderNo,
		UserID:         userID,
		PlanKey:        plan.Key,
		AmountFen:      finalAmount,
		Status:         PaymentOrderStatusPending,
		PayMethod:      payMethodStr,
		OrderType:      orderType,
		PromoCode:      promoCode,
		DiscountAmount: discountFen,
		ExpiredAt:      time.Now().Add(30 * time.Minute), // 30分钟过期
	}

	if orderType == "balance" {
		order.GroupID = 0
		order.ValidityDays = 0
		order.BalanceAmount = plan.BalanceAmount
	} else {
		order.GroupID = plan.GroupID
		order.ValidityDays = plan.ValidityDays
	}

	// 根据支付方式调用不同的下单接口
	var codeURL string
	if payMethod == "alipay" {
		codeURL, err = s.createAlipayNativeOrder(ctx, order, plan.Name)
	} else {
		codeURL, err = s.createWechatNativeOrder(ctx, order, plan.Name)
	}
	if err != nil {
		log.Printf("[Payment] Failed to create %s native order: %v", payMethod, err)
		return nil, err
	}
	order.CodeURL = &codeURL

	// 保存订单到数据库
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("save order: %w", err)
	}

	return order, nil
}

// CreateRechargeOrder 创建余额充值订单（支持自定义金额）
func (s *PaymentService) CreateRechargeOrder(ctx context.Context, userID int64, amountYuan float64, promoCode string, payMethod string) (*PaymentOrder, error) {
	// 检查支付是否启用
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return nil, ErrPaymentDisabled
	}

	if amountYuan <= 0 {
		return nil, infraerrors.BadRequest("INVALID_AMOUNT", "recharge amount must be positive")
	}

	// 检查最低充值金额
	if minAmountStr, _ := s.settingService.GetSettingValue(ctx, SettingKeyRechargeMinAmount); minAmountStr != "" {
		if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil && minAmount > 0 {
			if amountYuan < minAmount {
				return nil, infraerrors.BadRequest("AMOUNT_TOO_LOW", fmt.Sprintf("minimum recharge amount is ¥%.0f", minAmount))
			}
		}
	}

	amountFen := int(amountYuan * 100)

	// 计算优惠码折扣
	var discountFen int
	promoCode = strings.TrimSpace(promoCode)
	if promoCode != "" && s.promoService != nil && s.settingService.IsPromoCodeEnabled(ctx) {
		discount, _, err := s.promoService.ValidateAndCalculateDiscount(ctx, userID, promoCode, amountFen)
		if err != nil {
			return nil, err
		}
		discountFen = discount
	}

	finalAmount := amountFen - discountFen
	if finalAmount < 1 {
		finalAmount = 1
	}

	orderNo := generateOrderNo()

	// 确定支付方式
	if payMethod == "" {
		payMethod = "wechat"
	}
	var payMethodStr string
	if payMethod == "alipay" {
		payMethodStr = PaymentMethodAlipayNative
	} else {
		payMethodStr = PaymentMethodWechatNative
	}

	order := &PaymentOrder{
		OrderNo:        orderNo,
		UserID:         userID,
		PlanKey:        "recharge_custom",
		AmountFen:      finalAmount,
		GroupID:        0,
		ValidityDays:   0,
		OrderType:      "balance",
		BalanceAmount:  amountYuan,
		PromoCode:      promoCode,
		DiscountAmount: discountFen,
		Status:         PaymentOrderStatusPending,
		PayMethod:      payMethodStr,
		ExpiredAt:      time.Now().Add(30 * time.Minute),
	}

	description := fmt.Sprintf("余额充值 %.2f 元", amountYuan)
	var codeURL string
	var err error
	if payMethod == "alipay" {
		codeURL, err = s.createAlipayNativeOrder(ctx, order, description)
	} else {
		codeURL, err = s.createWechatNativeOrder(ctx, order, description)
	}
	if err != nil {
		log.Printf("[Payment] Failed to create %s native order for recharge: %v", payMethod, err)
		return nil, err
	}
	order.CodeURL = &codeURL

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("save recharge order: %w", err)
	}

	return order, nil
}

// QueryOrder 查询订单状态
func (s *PaymentService) QueryOrder(ctx context.Context, userID int64, orderNo string) (*PaymentOrder, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, ErrPaymentOrderNotFound
	}
	return order, nil
}

// ListOrders 列出用户订单
func (s *PaymentService) ListOrders(ctx context.Context, userID int64, params pagination.PaginationParams) ([]PaymentOrder, *pagination.PaginationResult, error) {
	return s.orderRepo.ListByUserID(ctx, userID, params)
}

// ListAllOrders 列出所有订单（管理员）
func (s *PaymentService) ListAllOrders(ctx context.Context, params pagination.PaginationParams, status string, orderType string) ([]PaymentOrder, *pagination.PaginationResult, error) {
	return s.orderRepo.ListAll(ctx, params, status, orderType)
}

// HandleWechatNotify 处理微信支付回调通知
func (s *PaymentService) HandleWechatNotify(ctx context.Context, body []byte, wechatpayTimestamp, wechatpayNonce, wechatpaySignature, wechatpaySerial string) error {
	// 1. 验证签名（使用微信支付公钥）
	if err := s.verifyWechatSignature(ctx, wechatpayTimestamp, wechatpayNonce, string(body), wechatpaySignature, wechatpaySerial); err != nil {
		log.Printf("[Payment] WechatNotify signature verification failed: %v (serial=%s)", err, wechatpaySerial)
		return ErrPaymentSignature
	}

	// 2. 解密通知内容
	var notification struct {
		EventType string `json:"event_type"`
		Resource  struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			Nonce          string `json:"nonce"`
			AssociatedData string `json:"associated_data"`
			OriginalType   string `json:"original_type"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &notification); err != nil {
		return fmt.Errorf("parse notification: %w", err)
	}

	if notification.EventType != "TRANSACTION.SUCCESS" {
		return nil // 忽略非成功通知
	}

	// 3. 使用 APIv3 密钥解密
	apiKey, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayAPIv3Key)
	plaintext, err := decryptAEAD(apiKey, notification.Resource.Nonce, notification.Resource.Ciphertext, notification.Resource.AssociatedData)
	if err != nil {
		return fmt.Errorf("decrypt notification: %w", err)
	}

	// 4. 解析交易信息
	var transaction struct {
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		TradeState    string `json:"trade_state"`
		SuccessTime   string `json:"success_time"`
	}
	if err := json.Unmarshal(plaintext, &transaction); err != nil {
		return fmt.Errorf("parse transaction: %w", err)
	}

	if transaction.TradeState != "SUCCESS" {
		return nil
	}

	// 5. 查询订单
	order, err := s.orderRepo.GetByOrderNo(ctx, transaction.OutTradeNo)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// 幂等检查：已支付的订单不再处理
	if order.Status == PaymentOrderStatusPaid {
		return nil
	}

	// 6. 更新订单状态
	now := time.Now()
	if err := s.orderRepo.UpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPaid, &transaction.TransactionID, &now); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	// 7. 根据订单类型处理
	if order.OrderType == "balance" {
		// 充值余额
		if err := s.userRepo.UpdateBalance(ctx, order.UserID, order.BalanceAmount); err != nil {
			log.Printf("[Payment] Failed to update balance for order %s, user %d: %v", order.OrderNo, order.UserID, err)
			return fmt.Errorf("update balance: %w", err)
		}
		// 失效余额缓存
		if err := s.billingCache.InvalidateUserBalance(ctx, order.UserID); err != nil {
			log.Printf("[Payment] Failed to invalidate balance cache for user %d: %v", order.UserID, err)
		}
		log.Printf("[Payment] Order %s paid successfully, balance +%.2f for user %d",
			order.OrderNo, order.BalanceAmount, order.UserID)
	} else {
		// 分配订阅
		_, _, err = s.subscriptionService.AssignOrExtendSubscription(ctx, &AssignSubscriptionInput{
			UserID:       order.UserID,
			GroupID:      order.GroupID,
			ValidityDays: order.ValidityDays,
			AssignedBy:   0, // 系统自动分配
			Notes:        fmt.Sprintf("微信支付订单 %s", order.OrderNo),
		})
		if err != nil {
			log.Printf("[Payment] Failed to assign subscription for order %s, user %d: %v", order.OrderNo, order.UserID, err)
			return fmt.Errorf("assign subscription: %w", err)
		}
		log.Printf("[Payment] Order %s paid successfully, subscription assigned for user %d, group %d, %d days",
			order.OrderNo, order.UserID, order.GroupID, order.ValidityDays)
	}

	// 8. 记录优惠码使用（支付成功后才记录）
	if order.PromoCode != "" && order.DiscountAmount > 0 && s.promoService != nil {
		if err := s.promoService.ApplyPromoCode(ctx, order.UserID, order.PromoCode, order.OrderNo, float64(order.DiscountAmount)); err != nil {
			log.Printf("[Payment] Failed to apply promo code for order %s: %v", order.OrderNo, err)
			// 不阻断支付流程
		}
	}

	return nil
}

// ===========================
// WeChat Pay API Integration
// ===========================

// createWechatNativeOrder 调用微信支付 Native 下单 API
func (s *PaymentService) createWechatNativeOrder(ctx context.Context, order *PaymentOrder, planName string) (string, error) {
	appID, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayAppID)
	mchID, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayMchID)
	notifyURL, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayNotifyURL)
	privateKeyPEM, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayPrivateKey)
	mchSerialNo, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayMchSerialNo)

	if appID == "" || mchID == "" || notifyURL == "" || privateKeyPEM == "" || mchSerialNo == "" {
		missing := make([]string, 0, 5)
		if appID == "" {
			missing = append(missing, "appid")
		}
		if mchID == "" {
			missing = append(missing, "mch_id")
		}
		if notifyURL == "" {
			missing = append(missing, "notify_url")
		}
		if privateKeyPEM == "" {
			missing = append(missing, "private_key")
		}
		if mchSerialNo == "" {
			missing = append(missing, "mch_serial_no")
		}
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", fmt.Sprintf("payment configuration incomplete, missing: %s", strings.Join(missing, ", ")))
	}

	// 过滤 planName 中的 emoji 和非 BMP 字符（微信支付不接受）
	safeDesc := regexp.MustCompile(`[^\x{0000}-\x{FFFF}]`).ReplaceAllString(planName, "")
	safeDesc = strings.TrimSpace(safeDesc)
	if safeDesc == "" {
		safeDesc = "订单支付"
	}

	// 构造请求体
	reqBody := map[string]any{
		"appid":        appID,
		"mchid":        mchID,
		"description":  safeDesc,
		"out_trade_no": order.OrderNo,
		"time_expire":  order.ExpiredAt.Format(time.RFC3339),
		"notify_url":   notifyURL,
		"amount": map[string]any{
			"total":    order.AmountFen,
			"currency": "CNY",
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// 签名
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := generateNonce()
	method := "POST"
	url := "/v3/pay/transactions/native"

	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", method, url, timestamp, nonce, string(bodyBytes))

	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "invalid private key PEM format")
	}

	signature, err := signSHA256WithRSA(privateKey, []byte(signStr))
	if err != nil {
		return "", fmt.Errorf("sign request: %w", err)
	}

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, method, "https://api.mch.weixin.qq.com"+url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(
		`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`,
		mchID, nonce, timestamp, mchSerialNo, signature,
	))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("wechat api request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Payment] WeChat API error: status=%d body=%s", resp.StatusCode, string(respBody))
		return "", infraerrors.BadRequest("WECHAT_API_ERROR", fmt.Sprintf("wechat pay api error: %s", string(respBody)))
	}

	var result struct {
		CodeURL string `json:"code_url"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse wechat response: %w", err)
	}

	return result.CodeURL, nil
}

// verifyWechatSignature 验证微信支付回调签名
func (s *PaymentService) verifyWechatSignature(ctx context.Context, timestamp, nonce, body, signature, serial string) error {
	// 获取微信支付公钥
	publicKeyPEM, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayPublicKey)
	if publicKeyPEM == "" {
		return fmt.Errorf("wechat pay public key not configured")
	}

	// 解析公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return fmt.Errorf("failed to decode public key PEM")
	}

	var rsaPub *rsa.PublicKey

	switch block.Type {
	case "CERTIFICATE":
		// 微信平台证书模式：从 X.509 证书中提取公钥
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("parse certificate: %w", err)
		}
		var ok bool
		rsaPub, ok = cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("certificate does not contain RSA public key")
		}
	case "RSA PUBLIC KEY":
		// PKCS#1 格式公钥
		var err error
		rsaPub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("parse PKCS1 public key: %w", err)
		}
	default:
		// PKIX 格式公钥 (PUBLIC KEY)
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			// 再尝试 PKCS#1 格式
			rsaPub2, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
			if err2 != nil {
				return fmt.Errorf("parse public key: PKIX=%w, PKCS1=%v", err, err2)
			}
			rsaPub = rsaPub2
		} else {
			var ok bool
			rsaPub, ok = pub.(*rsa.PublicKey)
			if !ok {
				return fmt.Errorf("not an RSA public key")
			}
		}
	}

	log.Printf("[Payment] verifyWechatSignature: PEM block type=%q, key bits=%d, serial=%s", block.Type, rsaPub.N.BitLen(), serial)
	log.Printf("[Payment] verifyWechatSignature: timestamp=%q nonce=%q signature_len=%d body_len=%d", timestamp, nonce, len(signature), len(body))

	// 构造签名串
	signStr := fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, body)
	hash := sha256.Sum256([]byte(signStr))

	// 解码签名
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	// 验证签名
	if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hash[:], sig); err != nil {
		log.Printf("[Payment] RSA verify failed: %v (sig_bytes=%d, hash=%x)", err, len(sig), hash[:8])
		return err
	}
	return nil
}

// ===========================
// Crypto Helper Functions
// ===========================

func parsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}

	// 尝试 PKCS8 格式
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// 尝试 PKCS1 格式
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return rsaKey, nil
}

func signSHA256WithRSA(key *rsa.PrivateKey, data []byte) (string, error) {
	hash := sha256.Sum256(data)
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func decryptAEAD(apiKey, nonce, ciphertext, associatedData string) ([]byte, error) {
	key := []byte(apiKey)
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, []byte(nonce), ciphertextBytes, []byte(associatedData))
}

func generateOrderNo() string {
	now := time.Now()
	b := make([]byte, 7)
	rand.Read(b)
	return fmt.Sprintf("P%s%X", now.Format("20060102150405"), b)
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return strings.ToUpper(fmt.Sprintf("%x", b))
}

// ===========================
// Alipay API Integration
// ===========================

// initAlipayClient 初始化支付宝客户端
func (s *PaymentService) initAlipayClient(ctx context.Context) (*alipay.Client, error) {
	appID, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayAppID)
	privateKey, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayPrivateKey)
	alipayPublicKey, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayPublicKey)
	isProduction := s.settingService.GetSettingValue(ctx, SettingKeyAlipayIsProduction) == "true"

	if appID == "" || privateKey == "" || alipayPublicKey == "" {
		missing := make([]string, 0, 3)
		if appID == "" {
			missing = append(missing, "app_id")
		}
		if privateKey == "" {
			missing = append(missing, "private_key")
		}
		if alipayPublicKey == "" {
			missing = append(missing, "public_key")
		}
		return nil, infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", fmt.Sprintf("alipay configuration incomplete, missing: %s", strings.Join(missing, ", ")))
	}

	client, err := alipay.New(appID, privateKey, isProduction)
	if err != nil {
		return nil, fmt.Errorf("init alipay client: %w", err)
	}

	err = client.LoadAliPayPublicKey(alipayPublicKey)
	if err != nil {
		return nil, fmt.Errorf("load alipay public key: %w", err)
	}

	return client, nil
}

// createAlipayNativeOrder 调用支付宝当面付（扫码支付）API
func (s *PaymentService) createAlipayNativeOrder(ctx context.Context, order *PaymentOrder, subject string) (string, error) {
	client, err := s.initAlipayClient(ctx)
	if err != nil {
		return "", err
	}

	notifyURL, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayNotifyURL)
	if notifyURL == "" {
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "alipay notify_url is not configured")
	}

	var p = alipay.TradePrecreate{}
	p.NotifyURL = notifyURL
	p.Subject = subject
	p.OutTradeNo = order.OrderNo
	p.TotalAmount = fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)
	p.TimeExpire = order.ExpiredAt.Format("2006-01-02 15:04:05")

	rsp, err := client.TradePrecreate(ctx, p)
	if err != nil {
		return "", fmt.Errorf("alipay trade precreate: %w", err)
	}

	if rsp.IsFailure() {
		return "", fmt.Errorf("alipay error: %s - %s", rsp.Content.Code, rsp.Content.Msg)
	}

	return rsp.Content.QRCode, nil
}

// HandleAlipayNotify 处理支付宝支付回调通知
func (s *PaymentService) HandleAlipayNotify(ctx context.Context, req *http.Request) error {
	client, err := s.initAlipayClient(ctx)
	if err != nil {
		return err
	}

	// SDK 自动验签
	notification, err := client.GetTradeNotification(req)
	if err != nil {
		log.Printf("[Payment] AlipayNotify signature verification failed: %v", err)
		return ErrPaymentSignature
	}

	// 检查交易状态
	if notification.TradeStatus != "TRADE_SUCCESS" && notification.TradeStatus != "TRADE_FINISHED" {
		return nil // 非成功状态，忽略
	}

	// 查询订单
	order, err := s.orderRepo.GetByOrderNo(ctx, notification.OutTradeNo)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// 幂等检查：已支付的订单不再处理
	if order.Status == PaymentOrderStatusPaid {
		return nil
	}

	// 更新订单状态
	paidAt := time.Now()
	tradeNo := notification.TradeNo
	err = s.orderRepo.UpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPaid, &tradeNo, &paidAt)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	log.Printf("[Payment] Order %s paid via Alipay, trade_no=%s", order.OrderNo, tradeNo)

	// 处理业务逻辑（订阅或余额）
	if order.OrderType == "subscription" {
		// 分配订阅
		err = s.subscriptionService.AssignOrExtendSubscription(ctx, order.UserID, order.GroupID, order.ValidityDays)
		if err != nil {
			log.Printf("[Payment] Failed to assign subscription for order %s: %v", order.OrderNo, err)
			return fmt.Errorf("assign subscription: %w", err)
		}
		log.Printf("[Payment] Order %s paid successfully, subscription assigned for user %d, group %d, %d days",
			order.OrderNo, order.UserID, order.GroupID, order.ValidityDays)
	} else if order.OrderType == "balance" {
		// 充值余额
		err = s.userRepo.UpdateBalance(ctx, order.UserID, order.BalanceAmount)
		if err != nil {
			log.Printf("[Payment] Failed to update balance for order %s: %v", order.OrderNo, err)
			return fmt.Errorf("update balance: %w", err)
		}
		// 清除余额缓存
		if s.billingCache != nil {
			s.billingCache.InvalidateUserBalance(ctx, order.UserID)
		}
		log.Printf("[Payment] Order %s paid successfully, balance recharged for user %d, amount=%.2f",
			order.OrderNo, order.UserID, order.BalanceAmount)
	}

	// 记录优惠码使用（支付成功后才记录）
	if order.PromoCode != "" && order.DiscountAmount > 0 && s.promoService != nil {
		if err := s.promoService.ApplyPromoCode(ctx, order.UserID, order.PromoCode, order.OrderNo, float64(order.DiscountAmount)); err != nil {
			log.Printf("[Payment] Failed to apply promo code for order %s: %v", order.OrderNo, err)
			// 不阻断支付流程
		}
	}

	return nil
}
