package wrong

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	wrongRepo "questionhelper-server/internal/repository/wrong"
	"questionhelper-server/pkg/logger"
)

// 艾宾浩斯遗忘曲线间隔（单位：分钟）
var ebbinghausIntervals = []int{
	1,           // 20分钟后
	60,          // 1小时后
	720,         // 12小时后 (约半天)
	1440,        // 1天后
	2880,        // 2天后
	4320,        // 3天后 (原为5天，更紧凑)
	10080,       // 7天后 (约1周)
	30240,       // 21天后 (原为15天，更紧凑)
	43200,       // 30天后
}

// ListWrongQuestions 错题列表
func ListWrongQuestions(userID uint, req *dto.WrongListRequest) ([]dto.WrongQuestionInfo, int64, error) {
	wrongs, total, err := wrongRepo.List(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询错题列表失败: %w", err)
	}

	list := make([]dto.WrongQuestionInfo, 0, len(wrongs))
	for _, w := range wrongs {
		list = append(list, toWrongQuestionInfo(&w))
	}
	return list, total, nil
}

// SearchWrongQuestions 搜索错题
func SearchWrongQuestions(userID uint, req *dto.WrongSearchRequest) ([]dto.WrongQuestionInfo, int64, error) {
	wrongs, total, err := wrongRepo.Search(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索错题失败: %w", err)
	}

	list := make([]dto.WrongQuestionInfo, 0, len(wrongs))
	for _, w := range wrongs {
		list = append(list, toWrongQuestionInfo(&w))
	}
	return list, total, nil
}

// GetWrongQuestion 获取错题详情
func GetWrongQuestion(id, userID uint) (*dto.WrongQuestionInfo, error) {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("错题不存在")
		}
		return nil, fmt.Errorf("查询错题失败: %w", err)
	}

	if wrong.UserID != userID {
		return nil, errors.New("无权查看此错题")
	}

	info := toWrongQuestionInfo(wrong)
	return &info, nil
}

// ReviewWrongQuestion 复习错题（兼容旧接口）
func ReviewWrongQuestion(id, userID uint, answer string) (bool, error) {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		return false, errors.New("错题不存在")
	}

	if wrong.UserID != userID {
		return false, errors.New("无权操作此错题")
	}

	isCorrect := wrong.Question.Answer == answer

	req := &dto.ReviewWrongAdvancedRequest{
		IsCorrect:  isCorrect,
		Answer:     answer,
		ReviewType: 1,
	}

	if err := processReview(wrong, req); err != nil {
		return false, err
	}

	return isCorrect, nil
}

// ReviewWrongQuestionAdvanced 高级复习错题（支持艾宾浩斯算法）
func ReviewWrongQuestionAdvanced(id, userID uint, req *dto.ReviewWrongAdvancedRequest) error {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		return errors.New("错题不存在")
	}

	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	return processReview(wrong, req)
}

// processReview 处理复习逻辑
func processReview(wrong *model.WrongQuestion, req *dto.ReviewWrongAdvancedRequest) error {
	now := time.Now()

	// 创建复习记录
	review := &model.WrongQuestionReview{
		WrongQuestionID: wrong.ID,
		UserID:          wrong.UserID,
		IsCorrect:       req.IsCorrect,
		Answer:          req.Answer,
		Duration:        req.Duration,
		ReviewType:      req.ReviewType,
	}

	if req.IsCorrect {
		wrong.CorrectStreak++
		wrong.Mastered = wrong.CorrectStreak >= len(ebbinghausIntervals)
	} else {
		wrong.CorrectStreak = 0
		wrong.Mastered = false
		wrong.WrongCount++
	}

	wrong.LastAnswer = req.Answer
	wrong.LastReviewAt = &now
	wrong.ReviewCount++

	// 计算下次复习时间（艾宾浩斯遗忘曲线）
	nextReview := calculateNextReview(wrong.CorrectStreak)
	wrong.NextReviewAt = &nextReview
	review.NextReviewAt = &nextReview

	if err := wrongRepo.Update(wrong); err != nil {
		return fmt.Errorf("更新错题失败: %w", err)
	}

	if err := wrongRepo.CreateReview(review); err != nil {
		logger.Errorf("创建复习记录失败: %v", err)
	}

	logger.Infof("用户 %d 复习错题 %d, 结果: %v, 连续正确: %d, 下次复习: %v",
		wrong.UserID, wrong.ID, req.IsCorrect, wrong.CorrectStreak, nextReview)
	return nil
}

