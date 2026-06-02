package practice

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/response"
)

type PracticeAdminController struct{}

func NewPracticeAdminController() *PracticeAdminController {
	return &PracticeAdminController{}
}

// GetPracticeStatistics 练习统计概览
func (c *PracticeAdminController) GetPracticeStatistics(ctx *gin.Context) {
	var stats struct {
		TotalSessions int64   `json:"total_sessions"`
		TotalUsers    int64   `json:"total_users"`
		AvgAccuracy   float64 `json:"avg_accuracy"`
		TodaySessions int64   `json:"today_sessions"`
	}
	database.DB.Table("practice_sessions").Count(&stats.TotalSessions)
	database.DB.Table("practice_sessions").Distinct("user_id").Count(&stats.TotalUsers)
	database.DB.Table("practice_sessions").Select("COALESCE(AVG(accuracy), 0)").Scan(&stats.AvgAccuracy)
	database.DB.Table("practice_sessions").Where("DATE(created_at) = CURDATE()").Count(&stats.TodaySessions)
	response.Success(ctx, stats)
}

// ListPracticeUsers 练习用户列表
func (c *PracticeAdminController) ListPracticeUsers(ctx *gin.Context) {
	var req struct {
		Page     int `form:"page,default=1"`
		PageSize int `form:"page_size,default=20"`
	}
	ctx.ShouldBindQuery(&req)

	var users []struct {
		UserID       uint    `json:"user_id"`
		Username     string  `json:"username"`
		TotalCount   int64   `json:"total_count"`
		CorrectCount int64   `json:"correct_count"`
		Accuracy     float64 `json:"accuracy"`
		LastPractice string  `json:"last_practice"`
	}

	database.DB.Raw(`
		SELECT ps.user_id, u.username,
			   COUNT(*) as total_count,
			   SUM(ps.correct_count) as correct_count,
			   ROUND(COALESCE(AVG(ps.accuracy), 0), 2) as accuracy,
			   MAX(ps.created_at) as last_practice
		FROM practice_sessions ps
		LEFT JOIN users u ON u.id = ps.user_id
		GROUP BY ps.user_id, u.username
		ORDER BY total_count DESC
		LIMIT ? OFFSET ?
	`, req.PageSize, (req.Page-1)*req.PageSize).Scan(&users)

	response.Success(ctx, users)
}

// GetPracticeUserDetail 用户练习详情
func (c *PracticeAdminController) GetPracticeUserDetail(ctx *gin.Context) {
	userID := ctx.Param("id")
	// 查询用户练习详情
	response.Success(ctx, gin.H{"user_id": userID})
}

// GetHotQuestions 热门题目
func (c *PracticeAdminController) GetHotQuestions(ctx *gin.Context) {
	var questions []struct {
		QuestionID uint    `json:"question_id"`
		Count      int64   `json:"count"`
		Accuracy   float64 `json:"accuracy"`
	}
	database.DB.Raw(`
		SELECT question_id, COUNT(*) as count,
			   ROUND(SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) as accuracy
		FROM practice_records
		GROUP BY question_id
		ORDER BY count DESC
		LIMIT 20
	`).Scan(&questions)
	response.Success(ctx, questions)
}

// GetAccuracyAnalysis 正确率分析
func (c *PracticeAdminController) GetAccuracyAnalysis(ctx *gin.Context) {
	var analysis []struct {
		CategoryID   uint    `json:"category_id"`
		CategoryName string  `json:"category_name"`
		TotalCount   int64   `json:"total_count"`
		Accuracy     float64 `json:"accuracy"`
	}
	database.DB.Raw(`
		SELECT q.category_id, c.name as category_name,
			   COUNT(*) as total_count,
			   ROUND(SUM(CASE WHEN pr.is_correct THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) as accuracy
		FROM practice_records pr
		LEFT JOIN questions q ON q.id = pr.question_id
		LEFT JOIN categories c ON c.id = q.category_id
		GROUP BY q.category_id, c.name
		HAVING total_count > 0
		ORDER BY accuracy ASC
	`).Scan(&analysis)
	response.Success(ctx, analysis)
}

// GetDifficultyDistribution 难度分布
func (c *PracticeAdminController) GetDifficultyDistribution(ctx *gin.Context) {
	var distribution []struct {
		Difficulty int   `json:"difficulty"`
		Count      int64 `json:"count"`
	}
	database.DB.Raw(`
		SELECT q.difficulty, COUNT(*) as count
		FROM practice_records pr
		LEFT JOIN questions q ON q.id = pr.question_id
		GROUP BY q.difficulty
		ORDER BY q.difficulty
	`).Scan(&distribution)
	response.Success(ctx, distribution)
}
