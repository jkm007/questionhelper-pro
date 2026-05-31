package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	notificationRepo "questionhelper-server/internal/repository/notification"
	"questionhelper-server/pkg/logger"
)

// ==================== 基础通知功能 ====================

// ListNotifications 通知列表
func ListNotifications(userID uint, req *dto.NotificationListRequest) ([]dto.NotificationInfo, int64, error) {
	notifications, total, err := notificationRepo.List(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询通知列表失败: %w", err)
	}

	list := make([]dto.NotificationInfo, 0, len(notifications))
	for _, n := range notifications {
		list = append(list, toNotificationInfo(&n))
	}
	return list, total, nil
}

// GetUnreadCount 获取未读数量
func GetUnreadCount(userID uint) (int, error) {
	count, err := notificationRepo.CountUnread(userID)
	if err != nil {
		return 0, fmt.Errorf("查询未读数量失败: %w", err)
	}
	return int(count), nil
}

// MarkAsRead 标记已读
func MarkAsRead(id, userID uint) error {
	notification, err := notificationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return fmt.Errorf("查询通知失败: %w", err)
	}

	if notification.UserID != userID {
		return errors.New("无权操作此通知")
	}

	if err := notificationRepo.MarkAsRead(id); err != nil {
		return fmt.Errorf("标记已读失败: %w", err)
	}

	return nil
}

// MarkAllAsRead 全部标记已读
func MarkAllAsRead(userID uint) error {
	if err := notificationRepo.MarkAllAsRead(userID); err != nil {
		return fmt.Errorf("全部标记已读失败: %w", err)
	}

	logger.Infof("用户 %d 全部标记已读", userID)
	return nil
}

// DeleteNotification 删除通知
func DeleteNotification(id, userID uint) error {
	notification, err := notificationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return fmt.Errorf("查询通知失败: %w", err)
	}

	if notification.UserID != userID {
		return errors.New("无权删除此通知")
	}

	if err := notificationRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除通知失败: %w", err)
	}

	return nil
}

// CreateNotification 创建通知（内部调用）
func CreateNotification(userID uint, notificationType int8, title, content string, targetType string, targetID uint) error {
	notification := &model.Notification{
		UserID:     userID,
		Type:       notificationType,
		Title:      title,
		Content:    content,
		TargetType: targetType,
		TargetID:   targetID,
		IsRead:     false,
	}

	if err := notificationRepo.Create(notification); err != nil {
		return fmt.Errorf("创建通知失败: %w", err)
	}

	return nil
}

// ==================== 撤回通知 ====================

// RecallNotification 撤回通知（只能撤回未读的）
func RecallNotification(id, senderID uint) error {
	notification, err := notificationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return fmt.Errorf("查询通知失败: %w", err)
	}

	// 只有发送者或系统管理员可以撤回
	if notification.SenderID != nil && *notification.SenderID != senderID {
		return errors.New("无权撤回此通知")
	}

	if notification.IsRead {
		return errors.New("已读通知无法撤回")
	}

	if notification.IsRecalled {
		return errors.New("通知已被撤回")
	}

	if err := notificationRepo.Recall(id); err != nil {
		return fmt.Errorf("撤回通知失败: %w", err)
	}

	logger.Infof("通知 %d 已被用户 %d 撤回", id, senderID)
	return nil
}

// ==================== 群发通知 ====================

