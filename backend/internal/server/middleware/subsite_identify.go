package middleware

import (
	"context"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeySubSite   ContextKey = "sub_site"
	ContextKeyIsSubSite ContextKey = "is_sub_site"
)

func SubSiteIdentify(subSiteService *service.SubSiteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ctxkey.SubSiteCacheKey, "main")
		c.Request = c.Request.WithContext(ctx)
		if subSiteService == nil {
			c.Next()
			return
		}

		host := strings.TrimSpace(c.GetHeader("X-Original-Host"))
		if host == "" {
			host = strings.TrimSpace(c.GetHeader("X-Forwarded-Host"))
		}
		if host == "" {
			host = strings.TrimSpace(c.Request.Host)
		}
		if host == "" {
			c.Next()
			return
		}
		if idx := strings.Index(host, ","); idx >= 0 {
			host = strings.TrimSpace(host[:idx])
		}

		site, err := subSiteService.ResolveByHost(c.Request.Context(), host)
		if err != nil || site == nil {
			c.Next()
			return
		}

		ctx = context.WithValue(c.Request.Context(), ctxkey.SubSite, site)
		ctx = context.WithValue(ctx, ctxkey.IsSubSite, true)
		ctx = context.WithValue(ctx, ctxkey.SubSiteCacheKey, "subsite:"+site.Slug+":"+site.UpdatedAt.UTC().Format("20060102150405"))
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(ContextKeySubSite), site)
		c.Set(string(ContextKeyIsSubSite), true)
		c.Next()
	}
}
