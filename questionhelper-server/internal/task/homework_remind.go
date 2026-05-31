package task

import (
	"time"

	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// HomeworkRemindTask 作业提醒任务
type HomeworkRemindTask struct{}

// Run 执行作业提醒
func (t *HomeworkRemindTask) Run() {
	logger.Info("执行作业提醒任务")

	now := time.Now()

	// 查找即将截止的作业（24小时内截止）
	var homeworks []struct {
		ID        uint
		Title     string
		ClassID   uint
		Deadline  time.Time
	}
	database.DB.Raw(`
		SELECT id, title, class_id, deadline
		FROM homeworks
		WHERE deleted_at IS NULL
		AND deadline BETWEEN ? AND ?
		AND remind_sent = false
	`, now, now.Add(24*time.Hour)).Scan(&homeworks)

	for _, hw := range homeworks {
		// 查找未提交作业的学生
		var students []uint
		database.DB.Raw(`
			SELECT cm.user_id
			FROM class_members cm
			WHERE cm.class_id = ?
			AND cm.role = 1
			AND cm.user_id NOT IN (
				SELECT user_id FROM homework_submissions
				WHERE homework_id = ?
			)
		`, hw.ClassID, hw.ID).Scan(&students)

		if len(students) > 0 {
			logger.Infof("发送作业提醒: homework_id=%d, title=%s, 未提交人数=%d", hw.ID, hw.Title, len(students))
			// TODO: 发送通知给未提交的学生
		}

		database.DB.Exec("UPDATE homeworks SET remind_sent = true WHERE id = ?", hw.ID)
	}
}
