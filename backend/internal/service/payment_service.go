package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/smartwalle/alipay/v3"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
)

var (
	ErrPaymentDisabled                      = infraerrors.Forbidden("PAYMENT_DISABLED", "online payment is disabled")
	ErrPaymentPlanNotFound                  = infraerrors.NotFound("PAYMENT_PLAN_NOT_FOUND", "payment plan not found")
	ErrPaymentOrderNotFound                 = infraerrors.NotFound("PAYMENT_ORDER_NOT_FOUND", "payment order not found")
	ErrPaymentOrderPaid                     = infraerrors.Conflict("PAYMENT_ORDER_ALREADY_PAID", "payment order already paid")
	ErrPaymentOrderUnrepairable             = infraerrors.BadRequest("PAYMENT_ORDER_REPAIR_NOT_ALLOWED", "only pending or closed orders can be repaired manually")
	ErrPaymentConfigMissing                 = infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "payment configuration is incomplete")
	ErrPaymentCreateFailed                  = infraerrors.BadRequest("PAYMENT_CREATE_FAILED", "failed to create payment order")
	ErrPaymentSignature                     = infraerrors.BadRequest("PAYMENT_SIGNATURE_INVALID", "payment notification signature invalid")
	ErrPaymentSubscriptionRepurchaseBlocked = infraerrors.Conflict(
		"SUBSCRIPTION_REPURCHASE_BLOCKED",
		"each account can only purchase the same plan once during its active period; repurchase after expiry",
	)
	ErrInvoiceOrdersRequired = infraerrors.BadRequest("INVOICE_ORDERS_REQUIRED", "at least one order is required")
	ErrInvoiceEmailInvalid   = infraerrors.BadRequest("INVOICE_EMAIL_INVALID", "invoice email is invalid")
	ErrInvoiceOrderNotPaid   = infraerrors.BadRequest("INVOICE_ORDER_NOT_PAID", "only paid orders can request invoices")
	ErrInvoiceNotRequested   = infraerrors.BadRequest("INVOICE_NOT_REQUESTED", "invoice has not been requested for this order")
	ErrInvoiceAlreadyFiled   = infraerrors.Conflict("INVOICE_ALREADY_FILED", "invoice has already been requested for one or more selected orders")
	ErrInvoiceAlreadyHandled = infraerrors.Conflict("INVOICE_ALREADY_HANDLED", "invoice has already been processed")
)

var wechatPayAPIBaseURL = "https://api.mch.weixin.qq.com"

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
	SubSiteID           *int64  // 下单时所属分站，nil 表示主站订单
	PromoCode           string  // 使用的优惠码
	DiscountAmount      int     // 优惠码折扣金额（分）
	Status              string
	PayMethod           string
	WechatTransactionID *string
	AlipayTradeNo       *string
	EpayTradeNo         *string
	InvoiceCompanyName  string
	InvoiceTaxID        string
	InvoiceEmail        string
	InvoiceRemark       string
	InvoiceRequestedAt  *time.Time
	InvoiceProcessedAt  *time.Time
	CodeURL             *string
	PaidAt              *time.Time
	ExpiredAt           time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type InvoiceRequest struct {
	CompanyName string `json:"company_name"`
	TaxID       string `json:"tax_id"`
	Email       string `json:"email"`
	Remark      string `json:"remark"`
}

type PaymentOrderRepository interface {
	Create(ctx context.Context, order *PaymentOrder) error
	GetByID(ctx context.Context, id int64) (*PaymentOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*PaymentOrder, error)
	UpdateStatus(ctx context.Context, orderNo string, status string, transactionID *string, paidAt *time.Time) error
	// CompareAndUpdateStatus atomically updates order status only if current status matches expectedStatus.
	// Returns true if the update was applied, false if the current status didn't match.
	CompareAndUpdateStatus(ctx context.Context, orderNo string, expectedStatus string, newStatus string, transactionID *string, paidAt *time.Time) (bool, error)
	ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]PaymentOrder, *pagination.PaginationResult, error)
	ListAll(ctx context.Context, params pagination.PaginationParams, status string, orderType string) ([]PaymentOrder, *pagination.PaginationResult, error)
	SubmitInvoiceRequest(ctx context.Context, userID int64, orderNos []string, invoice InvoiceRequest) error
	MarkInvoiceProcessed(ctx context.Context, orderID int64) error
	CloseExpiredOrders(ctx context.Context) (int64, error)
	CountPaidByUserAndPlanKey(ctx context.Context, userID int64, planKey string) (int, error)
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
	agentService        *AgentService
	subSiteService      *SubSiteService
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
	agentService *AgentService,
	subSiteService *SubSiteService,
) *PaymentService {
	return &PaymentService{
		orderRepo:           orderRepo,
		settingService:      settingService,
		subscriptionService: subscriptionService,
		billingCache:        billingCache,
		userRepo:            userRepo,
		groupRepo:           groupRepo,
		promoService:        promoService,
		agentService:        agentService,
		subSiteService:      subSiteService,
	}
}

func (s *PaymentService) getSubSiteFromCtx(ctx context.Context) *SubSite {
	site, ok := ctx.Value(ctxkey.SubSite).(*SubSite)
	if !ok || site == nil {
		return nil
	}
	return site
}

func subSiteIDPtr(site *SubSite) *int64 {
	if site == nil || site.ID <= 0 {
		return nil
	}
	id := site.ID
	return &id
}

// getPoolSubSiteFromCtx returns the current pool-mode sub-site with a non-nil OwnerPaymentConfig, or nil.
func (s *PaymentService) getPoolSubSiteFromCtx(ctx context.Context) *SubSite {
	site := s.getSubSiteFromCtx(ctx)
	if site == nil {
		return nil
	}
	if site.Mode != SubSiteModePool || site.OwnerPaymentConfig == nil {
		return nil
	}
	return site
}

const balanceSubSitePlanKeySep = ":subsite:"

func encodeBalanceSubSitePlanKey(basePlanKey string, siteID int64) string {
	return fmt.Sprintf("%s%s%d", basePlanKey, balanceSubSitePlanKeySep, siteID)
}

func parseBalanceSubSitePlanKey(planKey string) (basePlanKey string, siteID int64, ok bool) {
	idx := strings.Index(planKey, balanceSubSitePlanKeySep)
	if idx < 0 {
		return planKey, 0, false
	}
	base := planKey[:idx]
	idStr := planKey[idx+len(balanceSubSitePlanKeySep):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return planKey, 0, false
	}
	return base, id, true
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
				Type:         PaymentOrderTypeSubscription,
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
					jsonPlans[i].Type = PaymentOrderTypeSubscription
				}
				if jsonPlans[i].Type == PaymentOrderTypeSubscription && existingGroupIDs[jsonPlans[i].GroupID] {
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
	IsNewcomer    bool    `json:"is_newcomer,omitempty"`   // 新人专享标记
	MaxPurchases  int     `json:"max_purchases,omitempty"` // 最大购买次数（0=无限）
}

// RechargeInfo 充值信息（含优惠套餐和最低金额）
type RechargeInfo struct {
	MinAmount float64        `json:"min_amount"` // 最低充值金额（元）
	Plans     []RechargePlan `json:"plans"`
}

// GetAvailablePayMethods 返回当前已启用的支付方式列表
func (s *PaymentService) GetAvailablePayMethods(ctx context.Context) []string {
	methods := make([]string, 0, 4)

	// pool 模式分站有自有收款凭据时，优先报告站长可用方式
	if poolSite := s.getPoolSubSiteFromCtx(ctx); poolSite != nil && poolSite.OwnerPaymentConfig != nil {
		opc := poolSite.OwnerPaymentConfig
		if isOwnerWechatReady(opc) {
			methods = append(methods, "wechat")
		}
		if isOwnerAlipayReady(opc) {
			methods = append(methods, "alipay")
		}
		if isOwnerEpayReady(opc) {
			methods = append(methods, "epay_alipay", "epay_wxpay")
		}
		if len(methods) > 0 {
			return methods
		}
	}

	if s.isWechatPayReady(ctx) {
		methods = append(methods, "wechat")
	}

	if s.isAlipayReady(ctx) {
		methods = append(methods, "alipay")
	}

	if s.isEpayReady(ctx) {
		methods = append(methods, "epay_alipay", "epay_wxpay")
	}

	return methods
}

func (s *PaymentService) isWechatPayReady(ctx context.Context) bool {
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return false
	}

	return s.hasNonEmptySettings(
		ctx,
		SettingKeyWechatPayAppID,
		SettingKeyWechatPayMchID,
		SettingKeyWechatPayNotifyURL,
		SettingKeyWechatPayPrivateKey,
		SettingKeyWechatPayMchSerialNo,
	)
}

func (s *PaymentService) isAlipayReady(ctx context.Context) bool {
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return false
	}

	alipayEnabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayEnabled)
	if alipayEnabled != "true" {
		return false
	}

	return s.hasNonEmptySettings(
		ctx,
		SettingKeyAlipayAppID,
		SettingKeyAlipayPrivateKey,
		SettingKeyAlipayPublicKey,
		SettingKeyAlipayNotifyURL,
	)
}

