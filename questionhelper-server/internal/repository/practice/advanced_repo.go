package practice

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== MockExam ====================

func FindMockExamConfigByID(id uint) (*model.MockExamConfig, error) {
	var config model.MockExamConfig
	err := database.DB.First(&config, id).Error
	return &config, err
}

func CreateMockExamSession(session *model.PracticeSession) error {
	return database.DB.Create(session).Error
}

func FindMockExamSessionByID(id uint) (*model.PracticeSession, error) {
	var session model.PracticeSession
	err := database.DB.First(&session, id).Error
	return &session, err
}

func UpdateMockExamSession(session *model.PracticeSession) error {
	return database.DB.Save(session).Error
}

func ListMockExamHistory(userID uint, req *dto.MockExamHistoryRequest) ([]model.PracticeSession, int64, error) {
	var sessions []model.PracticeSession
	var total int64

	db := database.DB.Model(&model.PracticeSession{}).Where("user_id = ?", userID)

	if req.ConfigID != nil {
		db = db.Where("category_id = ?", *req.ConfigID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&sessions).Error

	return sessions, total, err
}

// ==================== PracticePlan ====================

func CreatePlan(plan *model.PracticePlan) error {
	return database.DB.Create(plan).Error
}

func FindPlanByID(id uint) (*model.PracticePlan, error) {
	var plan model.PracticePlan
	err := database.DB.First(&plan, id).Error
	return &plan, err
}

func UpdatePlan(plan *model.PracticePlan) error {
	return database.DB.Save(plan).Error
}

func DeletePlan(id uint) error {
	return database.DB.Delete(&model.PracticePlan{}, id).Error
}

func ListPlans(userID uint, req *dto.PlanListRequest) ([]model.PracticePlan, int64, error) {
	var plans []model.PracticePlan
	var total int64

	db := database.DB.Model(&model.PracticePlan{}).Where("user_id = ?", userID)

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&plans).Error

	return plans, total, err
}

func CreatePlanLog(log *model.PracticePlanLog) error {
	return database.DB.Create(log).Error
}

// ==================== DailyPractice ====================

func FindDailyPractice(userID uint, date string) (*model.DailyPractice, error) {
	var dp model.DailyPractice
	err := database.DB.Where("user_id = ? AND date = ?", userID, date).First(&dp).Error
	return &dp, err
}

func CreateDailyPractice(dp *model.DailyPractice) error {
	return database.DB.Create(dp).Error
}

func UpdateDailyPractice(dp *model.DailyPractice) error {
	return database.DB.Save(dp).Error
}

// ==================== PracticeCheckin ====================

func FindCheckin(userID uint, date string) (*model.PracticeCheckin, error) {
	var checkin model.PracticeCheckin
	err := database.DB.Where("user_id = ? AND date = ?", userID, date).First(&checkin).Error
	return &checkin, err
}

func CreateCheckin(checkin *model.PracticeCheckin) error {
	return database.DB.Create(checkin).Error
}

func UpdateCheckin(checkin *model.PracticeCheckin) error {
	return database.DB.Save(checkin).Error
}

func ListCheckins(userID uint, year, month int) ([]model.PracticeCheckin, error) {
	var checkins []model.PracticeCheckin
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	err := database.DB.Where("user_id = ? AND date >= ? AND date < ?",
		userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date ASC").Find(&checkins).Error

	return checkins, err
}

// ==================== Leaderboard ====================

func GetLeaderboard(rankType int8, limit int) ([]model.PracticeLeaderboard, error) {
	var items []model.PracticeLeaderboard
	today := time.Now().Format("2006-01-02")

	err := database.DB.Where("rank_type = ? AND rank_date = ?", rankType, today).
		Order("rank_pos ASC").Limit(limit).Find(&items).Error

	return items, err
}

func CreateLeaderboard(item *model.PracticeLeaderboard) error {
	return database.DB.Create(item).Error
}

// ==================== UserAbilityProfile ====================

func FindAbilityProfile(userID uint) (*model.UserAbilityProfile, error) {
	var profile model.UserAbilityProfile
	err := database.DB.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func CreateAbilityProfile(profile *model.UserAbilityProfile) error {
	return database.DB.Create(profile).Error
}

func UpdateAbilityProfile(profile *model.UserAbilityProfile) error {
	return database.DB.Save(profile).Error
}

// ==================== ChallengeLevel ====================

func ListChallengeLevels() ([]model.ChallengeLevel, error) {
	var levels []model.ChallengeLevel
	err := database.DB.Where("status = 1").Order("sort ASC, level ASC").Find(&levels).Error
	return levels, err
}

func FindChallengeLevelByID(id uint) (*model.ChallengeLevel, error) {
	var level model.ChallengeLevel
	err := database.DB.First(&level, id).Error
	return &level, err
}

func FindChallengeLevelByLevel(level int) (*model.ChallengeLevel, error) {
	var cl model.ChallengeLevel
	err := database.DB.Where("level = ? AND status = 1", level).First(&cl).Error
	return &cl, err
}

// ==================== UserChallengeProgress ====================

func FindChallengeProgress(userID, levelID uint) (*model.UserChallengeProgress, error) {
	var progress model.UserChallengeProgress
	err := database.DB.Where("user_id = ? AND challenge_level = ?", userID, levelID).First(&progress).Error
	return &progress, err
}

func CreateChallengeProgress(progress *model.UserChallengeProgress) error {
	return database.DB.Create(progress).Error
}

func UpdateChallengeProgress(progress *model.UserChallengeProgress) error {
	return database.DB.Save(progress).Error
}

func ListChallengeProgress(userID uint) ([]model.UserChallengeProgress, error) {
	var progress []model.UserChallengeProgress
	err := database.DB.Where("user_id = ?", userID).Find(&progress).Error
	return progress, err
}
