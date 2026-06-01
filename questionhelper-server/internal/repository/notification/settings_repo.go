package notification

import (
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

func FindSettingsByUserID(userID uint) ([]model.NotificationSetting, error) {
	var settings []model.NotificationSetting
	err := database.DB.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

func UpsertSetting(setting *model.NotificationSetting) error {
	return database.DB.Where("user_id = ? AND type = ? AND channel = ?",
		setting.UserID, setting.Type, setting.Channel).
		Assign(map[string]interface{}{
			"enabled":        setting.Enabled,
			"do_not_disturb": setting.DoNotDisturb,
			"disturb_start":  setting.DisturbStart,
			"disturb_end":    setting.DisturbEnd,
		}).FirstOrCreate(setting).Error
}