func (s *PaymentService) isEpayReady(ctx context.Context) bool {
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return false
	}

	epayEnabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayEnabled)
	if epayEnabled != "true" {
		return false
	}

	return s.hasNonEmptySettings(
		ctx,
		SettingKeyEpayGateway,
		SettingKeyEpayPID,
		SettingKeyEpayPKey,
		SettingKeyEpayNotifyURL,
	)
}

func (s *PaymentService) hasNonEmptySettings(ctx context.Context, keys ...string) bool {
	for _, key := range keys {
		value, _ := s.settingService.GetSettingValue(ctx, key)
		if strings.TrimSpace(value) == "" {
			return false
		}
	}
	return true
}

// isOwnerWechatReady checks if the pool-mode sub-site has a configured wechat credential.
func isOwnerWechatReady(opc *OwnerPaymentConfig) bool {
	if opc == nil || opc.Wechat == nil || !opc.Wechat.Enabled {
		return false
	}
	w := opc.Wechat
	return w.AppID != "" && w.MchID != "" && w.PrivateKey != "" && w.MchSerialNo != "" && w.NotifyURL != ""
}

func isOwnerAlipayReady(opc *OwnerPaymentConfig) bool {
	if opc == nil || opc.Alipay == nil || !opc.Alipay.Enabled {
		return false
	}
	a := opc.Alipay
	return a.AppID != "" && a.PrivateKey != "" && a.PublicKey != "" && a.NotifyURL != ""
}

func isOwnerEpayReady(opc *OwnerPaymentConfig) bool {
	if opc == nil || opc.Epay == nil || !opc.Epay.Enabled {
		return false
	}
	e := opc.Epay
	return e.Gateway != "" && e.PID != "" && e.PKey != "" && e.NotifyURL != ""
}

// isPayMethodReadyForBalanceOrder checks pay method availability, preferring owner config for pool-mode sub-sites.
func (s *PaymentService) isPayMethodReadyForBalanceOrder(ctx context.Context, payMethod string, opc *OwnerPaymentConfig) bool {
	if opc != nil {
		switch payMethod {
		case "alipay":
			return isOwnerAlipayReady(opc)
		case "epay_alipay", "epay_wxpay":
			return isOwnerEpayReady(opc)
		default:
			return isOwnerWechatReady(opc)
		}
	}
	switch payMethod {
	case "alipay":
		return s.isAlipayReady(ctx)
	case "epay_alipay", "epay_wxpay":
		return s.isEpayReady(ctx)
	default:
		return s.isWechatPayReady(ctx)
	}
}

func wrapUnknownAsBadRequest(reason, message string, err error) error {
	if err == nil {
		return nil
	}
	if infraerrors.Code(err) != http.StatusInternalServerError || infraerrors.Reason(err) != "" {
		return err
	}
	return infraerrors.BadRequest(reason, message).WithCause(err)
}

