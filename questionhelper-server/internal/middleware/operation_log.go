package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/logger"
)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 读取请求体
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(startTime).Milliseconds()

		logger.Info("操作日志",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"status", c.Writer.Status(),
			"latency", latency,
		)

		// TODO: 异步写入数据库
	}
}
