package system

import (
	"fmt"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// GetSettings 获取系统设置
func GetSettings() (map[string]string, error) {
	var settings []model.SystemSetting
	if err := database.DB.Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("获取系统设置失败: %w", err)
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// UpdateSettings 更新系统设置
func UpdateSettings(settings map[string]string) error {
	for key, value := range settings {
		setting := model.SystemSetting{
			Key:   key,
			Value: value,
		}
		if err := database.DB.Where("key = ?", key).Assign(setting).FirstOrCreate(&setting).Error; err != nil {
			return fmt.Errorf("更新设置 %s 失败: %w", key, err)
		}
	}
	return nil
}

// ListOperationLogs 操作日志列表
func ListOperationLogs(page, pageSize int) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	db := database.DB.Model(&model.OperationLog{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	if err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ListLoginLogs 登录日志列表
func ListLoginLogs(page, pageSize int) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	db := database.DB.Model(&model.LoginLog{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	if err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
