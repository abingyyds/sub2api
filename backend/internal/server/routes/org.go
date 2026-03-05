package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterOrgRoutes registers organization management routes
func RegisterOrgRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
	orgAuth middleware.OrgAuthMiddleware,
) {
	org := v1.Group("/org")
	org.Use(gin.HandlerFunc(jwtAuth))
	org.Use(gin.HandlerFunc(orgAuth))
	{
		// Dashboard
		org.GET("/dashboard", h.Org.Dashboard.GetDashboard)

		// Member management
		members := org.Group("/members")
		{
			members.GET("", h.Org.Member.List)
			members.POST("", h.Org.Member.Create)
			members.GET("/:id", h.Org.Member.GetByID)
			members.PUT("/:id", h.Org.Member.Update)
			members.DELETE("/:id", h.Org.Member.Delete)
			members.POST("/:id/suspend", h.Org.Member.Suspend)
		}

		// Project management
		projects := org.Group("/projects")
		{
			projects.GET("", h.Org.Project.List)
			projects.POST("", h.Org.Project.Create)
			projects.GET("/:id", h.Org.Project.GetByID)
			projects.PUT("/:id", h.Org.Project.Update)
			projects.DELETE("/:id", h.Org.Project.Delete)
		}

		// Audit logs
		auditLogs := org.Group("/audit-logs")
		{
			auditLogs.GET("", h.Org.AuditLog.List)
			auditLogs.GET("/flagged-count", h.Org.AuditLog.CountFlagged)
			auditLogs.GET("/:id", h.Org.AuditLog.GetByID)
		}

		// Audit config
		org.GET("/audit-config", h.Org.AuditLog.GetAuditConfig)
		org.PUT("/audit-config", h.Org.AuditLog.UpdateAuditConfig)
	}
}
