package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNoStore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("sets no-store headers for sensitive responses", func(t *testing.T) {
		router := gin.New()
		router.Use(NoStore())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, "no-store, no-cache, must-revalidate, private", rec.Header().Get("Cache-Control"))
		assert.Equal(t, "no-store", rec.Header().Get("CDN-Cache-Control"))
		assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
		assert.Equal(t, "0", rec.Header().Get("Expires"))
		assert.Equal(t, "no-store", rec.Header().Get("Surrogate-Control"))
		assert.Equal(t, "Authorization, Cookie", rec.Header().Get("Vary"))
	})

	t.Run("preserves existing vary values without duplicates", func(t *testing.T) {
		router := gin.New()
		router.Use(NoStore())
		router.GET("/test", func(c *gin.Context) {
			c.Header("Vary", "Origin, Authorization")
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, "Origin, Authorization, Cookie", rec.Header().Get("Vary"))
	})
}