// BatchSend 群发通知
func BatchSend(senderID uint, req *dto.BatchSendRequest) (*dto.BatchSendResponse, error) {
	batchID := uuid.New().String()

	// 收集目标用户ID
	userIDSet := make(map[uint]struct{})

	// 添加直接指定的用户
	for _, uid := range req.UserIDs {
		userIDSet[uid] = struct{}{}
	}

	// 根据目标类型收集用户
	switch req.TargetType {
	case "all":
		ids, err := notificationRepo.FindAllUserIDs()
		if err != nil {
			return nil, fmt.Errorf("查询全部用户失败: %w", err)
		}
		for _, uid := range ids {
			userIDSet[uid] = struct{}{}
		}
	case "role":
		for _, roleID := range req.TargetIDs {
			ids, err := notificationRepo.FindUserIDsByRole(roleID)
			if err != nil {
				logger.Errorf("查询角色 %d 用户失败: %v", roleID, err)
				continue
			}
			for _, uid := range ids {
				userIDSet[uid] = struct{}{}
			}
		}
	case "class":
		for _, classID := range req.TargetIDs {
			ids, err := notificationRepo.FindUserIDsByClass(classID)
			if err != nil {
				logger.Errorf("查询班级 %d 用户失败: %v", classID, err)
				continue
			}
			for _, uid := range ids {
				userIDSet[uid] = struct{}{}
			}
		}
	case "group":
		// 用户组目前等同于直接指定用户ID列表
		for _, uid := range req.TargetIDs {
			userIDSet[uid] = struct{}{}
		}
	}

	// 移除发送者自己
	delete(userIDSet, senderID)

	if len(userIDSet) == 0 {
		return &dto.BatchSendResponse{
			BatchID: batchID,
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	// 构建通知列表
	notifications := make([]model.Notification, 0, len(userIDSet))
	targetTypeJSON, _ := json.Marshal(req.TargetIDs)

	for uid := range userIDSet {
		notifications = append(notifications, model.Notification{
			UserID:     uid,
			SenderID:   &senderID,
			Type:       req.Type,
			Title:      req.Title,
			Content:    req.Content,
			TargetType: req.TargetType,
			Channel:    req.Channel,
			BatchID:    batchID,
			Extra:      string(targetTypeJSON),
			IsRead:     false,
		})
	}

	// 批量创建
	total := len(notifications)
	if err := notificationRepo.CreateBatch(notifications); err != nil {
		return nil, fmt.Errorf("批量创建通知失败: %w", err)
	}

	logger.Infof("用户 %d 群发通知 batch=%s, 共 %d 条", senderID, batchID, total)

	return &dto.BatchSendResponse{
		BatchID: batchID,
		Total:   total,
		Success: total,
		Failed:  0,
	}, nil
}

// ==================== 批量操作 ====================

// BatchMarkAsRead 批量标记已读
func BatchMarkAsRead(userID uint, ids []uint) error {
	if len(ids) == 0 {
		return errors.New("通知ID列表不能为空")
	}
	if err := notificationRepo.BatchMarkAsRead(ids, userID); err != nil {
		return fmt.Errorf("批量标记已读失败: %w", err)
	}
	logger.Infof("用户 %d 批量标记已读 %d 条", userID, len(ids))
	return nil
}

// BatchDeleteNotifications 批量删除
func BatchDeleteNotifications(userID uint, ids []uint) error {
	if len(ids) == 0 {
		return errors.New("通知ID列表不能为空")
	}
	if err := notificationRepo.BatchDelete(ids, userID); err != nil {
		return fmt.Errorf("批量删除失败: %w", err)
	}
	logger.Infof("用户 %d 批量删除 %d 条通知", userID, len(ids))
	return nil
}

// ==================== 通知模板管理 ====================

// ListTemplates 模板列表
func ListTemplates(req *dto.NotificationTemplateListRequest) ([]dto.TemplateInfo, int64, error) {
	templates, total, err := notificationRepo.ListTemplates(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询模板列表失败: %w", err)
	}

	list := make([]dto.TemplateInfo, 0, len(templates))
	for _, t := range templates {
		list = append(list, toTemplateInfo(&t))
	}
	return list, total, nil
}

// CreateTemplate 创建模板
func CreateTemplate(req *dto.CreateTemplateRequest) (*dto.TemplateInfo, error) {
	// 检查编码唯一性
	existing, _ := notificationRepo.FindTemplateByCode(req.Code)
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("模板编码已存在")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	tpl := &model.NotificationTemplate{
		Code:       req.Code,
		Name:       req.Name,
		TitleTpl:   req.TitleTpl,
		ContentTpl: req.ContentTpl,
		Type:       req.Type,
		Channel:    req.Channel,
		IsActive:   isActive,
		Remark:     req.Remark,
	}

	if err := notificationRepo.CreateTemplate(tpl); err != nil {
		return nil, fmt.Errorf("创建模板失败: %w", err)
	}

	info := toTemplateInfo(tpl)
	return &info, nil
}

// UpdateTemplate 更新模板
func UpdateTemplate(id uint, req *dto.UpdateTemplateRequest) (*dto.TemplateInfo, error) {
	tpl, err := notificationRepo.FindTemplateByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("模板不存在")
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.TitleTpl != "" {
		updates["title_tpl"] = req.TitleTpl
	}
	if req.ContentTpl != "" {
		updates["content_tpl"] = req.ContentTpl
	}
	if req.Type != 0 {
		updates["type"] = req.Type
	}
	if req.Channel != "" {
		updates["channel"] = req.Channel
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if len(updates) > 0 {
		if err := notificationRepo.UpdateTemplate(id, updates); err != nil {
			return nil, fmt.Errorf("更新模板失败: %w", err)
		}
	}

	// 重新查询
	tpl, _ = notificationRepo.FindTemplateByID(id)
	info := toTemplateInfo(tpl)
	return &info, nil
}

// DeleteTemplate 删除模板
func DeleteTemplate(id uint) error {
	_, err := notificationRepo.FindTemplateByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("模板不存在")
		}
		return fmt.Errorf("查询模板失败: %w", err)
	}

	if err := notificationRepo.DeleteTemplate(id); err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}
	return nil
}

