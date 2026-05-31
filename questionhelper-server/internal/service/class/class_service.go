package class

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	classRepo "questionhelper-server/internal/repository/class"
	"questionhelper-server/pkg/logger"
)

// ==================== Class ====================

// ListClasses 班级列表
func ListClasses(req *dto.ClassListRequest) ([]dto.ClassInfo, int64, error) {
	classes, total, err := classRepo.List(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询班级列表失败: %w", err)
	}

	list := make([]dto.ClassInfo, 0, len(classes))
	for _, c := range classes {
		list = append(list, toClassInfo(&c))
	}
	return list, total, nil
}

// ListMyClasses 我的班级列表
func ListMyClasses(userID uint, req *dto.PageRequest) ([]dto.ClassInfo, int64, error) {
	classes, total, err := classRepo.ListByUser(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询我的班级失败: %w", err)
	}

	list := make([]dto.ClassInfo, 0, len(classes))
	for _, c := range classes {
		list = append(list, toClassInfo(&c))
	}
	return list, total, nil
}

// GetClass 获取班级详情
func GetClass(id uint) (*dto.ClassInfo, error) {
	class, err := classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("班级不存在")
		}
		return nil, fmt.Errorf("查询班级失败: %w", err)
	}
	info := toClassInfo(class)
	return &info, nil
}

// CreateClass 创建班级
func CreateClass(creatorID uint, req *dto.CreateClassRequest) error {
	code := generateClassCode()

	class := &model.Class{
		Name:        req.Name,
		Description: req.Description,
		Cover:       req.Cover,
		Code:        code,
		CreatorID:   creatorID,
		MemberCount: 1,
		Status:      1,
	}

	if err := classRepo.Create(class); err != nil {
		return fmt.Errorf("创建班级失败: %w", err)
	}

	// 创建者自动成为管理员
	member := &model.ClassMember{
		ClassID:  class.ID,
		UserID:   creatorID,
		Role:     3,
		JoinedAt: time.Now(),
	}
	if err := classRepo.CreateMember(member); err != nil {
		logger.Errorf("创建班级成员失败: %v", err)
	}

	logger.Infof("创建班级成功: %s (加入码: %s)", class.Name, code)
	return nil
}

// UpdateClass 更新班级
func UpdateClass(id uint, req *dto.UpdateClassRequest) error {
	class, err := classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("班级不存在")
		}
		return fmt.Errorf("查询班级失败: %w", err)
	}

	if req.Name != "" {
		class.Name = req.Name
	}
	if req.Description != "" {
		class.Description = req.Description
	}
	if req.Cover != "" {
		class.Cover = req.Cover
	}

	if err := classRepo.Update(class); err != nil {
		return fmt.Errorf("更新班级失败: %w", err)
	}

	logger.Infof("更新班级 %d 成功", id)
	return nil
}

// DeleteClass 删除班级
func DeleteClass(id uint) error {
	_, err := classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("班级不存在")
		}
		return fmt.Errorf("查询班级失败: %w", err)
	}

	if err := classRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除班级失败: %w", err)
	}

	logger.Infof("删除班级 %d 成功", id)
	return nil
}

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
func RemoveMember(classID, userID uint) error {
	_, err := classRepo.FindMember(classID, userID)
	if err != nil {
		return errors.New("不是班级成员")
	}

	if err := classRepo.DeleteMember(classID, userID); err != nil {
		return fmt.Errorf("移除成员失败: %w", err)
	}

	classRepo.DecrementMemberCount(classID)
	logger.Infof("移除班级 %d 成员 %d", classID, userID)
	return nil
}

// ==================== Homework ====================

// ListHomework 作业列表
func ListHomework(classID uint, req *dto.PageRequest) ([]dto.HomeworkInfo, int64, error) {
	homework, total, err := classRepo.ListHomework(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询作业列表失败: %w", err)
	}

	list := make([]dto.HomeworkInfo, 0, len(homework))
	for _, h := range homework {
		list = append(list, dto.HomeworkInfo{
			ID:          h.ID,
			ClassID:     h.ClassID,
			Title:       h.Title,
			Description: h.Description,
			Deadline:    h.Deadline,
			CreatorID:   h.CreatorID,
		})
	}
	return list, total, nil
}

// CreateHomework 创建作业
func CreateHomework(classID, creatorID uint, req *dto.CreateHomeworkRequest) error {
	hw := &model.Homework{
		ClassID:     classID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		CreatorID:   creatorID,
	}

	if err := classRepo.CreateHomework(hw); err != nil {
		return fmt.Errorf("创建作业失败: %w", err)
	}

	logger.Infof("创建作业成功: %d", hw.ID)
	return nil
}

// ==================== Notice ====================

// ListNotices 公告列表
func ListNotices(classID uint, req *dto.PageRequest) ([]dto.ClassNoticeInfo, int64, error) {
	notices, total, err := classRepo.ListNotices(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询公告列表失败: %w", err)
	}

	list := make([]dto.ClassNoticeInfo, 0, len(notices))
	for _, n := range notices {
		list = append(list, dto.ClassNoticeInfo{
			ID:        n.ID,
			ClassID:   n.ClassID,
			Title:     n.Title,
			Content:   n.Content,
			CreatorID: n.CreatorID,
			CreatedAt: n.CreatedAt,
		})
	}
	return list, total, nil
}

// CreateNotice 创建公告
func CreateNotice(classID, creatorID uint, req *dto.CreateNoticeRequest) error {
	notice := &model.ClassNotice{
		ClassID:   classID,
		Title:     req.Title,
		Content:   req.Content,
		CreatorID: creatorID,
	}

	if err := classRepo.CreateNotice(notice); err != nil {
		return fmt.Errorf("创建公告失败: %w", err)
	}

	logger.Infof("创建公告成功: %d", notice.ID)
	return nil
}

// ==================== Helpers ====================

func toClassInfo(c *model.Class) dto.ClassInfo {
	info := dto.ClassInfo{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Cover:       c.Cover,
		Code:        c.Code,
		CreatorID:   c.CreatorID,
		MemberCount: c.MemberCount,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
	}
	if c.Creator.ID > 0 {
		info.CreatorName = c.Creator.Nickname
	}
	return info
}

func generateClassCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
