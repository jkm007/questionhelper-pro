package class

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	classRepo "questionhelper-server/internal/repository/class"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ==================== Homework Extension ====================

// UpdateHomework 更新作业
func UpdateHomework(classID, homeworkID, operatorID uint, req *dto.UpdateHomeworkRequest) error {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作业不存在")
		}
		return fmt.Errorf("查询作业失败: %w", err)
	}

	if hw.ClassID != classID {
		return errors.New("作业不属于此班级")
	}

	// 检查权限
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	if req.Title != "" {
		hw.Title = req.Title
	}
	if req.Description != "" {
		hw.Description = req.Description
	}
	if !req.Deadline.IsZero() {
		hw.Deadline = req.Deadline
	}

	if err := classRepo.UpdateHomework(hw); err != nil {
		return fmt.Errorf("更新作业失败: %w", err)
	}

	logger.Infof("更新作业 %d 成功", homeworkID)
	return nil
}

// DeleteHomework 删除作业
func DeleteHomework(classID, homeworkID, operatorID uint) error {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作业不存在")
		}
		return fmt.Errorf("查询作业失败: %w", err)
	}

	if hw.ClassID != classID {
		return errors.New("作业不属于此班级")
	}

	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	if err := classRepo.DeleteHomework(homeworkID); err != nil {
		return fmt.Errorf("删除作业失败: %w", err)
	}

	logger.Infof("删除作业 %d 成功", homeworkID)
	return nil
}

// GetHomework 获取作业详情
func GetHomework(classID, homeworkID uint) (*dto.HomeworkInfo, error) {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("作业不存在")
		}
		return nil, fmt.Errorf("查询作业失败: %w", err)
	}

	if hw.ClassID != classID {
		return nil, errors.New("作业不属于此班级")
	}

	return &dto.HomeworkInfo{
		ID:          hw.ID,
		ClassID:     hw.ClassID,
		Title:       hw.Title,
		Description: hw.Description,
		Deadline:    hw.Deadline,
		CreatorID:   hw.CreatorID,
	}, nil
}

// SubmitHomework 提交作业
func SubmitHomework(classID, homeworkID, userID uint, req *dto.SubmitHomeworkRequest) error {
	// 验证作业存在
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作业不存在")
		}
		return fmt.Errorf("查询作业失败: %w", err)
	}

	if hw.ClassID != classID {
		return errors.New("作业不属于此班级")
	}

	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	// 检查截止时间
	if time.Now().After(hw.Deadline) {
		return errors.New("作业已过截止时间")
	}

	// 序列化附件
	attachments := ""
	if len(req.Attachments) > 0 {
		attachments = marshalAttachments(req.Attachments)
	}

	sub := &model.HomeworkSubmission{
		HomeworkID:  homeworkID,
		UserID:      userID,
		Content:     req.Content,
		Attachments: attachments,
		Status:      0,
		SubmittedAt: time.Now(),
	}

	if err := classRepo.CreateSubmission(sub); err != nil {
		return fmt.Errorf("提交作业失败: %w", err)
	}

	logger.Infof("用户 %d 提交作业 %d", userID, homeworkID)
	return nil
}

// GradeHomework 批改作业
func GradeHomework(classID, homeworkID, submissionID, operatorID uint, req *dto.GradeHomeworkRequest) error {
	// 验证作业存在
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作业不存在")
		}
		return fmt.Errorf("查询作业失败: %w", err)
	}

	if hw.ClassID != classID {
		return errors.New("作业不属于此班级")
	}

	// 检查权限
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	sub, err := classRepo.FindSubmission(submissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("提交记录不存在")
		}
		return fmt.Errorf("查询提交记录失败: %w", err)
	}

	if sub.HomeworkID != homeworkID {
		return errors.New("提交记录不属于此作业")
	}

	now := time.Now()
	sub.Score = &req.Score
	sub.Feedback = req.Feedback
	sub.GradedBy = &operatorID
	sub.Status = req.Status
	sub.GradedAt = &now

	if err := classRepo.UpdateSubmission(sub); err != nil {
		return fmt.Errorf("批改作业失败: %w", err)
	}

	logger.Infof("批改作业 %d 提交 %d 成功", homeworkID, submissionID)
	return nil
}

// ListHomeworkSubmissions 作业提交列表
func ListHomeworkSubmissions(classID, homeworkID uint, req *dto.PageRequest) ([]dto.HomeworkSubmissionInfo, int64, error) {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		return nil, 0, errors.New("作业不存在")
	}
	if hw.ClassID != classID {
		return nil, 0, errors.New("作业不属于此班级")
	}

	subs, total, err := classRepo.ListSubmissions(homeworkID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询提交列表失败: %w", err)
	}

	list := make([]dto.HomeworkSubmissionInfo, 0, len(subs))
	for _, s := range subs {
		userName := ""
		if user, err := userRepo.FindByID(s.UserID); err == nil {
			userName = user.Nickname
		}
		list = append(list, dto.HomeworkSubmissionInfo{
			ID:          s.ID,
			HomeworkID:  s.HomeworkID,
			UserID:      s.UserID,
			UserName:    userName,
			Content:     s.Content,
			Attachments: s.Attachments,
			Score:       s.Score,
			Feedback:    s.Feedback,
			Status:      s.Status,
			SubmittedAt: s.SubmittedAt,
			GradedAt:    s.GradedAt,
		})
	}
	return list, total, nil
}

// ==================== Peer Review ====================

// AssignPeerReview 分配互评
func AssignPeerReview(classID, homeworkID, operatorID uint, req *dto.AssignPeerReviewRequest) error {
	// 检查权限
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	// 验证作业存在
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		return errors.New("作业不存在")
	}
	if hw.ClassID != classID {
		return errors.New("作业不属于此班级")
	}

	var reviews []model.ClassPeerReview

	if req.Strategy == "assign" && len(req.Pairs) > 0 {
		// 指定配对
		for _, pair := range req.Pairs {
			reviews = append(reviews, model.ClassPeerReview{
				TargetType: req.TargetType,
				TargetID:   req.TargetID,
				ReviewerID: pair.ReviewerID,
				RevieweeID: pair.RevieweeID,
				Status:     0,
			})
		}
	} else {
		// 随机分配：获取所有已提交作业的学生
		subs, _, err := classRepo.ListSubmissions(homeworkID, &dto.PageRequest{Page: 1, PageSize: 1000})
		if err != nil {
			return fmt.Errorf("查询提交列表失败: %w", err)
		}

		if len(subs) < 2 {
			return errors.New("提交人数不足，无法分配互评")
		}

		// 简单随机分配：每人评下一个人
		for i := 0; i < len(subs); i++ {
			revieweeIdx := (i + 1) % len(subs)
			reviews = append(reviews, model.ClassPeerReview{
				TargetType: req.TargetType,
				TargetID:   req.TargetID,
				ReviewerID: subs[i].UserID,
				RevieweeID: subs[revieweeIdx].UserID,
				Status:     0,
			})
		}
	}

	if err := classRepo.BatchCreatePeerReviews(reviews); err != nil {
		return fmt.Errorf("分配互评失败: %w", err)
	}

	logger.Infof("为作业 %d 分配 %d 条互评", homeworkID, len(reviews))
	return nil
}

// ListPeerReviews 互评列表
func ListPeerReviews(classID, homeworkID uint, req *dto.PageRequest) ([]dto.PeerReviewInfo, int64, error) {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		return nil, 0, errors.New("作业不存在")
	}
	if hw.ClassID != classID {
		return nil, 0, errors.New("作业不属于此班级")
	}

	reviews, total, err := classRepo.ListPeerReviews(1, homeworkID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询互评列表失败: %w", err)
	}

	list := make([]dto.PeerReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		reviewerName := ""
		if user, err := userRepo.FindByID(r.ReviewerID); err == nil {
			reviewerName = user.Nickname
		}
		revieweeName := ""
		if user, err := userRepo.FindByID(r.RevieweeID); err == nil {
			revieweeName = user.Nickname
		}
		list = append(list, dto.PeerReviewInfo{
			ID:           r.ID,
			TargetType:   r.TargetType,
			TargetID:     r.TargetID,
			ReviewerID:   r.ReviewerID,
			ReviewerName: reviewerName,
			RevieweeID:   r.RevieweeID,
			RevieweeName: revieweeName,
			Score:        r.Score,
			Content:      r.Content,
			Status:       r.Status,
			ReviewedAt:   r.ReviewedAt,
		})
	}
	return list, total, nil
}