// ==================== 定时通知 ====================

// CreateScheduled 创建定时通知
func CreateScheduled(senderID uint, req *dto.CreateScheduledRequest) (*dto.ScheduledNotificationInfo, error) {
	if req.ScheduledAt.Before(time.Now()) {
		return nil, errors.New("计划发送时间不能早于当前时间")
	}

	targetIDsJSON, _ := json.Marshal(req.TargetIDs)
	batchID := uuid.New().String()

	scheduled := &model.ScheduledNotification{
		SenderID:    senderID,
		Title:       req.Title,
		Content:     req.Content,
		Type:        req.Type,
		Channel:     req.Channel,
		TargetType:  req.TargetType,
		TargetIDs:   string(targetIDsJSON),
		ScheduledAt: req.ScheduledAt,
		Status:      0,
		BatchID:     batchID,
	}

	if err := notificationRepo.CreateScheduled(scheduled); err != nil {
		return nil, fmt.Errorf("创建定时通知失败: %w", err)
	}

	logger.Infof("用户 %d 创建定时通知 id=%d, 计划时间=%s", senderID, scheduled.ID, req.ScheduledAt.Format(time.RFC3339))

	info := toScheduledInfo(scheduled)
	return &info, nil
}

// ListScheduled 定时通知列表
func ListScheduled(userID uint, req *dto.ScheduledListRequest) ([]dto.ScheduledNotificationInfo, int64, error) {
	list, total, err := notificationRepo.ListScheduled(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询定时通知列表失败: %w", err)
	}

	result := make([]dto.ScheduledNotificationInfo, 0, len(list))
	for _, s := range list {
		result = append(result, toScheduledInfo(&s))
	}
	return result, total, nil
}

// DeleteScheduled 取消定时通知
func DeleteScheduled(id, userID uint) error {
	scheduled, err := notificationRepo.FindScheduledByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("定时通知不存在")
		}
		return fmt.Errorf("查询定时通知失败: %w", err)
	}

	if scheduled.SenderID != userID {
		return errors.New("无权取消此定时通知")
	}

	if scheduled.Status != 0 {
		return errors.New("只能取消待发送的定时通知")
	}

	if err := notificationRepo.DeleteScheduled(id); err != nil {
		return fmt.Errorf("取消定时通知失败: %w", err)
	}

	logger.Infof("用户 %d 取消定时通知 %d", userID, id)
	return nil
}

// ==================== 通知设置 ====================

// GetNotificationSettings 获取通知设置
func GetNotificationSettings(userID uint) ([]dto.NotificationSettingInfo, error) {
	settings, err := notificationRepo.FindSettingsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询通知设置失败: %w", err)
	}

	list := make([]dto.NotificationSettingInfo, 0, len(settings))
	for _, s := range settings {
		list = append(list, dto.NotificationSettingInfo{
			ID:           s.ID,
			Type:         s.Type,
			Channel:      s.Channel,
			Enabled:      s.Enabled,
			DoNotDisturb: s.DoNotDisturb,
			DisturbStart: s.DisturbStart,
			DisturbEnd:   s.DisturbEnd,
		})
	}
	return list, nil
}

// UpdateNotificationSettings 更新通知设置
func UpdateNotificationSettings(userID uint, req *dto.UpdateNotificationSettingRequest) error {
	for _, item := range req.Settings {
		enabled := true
		if item.Enabled != nil {
			enabled = *item.Enabled
		}
		doNotDisturb := false
		if item.DoNotDisturb != nil {
			doNotDisturb = *item.DoNotDisturb
		}

		setting := &model.NotificationSetting{
			UserID:       userID,
			Type:         item.Type,
			Channel:      item.Channel,
			Enabled:      enabled,
			DoNotDisturb: doNotDisturb,
			DisturbStart: item.DisturbStart,
			DisturbEnd:   item.DisturbEnd,
		}

		if err := notificationRepo.UpsertSetting(setting); err != nil {
			return fmt.Errorf("更新通知设置失败: %w", err)
		}
	}

	logger.Infof("用户 %d 更新通知设置 %d 条", userID, len(req.Settings))
	return nil
}

