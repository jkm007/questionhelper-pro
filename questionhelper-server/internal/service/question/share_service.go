package question

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// GenerateShareCode 生成12位分享码
func GenerateShareCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// CreateShare 创建分享
func CreateShare(userID uint, req *dto.CreateShareRequest) (*dto.ShareInfo, error) {
	// 验证题目存在
	_, err := questionRepo.FindByID(req.QuestionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	shareCode := GenerateShareCode()
	share := &model.QuestionShare{
		ShareCode:  shareCode,
		QuestionID: req.QuestionID,
		UserID:     userID,
		ShareType:  req.ShareType,
		Password:   req.Password,
		Status:     1,
	}

	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
		share.ExpiresAt = &expiresAt
	}

	if share.ShareType == 0 {
		share.ShareType = 1
	}

	if err := questionRepo.CreateShare(share); err != nil {
		return nil, fmt.Errorf("创建分享失败: %w", err)
	}

	logger.Infof("用户 %d 创建题目分享成功: %s", userID, shareCode)
	return toShareInfo(share), nil
}

// GetShare 获取分享
func GetShare(code string) (*dto.ShareInfo, error) {
	share, err := questionRepo.FindShareByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分享不存在")
		}
		return nil, fmt.Errorf("查询分享失败: %w", err)
	}

	// 检查是否过期
	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("分享已过期")
	}

	// 检查状态
	if share.Status == 0 {
		return nil, errors.New("分享已撤销")
	}

	// 增加查看次数
	questionRepo.IncrementShareViewCount(share.ID)

	return toShareInfo(share), nil
}

// RevokeShare 撤销分享
func RevokeShare(id, userID uint) error {
	share, err := questionRepo.FindShareByID(id)
	if err != nil {
		return fmt.Errorf("查询分享失败: %w", err)
	}

	if share.UserID != userID {
		return errors.New("无权撤销此分享")
	}

	share.Status = 0
	if err := questionRepo.UpdateShare(share); err != nil {
		return fmt.Errorf("撤销分享失败: %w", err)
	}

	logger.Infof("用户 %d 撤销分享 %d 成功", userID, id)
	return nil
}

// ListSharesByQuestion 获取题目分享列表
func ListSharesByQuestion(questionID uint) ([]dto.ShareInfo, error) {
	shares, err := questionRepo.ListSharesByQuestionID(questionID)
	if err != nil {
		return nil, fmt.Errorf("查询分享列表失败: %w", err)
	}

	list := make([]dto.ShareInfo, 0, len(shares))
	for _, s := range shares {
		list = append(list, *toShareInfo(&s))
	}

	return list, nil
}

// ListMyShares 我的分享列表
func ListMyShares(userID uint, page, pageSize int) ([]dto.ShareInfo, int64, error) {
	shares, total, err := questionRepo.ListSharesByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询分享列表失败: %w", err)
	}

	list := make([]dto.ShareInfo, 0, len(shares))
	for _, s := range shares {
		list = append(list, *toShareInfo(&s))
	}

	return list, total, nil
}

// toShareInfo 转换为 ShareInfo DTO
func toShareInfo(share *model.QuestionShare) *dto.ShareInfo {
	return &dto.ShareInfo{
		ID:          share.ID,
		ShareCode:   share.ShareCode,
		ShareURL:    fmt.Sprintf("/share/%s", share.ShareCode),
		QuestionID:  share.QuestionID,
		ShareType:   share.ShareType,
		HasPassword: share.Password != "",
		ViewCount:   share.ViewCount,
		Status:      share.Status,
		ExpiresAt:   share.ExpiresAt,
		CreatedAt:   share.CreatedAt,
	}
}
