package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/admin"
	"github.com/Wei-Shaw/sub2api/internal/handler/org"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/google/wire"
)

// ProvideAdminHandlers creates the AdminHandlers struct
func ProvideAdminHandlers(
	dashboardHandler *admin.DashboardHandler,
	userHandler *admin.UserHandler,
	groupHandler *admin.GroupHandler,
	accountHandler *admin.AccountHandler,
	oauthHandler *admin.OAuthHandler,
	openaiOAuthHandler *admin.OpenAIOAuthHandler,
	geminiOAuthHandler *admin.GeminiOAuthHandler,
	antigravityOAuthHandler *admin.AntigravityOAuthHandler,
	proxyHandler *admin.ProxyHandler,
	redeemHandler *admin.RedeemHandler,
	promoHandler *admin.PromoHandler,
	settingHandler *admin.SettingHandler,
	opsHandler *admin.OpsHandler,
	systemHandler *admin.SystemHandler,
	subscriptionHandler *admin.SubscriptionHandler,
	usageHandler *admin.UsageHandler,
	userAttributeHandler *admin.UserAttributeHandler,
	referralHandler *admin.ReferralHandler,
	announcementHandler *admin.AnnouncementHandler,
	organizationHandler *admin.OrganizationHandler,
	adminInviteCodeHandler *admin.AdminInviteCodeHandler,
	discoverySourceStatsHandler *admin.DiscoverySourceStatsHandler,
	paymentOrderHandler *admin.PaymentOrderHandler,
	agentHandler *admin.AgentHandler,
	subSiteHandler *admin.SubSiteHandler,
) *AdminHandlers {
	return &AdminHandlers{
		Dashboard:            dashboardHandler,
		User:                 userHandler,
		Group:                groupHandler,
		Account:              accountHandler,
		OAuth:                oauthHandler,
		OpenAIOAuth:          openaiOAuthHandler,
		GeminiOAuth:          geminiOAuthHandler,
		AntigravityOAuth:     antigravityOAuthHandler,
		Proxy:                proxyHandler,
		Redeem:               redeemHandler,
		Promo:                promoHandler,
		Setting:              settingHandler,
		Ops:                  opsHandler,
		System:               systemHandler,
		Subscription:         subscriptionHandler,
		Usage:                usageHandler,
		UserAttribute:        userAttributeHandler,
		Referral:             referralHandler,
		Announcement:         announcementHandler,
		Organization:         organizationHandler,
		AdminInviteCode:      adminInviteCodeHandler,
		DiscoverySourceStats: discoverySourceStatsHandler,
		PaymentOrder:         paymentOrderHandler,
		Agent:                agentHandler,
		SubSite:              subSiteHandler,
	}
}

// ProvideSystemHandler creates admin.SystemHandler with version from BuildInfo
func ProvideSystemHandler(buildInfo BuildInfo) *admin.SystemHandler {
	return admin.NewSystemHandler(buildInfo.Version)
}

// ProvideSettingHandler creates SettingHandler with version from BuildInfo
func ProvideSettingHandler(settingService *service.SettingService, buildInfo BuildInfo) *SettingHandler {
	return NewSettingHandler(settingService, buildInfo.Version)
}

// ProvideOrgHandlers creates the OrgHandlers struct
func ProvideOrgHandlers(
	dashboardHandler *org.DashboardHandler,
	memberHandler *org.MemberHandler,
	projectHandler *org.ProjectHandler,
	auditLogHandler *org.AuditLogHandler,
) *org.OrgHandlers {
	return &org.OrgHandlers{
		Dashboard: dashboardHandler,
		Member:    memberHandler,
		Project:   projectHandler,
		AuditLog:  auditLogHandler,
	}
}

// ProvideHandlers creates the Handlers struct
func ProvideHandlers(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	apiKeyHandler *APIKeyHandler,
	usageHandler *UsageHandler,
	redeemHandler *RedeemHandler,
	subscriptionHandler *SubscriptionHandler,
	adminHandlers *AdminHandlers,
	orgHandlers *org.OrgHandlers,
	gatewayHandler *GatewayHandler,
	openaiGatewayHandler *OpenAIGatewayHandler,
	settingHandler *SettingHandler,
	totpHandler *TotpHandler,
	modelPlazaHandler *ModelPlazaHandler,
	referralHandler *ReferralHandler,
	announcementHandler *AnnouncementHandler,
	paymentHandler *PaymentHandler,
	agentHandler *AgentHandler,
	subSiteHandler *SubSiteHandler,
	subSiteAdminHandler *SubSiteAdminHandler,
	withdrawHandler *WithdrawHandler,
	wechatNotificationHandler *WechatNotificationHandler,
) *Handlers {
	return &Handlers{
		Auth:          authHandler,
		User:          userHandler,
		APIKey:        apiKeyHandler,
		Usage:         usageHandler,
		Redeem:        redeemHandler,
		Subscription:  subscriptionHandler,
		Admin:         adminHandlers,
		Org:           orgHandlers,
		Gateway:       gatewayHandler,
		OpenAIGateway: openaiGatewayHandler,
		Setting:       settingHandler,
		Totp:          totpHandler,
		ModelPlaza:    modelPlazaHandler,
		Referral:      referralHandler,
		Announcement:  announcementHandler,
		Payment:       paymentHandler,
		Agent:         agentHandler,
		SubSite:       subSiteHandler,
		SubSiteAdmin:  subSiteAdminHandler,
		Withdraw:      withdrawHandler,
		WechatNotify:  wechatNotificationHandler,
	}
}

// ProviderSet is the Wire provider set for all handlers
var ProviderSet = wire.NewSet(
	// Top-level handlers
	NewAuthHandler,
	NewUserHandler,
	NewAPIKeyHandler,
	NewUsageHandler,
	NewRedeemHandler,
	NewSubscriptionHandler,
	NewGatewayHandler,
	NewOpenAIGatewayHandler,
	NewTotpHandler,
	NewReferralHandler,
	NewAnnouncementHandler,
	NewModelPlazaHandler,
	NewPaymentHandler,
	NewAgentHandler,
	NewSubSiteHandler,
	NewSubSiteAdminHandler,
	NewWithdrawHandler,
	NewWechatNotificationHandler,
	ProvideSettingHandler,

	// Admin handlers
	admin.NewDashboardHandler,
	admin.NewUserHandler,
	admin.NewGroupHandler,
	admin.NewAccountHandler,
	admin.NewOAuthHandler,
	admin.NewOpenAIOAuthHandler,
	admin.NewGeminiOAuthHandler,
	admin.NewAntigravityOAuthHandler,
	admin.NewProxyHandler,
	admin.NewRedeemHandler,
	admin.NewPromoHandler,
	admin.NewSettingHandler,
	admin.NewOpsHandler,
	ProvideSystemHandler,
	admin.NewSubscriptionHandler,
	admin.NewUsageHandler,
	admin.NewUserAttributeHandler,
	admin.NewReferralHandler,
	admin.NewAnnouncementHandler,
	admin.NewOrganizationHandler,
	admin.NewAdminInviteCodeHandler,
	admin.NewDiscoverySourceStatsHandler,
	admin.NewPaymentOrderHandler,
	admin.NewAgentHandler,
	admin.NewSubSiteHandler,

	// Org handlers
	org.NewDashboardHandler,
	org.NewMemberHandler,
	org.NewProjectHandler,
	org.NewAuditLogHandler,

	// AdminHandlers, OrgHandlers and Handlers constructors
	ProvideAdminHandlers,
	ProvideOrgHandlers,
	ProvideHandlers,
)
