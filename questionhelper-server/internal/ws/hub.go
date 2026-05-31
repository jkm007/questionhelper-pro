package ws

import (
	"sync"
)

// Hub WebSocket 连接管理中心
type Hub struct {
	// 已注册的客户端
	clients map[uint]map[*Client]bool

	// 注册请求
	Register chan *Client

	// 注销请求
	Unregister chan *Client

	// 广播消息
	Broadcast chan *Message

	// 互斥锁
	mu sync.RWMutex
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.register(client)
		case client := <-h.Unregister:
			h.unregister(client)
		case message := <-h.Broadcast:
			h.broadcast(message)
		}
	}
}

// register 注册客户端
func (h *Hub) register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.UserID]; !ok {
		h.clients[client.UserID] = make(map[*Client]bool)
	}
	h.clients[client.UserID][client] = true
}

// unregister 注销客户端
func (h *Hub) unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[client.UserID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)
			if len(clients) == 0 {
				delete(h.clients, client.UserID)
			}
		}
	}
}

// broadcast 广播消息
func (h *Hub) broadcast(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if message.UserID > 0 {
		// 发送给指定用户
		if clients, ok := h.clients[message.UserID]; ok {
			for client := range clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
		}
	} else {
		// 广播给所有用户
		for _, clients := range h.clients {
			for client := range clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID uint, message *Message) {
	message.UserID = userID
	h.Broadcast <- message
}

// GetOnlineUsers 获取在线用户列表
func (h *Hub) GetOnlineUsers() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]uint, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}
