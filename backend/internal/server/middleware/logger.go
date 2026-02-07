package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 在请求开始时记录关键信息（用于调试 502 问题）
		path := c.Request.URL.Path
		if path == "/v1/messages" || path == "/antigravity/v1/messages" {
			log.Printf("[GIN-DEBUG] Request started: %s %s | Proto=%s | UA=%s | ContentType=%s",
				c.Request.Method,
				path,
				c.Request.Proto,
				c.GetHeader("User-Agent"),
				c.GetHeader("Content-Type"),
			)
		}

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方法
		method := c.Request.Method

		// 请求路径
		path := c.Request.URL.Path

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 日志格式: [时间] 状态码 | 延迟 | IP | 方法 路径
		log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		// 如果有错误，额外记录错误信息
		if len(c.Errors) > 0 {
			log.Printf("[GIN] Errors: %v", c.Errors.String())
		}
	}
}
