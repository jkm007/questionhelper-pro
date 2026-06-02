package comment

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	commentRepo "questionhelper-server/internal/repository/comment"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/sensitive"
)

// ListComments 评论列表
func ListComments(req *dto.CommentListRequest) ([]dto.CommentInfo, int64, error) {
	comments, total, err := commentRepo.List(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论列表失败: %w", err)
	}

	list := make([]dto.CommentInfo, 0, len(comments))
	for _, c := range comments {
		list = append(list, toCommentInfo(&c))
	}

	// 构建评论树
	tree := buildCommentTree(list)

	return tree, total, nil
}

// ListCommentsAdmin 管理员评论列表
func ListCommentsAdmin(req *dto.CommentAdminListRequest) ([]dto.CommentInfo, int64, error) {
	comments, total, err := commentRepo.ListByAdmin(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论列表失败: %w", err)
	}

	list := make([]dto.CommentInfo, 0, len(comments))
	for _, c := range comments {
		list = append(list, toCommentInfo(&c))
	}

	return list, total, nil
}

// CreateComment 创建评论
func CreateComment(userID uint, req *dto.CreateCommentRequest) error {
	// 检查回复层级限制
	if req.ParentID != nil && *req.ParentID > 0 {
		depth := getCommentDepth(*req.ParentID)
		if depth >= 2 {
			return errors.New("评论回复最多支持2层嵌套")
		}
	}

	// 检查黑名单
	if err := checkBlacklist(userID, req.TargetType, req.TargetID); err != nil {
		return err
	}

	// 敏感词检查
	filter := sensitive.NewFilter()
	if filter.HasSensitive(req.Content) {
		return errors.New("评论内容包含敏感词，请修改后重试")
	}

	comment := &model.Comment{
		TargetType:    req.TargetType,
		TargetID:      req.TargetID,
		UserID:        userID,
		Content:       req.Content,
		ParentID:      req.ParentID,
		Images:        req.Images,
		Mentions:      req.Mentions,
		ReplyToUserID: req.ReplyToUserID,
		Status:        1,
	}

	// 审核规则匹配
	if err := applyAuditRules(comment); err != nil {
		logger.Warnf("审核规则匹配出错: %v", err)
	}

	if err := commentRepo.Create(comment); err != nil {
		return fmt.Errorf("创建评论失败: %w", err)
	}

	// 评论回复通知
	if req.ParentID != nil && *req.ParentID > 0 {
		parentComment, err := commentRepo.FindByID(*req.ParentID)
		if err == nil && parentComment.UserID != userID {
			truncated := truncateContent(comment.Content, 50)
			notification := &model.Notification{
				UserID:     parentComment.UserID,
				Type:       1, // comment_reply
				Title:      "评论回复",
				Content:    fmt.Sprintf("有人回复了您的评论：%s", truncated),
				TargetType: fmt.Sprintf("%d", req.TargetType),
				TargetID:   req.TargetID,
			}
			database.DB.Create(notification)
		}
	}

	logger.Infof("用户 %d 创建评论成功", userID)
	return nil
}

// applyAuditRules 对评论应用审核规则
func applyAuditRules(comment *model.Comment) error {
	rules, err := commentRepo.FindEnabledAuditRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if matchAuditRule(comment.Content, rule) {
			switch rule.Action {
			case 1: // 标记待审
				comment.Status = 0
			case 2: // 自动隐藏
				comment.Status = 0
			case 3: // 自动删除（标记为隐藏）
				comment.Status = 0
			case 4: // 拒绝发布
				comment.Status = 0
			}
			break
		}
	}
	return nil
}