// calculateNextReview 根据艾宾浩斯遗忘曲线计算下次复习时间
func calculateNextReview(correctStreak int) time.Time {
	if correctStreak <= 0 {
		// 答错了，20分钟后重新复习
		return time.Now().Add(time.Duration(ebbinghausIntervals[0]) * time.Minute)
	}

	// 使用正确的间隔索引（连续答对次数-1，因为第一次答对对应第1个间隔）
	index := correctStreak - 1
	if index >= len(ebbinghausIntervals) {
		index = len(ebbinghausIntervals) - 1
	}

	interval := ebbinghausIntervals[index]
	return time.Now().Add(time.Duration(interval) * time.Minute)
}

// RemoveWrongQuestion 移除错题
func RemoveWrongQuestion(id, userID uint) error {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("错题不存在")
		}
		return fmt.Errorf("查询错题失败: %w", err)
	}

	if wrong.UserID != userID {
		return errors.New("无权删除此错题")
	}

	if err := wrongRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除错题失败: %w", err)
	}

	logger.Infof("用户 %d 删除错题 %d", userID, id)
	return nil
}

// BatchDeleteWrongQuestions 批量删除错题
func BatchDeleteWrongQuestions(ids []uint, userID uint) (*dto.BatchResult, error) {
	result := &dto.BatchResult{
		Total: len(ids),
	}

	for _, id := range ids {
		wrong, err := wrongRepo.FindByID(id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: "错题不存在",
			})
			continue
		}

		if wrong.UserID != userID {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: "无权删除此错题",
			})
			continue
		}

		if err := wrongRepo.DeleteByID(id); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	logger.Infof("用户 %d 批量删除错题完成: %d/%d", userID, result.Success, result.Total)
	return result, nil
}

// BatchExportWrongQuestions 批量导出错题（CSV格式）
func BatchExportWrongQuestions(ids []uint, userID uint) (*bytes.Buffer, error) {
	wrongs, err := wrongRepo.FindByIDs(ids, userID)
	if err != nil {
		return nil, fmt.Errorf("查询错题失败: %w", err)
	}

	return exportToCSV(wrongs)
}

// ExportWrongQuestions 导出错题（支持筛选条件）
func ExportWrongQuestions(userID uint, req *dto.WrongExportRequest) (*bytes.Buffer, error) {
	var wrongs []model.WrongQuestion
	var err error

	if len(req.IDs) > 0 {
		wrongs, err = wrongRepo.FindByIDs(req.IDs, userID)
		if err != nil {
			return nil, fmt.Errorf("查询错题失败: %w", err)
		}
	} else {
		searchReq := &dto.WrongSearchRequest{
			Keyword:    req.Keyword,
			CategoryID: req.CategoryID,
			Difficulty: req.Difficulty,
			Mastered:   req.Mastered,
		}
		searchReq.PageSize = 10000 // 导出上限
		var total int64
		wrongs, total, err = wrongRepo.Search(userID, searchReq)
		if err != nil {
			return nil, fmt.Errorf("查询错题失败: %w", err)
		}
		_ = total
	}

	return exportToCSV(wrongs)
}