// GetMyPeerReviews 我的互评任务
func GetMyPeerReviews(classID, homeworkID, userID uint) ([]dto.PeerReviewInfo, error) {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		return nil, errors.New("作业不存在")
	}
	if hw.ClassID != classID {
		return nil, errors.New("作业不属于此班级")
	}

	reviews, err := classRepo.ListPeerReviewsByReviewer(1, homeworkID, userID)
	if err != nil {
		return nil, fmt.Errorf("查询互评任务失败: %w", err)
	}

	list := make([]dto.PeerReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		revieweeName := ""
		if user, err := userRepo.FindByID(r.RevieweeID); err == nil {
			revieweeName = user.Nickname
		}
		list = append(list, dto.PeerReviewInfo{
			ID:           r.ID,
			TargetType:   r.TargetType,
			TargetID:     r.TargetID,
			ReviewerID:   r.ReviewerID,
			RevieweeID:   r.RevieweeID,
			RevieweeName: revieweeName,
			Score:        r.Score,
			Content:      r.Content,
			Status:       r.Status,
			ReviewedAt:   r.ReviewedAt,
		})
	}
	return list, nil
}

// SubmitPeerReview 提交互评
func SubmitPeerReview(classID, homeworkID, reviewID, userID uint, req *dto.SubmitPeerReviewRequest) error {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		return errors.New("作业不存在")
	}
	if hw.ClassID != classID {
		return errors.New("作业不属于此班级")
	}

	review, err := classRepo.FindPeerReview(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("互评记录不存在")
		}
		return fmt.Errorf("查询互评记录失败: %w", err)
	}

	if review.ReviewerID != userID {
		return errors.New("无权操作此互评")
	}

	if review.Status == 1 {
		return errors.New("已提交过互评")
	}

	now := time.Now()
	review.Score = &req.Score
	review.Content = req.Content
	review.Status = 1
	review.ReviewedAt = &now

	if err := classRepo.UpdatePeerReview(review); err != nil {
		return fmt.Errorf("提交互评失败: %w", err)
	}

	logger.Infof("用户 %d 提交互评 %d", userID, reviewID)
	return nil
}

// GetPeerReviewResult 互评结果
func GetPeerReviewResult(classID, homeworkID uint) ([]dto.PeerReviewResult, error) {
	hw, err := classRepo.FindHomework(homeworkID)
	if err != nil {
		return nil, errors.New("作业不存在")
	}
	if hw.ClassID != classID {
		return nil, errors.New("作业不属于此班级")
	}

	// 获取所有已完成的互评
	reviews, _, err := classRepo.ListPeerReviews(1, homeworkID, &dto.PageRequest{Page: 1, PageSize: 10000})
	if err != nil {
		return nil, fmt.Errorf("查询互评列表失败: %w", err)
	}

	// 按被评价人聚合
	resultMap := make(map[uint]*dto.PeerReviewResult)
	for _, r := range reviews {
		if r.Status != 1 {
			continue
		}

		result, exists := resultMap[r.RevieweeID]
		if !exists {
			revieweeName := ""
			if user, err := userRepo.FindByID(r.RevieweeID); err == nil {
				revieweeName = user.Nickname
			}
			result = &dto.PeerReviewResult{
				UserID:   r.RevieweeID,
				UserName: revieweeName,
			}
			resultMap[r.RevieweeID] = result
		}

		if r.Score != nil {
			result.AvgScore += *r.Score
		}
		result.ReviewCount++
		result.Reviews = append(result.Reviews, dto.PeerReviewInfo{
			ID:       r.ID,
			Score:    r.Score,
			Content:  r.Content,
			Status:   r.Status,
		})
	}

	// 计算平均分
	results := make([]dto.PeerReviewResult, 0, len(resultMap))
	for _, r := range resultMap {
		if r.ReviewCount > 0 {
			r.AvgScore = r.AvgScore / float64(r.ReviewCount)
		}
		results = append(results, *r)
	}

	return results, nil
}

// ==================== Group Management ====================

// ListGroups 分组列表
func ListGroups(classID uint) ([]dto.GroupInfo, error) {
	groups, err := classRepo.ListGroups(classID)
	if err != nil {
		return nil, fmt.Errorf("查询分组列表失败: %w", err)
	}

	list := make([]dto.GroupInfo, 0, len(groups))
	for _, g := range groups {
		members, _ := classRepo.ListGroupMembers(g.ID)
		memberInfos := make([]dto.GroupMemberInfo, 0, len(members))
		for _, m := range members {
			userName := ""
			if user, err := userRepo.FindByID(m.UserID); err == nil {
				userName = user.Nickname
			}
			memberInfos = append(memberInfos, dto.GroupMemberInfo{
				ID:       m.ID,
				GroupID:  m.GroupID,
				UserID:   m.UserID,
				UserName: userName,
				Role:     m.Role,
				JoinedAt: m.JoinedAt,
			})
		}

		list = append(list, dto.GroupInfo{
			ID:          g.ID,
			ClassID:     g.ClassID,
			Name:        g.Name,
			Description: g.Description,
			LeaderID:    g.LeaderID,
			MaxMembers:  g.MaxMembers,
			MemberCount: len(members),
			Members:     memberInfos,
			CreatedAt:   g.CreatedAt,
		})
	}
	return list, nil
}

// CreateGroup 创建分组
func CreateGroup(classID, operatorID uint, req *dto.CreateGroupRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	group := &model.ClassGroup{
		ClassID:     classID,
		Name:        req.Name,
		Description: req.Description,
		LeaderID:    req.LeaderID,
		MaxMembers:  req.MaxMembers,
	}

	if err := classRepo.CreateGroup(group); err != nil {
		return fmt.Errorf("创建分组失败: %w", err)
	}

	logger.Infof("创建分组成功: %d", group.ID)
	return nil
}

// UpdateGroup 更新分组
func UpdateGroup(classID, groupID, operatorID uint, req *dto.UpdateGroupRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	group, err := classRepo.FindGroup(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分组不存在")
		}
		return fmt.Errorf("查询分组失败: %w", err)
	}

	if group.ClassID != classID {
		return errors.New("分组不属于此班级")
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}
	if req.LeaderID != nil {
		group.LeaderID = req.LeaderID
	}
	if req.MaxMembers > 0 {
		group.MaxMembers = req.MaxMembers
	}
	group.SortOrder = req.SortOrder

	if err := classRepo.UpdateGroup(group); err != nil {
		return fmt.Errorf("更新分组失败: %w", err)
	}

	logger.Infof("更新分组 %d 成功", groupID)
	return nil
}

// DeleteGroup 删除分组
func DeleteGroup(classID, groupID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	group, err := classRepo.FindGroup(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分组不存在")
		}
		return fmt.Errorf("查询分组失败: %w", err)
	}

	if group.ClassID != classID {
		return errors.New("分组不属于此班级")
	}

	if err := classRepo.DeleteGroup(groupID); err != nil {
		return fmt.Errorf("删除分组失败: %w", err)
	}

	logger.Infof("删除分组 %d 成功", groupID)
	return nil
}

