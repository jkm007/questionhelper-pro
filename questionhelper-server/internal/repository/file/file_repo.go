package file

import (
	"fmt"
	"time"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== File ====================

// FindFileByID 根据ID查找文件
func FindFileByID(id uint) (*model.File, error) {
	var file model.File
	err := database.DB.First(&file, id).Error
	return &file, err
}

// FindFileByMD5 根据MD5查找已有文件
func FindFileByMD5(md5 string) (*model.File, error) {
	var file model.File
	err := database.DB.Where("md5 = ? AND status = ?", md5, "active").First(&file).Error
	return &file, err
}

// CreateFile 创建文件记录
func CreateFile(file *model.File) error {
	return database.DB.Create(file).Error
}

// UpdateFile 更新文件记录
func UpdateFile(file *model.File) error {
	return database.DB.Save(file).Error
}

// DeleteFileByID 删除文件记录（软删除）
func DeleteFileByID(id uint) error {
	return database.DB.Delete(&model.File{}, id).Error
}

// ListFiles 文件列表
func ListFiles(uploaderID *uint, page, pageSize int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	db := database.DB.Model(&model.File{}).Where("status = ?", "active")
	if uploaderID != nil {
		db = db.Where("uploader_id = ?", *uploaderID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	if err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// ==================== File Reference ====================

// CreateFileReference 创建文件引用
func CreateFileReference(ref *model.FileReference) error {
	return database.DB.Create(ref).Error
}

// FindFileReferenceByID 根据ID查找文件引用
func FindFileReferenceByID(id uint) (*model.FileReference, error) {
	var ref model.FileReference
	err := database.DB.First(&ref, id).Error
	return &ref, err
}

// DeleteFileReferenceByID 删除文件引用
func DeleteFileReferenceByID(id uint) error {
	return database.DB.Delete(&model.FileReference{}, id).Error
}

// ListFileReferences 获取文件引用列表
func ListFileReferences(fileID uint) ([]model.FileReference, error) {
	var refs []model.FileReference
	err := database.DB.Where("file_id = ?", fileID).
		Order("created_at DESC").Find(&refs).Error
	return refs, err
}

// IncrementReferenceCount 增加引用计数
func IncrementReferenceCount(fileID uint) error {
	return database.DB.Model(&model.File{}).
		Where("id = ?", fileID).
		UpdateColumn("reference_count", database.DB.Raw("reference_count + 1")).Error
}

// DecrementReferenceCount 减少引用计数
func DecrementReferenceCount(fileID uint) error {
	return database.DB.Model(&model.File{}).
		Where("id = ? AND reference_count > 0", fileID).
		UpdateColumn("reference_count", database.DB.Raw("reference_count - 1")).Error
}

// CountFileReferences 统计文件引用数
func CountFileReferences(fileID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.FileReference{}).
		Where("file_id = ?", fileID).Count(&count).Error
	return count, err
}

// ==================== File Access Log ====================

// CreateFileAccessLog 创建文件访问日志
func CreateFileAccessLog(log *model.FileAccessLog) error {
	return database.DB.Create(log).Error
}

// ListFileAccessLogs 获取文件访问日志
func ListFileAccessLogs(fileID uint, page, pageSize int) ([]model.FileAccessLog, int64, error) {
	var logs []model.FileAccessLog
	var total int64

	db := database.DB.Model(&model.FileAccessLog{}).Where("file_id = ?", fileID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	if err := db.Offset(offset).Limit(pageSize).
		Order("accessed_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ==================== Hotlink Protection Rule ====================

// ListHotlinkRules 获取防盗链规则列表
func ListHotlinkRules() ([]model.HotlinkProtectionRule, error) {
	var rules []model.HotlinkProtectionRule
	err := database.DB.Order("created_at DESC").Find(&rules).Error
	return rules, err
}

// FindHotlinkRuleByID 根据ID查找防盗链规则
func FindHotlinkRuleByID(id uint) (*model.HotlinkProtectionRule, error) {
	var rule model.HotlinkProtectionRule
	err := database.DB.First(&rule, id).Error
	return &rule, err
}

// CreateHotlinkRule 创建防盗链规则
func CreateHotlinkRule(rule *model.HotlinkProtectionRule) error {
	return database.DB.Create(rule).Error
}

// UpdateHotlinkRule 更新防盗链规则
func UpdateHotlinkRule(rule *model.HotlinkProtectionRule) error {
	return database.DB.Save(rule).Error
}

// DeleteHotlinkRuleByID 删除防盗链规则
func DeleteHotlinkRuleByID(id uint) error {
	return database.DB.Delete(&model.HotlinkProtectionRule{}, id).Error
}

// ==================== Watermark Config ====================

// ListWatermarkConfigs 获取水印配置列表
func ListWatermarkConfigs() ([]model.WatermarkConfig, error) {
	var configs []model.WatermarkConfig
	err := database.DB.Order("created_at DESC").Find(&configs).Error
	return configs, err
}

// FindWatermarkConfigByID 根据ID查找水印配置
func FindWatermarkConfigByID(id uint) (*model.WatermarkConfig, error) {
	var config model.WatermarkConfig
	err := database.DB.First(&config, id).Error
	return &config, err
}

// CreateWatermarkConfig 创建水印配置
func CreateWatermarkConfig(config *model.WatermarkConfig) error {
	return database.DB.Create(config).Error
}

// UpdateWatermarkConfig 更新水印配置
func UpdateWatermarkConfig(config *model.WatermarkConfig) error {
	return database.DB.Save(config).Error
}

// ==================== File Cleanup Log ====================

// CreateCleanupLog 创建清理日志
func CreateCleanupLog(log *model.FileCleanupLog) error {
	return database.DB.Create(log).Error
}

// UpdateCleanupLog 更新清理日志
func UpdateCleanupLog(log *model.FileCleanupLog) error {
	return database.DB.Save(log).Error
}

// ListCleanupLogs 获取清理日志列表
func ListCleanupLogs(page, pageSize int) ([]model.FileCleanupLog, int64, error) {
	var logs []model.FileCleanupLog
	var total int64

	db := database.DB.Model(&model.FileCleanupLog{})

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

// FindOrphanFiles 查找孤立文件（无引用且超过指定天数）
func FindOrphanFiles(days int) ([]model.File, error) {
	var files []model.File
	cutoff := time.Now().AddDate(0, 0, -days)
	err := database.DB.Where("reference_count = 0 AND status = ? AND created_at < ?", "active", cutoff).
		Find(&files).Error
	return files, err
}

// ==================== Storage Statistics ====================

// CountFilesByType 按类型统计文件数量和大小
func CountFilesByType(extensions []string) (count int64, size int64, err error) {
	db := database.DB.Model(&model.File{}).Where("status = ?", "active")
	if len(extensions) > 0 {
		db = db.Where("extension IN ?", extensions)
	}
	if err = db.Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = db.Select("COALESCE(SUM(size), 0)").Scan(&size).Error
	}
	return
}

// CountTotalFiles 统计总文件数和大小
func CountTotalFiles() (count int64, size int64, err error) {
	db := database.DB.Model(&model.File{}).Where("status = ?", "active")
	if err = db.Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = db.Select("COALESCE(SUM(size), 0)").Scan(&size).Error
	}
	return
}

// CountOrphanFiles 统计孤立文件数和大小
func CountOrphanFiles() (count int64, size int64, err error) {
	db := database.DB.Model(&model.File{}).Where("reference_count = 0 AND status = ?", "active")
	if err = db.Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = db.Select("COALESCE(SUM(size), 0)").Scan(&size).Error
	}
	return
}

// GetDailyUploadStats 获取每日上传统计
func GetDailyUploadStats(days int) ([]struct {
	Date  string
	Count int64
	Size  int64
}, error) {
	var results []struct {
		Date  string
		Count int64
		Size  int64
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	err := database.DB.Model(&model.File{}).
		Where("status = ? AND created_at >= ?", "active", cutoff).
		Select("DATE(created_at) as date, COUNT(*) as count, COALESCE(SUM(size), 0) as size").
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&results).Error
	return results, err
}

// GetTopUploaders 获取上传文件最多的用户
func GetTopUploaders(limit int) ([]struct {
	UserID    uint
	Count     int64
	TotalSize int64
}, error) {
	var results []struct {
		UserID    uint
		Count     int64
		TotalSize int64
	}
	err := database.DB.Model(&model.File{}).
		Where("status = ?", "active").
		Select("uploader_id as user_id, COUNT(*) as count, COALESCE(SUM(size), 0) as total_size").
		Group("uploader_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

// ListFilesWithFilter 带过滤条件的文件列表
func ListFilesWithFilter(page, pageSize int, fileType, keyword, sortBy, sortOrder string) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	db := database.DB.Model(&model.File{}).Where("status = ?", "active")

	if fileType != "" {
		db = db.Where("extension = ?", fileType)
	}
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 白名单校验排序字段，防止 SQL 注入
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"size":       true,
		"name":       true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	if err := db.Offset(offset).Limit(pageSize).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// BatchDeleteFiles 批量删除文件记录
func BatchDeleteFiles(ids []uint) error {
	return database.DB.Delete(&model.File{}, ids).Error
}
