package comment

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	commentRepo "questionhelper-server/internal/repository/comment"
	"questionhelper-server/pkg/logger"
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

// CreateComment 创建评论
func CreateComment(userID uint, req *dto.CreateCommentRequest) error {
	comment := &model.Comment{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		UserID:     userID,
		Content:    req.Content,
		ParentID:   req.ParentID,
		Status:     1,
	}

	if err := commentRepo.Create(comment); err != nil {
		return fmt.Errorf("创建评论失败: %w", err)
	}

	logger.Infof("用户 %d 创建评论成功", userID)
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

	// 检查是否已点赞
	_, err = commentRepo.FindLike(commentID, userID)
	if err == nil {
		// 已点赞，取消点赞
		if err := commentRepo.DeleteLike(commentID, userID); err != nil {
			return fmt.Errorf("取消点赞失败: %w", err)
		}
		commentRepo.DecrementLikeCount(commentID)
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
	commentRepo.IncrementLikeCount(commentID)

	return nil
}

// ReportComment 举报评论
func ReportComment(commentID, userID uint, reason string) error {
	_, err := commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	report := &model.CommentReport{
		CommentID: commentID,
		UserID:    userID,
		Reason:    reason,
		Status:    0,
	}

	if err := commentRepo.CreateReport(report); err != nil {
		return fmt.Errorf("举报失败: %w", err)
	}

	logger.Infof("用户 %d 举报评论 %d", userID, commentID)
	return nil
}

// ListReports 举报列表（管理员）
func ListReports(req *dto.PageRequest) ([]model.CommentReport, int64, error) {
	return commentRepo.ListReports(req)
}

// HandleReport 处理举报（管理员）
func HandleReport(reportID uint, status int8) error {
	report, err := commentRepo.FindReport(reportID)
	if err != nil {
		return errors.New("举报不存在")
	}

	report.Status = status
	if err := commentRepo.UpdateReport(report); err != nil {
		return fmt.Errorf("处理举报失败: %w", err)
	}

	// 如果举报成立，隐藏评论
	if status == 1 {
		comment, err := commentRepo.FindByID(report.CommentID)
		if err == nil {
			comment.Status = 0
			commentRepo.Update(comment)
		}
	}

	logger.Infof("处理举报 %d，状态: %d", reportID, status)
	return nil
}

// toCommentInfo 转换为 CommentInfo DTO
func toCommentInfo(c *model.Comment) dto.CommentInfo {
	info := dto.CommentInfo{
		ID:         c.ID,
		TargetType: c.TargetType,
		TargetID:   c.TargetID,
		UserID:     c.UserID,
		Content:    c.Content,
		ParentID:   c.ParentID,
		LikeCount:  c.LikeCount,
		Status:     c.Status,
		CreatedAt:  c.CreatedAt,
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
