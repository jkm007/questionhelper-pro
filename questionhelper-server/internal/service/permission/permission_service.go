package permission

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ListButtonPermissions 获取按钮权限列表
func ListButtonPermissions(req *dto.ButtonPermissionListRequest) ([]dto.ButtonPermissionInfo, int64, error) {
	var total int64

	db := database.DB.Model(&model.ButtonPermission{})

	if req.MenuID != nil {
		db = db.Where("menu_id = ?", *req.MenuID)
	}
	if req.PermissionCode != "" {
		db = db.Where("permission_code LIKE ?", "%"+req.PermissionCode+"%")
	}
	if req.PermissionName != "" {
		db = db.Where("permission_name LIKE ?", "%"+req.PermissionName+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询按钮权限总数失败: %w", err)
	}

	var list []model.ButtonPermission
	err := db.Order("sort ASC, id ASC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&list).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询按钮权限列表失败: %w", err)
	}

	result := make([]dto.ButtonPermissionInfo, 0, len(list))
	for _, item := range list {
		result = append(result, dto.ButtonPermissionInfo{
			ID:             item.ID,
			MenuID:         item.MenuID,
			PermissionCode: item.PermissionCode,
			PermissionName: item.PermissionName,
			Sort:           item.Sort,
			Status:         item.Status,
			Remark:         item.Remark,
		})
	}

	return result, total, nil
}

// GetButtonPermission 获取按钮权限详情
func GetButtonPermission(id uint) (*dto.ButtonPermissionInfo, error) {
	var item model.ButtonPermission
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("按钮权限不存在")
		}
		return nil, fmt.Errorf("查询按钮权限失败: %w", err)
	}

	return &dto.ButtonPermissionInfo{
		ID:             item.ID,
		MenuID:         item.MenuID,
		PermissionCode: item.PermissionCode,
		PermissionName: item.PermissionName,
		Sort:           item.Sort,
		Status:         item.Status,
		Remark:         item.Remark,
	}, nil
}

// CreateButtonPermission 创建按钮权限
func CreateButtonPermission(req *dto.CreateButtonPermissionRequest) error {
	// 检查权限编码是否已存在
	var count int64
	if err := database.DB.Model(&model.ButtonPermission{}).
		Where("permission_code = ?", req.PermissionCode).
		Count(&count).Error; err != nil {
		return fmt.Errorf("检查权限编码失败: %w", err)
	}
	if count > 0 {
		return errors.New("权限编码已存在")
	}

	item := &model.ButtonPermission{
		MenuID:         req.MenuID,
		PermissionCode: req.PermissionCode,
		PermissionName: req.PermissionName,
		Sort:           req.Sort,
		Status:         true,
		Remark:         req.Remark,
	}

	if err := database.DB.Create(item).Error; err != nil {
		return fmt.Errorf("创建按钮权限失败: %w", err)
	}

	logger.Infof("创建按钮权限成功: %s", req.PermissionCode)
	return nil
}

// UpdateButtonPermission 更新按钮权限
func UpdateButtonPermission(id uint, req *dto.UpdateButtonPermissionRequest) error {
	var item model.ButtonPermission
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("按钮权限不存在")
		}
		return fmt.Errorf("查询按钮权限失败: %w", err)
	}

	// 如果修改了权限编码，检查是否重复
	if req.PermissionCode != "" && req.PermissionCode != item.PermissionCode {
		var count int64
		if err := database.DB.Model(&model.ButtonPermission{}).
			Where("permission_code = ? AND id != ?", req.PermissionCode, id).
			Count(&count).Error; err != nil {
			return fmt.Errorf("检查权限编码失败: %w", err)
		}
		if count > 0 {
			return errors.New("权限编码已存在")
		}
		item.PermissionCode = req.PermissionCode
	}

	if req.MenuID != nil {
		item.MenuID = *req.MenuID
	}
	if req.PermissionName != "" {
		item.PermissionName = req.PermissionName
	}
	if req.Sort != nil {
		item.Sort = *req.Sort
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Remark != "" {
		item.Remark = req.Remark
	}

	if err := database.DB.Save(&item).Error; err != nil {
		return fmt.Errorf("更新按钮权限失败: %w", err)
	}

	logger.Infof("更新按钮权限 %d 成功", id)
	return nil
}

// DeleteButtonPermission 删除按钮权限
func DeleteButtonPermission(id uint) error {
	var item model.ButtonPermission
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("按钮权限不存在")
		}
		return fmt.Errorf("查询按钮权限失败: %w", err)
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		return fmt.Errorf("删除按钮权限失败: %w", err)
	}

	logger.Infof("删除按钮权限 %d 成功", id)
	return nil
}

// ListButtonPermissionsByMenuID 获取指定菜单下的按钮权限
func ListButtonPermissionsByMenuID(menuID uint) ([]dto.ButtonPermissionInfo, error) {
	var list []model.ButtonPermission
	err := database.DB.Where("menu_id = ? AND status = ?", menuID, true).
		Order("sort ASC, id ASC").
		Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("查询菜单按钮权限失败: %w", err)
	}

	result := make([]dto.ButtonPermissionInfo, 0, len(list))
	for _, item := range list {
		result = append(result, dto.ButtonPermissionInfo{
			ID:             item.ID,
			MenuID:         item.MenuID,
			PermissionCode: item.PermissionCode,
			PermissionName: item.PermissionName,
			Sort:           item.Sort,
			Status:         item.Status,
			Remark:         item.Remark,
		})
	}

	return result, nil
}
