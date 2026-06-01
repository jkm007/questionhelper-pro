package class

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== ClassMember ====================

func FindMember(classID, userID uint) (*model.ClassMember, error) {
	var member model.ClassMember
	err := database.DB.Where("class_id = ? AND user_id = ?", classID, userID).First(&member).Error
	return &member, err
}

func CreateMember(member *model.ClassMember) error {
	return database.DB.Create(member).Error
}

func DeleteMember(classID, userID uint) error {
	return database.DB.Where("class_id = ? AND user_id = ?", classID, userID).Delete(&model.ClassMember{}).Error
}

func ListMembers(classID uint, req *dto.ClassMemberListRequest) ([]model.ClassMember, int64, error) {
	var members []model.ClassMember
	var total int64

	db := database.DB.Model(&model.ClassMember{}).Where("class_id = ?", classID)

	if req.Role != nil {
		db = db.Where("role = ?", *req.Role)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("joined_at DESC").Find(&members).Error

	return members, total, err
}

func IncrementMemberCount(classID uint) error {
	return database.DB.Model(&model.Class{}).Where("id = ?", classID).
		UpdateColumn("member_count", database.DB.Raw("member_count + 1")).Error
}

func DecrementMemberCount(classID uint) error {
	return database.DB.Model(&model.Class{}).Where("id = ?", classID).
		UpdateColumn("member_count", database.DB.Raw("GREATEST(member_count - 1, 0)")).Error
}
