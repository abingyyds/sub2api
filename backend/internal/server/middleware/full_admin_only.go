package middleware

import (
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// FullAdminOnly 仅允许完整管理员（admin）访问，sub_admin 不可访问
func FullAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRoleFromContext(c)
		if !ok {
			AbortWithError(c, 401, "UNAUTHORIZED", "User not found in context")
			return
		}

		if role != service.RoleAdmin {
			AbortWithError(c, 403, "FORBIDDEN", "Full admin access required")
			return
		}

		c.Next()
	}
}
