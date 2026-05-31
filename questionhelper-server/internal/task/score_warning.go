package task

import (
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ScoreWarningTask 成绩预警任务
type ScoreWarningTask struct{}

// Run 执行成绩预警
func (t *ScoreWarningTask) Run() {
	logger.Info("执行成绩预警任务")

	// 查找最近完成的考试记录
	var records []struct {
		ID       uint
		ExamID   uint
		UserID   uint
		Score    float64
		ExamTitle string
	}
	database.DB.Raw(`
		SELECT er.id, er.exam_id, er.user_id, er.score, e.title as exam_title
		FROM exam_records er
		JOIN exams e ON e.id = er.exam_id
		WHERE er.status = 2
		AND er.score_warning_sent = false
		AND er.score < e.pass_score
	`).Scan(&records)

	for _, record := range records {
		logger.Infof("发送成绩预警: record_id=%d, user_id=%d, score=%.1f", record.ID, record.UserID, record.Score)
		// TODO: 发送成绩预警通知

		database.DB.Exec("UPDATE exam_records SET score_warning_sent = true WHERE id = ?", record.ID)
	}
}