// exportToCSV 导出为CSV格式
func exportToCSV(wrongs []model.WrongQuestion) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("\xEF\xBB\xBF") // BOM for Excel Chinese support

	w := csv.NewWriter(buf)

	// 写入表头
	if err := w.Write([]string{"ID", "题目标题", "题目类型", "难度", "来源", "错误次数", "是否掌握", "连续正确", "复习次数", "上次答案", "创建时间"}); err != nil {
		return nil, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	for _, wr := range wrongs {
		typeName := getTypeName(wr.Question.Type)
		diffName := getDifficultyName(wr.Question.Difficulty)
		sourceName := getSourceName(wr.Source)
		masteredStr := "否"
		if wr.Mastered {
			masteredStr = "是"
		}

		record := []string{
			fmt.Sprintf("%d", wr.ID),
			wr.Question.Title,
			typeName,
			diffName,
			sourceName,
			fmt.Sprintf("%d", wr.WrongCount),
			masteredStr,
			fmt.Sprintf("%d", wr.CorrectStreak),
			fmt.Sprintf("%d", wr.ReviewCount),
			wr.LastAnswer,
			wr.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(record); err != nil {
			return nil, fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("刷新CSV写入器失败: %w", err)
	}

	return buf, nil
}

func getTypeName(t int8) string {
	switch t {
	case 1:
		return "单选题"
	case 2:
		return "多选题"
	case 3:
		return "判断题"
	case 4:
		return "填空题"
	case 5:
		return "简答题"
	default:
		return "未知"
	}
}

func getDifficultyName(d int8) string {
	switch d {
	case 1:
		return "简单"
	case 2:
		return "中等"
	case 3:
		return "困难"
	default:
		return "未知"
	}
}

func getSourceName(s int8) string {
	switch s {
	case 1:
		return "练习"
	case 2:
		return "考试"
	default:
		return "未知"
	}
}

// BatchTagWrongQuestions 批量为错题打标签
func BatchTagWrongQuestions(ids []uint, tagIDs []uint, userID uint) (*dto.BatchResult, error) {
	result := &dto.BatchResult{
		Total: len(ids),
	}

	for _, id := range ids {
		wrong, err := wrongRepo.FindByID(id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: "错题不存在",
			})
			continue
		}

		if wrong.UserID != userID {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: "无权操作此错题",
			})
			continue
		}

		if err := wrongRepo.BatchAddTagsToWrongQuestion(id, tagIDs); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	// 更新标签计数
	for _, tagID := range tagIDs {
		wrongRepo.UpdateTagCount(tagID)
	}

	logger.Infof("用户 %d 批量打标签完成: %d/%d", userID, result.Success, result.Total)
	return result, nil
}

// ==================== 错题标签 ====================

// ListWrongTags 获取用户错题标签列表
func ListWrongTags(userID uint) ([]dto.WrongTagInfo, error) {
	tags, err := wrongRepo.FindTagsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}

	list := make([]dto.WrongTagInfo, 0, len(tags))
	for _, tag := range tags {
		list = append(list, toWrongTagInfo(&tag))
	}
	return list, nil
}

// CreateWrongTag 创建错题标签
func CreateWrongTag(userID uint, req *dto.CreateWrongTagRequest) error {
	exists, err := wrongRepo.ExistsTagByName(userID, req.Name)
	if err != nil {
		return fmt.Errorf("检查标签名称失败: %w", err)
	}
	if exists {
		return errors.New("标签名称已存在")
	}

	tag := &model.WrongQuestionTag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Icon:   req.Icon,
		Sort:   req.Sort,
	}

	if err := wrongRepo.CreateTag(tag); err != nil {
		return fmt.Errorf("创建标签失败: %w", err)
	}

	logger.Infof("用户 %d 创建错题标签: %s", userID, req.Name)
	return nil
}

// UpdateWrongTag 更新错题标签
func UpdateWrongTag(id, userID uint, req *dto.UpdateWrongTagRequest) error {
	tag, err := wrongRepo.FindTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	if tag.UserID != userID {
		return errors.New("无权修改此标签")
	}

	if req.Name != "" {
		if req.Name != tag.Name {
			exists, err := wrongRepo.ExistsTagByName(userID, req.Name)
			if err != nil {
				return fmt.Errorf("检查标签名称失败: %w", err)
			}
			if exists {
				return errors.New("标签名称已存在")
			}
		}
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}
	if req.Icon != "" {
		tag.Icon = req.Icon
	}
	tag.Sort = req.Sort

	if err := wrongRepo.UpdateTag(tag); err != nil {
		return fmt.Errorf("更新标签失败: %w", err)
	}

	logger.Infof("用户 %d 更新错题标签 %d", userID, id)
	return nil
}