// matchAuditRule 匹配审核规则
func matchAuditRule(content string, rule model.CommentAuditRule) bool {
	switch rule.RuleType {
	case "keyword":
		// 关键词匹配：检查内容是否包含规则中的关键词
		keywords := strings.Split(rule.Pattern, ",")
		for _, kw := range keywords {
			kw = strings.TrimSpace(kw)
			if kw != "" && strings.Contains(content, kw) {
				return true
			}
		}
	case "regex":
		// 正则匹配
		matched, _ := regexp.MatchString(rule.Pattern, content)
		return matched
	case "length":
		// 长度匹配：规则pattern为最大长度值
		maxLen, err := strconv.Atoi(strings.TrimSpace(rule.Pattern))
		if err == nil && len(content) > maxLen {
			return true
		}
	case "repeat":
		// 重复字符检测：规则pattern为重复次数阈值
		threshold, err := strconv.Atoi(strings.TrimSpace(rule.Pattern))
		if err == nil && hasRepeatedChars(content, threshold) {
			return true
		}
	}
	return false
}

// hasRepeatedChars 检测内容中是否有连续重复的字符
func hasRepeatedChars(content string, threshold int) bool {
	if threshold <= 0 || len(content) == 0 {
		return false
	}
	runes := []rune(content)
	count := 1
	for i := 1; i < len(runes); i++ {
		if runes[i] == runes[i-1] {
			count++
			if count >= threshold {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// truncateContent 截断内容到指定长度
func truncateContent(content string, maxLen int) string {
	runes := []rune(content)
	if len(runes) <= maxLen {
		return content
	}
	return string(runes[:maxLen]) + "..."
}

// EditComment 编辑评论
func EditComment(id, userID uint, req *dto.UpdateCommentRequest) error {
	comment, err := commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return fmt.Errorf("查询评论失败: %w", err)
	}

	// 只能编辑自己的评论
	if comment.UserID != userID {
		return errors.New("无权编辑此评论")
	}

	// 敏感词检查
	filter := sensitive.NewFilter()
	if filter.HasSensitive(req.Content) {
		return errors.New("评论内容包含敏感词，请修改后重试")
	}

	comment.Content = req.Content
	if req.Images != "" {
		comment.Images = req.Images
	}

	if err := commentRepo.Update(comment); err != nil {
		return fmt.Errorf("编辑评论失败: %w", err)
	}

	logger.Infof("用户 %d 编辑评论 %d", userID, id)
	return nil
}

// DeleteComment 删除评论
func DeleteComment(id, userID uint) error {
	comment, err := commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return fmt.Errorf("查询评论失败: %w", err)
	}

	// 只能删除自己的评论
	if comment.UserID != userID {
		return errors.New("无权删除此评论")
	}

	if err := commentRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除评论失败: %w", err)
	}

	logger.Infof("用户 %d 删除评论 %d", userID, id)
	return nil
}

// LikeComment 点赞/取消点赞评论
func LikeComment(commentID, userID uint) error {
	_, err := commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 使用事务确保点赞/取消点赞操作的原子性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否已点赞
		_, err = commentRepo.FindLike(commentID, userID)
		if err == nil {
			// 已点赞，取消点赞
			if err := commentRepo.DeleteLike(commentID, userID); err != nil {
				return fmt.Errorf("取消点赞失败: %w", err)
			}
			if err := commentRepo.DecrementLikeCount(commentID); err != nil {
				return fmt.Errorf("更新点赞计数失败: %w", err)
			}
			return nil
		}

		// 未点赞，添加点赞
		like := &model.CommentLike{
			CommentID: commentID,
			UserID:    userID,
		}
		if err := commentRepo.CreateLike(like); err != nil {
			return fmt.Errorf("点赞失败: %w", err)
		}
		if err := commentRepo.IncrementLikeCount(commentID); err != nil {
			return fmt.Errorf("更新点赞计数失败: %w", err)
		}

		return nil
	})
}

// ReportComment 举报评论
func ReportComment(commentID, userID uint, req *dto.ReportCommentRequest) error {
	comment, err := commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 检查是否举报自己的评论
	if comment.UserID == userID {
		return errors.New("不能举报自己的评论")
	}

	// 检查是否重复举报
	existing, _ := commentRepo.FindReportByUserAndComment(userID, commentID)
	if existing != nil && existing.ID > 0 {
		return errors.New("您已举报过该评论")
	}

	report := &model.CommentReport{
		CommentID:  commentID,
		UserID:     userID,
		Reason:     req.Reason,
		ReasonType: req.ReasonType,
		Status:     0,
	}

	if err := commentRepo.CreateReport(report); err != nil {
		return fmt.Errorf("举报失败: %w", err)
	}

	logger.Infof("用户 %d 举报评论 %d", userID, commentID)
	return nil
}

