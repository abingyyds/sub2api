package service

// Status constants
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusError    = "error"
	StatusUnused   = "unused"
	StatusUsed     = "used"
	StatusExpired  = "expired"
)

// Role constants
const (
	RoleAdmin    = "admin"
	RoleSubAdmin = "sub_admin"
	RoleUser     = "user"
	RoleOrgAdmin = "org_admin"
)

// Platform constants
const (
	PlatformAnthropic   = "anthropic"
	PlatformOpenAI      = "openai"
	PlatformGemini      = "gemini"
	PlatformAntigravity = "antigravity"
	PlatformMulti       = "multi" // 多平台分组，根据请求路由自动选择账号
)

// Account type constants
const (
	AccountTypeOAuth      = "oauth"       // OAuth类型账号（full scope: profile + inference）
	AccountTypeSetupToken = "setup-token" // Setup Token类型账号（inference only scope）
	AccountTypeAPIKey     = "apikey"      // API Key类型账号
)

// Redeem type constants
const (
	RedeemTypeBalance      = "balance"
	RedeemTypeConcurrency  = "concurrency"
	RedeemTypeSubscription = "subscription"
)

// PromoCode status constants
const (
	PromoCodeStatusActive   = "active"
	PromoCodeStatusDisabled = "disabled"
)

// PromoCode discount type constants
const (
	PromoCodeDiscountTypeFixed      = "fixed"      // 固定金额减免（分）
	PromoCodeDiscountTypePercentage = "percentage" // 百分比折扣
)

// Admin adjustment type constants
const (
	AdjustmentTypeAdminBalance     = "admin_balance"     // 管理员调整余额
	AdjustmentTypeAdminConcurrency = "admin_concurrency" // 管理员调整并发数
)

// Group subscription type constants
const (
	SubscriptionTypeStandard     = "standard"     // 标准计费模式（按余额扣费）
	SubscriptionTypeSubscription = "subscription" // 订阅模式（按限额控制）
)

// Subscription status constants
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusSuspended = "suspended"
)

// Referral reward status constants
const (
	ReferralRewardPending  = "pending"
	ReferralRewardRewarded = "rewarded"
)

// Agent status constants
const (
	AgentStatusPending  = "pending"
	AgentStatusApproved = "approved"
	AgentStatusRejected = "rejected"
)

// Agent commission source type constants
const (
	AgentCommissionSourcePayment = "payment"
)

// Agent commission status constants
const (
	AgentCommissionStatusPending = "pending"
	AgentCommissionStatusSettled = "settled"
)

// LinuxDoConnectSyntheticEmailDomain 是 LinuxDo Connect 用户的合成邮箱后缀（RFC 保留域名）。
const LinuxDoConnectSyntheticEmailDomain = "@linuxdo-connect.invalid"

