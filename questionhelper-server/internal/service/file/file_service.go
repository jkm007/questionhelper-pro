package file

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	fileRepo "questionhelper-server/internal/repository/file"
	"questionhelper-server/pkg/logger"
)

const (
	uploadDir   = "./uploads"
	maxFileSize = 10 << 20 // 10MB

	imageMaxSize    = 5 << 20 // 5MB for images
	batchMaxCount   = 10      // max files per batch
	thumbnailDir    = "./uploads/thumbnails"
	thumbnailWidth  = 200
	thumbnailHeight = 200
)

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".txt": true, ".csv": true,
	".mp3": true, ".mp4": true,
}

var imageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

var documentExtensions = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".txt": true, ".csv": true,
}

var videoExtensions = map[string]bool{
	".mp4": true,
}

var audioExtensions = map[string]bool{
	".mp3": true,
}

// ==================== Basic Upload ====================

// UploadFile 上传文件
func UploadFile(uploaderID uint, fileName string, fileSize int64, fileType string, reader io.Reader) (*model.File, error) {
	// 检查文件大小限制
	if fileSize > maxFileSize {
		return nil, fmt.Errorf("文件大小超过限制: %d bytes (最大 %d bytes)", fileSize, maxFileSize)
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 检查文件扩展名白名单
	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedExtensions[ext] {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 使用UUID生成安全的文件名，避免路径遍历
	newFileName := uuid.New().String() + ext
	filePath := filepath.Clean(filepath.Join(uploadDir, newFileName))

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 计算MD5并保存
	hash := md5.New()
	writer := io.MultiWriter(dst, hash)

	if _, err := io.Copy(writer, reader); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	// 检查MD5去重
	existing, err := fileRepo.FindFileByMD5(md5Hash)
	if err == nil && existing != nil {
		// 文件已存在，删除刚保存的文件
		os.Remove(filePath)
		logger.Infof("文件已存在(MD5去重): %s -> %s", fileName, existing.Name)
		return existing, nil
	}

	// 保存文件信息到数据库
	file := &model.File{
		Name:       newFileName,
		Original:   fileName,
		Path:       filePath,
		URL:        "/uploads/" + newFileName,
		Size:       fileSize,
		Type:       fileType,
		Extension:  ext,
		MD5:        md5Hash,
		UploaderID: uploaderID,
		Status:     "active",
	}

	if err := fileRepo.CreateFile(file); err != nil {
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	logger.Infof("文件上传成功: %s (MD5: %s)", fileName, md5Hash)
	return file, nil
}

// ==================== Image Upload ====================

// UploadImage 上传图片（带压缩和缩略图生成）
func UploadImage(uploaderID uint, fileName string, fileSize int64, fileType string, reader io.Reader) (*dto.ImageUploadResponse, error) {
	// 检查文件大小限制
	if fileSize > imageMaxSize {
		return nil, fmt.Errorf("图片大小超过限制: %d bytes (最大 %d bytes)", fileSize, imageMaxSize)
	}

	// 检查是否为图片类型
	ext := strings.ToLower(filepath.Ext(fileName))
	if !imageExtensions[ext] {
		return nil, fmt.Errorf("不支持的图片类型: %s，仅支持 jpg/jpeg/png/gif/webp", ext)
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		return nil, fmt.Errorf("创建缩略图目录失败: %w", err)
	}

	// 使用UUID生成安全的文件名
	newFileName := uuid.New().String() + ext
	filePath := filepath.Clean(filepath.Join(uploadDir, newFileName))
	thumbnailName := uuid.New().String() + "_thumb" + ext
	thumbnailPath := filepath.Clean(filepath.Join(thumbnailDir, thumbnailName))

	// 保存原始文件并计算MD5
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	hash := md5.New()
	writer := io.MultiWriter(dst, hash)

	if _, err := io.Copy(writer, reader); err != nil {
		dst.Close()
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}
	dst.Close()
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	// 检查MD5去重
	existing, err := fileRepo.FindFileByMD5(md5Hash)
	if err == nil && existing != nil {
		os.Remove(filePath)
		logger.Infof("图片已存在(MD5去重): %s -> %s", fileName, existing.Name)
		return &dto.ImageUploadResponse{
			ID:       existing.ID,
			Name:     existing.Name,
			Original: existing.Original,
			URL:      existing.URL,
			Size:     existing.Size,
			Width:    existing.Width,
			Height:   existing.Height,
		}, nil
	}

	// 生成缩略图（简单复制作为占位实现，实际项目中应使用图片处理库）
	width, height := getImageDimensions(filePath)
	if err := copyFile(filePath, thumbnailPath); err != nil {
		logger.Errorf("生成缩略图失败: %v", err)
	}

	// 保存文件信息到数据库
	thumbnailURL := "/uploads/thumbnails/" + thumbnailName
	thumbnailIDsJSON, _ := json.Marshal([]string{thumbnailName})

	file := &model.File{
		Name:         newFileName,
		Original:     fileName,
		Path:         filePath,
		URL:          "/uploads/" + newFileName,
		Size:         fileSize,
		Type:         fileType,
		Extension:    ext,
		MD5:          md5Hash,
		Width:        width,
		Height:       height,
		ThumbnailIDs: string(thumbnailIDsJSON),
		UploaderID:   uploaderID,
		Status:       "active",
	}

	if err := fileRepo.CreateFile(file); err != nil {
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	logger.Infof("图片上传成功: %s (%dx%d)", fileName, width, height)
	return &dto.ImageUploadResponse{
		ID:        file.ID,
		Name:      file.Name,
		Original:  file.Original,
		URL:       file.URL,
		Thumbnail: thumbnailURL,
		Size:      file.Size,
		Width:     width,
		Height:    height,
	}, nil
}

// ==================== Batch Upload ====================

// BatchUpload 批量上传文件
func BatchUpload(uploaderID uint, files []*multipart.FileHeader) *dto.BatchUploadResponse {
	result := &dto.BatchUploadResponse{
		Total: len(files),
	}

	if len(files) > batchMaxCount {
		result.Failed = make([]dto.FileUploadError, len(files))
		for i, f := range files {
			result.Failed[i] = dto.FileUploadError{
				Filename: f.Filename,
				Error:    fmt.Sprintf("批量上传最多支持 %d 个文件", batchMaxCount),
			}
		}
		return result
	}

	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			result.Failed = append(result.Failed, dto.FileUploadError{
				Filename: header.Filename,
				Error:    "打开文件失败: " + err.Error(),
			})
			continue
		}

		f, err := UploadFile(uploaderID, header.Filename, header.Size, header.Header.Get("Content-Type"), file)
		file.Close()

		if err != nil {
			result.Failed = append(result.Failed, dto.FileUploadError{
				Filename: header.Filename,
				Error:    err.Error(),
			})
			continue
		}

		result.Success = append(result.Success, dto.FileUploadResult{
			ID:       f.ID,
			Name:     f.Name,
			Original: f.Original,
			URL:      f.URL,
			Size:     f.Size,
		})
	}

	logger.Infof("批量上传完成: 成功 %d, 失败 %d", len(result.Success), len(result.Failed))
	return result
}

// ==================== Download ====================

// GetFileForDownload 获取文件下载信息
func GetFileForDownload(fileID, userID uint) (*model.File, error) {
	file, err := fileRepo.FindFileByID(fileID)
	if err != nil {
		return nil, errors.New("文件不存在")
	}

	if file.Status != "active" {
		return nil, errors.New("文件不可用")
	}

	// 记录下载日志
	log := &model.FileAccessLog{
		FileID:     fileID,
		UserID:     userID,
		AccessType: "download",
		AccessedAt: time.Now(),
	}
	if err := fileRepo.CreateFileAccessLog(log); err != nil {
		logger.Errorf("记录下载日志失败: %v", err)
	}

	return file, nil
}

// ==================== Thumbnail ====================

// GetThumbnail 获取缩略图路径
func GetThumbnail(fileID uint, size string) (string, string, error) {
	file, err := fileRepo.FindFileByID(fileID)
	if err != nil {
		return "", "", errors.New("文件不存在")
	}

	if file.Status != "active" {
		return "", "", errors.New("文件不可用")
	}

	ext := strings.ToLower(file.Extension)
	if !imageExtensions[ext] {
		return "", "", errors.New("仅图片文件支持缩略图")
	}

	// 解析缩略图ID列表
	var thumbnailIDs []string
	if file.ThumbnailIDs != "" {
		json.Unmarshal([]byte(file.ThumbnailIDs), &thumbnailIDs)
	}

	if len(thumbnailIDs) == 0 {
		// 没有缩略图，返回原图
		return file.Path, file.Type, nil
	}

	// 根据size参数选择缩略图
	thumbName := thumbnailIDs[0]
	thumbPath := filepath.Clean(filepath.Join(thumbnailDir, thumbName))

	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		// 缩略图不存在，返回原图
		return file.Path, file.Type, nil
	}

	return thumbPath, file.Type, nil
}

// ==================== File Reference ====================

// AddReference 添加文件引用
func AddReference(fileID, userID uint, req *dto.AddReferenceRequest) error {
	file, err := fileRepo.FindFileByID(fileID)
	if err != nil {
		return errors.New("文件不存在")
	}

	if file.Status != "active" {
		return errors.New("文件不可用")
	}

	ref := &model.FileReference{
		FileID:       fileID,
		BusinessType: req.BusinessType,
		BusinessID:   req.BusinessID,
		FieldName:    req.FieldName,
		CreatorID:    userID,
	}

	if err := fileRepo.CreateFileReference(ref); err != nil {
		return fmt.Errorf("创建引用失败: %w", err)
	}

	// 更新引用计数
	if err := fileRepo.IncrementReferenceCount(fileID); err != nil {
		logger.Errorf("更新引用计数失败: %v", err)
	}

	logger.Infof("添加文件引用: file=%d, business=%s/%d", fileID, req.BusinessType, req.BusinessID)
	return nil
}

// GetReferences 获取文件引用列表
func GetReferences(fileID uint) ([]dto.FileReferenceInfo, error) {
	_, err := fileRepo.FindFileByID(fileID)
	if err != nil {
		return nil, errors.New("文件不存在")
	}

	refs, err := fileRepo.ListFileReferences(fileID)
	if err != nil {
		return nil, fmt.Errorf("查询引用列表失败: %w", err)
	}

	result := make([]dto.FileReferenceInfo, 0, len(refs))
	for _, ref := range refs {
		result = append(result, dto.FileReferenceInfo{
			ID:           ref.ID,
			FileID:       ref.FileID,
			BusinessType: ref.BusinessType,
			BusinessID:   ref.BusinessID,
			FieldName:    ref.FieldName,
			CreatorID:    ref.CreatorID,
			CreatedAt:    ref.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// DeleteReference 删除文件引用
func DeleteReference(fileID, refID, userID uint) error {
	file, err := fileRepo.FindFileByID(fileID)
	if err != nil {
		return errors.New("文件不存在")
	}

	ref, err := fileRepo.FindFileReferenceByID(refID)
	if err != nil {
		return errors.New("引用不存在")
	}

	if ref.FileID != fileID {
		return errors.New("引用与文件不匹配")
	}

	if err := fileRepo.DeleteFileReferenceByID(refID); err != nil {
		return fmt.Errorf("删除引用失败: %w", err)
	}

	// 更新引用计数
	_ = file

	if err := fileRepo.DecrementReferenceCount(fileID); err != nil {
		logger.Errorf("更新引用计数失败: %v", err)
	}

	logger.Infof("删除文件引用: ref=%d, file=%d", refID, fileID)
	return nil
}

// ==================== Access Log ====================

// GetAccessLogs 获取文件访问日志
func GetAccessLogs(fileID uint, page, pageSize int) ([]dto.FileAccessLogInfo, int64, error) {
	_, err := fileRepo.FindFileByID(fileID)
	if err != nil {
		return nil, 0, errors.New("文件不存在")
	}

	logs, total, err := fileRepo.ListFileAccessLogs(fileID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询访问日志失败: %w", err)
	}

	result := make([]dto.FileAccessLogInfo, 0, len(logs))
	for _, l := range logs {
		result = append(result, dto.FileAccessLogInfo{
			ID:         l.ID,
			FileID:     l.FileID,
			UserID:     l.UserID,
			AccessType: l.AccessType,
			IP:         l.IP,
			UserAgent:  l.UserAgent,
			AccessedAt: l.AccessedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// ==================== Hotlink Protection (Admin) ====================

// ListHotlinkRules 获取防盗链规则列表
func ListHotlinkRules() ([]dto.HotlinkRuleInfo, error) {
	rules, err := fileRepo.ListHotlinkRules()
	if err != nil {
		return nil, fmt.Errorf("查询防盗链规则失败: %w", err)
	}

	result := make([]dto.HotlinkRuleInfo, 0, len(rules))
	for _, r := range rules {
		var allowed, blocked []string
		if r.AllowedDomains != "" {
			json.Unmarshal([]byte(r.AllowedDomains), &allowed)
		}
		if r.BlockedDomains != "" {
			json.Unmarshal([]byte(r.BlockedDomains), &blocked)
		}

		result = append(result, dto.HotlinkRuleInfo{
			ID:             r.ID,
			Name:           r.Name,
			AllowedDomains: allowed,
			BlockedDomains: blocked,
			IsActive:       r.IsActive,
			DefaultAction:  r.DefaultAction,
			CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      r.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// CreateHotlinkRule 创建防盗链规则
func CreateHotlinkRule(req *dto.HotlinkRuleRequest) error {
	allowedJSON, _ := json.Marshal(req.AllowedDomains)
	blockedJSON, _ := json.Marshal(req.BlockedDomains)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	defaultAction := req.DefaultAction
	if defaultAction == "" {
		defaultAction = "allow"
	}

	rule := &model.HotlinkProtectionRule{
		Name:           req.Name,
		AllowedDomains: string(allowedJSON),
		BlockedDomains: string(blockedJSON),
		IsActive:       isActive,
		DefaultAction:  defaultAction,
	}

	if err := fileRepo.CreateHotlinkRule(rule); err != nil {
		return fmt.Errorf("创建防盗链规则失败: %w", err)
	}

	logger.Infof("创建防盗链规则: %s", req.Name)
	return nil
}

// UpdateHotlinkRule 更新防盗链规则
func UpdateHotlinkRule(id uint, req *dto.HotlinkRuleRequest) error {
	rule, err := fileRepo.FindHotlinkRuleByID(id)
	if err != nil {
		return errors.New("防盗链规则不存在")
	}

	allowedJSON, _ := json.Marshal(req.AllowedDomains)
	blockedJSON, _ := json.Marshal(req.BlockedDomains)

	rule.Name = req.Name
	rule.AllowedDomains = string(allowedJSON)
	rule.BlockedDomains = string(blockedJSON)
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	if req.DefaultAction != "" {
		rule.DefaultAction = req.DefaultAction
	}

	if err := fileRepo.UpdateHotlinkRule(rule); err != nil {
		return fmt.Errorf("更新防盗链规则失败: %w", err)
	}

	logger.Infof("更新防盗链规则: %d - %s", id, req.Name)
	return nil
}

// DeleteHotlinkRule 删除防盗链规则
func DeleteHotlinkRule(id uint) error {
	if err := fileRepo.DeleteHotlinkRuleByID(id); err != nil {
		return fmt.Errorf("删除防盗链规则失败: %w", err)
	}

	logger.Infof("删除防盗链规则: %d", id)
	return nil
}

// ==================== Watermark Config (Admin) ====================

// ListWatermarkConfigs 获取水印配置列表
func ListWatermarkConfigs() ([]dto.WatermarkConfigInfo, error) {
	configs, err := fileRepo.ListWatermarkConfigs()
	if err != nil {
		return nil, fmt.Errorf("查询水印配置失败: %w", err)
	}

	result := make([]dto.WatermarkConfigInfo, 0, len(configs))
	for _, c := range configs {
		result = append(result, dto.WatermarkConfigInfo{
			ID:        c.ID,
			Name:      c.Name,
			Type:      c.Type,
			Position:  c.Position,
			Opacity:   c.Opacity,
			Text:      c.Text,
			ImagePath: c.ImagePath,
			FontSize:  c.FontSize,
			FontColor: c.FontColor,
			IsActive:  c.IsActive,
			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// CreateWatermarkConfig 创建水印配置
func CreateWatermarkConfig(req *dto.WatermarkConfigRequest) error {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	position := req.Position
	if position == "" {
		position = "bottom-right"
	}

	opacity := req.Opacity
	if opacity <= 0 || opacity > 1 {
		opacity = 0.5
	}

	fontSize := req.FontSize
	if fontSize <= 0 {
		fontSize = 14
	}

	fontColor := req.FontColor
	if fontColor == "" {
		fontColor = "#000000"
	}

	config := &model.WatermarkConfig{
		Name:      req.Name,
		Type:      req.Type,
		Position:  position,
		Opacity:   opacity,
		Text:      req.Text,
		ImagePath: req.ImagePath,
		FontSize:  fontSize,
		FontColor: fontColor,
		IsActive:  isActive,
	}

	if err := fileRepo.CreateWatermarkConfig(config); err != nil {
		return fmt.Errorf("创建水印配置失败: %w", err)
	}

	logger.Infof("创建水印配置: %s", req.Name)
	return nil
}

// UpdateWatermarkConfig 更新水印配置
func UpdateWatermarkConfig(id uint, req *dto.WatermarkConfigRequest) error {
	config, err := fileRepo.FindWatermarkConfigByID(id)
	if err != nil {
		return errors.New("水印配置不存在")
	}

	config.Name = req.Name
	config.Type = req.Type
	if req.Position != "" {
		config.Position = req.Position
	}
	if req.Opacity > 0 && req.Opacity <= 1 {
		config.Opacity = req.Opacity
	}
	config.Text = req.Text
	config.ImagePath = req.ImagePath
	if req.FontSize > 0 {
		config.FontSize = req.FontSize
	}
	if req.FontColor != "" {
		config.FontColor = req.FontColor
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}

	if err := fileRepo.UpdateWatermarkConfig(config); err != nil {
		return fmt.Errorf("更新水印配置失败: %w", err)
	}

	logger.Infof("更新水印配置: %d - %s", id, req.Name)
	return nil
}

// ==================== Orphan Cleanup (Admin) ====================

// CleanupOrphanFiles 清理孤立文件
func CleanupOrphanFiles() (*dto.CleanupOrphanResponse, error) {
	// 创建清理日志
	now := time.Now()
	log := &model.FileCleanupLog{
		Status:    "running",
		StartedAt: &now,
	}
	if err := fileRepo.CreateCleanupLog(log); err != nil {
		return nil, fmt.Errorf("创建清理日志失败: %w", err)
	}

	// 查找孤立文件（超过7天无引用的文件）
	orphanFiles, err := fileRepo.FindOrphanFiles(7)
	if err != nil {
		log.Status = "failed"
		fileRepo.UpdateCleanupLog(log)
		return nil, fmt.Errorf("查找孤立文件失败: %w", err)
	}

	var cleanedCount int
	var freedSpace int64
	var cleanedIDs []uint

	for _, f := range orphanFiles {
		// 检查文件是否确实没有引用
		refCount, err := fileRepo.CountFileReferences(f.ID)
		if err != nil {
			logger.Errorf("检查文件引用失败: %d, %v", f.ID, err)
			continue
		}
		if refCount > 0 {
			continue
		}

		// 删除物理文件
		if err := os.Remove(f.Path); err != nil && !os.IsNotExist(err) {
			logger.Errorf("删除物理文件失败: %s, %v", f.Path, err)
			continue
		}

		// 删除缩略图
		var thumbIDs []string
		if f.ThumbnailIDs != "" {
			json.Unmarshal([]byte(f.ThumbnailIDs), &thumbIDs)
		}
		for _, thumbID := range thumbIDs {
			thumbPath := filepath.Clean(filepath.Join(thumbnailDir, thumbID))
			os.Remove(thumbPath)
		}

		// 标记为已删除
		f.Status = "deleted"
		if err := fileRepo.UpdateFile(&f); err != nil {
			logger.Errorf("更新文件状态失败: %d, %v", f.ID, err)
			continue
		}

		cleanedCount++
		freedSpace += f.Size
		cleanedIDs = append(cleanedIDs, f.ID)
	}

	// 更新清理日志
	completedAt := time.Now()
	idsJSON, _ := json.Marshal(cleanedIDs)
	log.CleanedCount = cleanedCount
	log.FreedSpace = freedSpace
	log.FileIDs = string(idsJSON)
	log.Status = "completed"
	log.CompletedAt = &completedAt
	if err := fileRepo.UpdateCleanupLog(log); err != nil {
		logger.Errorf("更新清理日志失败: %v", err)
	}

	logger.Infof("孤立文件清理完成: 清理 %d 个文件，释放 %d 字节", cleanedCount, freedSpace)

	return &dto.CleanupOrphanResponse{
		CleanedCount: cleanedCount,
		FreedSpace:   freedSpace,
		LogID:        log.ID,
	}, nil
}

// GetCleanupLogs 获取清理日志列表
func GetCleanupLogs(page, pageSize int) ([]dto.CleanupLogInfo, int64, error) {
	logs, total, err := fileRepo.ListCleanupLogs(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询清理日志失败: %w", err)
	}

	result := make([]dto.CleanupLogInfo, 0, len(logs))
	for _, l := range logs {
		info := dto.CleanupLogInfo{
			ID:           l.ID,
			CleanedCount: l.CleanedCount,
			FreedSpace:   l.FreedSpace,
			FileIDs:      l.FileIDs,
			Status:       l.Status,
			CreatedAt:    l.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if l.StartedAt != nil {
			info.StartedAt = l.StartedAt.Format("2006-01-02 15:04:05")
		}
		if l.CompletedAt != nil {
			info.CompletedAt = l.CompletedAt.Format("2006-01-02 15:04:05")
		}
		result = append(result, info)
	}

	return result, total, nil
}

// ==================== Storage Statistics (Admin) ====================

// GetStorageStatistics 获取存储统计
func GetStorageStatistics() (*dto.StorageStatistics, error) {
	stats := &dto.StorageStatistics{}

	// 总文件数和大小
	totalCount, totalSize, err := fileRepo.CountTotalFiles()
	if err != nil {
		return nil, fmt.Errorf("统计总文件失败: %w", err)
	}
	stats.TotalFiles = totalCount
	stats.TotalSize = totalSize

	// 图片统计
	imageCount, imageSize, err := fileRepo.CountFilesByType([]string{".jpg", ".jpeg", ".png", ".gif", ".webp"})
	if err != nil {
		return nil, fmt.Errorf("统计图片失败: %w", err)
	}
	stats.ImageCount = imageCount
	stats.ImageSize = imageSize

	// 文档统计
	docCount, docSize, err := fileRepo.CountFilesByType([]string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".csv"})
	if err != nil {
		return nil, fmt.Errorf("统计文档失败: %w", err)
	}
	stats.DocumentCount = docCount
	stats.DocumentSize = docSize

	// 视频统计
	videoCount, videoSize, err := fileRepo.CountFilesByType([]string{".mp4"})
	if err != nil {
		return nil, fmt.Errorf("统计视频失败: %w", err)
	}
	stats.VideoCount = videoCount
	stats.VideoSize = videoSize

	// 音频统计
	audioCount, audioSize, err := fileRepo.CountFilesByType([]string{".mp3"})
	if err != nil {
		return nil, fmt.Errorf("统计音频失败: %w", err)
	}
	stats.AudioCount = audioCount
	stats.AudioSize = audioSize

	// 其他类型统计
	knownExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".csv", ".mp4", ".mp3"}
	otherCount := totalCount - imageCount - docCount - videoCount - audioCount
	if otherCount < 0 {
		otherCount = 0
	}
	otherSize := totalSize - imageSize - docSize - videoSize - audioSize
	if otherSize < 0 {
		otherSize = 0
	}
	_ = knownExts
	stats.OtherCount = otherCount
	stats.OtherSize = otherSize

	// 孤立文件统计
	orphanCount, orphanSize, err := fileRepo.CountOrphanFiles()
	if err != nil {
		return nil, fmt.Errorf("统计孤立文件失败: %w", err)
	}
	stats.OrphanCount = orphanCount
	stats.OrphanSize = orphanSize

	// 每日上传统计（最近30天）
	dailyStats, err := fileRepo.GetDailyUploadStats(30)
	if err != nil {
		return nil, fmt.Errorf("统计每日上传失败: %w", err)
	}
	stats.DailyUploads = make([]dto.DailyUploadStat, 0, len(dailyStats))
	for _, d := range dailyStats {
		stats.DailyUploads = append(stats.DailyUploads, dto.DailyUploadStat{
			Date:  d.Date,
			Count: d.Count,
			Size:  d.Size,
		})
	}

	// 上传者排行（前10）
	topUploaders, err := fileRepo.GetTopUploaders(10)
	if err != nil {
		return nil, fmt.Errorf("统计上传者失败: %w", err)
	}
	stats.TopUploaders = make([]dto.UploaderStat, 0, len(topUploaders))
	for _, u := range topUploaders {
		stats.TopUploaders = append(stats.TopUploaders, dto.UploaderStat{
			UserID:    u.UserID,
			Count:     u.Count,
			TotalSize: u.TotalSize,
		})
	}

	return stats, nil
}

// ==================== Delete ====================

// DeleteFile 删除文件
func DeleteFile(id, userID uint) error {
	file, err := fileRepo.FindFileByID(id)
	if err != nil {
		return errors.New("文件不存在")
	}

	// 检查权限（只能删除自己上传的文件）
	if file.UploaderID != userID {
		return errors.New("无权删除此文件")
	}

	// 检查是否有引用
	refCount, _ := fileRepo.CountFileReferences(id)
	if refCount > 0 {
		return fmt.Errorf("文件正在被 %d 处引用，无法删除", refCount)
	}

	// 删除物理文件
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		logger.Errorf("删除物理文件失败: %v", err)
	}

	// 删除缩略图
	var thumbIDs []string
	if file.ThumbnailIDs != "" {
		json.Unmarshal([]byte(file.ThumbnailIDs), &thumbIDs)
	}
	for _, thumbID := range thumbIDs {
		thumbPath := filepath.Clean(filepath.Join(thumbnailDir, thumbID))
		os.Remove(thumbPath)
	}

	// 删除数据库记录
	if err := fileRepo.DeleteFileByID(id); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	logger.Infof("文件删除成功: %d", id)
	return nil
}

// GetFile 获取文件信息
func GetFile(id uint) (*model.File, error) {
	return fileRepo.FindFileByID(id)
}

// ListFiles 文件列表
func ListFiles(uploaderID *uint, page, pageSize int) ([]model.File, int64, error) {
	return fileRepo.ListFiles(uploaderID, page, pageSize)
}

// ==================== Helper Functions ====================

// getImageDimensions 获取图片尺寸（简单实现，返回0表示无法获取）
func getImageDimensions(filePath string) (int, int) {
	// 打开文件读取图片头信息
	f, err := os.Open(filePath)
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	// 读取前512字节判断格式
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil || n < 4 {
		return 0, 0
	}

	// 简单的图片尺寸检测
	// PNG: 宽高在第16-23字节
	if n >= 24 && buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
		width := int(buf[16])<<24 | int(buf[17])<<16 | int(buf[18])<<8 | int(buf[19])
		height := int(buf[20])<<24 | int(buf[21])<<16 | int(buf[22])<<8 | int(buf[23])
		return width, height
	}

	// JPEG: 需要解析SOI标记，简化处理
	if buf[0] == 0xFF && buf[1] == 0xD8 {
		// JPEG解析较复杂，返回0表示需要专业库处理
		return 0, 0
	}

	// GIF: 宽高在第6-9字节
	if n >= 10 && buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46 {
		width := int(buf[6]) | int(buf[7])<<8
		height := int(buf[8]) | int(buf[9])<<8
		return width, height
	}

	// WebP: 需要解析RIFF头
	if n >= 30 && buf[0] == 0x52 && buf[1] == 0x49 && buf[2] == 0x46 && buf[3] == 0x46 {
		if buf[8] == 0x57 && buf[9] == 0x45 && buf[10] == 0x42 && buf[11] == 0x50 {
			if buf[12] == 0x56 && buf[13] == 0x50 && buf[14] == 0x38 {
				width := int(buf[26]) | int(buf[27])<<8 + 1
				height := int(buf[28]) | int(buf[29])<<8 + 1
				return width, height
			}
		}
	}

	return 0, 0
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// LogFileAccess 记录文件访问日志
func LogFileAccess(fileID, userID uint, accessType, ip, userAgent string) {
	log := &model.FileAccessLog{
		FileID:     fileID,
		UserID:     userID,
		AccessType: accessType,
		IP:         ip,
		UserAgent:  userAgent,
		AccessedAt: time.Now(),
	}
	if err := fileRepo.CreateFileAccessLog(log); err != nil {
		logger.Errorf("记录文件访问日志失败: %v", err)
	}
}

// CheckHotlink 检查防盗链
func CheckHotlink(referer string) bool {
	rules, err := fileRepo.ListHotlinkRules()
	if err != nil {
		return true // 无规则时允许访问
	}

	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		var allowed, blocked []string
		if rule.AllowedDomains != "" {
			json.Unmarshal([]byte(rule.AllowedDomains), &allowed)
		}
		if rule.BlockedDomains != "" {
			json.Unmarshal([]byte(rule.BlockedDomains), &blocked)
		}

		// 检查禁止列表
		for _, domain := range blocked {
			if strings.Contains(referer, domain) {
				return false
			}
		}

		// 检查允许列表
		if len(allowed) > 0 {
			isAllowed := false
			for _, domain := range allowed {
				if strings.Contains(referer, domain) {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				return rule.DefaultAction == "block"
			}
		}
	}

	return true
}