// ==================== Comment Management (Admin) ====================

// PinComment 置顶/取消置顶评论
func PinComment(id uint, pinned bool) error {
	comment, err := commentRepo.FindByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	comment.IsPinned = pinned
	if err := commentRepo.Update(comment); err != nil {
		return fmt.Errorf("更新置顶状态失败: %w", err)
	}

	logger.Infof("评论 %d 置顶状态: %v", id, pinned)
	return nil
}

// FeatureComment 精选/取消精选评论
func FeatureComment(id uint, featured bool) error {
	comment, err := commentRepo.FindByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	comment.IsFeatured = featured
	if err := commentRepo.Update(comment); err != nil {
		return fmt.Errorf("更新精选状态失败: %w", err)
	}

	logger.Infof("评论 %d 精选状态: %v", id, featured)
	return nil
}

// OfficialComment 标记/取消官方解答
func OfficialComment(id uint, official bool) error {
	comment, err := commentRepo.FindByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	comment.IsOfficial = official
	if err := commentRepo.Update(comment); err != nil {
		return fmt.Errorf("更新官方解答状态失败: %w", err)
	}

	logger.Infof("评论 %d 官方解答状态: %v", id, official)
	return nil
}

// ==================== Sticker ====================

// ListStickers 获取表情包列表
func ListStickers(category string) ([]dto.StickerInfo, error) {
	stickers, err := commentRepo.ListStickers(category)
	if err != nil {
		return nil, fmt.Errorf("查询表情包失败: %w", err)
	}

	list := make([]dto.StickerInfo, 0, len(stickers))
	for _, s := range stickers {
		list = append(list, dto.StickerInfo{
			ID:        s.ID,
			Name:      s.Name,
			Category:  s.Category,
			ImageURL:  s.ImageURL,
			SortOrder: s.SortOrder,
		})
	}
	return list, nil
}

// ListStickerCategories 获取表情分类
func ListStickerCategories() ([]dto.StickerCategoryInfo, error) {
	return commentRepo.ListStickerCategories()
}

// ==================== User Search ====================