// Setting keys
const (
	// 注册设置
	SettingKeyRegistrationEnabled  = "registration_enabled"   // 是否开放注册
	SettingKeyEmailVerifyEnabled   = "email_verify_enabled"   // 是否开启邮件验证
	SettingKeyPromoCodeEnabled     = "promo_code_enabled"     // 是否启用优惠码功能
	SettingKeyPasswordResetEnabled = "password_reset_enabled" // 是否启用忘记密码功能（需要先开启邮件验证）

	// 邮件服务设置
	SettingKeySMTPHost     = "smtp_host"      // SMTP服务器地址
	SettingKeySMTPPort     = "smtp_port"      // SMTP端口
	SettingKeySMTPUsername = "smtp_username"  // SMTP用户名
	SettingKeySMTPPassword = "smtp_password"  // SMTP密码（加密存储）
	SettingKeySMTPFrom     = "smtp_from"      // 发件人地址
	SettingKeySMTPFromName = "smtp_from_name" // 发件人名称
	SettingKeySMTPUseTLS   = "smtp_use_tls"   // 是否使用TLS

	// Cloudflare Turnstile 设置
	SettingKeyTurnstileEnabled   = "turnstile_enabled"    // 是否启用 Turnstile 验证
	SettingKeyTurnstileSiteKey   = "turnstile_site_key"   // Turnstile Site Key
	SettingKeyTurnstileSecretKey = "turnstile_secret_key" // Turnstile Secret Key

	// TOTP 双因素认证设置
	SettingKeyTotpEnabled = "totp_enabled" // 是否启用 TOTP 2FA 功能

	// LinuxDo Connect OAuth 登录设置
	SettingKeyLinuxDoConnectEnabled      = "linuxdo_connect_enabled"
	SettingKeyLinuxDoConnectClientID     = "linuxdo_connect_client_id"
	SettingKeyLinuxDoConnectClientSecret = "linuxdo_connect_client_secret"
	SettingKeyLinuxDoConnectRedirectURL  = "linuxdo_connect_redirect_url"

	// OEM设置
	SettingKeySiteName            = "site_name"              // 网站名称
	SettingKeySiteLogo            = "site_logo"              // 网站Logo (base64)
	SettingKeySiteSubtitle        = "site_subtitle"          // 网站副标题
	SettingKeyAPIBaseURL          = "api_base_url"           // API端点地址（用于客户端配置和导入）
	SettingKeyContactInfo         = "contact_info"           // 客服联系方式
	SettingKeyDocURL              = "doc_url"                // 文档链接
	SettingKeyHomeContent         = "home_content"           // 首页内容（支持 Markdown/HTML，或 URL 作为 iframe src）
	SettingKeyHideCcsImportButton = "hide_ccs_import_button" // 是否隐藏 API Keys 页面的导入 CCS 按钮

	// 默认配置
	SettingKeyDefaultConcurrency = "default_concurrency" // 新用户默认并发量
	SettingKeyDefaultBalance     = "default_balance"     // 新用户默认余额
	SettingKeyMaxRetryRounds     = "max_retry_rounds"    // 失败重试轮数

	// 管理员 API Key
	SettingKeyAdminAPIKey = "admin_api_key" // 全局管理员 API Key（用于外部系统集成）

	// Gemini 配额策略（JSON）
	SettingKeyGeminiQuotaPolicy = "gemini_quota_policy"

	// Model fallback settings
	SettingKeyEnableModelFallback      = "enable_model_fallback"
	SettingKeyFallbackModelAnthropic   = "fallback_model_anthropic"
	SettingKeyFallbackModelOpenAI      = "fallback_model_openai"
	SettingKeyFallbackModelGemini      = "fallback_model_gemini"
	SettingKeyFallbackModelAntigravity = "fallback_model_antigravity"

	// Request identity patch (Claude -> Gemini systemInstruction injection)
	SettingKeyEnableIdentityPatch = "enable_identity_patch"
	SettingKeyIdentityPatchPrompt = "identity_patch_prompt"

	// =========================
	// Ops Monitoring (vNext)
	// =========================

	// SettingKeyOpsMonitoringEnabled is a DB-backed soft switch to enable/disable ops module at runtime.
	SettingKeyOpsMonitoringEnabled = "ops_monitoring_enabled"

	// SettingKeyOpsRealtimeMonitoringEnabled controls realtime features (e.g. WS/QPS push).
	SettingKeyOpsRealtimeMonitoringEnabled = "ops_realtime_monitoring_enabled"

	// SettingKeyOpsQueryModeDefault controls the default query mode for ops dashboard (auto/raw/preagg).
	SettingKeyOpsQueryModeDefault = "ops_query_mode_default"

	// SettingKeyOpsEmailNotificationConfig stores JSON config for ops email notifications.
	SettingKeyOpsEmailNotificationConfig = "ops_email_notification_config"

	// SettingKeyOpsAlertRuntimeSettings stores JSON config for ops alert evaluator runtime settings.
	SettingKeyOpsAlertRuntimeSettings = "ops_alert_runtime_settings"

	// SettingKeyOpsMetricsIntervalSeconds controls the ops metrics collector interval (>=60).
	SettingKeyOpsMetricsIntervalSeconds = "ops_metrics_interval_seconds"

	// SettingKeyOpsAdvancedSettings stores JSON config for ops advanced settings (data retention, aggregation).
	SettingKeyOpsAdvancedSettings = "ops_advanced_settings"

	// =========================
	// Referral / Invite Reward
	// =========================

	SettingKeyReferralEnabled      = "referral_enabled"       // 是否启用邀请返利功能
	SettingKeyReferralRewardAmount = "referral_reward_amount" // 邀请奖励金额
	SettingKeyInviteeRewardAmount  = "invitee_reward_amount"  // 被邀请人注册奖励金额

	// =========================
	// Stream Timeout Handling
	// =========================

	// SettingKeyStreamTimeoutSettings stores JSON config for stream timeout handling.
	SettingKeyStreamTimeoutSettings = "stream_timeout_settings"

	// =========================
	// Payment / WeChat Pay
	// =========================

	SettingKeyPaymentEnabled       = "payment_enabled"          // 是否启用在线支付
	SettingKeyPaymentPlans         = "payment_plans"            // 套餐配置 JSON
	SettingKeyRechargeMinAmount    = "recharge_min_amount"      // 充值最低金额（元）
	SettingKeyRechargePlans        = "recharge_plans"           // 充值优惠套餐配置 JSON
	SettingKeyWechatPayAppID       = "wechat_pay_appid"         // 微信支付关联的AppID
	SettingKeyWechatPayMchID       = "wechat_pay_mch_id"        // 微信支付商户号
	SettingKeyWechatPayAPIv3Key    = "wechat_pay_apiv3_key"     // APIv3 密钥
	SettingKeyWechatPayMchSerialNo = "wechat_pay_mch_serial_no" // 商户API证书序列号
	SettingKeyWechatPayPublicKeyID = "wechat_pay_public_key_id" // 微信支付公钥ID（验证回调）
	SettingKeyWechatPayPublicKey   = "wechat_pay_public_key"    // 微信支付公钥内容(PEM)
	SettingKeyWechatPayPrivateKey  = "wechat_pay_private_key"   // 商户API私钥(PEM)
	SettingKeyWechatPayNotifyURL   = "wechat_pay_notify_url"    // 支付回调通知URL

	// 初始余额有效期
	SettingKeyInitialBalanceExpiryDays = "initial_balance_expiry_days" // 新用户初始余额有效天数（0=永不过期）

	// =========================
	// Agent / Affiliate System
	// =========================

	SettingKeyAgentEnabled               = "agent_enabled"                 // 是否启用代理系统
	SettingKeyAgentDefaultCommissionRate = "agent_default_commission_rate" // 默认佣金比例
	SettingKeyAgentActivationFee         = "agent_activation_fee"          // 代理开通费（元）
	SettingKeyAgentContractVersion       = "agent_contract_version"        // 代理合同版本
	SettingKeyAgentContractTemplate      = "agent_contract_template"       // 代理合同 PDF 模板（data URL）
	SettingKeyAgentWithdrawFreezeDays    = "agent_withdraw_freeze_days"    // 提现冻结天数
	SettingKeyAgentWithdrawWeekday       = "agent_withdraw_weekday"        // 提现开放星期（1-7）
	SettingKeyAgentWithdrawStartHour     = "agent_withdraw_start_hour"     // 提现开始小时
	SettingKeyAgentWithdrawEndHour       = "agent_withdraw_end_hour"       // 提现结束小时

	// =========================
	// SubSite / 分站系统
	// =========================

	SettingKeySubSiteSelfServiceEnabled     = "subsite_self_service_enabled"
	SettingKeySubSiteEntryEnabled           = "subsite_entry_enabled"
	SettingKeySubSiteActivationPriceFen     = "subsite_activation_price_fen"
	SettingKeySubSiteActivationValidityDays = "subsite_activation_validity_days"
	SettingKeySubSiteDefaultThemeTemplate   = "subsite_default_theme_template"
	SettingKeySubSiteDefaultCustomConfig    = "subsite_default_custom_config"
)

