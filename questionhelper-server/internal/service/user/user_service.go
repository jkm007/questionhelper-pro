package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/encrypt"
	"questionhelper-server/pkg/logger"
)

// GetProfile 获取用户个人信息
func GetProfile(userID uint) (*dto.UserInfo, error) {
	u, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return toUserInfo(u), nil
}

// UpdateProfile 更新个人信息
func UpdateProfile(userID uint, req *dto.UpdateProfileRequest) error {
	u, err := userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	if req.Nickname != "" {
		u.Nickname = req.Nickname
	}
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		if req.Email != u.Email {
			exists, err := userRepo.ExistsByEmail(req.Email)
			if err != nil {
				return fmt.Errorf("检查邮箱失败: %w", err)
			}
			if exists {
				return errors.New("邮箱已被使用")
			}
		}
		u.Email = req.Email
	}
	if req.Gender != 0 {
		u.Gender = req.Gender
	}
	if req.Birthday != nil {
		u.Birthday = req.Birthday
	}
	if req.Bio != "" {
		u.Bio = req.Bio
	}

	if err := userRepo.Update(u); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	logger.Infof("用户 %d 更新个人信息成功", userID)
	return nil
}

// ChangePassword 修改密码
func ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	u, err := userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证旧密码
	if !encrypt.CheckPassword(req.OldPassword, u.Password) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := encrypt.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	if err := userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return fmt.Errorf("修改密码失败: %w", err)
	}

	logger.Infof("用户 %d 修改密码成功", userID)
	return nil
}

// RealNameAuth 实名认证
func RealNameAuth(userID uint, req *dto.RealNameAuthRequest) error {
	u, err := userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	if u.IsReal {
		return errors.New("已通过实名认证")
	}

	u.RealName = req.RealName
	u.IDCard = req.IDCard
	u.IsReal = true

	if err := userRepo.Update(u); err != nil {
		return fmt.Errorf("实名认证失败: %w", err)
	}

	logger.Infof("用户 %d 实名认证成功", userID)
	return nil
}

// ListUsers 用户列表（管理员）
func ListUsers(req *dto.UserListRequest) ([]dto.UserInfo, int64, error) {
	users, total, err := userRepo.List(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	list := make([]dto.UserInfo, 0, len(users))
	for _, u := range users {
		list = append(list, *toUserInfo(&u))
	}

	return list, total, nil
}

// GetUser 获取用户详情（管理员）
func GetUser(id uint) (*dto.UserInfo, error) {
	u, err := userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return toUserInfo(u), nil
}

// CreateUser 创建用户（管理员）
func CreateUser(req *dto.CreateUserRequest) error {
	// 检查用户名
	exists, err := userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查手机号
	if req.Phone != "" {
		exists, err = userRepo.ExistsByPhone(req.Phone)
		if err != nil {
			return fmt.Errorf("检查手机号失败: %w", err)
		}
		if exists {
			return errors.New("手机号已被注册")
		}
	}

	// 检查邮箱
	if req.Email != "" {
		exists, err = userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return fmt.Errorf("检查邮箱失败: %w", err)
		}
		if exists {
			return errors.New("邮箱已被注册")
		}
	}

	// 加密密码
	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	u := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   1,
	}

	// 关联角色
	if len(req.RoleIDs) > 0 {
		roles := make([]model.Role, 0, len(req.RoleIDs))
		for _, roleID := range req.RoleIDs {
			roles = append(roles, model.Role{ID: roleID})
		}
		u.Roles = roles
	}

	if err := userRepo.Create(u); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	logger.Infof("管理员创建用户成功: %s", req.Username)
	return nil
}