// SearchUsers 搜索用户（用于@功能）
func SearchUsers(keyword string, limit int) ([]dto.UserSearchInfo, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	users, err := commentRepo.SearchUsers(keyword, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	list := make([]dto.UserSearchInfo, 0, len(users))
	for _, u := range users {
		list = append(list, dto.UserSearchInfo{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
		})
	}
	return list, nil
}

// ==================== Blacklist (Admin) ====================

// ListBlacklists 黑名单列表
func ListBlacklists(req *dto.BlacklistListRequest) ([]dto.BlacklistInfo, int64, error) {
	list, total, err := commentRepo.ListBlacklists(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询黑名单失败: %w", err)
	}

	result := make([]dto.BlacklistInfo, 0, len(list))
	for _, bl := range list {
		info := dto.BlacklistInfo{
			ID:         bl.ID,
			UserID:     bl.UserID,
			TargetType: bl.TargetType,
			TargetID:   bl.TargetID,
			Reason:     bl.Reason,
			OperatorID: bl.OperatorID,
			ExpiresAt:  bl.ExpiresAt,
			CreatedAt:  bl.CreatedAt,
		}
		if bl.User.ID > 0 {
			info.UserName = bl.User.Nickname
			info.UserAvatar = bl.User.Avatar
		}
		result = append(result, info)
	}
	return result, total, nil
}

// AddBlacklist 添加黑名单
func AddBlacklist(operatorID uint, req *dto.AddBlacklistRequest) error {
	// 检查用户是否存在
	user, err := userRepo.FindByID(req.UserID)
	if err != nil || user.ID == 0 {
		return errors.New("用户不存在")
	}

	// 检查是否已在黑名单中
	existing, _ := commentRepo.FindBlacklistByUserAndTarget(req.UserID, req.TargetType, req.TargetID)
	if existing != nil && existing.ID > 0 {
		return errors.New("该用户已在黑名单中")
	}

	var expiresAt *time.Time
	if req.Duration > 0 {
		t := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		expiresAt = &t
	}

	bl := &model.CommentBlacklist{
		UserID:     req.UserID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
		OperatorID: operatorID,
		ExpiresAt:  expiresAt,
	}

	if err := commentRepo.CreateBlacklist(bl); err != nil {
		return fmt.Errorf("添加黑名单失败: %w", err)
	}

	logger.Infof("管理员 %d 将用户 %d 加入黑名单", operatorID, req.UserID)
	return nil
}

// RemoveBlacklist 移除黑名单
func RemoveBlacklist(id uint) error {
	bl, err := commentRepo.FindBlacklistByID(id)
	if err != nil {
		return errors.New("黑名单记录不存在")
	}

	if err := commentRepo.DeleteBlacklist(bl.ID); err != nil {
		return fmt.Errorf("移除黑名单失败: %w", err)
	}

	logger.Infof("移除黑名单 %d", id)
	return nil
}

// ==================== Audit Rules (Admin) ====================

// ListAuditRules 审核规则列表
func ListAuditRules() ([]dto.AuditRuleInfo, error) {
	rules, err := commentRepo.ListAuditRules()
	if err != nil {
		return nil, fmt.Errorf("查询审核规则失败: %w", err)
	}

	list := make([]dto.AuditRuleInfo, 0, len(rules))
	for _, r := range rules {
		list = append(list, dto.AuditRuleInfo{
			ID:          r.ID,
			Name:        r.Name,
			RuleType:    r.RuleType,
			Pattern:     r.Pattern,
			Action:      r.Action,
			Priority:    r.Priority,
			Description: r.Description,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   r.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return list, nil
}

// CreateAuditRule 创建审核规则
func CreateAuditRule(req *dto.CreateAuditRuleRequest) error {
	rule := &model.CommentAuditRule{
		Name:        req.Name,
		RuleType:    req.RuleType,
		Pattern:     req.Pattern,
		Action:      req.Action,
		Priority:    req.Priority,
		Description: req.Description,
		Status:      1,
	}

	if err := commentRepo.CreateAuditRule(rule); err != nil {
		return fmt.Errorf("创建审核规则失败: %w", err)
	}

	logger.Infof("创建审核规则: %s", req.Name)
	return nil
}

// UpdateAuditRule 更新审核规则
func UpdateAuditRule(id uint, req *dto.UpdateAuditRuleRequest) error {
	rule, err := commentRepo.FindAuditRuleByID(id)
	if err != nil {
		return errors.New("审核规则不存在")
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.RuleType != "" {
		rule.RuleType = req.RuleType
	}
	if req.Pattern != "" {
		rule.Pattern = req.Pattern
	}
	if req.Action != 0 {
		rule.Action = req.Action
	}
	if req.Priority != 0 {
		rule.Priority = req.Priority
	}
	if req.Description != "" {
		rule.Description = req.Description
	}
	if req.Status != nil {
		rule.Status = *req.Status
	}

	if err := commentRepo.UpdateAuditRule(rule); err != nil {
		return fmt.Errorf("更新审核规则失败: %w", err)
	}

	logger.Infof("更新审核规则 %d", id)
	return nil
}

// DeleteAuditRule 删除审核规则
func DeleteAuditRule(id uint) error {
	if err := commentRepo.DeleteAuditRule(id); err != nil {
		return fmt.Errorf("删除审核规则失败: %w", err)
	}

	logger.Infof("删除审核规则 %d", id)
	return nil
}

// ==================== Reports (Admin) ====================

// ListReports 举报列表（管理员）
func ListReports(req *dto.ReportListRequest) ([]dto.ReportInfo, int64, error) {
	reports, total, err := commentRepo.ListReports(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询举报列表失败: %w", err)
	}

	list := make([]dto.ReportInfo, 0, len(reports))
	for _, r := range reports {
		info := dto.ReportInfo{
			ID:           r.ID,
			CommentID:    r.CommentID,
			UserID:       r.UserID,
			Reason:       r.Reason,
			ReasonType:   r.ReasonType,
			Status:       r.Status,
			HandlerID:    r.HandlerID,
			HandleRemark: r.HandleRemark,
			HandledAt:    r.HandledAt,
			CreatedAt:    r.CreatedAt,
		}
		// 查询评论内容
		comment, err := commentRepo.FindByID(r.CommentID)
		if err == nil {
			content := comment.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			info.CommentContent = content
		}
		// 查询举报用户名
		user, err := userRepo.FindByID(r.UserID)
		if err == nil {
			info.UserName = user.Nickname
		}
		list = append(list, info)
	}
	return list, total, nil
}

// HandleReport 处理举报（管理员）
func HandleReport(reportID uint, handlerID uint, req *dto.HandleReportRequest) error {
	report, err := commentRepo.FindReport(reportID)
	if err != nil {
		return errors.New("举报不存在")
	}

	report.Status = req.Status
	report.HandlerID = &handlerID
	report.HandleRemark = req.HandleRemark
	now := time.Now()
	report.HandledAt = &now

	if err := commentRepo.UpdateReport(report); err != nil {
		return fmt.Errorf("处理举报失败: %w", err)
	}

	// 如果举报成立，隐藏评论
	if req.Status == 1 {
		comment, err := commentRepo.FindByID(report.CommentID)
		if err == nil {
			comment.Status = 0
			commentRepo.Update(comment)
		}
	}

	logger.Infof("管理员 %d 处理举报 %d，状态: %d", handlerID, reportID, req.Status)
	return nil
}

// ==================== Batch Operations (Admin) ====================

// BatchAudit 批量审核评论
func BatchAudit(req *dto.BatchAuditRequest) error {
	if err := commentRepo.BatchUpdateStatus(req.IDs, req.Status); err != nil {
		return fmt.Errorf("批量审核失败: %w", err)
	}

	logger.Infof("批量审核 %d 条评论，状态: %d", len(req.IDs), req.Status)
	return nil
}

// BatchDelete 批量删除评论
func BatchDelete(ids []uint) error {
	if err := commentRepo.DeleteByIDs(ids); err != nil {
		return fmt.Errorf("批量删除失败: %w", err)
	}

	logger.Infof("批量删除 %d 条评论", len(ids))
	return nil
}

// ==================== Stats (Admin) ====================

// GetCommentStats 获取评论统计
func GetCommentStats() (*dto.CommentStats, error) {
	totalCount, _ := commentRepo.CountTotal()
	pendingCount, _ := commentRepo.CountPendingReports()

	// 今日评论数
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	todayCount, _ := commentRepo.CountByDateRange(today, tomorrow)

	// 隐藏评论数
	hiddenCount, _ := commentRepo.CountByStatus(0)

	// 举报总数
	reportCount, _ := commentRepo.CountPendingReports()

	// 最近30天每日统计
	dailyStats, _ := commentRepo.DailyStats(30)

	// 评论最多的目标
	topTargets, _ := commentRepo.TopTargetStats(10)

	return &dto.CommentStats{
		TotalCount:     totalCount,
		TodayCount:     todayCount,
		PendingCount:   pendingCount,
		HiddenCount:    hiddenCount,
		ReportCount:    reportCount,
		DailyStats:     dailyStats,
		TopTargetStats: topTargets,
	}, nil
}

// ExportComments 导出评论为 Excel
func ExportComments(req *dto.CommentExportRequest) (*excelize.File, error) {
	comments, err := commentRepo.ListForExport(req)
	if err != nil {
		return nil, fmt.Errorf("导出评论失败: %w", err)
	}

	f := excelize.NewFile()
	sheet := "评论数据"
	f.NewSheet(sheet)

	// 表头
	headers := []string{"ID", "内容", "用户ID", "目标类型", "目标ID", "状态", "点赞数", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 数据行
	for i, c := range comments {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), c.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), c.Content)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), c.UserID)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), c.TargetType)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), c.TargetID)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), c.Status)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), c.LikeCount)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), c.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// 删除默认的 Sheet1
	f.DeleteSheet("Sheet1")

	return f, nil
}