// AddGroupMember 添加分组成员
func AddGroupMember(classID, groupID, operatorID uint, req *dto.AddGroupMemberRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	group, err := classRepo.FindGroup(groupID)
	if err != nil {
		return errors.New("分组不存在")
	}
	if group.ClassID != classID {
		return errors.New("分组不属于此班级")
	}

	for _, userID := range req.UserIDs {
		// 检查是否是班级成员
		if _, err := classRepo.FindMember(classID, userID); err != nil {
			continue
		}

		// 检查是否已在分组中
		if _, err := classRepo.FindGroupMember(groupID, userID); err == nil {
			continue
		}

		// 检查分组人数限制
		if group.MaxMembers > 0 {
			count := classRepo.CountGroupMembers(groupID)
			if count >= int64(group.MaxMembers) {
				return errors.New("分组已满员")
			}
		}

		member := &model.ClassGroupMember{
			GroupID:  groupID,
			UserID:   userID,
			Role:     1,
			JoinedAt: time.Now(),
		}
		if err := classRepo.CreateGroupMember(member); err != nil {
			logger.Errorf("添加分组成员失败: %v", err)
		}
	}

	logger.Infof("为分组 %d 添加 %d 名成员", groupID, len(req.UserIDs))
	return nil
}

// RemoveGroupMember 移除分组成员
func RemoveGroupMember(classID, groupID, userID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	group, err := classRepo.FindGroup(groupID)
	if err != nil {
		return errors.New("分组不存在")
	}
	if group.ClassID != classID {
		return errors.New("分组不属于此班级")
	}

	if _, err := classRepo.FindGroupMember(groupID, userID); err != nil {
		return errors.New("不是分组成员")
	}

	if err := classRepo.DeleteGroupMember(groupID, userID); err != nil {
		return fmt.Errorf("移除分组成员失败: %w", err)
	}

	logger.Infof("从分组 %d 移除成员 %d", groupID, userID)
	return nil
}

// ==================== Join Application ====================

// ApplyClass 提交加入申请
func ApplyClass(classID, userID uint, req *dto.ClassApplyRequest) error {
	// 检查班级是否存在
	class, err := classRepo.FindByID(classID)
	if err != nil {
		return errors.New("班级不存在")
	}

	if class.Status != 1 {
		return errors.New("班级已归档或不可用")
	}

	// 检查是否已是成员
	if _, err := classRepo.FindMember(classID, userID); err == nil {
		return errors.New("已是班级成员")
	}

	// 检查是否有待审批的申请
	if _, err := classRepo.FindPendingApplication(classID, userID); err == nil {
		return errors.New("已有待审批的申请")
	}

	app := &model.ClassJoinApplication{
		ClassID: classID,
		UserID:  userID,
		Reason:  req.Reason,
		Status:  0,
	}

	if err := classRepo.CreateApplication(app); err != nil {
		return fmt.Errorf("提交申请失败: %w", err)
	}

	logger.Infof("用户 %d 申请加入班级 %d", userID, classID)
	return nil
}

// ListApplications 申请列表
func ListApplications(classID uint, req *dto.PageRequest) ([]dto.ClassApplicationInfo, int64, error) {
	apps, total, err := classRepo.ListApplications(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询申请列表失败: %w", err)
	}

	list := make([]dto.ClassApplicationInfo, 0, len(apps))
	for _, a := range apps {
		userName := ""
		if a.User.ID > 0 {
			userName = a.User.Nickname
		}
		reviewName := ""
		if a.ReviewBy != nil {
			if user, err := userRepo.FindByID(*a.ReviewBy); err == nil {
				reviewName = user.Nickname
			}
		}
		list = append(list, dto.ClassApplicationInfo{
			ID:         a.ID,
			ClassID:    a.ClassID,
			UserID:     a.UserID,
			UserName:   userName,
			Reason:     a.Reason,
			Status:     a.Status,
			Remark:     a.Remark,
			ReviewBy:   a.ReviewBy,
			ReviewName: reviewName,
			ReviewAt:   a.ReviewAt,
			CreatedAt:  a.CreatedAt,
		})
	}
	return list, total, nil
}

// ApproveApplication 审批通过
func ApproveApplication(classID, appID, operatorID uint, req *dto.ClassReviewApplicationRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	app, err := classRepo.FindApplication(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("申请不存在")
		}
		return fmt.Errorf("查询申请失败: %w", err)
	}

	if app.ClassID != classID {
		return errors.New("申请不属于此班级")
	}

	if app.Status != 0 {
		return errors.New("申请已处理")
	}

	now := time.Now()
	app.Status = 1
	app.Remark = req.Remark
	app.ReviewBy = &operatorID
	app.ReviewAt = &now

	if err := classRepo.UpdateApplication(app); err != nil {
		return fmt.Errorf("审批申请失败: %w", err)
	}

	// 添加为班级成员
	member := &model.ClassMember{
		ClassID:  classID,
		UserID:   app.UserID,
		Role:     1,
		JoinedAt: time.Now(),
	}
	if err := classRepo.CreateMember(member); err != nil {
		logger.Errorf("添加班级成员失败: %v", err)
	}

	classRepo.IncrementMemberCount(classID)
	logger.Infof("审批通过用户 %d 加入班级 %d", app.UserID, classID)
	return nil
}

// RejectApplication 审批驳回
func RejectApplication(classID, appID, operatorID uint, req *dto.ClassReviewApplicationRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	app, err := classRepo.FindApplication(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("申请不存在")
		}
		return fmt.Errorf("查询申请失败: %w", err)
	}

	if app.ClassID != classID {
		return errors.New("申请不属于此班级")
	}

	if app.Status != 0 {
		return errors.New("申请已处理")
	}

	now := time.Now()
	app.Status = 2
	app.Remark = req.Remark
	app.ReviewBy = &operatorID
	app.ReviewAt = &now

	if err := classRepo.UpdateApplication(app); err != nil {
		return fmt.Errorf("审批申请失败: %w", err)
	}

	logger.Infof("审批驳回用户 %d 加入班级 %d", app.UserID, classID)
	return nil
}

// ==================== Attendance Management ====================

// ListAttendances 考勤列表
func ListAttendances(classID uint, req *dto.PageRequest) ([]dto.AttendanceInfo, int64, error) {
	atts, total, err := classRepo.ListAttendances(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询考勤列表失败: %w", err)
	}

	list := make([]dto.AttendanceInfo, 0, len(atts))
	for _, a := range atts {
		list = append(list, dto.AttendanceInfo{
			ID:          a.ID,
			ClassID:     a.ClassID,
			Title:       a.Title,
			Description: a.Description,
			Type:        a.Type,
			Deadline:    a.Deadline,
			CreatorID:   a.CreatorID,
			Status:      a.Status,
			CreatedAt:   a.CreatedAt,
		})
	}
	return list, total, nil
}

// CreateAttendance 创建考勤
func CreateAttendance(classID, operatorID uint, req *dto.CreateAttendanceRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	att := &model.ClassAttendance{
		ClassID:     classID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Deadline:    req.Deadline,
		CreatorID:   operatorID,
		Status:      1,
	}

	if err := classRepo.CreateAttendance(att); err != nil {
		return fmt.Errorf("创建考勤失败: %w", err)
	}

	logger.Infof("创建考勤成功: %d", att.ID)
	return nil
}

// UpdateAttendance 编辑考勤
func UpdateAttendance(classID, attID, operatorID uint, req *dto.UpdateAttendanceRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	att, err := classRepo.FindAttendance(attID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考勤不存在")
		}
		return fmt.Errorf("查询考勤失败: %w", err)
	}

	if att.ClassID != classID {
		return errors.New("考勤不属于此班级")
	}

	if req.Title != "" {
		att.Title = req.Title
	}
	if req.Description != "" {
		att.Description = req.Description
	}
	if req.Type > 0 {
		att.Type = req.Type
	}
	if req.Deadline != nil {
		att.Deadline = req.Deadline
	}
	if req.Status >= 0 {
		att.Status = req.Status
	}

	if err := classRepo.UpdateAttendance(att); err != nil {
		return fmt.Errorf("更新考勤失败: %w", err)
	}

	logger.Infof("更新考勤 %d 成功", attID)
	return nil
}

// DeleteAttendance 删除考勤
func DeleteAttendance(classID, attID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	att, err := classRepo.FindAttendance(attID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考勤不存在")
		}
		return fmt.Errorf("查询考勤失败: %w", err)
	}

	if att.ClassID != classID {
		return errors.New("考勤不属于此班级")
	}

	if err := classRepo.DeleteAttendance(attID); err != nil {
		return fmt.Errorf("删除考勤失败: %w", err)
	}

	logger.Infof("删除考勤 %d 成功", attID)
	return nil
}