// ==================== 通知统计 ====================

// GetNotificationStats 获取通知统计
func GetNotificationStats(userID uint) (*dto.NotificationStats, error) {
	total, unread, read, err := notificationRepo.StatsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("查询通知统计失败: %w", err)
	}

	typeStats, err := notificationRepo.StatsByType(userID)
	if err != nil {
		return nil, fmt.Errorf("查询类型统计失败: %w", err)
	}

	channelStats, err := notificationRepo.StatsByChannel(userID)
	if err != nil {
		return nil, fmt.Errorf("查询渠道统计失败: %w", err)
	}

	dailyStats, err := notificationRepo.StatsDaily(userID, 30)
	if err != nil {
		return nil, fmt.Errorf("查询每日统计失败: %w", err)
	}

	return &dto.NotificationStats{
		TotalCount:   total,
		UnreadCount:  unread,
		ReadCount:    read,
		TypeStats:    typeStats,
		ChannelStats: channelStats,
		DailyStats:   dailyStats,
	}, nil
}

// ==================== 通知渠道管理 ====================

// ListChannels 渠道列表
func ListChannels(req *dto.ChannelListRequest) ([]dto.ChannelInfo, int64, error) {
	channels, total, err := notificationRepo.ListChannels(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询渠道列表失败: %w", err)
	}

	list := make([]dto.ChannelInfo, 0, len(channels))
	for _, ch := range channels {
		list = append(list, toChannelInfo(&ch))
	}
	return list, total, nil
}

// UpdateChannel 更新渠道配置
func UpdateChannel(id uint, req *dto.UpdateChannelRequest) (*dto.ChannelInfo, error) {
	channel, err := notificationRepo.FindChannelByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("渠道不存在")
		}
		return nil, fmt.Errorf("查询渠道失败: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}

	if len(updates) > 0 {
		if err := notificationRepo.UpdateChannel(id, updates); err != nil {
			return nil, fmt.Errorf("更新渠道配置失败: %w", err)
		}
	}

	channel, _ = notificationRepo.FindChannelByID(id)
	info := toChannelInfo(channel)
	return &info, nil
}

// ==================== 转换函数 ====================

func toNotificationInfo(n *model.Notification) dto.NotificationInfo {
	info := dto.NotificationInfo{
		ID:         n.ID,
		UserID:     n.UserID,
		Type:       n.Type,
		Title:      n.Title,
		Content:    n.Content,
		TargetType: n.TargetType,
		TargetID:   n.TargetID,
		IsRead:     n.IsRead,
		IsRecalled: n.IsRecalled,
		Channel:    n.Channel,
		BatchID:    n.BatchID,
		CreatedAt:  n.CreatedAt,
	}
	if n.SenderID != nil {
		// SenderID is stored but not exposed in NotificationInfo
	}
	return info
}

func toTemplateInfo(t *model.NotificationTemplate) dto.TemplateInfo {
	return dto.TemplateInfo{
		ID:         t.ID,
		Code:       t.Code,
		Name:       t.Name,
		TitleTpl:   t.TitleTpl,
		ContentTpl: t.ContentTpl,
		Type:       t.Type,
		Channel:    t.Channel,
		IsActive:   t.IsActive,
		Remark:     t.Remark,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func toScheduledInfo(s *model.ScheduledNotification) dto.ScheduledNotificationInfo {
	return dto.ScheduledNotificationInfo{
		ID:          s.ID,
		SenderID:    s.SenderID,
		Title:       s.Title,
		Content:     s.Content,
		Type:        s.Type,
		Channel:     s.Channel,
		TargetType:  s.TargetType,
		TargetIDs:   s.TargetIDs,
		ScheduledAt: s.ScheduledAt,
		Status:      s.Status,
		BatchID:     s.BatchID,
		ErrorMsg:    s.ErrorMsg,
		CreatedAt:   s.CreatedAt,
	}
}

func toChannelInfo(ch *model.NotificationChannel) dto.ChannelInfo {
	return dto.ChannelInfo{
		ID:        ch.ID,
		Name:      ch.Name,
		Type:      ch.Type,
		Config:    ch.Config,
		IsEnabled: ch.IsEnabled,
		Priority:  ch.Priority,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
	}
}
