package question

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
)

// ==================== 题目笔记 ====================

// GetQuestionNotes 获取题目笔记
func GetQuestionNotes(questionID uint, req *dto.NoteListRequest) ([]dto.NoteInfo, int64, error) {
	notes, total, err := questionRepo.ListNotes(questionID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询笔记失败: %w", err)
	}

	items := make([]dto.NoteInfo, 0, len(notes))
	for _, n := range notes {
		items = append(items, dto.NoteInfo{
			ID:         n.ID,
			UserID:     n.UserID,
			QuestionID: n.QuestionID,
			Content:    n.Content,
			IsPublic:   n.IsPublic,
			LikeCount:  n.LikeCount,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
		})
	}

	return items, total, nil
}

// CreateNote 创建笔记
func CreateNote(userID, questionID uint, req *dto.CreateNoteRequest) (*dto.NoteInfo, error) {
	// 检查题目是否存在
	_, err := questionRepo.FindByID(questionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	// 检查是否已有笔记
	existing, _ := questionRepo.FindNoteByUserAndQuestion(userID, questionID)
	if existing != nil {
		return nil, errors.New("您已为该题目创建过笔记，请编辑已有笔记")
	}

	note := &model.QuestionNote{
		UserID:     userID,
		QuestionID: questionID,
		Content:    req.Content,
		IsPublic:   req.IsPublic,
	}

	if err := questionRepo.CreateNote(note); err != nil {
		return nil, fmt.Errorf("创建笔记失败: %w", err)
	}

	return &dto.NoteInfo{
		ID:         note.ID,
		UserID:     note.UserID,
		QuestionID: note.QuestionID,
		Content:    note.Content,
		IsPublic:   note.IsPublic,
		LikeCount:  note.LikeCount,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	}, nil
}

// UpdateNote 更新笔记
func UpdateNote(userID, questionID, noteID uint, req *dto.UpdateNoteRequest) error {
	note, err := questionRepo.FindNoteByID(noteID)
	if err != nil {
		return errors.New("笔记不存在")
	}

	if note.UserID != userID {
		return errors.New("无权修改此笔记")
	}

	if note.QuestionID != questionID {
		return errors.New("笔记与题目不匹配")
	}

	if req.Content != "" {
		note.Content = req.Content
	}
	if req.IsPublic != nil {
		note.IsPublic = *req.IsPublic
	}

	return questionRepo.UpdateNote(note)
}

// DeleteNote 删除笔记
func DeleteNote(userID, questionID, noteID uint) error {
	note, err := questionRepo.FindNoteByID(noteID)
	if err != nil {
		return errors.New("笔记不存在")
	}

	if note.UserID != userID {
		return errors.New("无权删除此笔记")
	}

	if note.QuestionID != questionID {
		return errors.New("笔记与题目不匹配")
	}

	return questionRepo.DeleteNote(noteID)
}

// ==================== 题目评价 ====================

// RateDifficulty 评价难度
func RateDifficulty(userID, questionID uint, req *dto.DifficultyRatingRequest) error {
	// 检查题目是否存在
	_, err := questionRepo.FindByID(questionID)
	if err != nil {
		return errors.New("题目不存在")
	}

	existing, _ := questionRepo.FindDifficultyRating(userID, questionID)
	if existing != nil {
		// 更新已有评价
		existing.Rating = req.Rating
		existing.IsCorrect = req.IsCorrect
		return questionRepo.UpdateDifficultyRating(existing)
	}

	// 创建新评价
	rating := &model.QuestionDifficultyRating{
		UserID:     userID,
		QuestionID: questionID,
		Rating:     req.Rating,
		IsCorrect:  req.IsCorrect,
	}

	return questionRepo.CreateDifficultyRating(rating)
}

// RateQuality 评价质量
func RateQuality(userID, questionID uint, req *dto.QualityRatingRequest) error {
	// 检查题目是否存在
	_, err := questionRepo.FindByID(questionID)
	if err != nil {
		return errors.New("题目不存在")
	}

	existing, _ := questionRepo.FindQualityRating(userID, questionID)
	if existing != nil {
		// 更新已有评价
		existing.Score = req.Score
		if req.ClarityScore > 0 {
			existing.ClarityScore = req.ClarityScore
		}
		if req.RelevanceScore > 0 {
			existing.RelevanceScore = req.RelevanceScore
		}
		if req.Comment != "" {
			existing.Comment = req.Comment
		}
		return questionRepo.UpdateQualityRating(existing)
	}

	// 创建新评价
	rating := &model.QuestionQualityRating{
		UserID:         userID,
		QuestionID:     questionID,
		Score:          req.Score,
		ClarityScore:   req.ClarityScore,
		RelevanceScore: req.RelevanceScore,
		Comment:        req.Comment,
	}

	return questionRepo.CreateQualityRating(rating)
}

// ==================== 题目纠错 ====================

// CreateCorrection 提交纠错
func CreateCorrection(userID, questionID uint, req *dto.CreateCorrectionRequest) (*dto.CorrectionInfo, error) {
	// 检查题目是否存在
	_, err := questionRepo.FindByID(questionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	correction := &model.QuestionCorrection{
		UserID:         userID,
		QuestionID:     questionID,
		CorrectionType: req.CorrectionType,
		Description:    req.Description,
		Screenshot:     req.Screenshot,
		Status:         0,
	}

	if err := questionRepo.CreateCorrection(correction); err != nil {
		return nil, fmt.Errorf("提交纠错失败: %w", err)
	}

	return &dto.CorrectionInfo{
		ID:             correction.ID,
		UserID:         correction.UserID,
		QuestionID:     correction.QuestionID,
		CorrectionType: correction.CorrectionType,
		Description:    correction.Description,
		Screenshot:     correction.Screenshot,
		Status:         correction.Status,
		CreatedAt:      correction.CreatedAt,
	}, nil
}

// GetCorrections 获取纠错列表
func GetCorrections(questionID uint, req *dto.CorrectionListRequest) ([]dto.CorrectionInfo, int64, error) {
	corrections, total, err := questionRepo.ListCorrections(questionID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询纠错列表失败: %w", err)
	}

	items := make([]dto.CorrectionInfo, 0, len(corrections))
	for _, c := range corrections {
		items = append(items, dto.CorrectionInfo{
			ID:             c.ID,
			UserID:         c.UserID,
			QuestionID:     c.QuestionID,
			CorrectionType: c.CorrectionType,
			Description:    c.Description,
			Screenshot:     c.Screenshot,
			Status:         c.Status,
			AdminReply:     c.AdminReply,
			ProcessedBy:    c.ProcessedBy,
			ProcessedAt:    c.ProcessedAt,
			Reward:         c.Reward,
			CreatedAt:      c.CreatedAt,
		})
	}

	return items, total, nil
}
