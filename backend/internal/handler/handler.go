package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/admin"
	"github.com/Wei-Shaw/sub2api/internal/handler/org"
)

// AdminHandlers contains all admin-related HTTP handlers
type AdminHandlers struct {
	Dashboard            *admin.DashboardHandler
	User                 *admin.UserHandler
	Group                *admin.GroupHandler
	Account              *admin.AccountHandler
	OAuth                *admin.OAuthHandler
	OpenAIOAuth          *admin.OpenAIOAuthHandler
	GeminiOAuth          *admin.GeminiOAuthHandler
	AntigravityOAuth     *admin.AntigravityOAuthHandler
	Proxy                *admin.ProxyHandler
	Redeem               *admin.RedeemHandler
	Promo                *admin.PromoHandler
	Setting              *admin.SettingHandler
	Ops                  *admin.OpsHandler
	System               *admin.SystemHandler
	Subscription         *admin.SubscriptionHandler
	Usage                *admin.UsageHandler
	UserAttribute        *admin.UserAttributeHandler
	Referral             *admin.ReferralHandler
	Announcement         *admin.AnnouncementHandler
	Organization         *admin.OrganizationHandler
	AdminInviteCode      *admin.AdminInviteCodeHandler
	DiscoverySourceStats *admin.DiscoverySourceStatsHandler
	PaymentOrder         *admin.PaymentOrderHandler
	Agent                *admin.AgentHandler
	SubSite              *admin.SubSiteHandler
}

// Handlers contains all HTTP handlers
type Handlers struct {
	Auth          *AuthHandler
	User          *UserHandler
	APIKey        *APIKeyHandler
	Usage         *UsageHandler
	Redeem        *RedeemHandler
	Subscription  *SubscriptionHandler
	Admin         *AdminHandlers
	Org           *org.OrgHandlers
	Gateway       *GatewayHandler
	OpenAIGateway *OpenAIGatewayHandler
	Setting       *SettingHandler
	Totp          *TotpHandler
	Referral      *ReferralHandler
	Announcement  *AnnouncementHandler
	ModelPlaza    *ModelPlazaHandler
	Payment       *PaymentHandler
	Agent         *AgentHandler
	SubSite       *SubSiteHandler
	SubSiteAdmin  *SubSiteAdminHandler
	Withdraw      *WithdrawHandler
	WechatNotify  *WechatNotificationHandler
}

// BuildInfo contains build-time information
type BuildInfo struct {
	Version   string
	BuildType string // "source" for manual builds, "release" for CI builds
}
