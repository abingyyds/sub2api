package dto

// SystemSettings represents the admin settings API response payload.
type SystemSettings struct {
	RegistrationEnabled         bool `json:"registration_enabled"`
	EmailVerifyEnabled          bool `json:"email_verify_enabled"`
	PromoCodeEnabled            bool `json:"promo_code_enabled"`
	PasswordResetEnabled        bool `json:"password_reset_enabled"`
	TotpEnabled                 bool `json:"totp_enabled"`                   // TOTP 双因素认证
	TotpEncryptionKeyConfigured bool `json:"totp_encryption_key_configured"` // TOTP 加密密钥是否已配置

	SMTPHost               string `json:"smtp_host"`
	SMTPPort               int    `json:"smtp_port"`
	SMTPUsername           string `json:"smtp_username"`
	SMTPPasswordConfigured bool   `json:"smtp_password_configured"`
	SMTPFrom               string `json:"smtp_from_email"`
	SMTPFromName           string `json:"smtp_from_name"`
	SMTPUseTLS             bool   `json:"smtp_use_tls"`

	TurnstileEnabled             bool   `json:"turnstile_enabled"`
	TurnstileSiteKey             string `json:"turnstile_site_key"`
	TurnstileSecretKeyConfigured bool   `json:"turnstile_secret_key_configured"`

	LinuxDoConnectEnabled                bool   `json:"linuxdo_connect_enabled"`
	LinuxDoConnectClientID               string `json:"linuxdo_connect_client_id"`
	LinuxDoConnectClientSecretConfigured bool   `json:"linuxdo_connect_client_secret_configured"`
	LinuxDoConnectRedirectURL            string `json:"linuxdo_connect_redirect_url"`

	SiteName            string `json:"site_name"`
	SiteLogo            string `json:"site_logo"`
	SiteSubtitle        string `json:"site_subtitle"`
	APIBaseURL          string `json:"api_base_url"`
	ContactInfo         string `json:"contact_info"`
	DocURL              string `json:"doc_url"`
	HomeContent         string `json:"home_content"`
	HideCcsImportButton bool   `json:"hide_ccs_import_button"`

	DefaultConcurrency int     `json:"default_concurrency"`
	DefaultBalance     float64 `json:"default_balance"`
	MaxRetryRounds     int     `json:"max_retry_rounds"`

	// 初始余额有效期
	InitialBalanceExpiryDays int `json:"initial_balance_expiry_days"`

	// Model fallback configuration
	EnableModelFallback      bool   `json:"enable_model_fallback"`
	FallbackModelAnthropic   string `json:"fallback_model_anthropic"`
	FallbackModelOpenAI      string `json:"fallback_model_openai"`
	FallbackModelGemini      string `json:"fallback_model_gemini"`
	FallbackModelAntigravity string `json:"fallback_model_antigravity"`

	// Identity patch configuration (Claude -> Gemini)
	EnableIdentityPatch bool   `json:"enable_identity_patch"`
	IdentityPatchPrompt string `json:"identity_patch_prompt"`

	// Ops monitoring (vNext)
	OpsMonitoringEnabled         bool   `json:"ops_monitoring_enabled"`
	OpsRealtimeMonitoringEnabled bool   `json:"ops_realtime_monitoring_enabled"`
	OpsQueryModeDefault          string `json:"ops_query_mode_default"`
	OpsMetricsIntervalSeconds    int    `json:"ops_metrics_interval_seconds"`

	// Referral / Invite Reward
	ReferralEnabled            bool    `json:"referral_enabled"`
	ReferralRewardAmount       float64 `json:"referral_reward_amount"`
	InviteeRewardAmount        float64 `json:"invitee_reward_amount"`
	AgentEnabled               bool    `json:"agent_enabled"`
	AgentDefaultCommissionRate float64 `json:"agent_default_commission_rate"`
	AgentActivationFee         float64 `json:"agent_activation_fee"`
	AgentContractVersion       string  `json:"agent_contract_version"`
	AgentContractTemplate      string  `json:"agent_contract_template"`
	AgentWithdrawFreezeDays    int     `json:"agent_withdraw_freeze_days"`
	AgentWithdrawWeekday       int     `json:"agent_withdraw_weekday"`
	AgentWithdrawStartHour     int     `json:"agent_withdraw_start_hour"`
	AgentWithdrawEndHour       int     `json:"agent_withdraw_end_hour"`

	// Payment / WeChat Pay
	PaymentEnabled                bool    `json:"payment_enabled"`
	WechatPayAppID                string  `json:"wechat_pay_appid"`
	WechatPayMchID                string  `json:"wechat_pay_mch_id"`
	WechatPayAPIv3KeyConfigured   bool    `json:"wechat_pay_apiv3_key_configured"`
	WechatPayMchSerialNo          string  `json:"wechat_pay_mch_serial_no"`
	WechatPayPublicKeyID          string  `json:"wechat_pay_public_key_id"`
	WechatPayPublicKeyConfigured  bool    `json:"wechat_pay_public_key_configured"`
	WechatPayPrivateKeyConfigured bool    `json:"wechat_pay_private_key_configured"`
	WechatPayNotifyURL            string  `json:"wechat_pay_notify_url"`
	PaymentPlans                  string  `json:"payment_plans"`
	RechargeMinAmount             float64 `json:"recharge_min_amount"`
	RechargePlans                 string  `json:"recharge_plans"`

	// Alipay
	AlipayEnabled              bool   `json:"alipay_enabled"`
	AlipayAppID                string `json:"alipay_app_id"`
	AlipayPrivateKeyConfigured bool   `json:"alipay_private_key_configured"`
	AlipayPublicKeyConfigured  bool   `json:"alipay_public_key_configured"`
	AlipayNotifyURL            string `json:"alipay_notify_url"`
	AlipayIsProduction         bool   `json:"alipay_is_production"`

	// Epay (易支付)
	EpayEnabled        bool   `json:"epay_enabled"`
	EpayGateway        string `json:"epay_gateway"`
	EpayPID            string `json:"epay_pid"`
	EpayPKeyConfigured bool   `json:"epay_pkey_configured"`
	EpayNotifyURL      string `json:"epay_notify_url"`
}

