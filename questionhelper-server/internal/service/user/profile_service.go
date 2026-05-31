package user

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/encrypt"
	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/validator"
)

// GetPrivacy 获取用户隐私设置
func GetPrivacy(userID uint) (*dto.PrivacyInfo, error) {
	privacy, err := userRepo.FindPrivacyByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回默认隐私设置
			return &dto.PrivacyInfo{
				ProfileVisible:  1,
				RealnameVisible: 1,
				EmailVisible:    1,
				StatsVisible:    1,
				ClassVisible:    1,
			}, nil
		}
		return nil, fmt.Errorf("查询隐私设置失败: %w", err)
	}

	return &dto.PrivacyInfo{
		ProfileVisible:  privacy.ProfileVisible,
		RealnameVisible: privacy.RealnameVisible,
		EmailVisible:    privacy.EmailVisible,
		StatsVisible:    privacy.StatsVisible,
		ClassVisible:    privacy.ClassVisible,
	}, nil
}

// UpdatePrivacy 更新用户隐私设置
func UpdatePrivacy(userID uint, req *dto.UpdatePrivacyRequest) error {
	privacy, err := userRepo.FindPrivacyByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建隐私设置
			privacy = &model.UserPrivacy{
				UserID: userID,
			}
		} else {
			return fmt.Errorf("查询隐私设置失败: %w", err)
		}
	}

	if req.ProfileVisible != nil {
		privacy.ProfileVisible = *req.ProfileVisible
	}
	if req.RealnameVisible != nil {
		privacy.RealnameVisible = *req.RealnameVisible
	}
	if req.EmailVisible != nil {
		privacy.EmailVisible = *req.EmailVisible
	}
	if req.StatsVisible != nil {
		privacy.StatsVisible = *req.StatsVisible
	}
	if req.ClassVisible != nil {
		privacy.ClassVisible = *req.ClassVisible
	}

	if privacy.ID == 0 {
		if err := userRepo.CreatePrivacy(privacy); err != nil {
			return fmt.Errorf("创建隐私设置失败: %w", err)
		}
	} else {
		if err := userRepo.UpdatePrivacy(privacy); err != nil {
			return fmt.Errorf("更新隐私设置失败: %w", err)
		}
	}

	logger.Infof("用户 %d 更新隐私设置成功", userID)
	return nil
}

// BindPhone 绑定手机号
func BindPhone(userID uint, req *dto.BindPhoneRequest) error {
	// TODO: 验证短信验证码

	// 检查手机号是否已被使用
	exists, err := userRepo.ExistsByPhone(req.Phone)
	if err != nil {
		return fmt.Errorf("检查手机号失败: %w", err)
	}
	if exists {
		return errors.New("手机号已被其他账号绑定")
	}

	if err := userRepo.UpdateByID(userID, map[string]interface{}{
		"phone": req.Phone,
	}); err != nil {
		return fmt.Errorf("绑定手机号失败: %w", err)
	}

	logger.Infof("用户 %d 绑定手机号成功", userID)
	return nil
}

// BindEmail 绑定邮箱
func BindEmail(userID uint, req *dto.BindEmailRequest) error {
	// TODO: 验证邮箱验证码

	// 检查邮箱是否已被使用
	exists, err := userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return errors.New("邮箱已被其他账号绑定")
	}

	if err := userRepo.UpdateByID(userID, map[string]interface{}{
		"email": req.Email,
	}); err != nil {
		return fmt.Errorf("绑定邮箱失败: %w", err)
	}

	logger.Infof("用户 %d 绑定邮箱成功", userID)
	return nil
}

// SubmitRealName 提交实名认证
func SubmitRealName(userID uint, req *dto.RealNameSubmitRequest) error {
	// 验证身份证号格式
	if err := validator.ValidateIDCard(req.IDCard); err != nil {
		return fmt.Errorf("身份证号格式错误: %w", err)
	}

	// 检查是否已提交过
	existing, err := FindRealNameByUserID(userID)
	if err == nil && existing != nil {
		if existing.Status == 1 {
			return errors.New("已通过实名认证")
		}
		if existing.Status == 0 {
			return errors.New("实名认证正在审核中")
		}
	}

	// 检查身份证号是否已被使用(通过哈希查询)
	// TODO: 实现身份证号哈希查询

	// AES加密身份证号
	encryptedIDCard, err := encrypt.AESEncrypt(req.IDCard)
	if err != nil {
		return fmt.Errorf("加密身份证号失败: %w", err)
	}

	// 创建实名认证记录
	realName := &model.UserRealName{
		UserID:     userID,
		RealName:   req.RealName,
		IDCard:     encryptedIDCard,
		Status:     0,
		SubmittedAt: time.Now(),
	}

	if err := CreateRealName(realName); err != nil {
		return fmt.Errorf("提交实名认证失败: %w", err)
	}

	logger.Infof("用户 %d 提交实名认证成功", userID)
	return nil
}

