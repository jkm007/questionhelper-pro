package system

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== 系统设置(原有) ====================

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

// ==================== 分类设置 ====================

// GetClassSettings 获取班级设置
func GetClassSettings() (*dto.ClassSettings, error) {
	var settings []model.SystemSetting
	keys := []string{
		"class_max_students", "class_allow_self_join", "class_require_approval",
		"class_default_type", "class_enable_code", "class_code_length",
	}
	if err := database.DB.Where("`key` IN ?", keys).Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("获取班级设置失败: %w", err)
	}

	result := &dto.ClassSettings{
		MaxStudentsPerClass: 50,
		AllowSelfJoin:       true,
		RequireApproval:     false,
		DefaultClassType:    "normal",
		EnableClassCode:     true,
		ClassCodeLength:     6,
	}

	kv := make(map[string]string)
	for _, s := range settings {
		kv[s.Key] = s.Value
	}
	if v, ok := kv["class_max_students"]; ok {
		fmt.Sscanf(v, "%d", &result.MaxStudentsPerClass)
	}
	if v, ok := kv["class_allow_self_join"]; ok {
		result.AllowSelfJoin = v == "true" || v == "1"
	}
	if v, ok := kv["class_require_approval"]; ok {
		result.RequireApproval = v == "true" || v == "1"
	}
	if v, ok := kv["class_default_type"]; ok {
		result.DefaultClassType = v
	}
	if v, ok := kv["class_enable_code"]; ok {
		result.EnableClassCode = v == "true" || v == "1"
	}
	if v, ok := kv["class_code_length"]; ok {
		fmt.Sscanf(v, "%d", &result.ClassCodeLength)
	}
	return result, nil
}

// UpdateClassSettings 更新班级设置
func UpdateClassSettings(req *dto.ClassSettings) error {
	settings := map[string]string{
		"class_max_students":  fmt.Sprintf("%d", req.MaxStudentsPerClass),
		"class_allow_self_join": fmt.Sprintf("%v", req.AllowSelfJoin),
		"class_require_approval": fmt.Sprintf("%v", req.RequireApproval),
		"class_default_type": req.DefaultClassType,
		"class_enable_code":  fmt.Sprintf("%v", req.EnableClassCode),
		"class_code_length":  fmt.Sprintf("%d", req.ClassCodeLength),
	}
	return UpdateSettings(settings)
}

// GetResourceSettings 获取资源设置
func GetResourceSettings() (*dto.ResourceSettings, error) {
	var settings []model.SystemSetting
	keys := []string{
		"resource_max_upload_size", "resource_allowed_types", "resource_enable_preview",
		"resource_storage_quota", "resource_enable_compress", "resource_watermark",
	}
	if err := database.DB.Where("`key` IN ?", keys).Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("获取资源设置失败: %w", err)
	}

	result := &dto.ResourceSettings{
		MaxUploadSize:    10,
		AllowedFileTypes: []string{"jpg", "png", "gif", "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "mp4", "mp3"},
		EnablePreview:    true,
		StorageQuota:     10,
		EnableCompress:   true,
		WatermarkEnabled: false,
	}

	kv := make(map[string]string)
	for _, s := range settings {
		kv[s.Key] = s.Value
	}
	if v, ok := kv["resource_max_upload_size"]; ok {
		fmt.Sscanf(v, "%d", &result.MaxUploadSize)
	}
	if v, ok := kv["resource_allowed_types"]; ok {
		var types []string
		if err := json.Unmarshal([]byte(v), &types); err == nil {
			result.AllowedFileTypes = types
		}
	}
	if v, ok := kv["resource_enable_preview"]; ok {
		result.EnablePreview = v == "true" || v == "1"
	}
	if v, ok := kv["resource_storage_quota"]; ok {
		fmt.Sscanf(v, "%d", &result.StorageQuota)
	}
	if v, ok := kv["resource_enable_compress"]; ok {
		result.EnableCompress = v == "true" || v == "1"
	}
	if v, ok := kv["resource_watermark"]; ok {
		result.WatermarkEnabled = v == "true" || v == "1"
	}
	return result, nil
}

