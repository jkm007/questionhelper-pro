package question

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== FavoriteFolder ====================

// CreateFolder 创建收藏夹
func CreateFolder(folder *model.FavoriteFolder) error {
	return database.DB.Create(folder).Error
}

// UpdateFolder 更新收藏夹
func UpdateFolder(folder *model.FavoriteFolder) error {
	return database.DB.Save(folder).Error
}

// DeleteFolder 删除收藏夹
func DeleteFolder(id uint) error {
	return database.DB.Delete(&model.FavoriteFolder{}, id).Error
}

// FindFolderByID 获取收藏夹
func FindFolderByID(id uint) (*model.FavoriteFolder, error) {
	var folder model.FavoriteFolder
	err := database.DB.First(&folder, id).Error
	return &folder, err
}

// FindFoldersByUserID 获取用户收藏夹列表
func FindFoldersByUserID(userID uint) ([]model.FavoriteFolder, error) {
	var folders []model.FavoriteFolder
	err := database.DB.Where("user_id = ?", userID).
		Order("sort ASC, id ASC").
		Find(&folders).Error
	return folders, err
}

// FindDefaultFolder 获取默认收藏夹
func FindDefaultFolder(userID uint) (*model.FavoriteFolder, error) {
	var folder model.FavoriteFolder
	err := database.DB.Where("user_id = ? AND is_default = ?", userID, true).
		First(&folder).Error
	return &folder, err
}

// ==================== QuestionFavorite ====================

// AddFavorite 添加收藏
func AddFavorite(fav *model.QuestionFavorite) error {
	return database.DB.Create(fav).Error
}

// RemoveFavorite 取消收藏
func RemoveFavorite(userID, questionID uint) error {
	return database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).
		Delete(&model.QuestionFavorite{}).Error
}

// FindFavorite 查找收藏
func FindFavorite(userID, questionID uint) (*model.QuestionFavorite, error) {
	var fav model.QuestionFavorite
	err := database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).
		First(&fav).Error
	return &fav, err
}

// ListFavorites 收藏列表
func ListFavorites(userID uint, req *dto.FavoriteQueryRequest) ([]model.QuestionFavorite, int64, error) {
	var favorites []model.QuestionFavorite
	var total int64

	db := database.DB.Model(&model.QuestionFavorite{}).Where("user_id = ?", userID)
	if req.FolderID != nil {
		db = db.Where("folder_id = ?", *req.FolderID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("id DESC").
		Find(&favorites).Error

	return favorites, total, err
}

// IsFavorited 是否已收藏
func IsFavorited(userID, questionID uint) bool {
	var count int64
	database.DB.Model(&model.QuestionFavorite{}).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		Count(&count)
	return count > 0
}

// UpdateFolderCount 更新收藏夹数量
func UpdateFolderCount(folderID uint) error {
	var count int64
	database.DB.Model(&model.QuestionFavorite{}).Where("folder_id = ?", folderID).Count(&count)
	return database.DB.Model(&model.FavoriteFolder{}).Where("id = ?", folderID).
		Update("count", count).Error
}

// MoveFavoritesToFolder 将某收藏夹中的收藏移到目标收藏夹
func MoveFavoritesToFolder(fromFolderID, toFolderID uint) error {
	return database.DB.Model(&model.QuestionFavorite{}).
		Where("folder_id = ?", fromFolderID).
		Update("folder_id", toFolderID).Error
}

// ==================== QuestionLike ====================

// AddLike 添加点赞记录
func AddLike(like *model.QuestionLike) error {
	return database.DB.Create(like).Error
}

// RemoveLike 移除点赞记录
func RemoveLike(userID, questionID uint) error {
	return database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).
		Delete(&model.QuestionLike{}).Error
}

// IsLiked 是否已点赞
func IsLiked(userID, questionID uint) bool {
	var count int64
	database.DB.Model(&model.QuestionLike{}).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		Count(&count)
	return count > 0
}