// Checkin 学生签到
func Checkin(classID, attID, userID uint, req *dto.AttendanceCheckinRequest) error {
	att, err := classRepo.FindAttendance(attID)
	if err != nil {
		return errors.New("考勤不存在")
	}
	if att.ClassID != classID {
		return errors.New("考勤不属于此班级")
	}
	if att.Status != 1 {
		return errors.New("考勤已结束")
	}

	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	// 检查是否已签到
	if _, err := classRepo.FindAttendanceRecord(attID, userID); err == nil {
		return errors.New("已签到")
	}

	// 检查截止时间
	status := int8(1) // 正常
	if att.Deadline != nil && time.Now().After(*att.Deadline) {
		status = 2 // 迟到
	}

	rec := &model.ClassAttendanceRecord{
		AttendanceID: attID,
		UserID:       userID,
		Status:       status,
		IP:           req.IP,
		Location:     req.Location,
		SignedAt:     time.Now(),
	}

	if err := classRepo.CreateAttendanceRecord(rec); err != nil {
		return fmt.Errorf("签到失败: %w", err)
	}

	logger.Infof("用户 %d 签到考勤 %d", userID, attID)
	return nil
}

// Checkout 学生签退
func Checkout(classID, attID, userID uint, req *dto.AttendanceCheckinRequest) error {
	att, err := classRepo.FindAttendance(attID)
	if err != nil {
		return errors.New("考勤不存在")
	}
	if att.ClassID != classID {
		return errors.New("考勤不属于此班级")
	}

	rec, err := classRepo.FindAttendanceRecord(attID, userID)
	if err != nil {
		return errors.New("未签到，无法签退")
	}

	// 更新签退信息（复用 Location 字段记录签退位置）
	rec.Location = req.Location
	if err := classRepo.UpdateAttendanceRecord(rec); err != nil {
		return fmt.Errorf("签退失败: %w", err)
	}

	logger.Infof("用户 %d 签退考勤 %d", userID, attID)
	return nil
}

// ListAttendanceRecords 考勤记录
func ListAttendanceRecords(classID, attID uint, req *dto.PageRequest) ([]dto.AttendanceRecordInfo, int64, error) {
	att, err := classRepo.FindAttendance(attID)
	if err != nil {
		return nil, 0, errors.New("考勤不存在")
	}
	if att.ClassID != classID {
		return nil, 0, errors.New("考勤不属于此班级")
	}

	recs, total, err := classRepo.ListAttendanceRecordsPaged(attID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询考勤记录失败: %w", err)
	}

	list := make([]dto.AttendanceRecordInfo, 0, len(recs))
	for _, r := range recs {
		userName := ""
		if user, err := userRepo.FindByID(r.UserID); err == nil {
			userName = user.Nickname
		}
		list = append(list, dto.AttendanceRecordInfo{
			ID:           r.ID,
			AttendanceID: r.AttendanceID,
			UserID:       r.UserID,
			UserName:     userName,
			Status:       r.Status,
			Remark:       r.Remark,
			IP:           r.IP,
			Location:     r.Location,
			SignedAt:     r.SignedAt,
		})
	}
	return list, total, nil
}

// ExportAttendance 导出考勤
func ExportAttendance(classID, attID uint) ([]dto.AttendanceRecordInfo, error) {
	att, err := classRepo.FindAttendance(attID)
	if err != nil {
		return nil, errors.New("考勤不存在")
	}
	if att.ClassID != classID {
		return nil, errors.New("考勤不属于此班级")
	}

	recs, err := classRepo.ListAttendanceRecords(attID)
	if err != nil {
		return nil, fmt.Errorf("查询考勤记录失败: %w", err)
	}

	list := make([]dto.AttendanceRecordInfo, 0, len(recs))
	for _, r := range recs {
		userName := ""
		if user, err := userRepo.FindByID(r.UserID); err == nil {
			userName = user.Nickname
		}
		list = append(list, dto.AttendanceRecordInfo{
			ID:           r.ID,
			AttendanceID: r.AttendanceID,
			UserID:       r.UserID,
			UserName:     userName,
			Status:       r.Status,
			Remark:       r.Remark,
			IP:           r.IP,
			Location:     r.Location,
			SignedAt:     r.SignedAt,
		})
	}
	return list, nil
}

// ==================== Study Plan ====================

// ListStudyPlans 计划列表
func ListStudyPlans(classID uint, req *dto.PageRequest) ([]dto.StudyPlanInfo, int64, error) {
	plans, total, err := classRepo.ListStudyPlans(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询学习计划列表失败: %w", err)
	}

	list := make([]dto.StudyPlanInfo, 0, len(plans))
	for _, p := range plans {
		list = append(list, dto.StudyPlanInfo{
			ID:          p.ID,
			ClassID:     p.ClassID,
			Title:       p.Title,
			Description: p.Description,
			StartDate:   p.StartDate,
			EndDate:     p.EndDate,
			CreatorID:   p.CreatorID,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
		})
	}
	return list, total, nil
}

// CreateStudyPlan 创建计划
func CreateStudyPlan(classID, operatorID uint, req *dto.CreateStudyPlanRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	plan := &model.ClassStudyPlan{
		ClassID:     classID,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatorID:   operatorID,
		Status:      1,
	}

	if err := classRepo.CreateStudyPlan(plan); err != nil {
		return fmt.Errorf("创建学习计划失败: %w", err)
	}

	logger.Infof("创建学习计划成功: %d", plan.ID)
	return nil
}

// GetStudyPlan 计划详情
func GetStudyPlan(classID, planID uint) (*dto.StudyPlanInfo, error) {
	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学习计划不存在")
		}
		return nil, fmt.Errorf("查询学习计划失败: %w", err)
	}

	if plan.ClassID != classID {
		return nil, errors.New("学习计划不属于此班级")
	}

	items, _ := classRepo.ListStudyPlanItems(planID)
	itemInfos := make([]dto.StudyPlanItemInfo, 0, len(items))
	for _, item := range items {
		itemInfos = append(itemInfos, dto.StudyPlanItemInfo{
			ID:           item.ID,
			PlanID:       item.PlanID,
			Title:        item.Title,
			Description:  item.Description,
			ItemType:     item.ItemType,
			ResourceType: item.ResourceType,
			ResourceID:   item.ResourceID,
			TargetCount:  item.TargetCount,
			Required:     item.Required,
			SortOrder:    item.SortOrder,
			DueDate:      item.DueDate,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &dto.StudyPlanInfo{
		ID:          plan.ID,
		ClassID:     plan.ClassID,
		Title:       plan.Title,
		Description: plan.Description,
		StartDate:   plan.StartDate,
		EndDate:     plan.EndDate,
		CreatorID:   plan.CreatorID,
		Status:      plan.Status,
		Items:       itemInfos,
		CreatedAt:   plan.CreatedAt,
	}, nil
}

// UpdateStudyPlan 更新计划
func UpdateStudyPlan(classID, planID, operatorID uint, req *dto.UpdateStudyPlanRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学习计划不存在")
		}
		return fmt.Errorf("查询学习计划失败: %w", err)
	}

	if plan.ClassID != classID {
		return errors.New("学习计划不属于此班级")
	}

	if req.Title != "" {
		plan.Title = req.Title
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if !req.StartDate.IsZero() {
		plan.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		plan.EndDate = req.EndDate
	}
	if req.Status >= 0 {
		plan.Status = req.Status
	}

	if err := classRepo.UpdateStudyPlan(plan); err != nil {
		return fmt.Errorf("更新学习计划失败: %w", err)
	}

	logger.Infof("更新学习计划 %d 成功", planID)
	return nil
}

// DeleteStudyPlan 删除计划
func DeleteStudyPlan(classID, planID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学习计划不存在")
		}
		return fmt.Errorf("查询学习计划失败: %w", err)
	}

	if plan.ClassID != classID {
		return errors.New("学习计划不属于此班级")
	}

	if err := classRepo.DeleteStudyPlan(planID); err != nil {
		return fmt.Errorf("删除学习计划失败: %w", err)
	}

	logger.Infof("删除学习计划 %d 成功", planID)
	return nil
}