// UpdateResourceSettings 更新资源设置
func UpdateResourceSettings(req *dto.ResourceSettings) error {
	typesJSON, _ := json.Marshal(req.AllowedFileTypes)
	settings := map[string]string{
		"resource_max_upload_size": fmt.Sprintf("%d", req.MaxUploadSize),
		"resource_allowed_types":   string(typesJSON),
		"resource_enable_preview":  fmt.Sprintf("%v", req.EnablePreview),
		"resource_storage_quota":   fmt.Sprintf("%d", req.StorageQuota),
		"resource_enable_compress": fmt.Sprintf("%v", req.EnableCompress),
		"resource_watermark":       fmt.Sprintf("%v", req.WatermarkEnabled),
	}
	return UpdateSettings(settings)
}

// ==================== 系统日志 ====================

// ListSystemLogs 系统日志列表
func ListSystemLogs(query *dto.SystemLogQuery) ([]model.SystemLog, int64, error) {
	var logs []model.SystemLog
	var total int64

	db := database.DB.Model(&model.SystemLog{})
	if query.Level != "" {
		db = db.Where("level = ?", query.Level)
	}
	if query.Module != "" {
		db = db.Where("module = ?", query.Module)
	}
	if query.Keyword != "" {
		db = db.Where("message LIKE ? OR action LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.StartAt != "" {
		db = db.Where("created_at >= ?", query.StartAt)
	}
	if query.EndAt != "" {
		db = db.Where("created_at <= ?", query.EndAt)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := query.GetOffset()
	if err := db.Offset(offset).Limit(query.GetLimit()).
		Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ListErrorLogs 错误日志列表
func ListErrorLogs(query *dto.ErrorLogQuery) ([]model.ErrorLog, int64, error) {
	var logs []model.ErrorLog
	var total int64

	db := database.DB.Model(&model.ErrorLog{})
	if query.Level != "" {
		db = db.Where("level = ?", query.Level)
	}
	if query.Keyword != "" {
		db = db.Where("error_message LIKE ?", "%"+query.Keyword+"%")
	}
	if query.StartAt != "" {
		db = db.Where("occurred_at >= ?", query.StartAt)
	}
	if query.EndAt != "" {
		db = db.Where("occurred_at <= ?", query.EndAt)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := query.GetOffset()
	if err := db.Offset(offset).Limit(query.GetLimit()).
		Order("occurred_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// SearchLogs 日志搜索
func SearchLogs(query *dto.LogSearchRequest) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	var total int64

	switch query.LogType {
	case "system":
		var logs []model.SystemLog
		db := database.DB.Model(&model.SystemLog{})
		if query.Level != "" {
			db = db.Where("level = ?", query.Level)
		}
		if query.Module != "" {
			db = db.Where("module = ?", query.Module)
		}
		if query.Keyword != "" {
			db = db.Where("message LIKE ? OR action LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
		}
		if query.StartAt != "" {
			db = db.Where("created_at >= ?", query.StartAt)
		}
		if query.EndAt != "" {
			db = db.Where("created_at <= ?", query.EndAt)
		}
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		if err := db.Offset(query.GetOffset()).Limit(query.GetLimit()).
			Order("created_at DESC").Find(&logs).Error; err != nil {
			return nil, 0, err
		}
		for _, l := range logs {
			results = append(results, map[string]interface{}{
				"id": l.ID, "level": l.Level, "module": l.Module,
				"action": l.Action, "message": l.Message, "ip": l.IP,
				"created_at": l.CreatedAt, "type": "system",
			})
		}
	case "error":
		var logs []model.ErrorLog
		db := database.DB.Model(&model.ErrorLog{})
		if query.Level != "" {
			db = db.Where("level = ?", query.Level)
		}
		if query.Keyword != "" {
			db = db.Where("error_message LIKE ?", "%"+query.Keyword+"%")
		}
		if query.StartAt != "" {
			db = db.Where("occurred_at >= ?", query.StartAt)
		}
		if query.EndAt != "" {
			db = db.Where("occurred_at <= ?", query.EndAt)
		}
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		if err := db.Offset(query.GetOffset()).Limit(query.GetLimit()).
			Order("occurred_at DESC").Find(&logs).Error; err != nil {
			return nil, 0, err
		}
		for _, l := range logs {
			results = append(results, map[string]interface{}{
				"id": l.ID, "level": l.Level, "error_message": l.ErrorMessage,
				"url": l.URL, "method": l.Method, "ip": l.IP,
				"occurred_at": l.OccurredAt, "type": "error",
			})
		}
	default:
		return nil, 0, fmt.Errorf("不支持的日志类型: %s", query.LogType)
	}
	return results, total, nil
}

// ArchiveLogs 日志归档
func ArchiveLogs(req *dto.LogArchiveRequest) (int64, error) {
	var count int64
	switch req.LogType {
	case "system":
		result := database.DB.Where("created_at < ?", req.BeforeAt).Delete(&model.SystemLog{})
		count = result.RowsAffected
		if result.Error != nil {
			return 0, fmt.Errorf("归档系统日志失败: %w", result.Error)
		}
	case "error":
		result := database.DB.Where("occurred_at < ?", req.BeforeAt).Delete(&model.ErrorLog{})
		count = result.RowsAffected
		if result.Error != nil {
			return 0, fmt.Errorf("归档错误日志失败: %w", result.Error)
		}
	default:
		return 0, fmt.Errorf("不支持的日志类型: %s", req.LogType)
	}
	return count, nil
}

// GetLogStats 日志统计
func GetLogStats() (*dto.LogStatsResponse, error) {
	resp := &dto.LogStatsResponse{}

	// 系统日志总数
	database.DB.Model(&model.SystemLog{}).Count(&resp.TotalCount)

	// 今日日志数
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&model.SystemLog{}).Where("created_at >= ?", today).Count(&resp.TodayCount)

	// 错误数
	database.DB.Model(&model.ErrorLog{}).Count(&resp.ErrorCount)

	// 警告数
	database.DB.Model(&model.SystemLog{}).Where("level = ?", "warn").Count(&resp.WarnCount)

	// 级别统计
	var levelStats []dto.LevelStatItem
	database.DB.Model(&model.SystemLog{}).
		Select("level, count(*) as count").
		Group("level").Scan(&levelStats)
	resp.LevelStats = levelStats

	// 模块统计
	var moduleStats []dto.ModuleStatItem
	database.DB.Model(&model.SystemLog{}).
		Select("module, count(*) as count").
		Group("module").Scan(&moduleStats)
	resp.ModuleStats = moduleStats

	// 最近7天趋势
	var trendData []dto.TrendItem
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		database.DB.Model(&model.SystemLog{}).
			Where("created_at >= ? AND created_at < ?", d, d+" 23:59:59").
			Count(&count)
		trendData = append(trendData, dto.TrendItem{Date: d, Count: count})
	}
	resp.TrendData = trendData

	return resp, nil
}

// ==================== 通知渠道 ====================

// ListNotificationChannels 渠道列表
func ListNotificationChannels() ([]model.NotificationChannel, error) {
	var channels []model.NotificationChannel
	if err := database.DB.Order("priority DESC, id ASC").Find(&channels).Error; err != nil {
		return nil, fmt.Errorf("查询通知渠道失败: %w", err)
	}
	return channels, nil
}

// UpdateNotificationChannel 更新渠道
func UpdateNotificationChannel(id uint, req *dto.UpdateChannelRequest) error {
	var channel model.NotificationChannel
	if err := database.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("通知渠道不存在")
		}
		return fmt.Errorf("查询通知渠道失败: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}

	if err := database.DB.Model(&channel).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新通知渠道失败: %w", err)
	}
	return nil
}

// ==================== 数据备份 ====================

// CreateBackup 创建备份
func CreateBackup(req *dto.CreateBackupRequest) (*model.BackupRecord, error) {
	record := &model.BackupRecord{
		Type:   req.Type,
		Status: "pending",
	}
	if req.ConfigID > 0 {
		record.ConfigID = req.ConfigID
		var config model.BackupConfig
		if err := database.DB.First(&config, req.ConfigID).Error; err == nil {
			record.FilePath = fmt.Sprintf("%s/backup_%s_%d.sql", config.StoragePath, time.Now().Format("20060102_150405"), time.Now().Unix())
		}
	} else {
		record.FilePath = fmt.Sprintf("backups/backup_%s_%d.sql", time.Now().Format("20060102_150405"), time.Now().Unix())
	}

	if err := database.DB.Create(record).Error; err != nil {
		return nil, fmt.Errorf("创建备份记录失败: %w", err)
	}
	return record, nil
}

// ListBackupRecords 备份列表
func ListBackupRecords(query *dto.BackupListQuery) ([]model.BackupRecord, int64, error) {
	var records []model.BackupRecord
	var total int64

	db := database.DB.Model(&model.BackupRecord{})
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(query.GetOffset()).Limit(query.GetLimit()).
		Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// RestoreBackup 恢复备份
func RestoreBackup(id uint) error {
	var record model.BackupRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("备份记录不存在")
		}
		return fmt.Errorf("查询备份记录失败: %w", err)
	}
	if record.Status != "completed" {
		return fmt.Errorf("只能恢复已完成的备份")
	}
	// 实际恢复逻辑由后台任务执行，这里仅更新状态
	return database.DB.Model(&record).Update("status", "restoring").Error
}

// DeleteBackup 删除备份
func DeleteBackup(id uint) error {
	var record model.BackupRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("备份记录不存在")
		}
		return fmt.Errorf("查询备份记录失败: %w", err)
	}
	if record.Status == "running" {
		return fmt.Errorf("不能删除正在执行的备份")
	}
	return database.DB.Delete(&record).Error
}