// DeleteWrongTag 删除错题标签
func DeleteWrongTag(id, userID uint) error {
	tag, err := wrongRepo.FindTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	if tag.UserID != userID {
		return errors.New("无权删除此标签")
	}

	if err := wrongRepo.DeleteTag(id); err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}

	logger.Infof("用户 %d 删除错题标签 %d", userID, id)
	return nil
}

// AddTagToWrongQuestion 为错题添加标签
func AddTagToWrongQuestion(wrongID, userID uint, tagIDs []uint) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	for _, tagID := range tagIDs {
		tag, err := wrongRepo.FindTagByID(tagID)
		if err != nil {
			return fmt.Errorf("标签 %d 不存在", tagID)
		}
		if tag.UserID != userID {
			return fmt.Errorf("无权使用标签 %d", tagID)
		}
	}

	if err := wrongRepo.BatchAddTagsToWrongQuestion(wrongID, tagIDs); err != nil {
		return fmt.Errorf("添加标签失败: %w", err)
	}

	for _, tagID := range tagIDs {
		wrongRepo.UpdateTagCount(tagID)
	}

	logger.Infof("用户 %d 为错题 %d 添加标签 %v", userID, wrongID, tagIDs)
	return nil
}

// RemoveTagFromWrongQuestion 移除错题标签
func RemoveTagFromWrongQuestion(wrongID, tagID, userID uint) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	if err := wrongRepo.RemoveTagFromWrongQuestion(wrongID, tagID); err != nil {
		return fmt.Errorf("移除标签失败: %w", err)
	}

	wrongRepo.UpdateTagCount(tagID)

	logger.Infof("用户 %d 移除错题 %d 的标签 %d", userID, wrongID, tagID)
	return nil
}

// ==================== 错题备注 ====================

// ListWrongNotes 获取错题备注列表
func ListWrongNotes(wrongID, userID uint) ([]dto.WrongNoteInfo, error) {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return nil, errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return nil, errors.New("无权查看此错题")
	}

	notes, err := wrongRepo.FindNotesByWrongQuestionID(wrongID)
	if err != nil {
		return nil, fmt.Errorf("查询备注失败: %w", err)
	}

	list := make([]dto.WrongNoteInfo, 0, len(notes))
	for _, note := range notes {
		list = append(list, dto.WrongNoteInfo{
			ID:              note.ID,
			WrongQuestionID: note.WrongQuestionID,
			UserID:          note.UserID,
			Content:         note.Content,
			IsPinned:        note.IsPinned,
			CreatedAt:       note.CreatedAt,
			UpdatedAt:       note.UpdatedAt,
		})
	}
	return list, nil
}

// CreateWrongNote 创建错题备注
func CreateWrongNote(wrongID, userID uint, req *dto.CreateWrongNoteRequest) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	note := &model.WrongQuestionNote{
		WrongQuestionID: wrongID,
		UserID:          userID,
		Content:         req.Content,
	}

	if err := wrongRepo.CreateNote(note); err != nil {
		return fmt.Errorf("创建备注失败: %w", err)
	}

	logger.Infof("用户 %d 为错题 %d 创建备注", userID, wrongID)
	return nil
}

// UpdateWrongNote 更新错题备注
func UpdateWrongNote(wrongID, noteID, userID uint, req *dto.UpdateWrongNoteRequest) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	note, err := wrongRepo.FindNoteByID(noteID)
	if err != nil {
		return errors.New("备注不存在")
	}
	if note.WrongQuestionID != wrongID {
		return errors.New("备注不属于此错题")
	}
	if note.UserID != userID {
		return errors.New("无权修改此备注")
	}

	if req.Content != "" {
		note.Content = req.Content
	}
	if req.IsPinned != nil {
		note.IsPinned = *req.IsPinned
	}

	if err := wrongRepo.UpdateNote(note); err != nil {
		return fmt.Errorf("更新备注失败: %w", err)
	}

	logger.Infof("用户 %d 更新错题备注 %d", userID, noteID)
	return nil
}