// Agent profile status constants.
const (
	AgentIdentityStatusUnsubmitted = "unsubmitted"
	AgentIdentityStatusSubmitted   = "submitted"

	AgentContractStatusUnsigned = "unsigned"
	AgentContractStatusSigned   = "signed"
)

// Agent wallet balance type constants.
const (
	AgentBalanceTypeSite         = "site"
	AgentBalanceTypeFrozen       = "frozen"
	AgentBalanceTypeWithdrawable = "withdrawable"
)

// Agent wallet change type constants.
const (
	AgentWalletChangeInviteCommission   = "invite_commission"
	AgentWalletChangeConsumptionRevenue = "consumption_commission"
	AgentWalletChangeFreeze             = "freeze"
	AgentWalletChangeUnfreeze           = "unfreeze"
	AgentWalletChangeWithdrawApply      = "withdraw_apply"
	AgentWalletChangeWithdrawReject     = "withdraw_reject"
	AgentWalletChangeWithdrawPaid       = "withdraw_paid"
	AgentWalletChangeManualAdjust       = "manual_adjust"
	AgentWalletChangeSpend              = "spend"
)

// Agent withdraw status constants.
const (
	AgentWithdrawStatusPending  = "pending"
	AgentWithdrawStatusApproved = "approved"
	AgentWithdrawStatusPaid     = "paid"
	AgentWithdrawStatusRejected = "rejected"
	AgentWithdrawStatusCanceled = "canceled"
)

