package wrong

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== WrongQuestion ====================

func FindByID(id uint) (*model.WrongQuestion, error) {
	var wrong model.WrongQuestion
	err := database.DB.Preload("Question").Preload("Question.Options").First(&wrong, id).Error
	return &wrong, err
}

func FindByIDAndUser(id, userID uint) (*model.WrongQuestion, error) {
	var wrong model.WrongQuestion
	err := database.DB.Preload("Question").Preload("Question.Options").
		Where("id = ? AND user_id = ?", id, userID).First(&wrong).Error
	return &wrong, err
}

func FindByUserAndQuestion(userID, questionID uint) (*model.WrongQuestion, error) {
	var wrong model.WrongQuestion
	err := database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).First(&wrong).Error
	return &wrong, err
}

func Create(wrong *model.WrongQuestion) error {
	return database.DB.Create(wrong).Error
}

func Update(wrong *model.WrongQuestion) error {
	return database.DB.Save(wrong).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.WrongQuestion{}, id).Error
}

func DeleteByIDs(ids []uint) error {
	return database.DB.Delete(&model.WrongQuestion{}, ids).Error
}

func List(userID uint, req *dto.WrongListRequest) ([]model.WrongQuestion, int64, error) {
	var wrongs []model.WrongQuestion
	var total int64

	db := database.DB.Model(&model.WrongQuestion{}).Where("user_id = ?", userID)

	if req.Mastered != nil {
		db = db.Where("mastered = ?", *req.Mastered)
	}
	if req.Source != nil {
		db = db.Where("source = ?", *req.Source)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Question").Preload("Question.Options").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("updated_at DESC").Find(&wrongs).Error

	return wrongs, total, err
}

// Search 搜索错题
func Search(userID uint, req *dto.WrongSearchRequest) ([]model.WrongQuestion, int64, error) {
	var wrongs []model.WrongQuestion
	var total int64

	db := database.DB.Model(&model.WrongQuestion{}).Where("wrong_questions.user_id = ?", userID)

	if req.Keyword != "" {
		db = db.Joins("LEFT JOIN questions ON questions.id = wrong_questions.question_id").
			Where("(questions.title LIKE ? OR questions.content LIKE ?)",
				"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.CategoryID != nil {
		if req.Keyword == "" {
			db = db.Joins("LEFT JOIN questions ON questions.id = wrong_questions.question_id")
		}
		db = db.Where("questions.category_id = ?", *req.CategoryID)
	}
	if req.Difficulty != nil {
		if req.Keyword == "" && req.CategoryID == nil {
			db = db.Joins("LEFT JOIN questions ON questions.id = wrong_questions.question_id")
		}
		db = db.Where("questions.difficulty = ?", *req.Difficulty)
	}
	if req.Type != nil {
		if req.Keyword == "" && req.CategoryID == nil && req.Difficulty == nil {
			db = db.Joins("LEFT JOIN questions ON questions.id = wrong_questions.question_id")
		}
		db = db.Where("questions.type = ?", *req.Type)
	}
	if req.Source != nil {
		db = db.Where("wrong_questions.source = ?", *req.Source)
	}
	if req.Mastered != nil {
		db = db.Where("wrong_questions.mastered = ?", *req.Mastered)
	}
	if req.IsFavorite != nil {
		db = db.Where("wrong_questions.is_favorite = ?", *req.IsFavorite)
	}
	if req.TagID != nil {
		db = db.Joins("JOIN wrong_question_tag_relations ON wrong_question_tag_relations.wrong_question_id = wrong_questions.id").
			Where("wrong_question_tag_relations.tag_id = ?", *req.TagID)
	}
	if req.StartDate != "" {
		db = db.Where("wrong_questions.created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		db = db.Where("wrong_questions.created_at <= ?", req.EndDate+" 23:59:59")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := "wrong_questions.updated_at DESC"
	if req.SortBy != "" {
		direction := "DESC"
		if req.SortOrder == "asc" {
			direction = "ASC"
		}
		switch req.SortBy {
		case "created_at":
			orderClause = "wrong_questions.created_at " + direction
		case "wrong_count":
			orderClause = "wrong_questions.wrong_count " + direction
		case "difficulty":
			if req.Keyword == "" && req.CategoryID == nil && req.Difficulty == nil && req.Type == nil {
				db = db.Joins("LEFT JOIN questions ON questions.id = wrong_questions.question_id")
			}
			orderClause = "questions.difficulty " + direction
		}
	}

	err := db.Preload("Question").Preload("Question.Options").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order(orderClause).Find(&wrongs).Error

	return wrongs, total, err
}

// FindByIDs 根据ID列表查询错题
func FindByIDs(ids []uint, userID uint) ([]model.WrongQuestion, error) {
	var wrongs []model.WrongQuestion
	err := database.DB.Where("id IN ? AND user_id = ?", ids, userID).
		Preload("Question").Preload("Question.Options").
		Find(&wrongs).Error
	return wrongs, err
}

// FindDueForReview 查询今日待复习错题
func FindDueForReview(userID uint) ([]model.WrongQuestion, error) {
	var wrongs []model.WrongQuestion
	now := time.Now()
	err := database.DB.Where("user_id = ? AND mastered = ? AND (next_review_at IS NULL OR next_review_at <= ?)",
		userID, false, now).
		Preload("Question").Preload("Question.Options").
		Order("next_review_at ASC").
		Find(&wrongs).Error
	return wrongs, err
}

// GetAnalysis 错题分析
func GetAnalysis(userID uint) (map[string]interface{}, error) {
	var totalCount int64
	var masteredCount int64

	database.DB.Model(&model.WrongQuestion{}).Where("user_id = ?", userID).Count(&totalCount)
	database.DB.Model(&model.WrongQuestion{}).Where("user_id = ? AND mastered = ?", userID, true).Count(&masteredCount)

	// 按分类统计
	var byCategory []struct {
		CategoryID   uint   `json:"category_id"`
		CategoryName string `json:"category_name"`
		Count        int    `json:"count"`
	}
	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.category_id, categories.name as category_name, COUNT(*) as count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Joins("LEFT JOIN categories ON categories.id = questions.category_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.category_id, categories.name").
		Scan(&byCategory)

	// 按题型统计
	var byType []struct {
		Type  int8 `json:"type"`
		Count int  `json:"count"`
	}
	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.type, COUNT(*) as count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.type").
		Scan(&byType)

	// 按难度统计
	var byDifficulty []struct {
		Difficulty int8 `json:"difficulty"`
		Count      int  `json:"count"`
	}
	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.difficulty, COUNT(*) as count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.difficulty").
		Scan(&byDifficulty)

	return map[string]interface{}{
		"total_count":    totalCount,
		"mastered_count": masteredCount,
		"by_category":    byCategory,
		"by_type":        byType,
		"by_difficulty":  byDifficulty,
	}, nil
}

// GetTrendAnalysis 错题趋势分析
func GetTrendAnalysis(userID, startDate, endDate, interval string) ([]dto.WrongTrendInfo, error) {
	var results []dto.WrongTrendInfo

	dateFormat := "%Y-%m-%d"
	switch interval {
	case "week":
		dateFormat = "%Y-%u"
	case "month":
		dateFormat = "%Y-%m"
	}

	rows, err := database.DB.Model(&model.WrongQuestion{}).
		Select("DATE_FORMAT(created_at, ? as date, COUNT(*) as count", dateFormat).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate+" 23:59:59").
		Group("date").Order("date ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info dto.WrongTrendInfo
		if err := database.DB.ScanRows(rows, &info); err != nil {
			continue
		}
		results = append(results, info)
	}

	return results, nil
}

// GetCategoryAnalysis 错题分类分析
func GetCategoryAnalysis(userID uint) ([]dto.WrongCategoryAnalysisInfo, error) {
	var results []dto.WrongCategoryAnalysisInfo

	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.category_id, categories.name as category_name, "+
			"COUNT(*) as total_count, "+
			"SUM(CASE WHEN wrong_questions.mastered = true THEN 1 ELSE 0 END) as mastered_count, "+
			"AVG(wrong_questions.wrong_count) as avg_wrong_count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Joins("LEFT JOIN categories ON categories.id = questions.category_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.category_id, categories.name").
		Scan(&results)

	return results, nil
}

// GetAccuracyAnalysis 正确率分析
func GetAccuracyAnalysis(userID uint) (*dto.WrongAccuracyInfo, error) {
	info := &dto.WrongAccuracyInfo{}

	// 总体统计
	var totalReviews, correctReviews int64
	database.DB.Model(&model.WrongQuestionReview{}).
		Where("user_id = ?", userID).
		Count(&totalReviews)
	database.DB.Model(&model.WrongQuestionReview{}).
		Where("user_id = ? AND is_correct = ?", userID, true).
		Count(&correctReviews)

	info.TotalReviews = int(totalReviews)
	info.CorrectReviews = int(correctReviews)
	if totalReviews > 0 {
		info.AccuracyRate = float64(correctReviews) / float64(totalReviews) * 100
	}

	// 按难度统计
	var byDifficulty []dto.DifficultyAccuracyInfo
	database.DB.Model(&model.WrongQuestionReview{}).
		Select("questions.difficulty, "+
			"COUNT(*) as total_count, "+
			"SUM(CASE WHEN wrong_question_reviews.is_correct = true THEN 1 ELSE 0 END) as correct_count").
		Joins("JOIN wrong_questions ON wrong_questions.id = wrong_question_reviews.wrong_question_id").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Where("wrong_question_reviews.user_id = ?", userID).
		Group("questions.difficulty").
		Scan(&byDifficulty)
	for i := range byDifficulty {
		if byDifficulty[i].TotalCount > 0 {
			byDifficulty[i].AccuracyRate = float64(byDifficulty[i].CorrectCount) / float64(byDifficulty[i].TotalCount) * 100
		}
	}
	info.ByDifficulty = byDifficulty

	// 按题型统计
	var byType []dto.TypeAccuracyInfo
	database.DB.Model(&model.WrongQuestionReview{}).
		Select("questions.type, "+
			"COUNT(*) as total_count, "+
			"SUM(CASE WHEN wrong_question_reviews.is_correct = true THEN 1 ELSE 0 END) as correct_count").
		Joins("JOIN wrong_questions ON wrong_questions.id = wrong_question_reviews.wrong_question_id").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Where("wrong_question_reviews.user_id = ?", userID).
		Group("questions.type").
		Scan(&byType)
	for i := range byType {
		if byType[i].TotalCount > 0 {
			byType[i].AccuracyRate = float64(byType[i].CorrectCount) / float64(byType[i].TotalCount) * 100
		}
	}
	info.ByType = byType

	return info, nil
}

// ==================== WrongQuestionTag ====================

func FindTagByID(id uint) (*model.WrongQuestionTag, error) {
	var tag model.WrongQuestionTag
	err := database.DB.First(&tag, id).Error
	return &tag, err
}

func FindTagsByUserID(userID uint) ([]model.WrongQuestionTag, error) {
	var tags []model.WrongQuestionTag
	err := database.DB.Where("user_id = ?", userID).Order("sort ASC, id ASC").Find(&tags).Error
	return tags, err
}

func CreateTag(tag *model.WrongQuestionTag) error {
	return database.DB.Create(tag).Error
}

func UpdateTag(tag *model.WrongQuestionTag) error {
	return database.DB.Save(tag).Error
}

func DeleteTag(id uint) error {
	return database.DB.Delete(&model.WrongQuestionTag{}, id).Error
}

func ExistsTagByName(userID uint, name string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.WrongQuestionTag{}).
		Where("user_id = ? AND name = ?", userID, name).Count(&count).Error
	return count > 0, err
}

// AddTagToWrongQuestion 为错题添加标签
func AddTagToWrongQuestion(wrongQuestionID, tagID uint) error {
	relation := &model.WrongQuestionTagRelation{
		WrongQuestionID: wrongQuestionID,
		TagID:           tagID,
	}
	return database.DB.Create(relation).Error
}

// RemoveTagFromWrongQuestion 移除错题标签
func RemoveTagFromWrongQuestion(wrongQuestionID, tagID uint) error {
	return database.DB.Where("wrong_question_id = ? AND tag_id = ?", wrongQuestionID, tagID).
		Delete(&model.WrongQuestionTagRelation{}).Error
}

// FindTagsByWrongQuestionID 获取错题的标签列表
func FindTagsByWrongQuestionID(wrongQuestionID uint) ([]model.WrongQuestionTag, error) {
	var tags []model.WrongQuestionTag
	err := database.DB.Joins("JOIN wrong_question_tag_relations ON wrong_question_tag_relations.tag_id = wrong_question_tags.id").
		Where("wrong_question_tag_relations.wrong_question_id = ?", wrongQuestionID).
		Find(&tags).Error
	return tags, err
}

// UpdateTagCount 更新标签关联错题数量
func UpdateTagCount(tagID uint) error {
	var count int64
	database.DB.Model(&model.WrongQuestionTagRelation{}).Where("tag_id = ?", tagID).Count(&count)
	return database.DB.Model(&model.WrongQuestionTag{}).Where("id = ?", tagID).
		Update("count", count).Error
}

// BatchAddTagsToWrongQuestion 批量为错题添加标签
func BatchAddTagsToWrongQuestion(wrongQuestionID uint, tagIDs []uint) error {
	relations := make([]model.WrongQuestionTagRelation, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		relations = append(relations, model.WrongQuestionTagRelation{
			WrongQuestionID: wrongQuestionID,
			TagID:           tagID,
		})
	}
	return database.DB.CreateInBatches(relations, 100).Error
}

// ==================== WrongQuestionNote ====================

func FindNoteByID(id uint) (*model.WrongQuestionNote, error) {
	var note model.WrongQuestionNote
	err := database.DB.First(&note, id).Error
	return &note, err
}

func FindNotesByWrongQuestionID(wrongQuestionID uint) ([]model.WrongQuestionNote, error) {
	var notes []model.WrongQuestionNote
	err := database.DB.Where("wrong_question_id = ?", wrongQuestionID).
		Order("is_pinned DESC, created_at DESC").Find(&notes).Error
	return notes, err
}

func CreateNote(note *model.WrongQuestionNote) error {
	return database.DB.Create(note).Error
}

func UpdateNote(note *model.WrongQuestionNote) error {
	return database.DB.Save(note).Error
}

func DeleteNote(id uint) error {
	return database.DB.Delete(&model.WrongQuestionNote{}, id).Error
}

// ==================== WrongQuestionAttachment ====================

func FindAttachmentByID(id uint) (*model.WrongQuestionAttachment, error) {
	var attachment model.WrongQuestionAttachment
	err := database.DB.First(&attachment, id).Error
	return &attachment, err
}

func FindAttachmentsByWrongQuestionID(wrongQuestionID uint) ([]model.WrongQuestionAttachment, error) {
	var attachments []model.WrongQuestionAttachment
	err := database.DB.Where("wrong_question_id = ?", wrongQuestionID).
		Order("sort ASC, id ASC").Find(&attachments).Error
	return attachments, err
}

func CreateAttachment(attachment *model.WrongQuestionAttachment) error {
	return database.DB.Create(attachment).Error
}

func DeleteAttachment(id uint) error {
	return database.DB.Delete(&model.WrongQuestionAttachment{}, id).Error
}

// ==================== WrongQuestionFavorite ====================

func FindFavorite(userID, wrongQuestionID uint) (*model.WrongQuestionFavorite, error) {
	var fav model.WrongQuestionFavorite
	err := database.DB.Where("user_id = ? AND wrong_question_id = ?", userID, wrongQuestionID).
		First(&fav).Error
	return &fav, err
}

func CreateFavorite(fav *model.WrongQuestionFavorite) error {
	return database.DB.Create(fav).Error
}

func DeleteFavorite(userID, wrongQuestionID uint) error {
	return database.DB.Where("user_id = ? AND wrong_question_id = ?", userID, wrongQuestionID).
		Delete(&model.WrongQuestionFavorite{}).Error
}

func ListFavorites(userID uint, req *dto.WrongFavoriteListRequest) ([]model.WrongQuestionFavorite, int64, error) {
	var favorites []model.WrongQuestionFavorite
	var total int64

	db := database.DB.Model(&model.WrongQuestionFavorite{}).Where("user_id = ?", userID)
	if req.FolderID != nil {
		db = db.Where("folder_id = ?", *req.FolderID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&favorites).Error

	return favorites, total, err
}

// ==================== WrongQuestionReview ====================

func CreateReview(review *model.WrongQuestionReview) error {
	return database.DB.Create(review).Error
}

func FindReviewsByWrongQuestionID(wrongQuestionID uint) ([]model.WrongQuestionReview, error) {
	var reviews []model.WrongQuestionReview
	err := database.DB.Where("wrong_question_id = ?", wrongQuestionID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func ListReviews(userID uint, req *dto.WrongReviewListRequest) ([]model.WrongQuestionReview, int64, error) {
	var reviews []model.WrongQuestionReview
	var total int64

	db := database.DB.Model(&model.WrongQuestionReview{}).Where("user_id = ?", userID)
	if req.WrongQuestionID != nil {
		db = db.Where("wrong_question_id = ?", *req.WrongQuestionID)
	}
	if req.IsCorrect != nil {
		db = db.Where("is_correct = ?", *req.IsCorrect)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&reviews).Error

	return reviews, total, err
}

// FindReviewsByDateRange 按日期范围查询复习记录
func FindReviewsByDateRange(userID uint, startDate, endDate time.Time) ([]model.WrongQuestionReview, error) {
	var reviews []model.WrongQuestionReview
	err := database.DB.Where("user_id = ? AND created_at >= ? AND created_at <= ?",
		userID, startDate, endDate).
		Find(&reviews).Error
	return reviews, err
}