// DeleteWrongNote 删除错题备注
func DeleteWrongNote(wrongID, noteID, userID uint) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	note, err := wrongRepo.FindNoteByID(noteID)
	if err != nil {
		return errors.New("备注不存在")
	}
	if note.WrongQuestionID != wrongID {
		return errors.New("备注不属于此错题")
	}
	if note.UserID != userID {
		return errors.New("无权删除此备注")
	}

	if err := wrongRepo.DeleteNote(noteID); err != nil {
		return fmt.Errorf("删除备注失败: %w", err)
	}

	logger.Infof("用户 %d 删除错题备注 %d", userID, noteID)
	return nil
}

// ==================== 错题附件 ====================

// ListWrongAttachments 获取错题附件列表
func ListWrongAttachments(wrongID, userID uint) ([]dto.WrongAttachmentInfo, error) {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return nil, errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return nil, errors.New("无权查看此错题")
	}

	attachments, err := wrongRepo.FindAttachmentsByWrongQuestionID(wrongID)
	if err != nil {
		return nil, fmt.Errorf("查询附件失败: %w", err)
	}

	list := make([]dto.WrongAttachmentInfo, 0, len(attachments))
	for _, att := range attachments {
		list = append(list, dto.WrongAttachmentInfo{
			ID:              att.ID,
			WrongQuestionID: att.WrongQuestionID,
			UserID:          att.UserID,
			FileID:          att.FileID,
			FileName:        att.FileName,
			FileType:        att.FileType,
			FileURL:         att.FileURL,
			FileSize:        att.FileSize,
			Sort:            att.Sort,
			CreatedAt:       att.CreatedAt,
		})
	}
	return list, nil
}

// CreateWrongAttachment 创建错题附件
func CreateWrongAttachment(wrongID, userID uint, fileName, fileType, fileURL string, fileSize int64, fileID uint) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	attachment := &model.WrongQuestionAttachment{
		WrongQuestionID: wrongID,
		UserID:          userID,
		FileID:          fileID,
		FileName:        fileName,
		FileType:        fileType,
		FileURL:         fileURL,
		FileSize:        fileSize,
	}

	if err := wrongRepo.CreateAttachment(attachment); err != nil {
		return fmt.Errorf("创建附件失败: %w", err)
	}

	logger.Infof("用户 %d 为错题 %d 上传附件: %s", userID, wrongID, fileName)
	return nil
}

// DeleteWrongAttachment 删除错题附件
func DeleteWrongAttachment(wrongID, attachmentID, userID uint) error {
	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	attachment, err := wrongRepo.FindAttachmentByID(attachmentID)
	if err != nil {
		return errors.New("附件不存在")
	}
	if attachment.WrongQuestionID != wrongID {
		return errors.New("附件不属于此错题")
	}
	if attachment.UserID != userID {
		return errors.New("无权删除此附件")
	}

	if err := wrongRepo.DeleteAttachment(attachmentID); err != nil {
		return fmt.Errorf("删除附件失败: %w", err)
	}

	logger.Infof("用户 %d 删除错题附件 %d", userID, attachmentID)
	return nil
}

// ==================== 错题收藏 ====================

// FavoriteWrongQuestion 收藏错题
func FavoriteWrongQuestion(wrongID, userID uint, req *dto.FavoriteWrongRequest) error {
	// 检查是否已收藏
	existing, err := wrongRepo.FindFavorite(userID, wrongID)
	if err == nil && existing.ID > 0 {
		return errors.New("已收藏此错题")
	}

	wrong, err := wrongRepo.FindByID(wrongID)
	if err != nil {
		return errors.New("错题不存在")
	}
	if wrong.UserID != userID {
		return errors.New("无权操作此错题")
	}

	fav := &model.WrongQuestionFavorite{
		UserID:          userID,
		WrongQuestionID: wrongID,
		FolderID:        req.FolderID,
		Note:            req.Note,
	}

	if err := wrongRepo.CreateFavorite(fav); err != nil {
		return fmt.Errorf("收藏失败: %w", err)
	}

	// 更新错题收藏标记
	wrong.IsFavorite = true
	wrongRepo.Update(wrong)

	logger.Infof("用户 %d 收藏错题 %d", userID, wrongID)
	return nil
}

