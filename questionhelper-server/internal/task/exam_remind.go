package task

import (
	"time"

	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ExamRemindTask 考试提醒任务
type ExamRemindTask struct{}

// Run 执行考试提醒
func (t *ExamRemindTask) Run() {
	logger.Info("执行考试提醒任务")

	now := time.Now()

	// 查找即将开始的考试（30分钟后开始）
	var exams []struct {
		ID        uint
		Title     string
		StartTime time.Time
	}
	database.DB.Raw(`
		SELECT id, title, start_time
		FROM exams
		WHERE status = 1
		AND start_time BETWEEN ? AND ?
		AND remind_sent = false
	`, now, now.Add(30*time.Minute)).Scan(&exams)

	for _, exam := range exams {
		logger.Infof("发送考试提醒: exam_id=%d, title=%s", exam.ID, exam.Title)
		// TODO: 发送通知给参加考试的用户

		// 标记已发送
		database.DB.Exec("UPDATE exams SET remind_sent = true WHERE id = ?", exam.ID)
	}

	// 查找即将结束的考试（10分钟后结束）
	database.DB.Raw(`
		SELECT id, title, end_time
		FROM exams
		WHERE status = 1
		AND end_time BETWEEN ? AND ?
		AND end_remind_sent = false
	`, now, now.Add(10*time.Minute)).Scan(&exams)

	for _, exam := range exams {
		logger.Infof("发送考试结束提醒: exam_id=%d, title=%s", exam.ID, exam.Title)
		// TODO: 发送即将结束提醒

		database.DB.Exec("UPDATE exams SET end_remind_sent = true WHERE id = ?", exam.ID)
	}
}