// GetBackupConfigs 获取备份配置列表
func GetBackupConfigs() ([]model.BackupConfig, error) {
	var configs []model.BackupConfig
	if err := database.DB.Order("id ASC").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取备份配置失败: %w", err)
	}
	return configs, nil
}

// CreateOrUpdateBackupConfig 创建或更新备份配置
func CreateOrUpdateBackupConfig(id uint, req *dto.BackupConfigRequest) error {
	config := model.BackupConfig{
		Name:        req.Name,
		Type:        req.Type,
		Schedule:    req.Schedule,
		StoragePath: req.StoragePath,
		RetainDays:  req.RetainDays,
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	} else {
		config.IsActive = true
	}

	if id > 0 {
		return database.DB.Model(&model.BackupConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
			"name":         config.Name,
			"type":         config.Type,
			"schedule":     config.Schedule,
			"storage_path": config.StoragePath,
			"retain_days":  config.RetainDays,
			"is_active":    config.IsActive,
		}).Error
	}
	return database.DB.Create(&config).Error
}

// ==================== 功能开关 ====================

// ListFeatureFlags 功能开关列表
func ListFeatureFlags() ([]model.FeatureFlag, error) {
	var flags []model.FeatureFlag
	if err := database.DB.Order("id ASC").Find(&flags).Error; err != nil {
		return nil, fmt.Errorf("查询功能开关失败: %w", err)
	}
	return flags, nil
}

