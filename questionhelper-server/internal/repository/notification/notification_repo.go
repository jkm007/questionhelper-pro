package notification

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

func FindByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	err := database.DB.First(&notification, id).Error
	return &notification, err
}

func Create(notification *model.Notification) error {
	return database.DB.Create(notification).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.Notification{}, id).Error
}

func List(userID uint, req *dto.NotificationListRequest) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	db := database.DB.Model(&model.Notification{}).Where("user_id = ?", userID)

	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.IsRead != nil {
		db = db.Where("is_read = ?", *req.IsRead)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&notifications).Error

	return notifications, total, err
}

func CountUnread(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func MarkAsRead(id uint) error {
	return database.DB.Model(&model.Notification{}).Where("id = ?", id).
		Update("is_read", true).Error
}

func MarkAllAsRead(userID uint) error {
	return database.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
