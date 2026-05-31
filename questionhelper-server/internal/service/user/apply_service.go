package user

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	appRepo "questionhelper-server/internal/repository/application"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/logger"
)

// CreateApplication 创建角色申请
func CreateApplication(userID uint, req *dto.CreateApplicationRequest) error {
	// 验证用户存在
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证角色存在
	_, err = userRepo.FindRoleByID(req.RoleID)
	if err != nil {
		return errors.New("申请的角色不存在")
	}

	// 前置条件校验
	// 1. 检查是否已实名认证
	if !user.IsReal {
		return errors.New("请先完成实名认证")
	}

	// 2. 检查账号状态
	if user.Status != 1 {
		return errors.New("账号状态异常，无法申请")
	}

	// 3. 检查是否有待审核的申请
	pending, err := appRepo.FindPendingByUserID(userID)
	if err == nil && pending.ID > 0 {
		return errors.New("您有正在审核中的申请，请等待审核完成")
	}

	// 4. 检查驳回后24小时冷静期
	rejected, err := appRepo.FindRecentRejectedByUserID(userID, 24*time.Hour)
	if err == nil && rejected.ID > 0 {
		return errors.New("申请被驳回后需等待24小时才能重新申请")
	}

	// 创建申请
	app := &model.RoleApplication{
		UserID: userID,
		RoleID: req.RoleID,
		Reason: req.Reason,
		Status: 0,
	}

	if err := appRepo.Create(app); err != nil {
		return fmt.Errorf("创建角色申请失败: %w", err)
	}

	logger.Infof("用户 %d 创建角色申请成功，申请角色: %d", userID, req.RoleID)
	return nil
}

// ReviewApplication 审核角色申请
func ReviewApplication(id uint, reviewerID uint, req *dto.ReviewApplicationRequest) error {
	app, err := appRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("申请记录不存在")
		}
		return fmt.Errorf("查询申请失败: %w", err)
	}

	if app.Status != 0 {
		return errors.New("该申请已处理")
	}

	now := time.Now()
	app.Status = req.Status
	app.ReviewNote = req.Note
	app.ReviewedBy = &reviewerID
	app.ReviewedAt = &now

	if err := appRepo.Update(app); err != nil {
		return fmt.Errorf("审核申请失败: %w", err)
	}

	// 审核通过，给用户分配角色
	if req.Status == 1 {
		user, err := userRepo.FindByID(app.UserID)
		if err != nil {
			return fmt.Errorf("查询用户失败: %w", err)
		}

		// 添加角色
		roles := user.Roles
		roles = append(roles, model.Role{ID: app.RoleID})
		if err := userRepo.Update(user); err != nil {
			return fmt.Errorf("分配角色失败: %w", err)
		}

		logger.Infof("角色申请 %d 审核通过，用户 %d 获得角色 %d", id, app.UserID, app.RoleID)
	} else {
		logger.Infof("角色申请 %d 审核驳回，原因: %s", id, req.Note)
	}

	return nil
}

// GetApplication 获取申请详情
func GetApplication(id uint) (*dto.ApplicationInfo, error) {
	app, err := appRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("申请记录不存在")
		}
		return nil, fmt.Errorf("查询申请失败: %w", err)
	}
	return toApplicationInfo(app), nil
}

// ListApplications 申请列表
func ListApplications(req *dto.ApplicationListRequest) ([]dto.ApplicationInfo, int64, error) {
	apps, total, err := appRepo.List(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询申请列表失败: %w", err)
	}

	list := make([]dto.ApplicationInfo, 0, len(apps))
	for _, app := range apps {
		list = append(list, *toApplicationInfo(&app))
	}

	return list, total, nil
}

// GetUserApplications 获取用户的申请列表
func GetUserApplications(userID uint, page, pageSize int) ([]dto.ApplicationInfo, int64, error) {
	req := &dto.ApplicationListRequest{
		PageRequest: dto.PageRequest{
			Page:     page,
			PageSize: pageSize,
		},
		UserID: &userID,
	}
	return ListApplications(req)
}

// toApplicationInfo 转换为 ApplicationInfo DTO
func toApplicationInfo(app *model.RoleApplication) *dto.ApplicationInfo {
	info := &dto.ApplicationInfo{
		ID:         app.ID,
		UserID:     app.UserID,
		RoleID:     app.RoleID,
		Reason:     app.Reason,
		Status:     app.Status,
		ReviewNote: app.ReviewNote,
		ReviewedBy: app.ReviewedBy,
		ReviewedAt: app.ReviewedAt,
		CreatedAt:  app.CreatedAt,
	}

	if app.User.ID > 0 {
		info.Username = app.User.Username
		info.Nickname = app.User.Nickname
	}
	if app.Role.ID > 0 {
		info.RoleName = app.Role.Name
	}
	if app.Reviewer.ID > 0 {
		info.Reviewer = app.Reviewer.Nickname
	}

	return info
}
