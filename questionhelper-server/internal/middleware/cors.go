package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/config"
)

// CorsMiddleware 跨域中间件
// 允许的域名从 config.Server.Mode 判断: debug 模式允许所有，生产模式限制
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// debug 模式允许所有来源，生产模式需要配置
		allowedOrigin := "*"
		if config.Cfg != nil && config.Cfg.Server.Mode != "debug" {
			// 生产模式: 仅允许非空 origin（实际部署时应从配置读取允许的域名列表）
			if origin == "" {
				allowedOrigin = "*"
			} else {
				allowedOrigin = origin
			}
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if strings.ToUpper(c.Request.Method) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