// ==================== Helper Functions ====================

// getCommentDepth 获取评论的嵌套深度
// depth 0 = 顶级评论, 1 = 一级回复, 2 = 二级回复
func getCommentDepth(commentID uint) int {
	depth := 0
	currentID := commentID
	for depth < 3 { // 安全上限，防止无限循环
		comment, err := commentRepo.FindByID(currentID)
		if err != nil || comment.ParentID == nil {
			break
		}
		depth++
		currentID = *comment.ParentID
	}
	return depth
}

// checkBlacklist 检查用户是否在黑名单中
func checkBlacklist(userID uint, targetType int8, targetID uint) error {
	// 检查全局黑名单
	globalBL, _ := commentRepo.FindBlacklistByUserAndTarget(userID, 0, 0)
	if globalBL != nil && globalBL.ID > 0 {
		if globalBL.ExpiresAt == nil || globalBL.ExpiresAt.After(time.Now()) {
			return errors.New("您已被禁止评论")
		}
	}

	// 检查目标黑名单
	targetBL, _ := commentRepo.FindBlacklistByUserAndTarget(userID, targetType, targetID)
	if targetBL != nil && targetBL.ID > 0 {
		if targetBL.ExpiresAt == nil || targetBL.ExpiresAt.After(time.Now()) {
			return errors.New("您在该目标下已被禁止评论")
		}
	}

	return nil
}