// UnfavoriteWrongQuestion 取消收藏错题
func UnfavoriteWrongQuestion(wrongID, userID uint) error {
	if err := wrongRepo.DeleteFavorite(userID, wrongID); err != nil {
		return fmt.Errorf("取消收藏失败: %w", err)
	}

	// 更新错题收藏标记
	wrong, err := wrongRepo.FindByID(wrongID)
	if err == nil {
		wrong.IsFavorite = false
		wrongRepo.Update(wrong)
	}

	logger.Infof("用户 %d 取消收藏错题 %d", userID, wrongID)
	return nil
}

// ListWrongFavorites 收藏列表
func ListWrongFavorites(userID uint, req *dto.WrongFavoriteListRequest) ([]dto.WrongFavoriteInfo, int64, error) {
	favorites, total, err := wrongRepo.ListFavorites(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询收藏列表失败: %w", err)
	}

	list := make([]dto.WrongFavoriteInfo, 0, len(favorites))
	for _, fav := range favorites {
		info := dto.WrongFavoriteInfo{
			ID:              fav.ID,
			UserID:          fav.UserID,
			WrongQuestionID: fav.WrongQuestionID,
			FolderID:        fav.FolderID,
			Note:            fav.Note,
			CreatedAt:       fav.CreatedAt,
		}

		// 获取错题信息
		wrong, err := wrongRepo.FindByID(fav.WrongQuestionID)
		if err == nil {
			wInfo := toWrongQuestionInfo(wrong)
			info.WrongQuestion = &wInfo
		}

		list = append(list, info)
	}

	return list, total, nil
}

// ==================== 错题复习 ====================

// GetTodayReviewQuestions 获取今日待复习错题
func GetTodayReviewQuestions(userID uint) ([]dto.WrongQuestionInfo, error) {
	wrongs, err := wrongRepo.FindDueForReview(userID)
	if err != nil {
		return nil, fmt.Errorf("查询待复习错题失败: %w", err)
	}

	list := make([]dto.WrongQuestionInfo, 0, len(wrongs))
	for _, w := range wrongs {
		list = append(list, toWrongQuestionInfo(&w))
	}
	return list, nil
}

// GetReviewHistory 获取复习历史
func GetReviewHistory(userID uint, req *dto.WrongReviewListRequest) ([]dto.WrongReviewInfo, int64, error) {
	reviews, total, err := wrongRepo.ListReviews(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询复习历史失败: %w", err)
	}

	list := make([]dto.WrongReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		list = append(list, dto.WrongReviewInfo{
			ID:              r.ID,
			WrongQuestionID: r.WrongQuestionID,
			UserID:          r.UserID,
			IsCorrect:       r.IsCorrect,
			Answer:          r.Answer,
			Duration:        r.Duration,
			ReviewType:      r.ReviewType,
			NextReviewAt:    r.NextReviewAt,
			CreatedAt:       r.CreatedAt,
		})
	}
	return list, total, nil
}

// ==================== 错题分析 ====================

// GetWrongAnalysis 错题分析
func GetWrongAnalysis(userID uint) (map[string]interface{}, error) {
	return wrongRepo.GetAnalysis(userID)
}

// GetWrongTrendAnalysis 错题趋势分析
func GetWrongTrendAnalysis(userID uint, req *dto.WrongTrendRequest) ([]dto.WrongTrendInfo, error) {
	interval := req.Interval
	if interval == "" {
		interval = "day"
	}
	return wrongRepo.GetTrendAnalysis(fmt.Sprintf("%d", userID), req.StartDate, req.EndDate, interval)
}

