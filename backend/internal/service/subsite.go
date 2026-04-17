package service

import "time"

const (
	SubSiteStatusPending  = "pending"
	SubSiteStatusActive   = "active"
	SubSiteStatusDisabled = "disabled"
)

const (
	SubSiteRegistrationOpen   = "open"
	SubSiteRegistrationInvite = "invite"
	SubSiteRegistrationClosed = "closed"
)

const (
	SubSiteThemeTemplateStarter  = "starter"
	SubSiteThemeTemplateAurora   = "aurora"
	SubSiteThemeTemplateSummit   = "summit"
	SubSiteThemeTemplateTerminal = "terminal"
)

const (
	DefaultSubSiteActivationPriceFen = 38800
	DefaultSubSiteValidityDays       = 365
	DefaultSubSiteMaxLevel           = 2
	MaxSubSiteLevelHardLimit         = 10
	DefaultSubSiteConsumeRate        = 1.0
	MaxSubSiteConsumeRateMultiplier  = 10.0
)

// 分站模式：
//   pool — 经典资金池模式：分站主先向主站充值到 balance_fen，线下卖用户余额赚差价；
//          支持配置自有收款账号，用户线上充值直接到分站主账号并从池同步扣减（自动进货）。
//   rate — 倍率分成模式：无独立资金池；用户充值/消费都经主站，消费时按链上复合倍率
//          把 (compound_i - compound_{i-1}) × base 的利润分级入账到 balance_fen，
//          分站主通过提现把利润取走。
const (
	SubSiteModePool = "pool"
	SubSiteModeRate = "rate"
)

type SubSiteThemeTemplateOption struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

var DefaultSubSiteThemeTemplates = []SubSiteThemeTemplateOption{
	{
		Key:         SubSiteThemeTemplateStarter,
		Label:       "Starter",
		Description: "通用 SaaS 风格，适合大多数分站直接开箱使用",
	},
	{
		Key:         SubSiteThemeTemplateAurora,
		Label:       "Aurora",
		Description: "更强展示感的渐变品牌模板，适合营销导向站点",
	},
	{
		Key:         SubSiteThemeTemplateSummit,
		Label:       "Summit",
		Description: "偏企业感的信息型模板，适合强调稳定与专业服务",
	},
	{
		Key:         SubSiteThemeTemplateTerminal,
		Label:       "Terminal",
		Description: "偏开发者风格模板，适合 Claude Code / Codex 场景",
	},
}

type SubSiteOpenInfo struct {
	Enabled              bool                         `json:"enabled"`
	Scope                string                       `json:"scope"`
	ParentSubSiteID      *int64                       `json:"parent_sub_site_id,omitempty"`
	ParentSubSiteName    string                       `json:"parent_sub_site_name,omitempty"`
	Level                int                          `json:"level"`
	MaxLevel             int                          `json:"max_level"`
	PriceFen             int                          `json:"price_fen"`
	ValidityDays         int                          `json:"validity_days"`
	Currency             string                       `json:"currency"`
	AllowCustomDomain    bool                         `json:"allow_custom_domain"`
	DefaultThemeTemplate string                       `json:"default_theme_template"`
	ThemeTemplates       []SubSiteThemeTemplateOption `json:"theme_templates"`
}

type PlatformSubSiteConfig struct {
	EntryEnabled         bool                         `json:"entry_enabled"`
	Enabled              bool                         `json:"enabled"`
	ActivationPriceFen   int                          `json:"activation_price_fen"`
	ValidityDays         int                          `json:"validity_days"`
	MaxLevel             int                          `json:"max_level"`
	DefaultThemeTemplate string                       `json:"default_theme_template"`
	ThemeTemplates       []SubSiteThemeTemplateOption `json:"theme_templates"`
}

type UpdatePlatformSubSiteConfigInput struct {
	EntryEnabled         bool   `json:"entry_enabled"`
	Enabled              bool   `json:"enabled"`
	ActivationPriceFen   int    `json:"activation_price_fen"`
	ValidityDays         int    `json:"validity_days"`
	MaxLevel             int    `json:"max_level"`
	DefaultThemeTemplate string `json:"default_theme_template"`
}

