package dto

import "time"

// VersionInfo 版本信息
type VersionInfo struct {
	ID         uint      `json:"id"`
	QuestionID uint      `json:"question_id"`
	Version    int       `json:"version"`
	Title      string    `json:"title"`
	ChangeLog  string    `json:"change_log"`
	CreatorID  uint      `json:"creator_id"`
	Creator    string    `json:"creator"`
	CreatedAt  time.Time `json:"created_at"`
}

// VersionDetail 版本详情
type VersionDetail struct {
	VersionInfo
	Content    string         `json:"content"`
	Type       int8           `json:"type"`
	Difficulty int8           `json:"difficulty"`
	Answer     string         `json:"answer"`
	Analysis   string         `json:"analysis"`
	CategoryID uint           `json:"category_id"`
	Options    []OptionInfo   `json:"options"`
}