// GetWrongCategoryAnalysis 按分类统计
func GetWrongCategoryAnalysis(userID uint) ([]dto.WrongCategoryAnalysisInfo, error) {
	return wrongRepo.GetCategoryAnalysis(userID)
}

// GetWrongAccuracyAnalysis 正确率分析
func GetWrongAccuracyAnalysis(userID uint) (*dto.WrongAccuracyInfo, error) {
	return wrongRepo.GetAccuracyAnalysis(userID)
}

// ==================== 辅助函数 ====================

// toWrongQuestionInfo 转换为 WrongQuestionInfo DTO
func toWrongQuestionInfo(w *model.WrongQuestion) dto.WrongQuestionInfo {
	info := dto.WrongQuestionInfo{
		ID:            w.ID,
		UserID:        w.UserID,
		QuestionID:    w.QuestionID,
		Source:        w.Source,
		SourceID:      w.SourceID,
		WrongCount:    w.WrongCount,
		LastAnswer:    w.LastAnswer,
		Mastered:      w.Mastered,
		CorrectStreak: w.CorrectStreak,
		LastReviewAt:  w.LastReviewAt,
		NextReviewAt:  w.NextReviewAt,
		ReviewCount:   w.ReviewCount,
		IsFavorite:    w.IsFavorite,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}

	if w.Question.ID > 0 {
		info.Question = dto.QuestionInfo{
			ID:         w.Question.ID,
			Title:      w.Question.Title,
			Content:    w.Question.Content,
			Type:       w.Question.Type,
			Difficulty: w.Question.Difficulty,
			Answer:     w.Question.Answer,
			Analysis:   w.Question.Analysis,
		}

		if len(w.Question.Options) > 0 {
			info.Question.Options = make([]dto.OptionInfo, 0, len(w.Question.Options))
			for _, opt := range w.Question.Options {
				info.Question.Options = append(info.Question.Options, dto.OptionInfo{
					ID:        opt.ID,
					Label:     opt.Label,
					Content:   opt.Content,
					IsCorrect: opt.IsCorrect,
				})
			}
		}
	}

	// 加载标签
	tags, err := wrongRepo.FindTagsByWrongQuestionID(w.ID)
	if err == nil && len(tags) > 0 {
		info.Tags = make([]dto.WrongTagInfo, 0, len(tags))
		for _, tag := range tags {
			info.Tags = append(info.Tags, toWrongTagInfo(&tag))
		}
	}

	return info
}

// toWrongTagInfo 转换为 WrongTagInfo DTO
func toWrongTagInfo(tag *model.WrongQuestionTag) dto.WrongTagInfo {
	return dto.WrongTagInfo{
		ID:        tag.ID,
		UserID:    tag.UserID,
		Name:      tag.Name,
		Color:     tag.Color,
		Icon:      tag.Icon,
		Sort:      tag.Sort,
		Count:     tag.Count,
		CreatedAt: tag.CreatedAt,
	}
}

// CalculateNextReviewTime 计算下次复习时间（公开接口，供外部调用）
func CalculateNextReviewTime(correctStreak int) time.Time {
	return calculateNextReview(correctStreak)
}

// GetEbbinghausIntervals 获取艾宾浩斯间隔配置（分钟）
func GetEbbinghausIntervals() []int {
	result := make([]int, len(ebbinghausIntervals))
	copy(result, ebbinghausIntervals)
	return result
}

// RetryEbbinghaus 重新计算遗忘曲线（指数退避策略）
func retryEbbinghaus(streak int, wrongCount int) time.Duration {
	base := ebbinghausIntervals[0]
	if streak > 0 && streak <= len(ebbinghausIntervals) {
		base = ebbinghausIntervals[streak-1]
	}
	// 根据错误次数增加延迟
	factor := math.Min(float64(wrongCount), 5.0)
	return time.Duration(float64(base)*(1+factor*0.2)) * time.Minute
}
