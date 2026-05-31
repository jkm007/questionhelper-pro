package application

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// FindByID 根据ID查找角色申请
func FindByID(id uint) (*model.RoleApplication, error) {
	var app model.RoleApplication
	err := database.DB.Preload("User").Preload("Role").Preload("Reviewer").
		First(&app, id).Error
	return &app, err
}

// Create 创建角色申请
func Create(app *model.RoleApplication) error {
	return database.DB.Create(app).Error
}

// Update 更新角色申请
func Update(app *model.RoleApplication) error {
	return database.DB.Save(app).Error
}

// DeleteByID 删除角色申请
func DeleteByID(id uint) error {
	return database.DB.Delete(&model.RoleApplication{}, id).Error
}

// List 角色申请列表
func List(req *dto.ApplicationListRequest) ([]model.RoleApplication, int64, error) {
	var apps []model.RoleApplication
	var total int64

	db := database.DB.Model(&model.RoleApplication{})

	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.RoleID != nil {
		db = db.Where("role_id = ?", *req.RoleID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").Preload("Role").Preload("Reviewer").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("id DESC").
		Find(&apps).Error

	return apps, total, err
}

// FindPendingByUserID 查找用户待审核的申请
func FindPendingByUserID(userID uint) (*model.RoleApplication, error) {
	var app model.RoleApplication
	err := database.DB.Where("user_id = ? AND status = 0", userID).
		Order("id DESC").First(&app).Error
	return &app, err
}

// FindRecentRejectedByUserID 查找用户最近被驳回的申请
func FindRecentRejectedByUserID(userID uint, within time.Duration) (*model.RoleApplication, error) {
	var app model.RoleApplication
	since := time.Now().Add(-within)
	err := database.DB.Where("user_id = ? AND status = 2 AND reviewed_at > ?", userID, since).
		Order("id DESC").First(&app).Error
	return &app, err
}

// CountByStatus 统计各状态申请数量
func CountByStatus() (map[int8]int64, error) {
	var results []struct {
		Status int8
		Count  int64
	}
	err := database.DB.Model(&model.RoleApplication{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int8]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}
