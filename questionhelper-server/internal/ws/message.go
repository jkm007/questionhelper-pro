package ws

import "encoding/json"

// Message WebSocket 消息
type Message struct {
	UserID uint            `json:"user_id"`
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
}

// MessageType 消息类型
const (
	MessageTypeNotification = "notification" // 通知消息
	MessageTypeExam         = "exam"         // 考试消息
	MessageTypeChat         = "chat"         // 聊天消息
	MessageTypeSystem       = "system"       // 系统消息
)

// NotificationData 通知数据
type NotificationData struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    int    `json:"type"`
}

// ExamData 考试数据
type ExamData struct {
	ExamID  uint   `json:"exam_id"`
	Title   string `json:"title"`
	Action  string `json:"action"` // start, end, remind
	Message string `json:"message"`
}

// SystemData 系统数据
type SystemData struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

// NewNotificationMessage 创建通知消息
func NewNotificationMessage(data NotificationData) *Message {
	dataBytes, _ := json.Marshal(data)
	return &Message{
		Type: MessageTypeNotification,
		Data: dataBytes,
	}
}

// NewExamMessage 创建考试消息
func NewExamMessage(data ExamData) *Message {
	dataBytes, _ := json.Marshal(data)
	return &Message{
		Type: MessageTypeExam,
		Data: dataBytes,
	}
}

// NewSystemMessage 创建系统消息
func NewSystemMessage(data SystemData) *Message {
	dataBytes, _ := json.Marshal(data)
	return &Message{
		Type: MessageTypeSystem,
		Data: dataBytes,
	}
}
