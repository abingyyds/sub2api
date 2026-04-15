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
	MaxSubSiteLevel                  = 2
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

type SubSiteGroupPriceOverride struct {
	GroupID      int64  `json:"group_id"`
	GroupName    string `json:"group_name,omitempty"`
	PriceFen     int    `json:"price_fen"`
	BasePriceFen int    `json:"base_price_fen,omitempty"`
}

type SubSiteRechargePriceOverride struct {
	PlanKey          string  `json:"plan_key"`
	Name             string  `json:"name,omitempty"`
	PayAmountFen     int     `json:"pay_amount_fen"`
	BasePayAmountFen int     `json:"base_pay_amount_fen,omitempty"`
	BalanceAmount    float64 `json:"balance_amount,omitempty"`
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
	DefaultCustomConfig  string                       `json:"default_custom_config,omitempty"`
	ThemeTemplates       []SubSiteThemeTemplateOption `json:"theme_templates"`
}

type PlatformSubSiteConfig struct {
	Enabled              bool                         `json:"enabled"`
	ActivationPriceFen   int                          `json:"activation_price_fen"`
	ValidityDays         int                          `json:"validity_days"`
	DefaultThemeTemplate string                       `json:"default_theme_template"`
	DefaultCustomConfig  string                       `json:"default_custom_config,omitempty"`
	ThemeTemplates       []SubSiteThemeTemplateOption `json:"theme_templates"`
}

type UpdatePlatformSubSiteConfigInput struct {
	Enabled              bool   `json:"enabled"`
	ActivationPriceFen   int    `json:"activation_price_fen"`
	ValidityDays         int    `json:"validity_days"`
	DefaultThemeTemplate string `json:"default_theme_template"`
	DefaultCustomConfig  string `json:"default_custom_config"`
}

type SubSite struct {
	ID                     int64                          `json:"id"`
	OwnerUserID            int64                          `json:"owner_user_id"`
	OwnerEmail             string                         `json:"owner_email,omitempty"`
	ParentSubSiteID        *int64                         `json:"parent_sub_site_id,omitempty"`
	ParentSubSiteName      string                         `json:"parent_sub_site_name,omitempty"`
	Level                  int                            `json:"level"`
	Name                   string                         `json:"name"`
	Slug                   string                         `json:"slug"`
	CustomDomain           string                         `json:"custom_domain,omitempty"`
	Status                 string                         `json:"status"`
	SiteLogo               string                         `json:"site_logo,omitempty"`
	SiteFavicon            string                         `json:"site_favicon,omitempty"`
	SiteSubtitle           string                         `json:"site_subtitle,omitempty"`
	Announcement           string                         `json:"announcement,omitempty"`
	ContactInfo            string                         `json:"contact_info,omitempty"`
	DocURL                 string                         `json:"doc_url,omitempty"`
	HomeContent            string                         `json:"home_content,omitempty"`
	ThemeTemplate          string                         `json:"theme_template,omitempty"`
	ThemeConfig            string                         `json:"theme_config,omitempty"`
	CustomConfig           string                         `json:"custom_config,omitempty"`
	RegistrationMode       string                         `json:"registration_mode,omitempty"`
	EnableTopup            bool                           `json:"enable_topup"`
	AllowSubSite           bool                           `json:"allow_sub_site"`
	SubSitePriceFen        int                            `json:"sub_site_price_fen"`
	SubscriptionExpiredAt  *time.Time                     `json:"subscription_expired_at,omitempty"`
	UserCount              int64                          `json:"user_count,omitempty"`
	ChildSiteCount         int64                          `json:"child_site_count,omitempty"`
	EntryURL               string                         `json:"entry_url,omitempty"`
	GroupPriceOverrides    []SubSiteGroupPriceOverride    `json:"group_price_overrides,omitempty"`
	RechargePriceOverrides []SubSiteRechargePriceOverride `json:"recharge_price_overrides,omitempty"`
	CreatedAt              time.Time                      `json:"created_at"`
	UpdatedAt              time.Time                      `json:"updated_at"`
}