type PublicSettings struct {
	RegistrationEnabled  bool   `json:"registration_enabled"`
	EmailVerifyEnabled   bool   `json:"email_verify_enabled"`
	PromoCodeEnabled     bool   `json:"promo_code_enabled"`
	PasswordResetEnabled bool   `json:"password_reset_enabled"`
	TotpEnabled          bool   `json:"totp_enabled"` // TOTP 双因素认证
	TurnstileEnabled     bool   `json:"turnstile_enabled"`
	TurnstileSiteKey     string `json:"turnstile_site_key"`
	SiteName             string `json:"site_name"`
	SiteLogo             string `json:"site_logo"`
	SiteSubtitle         string `json:"site_subtitle"`
	APIBaseURL           string `json:"api_base_url"`
	ContactInfo          string `json:"contact_info"`
	DocURL               string `json:"doc_url"`
	HomeContent          string `json:"home_content"`
	HideCcsImportButton  bool   `json:"hide_ccs_import_button"`
	LinuxDoOAuthEnabled  bool   `json:"linuxdo_oauth_enabled"`
	ReferralEnabled      bool   `json:"referral_enabled"`
	AgentEnabled         bool   `json:"agent_enabled"`
	Version              string `json:"version"`
}

// StreamTimeoutSettings 流超时处理配置 DTO
type StreamTimeoutSettings struct {
	Enabled                bool   `json:"enabled"`
	Action                 string `json:"action"`
	TempUnschedMinutes     int    `json:"temp_unsched_minutes"`
	ThresholdCount         int    `json:"threshold_count"`
	ThresholdWindowMinutes int    `json:"threshold_window_minutes"`
}
