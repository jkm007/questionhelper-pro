package question

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// CreateFolder 创建收藏夹
func CreateFolder(userID uint, req *dto.CreateFolderRequest) error {
	folder := &model.FavoriteFolder{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}

	if err := questionRepo.CreateFolder(folder); err != nil {
		return fmt.Errorf("创建收藏夹失败: %w", err)
	}

	logger.Infof("用户 %d 创建收藏夹成功: %s", userID, req.Name)
	return nil
}

// UpdateFolder 更新收藏夹
func UpdateFolder(id, userID uint, req *dto.UpdateFolderRequest) error {
	folder, err := questionRepo.FindFolderByID(id)
	if err != nil {
		return fmt.Errorf("查询收藏夹失败: %w", err)
	}

	if folder.UserID != userID {
		return errors.New("无权修改此收藏夹")
	}

	if req.Name != "" {
		folder.Name = req.Name
	}
	if req.Description != "" {
		folder.Description = req.Description
	}
	if req.Icon != "" {
		folder.Icon = req.Icon
	}
	folder.Sort = req.Sort

	if err := questionRepo.UpdateFolder(folder); err != nil {
		return fmt.Errorf("更新收藏夹失败: %w", err)
	}

	return nil
}

// DeleteFolder 删除收藏夹
func DeleteFolder(id, userID uint) error {
	folder, err := questionRepo.FindFolderByID(id)
	if err != nil {
		return fmt.Errorf("查询收藏夹失败: %w", err)
	}

	if folder.UserID != userID {
		return errors.New("无权删除此收藏夹")
	}

	if folder.IsDefault {
		return errors.New("默认收藏夹不能删除")
	}

	// 将该收藏夹中的收藏移到默认收藏夹
	defaultFolder, err := questionRepo.FindDefaultFolder(userID)
	if err != nil {
		// 如果没有默认收藏夹，创建一个
		defaultFolder = &model.FavoriteFolder{
			UserID:    userID,
			Name:      "默认收藏夹",
			IsDefault: true,
		}
		if err := questionRepo.CreateFolder(defaultFolder); err != nil {
			return fmt.Errorf("创建默认收藏夹失败: %w", err)
		}
	}

	if err := questionRepo.MoveFavoritesToFolder(id, defaultFolder.ID); err != nil {
		return fmt.Errorf("移动收藏到默认收藏夹失败: %w", err)
	}

	// 更新默认收藏夹数量
	questionRepo.UpdateFolderCount(defaultFolder.ID)

	if err := questionRepo.DeleteFolder(id); err != nil {
		return fmt.Errorf("删除收藏夹失败: %w", err)
	}

	logger.Infof("用户 %d 删除收藏夹 %d 成功，收藏已移至默认收藏夹", userID, id)
	return nil
}

// ListFolders 获取收藏夹列表
func ListFolders(userID uint) ([]dto.FolderInfo, error) {
	folders, err := questionRepo.FindFoldersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询收藏夹列表失败: %w", err)
	}

	list := make([]dto.FolderInfo, 0, len(folders))
	for _, f := range folders {
		list = append(list, dto.FolderInfo{
			ID:          f.ID,
			Name:        f.Name,
			Description: f.Description,
			Icon:        f.Icon,
			Sort:        f.Sort,
			Count:       f.Count,
			IsDefault:   f.IsDefault,
			CreatedAt:   f.CreatedAt,
		})
	}

	return list, nil
}

// FavoriteQuestion 收藏题目
func FavoriteQuestion(userID, questionID uint, req *dto.FavoriteRequest) error {
	// 检查是否已收藏
	existing, err := questionRepo.FindFavorite(userID, questionID)
	if err == nil && existing.ID > 0 {
		return errors.New("已收藏此题目")
	}

	folderID := req.FolderID
	if folderID == 0 {
		// 获取默认收藏夹
		folder, err := questionRepo.FindDefaultFolder(userID)
		if err != nil {
			// 创建默认收藏夹
			folder = &model.FavoriteFolder{
				UserID:    userID,
				Name:      "默认收藏夹",
				IsDefault: true,
			}
			questionRepo.CreateFolder(folder)
		}
		folderID = folder.ID
	}

	fav := &model.QuestionFavorite{
		UserID:     userID,
		QuestionID: questionID,
		FolderID:   folderID,
		Note:       req.Note,
	}

	if err := questionRepo.AddFavorite(fav); err != nil {
		return fmt.Errorf("收藏失败: %w", err)
	}

	// 更新题目收藏数（原子操作）
	if err := questionRepo.IncrementFavoriteCount(questionID); err != nil {
		logger.Errorf("更新题目收藏数失败: %v", err)
	}

	// 更新收藏夹数量
	questionRepo.UpdateFolderCount(folderID)

	logger.Infof("用户 %d 收藏题目 %d 成功", userID, questionID)
	return nil
}

// UnfavoriteQuestion 取消收藏
func UnfavoriteQuestion(userID, questionID uint) error {
	if err := questionRepo.RemoveFavorite(userID, questionID); err != nil {
		return fmt.Errorf("取消收藏失败: %w", err)
	}

	// 更新题目收藏数（原子操作）
	if err := questionRepo.DecrementFavoriteCount(questionID); err != nil {
		logger.Errorf("更新题目收藏数失败: %v", err)
	}

	logger.Infof("用户 %d 取消收藏题目 %d 成功", userID, questionID)
	return nil
}

// ListFavorites 收藏列表
func ListFavorites(userID uint, req *dto.FavoriteQueryRequest) ([]dto.FavoriteInfo, int64, error) {
	favorites, total, err := questionRepo.ListFavorites(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询收藏列表失败: %w", err)
	}

	list := make([]dto.FavoriteInfo, 0, len(favorites))
	for _, fav := range favorites {
		info := dto.FavoriteInfo{
			ID:         fav.ID,
			QuestionID: fav.QuestionID,
			FolderID:   fav.FolderID,
			Note:       fav.Note,
			CreatedAt:  fav.CreatedAt,
		}

		// 获取题目信息
		question, err := questionRepo.FindByID(fav.QuestionID)
		if err == nil {
			qInfo := toQuestionInfo(question)
			info.Question = &qInfo
		}

		list = append(list, info)
	}

	return list, total, nil
}

// IsFavorited 是否已收藏
func IsFavorited(userID, questionID uint) bool {
	return questionRepo.IsFavorited(userID, questionID)
}
