//go:build unit

package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/handler"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRegisterAuthRoutesAuthMeUnauthorizedIncludesNoStoreHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	v1 := router.Group("/api/v1")

	h := &handler.Handlers{
		Auth:         &handler.AuthHandler{},
		Setting:      &handler.SettingHandler{},
		Announcement: &handler.AnnouncementHandler{},
	}

	jwtAuth := servermiddleware.JWTAuthMiddleware(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    "UNAUTHORIZED",
			"message": "Authorization header is required",
		})
	})

	RegisterAuthRoutes(v1, h, jwtAuth, (*redis.Client)(nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Equal(t, "no-store, no-cache, must-revalidate, private", rec.Header().Get("Cache-Control"))
	require.Equal(t, "no-store", rec.Header().Get("CDN-Cache-Control"))
	require.Equal(t, "no-cache", rec.Header().Get("Pragma"))
	require.Equal(t, "0", rec.Header().Get("Expires"))
	require.Equal(t, "no-store", rec.Header().Get("Surrogate-Control"))
	require.Equal(t, "Authorization, Cookie", rec.Header().Get("Vary"))
}
