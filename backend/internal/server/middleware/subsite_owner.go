package middleware

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// CtxKeyOwnedSite is the gin.Context key for the owned *service.SubSite attached
// by SubSiteOwnerRequired. Handlers can fetch it to avoid a second DB round-trip.
const CtxKeyOwnedSite = "currentOwnedSite"

// SubSiteOwnerRequired 校验当前 JWT 用户是 :siteId 对应分站的 owner。
// 必须在 JWTAuth 之后使用。siteIDParam 默认为 "siteId"。
func SubSiteOwnerRequired(svc *service.SubSiteService, siteIDParam ...string) gin.HandlerFunc {
	param := "siteId"
	if len(siteIDParam) > 0 && siteIDParam[0] != "" {
		param = siteIDParam[0]
	}
	return func(c *gin.Context) {
		subject, ok := GetAuthSubjectFromContext(c)
		if !ok {
			AbortWithError(c, 401, "UNAUTHORIZED", "User not found in context")
			return
		}
		siteID, err := strconv.ParseInt(c.Param(param), 10, 64)
		if err != nil || siteID <= 0 {
			response.BadRequest(c, "invalid sub-site ID")
			c.Abort()
			return
		}
		site, err := svc.AuthorizeOwner(c.Request.Context(), subject.UserID, siteID)
		if err != nil {
			response.ErrorFrom(c, err)
			c.Abort()
			return
		}
		c.Set(CtxKeyOwnedSite, site)
		c.Next()
	}
}

// OwnedSiteFromContext 从 gin.Context 取 SubSiteOwnerRequired 注入的 site。
func OwnedSiteFromContext(c *gin.Context) (*service.SubSite, bool) {
	v, ok := c.Get(CtxKeyOwnedSite)
	if !ok {
		return nil, false
	}
	site, ok := v.(*service.SubSite)
	return site, ok
}