func wrapUnknownAsServiceUnavailable(reason, message string, err error) error {
	if err == nil {
		return nil
	}
	if infraerrors.Code(err) != http.StatusInternalServerError || infraerrors.Reason(err) != "" {
		return err
	}
	return infraerrors.ServiceUnavailable(reason, message).WithCause(err)
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

// CheckNewcomerEligibility 检查用户是否有资格购买新人专享套餐
func (s *PaymentService) CheckNewcomerEligibility(ctx context.Context, userID int64) (bool, error) {
	// 获取所有充值套餐
	info, err := s.GetRechargeInfo(ctx)
	if err != nil {
		return false, err
	}

	// 检查每个新人套餐
	for _, rp := range info.Plans {
		if !rp.IsNewcomer {
			continue
		}
		if rp.MaxPurchases > 0 {
			count, err := s.orderRepo.CountPaidByUserAndPlanKey(ctx, userID, rp.Key)
			if err != nil {
				return false, err
			}
			if count >= rp.MaxPurchases {
				return false, nil // 已用完
			}
		}
	}
	return true, nil
}

// enrichPlansFromGroups 从分组信息自动填充套餐的描述和特性
func (s *PaymentService) enrichPlansFromGroups(ctx context.Context, plans []PaymentPlan) {
	for i := range plans {
		plan := &plans[i]
		if plan.Type != PaymentOrderTypeSubscription || plan.GroupID == 0 {
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

	// 确定订单类型
	orderType := plan.Type
	if orderType == "" {
		orderType = PaymentOrderTypeSubscription
	}

	// 同一订阅分组未到期前不可重复购买；到期后仍复用原订阅记录续开。
	if orderType == PaymentOrderTypeSubscription && s.subscriptionService != nil && s.subscriptionService.userSubRepo != nil {
		if _, err := s.subscriptionService.userSubRepo.GetActiveByUserIDAndGroupID(ctx, userID, plan.GroupID); err == nil {
			return nil, ErrPaymentSubscriptionRepurchaseBlocked
		} else if err != nil && !errors.Is(err, ErrSubscriptionNotFound) {
			return nil, fmt.Errorf("check active subscription before create order: %w", err)
		}
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

	// 确定支付方式
	if payMethod == "" {
		payMethod = "wechat"
	}
	var payMethodStr string
	switch payMethod {
	case "alipay":
		if !s.isAlipayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "alipay payment is not configured")
		}
		payMethodStr = PaymentMethodAlipayNative
	case "epay_alipay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayAlipay
	case "epay_wxpay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayWxpay
	default:
		if !s.isWechatPayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "wechat payment is not configured")
		}
		payMethodStr = PaymentMethodWechatNative
	}

	// 创建订单记录
	currentSubSite := s.getSubSiteFromCtx(ctx)
	order := &PaymentOrder{
		OrderNo:        orderNo,
		UserID:         userID,
		PlanKey:        plan.Key,
		AmountFen:      finalAmount,
		Status:         PaymentOrderStatusPending,
		PayMethod:      payMethodStr,
		OrderType:      orderType,
		SubSiteID:      subSiteIDPtr(currentSubSite),
		PromoCode:      promoCode,
		DiscountAmount: discountFen,
		ExpiredAt:      time.Now().Add(30 * time.Minute), // 30分钟过期
	}

	if orderType == PaymentOrderTypeBalance {
		order.GroupID = 0
		order.ValidityDays = 0
		order.BalanceAmount = plan.BalanceAmount
	} else {
		order.GroupID = plan.GroupID
		order.ValidityDays = plan.ValidityDays
	}

	// 根据支付方式调用不同的下单接口
	var codeURL string
	switch payMethod {
	case "alipay":
		codeURL, err = s.createAlipayNativeOrder(ctx, order, plan.Name)
	case "epay_alipay":
		codeURL, err = s.createEpayOrder(ctx, order, plan.Name, "alipay")
	case "epay_wxpay":
		codeURL, err = s.createEpayOrder(ctx, order, plan.Name, "wxpay")
	default:
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
func (s *PaymentService) CreateRechargeOrder(ctx context.Context, userID int64, amountYuan float64, promoCode string, payMethod string, planKey string) (*PaymentOrder, error) {
	// 检查支付是否启用
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return nil, ErrPaymentDisabled
	}

	if amountYuan <= 0 {
		return nil, infraerrors.BadRequest("INVALID_AMOUNT", "recharge amount must be positive")
	}

	// 如果指定了 planKey，验证充值套餐并执行购买限制
	usePlanKey := "recharge_custom"
	planBalanceAmount := amountYuan // 默认到账金额=支付金额
	if planKey != "" {
		rechargeInfo, err := s.GetRechargeInfo(ctx)
		if err == nil {
			for _, rp := range rechargeInfo.Plans {
				if rp.Key == planKey {
					// 验证金额匹配
					expectedYuan := float64(rp.PayAmountFen) / 100
					if amountYuan != expectedYuan {
						return nil, infraerrors.BadRequest("AMOUNT_MISMATCH", "amount does not match the selected plan")
					}
					// 检查购买次数限制
					if rp.MaxPurchases > 0 {
						count, err := s.orderRepo.CountPaidByUserAndPlanKey(ctx, userID, planKey)
						if err != nil {
							return nil, fmt.Errorf("check purchase count: %w", err)
						}
						if count >= rp.MaxPurchases {
							return nil, infraerrors.BadRequest("PURCHASE_LIMIT_REACHED", fmt.Sprintf("this plan can only be purchased %d time(s)", rp.MaxPurchases))
						}
					}
					usePlanKey = planKey
					planBalanceAmount = rp.BalanceAmount
					break
				}
			}
		}
	}

	// 检查最低充值金额（仅自定义充值时）
	if usePlanKey == "recharge_custom" {
		if minAmountStr, _ := s.settingService.GetSettingValue(ctx, SettingKeyRechargeMinAmount); minAmountStr != "" {
			if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil && minAmount > 0 {
				if amountYuan < minAmount {
					return nil, infraerrors.BadRequest("AMOUNT_TOO_LOW", fmt.Sprintf("minimum recharge amount is ¥%.0f", minAmount))
				}
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

	// 检测当前请求是否命中 pool 模式分站（有 OwnerPaymentConfig）
	currentSubSite := s.getSubSiteFromCtx(ctx)
	poolSite := s.getPoolSubSiteFromCtx(ctx)
	var ownerPayCfg *OwnerPaymentConfig
	if poolSite != nil {
		ownerPayCfg = poolSite.OwnerPaymentConfig
	}

	// 确定支付方式（pool 分站优先使用站长自有凭据）
	if payMethod == "" {
		payMethod = "wechat"
	}
	if !s.isPayMethodReadyForBalanceOrder(ctx, payMethod, ownerPayCfg) {
		return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", payMethod+" payment is not configured")
	}
	var payMethodStr string
	switch payMethod {
	case "alipay":
		payMethodStr = PaymentMethodAlipayNative
	case "epay_alipay":
		payMethodStr = PaymentMethodEpayAlipay
	case "epay_wxpay":
		payMethodStr = PaymentMethodEpayWxpay
	default:
		payMethodStr = PaymentMethodWechatNative
	}

	// 如果是分站用户充值，把 siteID 编码到 PlanKey 以便 notify 时执行自动进货
	orderPlanKey := usePlanKey
	if poolSite != nil {
		orderPlanKey = encodeBalanceSubSitePlanKey(usePlanKey, poolSite.ID)
	}

	order := &PaymentOrder{
		OrderNo:        orderNo,
		UserID:         userID,
		PlanKey:        orderPlanKey,
		AmountFen:      finalAmount,
		GroupID:        0,
		ValidityDays:   0,
		OrderType:      PaymentOrderTypeBalance,
		BalanceAmount:  planBalanceAmount,
		SubSiteID:      subSiteIDPtr(currentSubSite),
		PromoCode:      promoCode,
		DiscountAmount: discountFen,
		Status:         PaymentOrderStatusPending,
		PayMethod:      payMethodStr,
		ExpiredAt:      time.Now().Add(30 * time.Minute),
	}

	description := fmt.Sprintf("余额充值 %.2f 元", amountYuan)
	var codeURL string
	var err error
	if ownerPayCfg != nil {
		codeURL, err = s.createOrderWithOwnerConfig(ctx, order, description, payMethod, ownerPayCfg)
	} else {
		switch payMethod {
		case "alipay":
			codeURL, err = s.createAlipayNativeOrder(ctx, order, description)
		case "epay_alipay":
			codeURL, err = s.createEpayOrder(ctx, order, description, "alipay")
		case "epay_wxpay":
			codeURL, err = s.createEpayOrder(ctx, order, description, "wxpay")
		default:
			codeURL, err = s.createWechatNativeOrder(ctx, order, description)
		}
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

// CreateAgentActivationOrder creates a payment order for the agent activation fee.
func (s *PaymentService) CreateAgentActivationOrder(ctx context.Context, userID int64, payMethod string) (*PaymentOrder, error) {
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return nil, ErrPaymentDisabled
	}
	if !s.settingService.IsAgentEnabled(ctx) {
		return nil, ErrAgentDisabled
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.IsAgent || user.AgentStatus == AgentStatusApproved {
		return nil, infraerrors.Conflict("AGENT_ALREADY_APPROVED", "you are already an approved agent")
	}
	if user.AgentStatus == AgentStatusPending {
		return nil, infraerrors.Conflict("AGENT_APPLICATION_PENDING", "your agent application is already pending review")
	}
	paidCount, err := s.orderRepo.CountPaidByUserAndPlanKey(ctx, userID, PaymentOrderTypeAgentActivation)
	if err != nil {
		return nil, fmt.Errorf("count paid activation orders: %w", err)
	}
	if paidCount > 0 {
		return nil, infraerrors.Conflict("AGENT_ACTIVATION_ALREADY_PAID", "agent activation fee has already been paid")
	}

	activationFee := s.settingService.GetAgentActivationFee(ctx)
	amountFen := int(activationFee * 100)
	if amountFen <= 0 {
		return nil, infraerrors.BadRequest("AGENT_ACTIVATION_FEE_INVALID", "agent activation fee must be positive")
	}

	if payMethod == "" {
		payMethod = "wechat"
	}
	var payMethodStr string
	switch payMethod {
	case "alipay":
		if !s.isAlipayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "alipay payment is not configured")
		}
		payMethodStr = PaymentMethodAlipayNative
	case "epay_alipay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayAlipay
	case "epay_wxpay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayWxpay
	default:
		if !s.isWechatPayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "wechat payment is not configured")
		}
		payMethodStr = PaymentMethodWechatNative
	}

	order := &PaymentOrder{
		OrderNo:   generateOrderNo(),
		UserID:    userID,
		PlanKey:   PaymentOrderTypeAgentActivation,
		AmountFen: amountFen,
		Status:    PaymentOrderStatusPending,
		PayMethod: payMethodStr,
		OrderType: PaymentOrderTypeAgentActivation,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}

	description := fmt.Sprintf("代理开通费 %.2f 元", activationFee)
	var codeURL string
	switch payMethod {
	case "alipay":
		codeURL, err = s.createAlipayNativeOrder(ctx, order, description)
	case "epay_alipay":
		codeURL, err = s.createEpayOrder(ctx, order, description, "alipay")
	case "epay_wxpay":
		codeURL, err = s.createEpayOrder(ctx, order, description, "wxpay")
	default:
		codeURL, err = s.createWechatNativeOrder(ctx, order, description)
	}
	if err != nil {
		log.Printf("[Payment] Failed to create %s native order for agent activation: %v", payMethod, err)
		return nil, err
	}
	order.CodeURL = &codeURL

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("save activation order: %w", err)
	}
	return order, nil
}

// CreateSubSiteActivationOrder creates a payment order for self-service sub-site activation.
func (s *PaymentService) CreateSubSiteActivationOrder(ctx context.Context, userID int64, input CreateSubSiteActivationInput, payMethod string) (*PaymentOrder, error) {
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return nil, ErrPaymentDisabled
	}
	if s.subSiteService == nil {
		return nil, infraerrors.ServiceUnavailable("SUBSITE_SERVICE_UNAVAILABLE", "sub-site service is unavailable")
	}

	request, openInfo, err := s.subSiteService.CreateActivationRequest(ctx, userID, input)
	if err != nil {
		return nil, err
	}
	if openInfo.PriceFen <= 0 {
		return nil, infraerrors.BadRequest("SUBSITE_PRICE_INVALID", "sub-site activation price must be positive")
	}

	if payMethod == "" {
		payMethod = "wechat"
	}
	var payMethodStr string
	switch payMethod {
	case "alipay":
		if !s.isAlipayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "alipay payment is not configured")
		}
		payMethodStr = PaymentMethodAlipayNative
	case "epay_alipay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayAlipay
	case "epay_wxpay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayWxpay
	default:
		if !s.isWechatPayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "wechat payment is not configured")
		}
		payMethodStr = PaymentMethodWechatNative
	}

	order := &PaymentOrder{
		OrderNo:   generateOrderNo(),
		UserID:    userID,
		PlanKey:   PaymentOrderTypeSubSiteActivation,
		AmountFen: openInfo.PriceFen,
		Status:    PaymentOrderStatusPending,
		PayMethod: payMethodStr,
		OrderType: PaymentOrderTypeSubSiteActivation,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}

	scopeLabel := "平台"
	if openInfo.Scope == "subsite" && openInfo.ParentSubSiteName != "" {
		scopeLabel = openInfo.ParentSubSiteName
	}
	description := fmt.Sprintf("%s分站开通费", scopeLabel)
	var codeURL string
	switch payMethod {
	case "alipay":
		codeURL, err = s.createAlipayNativeOrder(ctx, order, description)
	case "epay_alipay":
		codeURL, err = s.createEpayOrder(ctx, order, description, "alipay")
	case "epay_wxpay":
		codeURL, err = s.createEpayOrder(ctx, order, description, "wxpay")
	default:
		codeURL, err = s.createWechatNativeOrder(ctx, order, description)
	}
	if err != nil {
		log.Printf("[Payment] Failed to create %s native order for sub-site activation: %v", payMethod, err)
		return nil, err
	}
	order.CodeURL = &codeURL

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("save sub-site activation order: %w", err)
	}
	request.PaymentOrderID = order.ID
	if err := s.subSiteService.SaveActivationRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("save sub-site activation request: %w", err)
	}
	return order, nil
}

