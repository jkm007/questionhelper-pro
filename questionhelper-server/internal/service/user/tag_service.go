package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	tagRepo "questionhelper-server/internal/repository/tag"
	"questionhelper-server/pkg/logger"
)

// CreateTag 创建标签
func CreateTag(req *dto.CreateTagRequest) error {
	// 检查名称是否重复
	exists, err := tagRepo.ExistsByName(req.Name)
	if err != nil {
		return fmt.Errorf("检查标签名称失败: %w", err)
	}
	if exists {
		return errors.New("标签名称已存在")
	}

	// 检查编码是否重复
	exists, err = tagRepo.ExistsByCode(req.Code)
	if err != nil {
		return fmt.Errorf("检查标签编码失败: %w", err)
	}
	if exists {
		return errors.New("标签编码已存在")
	}

	tag := &model.Tag{
		Name:  req.Name,
		Code:  req.Code,
		Color: req.Color,
		Icon:  req.Icon,
		Type:  req.Type,
		Sort:  req.Sort,
	}

	if tag.Type == 0 {
		tag.Type = 2 // 默认自定义标签
	}

	if err := tagRepo.Create(tag); err != nil {
		return fmt.Errorf("创建标签失败: %w", err)
	}

	logger.Infof("创建标签成功: %s", req.Name)
	return nil
}

// UpdateTag 更新标签
func UpdateTag(id uint, req *dto.UpdateTagRequest) error {
	tag, err := tagRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	// 系统标签不允许修改名称
	if tag.Type == 1 && req.Name != "" {
		return errors.New("系统标签不允许修改名称")
	}

	if req.Name != "" {
		if req.Name != tag.Name {
			exists, err := tagRepo.ExistsByName(req.Name)
			if err != nil {
				return fmt.Errorf("检查标签名称失败: %w", err)
			}
			if exists {
				return errors.New("标签名称已存在")
			}
		}
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}
	if req.Icon != "" {
		tag.Icon = req.Icon
	}
	tag.Sort = req.Sort
	if req.Status != nil {
		tag.Status = *req.Status
	}

	if err := tagRepo.Update(tag); err != nil {
		return fmt.Errorf("更新标签失败: %w", err)
	}

	logger.Infof("更新标签成功: %d", id)
	return nil
}

// DeleteTag 删除标签
func DeleteTag(id uint) error {
	tag, err := tagRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	// 系统标签不允许删除
	if tag.Type == 1 {
		return errors.New("系统标签不允许删除")
	}

	if err := tagRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}

	logger.Infof("删除标签成功: %d", id)
	return nil
}

// GetTag 获取标签详情
func GetTag(id uint) (*dto.TagInfo, error) {
	tag, err := tagRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, fmt.Errorf("查询标签失败: %w", err)
	}
	return toTagInfo(tag), nil
}

// ListTags 标签列表
func ListTags(req *dto.TagListRequest) ([]dto.TagInfo, int64, error) {
	tags, total, err := tagRepo.List(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询标签列表失败: %w", err)
	}

	list := make([]dto.TagInfo, 0, len(tags))
	for _, tag := range tags {
		list = append(list, *toTagInfo(&tag))
	}

	return list, total, nil
}

// GetAllTags 获取所有启用的标签
func GetAllTags() ([]dto.TagInfo, error) {
	tags, err := tagRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}

	list := make([]dto.TagInfo, 0, len(tags))
	for _, tag := range tags {
		list = append(list, *toTagInfo(&tag))
	}

	return list, nil
}

// AddUserTags 给用户添加标签
func AddUserTags(userID uint, tagIDs []uint) error {
	// 验证用户存在
	_, err := FindUserByID(userID)
	if err != nil {
		return err
	}

	// 验证标签存在
	for _, tagID := range tagIDs {
		_, err := tagRepo.FindByID(tagID)
		if err != nil {
			return fmt.Errorf("标签 %d 不存在", tagID)
		}
	}

	if err := tagRepo.BatchAddUserTags(userID, tagIDs); err != nil {
		return fmt.Errorf("添加用户标签失败: %w", err)
	}

	// 更新标签用户数
	for _, tagID := range tagIDs {
		tagRepo.UpdateTagUserCount(tagID)
	}

	logger.Infof("用户 %d 添加标签成功: %v", userID, tagIDs)
	return nil
}

// RemoveUserTags 移除用户标签
func RemoveUserTags(userID uint, tagIDs []uint) error {
	if err := tagRepo.BatchRemoveUserTags(userID, tagIDs); err != nil {
		return fmt.Errorf("移除用户标签失败: %w", err)
	}

	// 更新标签用户数
	for _, tagID := range tagIDs {
		tagRepo.UpdateTagUserCount(tagID)
	}

	logger.Infof("用户 %d 移除标签成功: %v", userID, tagIDs)
	return nil
}

// GetUserTags 获取用户的标签列表
func GetUserTags(userID uint) ([]dto.TagInfo, error) {
	tags, err := tagRepo.FindTagsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户标签失败: %w", err)
	}

	list := make([]dto.TagInfo, 0, len(tags))
	for _, tag := range tags {
		list = append(list, *toTagInfo(&tag))
	}

	return list, nil
}

// toTagInfo 转换为 TagInfo DTO
func toTagInfo(tag *model.Tag) *dto.TagInfo {
	return &dto.TagInfo{
		ID:        tag.ID,
		Name:      tag.Name,
		Code:      tag.Code,
		Color:     tag.Color,
		Icon:      tag.Icon,
		Type:      tag.Type,
		Sort:      tag.Sort,
		Status:    tag.Status,
		UserCount: tag.UserCount,
	}
}
