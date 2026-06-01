package user

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// AssignRoles 分配用户角色
func AssignRoles(userID uint, roleIDs []uint) error {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	roles := make([]model.Role, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		roles = append(roles, model.Role{ID: roleID})
	}

	return database.DB.Model(&user).Association("Roles").Replace(roles)
}

func FindRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	err := database.DB.First(&role, id).Error
	return &role, err
}

func FindRolesByUserID(userID uint) ([]model.Role, error) {
	var user model.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

func ListRoles(req *dto.RoleListRequest) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	db := database.DB.Model(&model.Role{})

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("sort ASC, id ASC").
		Find(&roles).Error

	return roles, total, err
}

func CreateRole(role *model.Role) error {
	return database.DB.Create(role).Error
}

func UpdateRole(role *model.Role) error {
	return database.DB.Save(role).Error
}

func DeleteRoleByID(id uint) error {
	return database.DB.Delete(&model.Role{}, id).Error
}

// CountUsersByRoleID 统计关联了指定角色的用户数量
func CountUsersByRoleID(roleID uint) (int64, error) {
	var count int64
	err := database.DB.Table("user_roles").Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}