type SubSite struct {
	ID                    int64               `json:"id"`
	OwnerUserID           int64               `json:"owner_user_id"`
	OwnerEmail            string              `json:"owner_email,omitempty"`
	ParentSubSiteID       *int64              `json:"parent_sub_site_id,omitempty"`
	ParentSubSiteName     string              `json:"parent_sub_site_name,omitempty"`
	Level                 int                 `json:"level"`
	Name                  string              `json:"name"`
	Slug                  string              `json:"slug"`
	CustomDomain          string              `json:"custom_domain,omitempty"`
	Status                string              `json:"status"`
	Mode                  string              `json:"mode"`
	SiteLogo              string              `json:"site_logo,omitempty"`
	SiteFavicon           string              `json:"site_favicon,omitempty"`
	SiteSubtitle          string              `json:"site_subtitle,omitempty"`
	Announcement          string              `json:"announcement,omitempty"`
	ContactInfo           string              `json:"contact_info,omitempty"`
	DocURL                string              `json:"doc_url,omitempty"`
	HomeContent           string              `json:"home_content,omitempty"`
	ThemeTemplate         string              `json:"theme_template,omitempty"`
	RegistrationMode      string              `json:"registration_mode,omitempty"`
	EnableTopup           bool                `json:"enable_topup"`
	AllowSubSite          bool                `json:"allow_sub_site"`
	SubSitePriceFen       int                 `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64             `json:"consume_rate_multiplier"`
	BalanceFen            int64               `json:"balance_fen"`
	TotalTopupFen         int64               `json:"total_topup_fen"`
	TotalConsumedFen      int64               `json:"total_consumed_fen"`
	TotalWithdrawnFen     int64               `json:"total_withdrawn_fen"`
	AllowOnlineTopup      bool                `json:"allow_online_topup"`
	AllowOfflineTopup     bool                `json:"allow_offline_topup"`
	OwnerPaymentConfig    *OwnerPaymentConfig `json:"owner_payment_config,omitempty"`
	SubscriptionExpiredAt *time.Time          `json:"subscription_expired_at,omitempty"`
	UserCount             int64               `json:"user_count,omitempty"`
	ChildSiteCount        int64               `json:"child_site_count,omitempty"`
	EntryURL              string              `json:"entry_url,omitempty"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
}

// OwnerPaymentConfig 分站主自有收款凭据（仅 pool 模式可用）。
// 任一渠道启用时对应 Enabled 字段为 true；未配置（nil 或 Enabled=false）则退回全局 settings。
type OwnerPaymentConfig struct {
	Wechat *WechatPayCredentials `json:"wechat,omitempty"`
	Alipay *AlipayCredentials    `json:"alipay,omitempty"`
	Epay   *EpayCredentials      `json:"epay,omitempty"`
}

type WechatPayCredentials struct {
	Enabled      bool   `json:"enabled"`
	AppID        string `json:"app_id"`
	MchID        string `json:"mch_id"`
	APIv3Key     string `json:"apiv3_key"`
	MchSerialNo  string `json:"mch_serial_no"`
	PublicKeyID  string `json:"public_key_id"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	NotifyURL    string `json:"notify_url"`
}

type AlipayCredentials struct {
	Enabled      bool   `json:"enabled"`
	AppID        string `json:"app_id"`
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`
	NotifyURL    string `json:"notify_url"`
	IsProduction bool   `json:"is_production"`
}

type EpayCredentials struct {
	Enabled   bool   `json:"enabled"`
	Gateway   string `json:"gateway"`
	PID       string `json:"pid"`
	PKey      string `json:"pkey"`
	NotifyURL string `json:"notify_url"`
}