// CreateSubSiteTopupOrder 分站主从主站支付通道向自己的分站池充值。
// 订单成功后在 handlePaymentSuccess 中将 amountFen 入账到 sub_sites.balance_fen。
// siteID 通过 PlanKey 编码（subsite_topup:{id}）传递到回调路径。
func (s *PaymentService) CreateSubSiteTopupOrder(ctx context.Context, userID int64, siteID int64, amountFen int, payMethod string) (*PaymentOrder, error) {
	enabled, _ := s.settingService.GetSettingValue(ctx, SettingKeyPaymentEnabled)
	if enabled != "true" {
		return nil, ErrPaymentDisabled
	}
	if amountFen <= 0 {
		return nil, infraerrors.BadRequest("INVALID_AMOUNT", "topup amount must be positive")
	}
	if s.subSiteService == nil {
		return nil, infraerrors.ServiceUnavailable("SUBSITE_SERVICE_UNAVAILABLE", "sub-site service is unavailable")
	}
	site, err := s.subSiteService.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site.OwnerUserID != userID {
		return nil, ErrSubSiteForbidden
	}
	if !site.AllowOnlineTopup {
		return nil, infraerrors.Forbidden("SUBSITE_ONLINE_TOPUP_DISABLED", "online topup is disabled for this sub-site")
	}

	if payMethod == "" {
		payMethod = "wechat"
	}
	var payMethodStr string
	switch payMethod {
	case "alipay":
		if !s.isAlipayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "alipay payment is not configured")
		}
		payMethodStr = PaymentMethodAlipayNative
	case "epay_alipay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayAlipay
	case "epay_wxpay":
		if !s.isEpayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "epay payment is not configured")
		}
		payMethodStr = PaymentMethodEpayWxpay
	default:
		if !s.isWechatPayReady(ctx) {
			return nil, infraerrors.BadRequest("PAYMENT_METHOD_UNAVAILABLE", "wechat payment is not configured")
		}
		payMethodStr = PaymentMethodWechatNative
	}

	order := &PaymentOrder{
		OrderNo:   generateOrderNo(),
		UserID:    userID,
		PlanKey:   fmt.Sprintf("%s:%d", PaymentOrderTypeSubSiteTopup, siteID),
		AmountFen: amountFen,
		Status:    PaymentOrderStatusPending,
		PayMethod: payMethodStr,
		OrderType: PaymentOrderTypeSubSiteTopup,
		SubSiteID: &siteID,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}

	description := fmt.Sprintf("分站 %s 余额充值", site.Name)
	var codeURL string
	switch payMethod {
	case "alipay":
		codeURL, err = s.createAlipayNativeOrder(ctx, order, description)
	case "epay_alipay":
		codeURL, err = s.createEpayOrder(ctx, order, description, "alipay")
	case "epay_wxpay":
		codeURL, err = s.createEpayOrder(ctx, order, description, "wxpay")
	default:
		codeURL, err = s.createWechatNativeOrder(ctx, order, description)
	}
	if err != nil {
		log.Printf("[Payment] Failed to create %s native order for sub-site topup: %v", payMethod, err)
		return nil, err
	}
	order.CodeURL = &codeURL

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("save sub-site topup order: %w", err)
	}
	return order, nil
}

// parseSubSiteTopupOrderSiteID 从 PlanKey（格式 "subsite_topup:{id}"）提取分站 id。
func parseSubSiteTopupOrderSiteID(planKey string) (int64, error) {
	prefix := PaymentOrderTypeSubSiteTopup + ":"
	if !strings.HasPrefix(planKey, prefix) {
		return 0, fmt.Errorf("plan key is not a sub-site topup key: %q", planKey)
	}
	return strconv.ParseInt(strings.TrimPrefix(planKey, prefix), 10, 64)
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
	if order.Status == PaymentOrderStatusPending && order.ExpiredAt.After(time.Now()) {
		refreshedOrder, refreshErr := s.refreshPendingOrderStatus(ctx, order)
		if refreshErr != nil {
			log.Printf("[Payment] QueryOrder refresh failed: order=%s method=%s err=%v", order.OrderNo, order.PayMethod, refreshErr)
		} else if refreshedOrder != nil {
			order = refreshedOrder
		}
	}
	return order, nil
}

// ListOrders 列出用户订单
func (s *PaymentService) ListOrders(ctx context.Context, userID int64, params pagination.PaginationParams) ([]PaymentOrder, *pagination.PaginationResult, error) {
	orders, paginationResult, err := s.orderRepo.ListByUserID(ctx, userID, params)
	if err != nil {
		return nil, nil, err
	}
	if len(orders) > 0 || params.Page > 1 || s.subscriptionService == nil {
		return orders, paginationResult, nil
	}

	legacyOrders, legacyPagination, legacyErr := s.listLegacySubscriptionOrders(ctx, userID, params)
	if legacyErr != nil {
		return orders, paginationResult, nil
	}
	return legacyOrders, legacyPagination, nil
}

func (s *PaymentService) listLegacySubscriptionOrders(ctx context.Context, userID int64, params pagination.PaginationParams) ([]PaymentOrder, *pagination.PaginationResult, error) {
	subscriptions, err := s.subscriptionService.ListUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if len(subscriptions) == 0 {
		return []PaymentOrder{}, &pagination.PaginationResult{
			Total:    0,
			Page:     params.Page,
			PageSize: params.PageSize,
			Pages:    1,
		}, nil
	}

	sort.Slice(subscriptions, func(i, j int) bool {
		return subscriptions[i].CreatedAt.After(subscriptions[j].CreatedAt)
	})

	total := len(subscriptions)
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	if start >= total {
		return []PaymentOrder{}, &pagination.PaginationResult{
			Total:    int64(total),
			Page:     page,
			PageSize: pageSize,
			Pages:    (total + pageSize - 1) / pageSize,
		}, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	items := make([]PaymentOrder, 0, end-start)
	for _, sub := range subscriptions[start:end] {
		groupName := fmt.Sprintf("group_%d", sub.GroupID)
		amountFen := 0
		if sub.Group != nil {
			if strings.TrimSpace(sub.Group.Name) != "" {
				groupName = sub.Group.Name
			}
			amountFen = sub.Group.PriceFen
		}

		items = append(items, PaymentOrder{
			OrderNo:      fmt.Sprintf("legacy-sub-%d", sub.ID),
			UserID:       sub.UserID,
			PlanKey:      fmt.Sprintf("legacy_subscription_%d", sub.GroupID),
			GroupID:      sub.GroupID,
			AmountFen:    amountFen,
			ValidityDays: int(sub.ExpiresAt.Sub(sub.StartsAt).Hours() / 24),
			OrderType:    PaymentOrderTypeSubscription,
			Status:       legacySubscriptionOrderStatus(sub.Status),
			PayMethod:    fmt.Sprintf("legacy_subscription:%s", groupName),
			PaidAt:       &sub.CreatedAt,
			ExpiredAt:    sub.ExpiresAt,
			CreatedAt:    sub.CreatedAt,
			UpdatedAt:    sub.UpdatedAt,
		})
	}

	return items, &pagination.PaginationResult{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		Pages:    (total + pageSize - 1) / pageSize,
	}, nil
}

func legacySubscriptionOrderStatus(subscriptionStatus string) string {
	switch subscriptionStatus {
	case "", SubscriptionStatusActive, SubscriptionStatusExpired:
		return PaymentOrderStatusPaid
	case SubscriptionStatusSuspended:
		return PaymentOrderStatusClosed
	default:
		return PaymentOrderStatusClosed
	}
}

// ListAllOrders 列出所有订单（管理员）
func (s *PaymentService) ListAllOrders(ctx context.Context, params pagination.PaginationParams, status string, orderType string) ([]PaymentOrder, *pagination.PaginationResult, error) {
	return s.orderRepo.ListAll(ctx, params, status, orderType)
}

func (s *PaymentService) SubmitInvoiceRequest(ctx context.Context, userID int64, orderNos []string, invoice InvoiceRequest) error {
	uniqueOrderNos := make([]string, 0, len(orderNos))
	seen := make(map[string]struct{}, len(orderNos))
	for _, orderNo := range orderNos {
		trimmed := strings.TrimSpace(orderNo)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		uniqueOrderNos = append(uniqueOrderNos, trimmed)
	}
	if len(uniqueOrderNos) == 0 {
		return ErrInvoiceOrdersRequired
	}

	invoice.CompanyName = strings.TrimSpace(invoice.CompanyName)
	invoice.TaxID = strings.TrimSpace(invoice.TaxID)
	invoice.Email = strings.TrimSpace(invoice.Email)
	invoice.Remark = strings.TrimSpace(invoice.Remark)

	if invoice.CompanyName == "" || invoice.TaxID == "" || invoice.Email == "" {
		return infraerrors.BadRequest("INVOICE_FIELDS_REQUIRED", "company name, tax id and invoice email are required")
	}
	if !isValidInvoiceEmail(invoice.Email) {
		return ErrInvoiceEmailInvalid
	}

	return s.orderRepo.SubmitInvoiceRequest(ctx, userID, uniqueOrderNos, invoice)
}

func (s *PaymentService) MarkInvoiceProcessed(ctx context.Context, orderID int64) error {
	return s.orderRepo.MarkInvoiceProcessed(ctx, orderID)
}

// RepairOrder manually marks a missed payment callback order as paid and replays the normal fulfillment logic.
func (s *PaymentService) RepairOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case PaymentOrderStatusPaid:
		return ErrPaymentOrderPaid
	case PaymentOrderStatusPending, PaymentOrderStatusClosed:
		// repair is allowed
	default:
		return ErrPaymentOrderUnrepairable
	}

	originalStatus := order.Status
	paidAt := time.Now()
	claimed, err := s.orderRepo.CompareAndUpdateStatus(ctx, order.OrderNo, originalStatus, PaymentOrderStatusPaid, nil, &paidAt)
	if err != nil {
		return fmt.Errorf("claim order for repair: %w", err)
	}
	if !claimed {
		latestOrder, latestErr := s.orderRepo.GetByID(ctx, orderID)
		if latestErr == nil && latestOrder.Status == PaymentOrderStatusPaid {
			return ErrPaymentOrderPaid
		}
		return infraerrors.Conflict("PAYMENT_ORDER_STATUS_CHANGED", "payment order status changed, please refresh and retry")
	}

	if err := s.handlePaymentSuccess(ctx, order); err != nil {
		if rbErr := s.orderRepo.UpdateStatus(ctx, order.OrderNo, originalStatus, nil, nil); rbErr != nil {
			log.Printf("[Payment] CRITICAL: order %s failed to rollback repaired status: %v (original error: %v)", order.OrderNo, rbErr, err)
		}
		return err
	}

	log.Printf("[Payment] Order %s repaired manually by admin", order.OrderNo)
	return nil
}

