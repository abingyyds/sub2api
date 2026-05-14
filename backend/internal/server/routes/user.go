package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由（需要认证）
func RegisterUserRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
) {
	// 支付回调通知（无需认证，微信服务器调用）
	payment := v1.Group("/payment")
	{
		payment.POST("/wechat/notify", h.Payment.WechatNotify)
		payment.POST("/alipay/notify", h.Payment.AlipayNotify)
		payment.POST("/epay/notify", h.Payment.EpayNotify)
		payment.GET("/epay/notify", h.Payment.EpayNotify)
		// 公开接口：获取套餐列表（无需认证）
		payment.GET("/plans", h.Payment.GetPlans)
		// 公开接口：获取充值信息（无需认证）
		payment.GET("/recharge-info", h.Payment.GetRechargeInfo)
		// 公开接口：获取可用支付方式（无需认证）
		payment.GET("/methods", h.Payment.GetPayMethods)
	}

	subSite := v1.Group("/subsite")
	{
		subSite.GET("/open-info", h.SubSite.GetOpenInfo)
	}

	authenticated := v1.Group("")
	authenticated.Use(middleware.NoStore())
	authenticated.Use(gin.HandlerFunc(jwtAuth))
	{
		// 用户接口
		user := authenticated.Group("/user")
		{
			user.GET("/profile", h.User.GetProfile)
			user.PUT("/password", h.User.ChangePassword)
			user.PUT("", h.User.UpdateProfile)
			user.PUT("/discovery-source", h.User.UpdateDiscoverySource)

			if h.WechatNotify != nil {
				wechat := user.Group("/wechat-official")
				{
					wechat.GET("/status", h.WechatNotify.Status)
					wechat.GET("/bind-url", h.WechatNotify.BindURL)
					wechat.POST("/unbind", h.WechatNotify.Unbind)
				}
			}

			// TOTP 双因素认证
			totp := user.Group("/totp")
			{
				totp.GET("/status", h.Totp.GetStatus)
				totp.GET("/verification-method", h.Totp.GetVerificationMethod)
				totp.POST("/send-code", h.Totp.SendVerifyCode)
				totp.POST("/setup", h.Totp.InitiateSetup)
				totp.POST("/enable", h.Totp.Enable)
				totp.POST("/disable", h.Totp.Disable)
			}
		}

		// API Key管理
		keys := authenticated.Group("/keys")
		{
			keys.GET("", h.APIKey.List)
			keys.GET("/:id", h.APIKey.GetByID)
			keys.POST("", h.APIKey.Create)
			keys.PUT("/:id", h.APIKey.Update)
			keys.DELETE("/:id", h.APIKey.Delete)
		}

		// 用户可用分组（非管理员接口）
		groups := authenticated.Group("/groups")
		{
			groups.GET("/available", h.APIKey.GetAvailableGroups)
		}

		// 模型广场
		authenticated.GET("/model-plaza", h.ModelPlaza.List)
		authenticated.GET("/model-plaza/pricing-table", h.ModelPlaza.PricingTable)

		// 使用记录
		usage := authenticated.Group("/usage")
		{
			usage.GET("", h.Usage.List)
			usage.GET("/:id", h.Usage.GetByID)
			usage.GET("/stats", h.Usage.Stats)
			// User dashboard endpoints
			usage.GET("/dashboard/stats", h.Usage.DashboardStats)
			usage.GET("/dashboard/trend", h.Usage.DashboardTrend)
			usage.GET("/dashboard/models", h.Usage.DashboardModels)
			usage.POST("/dashboard/api-keys-usage", h.Usage.DashboardAPIKeysUsage)
		}

		// 卡密兑换
		redeem := authenticated.Group("/redeem")
		{
			redeem.POST("", h.Redeem.Redeem)
			redeem.GET("/history", h.Redeem.GetHistory)
		}

		// 用户订阅
		subscriptions := authenticated.Group("/subscriptions")
		{
			subscriptions.GET("", h.Subscription.List)
			subscriptions.GET("/active", h.Subscription.GetActive)
			subscriptions.GET("/progress", h.Subscription.GetProgress)
			subscriptions.GET("/summary", h.Subscription.GetSummary)
		}

		// 邀请返利
		referral := authenticated.Group("/referral")
		{
			referral.GET("/code", h.Referral.GetInviteCode)
			referral.GET("/invitees", h.Referral.ListInvitees)
			referral.GET("/stats", h.Referral.GetStats)
		}

		// 在线支付（需要认证）
		authPayment := authenticated.Group("/payment")
		{
			authPayment.POST("/orders", h.Payment.CreateOrder)
			authPayment.POST("/recharge", h.Payment.CreateRecharge)
			authPayment.POST("/agent-activation", h.Payment.CreateAgentActivationOrder)
			authPayment.POST("/subsite-activation", h.Payment.CreateSubSiteActivationOrder)
			authPayment.POST("/subsite-topup", h.Payment.CreateSubSiteTopupOrder)
			authPayment.GET("/orders", h.Payment.ListOrders)
			authPayment.GET("/orders/:orderNo", h.Payment.QueryOrder)
			authPayment.POST("/invoice-requests", h.Payment.SubmitInvoice)
			authPayment.GET("/newcomer-status", h.Payment.GetNewcomerStatus)
		}

		// 代理中心
		agent := authenticated.Group("/agent")
		{
			agent.GET("/status", h.Agent.GetStatus)
			agent.POST("/profile", h.Agent.SaveProfile)
			agent.POST("/apply", h.Agent.Apply)
			agent.GET("/dashboard", h.Agent.Dashboard)
			agent.GET("/link", h.Agent.GetLink)
			agent.GET("/sub-users", h.Agent.ListSubUsers)
			agent.PUT("/sub-users/:id/rate", h.Agent.SetSubUserRate)
			agent.GET("/financial-logs", h.Agent.ListFinancialLogs)
			agent.GET("/commissions", h.Agent.ListCommissions)
		}

		// 提现
		withdraw := authenticated.Group("/withdraw")
		{
			withdraw.POST("/apply", h.Withdraw.Apply)
			withdraw.GET("/list", h.Withdraw.List)
			withdraw.POST("/:id/cancel", h.Withdraw.Cancel)
		}

		ownedSubSite := authenticated.Group("/subsite")
		{
			ownedSubSite.GET("/owned", h.SubSite.ListOwned)
			ownedSubSite.GET("/owned/:id", h.SubSite.GetOwned)
			ownedSubSite.PUT("/owned/:id", h.SubSite.UpdateOwned)
			ownedSubSite.POST("/owned/:id/offline-topup", h.SubSite.OfflineTopupOwnedUser)
			ownedSubSite.GET("/owned/:id/ledger", h.SubSite.OwnedLedger)
		}
	}
}
