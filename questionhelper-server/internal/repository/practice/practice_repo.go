package practice

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== PracticeSession ====================

func FindSessionByID(id uint) (*model.PracticeSession, error) {
	var session model.PracticeSession
	err := database.DB.First(&session, id).Error
	return &session, err
}

func CreateSession(session *model.PracticeSession) error {
	return database.DB.Create(session).Error
}

func UpdateSession(session *model.PracticeSession) error {
	return database.DB.Save(session).Error
}

func ListSessions(userID uint, req *dto.PracticeListRequest) ([]model.PracticeSession, int64, error) {
	var sessions []model.PracticeSession
	var total int64

	db := database.DB.Model(&model.PracticeSession{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&sessions).Error

	return sessions, total, err
}

// ==================== PracticeRecord ====================

func CreateRecords(records []model.PracticeRecord) error {
	if len(records) == 0 {
		return nil
	}
	return database.DB.Create(&records).Error
}

func GetRecords(sessionID uint) ([]model.PracticeRecord, error) {
	var records []model.PracticeRecord
	err := database.DB.Where("session_id = ?", sessionID).Find(&records).Error
	return records, err
}

// ==================== Stats ====================

func GetUserPracticeStats(userID uint) (map[string]interface{}, error) {
	var totalSessions int64
	var totalRecords int64
	var correctRecords int64

	database.DB.Model(&model.PracticeSession{}).Where("user_id = ?", userID).Count(&totalSessions)
	database.DB.Model(&model.PracticeRecord{}).
		Joins("JOIN practice_sessions ON practice_sessions.id = practice_records.session_id").
		Where("practice_sessions.user_id = ?", userID).Count(&totalRecords)
	database.DB.Model(&model.PracticeRecord{}).
		Joins("JOIN practice_sessions ON practice_sessions.id = practice_records.session_id").
		Where("practice_sessions.user_id = ? AND practice_records.is_correct = ?", userID, true).Count(&correctRecords)

	accuracy := float64(0)
	if totalRecords > 0 {
		accuracy = float64(correctRecords) / float64(totalRecords) * 100
	}

	return map[string]interface{}{
		"total_sessions": totalSessions,
		"total_records":  totalRecords,
		"correct_records": correctRecords,
		"accuracy":       accuracy,
	}, nil
}