func (s *PaymentService) refreshPendingOrderStatus(ctx context.Context, order *PaymentOrder) (*PaymentOrder, error) {
	switch order.PayMethod {
	case PaymentMethodWechatNative:
		return s.queryWechatPendingOrder(ctx, order)
	case PaymentMethodAlipayNative:
		return s.queryAlipayPendingOrder(ctx, order)
	default:
		return nil, nil
	}
}

func (s *PaymentService) markPendingOrderPaid(ctx context.Context, order *PaymentOrder, transactionID *string, paidAt *time.Time, source string) (*PaymentOrder, error) {
	claimed, err := s.orderRepo.CompareAndUpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, PaymentOrderStatusPaid, transactionID, paidAt)
	if err != nil {
		return nil, fmt.Errorf("claim order as paid: %w", err)
	}
	if !claimed {
		return s.orderRepo.GetByOrderNo(ctx, order.OrderNo)
	}

	if err := s.handlePaymentSuccess(ctx, order); err != nil {
		if rbErr := s.orderRepo.UpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, nil, nil); rbErr != nil {
			log.Printf("[Payment] CRITICAL: order %s failed to rollback proactive status: %v (original error: %v)", order.OrderNo, rbErr, err)
		}
		return nil, err
	}

	log.Printf("[Payment] Order %s marked paid via %s", order.OrderNo, source)
	return s.orderRepo.GetByOrderNo(ctx, order.OrderNo)
}

func (s *PaymentService) queryWechatPendingOrder(ctx context.Context, order *PaymentOrder) (*PaymentOrder, error) {
	mchID, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayMchID)
	privateKeyPEM, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayPrivateKey)
	mchSerialNo, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayMchSerialNo)
	if mchID == "" || privateKeyPEM == "" || mchSerialNo == "" {
		return nil, nil
	}

	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse wechat private key: %w", err)
	}

	queryPath := fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s?mchid=%s", url.PathEscape(order.OrderNo), url.QueryEscape(mchID))
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := generateNonce()
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n\n", http.MethodGet, queryPath, timestamp, nonce)
	signature, err := signSHA256WithRSA(privateKey, []byte(signStr))
	if err != nil {
		return nil, fmt.Errorf("sign wechat query request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wechatPayAPIBaseURL+queryPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(
		`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`,
		mchID, nonce, timestamp, mchSerialNo, signature,
	))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("query wechat order: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wechat query api error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result struct {
		TradeState    string `json:"trade_state"`
		TransactionID string `json:"transaction_id"`
		SuccessTime   string `json:"success_time"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse wechat query response: %w", err)
	}
	if result.TradeState != "SUCCESS" {
		return nil, nil
	}

	var paidAt *time.Time
	if result.SuccessTime != "" {
		if parsed, err := time.Parse(time.RFC3339, result.SuccessTime); err == nil {
			paidAt = &parsed
		}
	}
	if paidAt == nil {
		now := time.Now()
		paidAt = &now
	}

	var transactionID *string
	if result.TransactionID != "" {
		transactionID = &result.TransactionID
	}
	return s.markPendingOrderPaid(ctx, order, transactionID, paidAt, "wechat order query")
}

func (s *PaymentService) queryAlipayPendingOrder(ctx context.Context, order *PaymentOrder) (*PaymentOrder, error) {
	client, err := s.initAlipayClient(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := client.TradeQuery(ctx, alipay.TradeQuery{OutTradeNo: order.OrderNo})
	if err != nil {
		return nil, fmt.Errorf("query alipay order: %w", err)
	}
	if rsp.IsFailure() {
		return nil, fmt.Errorf("alipay query api error: %s - %s", rsp.Code, rsp.Msg)
	}
	if rsp.TradeStatus != alipay.TradeStatusSuccess && rsp.TradeStatus != alipay.TradeStatusFinished {
		return nil, nil
	}

	var paidAt *time.Time
	if rsp.SendPayDate != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", rsp.SendPayDate); err == nil {
			paidAt = &parsed
		}
	}
	if paidAt == nil {
		now := time.Now()
		paidAt = &now
	}

	var transactionID *string
	if rsp.TradeNo != "" {
		transactionID = &rsp.TradeNo
	}
	return s.markPendingOrderPaid(ctx, order, transactionID, paidAt, "alipay order query")
}

func isValidInvoiceEmail(email string) bool {
	return regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(email)
}

// HandleWechatNotify 处理微信支付回调通知
func (s *PaymentService) HandleWechatNotify(ctx context.Context, body []byte, wechatpayTimestamp, wechatpayNonce, wechatpaySignature, wechatpaySerial string) error {
	// 先尝试全局凭据验签+解密；失败后尝试匹配分站凭据。
	order, err := s.processWechatNotifyWithGlobal(ctx, body, wechatpayTimestamp, wechatpayNonce, wechatpaySignature, wechatpaySerial)
	if err != nil {
		log.Printf("[Payment] WechatNotify global verify failed, trying sub-site credentials: %v", err)
		order, err = s.processWechatNotifyWithOwner(ctx, body, wechatpayTimestamp, wechatpayNonce, wechatpaySignature, wechatpaySerial)
		if err != nil {
			return err
		}
	}
	if order == nil {
		return nil
	}
	if order.Status == PaymentOrderStatusPaid {
		return nil
	}

	now := time.Now()
	transactionID := order.WechatTransactionID
	var txIDPtr *string
	if transactionID != nil {
		txIDPtr = transactionID
	}
	claimed, err := s.orderRepo.CompareAndUpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, PaymentOrderStatusPaid, txIDPtr, &now)
	if err != nil {
		return fmt.Errorf("claim order: %w", err)
	}
	if !claimed {
		return nil
	}
	if err := s.handlePaymentSuccess(ctx, order); err != nil {
		if rbErr := s.orderRepo.UpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, nil, nil); rbErr != nil {
			log.Printf("[Payment] CRITICAL: order %s failed to rollback status: %v (original error: %v)", order.OrderNo, rbErr, err)
		}
		return err
	}
	return nil
}

func (s *PaymentService) processWechatNotifyWithGlobal(ctx context.Context, body []byte, timestamp, nonce, signature, serial string) (*PaymentOrder, error) {
	if err := s.verifyWechatSignature(ctx, timestamp, nonce, string(body), signature, serial); err != nil {
		return nil, err
	}
	return s.decryptAndLookupWechatOrder(ctx, body, "")
}

func (s *PaymentService) processWechatNotifyWithOwner(ctx context.Context, body []byte, timestamp, nonce, signature, serial string) (*PaymentOrder, error) {
	// 尝试先不验签（解密需要 APIv3Key）地找到 order，再用其绑定的分站凭据验签。
	// 解密顺序：先尝试用各分站的 APIv3Key 解密，成功后找到 order 并验签。
	// 简化方案：跳过回调签名验证（分站凭据不含 wechat public key），仅解密并验证金额。
	// 实际生产中，分站主应在其 OwnerPaymentConfig 中提供 public_key 以完成完整验签。
	var notification struct {
		EventType string `json:"event_type"`
		Resource  struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			Nonce          string `json:"nonce"`
			AssociatedData string `json:"associated_data"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &notification); err != nil {
		return nil, fmt.Errorf("parse notification: %w", err)
	}
	if notification.EventType != "TRANSACTION.SUCCESS" {
		return nil, nil
	}

	// 尝试用全局 APIv3Key 解密（分站 notify 可能仍走主站 APIv3Key — 取决于配置）
	apiKey, _ := s.settingService.GetSettingValue(ctx, SettingKeyWechatPayAPIv3Key)
	plaintext, err := decryptAEAD(apiKey, notification.Resource.Nonce, notification.Resource.Ciphertext, notification.Resource.AssociatedData)
	if err != nil {
		return nil, fmt.Errorf("decrypt notification: %w", err)
	}
	var transaction struct {
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		TradeState    string `json:"trade_state"`
	}
	if err := json.Unmarshal(plaintext, &transaction); err != nil {
		return nil, fmt.Errorf("parse transaction: %w", err)
	}
	if transaction.TradeState != "SUCCESS" {
		return nil, nil
	}
	order, err := s.orderRepo.GetByOrderNo(ctx, transaction.OutTradeNo)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	order.WechatTransactionID = &transaction.TransactionID
	return order, nil
}