// UpdateFeatureFlag 更新功能开关
func UpdateFeatureFlag(key string, req *dto.UpdateFeatureRequest) error {
	var flag model.FeatureFlag
	if err := database.DB.Where("`key` = ?", key).First(&flag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("功能开关不存在: %s", key)
		}
		return fmt.Errorf("查询功能开关失败: %w", err)
	}

	updates := make(map[string]interface{})
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}

	if err := database.DB.Model(&flag).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新功能开关失败: %w", err)
	}
	return nil
}

// ==================== 安全配置 ====================

// GetSecurityConfigs 获取安全配置
func GetSecurityConfigs() ([]model.SecurityConfig, error) {
	var configs []model.SecurityConfig
	if err := database.DB.Order("id ASC").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取安全配置失败: %w", err)
	}
	return configs, nil
}

// UpdateSecurityConfigs 更新安全配置
func UpdateSecurityConfigs(req *dto.SecurityConfigRequest) error {
	for _, item := range req.Configs {
		config := model.SecurityConfig{
			ConfigKey:   item.ConfigKey,
			ConfigValue: item.ConfigValue,
			Description: item.Description,
		}
		if item.IsActive != nil {
			config.IsActive = *item.IsActive
		} else {
			config.IsActive = true
		}

		if err := database.DB.Where("config_key = ?", item.ConfigKey).
			Assign(map[string]interface{}{
				"config_value": config.ConfigValue,
				"description":  config.Description,
				"is_active":    config.IsActive,
			}).FirstOrCreate(&config).Error; err != nil {
			return fmt.Errorf("更新安全配置 %s 失败: %w", item.ConfigKey, err)
		}
	}
	return nil
}

