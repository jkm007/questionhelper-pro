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
	"questionhelper-server/pkg/database"
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
func UpdateClass(id, operatorID uint, isAdmin bool, req *dto.UpdateClassRequest) error {
	class, err := classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("班级不存在")
		}
		return fmt.Errorf("查询班级失败: %w", err)
	}

	// 系统管理员可跳过权限检查
	if !isAdmin {
		// 检查权限：必须是班级创建者或班级管理员（Role >= 2）
		if class.CreatorID != operatorID {
			member, err := classRepo.FindMember(id, operatorID)
			if err != nil || member.Role < 2 {
				return errors.New("无权操作此班级")
			}
		}
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
func DeleteClass(id, operatorID uint, isAdmin bool) error {
	class, err := classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("班级不存在")
		}
		return fmt.Errorf("查询班级失败: %w", err)
	}

	// 系统管理员可跳过权限检查；否则只有创建者才能删除班级
	if !isAdmin && class.CreatorID != operatorID {
		return errors.New("只有创建者才能删除班级")
	}

	// 保护规则1: 检查班级成员数（除创建者外）
	var memberCount int64
	if err := database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND user_id != ?", id, class.CreatorID).
		Count(&memberCount).Error; err != nil {
		return fmt.Errorf("查询班级成员失败: %w", err)
	}
	if memberCount > 0 {
		return fmt.Errorf("班级仍有 %d 名成员，请先移除所有成员", memberCount)
	}

	// 保护规则2: 检查是否有进行中的考试
	var activeExamCount int64
	if err := database.DB.Model(&model.Exam{}).
		Where("class_id = ? AND status = 1", id).
		Count(&activeExamCount).Error; err != nil {
		return fmt.Errorf("查询班级考试失败: %w", err)
	}
	if activeExamCount > 0 {
		return errors.New("班级有进行中的考试，请先结束所有考试")
	}

	// 保护规则3: 检查是否有考试记录
	var examRecordCount int64
	if err := database.DB.Model(&model.ExamRecord{}).
		Where("exam_id IN (SELECT id FROM exams WHERE class_id = ?)", id).
		Count(&examRecordCount).Error; err != nil {
		return fmt.Errorf("查询考试记录失败: %w", err)
	}
	if examRecordCount > 0 {
		return errors.New("班级有考试记录，请先清理考试记录")
	}

	// 保护规则4: 检查是否有班级资源
	var resourceCount int64
	if err := database.DB.Model(&model.ClassResourceVersion{}).
		Where("resource_id = ?", id).
		Count(&resourceCount).Error; err != nil {
		return fmt.Errorf("查询班级资源失败: %w", err)
	}
	if resourceCount > 0 {
		return errors.New("班级有教学资源，请先清理班级资源")
	}

	// 保护规则5: 检查是否有未完成的作业
	var incompleteHwCount int64
	if err := database.DB.Model(&model.Homework{}).
		Where("class_id = ? AND deadline > ?", id, time.Now()).
		Count(&incompleteHwCount).Error; err != nil {
		return fmt.Errorf("查询作业失败: %w", err)
	}
	if incompleteHwCount > 0 {
		return errors.New("班级有未过截止时间的作业，请先处理作业")
	}

	if err := classRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除班级失败: %w", err)
	}

	logger.Infof("删除班级 %d 成功", id)
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
