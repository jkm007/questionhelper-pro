package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/response"
	"questionhelper-server/pkg/sensitive"
)

func SensitiveFilterMiddleware() gin.HandlerFunc {
	filter := sensitive.NewFilter()

	return func(c *gin.Context) {
		// 只对 POST/PUT 请求进行敏感词过滤
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.Next()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			content := string(body)
			if filter.HasSensitive(content) {
				response.Error(c, 400, "内容包含敏感词，请修改后重试")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ReplaceSensitive 替换敏感词
func ReplaceSensitive(content string) string {
	filter := sensitive.NewFilter()
	return filter.Replace(content, '*')
}