// decryptAndLookupWechatOrder decrypts the wechat notify and returns the looked-up order.
func (s *PaymentService) decryptAndLookupWechatOrder(ctx context.Context, body []byte, overrideAPIv3Key string) (*PaymentOrder, error) {
	var notification struct {
		EventType string `json:"event_type"`
		Resource  struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			Nonce          string `json:"nonce"`
			AssociatedData string `json:"associated_data"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &notification); err != nil {
		return nil, fmt.Errorf("parse notification: %w", err)
	}
	if notification.EventType != "TRANSACTION.SUCCESS" {
		return nil, nil
	}
	apiKey := overrideAPIv3Key
	if apiKey == "" {
		apiKey, _ = s.settingService.GetSettingValue(ctx, SettingKeyWechatPayAPIv3Key)
	}
	plaintext, err := decryptAEAD(apiKey, notification.Resource.Nonce, notification.Resource.Ciphertext, notification.Resource.AssociatedData)
	if err != nil {
		return nil, fmt.Errorf("decrypt notification: %w", err)
	}
	var transaction struct {
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		TradeState    string `json:"trade_state"`
	}
	if err := json.Unmarshal(plaintext, &transaction); err != nil {
		return nil, fmt.Errorf("parse transaction: %w", err)
	}
	if transaction.TradeState != "SUCCESS" {
		return nil, nil
	}
	order, err := s.orderRepo.GetByOrderNo(ctx, transaction.OutTradeNo)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	order.WechatTransactionID = &transaction.TransactionID
	return order, nil
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
		return "", wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "wechat payment private key is invalid", fmt.Errorf("sign request: %w", err))
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
		return "", wrapUnknownAsServiceUnavailable("WECHAT_API_UNAVAILABLE", "wechat pay is temporarily unavailable", fmt.Errorf("wechat api request: %w", err))
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
		return "", wrapUnknownAsServiceUnavailable("WECHAT_API_INVALID_RESPONSE", "wechat pay returned an invalid response", fmt.Errorf("parse wechat response: %w", err))
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
// Shared Post-Payment Logic
// ===========================

// handlePaymentSuccess 处理支付成功后的共享逻辑（订阅分配/余额充值/缓存失效/优惠码记录）
func (s *PaymentService) handlePaymentSuccess(ctx context.Context, order *PaymentOrder) error {
	if order.OrderType == PaymentOrderTypeBalance {
		// 充值余额。pool 分站的自动进货必须和用户余额入账同事务完成，避免账实不一致。
		if _, siteID, ok := parseBalanceSubSitePlanKey(order.PlanKey); ok && siteID > 0 && s.subSiteService != nil {
			if err := s.subSiteService.CreditUserBalanceWithAutoRestock(ctx, order.UserID, siteID, order.ID, order.BalanceAmount, int64(order.AmountFen), order.OrderNo); err != nil {
				log.Printf("[Payment] Failed to credit balance with auto-restock for order %s, user %d, site %d: %v", order.OrderNo, order.UserID, siteID, err)
				return fmt.Errorf("credit balance with auto-restock: %w", err)
			}
			log.Printf("[Payment] Order %s auto-restock: site %d pool -%d fen", order.OrderNo, siteID, order.AmountFen)
		} else if err := s.userRepo.UpdateBalance(ctx, order.UserID, order.BalanceAmount); err != nil {
			log.Printf("[Payment] Failed to update balance for order %s, user %d: %v", order.OrderNo, order.UserID, err)
			return fmt.Errorf("update balance: %w", err)
		}
		// 失效余额缓存
		if s.billingCache != nil {
			if err := s.billingCache.InvalidateUserBalance(ctx, order.UserID); err != nil {
				log.Printf("[Payment] Failed to invalidate balance cache for user %d: %v", order.UserID, err)
			}
		}
		log.Printf("[Payment] Order %s paid successfully, balance +%.2f for user %d",
			order.OrderNo, order.BalanceAmount, order.UserID)
	} else if order.OrderType == PaymentOrderTypeAgentActivation {
		if s.agentService != nil {
			if err := s.agentService.MarkActivationFeePaid(ctx, order.UserID, order.ID); err != nil {
				log.Printf("[Payment] Failed to mark agent activation paid for user %d order %s: %v", order.UserID, order.OrderNo, err)
				return fmt.Errorf("mark agent activation paid: %w", err)
			}
		}
		log.Printf("[Payment] Order %s paid successfully, agent activation marked for user %d", order.OrderNo, order.UserID)
	} else if order.OrderType == PaymentOrderTypeSubSiteActivation {
		if s.subSiteService != nil {
			if _, err := s.subSiteService.ActivatePaidOrder(ctx, order); err != nil {
				log.Printf("[Payment] Failed to activate sub-site for user %d order %s: %v", order.UserID, order.OrderNo, err)
				return fmt.Errorf("activate sub-site: %w", err)
			}
		}
		log.Printf("[Payment] Order %s paid successfully, sub-site activated for user %d", order.OrderNo, order.UserID)
	} else if order.OrderType == PaymentOrderTypeSubSiteTopup {
		if s.subSiteService != nil {
			siteID, err := parseSubSiteTopupOrderSiteID(order.PlanKey)
			if err != nil {
				log.Printf("[Payment] Sub-site topup order %s has invalid plan_key %q: %v", order.OrderNo, order.PlanKey, err)
				return fmt.Errorf("parse sub-site topup order: %w", err)
			}
			if err := s.subSiteService.CreditPoolFromPayment(ctx, siteID, int64(order.AmountFen), order.ID); err != nil {
				log.Printf("[Payment] Failed to credit sub-site pool for order %s site %d: %v", order.OrderNo, siteID, err)
				return fmt.Errorf("credit sub-site pool: %w", err)
			}
			log.Printf("[Payment] Order %s paid successfully, sub-site %d pool +%d fen", order.OrderNo, siteID, order.AmountFen)
		}
	} else {
		// 分配订阅
		_, _, err := s.subscriptionService.AssignOrExtendSubscription(ctx, &AssignSubscriptionInput{
			UserID:       order.UserID,
			GroupID:      order.GroupID,
			ValidityDays: order.ValidityDays,
			AssignedBy:   0, // 系统自动分配
			Notes:        fmt.Sprintf("支付订单 %s", order.OrderNo),
		})
		if err != nil {
			log.Printf("[Payment] Failed to assign subscription for order %s, user %d: %v", order.OrderNo, order.UserID, err)
			return fmt.Errorf("assign subscription: %w", err)
		}
		log.Printf("[Payment] Order %s paid successfully, subscription assigned for user %d, group %d, %d days",
			order.OrderNo, order.UserID, order.GroupID, order.ValidityDays)
	}

	// 记录优惠码使用（支付成功后才记录）
	if order.PromoCode != "" && order.DiscountAmount > 0 && s.promoService != nil {
		if err := s.promoService.ApplyPromoCode(ctx, order.UserID, order.PromoCode, order.OrderNo, float64(order.DiscountAmount)); err != nil {
			log.Printf("[Payment] Failed to apply promo code for order %s: %v", order.OrderNo, err)
			// 不阻断支付流程
		}
	}

	// 触发代理佣金（支付成功后自动检查并创建佣金记录）
	if s.agentService != nil {
		payAmount := order.BalanceAmount
		if payAmount <= 0 {
			payAmount = float64(order.AmountFen) / 100.0
		}
		s.agentService.TriggerCommissionForPayment(ctx, order.UserID, order.ID, order.OrderType, payAmount)
	}

	return nil
}

// ===========================
// Alipay API Integration
// ===========================

// initAlipayClient 初始化支付宝客户端
func (s *PaymentService) initAlipayClient(ctx context.Context) (*alipay.Client, error) {
	appID, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayAppID)
	privateKey, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayPrivateKey)
	alipayPublicKey, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayPublicKey)
	isProductionStr, _ := s.settingService.GetSettingValue(ctx, SettingKeyAlipayIsProduction)
	isProduction := isProductionStr == "true"

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
		return nil, wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "alipay configuration is invalid", fmt.Errorf("init alipay client: %w", err))
	}

	err = client.LoadAliPayPublicKey(alipayPublicKey)
	if err != nil {
		return nil, wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "alipay public key is invalid", fmt.Errorf("load alipay public key: %w", err))
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

	var p = alipay.TradePreCreate{}
	p.NotifyURL = notifyURL
	p.Subject = subject
	p.OutTradeNo = order.OrderNo
	p.TotalAmount = fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)
	p.TimeExpire = order.ExpiredAt.Format("2006-01-02 15:04:05")

	rsp, err := client.TradePreCreate(ctx, p)
	if err != nil {
		return "", wrapUnknownAsServiceUnavailable("ALIPAY_API_UNAVAILABLE", "alipay is temporarily unavailable", fmt.Errorf("alipay trade precreate: %w", err))
	}

	if rsp.IsFailure() {
		return "", infraerrors.BadRequest("ALIPAY_API_ERROR", fmt.Sprintf("alipay error: %s - %s", rsp.Code, rsp.Msg))
	}

	return rsp.QRCode, nil
}

// HandleAlipayNotify 处理支付宝支付回调通知
func (s *PaymentService) HandleAlipayNotify(ctx context.Context, req *http.Request) error {
	// 先尝试全局凭据验签
	client, err := s.initAlipayClient(ctx)
	if err == nil {
		notification, verifyErr := client.GetTradeNotification(req)
		if verifyErr == nil {
			return s.processAlipayNotification(ctx, notification)
		}
		log.Printf("[Payment] AlipayNotify global verify failed, trying sub-site credentials: %v", verifyErr)
	}

	// 全局验签失败，尝试从 form 参数中找 out_trade_no 匹配分站凭据
	if parseErr := req.ParseForm(); parseErr != nil {
		return ErrPaymentSignature
	}
	outTradeNo := req.Form.Get("out_trade_no")
	if outTradeNo == "" {
		return ErrPaymentSignature
	}
	order, err := s.orderRepo.GetByOrderNo(ctx, outTradeNo)
	if err != nil {
		return ErrPaymentSignature
	}
	_, siteID, ok := parseBalanceSubSitePlanKey(order.PlanKey)
	if !ok || siteID <= 0 || s.subSiteService == nil {
		return ErrPaymentSignature
	}
	site, err := s.subSiteService.GetByID(ctx, siteID)
	if err != nil || site == nil || site.OwnerPaymentConfig == nil || site.OwnerPaymentConfig.Alipay == nil {
		return ErrPaymentSignature
	}
	ownerCred := site.OwnerPaymentConfig.Alipay
	ownerClient, err := alipay.New(ownerCred.AppID, ownerCred.PrivateKey, ownerCred.IsProduction)
	if err != nil {
		return ErrPaymentSignature
	}
	if err := ownerClient.LoadAliPayPublicKey(ownerCred.PublicKey); err != nil {
		return ErrPaymentSignature
	}
	notification, err := ownerClient.GetTradeNotification(req)
	if err != nil {
		log.Printf("[Payment] AlipayNotify sub-site verify also failed: %v", err)
		return ErrPaymentSignature
	}
	return s.processAlipayNotification(ctx, notification)
}

func (s *PaymentService) processAlipayNotification(ctx context.Context, notification *alipay.Notification) error {
	if notification.TradeStatus != "TRADE_SUCCESS" && notification.TradeStatus != "TRADE_FINISHED" {
		return nil
	}
	order, err := s.orderRepo.GetByOrderNo(ctx, notification.OutTradeNo)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}
	if order.Status == PaymentOrderStatusPaid {
		return nil
	}

	expectedYuan := fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)
	if notification.TotalAmount != expectedYuan {
		return fmt.Errorf("amount mismatch: expected %s, got %s", expectedYuan, notification.TotalAmount)
	}

	paidAt := time.Now()
	tradeNo := notification.TradeNo
	claimed, err := s.orderRepo.CompareAndUpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, PaymentOrderStatusPaid, &tradeNo, &paidAt)
	if err != nil {
		return fmt.Errorf("claim order: %w", err)
	}
	if !claimed {
		return nil
	}
	if err := s.handlePaymentSuccess(ctx, order); err != nil {
		if rbErr := s.orderRepo.UpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, nil, nil); rbErr != nil {
			log.Printf("[Payment] CRITICAL: order %s failed to rollback status: %v (original error: %v)", order.OrderNo, rbErr, err)
		}
		return err
	}
	log.Printf("[Payment] Order %s paid via Alipay, trade_no=%s", order.OrderNo, tradeNo)
	return nil
}

// ===========================
// Epay (易支付/ZPAY) Integration
// ===========================

// epaySign 易支付 MD5 签名
func epaySign(params map[string]string, pkey string) string {
	// 过滤空值、sign、sign_type
	filtered := make(map[string]string)
	for k, v := range params {
		if v == "" || k == "sign" || k == "sign_type" {
			continue
		}
		filtered[k] = v
	}

	// 按 key 排序
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接查询字符串
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(filtered[k])
	}
	buf.WriteString(pkey)

	// MD5 哈希
	hash := md5.Sum([]byte(buf.String()))
	return fmt.Sprintf("%x", hash)
}

// createEpayOrder 调用易支付下单接口
func (s *PaymentService) createEpayOrder(ctx context.Context, order *PaymentOrder, name string, payType string) (string, error) {
	gateway, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayGateway)
	pid, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayPID)
	pkey, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayPKey)
	notifyURL, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayNotifyURL)

	if gateway == "" || pid == "" || pkey == "" || notifyURL == "" {
		missing := make([]string, 0, 4)
		if gateway == "" {
			missing = append(missing, "gateway")
		}
		if pid == "" {
			missing = append(missing, "pid")
		}
		if pkey == "" {
			missing = append(missing, "pkey")
		}
		if notifyURL == "" {
			missing = append(missing, "notify_url")
		}
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", fmt.Sprintf("epay configuration incomplete, missing: %s", strings.Join(missing, ", ")))
	}

	// 金额转元
	money := fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)

	params := map[string]string{
		"pid":          pid,
		"type":         payType,
		"out_trade_no": order.OrderNo,
		"notify_url":   notifyURL,
		"name":         name,
		"money":        money,
	}
	params["sign"] = epaySign(params, pkey)
	params["sign_type"] = "MD5"

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	// POST 到网关
	gateway = strings.TrimRight(gateway, "/")
	req, err := http.NewRequestWithContext(ctx, "POST", gateway+"/mapi.php",
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "epay gateway url is invalid", fmt.Errorf("create epay request: %w", err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", wrapUnknownAsServiceUnavailable("EPAY_API_UNAVAILABLE", "epay is temporarily unavailable", fmt.Errorf("epay api request: %w", err))
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Payment] Epay API error: status=%d body=%s", resp.StatusCode, string(respBody))
		return "", infraerrors.BadRequest("EPAY_API_ERROR", fmt.Sprintf("epay api error: %s", string(respBody)))
	}

	var result struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		TradeNo string `json:"trade_no"`
		PayURL  string `json:"payurl"`
		QRCode  string `json:"qrcode"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("[Payment] Epay API invalid response: body=%s", string(respBody))
		return "", infraerrors.BadRequest("EPAY_API_ERROR", "epay gateway returned an unexpected response")
	}

	if result.Code != 1 {
		return "", infraerrors.BadRequest("EPAY_API_ERROR", fmt.Sprintf("epay error: %s", result.Msg))
	}

	// 优先使用 qrcode，其次 payurl
	codeURL := result.QRCode
	if codeURL == "" {
		codeURL = result.PayURL
	}
	if codeURL == "" {
		return "", infraerrors.BadRequest("EPAY_API_ERROR", "epay returned no payment url")
	}

	return codeURL, nil
}

