package user

import (
	"gorm.io/gorm"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

func FindMenuTree() ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Where("parent_id IS NULL").
		Order("sort ASC, id ASC").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC, id ASC")
		}).
		Find(&menus).Error
	return menus, err
}

func FindMenusByRoleID(roleID uint) ([]model.Menu, error) {
	var role model.Role
	if err := database.DB.Preload("Menus").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Menus, nil
}

func FindAllMenus() ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func CreateMenu(menu *model.Menu) error {
	return database.DB.Create(menu).Error
}

func UpdateMenu(menu *model.Menu) error {
	return database.DB.Save(menu).Error
}

func DeleteMenuByID(id uint) error {
	return database.DB.Delete(&model.Menu{}, id).Error
}
