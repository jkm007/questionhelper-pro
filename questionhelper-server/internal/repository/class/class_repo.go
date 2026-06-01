package class

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Class ====================

func FindByID(id uint) (*model.Class, error) {
	var class model.Class
	err := database.DB.Preload("Creator").First(&class, id).Error
	return &class, err
}

func FindByCode(code string) (*model.Class, error) {
	var class model.Class
	err := database.DB.Where("code = ?", code).First(&class).Error
	return &class, err
}

func Create(class *model.Class) error {
	return database.DB.Create(class).Error
}

func Update(class *model.Class) error {
	return database.DB.Save(class).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.Class{}, id).Error
}

func List(req *dto.ClassListRequest) ([]model.Class, int64, error) {
	var classes []model.Class
	var total int64

	db := database.DB.Model(&model.Class{})

	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Creator").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&classes).Error

	return classes, total, err
}

func ListByUser(userID uint, req *dto.PageRequest) ([]model.Class, int64, error) {
	var classes []model.Class
	var total int64

	db := database.DB.Model(&model.Class{}).
		Joins("JOIN class_members ON class_members.class_id = classes.id").
		Where("class_members.user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Creator").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("class_members.joined_at DESC").Find(&classes).Error

	return classes, total, err
}

// ==================== Homework ====================

func FindHomework(id uint) (*model.Homework, error) {
	var hw model.Homework
	err := database.DB.First(&hw, id).Error
	return &hw, err
}

func CreateHomework(hw *model.Homework) error {
	return database.DB.Create(hw).Error
}

func UpdateHomework(hw *model.Homework) error {
	return database.DB.Save(hw).Error
}

func DeleteHomework(id uint) error {
	return database.DB.Delete(&model.Homework{}, id).Error
}

func ListHomework(classID uint, req *dto.PageRequest) ([]model.Homework, int64, error) {
	var homework []model.Homework
	var total int64

	db := database.DB.Model(&model.Homework{}).Where("class_id = ?", classID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("deadline DESC").Find(&homework).Error

	return homework, total, err
}

// ==================== ClassNotice ====================

func FindNotice(id uint) (*model.ClassNotice, error) {
	var notice model.ClassNotice
	err := database.DB.First(&notice, id).Error
	return &notice, err
}

func CreateNotice(notice *model.ClassNotice) error {
	return database.DB.Create(notice).Error
}

func UpdateNotice(notice *model.ClassNotice) error {
	return database.DB.Save(notice).Error
}

func DeleteNotice(id uint) error {
	return database.DB.Delete(&model.ClassNotice{}, id).Error
}

func ListNotices(classID uint, req *dto.PageRequest) ([]model.ClassNotice, int64, error) {
	var notices []model.ClassNotice
	var total int64

	db := database.DB.Model(&model.ClassNotice{}).Where("class_id = ?", classID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&notices).Error

	return notices, total, err
}