// AddStudyPlanItem 添加任务
func AddStudyPlanItem(classID, planID, operatorID uint, req *dto.CreateStudyPlanItemRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		return errors.New("学习计划不存在")
	}
	if plan.ClassID != classID {
		return errors.New("学习计划不属于此班级")
	}

	item := &model.ClassStudyPlanItem{
		PlanID:       planID,
		Title:        req.Title,
		Description:  req.Description,
		ItemType:     req.ItemType,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		TargetCount:  req.TargetCount,
		Required:     req.Required,
		DueDate:      req.DueDate,
	}

	if err := classRepo.CreateStudyPlanItem(item); err != nil {
		return fmt.Errorf("添加学习任务失败: %w", err)
	}

	logger.Infof("添加学习任务成功: %d", item.ID)
	return nil
}

// UpdateStudyPlanItem 更新任务
func UpdateStudyPlanItem(classID, planID, itemID, operatorID uint, req *dto.UpdateStudyPlanItemRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		return errors.New("学习计划不存在")
	}
	if plan.ClassID != classID {
		return errors.New("学习计划不属于此班级")
	}

	item, err := classRepo.FindStudyPlanItem(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学习任务不存在")
		}
		return fmt.Errorf("查询学习任务失败: %w", err)
	}

	if item.PlanID != planID {
		return errors.New("任务不属于此学习计划")
	}

	if req.Title != "" {
		item.Title = req.Title
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.ItemType > 0 {
		item.ItemType = req.ItemType
	}
	if req.ResourceType != "" {
		item.ResourceType = req.ResourceType
	}
	if req.ResourceID > 0 {
		item.ResourceID = req.ResourceID
	}
	if req.TargetCount > 0 {
		item.TargetCount = req.TargetCount
	}
	if req.Required != nil {
		item.Required = *req.Required
	}
	if req.DueDate != nil {
		item.DueDate = req.DueDate
	}

	if err := classRepo.UpdateStudyPlanItem(item); err != nil {
		return fmt.Errorf("更新学习任务失败: %w", err)
	}

	logger.Infof("更新学习任务 %d 成功", itemID)
	return nil
}

// DeleteStudyPlanItem 删除任务
func DeleteStudyPlanItem(classID, planID, itemID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		return errors.New("学习计划不存在")
	}
	if plan.ClassID != classID {
		return errors.New("学习计划不属于此班级")
	}

	item, err := classRepo.FindStudyPlanItem(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学习任务不存在")
		}
		return fmt.Errorf("查询学习任务失败: %w", err)
	}

	if item.PlanID != planID {
		return errors.New("任务不属于此学习计划")
	}

	if err := classRepo.DeleteStudyPlanItem(itemID); err != nil {
		return fmt.Errorf("删除学习任务失败: %w", err)
	}

	logger.Infof("删除学习任务 %d 成功", itemID)
	return nil
}

// CompleteStudyPlanItem 完成任务
func CompleteStudyPlanItem(classID, planID, itemID, userID uint) error {
	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		return errors.New("学习计划不存在")
	}
	if plan.ClassID != classID {
		return errors.New("学习计划不属于此班级")
	}

	item, err := classRepo.FindStudyPlanItem(itemID)
	if err != nil {
		return errors.New("学习任务不存在")
	}
	if item.PlanID != planID {
		return errors.New("任务不属于此学习计划")
	}

	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	// 查找或创建进度记录
	progress, err := classRepo.FindStudyPlanProgress(planID, itemID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now()
			progress = &model.ClassStudyPlanProgress{
				PlanID:      planID,
				ItemID:      itemID,
				UserID:      userID,
				Completed:   item.TargetCount,
				Status:      2,
				CompletedAt: &now,
			}
			if err := classRepo.CreateStudyPlanProgress(progress); err != nil {
				return fmt.Errorf("记录完成状态失败: %w", err)
			}
		} else {
			return fmt.Errorf("查询进度失败: %w", err)
		}
	} else {
		now := time.Now()
		progress.Completed = item.TargetCount
		progress.Status = 2
		progress.CompletedAt = &now
		if err := classRepo.UpdateStudyPlanProgress(progress); err != nil {
			return fmt.Errorf("更新完成状态失败: %w", err)
		}
	}

	logger.Infof("用户 %d 完成学习任务 %d", userID, itemID)
	return nil
}

// GetStudyPlanProgress 进度查看
func GetStudyPlanProgress(classID, planID uint) (*dto.StudyPlanProgressInfo, error) {
	plan, err := classRepo.FindStudyPlan(planID)
	if err != nil {
		return nil, errors.New("学习计划不存在")
	}
	if plan.ClassID != classID {
		return nil, errors.New("学习计划不属于此班级")
	}

	totalItems := classRepo.CountStudyPlanItems(planID)
	progressList, err := classRepo.ListStudyPlanProgressByPlan(planID)
	if err != nil {
		return nil, fmt.Errorf("查询进度失败: %w", err)
	}

	// 按用户聚合进度
	userMap := make(map[uint]*dto.UserStudyPlanProgress)
	for _, p := range progressList {
		up, exists := userMap[p.UserID]
		if !exists {
			userName := ""
			if user, err := userRepo.FindByID(p.UserID); err == nil {
				userName = user.Nickname
			}
			up = &dto.UserStudyPlanProgress{
				UserID:     p.UserID,
				UserName:   userName,
				TotalItems: int(totalItems),
			}
			userMap[p.UserID] = up
		}
		if p.Status == 2 {
			up.Completed++
		}
	}

	// 计算进度百分比
	userProgress := make([]dto.UserStudyPlanProgress, 0, len(userMap))
	for _, up := range userMap {
		if totalItems > 0 {
			up.Progress = float64(up.Completed) / float64(totalItems) * 100
		}
		userProgress = append(userProgress, *up)
	}

	// 计算整体完成数
	completedItems := 0
	for _, up := range userProgress {
		if up.Completed == int(totalItems) {
			completedItems++
		}
	}

	progressPercent := float64(0)
	if totalItems > 0 {
		progressPercent = float64(completedItems) / float64(len(userProgress)) * 100
	}

	return &dto.StudyPlanProgressInfo{
		TotalItems:     int(totalItems),
		CompletedItems: completedItems,
		Progress:       progressPercent,
		UserProgress:   userProgress,
	}, nil
}

// ==================== Class File ====================

// ListClassFiles 文件列表
func ListClassFiles(classID uint, req *dto.PageRequest) ([]dto.ClassFileInfo, int64, error) {
	files, total, err := classRepo.ListClassFiles(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询文件列表失败: %w", err)
	}

	list := make([]dto.ClassFileInfo, 0, len(files))
	for _, f := range files {
		list = append(list, dto.ClassFileInfo{
			ID:        f.ID,
			ClassID:   f.ClassID,
			Name:      f.Name,
			Path:      f.Path,
			Size:      f.Size,
			MimeType:  f.MimeType,
			FolderID:  f.FolderID,
			CreatorID: f.CreatorID,
			Downloads: f.Downloads,
			CreatedAt: f.CreatedAt,
		})
	}
	return list, total, nil
}

// UploadClassFile 上传文件
func UploadClassFile(classID, userID uint, name, path, mimeType string, size int64) error {
	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	f := &model.ClassFile{
		ClassID:   classID,
		Name:      name,
		Path:      path,
		Size:      size,
		MimeType:  mimeType,
		CreatorID: userID,
	}

	if err := classRepo.CreateClassFile(f); err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}

	logger.Infof("用户 %d 上传文件到班级 %d: %s", userID, classID, name)
	return nil
}