// GetRealName 获取实名认证信息
func GetRealName(userID uint) (*dto.RealNameInfo, error) {
	realName, err := FindRealNameByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询实名认证失败: %w", err)
	}

	return &dto.RealNameInfo{
		ID:           realName.ID,
		RealName:     realName.RealName,
		IDCardMasked: encrypt.MaskIDCard(realName.IDCard),
		Status:       realName.Status,
		RejectReason: realName.RejectReason,
		SubmittedAt:  realName.SubmittedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ReviewRealName 审核实名认证(管理员)
func ReviewRealName(id uint, reviewerID uint, req *dto.ReviewRealNameRequest) error {
	realName, err := FindRealNameByID(id)
	if err != nil {
		return fmt.Errorf("查询实名认证失败: %w", err)
	}

	if realName.Status != 0 {
		return errors.New("该认证已处理")
	}

	now := time.Now()
	realName.Status = req.Status
	realName.ReviewedBy = &reviewerID
	realName.ReviewedAt = &now

	if req.Status == 2 {
		realName.RejectReason = req.Reason
	}

	if err := UpdateRealName(realName); err != nil {
		return fmt.Errorf("审核实名认证失败: %w", err)
	}

	// 审核通过，更新用户表的实名状态
	if req.Status == 1 {
		if err := userRepo.UpdateByID(realName.UserID, map[string]interface{}{
			"is_real":    true,
			"real_name":  realName.RealName,
		}); err != nil {
			return fmt.Errorf("更新用户实名状态失败: %w", err)
		}
	}

	logger.Infof("实名认证 %d 审核完成，状态: %d", id, req.Status)
	return nil
}

// GetOAuthBindings 获取第三方账号绑定列表
func GetOAuthBindings(userID uint) ([]dto.OAuthInfo, error) {
	oauths, err := userRepo.FindOAuthUsersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询第三方绑定失败: %w", err)
	}

	list := make([]dto.OAuthInfo, 0, len(oauths))
	for _, oauth := range oauths {
		list = append(list, dto.OAuthInfo{
			Provider:       oauth.Provider,
			ProviderType:   oauth.ProviderType,
			ProviderUserID: oauth.ProviderUserID,
			BindTime:       oauth.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return list, nil
}

// UnbindOAuth 解绑第三方账号
func UnbindOAuth(userID uint, provider string) error {
	// 检查是否是最后一个绑定方式
	oauths, err := userRepo.FindOAuthUsersByUserID(userID)
	if err != nil {
		return fmt.Errorf("查询第三方绑定失败: %w", err)
	}

	user, err := userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 如果没有密码且只有一个第三方绑定，不允许解绑
	if user.Password == "" && len(oauths) <= 1 {
		return errors.New("当前账号未设置密码且只有一个第三方绑定，无法解绑")
	}

	if err := userRepo.DeleteOAuthUser(userID, provider); err != nil {
		return fmt.Errorf("解绑第三方账号失败: %w", err)
	}

	logger.Infof("用户 %d 解绑第三方账号 %s 成功", userID, provider)
	return nil
}

// FindUserByID 查找用户(内部辅助函数)
func FindUserByID(userID uint) (*model.User, error) {
	return userRepo.FindByID(userID)
}

// 以下是实名认证相关的辅助函数，需要在 real_name_repo 中实现

func FindRealNameByUserID(userID uint) (*model.UserRealName, error) {
	var realName model.UserRealName
	err := userRepo.GetDB().Where("user_id = ?", userID).First(&realName).Error
	return &realName, err
}

func FindRealNameByID(id uint) (*model.UserRealName, error) {
	var realName model.UserRealName
	err := userRepo.GetDB().First(&realName, id).Error
	return &realName, err
}

func CreateRealName(realName *model.UserRealName) error {
	return userRepo.GetDB().Create(realName).Error
}

func UpdateRealName(realName *model.UserRealName) error {
	return userRepo.GetDB().Save(realName).Error
}