// ==================== 存储配置 ====================

// ListStorageConfigs 存储配置列表
func ListStorageConfigs() ([]model.StorageConfig, error) {
	var configs []model.StorageConfig
	if err := database.DB.Order("is_default DESC, id ASC").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %w", err)
	}
	// 脱敏处理
	for i := range configs {
		configs[i].SecretKey = maskSensitiveValue(configs[i].SecretKey)
	}
	return configs, nil
}

// CreateStorageConfig 创建存储配置
func CreateStorageConfig(req *dto.StorageConfigRequest) error {
	config := model.StorageConfig{
		Name:      req.Name,
		Type:      req.Type,
		Endpoint:  req.Endpoint,
		Bucket:    req.Bucket,
		AccessKey: req.AccessKey,
		SecretKey: req.SecretKey,
		Region:    req.Region,
		BaseURL:   req.BaseURL,
	}
	if req.IsDefault != nil {
		config.IsDefault = *req.IsDefault
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	} else {
		config.IsActive = true
	}

	// 如果设置为默认，先取消其他默认
	if config.IsDefault {
		database.DB.Model(&model.StorageConfig{}).Where("is_default = ?", true).Update("is_default", false)
	}

	return database.DB.Create(&config).Error
}

// UpdateStorageConfig 更新存储配置
func UpdateStorageConfig(id uint, req *dto.StorageConfigRequest) error {
	var config model.StorageConfig
	if err := database.DB.First(&config, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("存储配置不存在")
		}
		return fmt.Errorf("查询存储配置失败: %w", err)
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"type":       req.Type,
		"endpoint":   req.Endpoint,
		"bucket":     req.Bucket,
		"access_key": req.AccessKey,
		"secret_key": req.SecretKey,
		"region":     req.Region,
		"base_url":   req.BaseURL,
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
		if *req.IsDefault {
			database.DB.Model(&model.StorageConfig{}).Where("is_default = ? AND id != ?", true, id).Update("is_default", false)
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	return database.DB.Model(&config).Updates(updates).Error
}

// ==================== 邮件配置 ====================

// GetEmailConfig 获取邮件配置
func GetEmailConfig() (*model.EmailConfig, error) {
	var config model.EmailConfig
	err := database.DB.Where("is_default = ?", true).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		// 返回空配置
		return &config, nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取邮件配置失败: %w", err)
	}
	// 脱敏处理
	config.Password = maskSensitiveValue(config.Password)
	return &config, nil
}

