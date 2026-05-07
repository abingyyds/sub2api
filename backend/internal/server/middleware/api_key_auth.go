package middleware

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// NewAPIKeyAuthMiddleware 创建 API Key 认证中间件
func NewAPIKeyAuthMiddleware(apiKeyService *service.APIKeyService, subscriptionService *service.SubscriptionService, quotaPackageRepo service.QuotaPackageRepository, orgService *service.OrganizationService, orgMemberService *service.OrgMemberService, orgProjectService *service.OrgProjectService, cfg *config.Config) APIKeyAuthMiddleware {
	return APIKeyAuthMiddleware(apiKeyAuthWithSubscription(apiKeyService, subscriptionService, quotaPackageRepo, orgService, orgMemberService, orgProjectService, cfg))
}

// apiKeyAuthWithSubscription API Key认证中间件（支持订阅验证）
func apiKeyAuthWithSubscription(apiKeyService *service.APIKeyService, subscriptionService *service.SubscriptionService, quotaPackageRepo service.QuotaPackageRepository, orgService *service.OrganizationService, orgMemberService *service.OrgMemberService, orgProjectService *service.OrgProjectService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryKey := strings.TrimSpace(c.Query("key"))
		queryApiKey := strings.TrimSpace(c.Query("api_key"))
		if queryKey != "" || queryApiKey != "" {
			AbortWithError(c, 400, "api_key_in_query_deprecated", "API key in query parameter is deprecated. Please use Authorization header instead.")
			return
		}

		// 尝试从Authorization header中提取API key (Bearer scheme)
		authHeader := c.GetHeader("Authorization")
		var apiKeyString string

		if authHeader != "" {
			// 验证Bearer scheme
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				apiKeyString = parts[1]
			}
		}

		// 如果Authorization header中没有，尝试从x-api-key header中提取
		if apiKeyString == "" {
			apiKeyString = c.GetHeader("x-api-key")
		}

		// 如果x-api-key header中没有，尝试从x-goog-api-key header中提取（Gemini CLI兼容）
		if apiKeyString == "" {
			apiKeyString = c.GetHeader("x-goog-api-key")
		}

		// 如果所有header都没有API key
		if apiKeyString == "" {
			AbortWithError(c, 401, "API_KEY_REQUIRED", "API key is required in Authorization header (Bearer scheme), x-api-key header, or x-goog-api-key header")
			return
		}

		// 从数据库验证API key
		apiKey, err := apiKeyService.GetByKey(c.Request.Context(), apiKeyString)
		if err != nil {
			if errors.Is(err, service.ErrAPIKeyNotFound) {
				AbortWithError(c, 401, "INVALID_API_KEY", "Invalid API key")
				return
			}
			AbortWithError(c, 500, "INTERNAL_ERROR", "Failed to validate API key")
			return
		}

		// 检查API key是否激活
		if !apiKey.IsActive() {
			AbortWithError(c, 401, "API_KEY_DISABLED", "API key is disabled")
			return
		}

		// 检查 IP 限制（白名单/黑名单）
		// 注意：错误信息故意模糊，避免暴露具体的 IP 限制机制
		if len(apiKey.IPWhitelist) > 0 || len(apiKey.IPBlacklist) > 0 {
			clientIP := ip.GetClientIP(c)
			allowed, _ := ip.CheckIPRestriction(clientIP, apiKey.IPWhitelist, apiKey.IPBlacklist)
			if !allowed {
				log.Printf("[Auth] 403 ACCESS_DENIED: user=%d path=%s ip=%s", apiKey.User.ID, c.Request.URL.Path, clientIP)
				AbortWithError(c, 403, "ACCESS_DENIED", "Access denied")
				return
			}
		}

		// 检查关联的用户
		if apiKey.User == nil {
			AbortWithError(c, 401, "USER_NOT_FOUND", "User associated with API key not found")
			return
		}

		// 检查用户状态
		if !apiKey.User.IsActive() {
			AbortWithError(c, 401, "USER_INACTIVE", "User account is not active")
			return
		}

		// === 企业计费前置检查 ===
		// 如果 API Key 绑定了组织，检查组织状态和成员限额
		if apiKey.OrgID != nil {
			org, err := orgService.GetByID(c.Request.Context(), *apiKey.OrgID)
			if err != nil {
				AbortWithError(c, 403, "ORG_NOT_FOUND", "Organization not found")
				return
			}
			if !org.IsActive() {
				AbortWithError(c, 403, "ORG_INACTIVE", "Organization is not active")
				return
			}

			// Check org balance (for balance billing mode)
			if org.BillingMode == service.OrgBillingModeBalance && org.Balance <= 0 {
				AbortWithError(c, 403, "ORG_INSUFFICIENT_BALANCE", "Organization has insufficient balance")
				return
			}

			// Check member quota (pre-check with 0 cost)
			member, err := orgMemberService.GetByOrgAndUser(c.Request.Context(), org.ID, apiKey.User.ID)
			if err != nil {
				AbortWithError(c, 403, "ORG_MEMBER_NOT_FOUND", "User is not a member of this organization")
				return
			}
			if !member.IsActive() {
				AbortWithError(c, 403, "ORG_MEMBER_INACTIVE", "Organization membership is not active")
				return
			}
			if err := orgMemberService.CheckMemberQuota(member, 0); err != nil {
				AbortWithError(c, 429, "ORG_MEMBER_QUOTA_EXCEEDED", err.Error())
				return
			}

			// Store org info in context for downstream billing
			c.Set(string(ContextKeyOrganization), org)
			c.Set(string(ContextKeyOrgMember), member)

			// Check project model whitelist if API key is bound to a project
			if apiKey.OrgProjectID != nil && orgProjectService != nil {
				project, err := orgProjectService.GetByID(c.Request.Context(), *apiKey.OrgProjectID)
				if err != nil {
					AbortWithError(c, 403, "ORG_PROJECT_NOT_FOUND", "Organization project not found")
					return
				}
				if !project.IsActive() {
					AbortWithError(c, 403, "ORG_PROJECT_INACTIVE", "Organization project is not active")
					return
				}
				c.Set(string(ContextKeyOrgProject), project)
			}
		}

		if cfg.RunMode == config.RunModeSimple {
			// 简易模式：跳过余额和订阅检查，但仍需设置必要的上下文
			c.Set(string(ContextKeyAPIKey), apiKey)
			c.Set(string(ContextKeyUser), AuthSubject{
				UserID:      apiKey.User.ID,
				Concurrency: apiKey.User.Concurrency,
			})
			c.Set(string(ContextKeyUserRole), apiKey.User.Role)
			setGroupContext(c, apiKey.Group)
			c.Next()
			return
		}

		// 判断计费方式：订阅模式 / 额度包 / 余额模式。
		// 额度包分组允许从既有订阅套餐平滑切换：已有活跃订阅继续生效；
		// 订阅不存在或当前窗口额度已用完时，再回落到额度包扣费。
		isSubscriptionType := apiKey.Group != nil && apiKey.Group.IsSubscriptionType()
		isQuotaPackageType := apiKey.Group != nil && apiKey.Group.IsQuotaPackage()

		subscriptionAuthorized := false
		if isSubscriptionType && subscriptionService != nil {
			subscription, err := validateSubscriptionForAPIKeyAuth(c.Request.Context(), subscriptionService, apiKey)
			if err == nil {
				c.Set(string(ContextKeySubscription), subscription)
				subscriptionAuthorized = true
			} else if !isQuotaPackageType {
				abortSubscriptionAuthError(c, apiKey, err)
				return
			} else {
				log.Printf("[Auth] quota package fallback after subscription check failed: user=%d group=%d path=%s err=%v", apiKey.User.ID, apiKey.Group.ID, c.Request.URL.Path, err)
			}
		}

		if isQuotaPackageType && !subscriptionAuthorized {
			if quotaPackageRepo == nil {
				AbortWithError(c, 503, "QUOTA_PACKAGE_UNAVAILABLE", "Quota package billing is unavailable")
				return
			}
			available, err := quotaPackageRepo.GetAvailableTotal(c.Request.Context(), apiKey.User.ID, apiKey.Group.ID)
			if err != nil {
				log.Printf("[Auth] 500 QUOTA_PACKAGE_CHECK_FAILED: user=%d group=%d path=%s err=%v", apiKey.User.ID, apiKey.Group.ID, c.Request.URL.Path, err)
				AbortWithError(c, 500, "QUOTA_PACKAGE_CHECK_FAILED", "Failed to validate quota package")
				return
			}
			if available <= 0 {
				log.Printf("[Auth] 403 QUOTA_PACKAGE_INSUFFICIENT: user=%d group=%d path=%s", apiKey.User.ID, apiKey.Group.ID, c.Request.URL.Path)
				AbortWithError(c, 403, "QUOTA_PACKAGE_INSUFFICIENT", "Quota package balance is insufficient")
				return
			}
		} else if !isSubscriptionType || subscriptionService == nil {
			// 余额模式：检查用户余额
			if apiKey.User.Balance <= 0 {
				log.Printf("[Auth] 403 INSUFFICIENT_BALANCE: user=%d balance=%.4f path=%s", apiKey.User.ID, apiKey.User.Balance, c.Request.URL.Path)
				AbortWithError(c, 403, "INSUFFICIENT_BALANCE", "Insufficient account balance")
				return
			}
		}

		// 将API key和用户信息存入上下文
		c.Set(string(ContextKeyAPIKey), apiKey)
		c.Set(string(ContextKeyUser), AuthSubject{
			UserID:      apiKey.User.ID,
			Concurrency: apiKey.User.Concurrency,
		})
		c.Set(string(ContextKeyUserRole), apiKey.User.Role)
		setGroupContext(c, apiKey.Group)

		c.Next()
	}
}

