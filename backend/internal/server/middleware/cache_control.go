package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const sensitiveCacheControl = "no-store, no-cache, must-revalidate, private"

// NoStore marks auth/session-bound responses as non-cacheable to prevent
// shared proxies or browsers from reusing one user's data for another user.
func NoStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		setNoStoreHeaders(c.Writer.Header())
		c.Next()
		setNoStoreHeaders(c.Writer.Header())
	}
}

func setNoStoreHeaders(headers http.Header) {
	if headers == nil {
		return
	}

	headers.Set("Cache-Control", sensitiveCacheControl)
	headers.Set("CDN-Cache-Control", "no-store")
	headers.Set("Pragma", "no-cache")
	headers.Set("Expires", "0")
	headers.Set("Surrogate-Control", "no-store")
	appendVary(headers, "Authorization")
	appendVary(headers, "Cookie")
}

func appendVary(headers http.Header, value string) {
	current := headers.Values("Vary")
	if len(current) == 0 {
		headers.Set("Vary", value)
		return
	}

	seen := make(map[string]struct{}, len(current)+1)
	values := make([]string, 0, len(current)+1)
	for _, item := range current {
		for _, part := range strings.Split(item, ",") {
			token := strings.TrimSpace(part)
			if token == "" {
				continue
			}
			if _, ok := seen[token]; ok {
				continue
			}
			seen[token] = struct{}{}
			values = append(values, token)
		}
	}

	if _, ok := seen[value]; !ok {
		values = append(values, value)
	}

	headers.Set("Vary", strings.Join(values, ", "))
}
