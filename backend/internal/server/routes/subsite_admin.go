package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterSubSiteAdminRoutes 注册分站站长后台专属路由。
// 中间件链：NoStore → JWTAuth → SubSiteOwnerRequired（三层鉴权保证 owner 只能看自己的分站）。
func RegisterSubSiteAdminRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
	subSiteService *service.SubSiteService,
) {
	sub := v1.Group("/subsite/admin/:siteId")
	sub.Use(middleware.NoStore())
	sub.Use(gin.HandlerFunc(jwtAuth))
	sub.Use(middleware.SubSiteOwnerRequired(subSiteService))

	sub.GET("/dashboard", h.SubSiteAdmin.Dashboard)
	sub.GET("/site", h.SubSiteAdmin.SiteDetail)
	sub.PUT("/site", h.SubSiteAdmin.UpdateSite)
	sub.GET("/users", h.SubSiteAdmin.ListUsers)
	sub.POST("/users/offline-topup", h.SubSiteAdmin.OfflineTopup)
	sub.GET("/orders", h.SubSiteAdmin.ListOrders)
	sub.GET("/usage", h.SubSiteAdmin.ListUsage)
	sub.GET("/ledger", h.SubSiteAdmin.Ledger)
}
