package log

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== OperationLog ====================

// CreateOperationLog 创建操作日志
func CreateOperationLog(log *model.OperationLog) error {
	return database.DB.Create(log).Error
}

// FindOperationLogByID 根据ID查找操作日志
func FindOperationLogByID(id uint) (*model.OperationLog, error) {
	var log model.OperationLog
	err := database.DB.First(&log, id).Error
	return &log, err
}

// ListOperationLogs 操作日志列表
func ListOperationLogs(req *dto.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	db := database.DB.Model(&model.OperationLog{})

	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.Module != "" {
		db = db.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.StartTime != nil {
		db = db.Where("created_at >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		db = db.Where("created_at <= ?", *req.EndTime)
	}
	if req.IP != "" {
		db = db.Where("ip = ?", req.IP)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("id DESC").
		Find(&logs).Error

	return logs, total, err
}

// DeleteOperationLogsBefore 删除指定时间之前的操作日志
func DeleteOperationLogsBefore(before time.Time) (int64, error) {
	result := database.DB.Where("created_at < ?", before).Delete(&model.OperationLog{})
	return result.RowsAffected, result.Error
}

// ==================== LoginLog ====================

// CreateLoginLog 创建登录日志
func CreateLoginLog(log *model.LoginLog) error {
	return database.DB.Create(log).Error
}

// ListLoginLogs 登录日志列表
func ListLoginLogs(req *dto.LoginLogListRequest) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	db := database.DB.Model(&model.LoginLog{})

	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.StartTime != nil {
		db = db.Where("created_at >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		db = db.Where("created_at <= ?", *req.EndTime)
	}
	if req.IP != "" {
		db = db.Where("ip = ?", req.IP)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("id DESC").
		Find(&logs).Error

	return logs, total, err
}

// DeleteLoginLogsBefore 删除指定时间之前的登录日志
func DeleteLoginLogsBefore(before time.Time) (int64, error) {
	result := database.DB.Where("created_at < ?", before).Delete(&model.LoginLog{})
	return result.RowsAffected, result.Error
}
