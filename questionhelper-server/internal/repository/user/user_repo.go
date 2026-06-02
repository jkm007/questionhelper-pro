package user

import (
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return database.DB
}

// ==================== User ====================

func FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Preload("Roles").Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.Preload("Roles").First(&user, id).Error
	return &user, err
}

func Create(user *model.User) error {
	return database.DB.Create(user).Error
}

func Update(user *model.User) error {
	return database.DB.Save(user).Error
}

// UpdateByID 更新用户指定字段
func UpdateByID(id uint, data map[string]interface{}) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}

// IncrementLoginFailCount 原子递增登录失败次数
func IncrementLoginFailCount(id uint) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("login_fail_count", gorm.Expr("login_fail_count + 1")).Error
}

func UpdatePassword(id uint, password string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

func FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func List(req *dto.UserListRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.DB.Model(&model.User{})

	if req.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.RoleID != nil {
		db = db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Where("user_roles.role_id = ?", *req.RoleID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Roles").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("id DESC").
		Find(&users).Error

	return users, total, err
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.User{}, id).Error
}

func DeleteByIDs(ids []uint) error {
	return database.DB.Delete(&model.User{}, ids).Error
}

func ExistsByUsername(username string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func ExistsByPhone(phone string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

func ExistsByEmail(email string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// BatchUpdateStatus 批量更新用户状态
func BatchUpdateStatus(ids []uint, status int8) error {
	return database.DB.Model(&model.User{}).Where("id IN ?", ids).
		Update("status", status).Error
}

// BatchSoftDelete 批量软删除用户
func BatchSoftDelete(ids []uint) error {
	return database.DB.Delete(&model.User{}, ids).Error
}

// ListWithTags 带标签筛选的用户列表
func ListWithTags(req *dto.UserListRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.DB.Model(&model.User{})

	if req.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.RoleID != nil {
		db = db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Where("user_roles.role_id = ?", *req.RoleID)
	}
	if req.RoleCode != "" {
		db = db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.code = ?", req.RoleCode)
	}
	if req.TagID != nil {
		db = db.Joins("JOIN user_tags ON user_tags.user_id = users.id").
			Where("user_tags.tag_id = ?", *req.TagID)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", *req.EndDate)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Roles").Preload("Tags").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("id DESC").
		Find(&users).Error

	return users, total, err
}

// FindByTagID 按标签查找用户
func FindByTagID(tagID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.DB.Model(&model.User{}).
		Joins("JOIN user_tags ON user_tags.user_id = users.id").
		Where("user_tags.tag_id = ?", tagID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Roles").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&users).Error

	return users, total, err
}

// ExportUsers 导出用户数据
func ExportUsers(req *dto.UserListRequest) ([]model.User, error) {
	var users []model.User

	db := database.DB.Model(&model.User{})

	if req.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.RoleID != nil {
		db = db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Where("user_roles.role_id = ?", *req.RoleID)
	}
	if req.TagID != nil {
		db = db.Joins("JOIN user_tags ON user_tags.user_id = users.id").
			Where("user_tags.tag_id = ?", *req.TagID)
	}

	err := db.Preload("Roles").Order("id ASC").Find(&users).Error
	return users, err
}

// ==================== UserPrivacy ====================

// CreatePrivacy 创建用户隐私设置
func CreatePrivacy(privacy *model.UserPrivacy) error {
	return database.DB.Create(privacy).Error
}

// UpdatePrivacy 更新用户隐私设置
func UpdatePrivacy(privacy *model.UserPrivacy) error {
	return database.DB.Save(privacy).Error
}

// FindPrivacyByUserID 查找用户隐私设置
func FindPrivacyByUserID(userID uint) (*model.UserPrivacy, error) {
	var privacy model.UserPrivacy
	err := database.DB.Where("user_id = ?", userID).First(&privacy).Error
	return &privacy, err
}

// ==================== LoginLog ====================

// CreateLoginLog 创建登录日志
func CreateLoginLog(log *model.LoginLog) error {
	return database.DB.Create(log).Error
}

// ListLoginLogs 查询登录日志
func ListLoginLogs(userID *uint, page, pageSize int) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	db := database.DB.Model(&model.LoginLog{})
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// ==================== LoginDevice ====================

// CreateDevice 创建登录设备
func CreateDevice(device *model.LoginDevice) error {
	return database.DB.Create(device).Error
}

// FindDevicesByUserID 查找用户所有设备
func FindDevicesByUserID(userID uint) ([]model.LoginDevice, error) {
	var devices []model.LoginDevice
	err := database.DB.Where("user_id = ?", userID).
		Order("last_active_at DESC").
		Find(&devices).Error
	return devices, err
}

// FindLastDeviceByUserID 查找用户最近一次登录设备
func FindLastDeviceByUserID(userID uint) (*model.LoginDevice, error) {
	var device model.LoginDevice
	err := database.DB.Where("user_id = ?", userID).
		Order("last_active_at DESC").
		First(&device).Error
	return &device, err
}

// FindDeviceByUserIDAndDeviceID 通过用户ID和设备标识查找设备
func FindDeviceByUserIDAndDeviceID(userID uint, deviceID string) (*model.LoginDevice, error) {
	var device model.LoginDevice
	err := database.DB.Where("user_id = ? AND device_id = ?", userID, deviceID).
		First(&device).Error
	return &device, err
}

// FindDeviceByID 查找设备
func FindDeviceByID(deviceID uint) (*model.LoginDevice, error) {
	var device model.LoginDevice
	err := database.DB.First(&device, deviceID).Error
	return &device, err
}

// FindDeviceByJTI 通过 JTI 查找设备
func FindDeviceByJTI(jti string) (*model.LoginDevice, error) {
	var device model.LoginDevice
	err := database.DB.Where("token_jti = ?", jti).First(&device).Error
	return &device, err
}

// UpdateDeviceLastActive 更新设备最后活跃时间
func UpdateDeviceLastActive(deviceID uint) error {
	now := time.Now()
	return database.DB.Model(&model.LoginDevice{}).
		Where("id = ?", deviceID).
		Update("last_active_at", now).Error
}

// DeleteDevice 删除设备
func DeleteDevice(deviceID uint) error {
	return database.DB.Delete(&model.LoginDevice{}, deviceID).Error
}

// DeleteDevicesByUserID 删除用户所有设备
func DeleteDevicesByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.LoginDevice{}).Error
}

// DeleteInactiveDevices 删除不活跃设备
func DeleteInactiveDevices(before time.Time) (int64, error) {
	result := database.DB.Where("last_active_at < ?", before).Delete(&model.LoginDevice{})
	return result.RowsAffected, result.Error
}

// ==================== SecurityLog ====================

// CreateSecurityLog 创建安全日志
func CreateSecurityLog(log *model.SecurityLog) error {
	return database.DB.Create(log).Error
}

// ListSecurityLogs 查询安全日志
func ListSecurityLogs(userID uint, page, pageSize int) ([]model.SecurityLog, int64, error) {
	var logs []model.SecurityLog
	var total int64

	db := database.DB.Model(&model.SecurityLog{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// ==================== PasswordHistory ====================

// CreatePasswordHistory 创建密码历史
func CreatePasswordHistory(history *model.PasswordHistory) error {
	return database.DB.Create(history).Error
}

// GetRecentPasswords 获取最近N个密码
func GetRecentPasswords(userID uint, limit int) ([]model.PasswordHistory, error) {
	var histories []model.PasswordHistory
	err := database.DB.Where("user_id = ?", userID).
		Order("id DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

// ==================== OAuthUser ====================

// CreateOAuthUser 创建第三方登录绑定
func CreateOAuthUser(oauth *model.OAuthUser) error {
	return database.DB.Create(oauth).Error
}

// FindOAuthUser 查找第三方登录绑定
func FindOAuthUser(provider, providerUserID string) (*model.OAuthUser, error) {
	var oauth model.OAuthUser
	err := database.DB.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&oauth).Error
	return &oauth, err
}

// FindOAuthUsersByUserID 查找用户所有第三方绑定
func FindOAuthUsersByUserID(userID uint) ([]model.OAuthUser, error) {
	var oauths []model.OAuthUser
	err := database.DB.Where("user_id = ?", userID).Find(&oauths).Error
	return oauths, err
}

// DeleteOAuthUser 删除第三方登录绑定
func DeleteOAuthUser(userID uint, provider string) error {
	return database.DB.Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&model.OAuthUser{}).Error
}

// UpdateOAuthUser 更新第三方登录绑定信息
func UpdateOAuthUser(oauth *model.OAuthUser) error {
	return database.DB.Save(oauth).Error
}
