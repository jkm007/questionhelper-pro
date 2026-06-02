package key

import "fmt"

// 考试相关缓存 Key（统一使用 qh: 前缀）

// ExamKey 考试信息缓存 Key
func ExamKey(examID uint) string {
	return fmt.Sprintf("qh:exam:%d", examID)
}

// ExamRecordKey 考试记录缓存 Key
func ExamRecordKey(recordID uint) string {
	return fmt.Sprintf("qh:exam:record:%d", recordID)
}

// ExamAnswerKey 答题缓存 Key（用于防作弊）
func ExamAnswerKey(recordID uint, questionID uint) string {
	return fmt.Sprintf("qh:exam:answer:%d:%d", recordID, questionID)
}

// ExamTimerKey 考试计时器 Key
func ExamTimerKey(recordID uint) string {
	return fmt.Sprintf("qh:exam:timer:%d", recordID)
}

// ExamMonitorKey 考试监控 Key
func ExamMonitorKey(examID uint) string {
	return fmt.Sprintf("qh:exam:monitor:%d", examID)
}

// ExamExpire 考试缓存过期时间
const ExamExpire = 86400        // 24小时
const ExamRecordExpire = 604800 // 7天