// DeleteClassFile 删除文件
func DeleteClassFile(classID, fileID, operatorID uint) error {
	f, err := classRepo.FindClassFile(fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件不存在")
		}
		return fmt.Errorf("查询文件失败: %w", err)
	}

	if f.ClassID != classID {
		return errors.New("文件不属于此班级")
	}

	// 只有上传者或班级管理员才能删除
	if f.CreatorID != operatorID {
		if err := checkClassPermission(classID, operatorID, 2); err != nil {
			return errors.New("无权删除此文件")
		}
	}

	if err := classRepo.DeleteClassFile(fileID); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	logger.Infof("删除文件 %d 成功", fileID)
	return nil
}

// GetClassFile 获取文件信息（用于下载）
func GetClassFile(classID, fileID uint) (*dto.ClassFileInfo, error) {
	f, err := classRepo.FindClassFile(fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文件不存在")
		}
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}

	if f.ClassID != classID {
		return nil, errors.New("文件不属于此班级")
	}

	// 增加下载次数
	classRepo.IncrementFileDownloads(fileID)

	return &dto.ClassFileInfo{
		ID:        f.ID,
		ClassID:   f.ClassID,
		Name:      f.Name,
		Path:      f.Path,
		Size:      f.Size,
		MimeType:  f.MimeType,
		FolderID:  f.FolderID,
		CreatorID: f.CreatorID,
		Downloads: f.Downloads + 1,
		CreatedAt: f.CreatedAt,
	}, nil
}

// ==================== Ranking ====================

// ListRanking 排名列表
func ListRanking(classID, userID uint, rankingType string) ([]dto.ClassRankingInfo, error) {
	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return nil, errors.New("不是班级成员")
	}

	rankings, err := classRepo.ListRankings(classID, rankingType)
	if err != nil {
		return nil, fmt.Errorf("查询排名失败: %w", err)
	}

	list := make([]dto.ClassRankingInfo, 0, len(rankings))
	for _, r := range rankings {
		userName := ""
		if user, err := userRepo.FindByID(r.UserID); err == nil {
			userName = user.Nickname
		}
		list = append(list, dto.ClassRankingInfo{
			UserID:      r.UserID,
			UserName:    userName,
			RankingType: r.RankingType,
			Score:       r.Score,
			Rank:        r.Rank,
		})
	}
	return list, nil
}

// CalculateRanking 触发排名计算
func CalculateRanking(classID, operatorID uint, rankingType string) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	// 获取班级成员
	members, _, err := classRepo.ListMembers(classID, &dto.ClassMemberListRequest{
		PageRequest: dto.PageRequest{Page: 1, PageSize: 1000},
	})
	if err != nil {
		return fmt.Errorf("查询班级成员失败: %w", err)
	}

	// 删除旧排名
	classRepo.DeleteRankings(classID, rankingType)

	// 计算每个成员的得分（简化实现：基于练习和考试记录）
	type userScore struct {
		UserID uint
		Score  float64
	}

	scores := make([]userScore, 0, len(members))
	for _, m := range members {
		// 这里简化处理，实际应该根据 rankingType 查询相应的统计数据
		score := 0.0
		scores = append(scores, userScore{UserID: m.UserID, Score: score})
	}

	// 按得分降序排序
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// 保存排名
	for i, s := range scores {
		ranking := &model.ClassRanking{
			ClassID:     classID,
			UserID:      s.UserID,
			RankingType: rankingType,
			Score:       s.Score,
			Rank:        i + 1,
		}
		classRepo.UpsertRanking(ranking)
	}

	logger.Infof("班级 %d 排名计算完成，类型: %s，共 %d 人", classID, rankingType, len(scores))
	return nil
}

// ==================== Tag Management ====================

// ListTags 标签列表
func ListTags() ([]dto.ClassTagInfo, error) {
	tags, err := classRepo.ListAllTags()
	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}

	list := make([]dto.ClassTagInfo, 0, len(tags))
	for _, t := range tags {
		list = append(list, dto.ClassTagInfo{
			ID:        t.ID,
			Name:      t.Name,
			CreatorID: t.CreatorID,
			SortOrder: t.SortOrder,
			Status:    t.Status,
			CreatedAt: t.CreatedAt,
		})
	}
	return list, nil
}

// CreateTag 创建标签
func CreateTag(operatorID uint, req *dto.ClassTagCreateRequest) error {
	tag := &model.ClassTag{
		Name:      req.Name,
		CreatorID: operatorID,
		SortOrder: req.SortOrder,
		Status:    1,
	}

	if err := classRepo.CreateTag(tag); err != nil {
		return fmt.Errorf("创建标签失败: %w", err)
	}

	logger.Infof("创建标签成功: %s", tag.Name)
	return nil
}

// UpdateTag 更新标签
func UpdateTag(tagID, operatorID uint, req *dto.ClassTagUpdateRequest) error {
	tag, err := classRepo.FindTag(tagID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.SortOrder > 0 {
		tag.SortOrder = req.SortOrder
	}
	if req.Status >= 0 {
		tag.Status = req.Status
	}

	if err := classRepo.UpdateTag(tag); err != nil {
		return fmt.Errorf("更新标签失败: %w", err)
	}

	logger.Infof("更新标签 %d 成功", tagID)
	return nil
}

// DeleteTag 删除标签
func DeleteTag(tagID uint) error {
	tag, err := classRepo.FindTag(tagID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	if err := classRepo.DeleteTag(tag.ID); err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}

	logger.Infof("删除标签 %d 成功", tagID)
	return nil
}

// AddClassTag 为班级添加标签
func AddClassTag(classID, operatorID uint, req *dto.AddClassTagRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	for _, tagID := range req.TagIDs {
		// 检查标签是否存在
		if _, err := classRepo.FindTag(tagID); err != nil {
			continue
		}

		if err := classRepo.AddClassTag(classID, tagID); err != nil {
			logger.Errorf("为班级 %d 添加标签 %d 失败: %v", classID, tagID, err)
		}
	}

	logger.Infof("为班级 %d 添加 %d 个标签", classID, len(req.TagIDs))
	return nil
}

// RemoveClassTag 移除班级标签
func RemoveClassTag(classID, tagID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	if err := classRepo.RemoveClassTag(classID, tagID); err != nil {
		return fmt.Errorf("移除班级标签失败: %w", err)
	}

	logger.Infof("移除班级 %d 标签 %d", classID, tagID)
	return nil
}

// ==================== Template Management ====================

// ListTemplates 模板列表
func ListTemplates(userID uint, req *dto.PageRequest, isAdmin bool) ([]dto.ClassTemplateInfo, int64, error) {
	templates, total, err := classRepo.ListTemplates(req, userID, isAdmin)
	if err != nil {
		return nil, 0, fmt.Errorf("查询模板列表失败: %w", err)
	}

	list := make([]dto.ClassTemplateInfo, 0, len(templates))
	for _, t := range templates {
		list = append(list, dto.ClassTemplateInfo{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Config:      t.Config,
			CreatorID:   t.CreatorID,
			IsPublic:    t.IsPublic,
			UsedCount:   t.UsedCount,
			Status:      t.Status,
			CreatedAt:   t.CreatedAt,
		})
	}
	return list, total, nil
}

// CreateTemplate 创建模板
func CreateTemplate(operatorID uint, req *dto.ClassTemplateCreateRequest) error {
	tpl := &model.ClassTemplate{
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		CreatorID:   operatorID,
		IsPublic:    req.IsPublic,
		Status:      1,
	}

	if err := classRepo.CreateTemplate(tpl); err != nil {
		return fmt.Errorf("创建模板失败: %w", err)
	}

	logger.Infof("创建模板成功: %s", tpl.Name)
	return nil
}

// GetTemplate 模板详情
func GetTemplate(templateID uint) (*dto.ClassTemplateInfo, error) {
	tpl, err := classRepo.FindTemplate(templateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("模板不存在")
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}

	return &dto.ClassTemplateInfo{
		ID:          tpl.ID,
		Name:        tpl.Name,
		Description: tpl.Description,
		Config:      tpl.Config,
		CreatorID:   tpl.CreatorID,
		IsPublic:    tpl.IsPublic,
		UsedCount:   tpl.UsedCount,
		Status:      tpl.Status,
		CreatedAt:   tpl.CreatedAt,
	}, nil
}

// UpdateTemplate 更新模板
func UpdateTemplate(templateID, operatorID uint, req *dto.ClassTemplateUpdateRequest) error {
	tpl, err := classRepo.FindTemplate(templateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模板不存在")
		}
		return fmt.Errorf("查询模板失败: %w", err)
	}

	// 只有创建者才能修改
	if tpl.CreatorID != operatorID {
		return errors.New("无权修改此模板")
	}

	if req.Name != "" {
		tpl.Name = req.Name
	}
	if req.Description != "" {
		tpl.Description = req.Description
	}
	if req.Config != "" {
		tpl.Config = req.Config
	}
	if req.IsPublic != nil {
		tpl.IsPublic = *req.IsPublic
	}
	if req.Status >= 0 {
		tpl.Status = req.Status
	}

	if err := classRepo.UpdateTemplate(tpl); err != nil {
		return fmt.Errorf("更新模板失败: %w", err)
	}

	logger.Infof("更新模板 %d 成功", templateID)
	return nil
}