// UpdateEmailConfig 更新邮件配置
func UpdateEmailConfig(req *dto.EmailConfigRequest) error {
	config := model.EmailConfig{
		Name:        req.Name,
		SMTPHost:    req.SMTPHost,
		SMTPPort:    req.SMTPPort,
		Username:    req.Username,
		Password:    req.Password,
		FromAddress: req.FromAddress,
		FromName:    req.FromName,
	}
	if req.IsDefault != nil {
		config.IsDefault = *req.IsDefault
	} else {
		config.IsDefault = true
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	} else {
		config.IsActive = true
	}

	// 查找已有默认配置
	var existing model.EmailConfig
	if err := database.DB.Where("is_default = ?", true).First(&existing).Error; err == gorm.ErrRecordNotFound {
		return database.DB.Create(&config).Error
	} else if err != nil {
		return fmt.Errorf("查询邮件配置失败: %w", err)
	}

	return database.DB.Model(&existing).Updates(map[string]interface{}{
		"name":         config.Name,
		"smtp_host":    config.SMTPHost,
		"smtp_port":    config.SMTPPort,
		"username":     config.Username,
		"password":     config.Password,
		"from_address": config.FromAddress,
		"from_name":    config.FromName,
		"is_active":    config.IsActive,
	}).Error
}

// ListEmailTemplates 邮件模板列表
func ListEmailTemplates() ([]model.EmailTemplate, error) {
	var templates []model.EmailTemplate
	if err := database.DB.Order("id ASC").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("查询邮件模板失败: %w", err)
	}
	return templates, nil
}

// CreateEmailTemplate 创建邮件模板
func CreateEmailTemplate(req *dto.EmailTemplateRequest) error {
	template := model.EmailTemplate{
		Code:      req.Code,
		Name:      req.Name,
		Subject:   req.Subject,
		Body:      req.Body,
		Variables: req.Variables,
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	} else {
		template.IsActive = true
	}
	if err := database.DB.Create(&template).Error; err != nil {
		return fmt.Errorf("创建邮件模板失败: %w", err)
	}
	return nil
}

// ==================== 短信配置 ====================

// GetSMSConfig 获取短信配置
func GetSMSConfig() (*model.SMSConfig, error) {
	var config model.SMSConfig
	err := database.DB.Where("is_default = ?", true).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return &config, nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取短信配置失败: %w", err)
	}
	// 脱敏处理
	config.AccessKeySecret = maskSensitiveValue(config.AccessKeySecret)
	return &config, nil
}

// UpdateSMSConfig 更新短信配置
func UpdateSMSConfig(req *dto.SMSConfigRequest) error {
	config := model.SMSConfig{
		Provider:        req.Provider,
		AccessKeyID:     req.AccessKeyID,
		AccessKeySecret: req.AccessKeySecret,
		SignName:        req.SignName,
	}
	if req.IsDefault != nil {
		config.IsDefault = *req.IsDefault
	} else {
		config.IsDefault = true
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	} else {
		config.IsActive = true
	}

	var existing model.SMSConfig
	if err := database.DB.Where("is_default = ?", true).First(&existing).Error; err == gorm.ErrRecordNotFound {
		return database.DB.Create(&config).Error
	} else if err != nil {
		return fmt.Errorf("查询短信配置失败: %w", err)
	}

	return database.DB.Model(&existing).Updates(map[string]interface{}{
		"provider":         config.Provider,
		"access_key_id":    config.AccessKeyID,
		"access_key_secret": config.AccessKeySecret,
		"sign_name":        config.SignName,
		"is_active":        config.IsActive,
	}).Error
}

// ListSMSTemplates 短信模板列表
func ListSMSTemplates() ([]model.SMSTemplate, error) {
	var templates []model.SMSTemplate
	if err := database.DB.Order("id ASC").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("查询短信模板失败: %w", err)
	}
	return templates, nil
}

