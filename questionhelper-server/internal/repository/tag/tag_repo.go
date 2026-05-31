package tag

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Tag ====================

func FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := database.DB.First(&tag, id).Error
	return &tag, err
}

func FindByCode(code string) (*model.Tag, error) {
	var tag model.Tag
	err := database.DB.Where("code = ?", code).First(&tag).Error
	return &tag, err
}

func FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := database.DB.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func Create(tag *model.Tag) error {
	return database.DB.Create(tag).Error
}

func Update(tag *model.Tag) error {
	return database.DB.Save(tag).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.Tag{}, id).Error
}

func List(req *dto.TagListRequest) ([]model.Tag, int64, error) {
	var tags []model.Tag
	var total int64

	db := database.DB.Model(&model.Tag{})

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("sort ASC, id ASC").
		Find(&tags).Error

	return tags, total, err
}

func FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Where("status = 1").Order("sort ASC").Find(&tags).Error
	return tags, err
}

func ExistsByName(name string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Tag{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func ExistsByCode(code string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Tag{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// ==================== UserTag ====================

// AddUserTag 给用户添加标签
func AddUserTag(userID, tagID uint) error {
	ut := &model.UserTag{
		UserID: userID,
		TagID:  tagID,
	}
	return database.DB.Create(ut).Error
}

// RemoveUserTag 移除用户标签
func RemoveUserTag(userID, tagID uint) error {
	return database.DB.Where("user_id = ? AND tag_id = ?", userID, tagID).
		Delete(&model.UserTag{}).Error
}

// FindTagsByUserID 获取用户的所有标签
func FindTagsByUserID(userID uint) ([]model.Tag, error) {
	var user model.User
	if err := database.DB.Preload("Tags").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Tags, nil
}

// FindUsersByTagID 获取标签下的所有用户ID
func FindUsersByTagID(tagID uint) ([]uint, error) {
	var userTags []model.UserTag
	if err := database.DB.Where("tag_id = ?", tagID).Find(&userTags).Error; err != nil {
		return nil, err
	}
	userIDs := make([]uint, 0, len(userTags))
	for _, ut := range userTags {
		userIDs = append(userIDs, ut.UserID)
	}
	return userIDs, nil
}

// BatchAddUserTags 批量给用户添加标签
func BatchAddUserTags(userID uint, tagIDs []uint) error {
	userTags := make([]model.UserTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		userTags = append(userTags, model.UserTag{
			UserID: userID,
			TagID:  tagID,
		})
	}
	return database.DB.CreateInBatches(userTags, 100).Error
}

// BatchRemoveUserTags 批量移除用户标签
func BatchRemoveUserTags(userID uint, tagIDs []uint) error {
	return database.DB.Where("user_id = ? AND tag_id IN ?", userID, tagIDs).
		Delete(&model.UserTag{}).Error
}

// UpdateTagUserCount 更新标签的用户数量
func UpdateTagUserCount(tagID uint) error {
	var count int64
	database.DB.Model(&model.UserTag{}).Where("tag_id = ?", tagID).Count(&count)
	return database.DB.Model(&model.Tag{}).Where("id = ?", tagID).
		Update("user_count", count).Error
}
