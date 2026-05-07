package server

import (
	"log"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/server/routes"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupRouter 配置路由器中间件和路由
func SetupRouter(
	r *gin.Engine,
	handlers *handler.Handlers,
	jwtAuth middleware2.JWTAuthMiddleware,
	adminAuth middleware2.AdminAuthMiddleware,
	apiKeyAuth middleware2.APIKeyAuthMiddleware,
	orgAuth middleware2.OrgAuthMiddleware,
	apiKeyService *service.APIKeyService,
	subscriptionService *service.SubscriptionService,
	quotaPackageRepo service.QuotaPackageRepository,
	opsService *service.OpsService,
	settingService *service.SettingService,
	subSiteService *service.SubSiteService,
	cfg *config.Config,
	redisClient *redis.Client,
) *gin.Engine {
	// 应用中间件
	r.Use(middleware2.Logger())
	r.Use(middleware2.CORS(cfg.CORS))
	r.Use(middleware2.SecurityHeaders(cfg.Security.CSP))
	r.Use(middleware2.SubSiteIdentify(subSiteService))

	// Serve embedded frontend with settings injection if available
	if web.HasEmbeddedFrontend() {
		frontendServer, err := web.NewFrontendServer(settingService)
		if err != nil {
			log.Printf("Warning: Failed to create frontend server with settings injection: %v, using legacy mode", err)
			r.Use(web.ServeEmbeddedFrontend())
		} else {
			// Register cache invalidation callback
			settingService.SetOnUpdateCallback(frontendServer.InvalidateCache)
			if subSiteService != nil {
				subSiteService.SetOnUpdateCallback(frontendServer.InvalidateCache)
			}
			r.Use(frontendServer.Middleware())
		}
	}

	// 注册路由
	registerRoutes(r, handlers, jwtAuth, adminAuth, apiKeyAuth, orgAuth, apiKeyService, subscriptionService, quotaPackageRepo, opsService, subSiteService, cfg, redisClient)

	return r
}

// registerRoutes 注册所有 HTTP 路由
func registerRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	jwtAuth middleware2.JWTAuthMiddleware,
	adminAuth middleware2.AdminAuthMiddleware,
	apiKeyAuth middleware2.APIKeyAuthMiddleware,
	orgAuth middleware2.OrgAuthMiddleware,
	apiKeyService *service.APIKeyService,
	subscriptionService *service.SubscriptionService,
	quotaPackageRepo service.QuotaPackageRepository,
	opsService *service.OpsService,
	subSiteService *service.SubSiteService,
	cfg *config.Config,
	redisClient *redis.Client,
) {
	// 通用路由（健康检查、状态等）
	routes.RegisterCommonRoutes(r, cfg)

	// API v1
	v1 := r.Group("/api/v1")

	// 注册各模块路由
	routes.RegisterAuthRoutes(v1, h, jwtAuth, redisClient)
	routes.RegisterUserRoutes(v1, h, jwtAuth)
	routes.RegisterAdminRoutes(v1, h, adminAuth)
	routes.RegisterOrgRoutes(v1, h, jwtAuth, orgAuth)
	routes.RegisterSubSiteAdminRoutes(v1, h, jwtAuth, subSiteService)
	routes.RegisterGatewayRoutes(r, h, apiKeyAuth, apiKeyService, subscriptionService, quotaPackageRepo, opsService, cfg)
}