// toCommentInfo 转换为 CommentInfo DTO
func toCommentInfo(c *model.Comment) dto.CommentInfo {
	info := dto.CommentInfo{
		ID:            c.ID,
		TargetType:    c.TargetType,
		TargetID:      c.TargetID,
		UserID:        c.UserID,
		Content:       c.Content,
		ParentID:      c.ParentID,
		Images:        c.Images,
		Mentions:      c.Mentions,
		ReplyToUserID: c.ReplyToUserID,
		LikeCount:     c.LikeCount,
		IsPinned:      c.IsPinned,
		IsFeatured:    c.IsFeatured,
		IsOfficial:    c.IsOfficial,
		Status:        c.Status,
		CreatedAt:     c.CreatedAt,
	}
	if c.User.ID > 0 {
		info.UserName = c.User.Nickname
		info.UserAvatar = c.User.Avatar
	}
	return info
}

// buildCommentTree 构建评论树
func buildCommentTree(comments []dto.CommentInfo) []dto.CommentInfo {
	// 按ParentID分组
	childMap := make(map[uint][]dto.CommentInfo)
	for _, c := range comments {
		if c.ParentID != nil {
			childMap[*c.ParentID] = append(childMap[*c.ParentID], c)
		}
	}

	// 只返回顶级评论，子评论放在Children中
	var tree []dto.CommentInfo
	for _, c := range comments {
		if c.ParentID == nil {
			c.Children = childMap[c.ID]
			tree = append(tree, c)
		}
	}

	return tree
}

// validateAndCleanImages 验证和清理图片JSON
func validateAndCleanImages(images string) (string, error) {
	if images == "" {
		return "", nil
	}

	var urls []string
	if err := json.Unmarshal([]byte(images), &urls); err != nil {
		return "", errors.New("图片格式无效，需要JSON数组")
	}

	// 限制最多9张图片
	if len(urls) > 9 {
		urls = urls[:9]
	}

	// 清理空URL
	var cleaned []string
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			cleaned = append(cleaned, u)
		}
	}

	if len(cleaned) == 0 {
		return "", nil
	}

	result, _ := json.Marshal(cleaned)
	return string(result), nil
}
