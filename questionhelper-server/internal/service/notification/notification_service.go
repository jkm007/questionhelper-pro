package notification

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	notificationRepo "questionhelper-server/internal/repository/notification"
	"questionhelper-server/pkg/logger"
)

// ListNotifications 通知列表
func ListNotifications(userID uint, req *dto.NotificationListRequest) ([]dto.NotificationInfo, int64, error) {
	notifications, total, err := notificationRepo.List(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询通知列表失败: %w", err)
	}

	list := make([]dto.NotificationInfo, 0, len(notifications))
	for _, n := range notifications {
		list = append(list, toNotificationInfo(&n))
	}
	return list, total, nil
}

// GetUnreadCount 获取未读数量
func GetUnreadCount(userID uint) (int, error) {
	count, err := notificationRepo.CountUnread(userID)
	if err != nil {
		return 0, fmt.Errorf("查询未读数量失败: %w", err)
	}
	return int(count), nil
}

// MarkAsRead 标记已读
func MarkAsRead(id, userID uint) error {
	notification, err := notificationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return fmt.Errorf("查询通知失败: %w", err)
	}

	if notification.UserID != userID {
		return errors.New("无权操作此通知")
	}

	if err := notificationRepo.MarkAsRead(id); err != nil {
		return fmt.Errorf("标记已读失败: %w", err)
	}

	return nil
}

// MarkAllAsRead 全部标记已读
func MarkAllAsRead(userID uint) error {
	if err := notificationRepo.MarkAllAsRead(userID); err != nil {
		return fmt.Errorf("全部标记已读失败: %w", err)
	}

	logger.Infof("用户 %d 全部标记已读", userID)
	return nil
}

// DeleteNotification 删除通知
func DeleteNotification(id, userID uint) error {
	notification, err := notificationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return fmt.Errorf("查询通知失败: %w", err)
	}

	if notification.UserID != userID {
		return errors.New("无权删除此通知")
	}

	if err := notificationRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除通知失败: %w", err)
	}

	return nil
}

// CreateNotification 创建通知（内部调用）
func CreateNotification(userID uint, notificationType int8, title, content string, targetType string, targetID uint) error {
	notification := &model.Notification{
		UserID:     userID,
		Type:       notificationType,
		Title:      title,
		Content:    content,
		TargetType: targetType,
		TargetID:   targetID,
		IsRead:     false,
	}

	if err := notificationRepo.Create(notification); err != nil {
		return fmt.Errorf("创建通知失败: %w", err)
	}

	return nil
}

// toNotificationInfo 转换为 NotificationInfo DTO
func toNotificationInfo(n *model.Notification) dto.NotificationInfo {
	return dto.NotificationInfo{
		ID:         n.ID,
		UserID:     n.UserID,
		Type:       n.Type,
		Title:      n.Title,
		Content:    n.Content,
		TargetType: n.TargetType,
		TargetID:   n.TargetID,
		IsRead:     n.IsRead,
		CreatedAt:  n.CreatedAt,
	}
}