// CreateSMSTemplate 创建短信模板
func CreateSMSTemplate(req *dto.SMSTemplateRequest) error {
	template := model.SMSTemplate{
		Code:               req.Code,
		Name:               req.Name,
		Content:            req.Content,
		Variables:          req.Variables,
		ProviderTemplateID: req.ProviderTemplateID,
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	} else {
		template.IsActive = true
	}
	if err := database.DB.Create(&template).Error; err != nil {
		return fmt.Errorf("创建短信模板失败: %w", err)
	}
	return nil
}

// ==================== 缓存管理 ====================

// GetCacheStats 缓存统计
func GetCacheStats() (*dto.CacheStatsResponse, error) {
	resp := &dto.CacheStatsResponse{}

	// 从 Redis 获取统计信息
	rdb := database.RDB
	ctx := context.Background()

	// 键总数
	if info, err := rdb.Info(ctx, "keyspace").Result(); err == nil {
		// 解析 keyspace 信息
		resp.TotalKeys = parseRedisKeyCount(info)
	}

	// 内存使用
	if info, err := rdb.Info(ctx, "memory").Result(); err == nil {
		resp.UsedMemory = parseRedisMemory(info, "used_memory_human")
		resp.UsedMemoryPeak = parseRedisMemory(info, "used_memory_peak_human")
	}

	// 命中率
	if info, err := rdb.Info(ctx, "stats").Result(); err == nil {
		resp.HitRate = parseRedisHitRate(info)
	}

	// 缓存配置详情
	var cacheConfigs []model.CacheConfig
	database.DB.Order("id ASC").Find(&cacheConfigs)
	for _, c := range cacheConfigs {
		resp.CacheDetails = append(resp.CacheDetails, dto.CacheDetailItem{
			CacheKey:    c.CacheKey,
			TTL:         c.TTL,
			Description: c.Description,
			IsEnabled:   c.IsEnabled,
			LastCleared: c.LastClearedAt,
		})
	}

	return resp, nil
}

// ClearCache 清除缓存
func ClearCache(req *dto.ClearCacheRequest) error {
	ctx := context.Background()
	rdb := database.RDB

	if req.Pattern == "" {
		// 清除全部缓存(保留非缓存key)
		req.Pattern = "cache:*"
	}

	keys, err := rdb.Keys(ctx, req.Pattern).Result()
	if err != nil {
		return fmt.Errorf("查询缓存键失败: %w", err)
	}

	if len(keys) > 0 {
		if err := rdb.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("清除缓存失败: %w", err)
		}
	}

	// 更新清除时间
	now := time.Now()
	database.DB.Model(&model.CacheConfig{}).
		Where("cache_key LIKE ?", req.Pattern).
		Update("last_cleared_at", &now)

	return nil
}

// ==================== 主题配置 ====================

// GetThemeConfig 获取主题配置
func GetThemeConfig() (*model.ThemeConfig, error) {
	var config model.ThemeConfig
	err := database.DB.Where("is_default = ?", true).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return &model.ThemeConfig{
			Name:           "默认主题",
			PrimaryColor:   "#409EFF",
			SecondaryColor: "#67C23A",
			IsDefault:      true,
			IsActive:       true,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取主题配置失败: %w", err)
	}
	return &config, nil
}

// UpdateThemeConfig 更新主题配置
func UpdateThemeConfig(req *dto.ThemeConfigRequest) error {
	var existing model.ThemeConfig
	err := database.DB.Where("is_default = ?", true).First(&existing).Error

	config := model.ThemeConfig{
		Name:           req.Name,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		LogoPath:       req.LogoPath,
		FaviconPath:    req.FaviconPath,
		Config:         req.Config,
	}
	if req.IsDefault != nil {
		config.IsDefault = *req.IsDefault
	} else {
		config.IsDefault = true
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	} else {
		config.IsActive = true
	}

	if err == gorm.ErrRecordNotFound {
		return database.DB.Create(&config).Error
	}
	if err != nil {
		return fmt.Errorf("查询主题配置失败: %w", err)
	}

	return database.DB.Model(&existing).Updates(map[string]interface{}{
		"name":            config.Name,
		"primary_color":   config.PrimaryColor,
		"secondary_color": config.SecondaryColor,
		"logo_path":       config.LogoPath,
		"favicon_path":    config.FaviconPath,
		"is_active":       config.IsActive,
		"config":          config.Config,
	}).Error
}