// DeleteTemplate 删除模板
func DeleteTemplate(templateID, operatorID uint) error {
	tpl, err := classRepo.FindTemplate(templateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模板不存在")
		}
		return fmt.Errorf("查询模板失败: %w", err)
	}

	if tpl.CreatorID != operatorID {
		return errors.New("无权删除此模板")
	}

	if err := classRepo.DeleteTemplate(templateID); err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}

	logger.Infof("删除模板 %d 成功", templateID)
	return nil
}

// CreateClassFromTemplate 从模板创建班级
func CreateClassFromTemplate(templateID, operatorID uint, req *dto.CreateClassFromTemplateRequest) error {
	tpl, err := classRepo.FindTemplate(templateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模板不存在")
		}
		return fmt.Errorf("查询模板失败: %w", err)
	}

	code := generateClassCode()
	class := &model.Class{
		Name:        req.Name,
		Description: tpl.Description,
		Code:        code,
		CreatorID:   operatorID,
		MemberCount: 1,
		Status:      1,
	}

	if err := classRepo.Create(class); err != nil {
		return fmt.Errorf("创建班级失败: %w", err)
	}

	// 创建者自动成为管理员
	member := &model.ClassMember{
		ClassID:  class.ID,
		UserID:   operatorID,
		Role:     3,
		JoinedAt: time.Now(),
	}
	if err := classRepo.CreateMember(member); err != nil {
		logger.Errorf("创建班级成员失败: %v", err)
	}

	// 增加模板使用次数
	classRepo.IncrementTemplateUsedCount(templateID)

	logger.Infof("从模板 %d 创建班级成功: %s", templateID, class.Name)
	return nil
}

// ==================== Discussion Management ====================

// ListDiscussions 讨论列表
func ListDiscussions(classID uint, req *dto.PageRequest) ([]dto.DiscussionInfo, int64, error) {
	discussions, total, err := classRepo.ListDiscussions(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询讨论列表失败: %w", err)
	}

	list := make([]dto.DiscussionInfo, 0, len(discussions))
	for _, d := range discussions {
		creatorName := ""
		if d.Creator.ID > 0 {
			creatorName = d.Creator.Nickname
		}
		list = append(list, dto.DiscussionInfo{
			ID:          d.ID,
			ClassID:     d.ClassID,
			Title:       d.Title,
			Content:     d.Content,
			CreatorID:   d.CreatorID,
			CreatorName: creatorName,
			ReplyCount:  d.ReplyCount,
			ViewCount:   d.ViewCount,
			IsTop:       d.IsTop,
			IsClosed:    d.IsClosed,
			LastReplyAt: d.LastReplyAt,
			CreatedAt:   d.CreatedAt,
		})
	}
	return list, total, nil
}

// CreateDiscussion 发布讨论
func CreateDiscussion(classID, userID uint, req *dto.CreateDiscussionRequest) error {
	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	d := &model.ClassDiscussion{
		ClassID:   classID,
		Title:     req.Title,
		Content:   req.Content,
		CreatorID: userID,
	}

	if err := classRepo.CreateDiscussion(d); err != nil {
		return fmt.Errorf("发布讨论失败: %w", err)
	}

	logger.Infof("用户 %d 发布讨论: %s", userID, d.Title)
	return nil
}

// GetDiscussion 讨论详情
func GetDiscussion(classID, discussionID uint) (*dto.DiscussionInfo, error) {
	d, err := classRepo.FindDiscussion(discussionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("讨论不存在")
		}
		return nil, fmt.Errorf("查询讨论失败: %w", err)
	}

	if d.ClassID != classID {
		return nil, errors.New("讨论不属于此班级")
	}

	// 增加浏览数
	classRepo.IncrementDiscussionViewCount(discussionID)

	creatorName := ""
	if d.Creator.ID > 0 {
		creatorName = d.Creator.Nickname
	}

	return &dto.DiscussionInfo{
		ID:          d.ID,
		ClassID:     d.ClassID,
		Title:       d.Title,
		Content:     d.Content,
		CreatorID:   d.CreatorID,
		CreatorName: creatorName,
		ReplyCount:  d.ReplyCount,
		ViewCount:   d.ViewCount + 1,
		IsTop:       d.IsTop,
		IsClosed:    d.IsClosed,
		LastReplyAt: d.LastReplyAt,
		CreatedAt:   d.CreatedAt,
	}, nil
}

// UpdateDiscussion 编辑讨论
func UpdateDiscussion(classID, discussionID, userID uint, req *dto.UpdateDiscussionRequest) error {
	d, err := classRepo.FindDiscussion(discussionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("讨论不存在")
		}
		return fmt.Errorf("查询讨论失败: %w", err)
	}

	if d.ClassID != classID {
		return errors.New("讨论不属于此班级")
	}

	// 只有创建者或管理员才能编辑
	if d.CreatorID != userID {
		if err := checkClassPermission(classID, userID, 2); err != nil {
			return errors.New("无权编辑此讨论")
		}
	}

	if req.Title != "" {
		d.Title = req.Title
	}
	if req.Content != "" {
		d.Content = req.Content
	}

	if err := classRepo.UpdateDiscussion(d); err != nil {
		return fmt.Errorf("编辑讨论失败: %w", err)
	}

	logger.Infof("编辑讨论 %d 成功", discussionID)
	return nil
}

// DeleteDiscussion 删除讨论
func DeleteDiscussion(classID, discussionID, userID uint) error {
	d, err := classRepo.FindDiscussion(discussionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("讨论不存在")
		}
		return fmt.Errorf("查询讨论失败: %w", err)
	}

	if d.ClassID != classID {
		return errors.New("讨论不属于此班级")
	}

	// 只有创建者或管理员才能删除
	if d.CreatorID != userID {
		if err := checkClassPermission(classID, userID, 2); err != nil {
			return errors.New("无权删除此讨论")
		}
	}

	if err := classRepo.DeleteDiscussion(discussionID); err != nil {
		return fmt.Errorf("删除讨论失败: %w", err)
	}

	logger.Infof("删除讨论 %d 成功", discussionID)
	return nil
}

// ToggleDiscussionPin 置顶/取消置顶
func ToggleDiscussionPin(classID, discussionID, userID uint) error {
	if err := checkClassPermission(classID, userID, 2); err != nil {
		return err
	}

	d, err := classRepo.FindDiscussion(discussionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("讨论不存在")
		}
		return fmt.Errorf("查询讨论失败: %w", err)
	}

	if d.ClassID != classID {
		return errors.New("讨论不属于此班级")
	}

	d.IsTop = !d.IsTop
	if err := classRepo.UpdateDiscussion(d); err != nil {
		return fmt.Errorf("更新讨论置顶状态失败: %w", err)
	}

	logger.Infof("讨论 %d 置顶状态: %v", discussionID, d.IsTop)
	return nil
}

// ==================== Resource Management ====================

