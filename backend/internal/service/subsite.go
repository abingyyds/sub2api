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
	ID                    int64      `json:"id"`
	OwnerUserID           int64      `json:"owner_user_id"`
	OwnerEmail            string     `json:"owner_email,omitempty"`
	ParentSubSiteID       *int64     `json:"parent_sub_site_id,omitempty"`
	ParentSubSiteName     string     `json:"parent_sub_site_name,omitempty"`
	Level                 int        `json:"level"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	CustomDomain          string     `json:"custom_domain,omitempty"`
	Status                string     `json:"status"`
	SiteLogo              string     `json:"site_logo,omitempty"`
	SiteFavicon           string     `json:"site_favicon,omitempty"`
	SiteSubtitle          string     `json:"site_subtitle,omitempty"`
	Announcement          string     `json:"announcement,omitempty"`
	ContactInfo           string     `json:"contact_info,omitempty"`
	DocURL                string     `json:"doc_url,omitempty"`
	HomeContent           string     `json:"home_content,omitempty"`
	ThemeTemplate         string     `json:"theme_template,omitempty"`
	RegistrationMode      string     `json:"registration_mode,omitempty"`
	EnableTopup           bool       `json:"enable_topup"`
	AllowSubSite          bool       `json:"allow_sub_site"`
	SubSitePriceFen       int        `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64    `json:"consume_rate_multiplier"`
	BalanceFen            int64      `json:"balance_fen"`
	TotalTopupFen         int64      `json:"total_topup_fen"`
	TotalConsumedFen      int64      `json:"total_consumed_fen"`
	AllowOnlineTopup      bool       `json:"allow_online_topup"`
	AllowOfflineTopup     bool       `json:"allow_offline_topup"`
	SubscriptionExpiredAt *time.Time `json:"subscription_expired_at,omitempty"`
	UserCount             int64      `json:"user_count,omitempty"`
	ChildSiteCount        int64      `json:"child_site_count,omitempty"`
	EntryURL              string     `json:"entry_url,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type CreateSubSiteInput struct {
	OwnerUserID           int64      `json:"owner_user_id"`
	ParentSubSiteID       *int64     `json:"parent_sub_site_id,omitempty"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	CustomDomain          string     `json:"custom_domain"`
	Status                string     `json:"status"`
	SiteLogo              string     `json:"site_logo"`
	SiteFavicon           string     `json:"site_favicon"`
	SiteSubtitle          string     `json:"site_subtitle"`
	Announcement          string     `json:"announcement"`
	ContactInfo           string     `json:"contact_info"`
	DocURL                string     `json:"doc_url"`
	HomeContent           string     `json:"home_content"`
	ThemeTemplate         string     `json:"theme_template"`
	RegistrationMode      string     `json:"registration_mode"`
	EnableTopup           *bool      `json:"enable_topup,omitempty"`
	AllowSubSite          *bool      `json:"allow_sub_site,omitempty"`
	SubSitePriceFen       int        `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64    `json:"consume_rate_multiplier"`
	AllowOnlineTopup      *bool      `json:"allow_online_topup,omitempty"`
	AllowOfflineTopup     *bool      `json:"allow_offline_topup,omitempty"`
	SubscriptionExpiredAt *time.Time `json:"subscription_expired_at,omitempty"`
}

type UpdateSubSiteInput struct {
	ID                    int64      `json:"id"`
	OwnerUserID           int64      `json:"owner_user_id"`
	ParentSubSiteID       *int64     `json:"parent_sub_site_id,omitempty"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	CustomDomain          string     `json:"custom_domain"`
	Status                string     `json:"status"`
	SiteLogo              string     `json:"site_logo"`
	SiteFavicon           string     `json:"site_favicon"`
	SiteSubtitle          string     `json:"site_subtitle"`
	Announcement          string     `json:"announcement"`
	ContactInfo           string     `json:"contact_info"`
	DocURL                string     `json:"doc_url"`
	HomeContent           string     `json:"home_content"`
	ThemeTemplate         string     `json:"theme_template"`
	RegistrationMode      string     `json:"registration_mode"`
	EnableTopup           *bool      `json:"enable_topup,omitempty"`
	AllowSubSite          *bool      `json:"allow_sub_site,omitempty"`
	SubSitePriceFen       int        `json:"sub_site_price_fen"`
	ConsumeRateMultiplier float64    `json:"consume_rate_multiplier"`
	AllowOnlineTopup      *bool      `json:"allow_online_topup,omitempty"`
	AllowOfflineTopup     *bool      `json:"allow_offline_topup,omitempty"`
	SubscriptionExpiredAt *time.Time `json:"subscription_expired_at,omitempty"`
}

type UpdateOwnedSubSiteInput = UpdateSubSiteInput

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