type CreateSubSiteInput struct {
	OwnerUserID           int64               `json:"owner_user_id"`
	ParentSubSiteID       *int64              `json:"parent_sub_site_id,omitempty"`
	Name                  string              `json:"name"`
	Slug                  string              `json:"slug"`
	CustomDomain          string              `json:"custom_domain"`
	Status                string              `json:"status"`
	Mode                  string              `json:"mode"`
	SiteLogo              string              `json:"site_logo"`
	SiteFavicon           string              `json:"site_favicon"`
	SiteSubtitle          string              `json:"site_subtitle"`
	Announcement          string              `json:"announcement"`
	ContactInfo           string              `json:"contact_info"`
	DocURL                string              `json:"doc_url"`
	HomeContent           string              `json:"home_content"`
	ThemeTemplate         string              `json:"theme_template"`
	RegistrationMode      string              `json:"registration_mode"`
	EnableTopup           *bool               `json:"enable_topup,omitempty"`
	AllowSubSite          *bool               `json:"allow_sub_site,omitempty"`
	SubSitePriceFen       int                 `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64             `json:"consume_rate_multiplier"`
	AllowOnlineTopup      *bool               `json:"allow_online_topup,omitempty"`
	AllowOfflineTopup     *bool               `json:"allow_offline_topup,omitempty"`
	OwnerPaymentConfig    *OwnerPaymentConfig `json:"owner_payment_config,omitempty"`
	SubscriptionExpiredAt *time.Time          `json:"subscription_expired_at,omitempty"`
}

// UpdateSubSiteInput 平台管理员可修改的完整字段集（含 Mode / SubscriptionExpiredAt 等敏感字段）。
// 分站主自助路径请使用 UpdateOwnedSubSiteInput，字段范围更窄。
type UpdateSubSiteInput struct {
	ID                    int64               `json:"id"`
	OwnerUserID           int64               `json:"owner_user_id"`
	ParentSubSiteID       *int64              `json:"parent_sub_site_id,omitempty"`
	Name                  string              `json:"name"`
	Slug                  string              `json:"slug"`
	CustomDomain          string              `json:"custom_domain"`
	Status                string              `json:"status"`
	Mode                  string              `json:"mode"`
	SiteLogo              string              `json:"site_logo"`
	SiteFavicon           string              `json:"site_favicon"`
	SiteSubtitle          string              `json:"site_subtitle"`
	Announcement          string              `json:"announcement"`
	ContactInfo           string              `json:"contact_info"`
	DocURL                string              `json:"doc_url"`
	HomeContent           string              `json:"home_content"`
	ThemeTemplate         string              `json:"theme_template"`
	RegistrationMode      string              `json:"registration_mode"`
	EnableTopup           *bool               `json:"enable_topup,omitempty"`
	AllowSubSite          *bool               `json:"allow_sub_site,omitempty"`
	SubSitePriceFen       int                 `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64             `json:"consume_rate_multiplier"`
	AllowOnlineTopup      *bool               `json:"allow_online_topup,omitempty"`
	AllowOfflineTopup     *bool               `json:"allow_offline_topup,omitempty"`
	OwnerPaymentConfig    *OwnerPaymentConfig `json:"owner_payment_config,omitempty"`
	SubscriptionExpiredAt *time.Time          `json:"subscription_expired_at,omitempty"`
}

// UpdateOwnedSubSiteInput 分站主自助可改字段；敏感字段（Mode/Status/ParentSubSiteID/
// SubscriptionExpiredAt/AllowSubSite/AllowOnlineTopup/EnableTopup/Level）由 admin 接口管理。
type UpdateOwnedSubSiteInput struct {
	ID                    int64               `json:"id"`
	Name                  string              `json:"name"`
	Slug                  string              `json:"slug"`
	CustomDomain          string              `json:"custom_domain"`
	SiteLogo              string              `json:"site_logo"`
	SiteFavicon           string              `json:"site_favicon"`
	SiteSubtitle          string              `json:"site_subtitle"`
	Announcement          string              `json:"announcement"`
	ContactInfo           string              `json:"contact_info"`
	DocURL                string              `json:"doc_url"`
	HomeContent           string              `json:"home_content"`
	ThemeTemplate         string              `json:"theme_template"`
	RegistrationMode      string              `json:"registration_mode"`
	SubSitePriceFen       int                 `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64             `json:"consume_rate_multiplier"`
	AllowOfflineTopup     *bool               `json:"allow_offline_topup,omitempty"`
	OwnerPaymentConfig    *OwnerPaymentConfig `json:"owner_payment_config,omitempty"`
}

type CreateSubSiteActivationInput struct {
	Name                  string  `json:"name"`
	Slug                  string  `json:"slug"`
	CustomDomain          string  `json:"custom_domain"`
	SiteLogo              string  `json:"site_logo"`
	SiteFavicon           string  `json:"site_favicon"`
	SiteSubtitle          string  `json:"site_subtitle"`
	Announcement          string  `json:"announcement"`
	ContactInfo           string  `json:"contact_info"`
	DocURL                string  `json:"doc_url"`
	HomeContent           string  `json:"home_content"`
	ThemeTemplate         string  `json:"theme_template"`
	RegistrationMode      string  `json:"registration_mode"`
	EnableTopup           *bool   `json:"enable_topup,omitempty"`
	AllowSubSite          *bool   `json:"allow_sub_site,omitempty"`
	SubSitePriceFen       int     `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64 `json:"consume_rate_multiplier"`
}

type SubSiteActivationRequest struct {
	PaymentOrderID     int64                        `json:"payment_order_id"`
	UserID             int64                        `json:"user_id"`
	ParentSubSiteID    *int64                       `json:"parent_sub_site_id,omitempty"`
	Level              int                          `json:"level"`
	ValidityDays       int                          `json:"validity_days"`
	Site               CreateSubSiteActivationInput `json:"site"`
	ActivatedSubSiteID *int64                       `json:"activated_sub_site_id,omitempty"`
	ActivatedAt        *time.Time                   `json:"activated_at,omitempty"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
}