// ListResources 资源列表
func ListResources(classID uint, req *dto.PageRequest) ([]dto.ResourceInfo, int64, error) {
	resources, total, err := classRepo.ListClassResources(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询资源列表失败: %w", err)
	}

	list := make([]dto.ResourceInfo, 0, len(resources))
	for _, r := range resources {
		list = append(list, dto.ResourceInfo{
			ID:           r.ID,
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
			Title:        r.Title,
			Version:      r.Version,
			CreatorID:    r.CreatorID,
			CreatedAt:    r.CreatedAt,
		})
	}
	return list, total, nil
}

// GetResourceStatistics 资源统计
func GetResourceStatistics(classID uint) (*dto.ResourceStatistics, error) {
	return &dto.ResourceStatistics{
		FileCount:     classRepo.CountClassFiles(classID),
		HomeworkCount: classRepo.CountClassHomework(classID),
		NoticeCount:   classRepo.CountClassNotices(classID),
		TotalSize:     classRepo.SumClassFileSize(classID),
	}, nil
}

// ImportResource 导入资源
func ImportResource(classID, operatorID uint, req *dto.ImportResourceRequest) error {
	if err := checkClassPermission(classID, operatorID, 2); err != nil {
		return err
	}

	r := &model.ClassResourceVersion{
		ResourceType: req.ResourceType,
		Title:        "批量导入",
		Content:      req.Data,
		CreatorID:    operatorID,
		Version:      1,
	}

	if err := classRepo.CreateClassResource(r); err != nil {
		return fmt.Errorf("导入资源失败: %w", err)
	}

	logger.Infof("用户 %d 导入资源到班级 %d", operatorID, classID)
	return nil
}

// ExportResource 导出资源
func ExportResource(classID uint) ([]dto.ResourceInfo, error) {
	resources, _, err := classRepo.ListClassResources(classID, &dto.PageRequest{Page: 1, PageSize: 10000})
	if err != nil {
		return nil, fmt.Errorf("查询资源失败: %w", err)
	}

	list := make([]dto.ResourceInfo, 0, len(resources))
	for _, r := range resources {
		list = append(list, dto.ResourceInfo{
			ID:           r.ID,
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
			Title:        r.Title,
			Version:      r.Version,
			CreatorID:    r.CreatorID,
			CreatedAt:    r.CreatedAt,
		})
	}
	return list, nil
}

// ==================== Management Enhancement ====================

// ArchiveClass 归档班级
func ArchiveClass(classID, operatorID uint) error {
	class, err := classRepo.FindByID(classID)
	if err != nil {
		return errors.New("班级不存在")
	}

	if class.CreatorID != operatorID {
		member, err := classRepo.FindMember(classID, operatorID)
		if err != nil || member.Role < 3 {
			return errors.New("只有管理员才能归档班级")
		}
	}

	class.Status = 2 // 归档
	if err := classRepo.Update(class); err != nil {
		return fmt.Errorf("归档班级失败: %w", err)
	}

	logger.Infof("归档班级 %d 成功", classID)
	return nil
}

// UnarchiveClass 取消归档
func UnarchiveClass(classID, operatorID uint) error {
	class, err := classRepo.FindByID(classID)
	if err != nil {
		return errors.New("班级不存在")
	}

	if class.CreatorID != operatorID {
		member, err := classRepo.FindMember(classID, operatorID)
		if err != nil || member.Role < 3 {
			return errors.New("只有管理员才能取消归档")
		}
	}

	class.Status = 1
	if err := classRepo.Update(class); err != nil {
		return fmt.Errorf("取消归档失败: %w", err)
	}

	logger.Infof("取消归档班级 %d 成功", classID)
	return nil
}

// PinClass 置顶班级
func PinClass(classID, userID uint) error {
	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	if err := database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND user_id = ?", classID, userID).
		Update("is_pinned", true).Error; err != nil {
		return fmt.Errorf("置顶班级失败: %w", err)
	}

	logger.Infof("用户 %d 置顶班级 %d", userID, classID)
	return nil
}

// UnpinClass 取消置顶
func UnpinClass(classID, userID uint) error {
	// 检查是否是班级成员
	if _, err := classRepo.FindMember(classID, userID); err != nil {
		return errors.New("不是班级成员")
	}

	if err := database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND user_id = ?", classID, userID).
		Update("is_pinned", false).Error; err != nil {
		return fmt.Errorf("取消置顶失败: %w", err)
	}

	logger.Infof("用户 %d 取消置顶班级 %d", userID, classID)
	return nil
}

// SearchClasses 搜索班级
func SearchClasses(keyword string, req *dto.PageRequest) ([]dto.ClassInfo, int64, error) {
	classes, total, err := classRepo.SearchClasses(keyword, req)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索班级失败: %w", err)
	}

	list := make([]dto.ClassInfo, 0, len(classes))
	for _, c := range classes {
		list = append(list, toClassInfo(&c))
	}
	return list, total, nil
}

// GenerateQRCode 生成二维码
func GenerateQRCode(classID uint) (*dto.QRCodeInfo, error) {
	class, err := classRepo.FindByID(classID)
	if err != nil {
		return nil, errors.New("班级不存在")
	}

	return &dto.QRCodeInfo{
		Code:  class.Code,
		URL:   fmt.Sprintf("/class/join?code=%s", class.Code),
		Image: "", // 实际实现需要生成二维码图片的 base64
	}, nil
}

// SetClassExpire 设置有效期
func SetClassExpire(classID, operatorID uint, req *dto.SetExpireRequest) error {
	class, err := classRepo.FindByID(classID)
	if err != nil {
		return errors.New("班级不存在")
	}

	if class.CreatorID != operatorID {
		member, err := classRepo.FindMember(classID, operatorID)
		if err != nil || member.Role < 3 {
			return errors.New("只有管理员才能设置有效期")
		}
	}

	if err := database.DB.Model(&model.Class{}).
		Where("id = ?", classID).
		Update("expires_at", req.ExpireAt).Error; err != nil {
		return fmt.Errorf("设置有效期失败: %w", err)
	}

	logger.Infof("设置班级 %d 有效期: %v", classID, req.ExpireAt)
	return nil
}

// ListClassExams 班级考试列表
func ListClassExams(classID uint, req *dto.PageRequest) ([]dto.ClassExamInfo, int64, error) {
	exams, total, err := classRepo.ListClassExams(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询班级考试失败: %w", err)
	}

	list := make([]dto.ClassExamInfo, 0, len(exams))
	for _, e := range exams {
		list = append(list, dto.ClassExamInfo{
			ID:        e.ID,
			Title:     e.Title,
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
			Status:    e.Status,
		})
	}
	return list, total, nil
}

// GetClassNotice 班级公告
func GetClassNotice(classID uint) (*dto.ClassNoticeDetailInfo, error) {
	notice, err := classRepo.FindLatestNotice(classID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询公告失败: %w", err)
	}

	creatorName := ""
	if user, err := userRepo.FindByID(notice.CreatorID); err == nil {
		creatorName = user.Nickname
	}

	return &dto.ClassNoticeDetailInfo{
		ID:          notice.ID,
		ClassID:     notice.ClassID,
		Title:       notice.Title,
		Content:     notice.Content,
		CreatorID:   notice.CreatorID,
		CreatorName: creatorName,
		CreatedAt:   notice.CreatedAt,
	}, nil
}

// ==================== Helper Functions ====================

// checkClassPermission 检查班级权限
// role: 1=学生, 2=教师, 3=管理员
func checkClassPermission(classID, userID uint, requiredRole int8) error {
	member, err := classRepo.FindMember(classID, userID)
	if err != nil {
		return errors.New("不是班级成员")
	}
	if member.Role < requiredRole {
		return errors.New("权限不足")
	}
	return nil
}

func marshalAttachments(attachments []string) string {
	if len(attachments) == 0 {
		return "[]"
	}
	result := "["
	for i, a := range attachments {
		if i > 0 {
			result += ","
		}
		result += `"` + a + `"`
	}
	result += "]"
	return result
}

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
