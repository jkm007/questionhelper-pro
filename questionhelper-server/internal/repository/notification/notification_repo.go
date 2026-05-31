package notification

import (
	"time"

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

func CreateBatch(notifications []model.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return database.DB.Create(&notifications).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.Notification{}, id).Error
}

func BatchDelete(ids []uint, userID uint) error {
	return database.DB.Where("id IN ? AND user_id = ?", ids, userID).
		Delete(&model.Notification{}).Error
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

func BatchMarkAsRead(ids []uint, userID uint) error {
	return database.DB.Model(&model.Notification{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Update("is_read", true).Error
}

func Recall(id uint) error {
	now := time.Now()
	return database.DB.Model(&model.Notification{}).Where("id = ? AND is_read = ?", id, false).
		Updates(map[string]interface{}{
			"is_recalled": true,
			"recalled_at": now,
		}).Error
}

// ==================== 通知模板 ====================

func FindTemplateByID(id uint) (*model.NotificationTemplate, error) {
	var tpl model.NotificationTemplate
	err := database.DB.First(&tpl, id).Error
	return &tpl, err
}

func FindTemplateByCode(code string) (*model.NotificationTemplate, error) {
	var tpl model.NotificationTemplate
	err := database.DB.Where("code = ?", code).First(&tpl).Error
	return &tpl, err
}

func CreateTemplate(tpl *model.NotificationTemplate) error {
	return database.DB.Create(tpl).Error
}

func UpdateTemplate(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&model.NotificationTemplate{}).Where("id = ?", id).
		Updates(updates).Error
}

func DeleteTemplate(id uint) error {
	return database.DB.Delete(&model.NotificationTemplate{}, id).Error
}

func ListTemplates(req *dto.NotificationTemplateListRequest) ([]model.NotificationTemplate, int64, error) {
	var templates []model.NotificationTemplate
	var total int64

	db := database.DB.Model(&model.NotificationTemplate{})

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.Channel != "" {
		db = db.Where("channel = ?", req.Channel)
	}
	if req.IsActive != nil {
		db = db.Where("is_active = ?", *req.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&templates).Error

	return templates, total, err
}

// ==================== 定时通知 ====================

func FindScheduledByID(id uint) (*model.ScheduledNotification, error) {
	var scheduled model.ScheduledNotification
	err := database.DB.First(&scheduled, id).Error
	return &scheduled, err
}

func CreateScheduled(scheduled *model.ScheduledNotification) error {
	return database.DB.Create(scheduled).Error
}

func DeleteScheduled(id uint) error {
	return database.DB.Model(&model.ScheduledNotification{}).
		Where("id = ? AND status = ?", id, 0).
		Update("status", 2).Error
}

func ListScheduled(userID uint, req *dto.ScheduledListRequest) ([]model.ScheduledNotification, int64, error) {
	var scheduled []model.ScheduledNotification
	var total int64

	db := database.DB.Model(&model.ScheduledNotification{}).
		Where("sender_id = ?", userID)

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("scheduled_at DESC").Find(&scheduled).Error

	return scheduled, total, err
}

func ListPendingScheduled(before time.Time) ([]model.ScheduledNotification, error) {
	var scheduled []model.ScheduledNotification
	err := database.DB.Model(&model.ScheduledNotification{}).
		Where("status = ? AND scheduled_at <= ?", 0, before).
		Find(&scheduled).Error
	return scheduled, err
}

func UpdateScheduledStatus(id uint, status int8, errorMsg string) error {
	updates := map[string]interface{}{"status": status}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return database.DB.Model(&model.ScheduledNotification{}).Where("id = ?", id).
		Updates(updates).Error
}

// ==================== 通知设置 ====================

func FindSettingsByUserID(userID uint) ([]model.NotificationSetting, error) {
	var settings []model.NotificationSetting
	err := database.DB.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

func UpsertSetting(setting *model.NotificationSetting) error {
	return database.DB.Where("user_id = ? AND type = ? AND channel = ?",
		setting.UserID, setting.Type, setting.Channel).
		Assign(map[string]interface{}{
			"enabled":       setting.Enabled,
			"do_not_disturb": setting.DoNotDisturb,
			"disturb_start": setting.DisturbStart,
			"disturb_end":   setting.DisturbEnd,
		}).FirstOrCreate(setting).Error
}

// ==================== 通知统计 ====================

func StatsByUser(userID uint) (total, unread, read int64, err error) {
	db := database.DB.Model(&model.Notification{}).Where("user_id = ?", userID)

	if err = db.Count(&total).Error; err != nil {
		return
	}
	if err = db.Where("is_read = ?", false).Count(&unread).Error; err != nil {
		return
	}
	read = total - unread
	return
}

func StatsByType(userID uint) ([]dto.TypeStatItem, error) {
	var results []dto.TypeStatItem
	err := database.DB.Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Select("type, count(*) as count").
		Group("type").
		Scan(&results).Error
	return results, err
}

func StatsByChannel(userID uint) ([]dto.ChannelStatItem, error) {
	var results []dto.ChannelStatItem
	err := database.DB.Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Select("channel, count(*) as count").
		Group("channel").
		Scan(&results).Error
	return results, err
}

func StatsDaily(userID uint, days int) ([]dto.DailyStatItem, error) {
	var results []dto.DailyStatItem
	since := time.Now().AddDate(0, 0, -days)
	err := database.DB.Model(&model.Notification{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Select("DATE(created_at) as date, count(*) as count").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}

// ==================== 用户查询（群发用）====================

func FindUserIDsByRole(roleID uint) ([]uint, error) {
	var ids []uint
	err := database.DB.Table("user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &ids).Error
	return ids, err
}

func FindUserIDsByClass(classID uint) ([]uint, error) {
	var ids []uint
	err := database.DB.Table("class_members").
		Where("class_id = ?", classID).
		Pluck("user_id", &ids).Error
	return ids, err
}

func FindAllUserIDs() ([]uint, error) {
	var ids []uint
	err := database.DB.Model(&model.User{}).
		Where("status = ?", 1).
		Pluck("id", &ids).Error
	return ids, err
}

// ==================== 通知渠道 ====================

func ListChannels(req *dto.ChannelListRequest) ([]model.NotificationChannel, int64, error) {
	var channels []model.NotificationChannel
	var total int64

	db := database.DB.Model(&model.NotificationChannel{})

	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *req.IsEnabled)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("priority DESC, id ASC").Find(&channels).Error

	return channels, total, err
}

func FindChannelByID(id uint) (*model.NotificationChannel, error) {
	var channel model.NotificationChannel
	err := database.DB.First(&channel, id).Error
	return &channel, err
}

func UpdateChannel(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&model.NotificationChannel{}).Where("id = ?", id).
		Updates(updates).Error
}

// IsRecalled 检查通知是否已撤回
func IsRecalled(id uint) (bool, error) {
	var n model.Notification
	err := database.DB.Select("is_recalled").First(&n, id).Error
	if err != nil {
		return false, err
	}
	return n.IsRecalled, nil
}

// ListAllNotifications 管理员查询所有通知
func ListAllNotifications(req *dto.NotificationListRequest) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	db := database.DB.Model(&model.Notification{})

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