// Payment order type constants.
const (
	PaymentOrderTypeSubscription      = "subscription"
	PaymentOrderTypeBalance           = "balance"
	PaymentOrderTypeAgentActivation   = "agent_activation"
	PaymentOrderTypeSubSiteActivation = "subsite_activation"
	PaymentOrderTypeSubSiteTopup      = "subsite_topup"
)

// AdminAPIKeyPrefix is the prefix for admin API keys (distinct from user "sk-" keys).
const AdminAPIKeyPrefix = "admin-"

// Organization billing mode constants
const (
	OrgBillingModeBalance      = "balance"
	OrgBillingModeSubscription = "subscription"
)

// Organization status constants
const (
	OrgStatusActive    = "active"
	OrgStatusSuspended = "suspended"
	OrgStatusDisabled  = "disabled"
)

// Organization member role constants
const (
	OrgMemberRoleAdmin  = "org_admin"
	OrgMemberRoleMember = "member"
)

// Organization audit mode constants
const (
	OrgAuditModeMetadata = "metadata"
	OrgAuditModeSummary  = "summary"
	OrgAuditModeFull     = "full"
)

// Payment order status constants
const (
	PaymentOrderStatusPending  = "pending"
	PaymentOrderStatusPaid     = "paid"
	PaymentOrderStatusClosed   = "closed"
	PaymentOrderStatusRefunded = "refunded"
)

// Payment method constants
const (
	PaymentMethodWechatNative = "wechat_native"
	PaymentMethodAlipayNative = "alipay_native"
	PaymentMethodEpayAlipay   = "epay_alipay"
	PaymentMethodEpayWxpay    = "epay_wxpay"
)

// Alipay settings keys
const (
	SettingKeyAlipayEnabled      = "alipay_enabled"
	SettingKeyAlipayAppID        = "alipay_app_id"
	SettingKeyAlipayPrivateKey   = "alipay_private_key"
	SettingKeyAlipayPublicKey    = "alipay_public_key"
	SettingKeyAlipayNotifyURL    = "alipay_notify_url"
	SettingKeyAlipayIsProduction = "alipay_is_production"
)

// Epay (易支付) settings keys
const (
	SettingKeyEpayEnabled   = "epay_enabled"
	SettingKeyEpayGateway   = "epay_gateway"
	SettingKeyEpayPID       = "epay_pid"
	SettingKeyEpayPKey      = "epay_pkey"
	SettingKeyEpayNotifyURL = "epay_notify_url"
)
