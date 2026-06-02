package class

import (
	"errors"
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	classRepo "questionhelper-server/internal/repository/class"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ==================== Creator Management ====================

// ListCreators 创作者列表
func ListCreators(classID uint) ([]dto.CreatorInfo, error) {
	perms, err := classRepo.ListCreatorPermissions(classID)
	if err != nil {
		return nil, fmt.Errorf("查询创作者列表失败: %w", err)
	}

	list := make([]dto.CreatorInfo, 0, len(perms))
	for _, p := range perms {
		userName := ""
		if p.User.ID > 0 {
			userName = p.User.Nickname
		}
		list = append(list, dto.CreatorInfo{
			UserID:     p.UserID,
			UserName:   userName,
			MaxClasses: p.MaxClasses,
			CanCreate:  p.CanCreate,
			ExpiresAt:  p.ExpiresAt,
			GrantedBy:  p.GrantedBy,
			Reason:     p.Reason,
			CreatedAt:  p.CreatedAt,
		})
	}
	return list, nil
}

// CreatorApply 申请创作者
func CreatorApply(classID, userID uint, req *dto.CreatorApplyRequest) error {
	// 检查是否已有权限
	if _, err := classRepo.FindCreatorPermission(userID); err == nil {
		return errors.New("已有创作者权限")
	}

	maxClasses := req.MaxClasses
	if maxClasses <= 0 {
		maxClasses = 5
	}

	perm := &model.ClassCreatorPermission{
		UserID:     userID,
		MaxClasses: maxClasses,
		CanCreate:  false, // 待审批
		GrantedBy:  0,
		Reason:     req.Reason,
	}

	if err := classRepo.CreateCreatorPermission(perm); err != nil {
		return fmt.Errorf("提交创作者申请失败: %w", err)
	}

	logger.Infof("用户 %d 申请创作者权限", userID)
	return nil
}

// ApproveCreatorApplication 审批通过
func ApproveCreatorApplication(classID, appID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 3); err != nil {
		return err
	}

	// 查询申请记录（ClassCreatorPermission 中 can_create=false 的视为待审批）
	var perm model.ClassCreatorPermission
	if err := database.DB.First(&perm, appID).Error; err != nil {
		return errors.New("申请记录不存在")
	}

	if perm.CanCreate {
		return errors.New("该申请已审批通过")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"can_create": true,
		"granted_by": operatorID,
		"updated_at": now,
	}

	if err := database.DB.Model(&perm).Updates(updates).Error; err != nil {
		return fmt.Errorf("审批通过失败: %w", err)
	}

	logger.Infof("审批通过创作者申请 %d，操作人 %d", appID, operatorID)
	return nil
}

// RejectCreatorApplication 审批驳回
func RejectCreatorApplication(classID, appID, operatorID uint, remark string) error {
	if err := checkClassPermission(classID, operatorID, 3); err != nil {
		return err
	}

	// 查询申请记录
	var perm model.ClassCreatorPermission
	if err := database.DB.First(&perm, appID).Error; err != nil {
		return errors.New("申请记录不存在")
	}

	if perm.CanCreate {
		return errors.New("该申请已审批通过，无法驳回")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"can_create": false,
		"granted_by": operatorID,
		"reason":     remark,
		"updated_at": now,
	}

	if err := database.DB.Model(&perm).Updates(updates).Error; err != nil {
		return fmt.Errorf("审批驳回失败: %w", err)
	}

	logger.Infof("审批驳回创作者申请 %d，操作人 %d", appID, operatorID)
	return nil
}

// RemoveCreator 撤销创作者
func RemoveCreator(classID, targetUserID, operatorID uint) error {
	if err := checkClassPermission(classID, operatorID, 3); err != nil {
		return err
	}

	if err := classRepo.DeleteCreatorPermission(targetUserID); err != nil {
		return fmt.Errorf("撤销创作者权限失败: %w", err)
	}

	logger.Infof("撤销用户 %d 创作者权限", targetUserID)
	return nil
}

// ListCreatorApplications 申请列表
func ListCreatorApplications(classID uint, req *dto.PageRequest) ([]dto.CreatorApplicationInfo, int64, error) {
	apps, total, err := classRepo.ListCreatorApplications(classID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询创作者申请列表失败: %w", err)
	}

	list := make([]dto.CreatorApplicationInfo, 0, len(apps))
	for _, a := range apps {
		userName := ""
		if a.User.ID > 0 {
			userName = a.User.Nickname
		}
		list = append(list, dto.CreatorApplicationInfo{
			ID:         a.ID,
			UserID:     a.UserID,
			UserName:   userName,
			Reason:     a.Reason,
			MaxClasses: a.MaxClasses,
			Status:     boolToInt8(a.CanCreate),
			Remark:     a.Reason,
			ReviewBy:   &a.GrantedBy,
			CreatedAt:  a.CreatedAt,
		})
	}
	return list, total, nil
}
