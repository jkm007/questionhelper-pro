package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源（生产环境应限制）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket 返回 WebSocket 处理函数
// 客户端通过 ws://host/ws?token=<jwt_token> 连接
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 query 参数获取 token
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "A0003", "msg": "缺少认证令牌"})
			return
		}

		// 解析 JWT
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "A0003", "msg": "认证令牌无效"})
			return
		}

		// 检查 Token 黑名单
		ctx := c.Request.Context()
		blacklistKey := "qh:token:bl:" + claims.JTI
		exists, _ := database.RDB.Exists(ctx, blacklistKey).Result()
		if exists > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "A0003", "msg": "令牌已失效"})
			return
		}

		// 升级 HTTP 连接为 WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Errorf("WebSocket 升级失败: %v", err)
			return
		}

		// 创建客户端并注册到 Hub
		client := NewClient(hub, claims.UserID, conn)
		hub.Register <- client

		logger.Infof("WebSocket 连接建立: userID=%d", claims.UserID)

		// 启动读写协程
		go client.WritePump()
		go client.ReadPump()
	}
}