func validateSubscriptionForAPIKeyAuth(ctx context.Context, subscriptionService *service.SubscriptionService, apiKey *service.APIKey) (*service.UserSubscription, error) {
	subscription, err := subscriptionService.GetActiveSubscription(
		ctx,
		apiKey.User.ID,
		apiKey.Group.ID,
	)
	if err != nil {
		return nil, err
	}
	if err := subscriptionService.ValidateSubscription(ctx, subscription); err != nil {
		return nil, err
	}
	if err := subscriptionService.CheckAndActivateWindow(ctx, subscription); err != nil {
		log.Printf("Failed to activate subscription windows: %v", err)
	}
	if err := subscriptionService.CheckAndResetWindows(ctx, subscription); err != nil {
		log.Printf("Failed to reset subscription windows: %v", err)
	}
	if err := subscriptionService.CheckUsageLimits(ctx, subscription, apiKey.Group, 0); err != nil {
		return nil, err
	}
	return subscription, nil
}

func abortSubscriptionAuthError(c *gin.Context, apiKey *service.APIKey, err error) {
	if errors.Is(err, service.ErrSubscriptionNotFound) {
		log.Printf("[Auth] 403 SUBSCRIPTION_NOT_FOUND: user=%d group=%d path=%s", apiKey.User.ID, apiKey.Group.ID, c.Request.URL.Path)
		AbortWithError(c, 403, "SUBSCRIPTION_NOT_FOUND", "No active subscription found for this group")
		return
	}
	if errors.Is(err, service.ErrDailyLimitExceeded) ||
		errors.Is(err, service.ErrWeeklyLimitExceeded) ||
		errors.Is(err, service.ErrMonthlyLimitExceeded) {
		AbortWithError(c, 429, "USAGE_LIMIT_EXCEEDED", err.Error()+", please switch to balance mode")
		return
	}
	log.Printf("[Auth] 403 SUBSCRIPTION_INVALID: user=%d group=%d path=%s err=%v", apiKey.User.ID, apiKey.Group.ID, c.Request.URL.Path, err)
	AbortWithError(c, 403, "SUBSCRIPTION_INVALID", err.Error())
}

// GetAPIKeyFromContext 从上下文中获取API key
func GetAPIKeyFromContext(c *gin.Context) (*service.APIKey, bool) {
	value, exists := c.Get(string(ContextKeyAPIKey))
	if !exists {
		return nil, false
	}
	apiKey, ok := value.(*service.APIKey)
	return apiKey, ok
}

// GetSubscriptionFromContext 从上下文中获取订阅信息
func GetSubscriptionFromContext(c *gin.Context) (*service.UserSubscription, bool) {
	value, exists := c.Get(string(ContextKeySubscription))
	if !exists {
		return nil, false
	}
	subscription, ok := value.(*service.UserSubscription)
	return subscription, ok
}

func setGroupContext(c *gin.Context, group *service.Group) {
	if !service.IsGroupContextValid(group) {
		return
	}
	if existing, ok := c.Request.Context().Value(ctxkey.Group).(*service.Group); ok && existing != nil && existing.ID == group.ID && service.IsGroupContextValid(existing) {
		return
	}
	ctx := context.WithValue(c.Request.Context(), ctxkey.Group, group)
	c.Request = c.Request.WithContext(ctx)
}