// UpdateUser 更新用户（管理员）
func UpdateUser(id uint, req *dto.UpdateUserRequest) error {
	u, err := userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	if req.Nickname != "" {
		u.Nickname = req.Nickname
	}
	if req.Phone != "" {
		if req.Phone != u.Phone {
			exists, err := userRepo.ExistsByPhone(req.Phone)
			if err != nil {
				return fmt.Errorf("检查手机号失败: %w", err)
			}
			if exists {
				return errors.New("手机号已被使用")
			}
		}
		u.Phone = req.Phone
	}
	if req.Email != "" {
		if req.Email != u.Email {
			exists, err := userRepo.ExistsByEmail(req.Email)
			if err != nil {
				return fmt.Errorf("检查邮箱失败: %w", err)
			}
			if exists {
				return errors.New("邮箱已被使用")
			}
		}
		u.Email = req.Email
	}
	if req.Status != nil {
		u.Status = *req.Status
	}

	// 更新角色关联
	if req.RoleIDs != nil {
		roles := make([]model.Role, 0, len(req.RoleIDs))
		for _, roleID := range req.RoleIDs {
			roles = append(roles, model.Role{ID: roleID})
		}
		u.Roles = roles
	}

	if err := userRepo.Update(u); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	logger.Infof("管理员更新用户 %d 成功", id)
	return nil
}

// DeleteUser 删除用户（管理员）
func DeleteUser(id uint) error {
	_, err := userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	if err := userRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	logger.Infof("管理员删除用户 %d 成功", id)
	return nil
}

// UpdateUserStatus 更新用户状态（管理员）
func UpdateUserStatus(id uint, status int8) error {
	u, err := userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	u.Status = status
	if err := userRepo.Update(u); err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	logger.Infof("管理员更新用户 %d 状态为 %d", id, status)
	return nil
}

// BatchUpdateStatus 批量更新用户状态
func BatchUpdateStatus(req *dto.BatchStatusRequest) error {
	if len(req.IDs) == 0 {
		return errors.New("请选择用户")
	}

	if err := userRepo.BatchUpdateStatus(req.IDs, req.Status); err != nil {
		return fmt.Errorf("批量更新状态失败: %w", err)
	}

	logger.Infof("管理员批量更新 %d 个用户状态为 %d", len(req.IDs), req.Status)
	return nil
}

// BatchDeleteUsers 批量删除用户
func BatchDeleteUsers(req *dto.BatchDeleteRequest) error {
	if len(req.IDs) == 0 {
		return errors.New("请选择用户")
	}

	if err := userRepo.BatchSoftDelete(req.IDs); err != nil {
		return fmt.Errorf("批量删除用户失败: %w", err)
	}

	logger.Infof("管理员批量删除 %d 个用户", len(req.IDs))
	return nil
}

// BatchAssignRoles 批量分配角色
func BatchAssignRoles(req *dto.BatchRoleRequest) error {
	if len(req.IDs) == 0 {
		return errors.New("请选择用户")
	}

	for _, userID := range req.IDs {
		if err := userRepo.AssignRoles(userID, req.RoleIDs); err != nil {
			return fmt.Errorf("分配角色失败: %w", err)
		}
	}

	logger.Infof("管理员批量为 %d 个用户分配角色 %v", len(req.IDs), req.RoleIDs)
	return nil
}

// ResetPassword 重置用户密码（管理员）
func ResetPassword(id uint, req *dto.AdminResetPasswordRequest) error {
	_, err := userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 加密新密码
	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	if err := userRepo.UpdatePassword(id, hashedPassword); err != nil {
		return fmt.Errorf("重置密码失败: %w", err)
	}

	logger.Infof("管理员重置用户 %d 密码成功", id)
	return nil
}

// AssignRoles 分配用户角色
func AssignRoles(userID uint, roleIDs []uint) error {
	if err := userRepo.AssignRoles(userID, roleIDs); err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	logger.Infof("分配用户 %d 角色成功: %v", userID, roleIDs)
	return nil
}

// ListUsersWithTags 带标签筛选的用户列表
func ListUsersWithTags(req *dto.UserListRequest) ([]dto.UserInfo, int64, error) {
	users, total, err := userRepo.ListWithTags(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	list := make([]dto.UserInfo, 0, len(users))
	for _, u := range users {
		list = append(list, *toUserInfo(&u))
	}

	return list, total, nil
}

// toUserInfo 转换为 UserInfo DTO
func toUserInfo(u *model.User) *dto.UserInfo {
	info := &dto.UserInfo{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		Status:    u.Status,
		IsReal:    u.IsReal,
		CreatedAt: u.CreatedAt,
	}
	info.Roles = make([]dto.RoleInfo, 0, len(u.Roles))
	for _, role := range u.Roles {
		info.Roles = append(info.Roles, dto.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		})
	}
	return info
}
