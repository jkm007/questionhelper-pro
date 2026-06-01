package class

import (
	"errors"
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	classRepo "questionhelper-server/internal/repository/class"
	"questionhelper-server/pkg/logger"
)

// ==================== Member Management ====================

// JoinClass 加入班级
func JoinClass(userID uint, req *dto.JoinClassRequest) error {
	class, err := classRepo.FindByCode(req.Code)
	if err != nil {
		return errors.New("加入码无效")
	}

	// 检查是否已是成员
	_, err = classRepo.FindMember(class.ID, userID)
	if err == nil {
		return errors.New("已是班级成员")
	}

	// 检查班级状态
	if class.Status != 1 {
		return errors.New("班级已归档或不可用")
	}

	// 检查班级是否满员
	if class.MemberCount >= 500 {
		return errors.New("班级已满员")
	}

	member := &model.ClassMember{
		ClassID:  class.ID,
		UserID:   userID,
		Role:     1,
		JoinedAt: time.Now(),
	}

	if err := classRepo.CreateMember(member); err != nil {
		return fmt.Errorf("加入班级失败: %w", err)
	}

	// 更新成员数
	classRepo.IncrementMemberCount(class.ID)

	logger.Infof("用户 %d 加入班级 %d", userID, class.ID)
	return nil
}

// LeaveClass 离开班级
func LeaveClass(classID, userID uint) error {
	_, err := classRepo.FindMember(classID, userID)
	if err != nil {
		return errors.New("不是班级成员")
	}

	if err := classRepo.DeleteMember(classID, userID); err != nil {
		return fmt.Errorf("离开班级失败: %w", err)
	}

	// 更新成员数
	classRepo.DecrementMemberCount(classID)

	logger.Infof("用户 %d 离开班级 %d", userID, classID)
	return nil
}

// ListMembers 成员列表
func ListMembers(classID uint, req *dto.ClassMemberListRequest) ([]dto.ClassMemberInfo, int64, error) {
	members, total, err := classRepo.ListMembers(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询成员列表失败: %w", err)
	}

	list := make([]dto.ClassMemberInfo, 0, len(members))
	for _, m := range members {
		list = append(list, dto.ClassMemberInfo{
			ID:       m.ID,
			ClassID:  m.ClassID,
			UserID:   m.UserID,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}
	return list, total, nil
}

// RemoveMember 移除成员
func RemoveMember(classID, userID, operatorID uint) error {
	// 检查被移除者是否是班级成员
	_, err := classRepo.FindMember(classID, userID)
	if err != nil {
		return errors.New("不是班级成员")
	}

	// 检查操作者权限
	operator, err := classRepo.FindMember(classID, operatorID)
	if err != nil {
		return errors.New("无权操作此班级")
	}
	if operator.Role < 2 {
		return errors.New("无权移除成员")
	}

	// 不允许移除创建者（Role == 3）
	target, _ := classRepo.FindMember(classID, userID)
	if target != nil && target.Role == 3 {
		return errors.New("不能移除班级创建者")
	}

	if err := classRepo.DeleteMember(classID, userID); err != nil {
		return fmt.Errorf("移除成员失败: %w", err)
	}

	classRepo.DecrementMemberCount(classID)
	logger.Infof("移除班级 %d 成员 %d", classID, userID)
	return nil
}