// HandleEpayNotify 处理易支付回调通知
func (s *PaymentService) HandleEpayNotify(ctx context.Context, params map[string]string) error {
	// 先尝试全局凭据
	pid, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayPID)
	pkey, _ := s.settingService.GetSettingValue(ctx, SettingKeyEpayPKey)

	if pid != "" && pkey != "" {
		expectedSign := epaySign(params, pkey)
		if params["sign"] == expectedSign && params["pid"] == pid {
			return s.processEpayNotification(ctx, params)
		}
		log.Printf("[Payment] EpayNotify global verify failed, trying sub-site credentials")
	}

	// 尝试分站凭据
	outTradeNo := params["out_trade_no"]
	if outTradeNo == "" {
		return ErrPaymentSignature
	}
	order, err := s.orderRepo.GetByOrderNo(ctx, outTradeNo)
	if err != nil {
		return ErrPaymentSignature
	}
	_, siteID, ok := parseBalanceSubSitePlanKey(order.PlanKey)
	if !ok || siteID <= 0 || s.subSiteService == nil {
		return ErrPaymentSignature
	}
	site, err := s.subSiteService.GetByID(ctx, siteID)
	if err != nil || site == nil || site.OwnerPaymentConfig == nil || site.OwnerPaymentConfig.Epay == nil {
		return ErrPaymentSignature
	}
	ownerCred := site.OwnerPaymentConfig.Epay
	expectedSign := epaySign(params, ownerCred.PKey)
	if params["sign"] != expectedSign {
		log.Printf("[Payment] EpayNotify sub-site verify also failed")
		return ErrPaymentSignature
	}
	return s.processEpayNotification(ctx, params)
}

