package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

const uploadDir = "./uploads"

// UploadFile 上传文件
func UploadFile(uploaderID uint, fileName string, fileSize int64, fileType string, reader io.Reader) (*model.File, error) {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(fileName)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, newFileName)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 保存文件信息到数据库
	file := &model.File{
		Name:       newFileName,
		Original:   fileName,
		Path:       filePath,
		Size:       fileSize,
		Type:       fileType,
		Extension:  ext,
		UploaderID: uploaderID,
	}

	if err := database.DB.Create(file).Error; err != nil {
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	logger.Infof("文件上传成功: %s", fileName)
	return file, nil
}

// DeleteFile 删除文件
func DeleteFile(id, userID uint) error {
	var file model.File
	if err := database.DB.First(&file, id).Error; err != nil {
		return errors.New("文件不存在")
	}

	// 检查权限（只能删除自己上传的文件）
	if file.UploaderID != userID {
		return errors.New("无权删除此文件")
	}

	// 删除物理文件
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		logger.Errorf("删除物理文件失败: %v", err)
	}

	// 删除数据库记录
	if err := database.DB.Delete(&file).Error; err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	logger.Infof("文件删除成功: %d", id)
	return nil
}

// GetFile 获取文件信息
func GetFile(id uint) (*model.File, error) {
	var file model.File
	if err := database.DB.First(&file, id).Error; err != nil {
		return nil, errors.New("文件不存在")
	}
	return &file, nil
}

// ListFiles 文件列表
func ListFiles(uploaderID *uint, page, pageSize int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	db := database.DB.Model(&model.File{})
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