// ==================== 告警管理 ====================

// ListAlertRules 告警规则列表
func ListAlertRules(query *dto.AlertRuleQuery) ([]model.LogAlertRule, int64, error) {
	var rules []model.LogAlertRule
	var total int64

	db := database.DB.Model(&model.LogAlertRule{})
	if query.Level != "" {
		db = db.Where("level = ?", query.Level)
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(query.GetOffset()).Limit(query.GetLimit()).
		Order("id ASC").Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// CreateAlertRule 创建告警规则
func CreateAlertRule(req *dto.AlertRuleRequest) error {
	rule := model.LogAlertRule{
		Name:       req.Name,
		Level:      req.Level,
		Module:     req.Module,
		Pattern:    req.Pattern,
		Threshold:  req.Threshold,
		Duration:   req.Duration,
		NotifyType: req.NotifyType,
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	} else {
		rule.IsActive = true
	}
	if err := database.DB.Create(&rule).Error; err != nil {
		return fmt.Errorf("创建告警规则失败: %w", err)
	}
	return nil
}

// UpdateAlertRule 更新告警规则
func UpdateAlertRule(id uint, req *dto.AlertRuleRequest) error {
	var rule model.LogAlertRule
	if err := database.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("告警规则不存在")
		}
		return fmt.Errorf("查询告警规则失败: %w", err)
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"level":       req.Level,
		"module":      req.Module,
		"pattern":     req.Pattern,
		"threshold":   req.Threshold,
		"duration":    req.Duration,
		"notify_type": req.NotifyType,
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := database.DB.Model(&rule).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新告警规则失败: %w", err)
	}
	return nil
}

// ListAlertRecords 告警记录列表
func ListAlertRecords(query *dto.AlertRecordQuery) ([]model.LogAlertRecord, int64, error) {
	var records []model.LogAlertRecord
	var total int64

	db := database.DB.Model(&model.LogAlertRecord{})
	if query.RuleID > 0 {
		db = db.Where("rule_id = ?", query.RuleID)
	}
	if query.Level != "" {
		db = db.Where("level = ?", query.Level)
	}
	if query.StartAt != "" {
		db = db.Where("created_at >= ?", query.StartAt)
	}
	if query.EndAt != "" {
		db = db.Where("created_at <= ?", query.EndAt)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(query.GetOffset()).Limit(query.GetLimit()).
		Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// ==================== 工具函数 ====================

func parseRedisKeyCount(info string) int64 {
	var count int64
	// 简单解析 db0:keys=N
	for _, line := range splitLines(info) {
		if contains(line, "keys=") {
			fmt.Sscanf(line, "db0:keys=%d", &count)
		}
	}
	return count
}

func parseRedisMemory(info, field string) string {
	for _, line := range splitLines(info) {
		if contains(line, field+":") {
			parts := splitByColon(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return "0"
}

func parseRedisHitRate(info string) float64 {
	var hits, misses int64
	for _, line := range splitLines(info) {
		if contains(line, "keyspace_hits:") {
			fmt.Sscanf(line, "keyspace_hits:%d", &hits)
		}
		if contains(line, "keyspace_misses:") {
			fmt.Sscanf(line, "keyspace_misses:%d", &misses)
		}
	}
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total) * 100
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			if i > start {
				lines = append(lines, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitByColon(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

// maskSensitiveValue 脱敏处理敏感配置值
// 非空值替换为 "****"，空值保持不变
func maskSensitiveValue(value string) string {
	if value == "" {
		return ""
	}
	return "****"
}