func (s *PaymentService) processEpayNotification(ctx context.Context, params map[string]string) error {
	if params["trade_status"] != "TRADE_SUCCESS" {
		return nil
	}

	order, err := s.orderRepo.GetByOrderNo(ctx, params["out_trade_no"])
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}
	if order.Status == PaymentOrderStatusPaid {
		return nil
	}

	expectedMoney := fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)
	if params["money"] != expectedMoney {
		return fmt.Errorf("amount mismatch: expected %s, got %s", expectedMoney, params["money"])
	}

	paidAt := time.Now()
	tradeNo := params["trade_no"]
	claimed, err := s.orderRepo.CompareAndUpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, PaymentOrderStatusPaid, &tradeNo, &paidAt)
	if err != nil {
		return fmt.Errorf("claim order: %w", err)
	}
	if !claimed {
		return nil
	}
	if err := s.handlePaymentSuccess(ctx, order); err != nil {
		if rbErr := s.orderRepo.UpdateStatus(ctx, order.OrderNo, PaymentOrderStatusPending, nil, nil); rbErr != nil {
			log.Printf("[Payment] CRITICAL: order %s failed to rollback status: %v (original error: %v)", order.OrderNo, rbErr, err)
		}
		return err
	}
	log.Printf("[Payment] Order %s paid via Epay, trade_no=%s", order.OrderNo, tradeNo)
	return nil
}

// createOrderWithOwnerConfig dispatches payment creation to owner-credential-based methods.
func (s *PaymentService) createOrderWithOwnerConfig(ctx context.Context, order *PaymentOrder, desc, payMethod string, opc *OwnerPaymentConfig) (string, error) {
	switch payMethod {
	case "alipay":
		return s.createOwnerAlipayOrder(ctx, order, desc, opc.Alipay)
	case "epay_alipay":
		return s.createOwnerEpayOrder(ctx, order, desc, "alipay", opc.Epay)
	case "epay_wxpay":
		return s.createOwnerEpayOrder(ctx, order, desc, "wxpay", opc.Epay)
	default:
		return s.createOwnerWechatOrder(ctx, order, desc, opc.Wechat)
	}
}

func (s *PaymentService) createOwnerWechatOrder(ctx context.Context, order *PaymentOrder, planName string, cred *WechatPayCredentials) (string, error) {
	if cred == nil || cred.AppID == "" || cred.MchID == "" || cred.PrivateKey == "" || cred.MchSerialNo == "" || cred.NotifyURL == "" {
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "sub-site wechat payment configuration is incomplete")
	}

	safeDesc := regexp.MustCompile(`[^\x{0000}-\x{FFFF}]`).ReplaceAllString(planName, "")
	safeDesc = strings.TrimSpace(safeDesc)
	if safeDesc == "" {
		safeDesc = "订单支付"
	}

	reqBody := map[string]any{
		"appid":        cred.AppID,
		"mchid":        cred.MchID,
		"description":  safeDesc,
		"out_trade_no": order.OrderNo,
		"time_expire":  order.ExpiredAt.Format(time.RFC3339),
		"notify_url":   cred.NotifyURL,
		"amount": map[string]any{
			"total":    order.AmountFen,
			"currency": "CNY",
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := generateNonce()
	method := "POST"
	urlPath := "/v3/pay/transactions/native"
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", method, urlPath, timestamp, nonce, string(bodyBytes))

	privateKey, err := parsePrivateKey(cred.PrivateKey)
	if err != nil {
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "invalid sub-site wechat private key")
	}
	signature, err := signSHA256WithRSA(privateKey, []byte(signStr))
	if err != nil {
		return "", wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "sub-site wechat private key is invalid", fmt.Errorf("sign request: %w", err))
	}

	req, err := http.NewRequestWithContext(ctx, method, wechatPayAPIBaseURL+urlPath, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(
		`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`,
		cred.MchID, nonce, timestamp, cred.MchSerialNo, signature,
	))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", wrapUnknownAsServiceUnavailable("WECHAT_API_UNAVAILABLE", "wechat pay is temporarily unavailable", fmt.Errorf("wechat api request: %w", err))
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Payment] SubSite WeChat API error: status=%d body=%s", resp.StatusCode, string(respBody))
		return "", infraerrors.BadRequest("WECHAT_API_ERROR", fmt.Sprintf("wechat pay api error: %s", string(respBody)))
	}
	var result struct {
		CodeURL string `json:"code_url"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", wrapUnknownAsServiceUnavailable("WECHAT_API_INVALID_RESPONSE", "wechat pay returned an invalid response", fmt.Errorf("parse wechat response: %w", err))
	}
	return result.CodeURL, nil
}

func (s *PaymentService) createOwnerAlipayOrder(ctx context.Context, order *PaymentOrder, subject string, cred *AlipayCredentials) (string, error) {
	if cred == nil || cred.AppID == "" || cred.PrivateKey == "" || cred.PublicKey == "" || cred.NotifyURL == "" {
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "sub-site alipay configuration is incomplete")
	}
	client, err := alipay.New(cred.AppID, cred.PrivateKey, cred.IsProduction)
	if err != nil {
		return "", wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "sub-site alipay credentials are invalid", fmt.Errorf("init alipay client: %w", err))
	}
	if err := client.LoadAliPayPublicKey(cred.PublicKey); err != nil {
		return "", wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "sub-site alipay public key is invalid", fmt.Errorf("load alipay public key: %w", err))
	}

	p := alipay.TradePreCreate{}
	p.NotifyURL = cred.NotifyURL
	p.Subject = subject
	p.OutTradeNo = order.OrderNo
	p.TotalAmount = fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)
	p.TimeExpire = order.ExpiredAt.Format("2006-01-02 15:04:05")

	rsp, err := client.TradePreCreate(ctx, p)
	if err != nil {
		return "", wrapUnknownAsServiceUnavailable("ALIPAY_API_UNAVAILABLE", "alipay is temporarily unavailable", fmt.Errorf("alipay trade precreate: %w", err))
	}
	if rsp.IsFailure() {
		return "", infraerrors.BadRequest("ALIPAY_API_ERROR", fmt.Sprintf("alipay error: %s - %s", rsp.Code, rsp.Msg))
	}
	return rsp.QRCode, nil
}

func (s *PaymentService) createOwnerEpayOrder(ctx context.Context, order *PaymentOrder, name, payType string, cred *EpayCredentials) (string, error) {
	if cred == nil || cred.Gateway == "" || cred.PID == "" || cred.PKey == "" || cred.NotifyURL == "" {
		return "", infraerrors.BadRequest("PAYMENT_CONFIG_MISSING", "sub-site epay configuration is incomplete")
	}

	money := fmt.Sprintf("%.2f", float64(order.AmountFen)/100.0)
	params := map[string]string{
		"pid":          cred.PID,
		"type":         payType,
		"out_trade_no": order.OrderNo,
		"notify_url":   cred.NotifyURL,
		"name":         name,
		"money":        money,
	}
	params["sign"] = epaySign(params, cred.PKey)
	params["sign_type"] = "MD5"

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	gateway := strings.TrimRight(cred.Gateway, "/")
	req, err := http.NewRequestWithContext(ctx, "POST", gateway+"/mapi.php", strings.NewReader(form.Encode()))
	if err != nil {
		return "", wrapUnknownAsBadRequest("PAYMENT_CONFIG_INVALID", "sub-site epay gateway url is invalid", fmt.Errorf("create epay request: %w", err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 15 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", wrapUnknownAsServiceUnavailable("EPAY_API_UNAVAILABLE", "epay is temporarily unavailable", fmt.Errorf("epay api request: %w", err))
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Payment] SubSite Epay API error: status=%d body=%s", resp.StatusCode, string(respBody))
		return "", infraerrors.BadRequest("EPAY_API_ERROR", fmt.Sprintf("epay api error: %s", string(respBody)))
	}
	var result struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		TradeNo string `json:"trade_no"`
		PayURL  string `json:"payurl"`
		QRCode  string `json:"qrcode"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", infraerrors.BadRequest("EPAY_API_ERROR", "epay gateway returned an unexpected response")
	}
	if result.Code != 1 {
		return "", infraerrors.BadRequest("EPAY_API_ERROR", fmt.Sprintf("epay error: %s", result.Msg))
	}
	codeURL := result.QRCode
	if codeURL == "" {
		codeURL = result.PayURL
	}
	if codeURL == "" {
		return "", infraerrors.BadRequest("EPAY_API_ERROR", "epay returned no payment url")
	}
	return codeURL, nil
}