type CreateSubSiteInput struct {
	OwnerUserID            int64                          `json:"owner_user_id"`
	ParentSubSiteID        *int64                         `json:"parent_sub_site_id,omitempty"`
	Name                   string                         `json:"name"`
	Slug                   string                         `json:"slug"`
	CustomDomain           string                         `json:"custom_domain"`
	Status                 string                         `json:"status"`
	SiteLogo               string                         `json:"site_logo"`
	SiteFavicon            string                         `json:"site_favicon"`
	SiteSubtitle           string                         `json:"site_subtitle"`
	Announcement           string                         `json:"announcement"`
	ContactInfo            string                         `json:"contact_info"`
	DocURL                 string                         `json:"doc_url"`
	HomeContent            string                         `json:"home_content"`
	ThemeTemplate          string                         `json:"theme_template"`
	ThemeConfig            string                         `json:"theme_config"`
	CustomConfig           string                         `json:"custom_config"`
	RegistrationMode       string                         `json:"registration_mode"`
	EnableTopup            *bool                          `json:"enable_topup,omitempty"`
	AllowSubSite           *bool                          `json:"allow_sub_site,omitempty"`
	SubSitePriceFen        int                            `json:"sub_site_price_fen"`
	SubscriptionExpiredAt  *time.Time                     `json:"subscription_expired_at,omitempty"`
	GroupPriceOverrides    []SubSiteGroupPriceOverride    `json:"group_price_overrides,omitempty"`
	RechargePriceOverrides []SubSiteRechargePriceOverride `json:"recharge_price_overrides,omitempty"`
}

type UpdateSubSiteInput struct {
	ID                     int64                          `json:"id"`
	OwnerUserID            int64                          `json:"owner_user_id"`
	ParentSubSiteID        *int64                         `json:"parent_sub_site_id,omitempty"`
	Name                   string                         `json:"name"`
	Slug                   string                         `json:"slug"`
	CustomDomain           string                         `json:"custom_domain"`
	Status                 string                         `json:"status"`
	SiteLogo               string                         `json:"site_logo"`
	SiteFavicon            string                         `json:"site_favicon"`
	SiteSubtitle           string                         `json:"site_subtitle"`
	Announcement           string                         `json:"announcement"`
	ContactInfo            string                         `json:"contact_info"`
	DocURL                 string                         `json:"doc_url"`
	HomeContent            string                         `json:"home_content"`
	ThemeTemplate          string                         `json:"theme_template"`
	ThemeConfig            string                         `json:"theme_config"`
	CustomConfig           string                         `json:"custom_config"`
	RegistrationMode       string                         `json:"registration_mode"`
	EnableTopup            *bool                          `json:"enable_topup,omitempty"`
	AllowSubSite           *bool                          `json:"allow_sub_site,omitempty"`
	SubSitePriceFen        int                            `json:"sub_site_price_fen"`
	SubscriptionExpiredAt  *time.Time                     `json:"subscription_expired_at,omitempty"`
	GroupPriceOverrides    []SubSiteGroupPriceOverride    `json:"group_price_overrides,omitempty"`
	RechargePriceOverrides []SubSiteRechargePriceOverride `json:"recharge_price_overrides,omitempty"`
}

type UpdateOwnedSubSiteInput = UpdateSubSiteInput

type CreateSubSiteActivationInput struct {
	Name                   string                         `json:"name"`
	Slug                   string                         `json:"slug"`
	CustomDomain           string                         `json:"custom_domain"`
	SiteLogo               string                         `json:"site_logo"`
	SiteFavicon            string                         `json:"site_favicon"`
	SiteSubtitle           string                         `json:"site_subtitle"`
	Announcement           string                         `json:"announcement"`
	ContactInfo            string                         `json:"contact_info"`
	DocURL                 string                         `json:"doc_url"`
	HomeContent            string                         `json:"home_content"`
	ThemeTemplate          string                         `json:"theme_template"`
	ThemeConfig            string                         `json:"theme_config"`
	CustomConfig           string                         `json:"custom_config"`
	RegistrationMode       string                         `json:"registration_mode"`
	EnableTopup            *bool                          `json:"enable_topup,omitempty"`
	AllowSubSite           *bool                          `json:"allow_sub_site,omitempty"`
	SubSitePriceFen        int                            `json:"sub_site_price_fen"`
	GroupPriceOverrides    []SubSiteGroupPriceOverride    `json:"group_price_overrides,omitempty"`
	RechargePriceOverrides []SubSiteRechargePriceOverride `json:"recharge_price_overrides,omitempty"`
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
