package exam

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	examRepo "questionhelper-server/internal/repository/exam"
)

// ==================== Score ====================

// ListScores 成绩列表
func ListScores(examID *uint, req *dto.PageRequest) ([]dto.ExamRecordInfo, int64, error) {
	records, total, err := examRepo.ListExamRecords(examID, nil, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询成绩列表失败: %w", err)
	}

	list := make([]dto.ExamRecordInfo, 0, len(records))
	for _, r := range records {
		list = append(list, dto.ExamRecordInfo{
			ID:         r.ID,
			ExamID:     r.ExamID,
			UserID:     r.UserID,
			Score:      r.Score,
			Status:     r.Status,
			StartTime:  r.StartTime,
			SubmitTime: r.SubmitTime,
		})
	}
	return list, total, nil
}

// GetScore 获取成绩详情
func GetScore(recordID uint) (interface{}, error) {
	record, err := examRepo.FindExamRecord(recordID)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}

	answers, _ := examRepo.GetAnswerRecords(recordID)

	return map[string]interface{}{
		"record":  record,
		"answers": answers,
	}, nil
}

// GetScoreAnalysis 成绩分析
func GetScoreAnalysis(examID uint) (map[string]interface{}, error) {
	return examRepo.GetExamScoreStats(examID)
}
