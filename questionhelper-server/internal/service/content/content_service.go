package content

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	contentRepo "questionhelper-server/internal/repository/content"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/logger"
)

// ==================== 创作者申请 ====================

// ApplyCreator 申请成为创作者
func ApplyCreator(userID uint, req *dto.ApplyCreatorRequest) error {
	// 检查用户是否存在
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查是否已经是创作者
	for _, role := range user.Roles {
		if role.Code == "creator" {
			return errors.New("您已经是创作者")
		}
	}

	// 查找创作者角色
	var creatorRole *model.Role
	for _, role := range user.Roles {
		if role.Code == "creator" {
			creatorRole = &role
			break
		}
	}

	// 如果没有预加载角色，需要查询
	if creatorRole == nil {
		role, err := findRoleByCode("creator")
		if err != nil {
			return fmt.Errorf("查找创作者角色失败: %w", err)
		}
		creatorRole = role
	}

	// 检查是否已有待审核的申请
	_, err = contentRepo.FindApplicationByUserAndRole(userID, creatorRole.ID)
	if err == nil {
		return errors.New("您已有待审核的申请")
	}

	// 创建申请
	app := &model.RoleApplication{
		UserID: userID,
		RoleID: creatorRole.ID,
		Reason: req.Reason,
		Status: 0,
	}

	if err := contentRepo.CreateApplication(app); err != nil {
		return fmt.Errorf("创建申请失败: %w", err)
	}

	logger.Infof("用户 %d 申请成为创作者", userID)
	return nil
}

// findRoleByCode 根据角色代码查找角色
func findRoleByCode(code string) (*model.Role, error) {
	var role model.Role
	err := userRepo.GetDB().Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ==================== 创作者信息 ====================

// GetCreatorProfile 获取创作者资料
func GetCreatorProfile(userID uint) (*dto.CreatorProfileInfo, error) {
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查是否是创作者
	isCreator := false
	for _, role := range user.Roles {
		if role.Code == "creator" {
			isCreator = true
			break
		}
	}

	// 获取积分信息
	point, err := contentRepo.FindCreatorPoint(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询积分失败: %w", err)
	}

	level := 1
	levelName := "新手创作者"
	totalPoints := 0
	if point != nil {
		level = point.Level
		totalPoints = point.TotalPoints
		levelInfo, err := contentRepo.FindCreatorLevelByLevel(level)
		if err == nil {
			levelName = levelInfo.Name
		}
	}

	// 检查是否签署协议
	agreement, err := contentRepo.FindActiveAgreement()
	agreementSigned := false
	if err == nil {
		_, err = contentRepo.FindAgreementSign(agreement.ID, userID)
		agreementSigned = err == nil
	}

	return &dto.CreatorProfileInfo{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		Level:       level,
		LevelName:   levelName,
		TotalPoints: totalPoints,
		IsCreator:   isCreator,
		Agreement:   agreementSigned,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetCreatorLevel 获取创作者等级
func GetCreatorLevel(userID uint) ([]dto.CreatorLevelInfo, error) {
	// 获取用户当前积分
	point, err := contentRepo.FindCreatorPoint(userID)
	currentLevel := 1
	if err == nil {
		currentLevel = point.Level
	}

	// 获取所有等级
	levels, err := contentRepo.ListCreatorLevels()
	if err != nil {
		return nil, fmt.Errorf("查询等级列表失败: %w", err)
	}

	result := make([]dto.CreatorLevelInfo, 0, len(levels))
	for _, level := range levels {
		result = append(result, dto.CreatorLevelInfo{
			ID:        level.ID,
			Level:     level.Level,
			Name:      level.Name,
			MinPoints: level.MinPoints,
			MaxPoints: level.MaxPoints,
			Benefits:  level.Benefits,
			Icon:      level.Icon,
			IsCurrent: level.Level == currentLevel,
		})
	}

	return result, nil
}

// GetCreatorPoints 获取积分信息
func GetCreatorPoints(userID uint) (*dto.CreatorPointsInfo, error) {
	point, err := contentRepo.FindCreatorPoint(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有积分记录，创建一个
			point = &model.CreatorPoint{
				UserID:          userID,
				TotalPoints:     0,
				AvailablePoints: 0,
				Level:           1,
			}
			if err := contentRepo.CreateCreatorPoint(point); err != nil {
				return nil, fmt.Errorf("创建积分记录失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("查询积分失败: %w", err)
		}
	}

	// 获取当前等级名称
	levelInfo, err := contentRepo.FindCreatorLevelByLevel(point.Level)
	levelName := "新手创作者"
	if err == nil {
		levelName = levelInfo.Name
	}

	// 获取下一等级信息
	nextLevel := point.Level + 1
	nextLevelName := ""
	pointsNeeded := 0
	nextLevelInfo, err := contentRepo.FindCreatorLevelByLevel(nextLevel)
	if err == nil {
		nextLevelName = nextLevelInfo.Name
		pointsNeeded = nextLevelInfo.MinPoints - point.TotalPoints
		if pointsNeeded < 0 {
			pointsNeeded = 0
		}
	}

	return &dto.CreatorPointsInfo{
		UserID:          point.UserID,
		TotalPoints:     point.TotalPoints,
		AvailablePoints: point.AvailablePoints,
		Level:           point.Level,
		LevelName:       levelName,
		NextLevel:       nextLevel,
		NextLevelName:   nextLevelName,
		PointsNeeded:    pointsNeeded,
	}, nil
}

// GetCreatorPointLogs 获取积分记录
func GetCreatorPointLogs(userID uint, req *dto.CreatorPointLogsRequest) ([]dto.CreatorPointLogInfo, int64, error) {
	logs, total, err := contentRepo.ListCreatorPointLogs(userID, req.ChangeType, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, 0, fmt.Errorf("查询积分记录失败: %w", err)
	}

	result := make([]dto.CreatorPointLogInfo, 0, len(logs))
	for _, log := range logs {
		result = append(result, dto.CreatorPointLogInfo{
			ID:         log.ID,
			ChangeType: log.ChangeType,
			Points:     log.Points,
			Balance:    log.Balance,
			Reason:     log.Reason,
			CreatedAt:  log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// ==================== 创作者协议 ====================

// GetCreatorAgreement 获取协议
func GetCreatorAgreement(userID uint) (*dto.CreatorAgreementInfo, error) {
	agreement, err := contentRepo.FindActiveAgreement()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("暂无可用协议")
		}
		return nil, fmt.Errorf("查询协议失败: %w", err)
	}

	// 检查是否已签署
	signed := false
	signedAt := ""
	sign, err := contentRepo.FindAgreementSign(agreement.ID, userID)
	if err == nil {
		signed = true
		signedAt = sign.SignedAt.Format("2006-01-02 15:04:05")
	}

	return &dto.CreatorAgreementInfo{
		ID:        agreement.ID,
		Title:     agreement.Title,
		Content:   agreement.Content,
		Version:   agreement.Version,
		IsActive:  agreement.IsActive,
		Signed:    signed,
		SignedAt:  signedAt,
		CreatedAt: agreement.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// SignAgreement 签署协议
func SignAgreement(userID uint, req *dto.SignAgreementRequest) error {
	// 检查协议是否存在
	agreement, err := contentRepo.FindAgreementByID(req.AgreementID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("协议不存在")
		}
		return fmt.Errorf("查询协议失败: %w", err)
	}

	// 检查协议是否生效
	if !agreement.IsActive {
		return errors.New("协议未生效")
	}

	// 检查是否已签署
	_, err = contentRepo.FindAgreementSign(agreement.ID, userID)
	if err == nil {
		return errors.New("您已签署该协议")
	}

	// 创建签署记录
	sign := &model.CreatorAgreementSign{
		AgreementID: agreement.ID,
		UserID:      userID,
		SignedAt:    time.Now(),
	}

	if err := contentRepo.CreateAgreementSign(sign); err != nil {
		return fmt.Errorf("签署协议失败: %w", err)
	}

	logger.Infof("用户 %d 签署创作者协议 %d", userID, agreement.ID)
	return nil
}

// ==================== 创作者作品集 ====================

// GetPortfolio 获取作品集详情
func GetPortfolio(id uint) (*dto.PortfolioInfo, error) {
	portfolio, err := contentRepo.FindPortfolioByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("作品集不存在")
		}
		return nil, fmt.Errorf("查询作品集失败: %w", err)
	}

	// 统计内容数量
	itemCount, _ := contentRepo.CountPortfolioItems(id)

	return &dto.PortfolioInfo{
		ID:          portfolio.ID,
		UserID:      portfolio.UserID,
		Name:        portfolio.Name,
		Description: portfolio.Description,
		CoverImage:  portfolio.CoverImage,
		IsPublic:    portfolio.IsPublic,
		ViewCount:   portfolio.ViewCount,
		LikeCount:   portfolio.LikeCount,
		ItemCount:   itemCount,
		CreatedAt:   portfolio.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   portfolio.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ListPortfolios 获取作品集列表
func ListPortfolios(req *dto.PortfolioListRequest) ([]dto.PortfolioInfo, int64, error) {
	portfolios, total, err := contentRepo.ListPortfolios(req.UserID, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, 0, fmt.Errorf("查询作品集列表失败: %w", err)
	}

	result := make([]dto.PortfolioInfo, 0, len(portfolios))
	for _, p := range portfolios {
		itemCount, _ := contentRepo.CountPortfolioItems(p.ID)
		result = append(result, dto.PortfolioInfo{
			ID:          p.ID,
			UserID:      p.UserID,
			Name:        p.Name,
			Description: p.Description,
			CoverImage:  p.CoverImage,
			IsPublic:    p.IsPublic,
			ViewCount:   p.ViewCount,
			LikeCount:   p.LikeCount,
			ItemCount:   itemCount,
			CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// CreatePortfolio 创建作品集
func CreatePortfolio(userID uint, req *dto.CreatePortfolioRequest) error {
	portfolio := &model.CreatorPortfolio{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		IsPublic:    *req.IsPublic,
	}

	if err := contentRepo.CreatePortfolio(portfolio); err != nil {
		return fmt.Errorf("创建作品集失败: %w", err)
	}

	logger.Infof("用户 %d 创建作品集 %d", userID, portfolio.ID)
	return nil
}

// UpdatePortfolio 更新作品集
func UpdatePortfolio(userID uint, id uint, req *dto.UpdatePortfolioRequest) error {
	portfolio, err := contentRepo.FindPortfolioByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作品集不存在")
		}
		return fmt.Errorf("查询作品集失败: %w", err)
	}

	// 检查权限
	if portfolio.UserID != userID {
		return errors.New("无权操作该作品集")
	}

	if req.Name != "" {
		portfolio.Name = req.Name
	}
	if req.Description != "" {
		portfolio.Description = req.Description
	}
	if req.CoverImage != "" {
		portfolio.CoverImage = req.CoverImage
	}
	if req.IsPublic != nil {
		portfolio.IsPublic = *req.IsPublic
	}

	if err := contentRepo.UpdatePortfolio(portfolio); err != nil {
		return fmt.Errorf("更新作品集失败: %w", err)
	}

	logger.Infof("用户 %d 更新作品集 %d", userID, id)
	return nil
}

// DeletePortfolio 删除作品集
func DeletePortfolio(userID uint, id uint) error {
	portfolio, err := contentRepo.FindPortfolioByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作品集不存在")
		}
		return fmt.Errorf("查询作品集失败: %w", err)
	}

	// 检查权限
	if portfolio.UserID != userID {
		return errors.New("无权操作该作品集")
	}

	if err := contentRepo.DeletePortfolio(id); err != nil {
		return fmt.Errorf("删除作品集失败: %w", err)
	}

	logger.Infof("用户 %d 删除作品集 %d", userID, id)
	return nil
}

// ==================== 内容版本管理 ====================

// ListContentVersions 获取版本列表
func ListContentVersions(contentType string, contentID uint) ([]dto.ContentVersionInfo, error) {
	versions, err := contentRepo.ListContentVersions(contentType, contentID)
	if err != nil {
		return nil, fmt.Errorf("查询版本列表失败: %w", err)
	}

	result := make([]dto.ContentVersionInfo, 0, len(versions))
	for _, v := range versions {
		// 获取创建者名称
		creatorName := ""
		user, err := userRepo.FindByID(v.CreatorID)
		if err == nil {
			creatorName = user.Nickname
		}

		result = append(result, dto.ContentVersionInfo{
			ID:          v.ID,
			ContentType: v.ContentType,
			ContentID:   v.ContentID,
			Version:     v.Version,
			Data:        v.Data,
			CreatorID:   v.CreatorID,
			CreatorName: creatorName,
			Remark:      v.Remark,
			CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// GetContentVersion 获取版本详情
func GetContentVersion(contentType string, contentID uint, versionID uint) (*dto.ContentVersionInfo, error) {
	version, err := contentRepo.FindContentVersionByID(versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("版本不存在")
		}
		return nil, fmt.Errorf("查询版本失败: %w", err)
	}

	// 验证版本是否属于该内容
	if version.ContentType != contentType || version.ContentID != contentID {
		return nil, errors.New("版本不属于该内容")
	}

	// 获取创建者名称
	creatorName := ""
	user, err := userRepo.FindByID(version.CreatorID)
	if err == nil {
		creatorName = user.Nickname
	}

	return &dto.ContentVersionInfo{
		ID:          version.ID,
		ContentType: version.ContentType,
		ContentID:   version.ContentID,
		Version:     version.Version,
		Data:        version.Data,
		CreatorID:   version.CreatorID,
		CreatorName: creatorName,
		Remark:      version.Remark,
		CreatedAt:   version.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// RollbackContentVersion 回滚版本
func RollbackContentVersion(userID uint, contentType string, contentID uint, req *dto.RollbackVersionRequest) error {
	version, err := contentRepo.FindContentVersionByID(req.VersionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("查询版本失败: %w", err)
	}

	// 验证版本是否属于该内容
	if version.ContentType != contentType || version.ContentID != contentID {
		return errors.New("版本不属于该内容")
	}

	// 获取最新版本号
	latestVersion, err := contentRepo.GetLatestVersion(contentType, contentID)
	if err != nil {
		return fmt.Errorf("获取最新版本号失败: %w", err)
	}

	// 创建新版本（回滚版本）
	newVersion := &model.ContentVersion{
		ContentType: contentType,
		ContentID:   contentID,
		Version:     latestVersion + 1,
		Data:        version.Data,
		CreatorID:   userID,
		Remark:      fmt.Sprintf("回滚到版本 %d", version.Version),
	}

	if err := contentRepo.CreateContentVersion(newVersion); err != nil {
		return fmt.Errorf("创建回滚版本失败: %w", err)
	}

	logger.Infof("用户 %d 回滚内容 %s:%d 到版本 %d", userID, contentType, contentID, version.Version)
	return nil
}

// ==================== 内容标签 ====================

// GetContentTags 获取内容标签
func GetContentTags(contentType string, contentID uint) ([]dto.ContentTagInfo, error) {
	tags, err := contentRepo.ListContentTags(contentType, contentID)
	if err != nil {
		return nil, fmt.Errorf("查询内容标签失败: %w", err)
	}

	result := make([]dto.ContentTagInfo, 0, len(tags))
	for _, tag := range tags {
		result = append(result, dto.ContentTagInfo{
			ID:          tag.ID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			Color:       tag.Color,
			UsageCount:  tag.UsageCount,
		})
	}

	return result, nil
}

// AddContentTag 添加内容标签
func AddContentTag(contentType string, contentID uint, req *dto.AddContentTagRequest) error {
	tagName := strings.TrimSpace(req.TagName)
	if tagName == "" {
		return errors.New("标签名称不能为空")
	}

	// 查找或创建标签
	tag, err := contentRepo.FindContentTagByName(tagName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新标签
			tag = &model.ContentTag{
				Name:       tagName,
				Slug:       generateSlug(tagName),
				UsageCount: 0,
				IsActive:   true,
			}
			if err := contentRepo.CreateContentTag(tag); err != nil {
				return fmt.Errorf("创建标签失败: %w", err)
			}
		} else {
			return fmt.Errorf("查询标签失败: %w", err)
		}
	}

	// 检查是否已关联
	_, err = contentRepo.FindContentTagRelation(tag.ID, contentType, contentID)
	if err == nil {
		return errors.New("标签已存在")
	}

	// 创建关联
	relation := &model.ContentTagRelation{
		TagID:       tag.ID,
		ContentType: contentType,
		ContentID:   contentID,
	}

	if err := contentRepo.CreateContentTagRelation(relation); err != nil {
		return fmt.Errorf("添加标签失败: %w", err)
	}

	// 更新使用次数
	tag.UsageCount++
	contentRepo.UpdateContentTag(tag)

	logger.Infof("为内容 %s:%d 添加标签 %s", contentType, contentID, tagName)
	return nil
}

// RemoveContentTag 移除内容标签
func RemoveContentTag(contentType string, contentID uint, tagID uint) error {
	// 检查标签是否存在
	tag, err := contentRepo.FindContentTagByID(tagID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return fmt.Errorf("查询标签失败: %w", err)
	}

	// 删除关联
	if err := contentRepo.DeleteContentTagRelation(tagID, contentType, contentID); err != nil {
		return fmt.Errorf("移除标签失败: %w", err)
	}

	// 更新使用次数
	if tag.UsageCount > 0 {
		tag.UsageCount--
		contentRepo.UpdateContentTag(tag)
	}

	logger.Infof("从内容 %s:%d 移除标签 %d", contentType, contentID, tagID)
	return nil
}

// GetHotTags 获取热门标签
func GetHotTags(limit int) ([]dto.HotTagInfo, error) {
	if limit <= 0 {
		limit = 20
	}

	tags, err := contentRepo.ListHotTags(limit)
	if err != nil {
		return nil, fmt.Errorf("查询热门标签失败: %w", err)
	}

	result := make([]dto.HotTagInfo, 0, len(tags))
	for i, tag := range tags {
		result = append(result, dto.HotTagInfo{
			ID:         tag.ID,
			Name:       tag.Name,
			Slug:       tag.Slug,
			UsageCount: tag.UsageCount,
			Rank:       i + 1,
		})
	}

	return result, nil
}

// generateSlug 生成标签Slug
func generateSlug(name string) string {
	// 简单的Slug生成，实际项目中可能需要更复杂的处理
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}

// ==================== 内容收藏 ====================

// CheckContentFavorite 检查是否收藏
func CheckContentFavorite(userID uint, contentType string, contentID uint) (bool, error) {
	_, err := contentRepo.FindContentFavorite(userID, contentType, contentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("查询收藏状态失败: %w", err)
	}
	return true, nil
}

// AddContentFavorite 添加收藏
func AddContentFavorite(userID uint, contentType string, contentID uint, req *dto.AddContentFavoriteRequest) error {
	// 检查是否已收藏
	_, err := contentRepo.FindContentFavorite(userID, contentType, contentID)
	if err == nil {
		return errors.New("已收藏该内容")
	}

	folderName := "default"
	if req.FolderName != "" {
		folderName = req.FolderName
	}

	favorite := &model.ContentFavorite{
		UserID:      userID,
		ContentType: contentType,
		ContentID:   contentID,
		FolderName:  folderName,
	}

	if err := contentRepo.CreateContentFavorite(favorite); err != nil {
		return fmt.Errorf("添加收藏失败: %w", err)
	}

	logger.Infof("用户 %d 收藏内容 %s:%d", userID, contentType, contentID)
	return nil
}

// RemoveContentFavorite 取消收藏
func RemoveContentFavorite(userID uint, contentType string, contentID uint) error {
	// 检查是否已收藏
	_, err := contentRepo.FindContentFavorite(userID, contentType, contentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未收藏该内容")
		}
		return fmt.Errorf("查询收藏状态失败: %w", err)
	}

	if err := contentRepo.DeleteContentFavorite(userID, contentType, contentID); err != nil {
		return fmt.Errorf("取消收藏失败: %w", err)
	}

	logger.Infof("用户 %d 取消收藏内容 %s:%d", userID, contentType, contentID)
	return nil
}

// ListContentFavorites 获取收藏列表
func ListContentFavorites(userID uint, req *dto.FavoriteListRequest) ([]dto.FavoriteListInfo, int64, error) {
	favorites, total, err := contentRepo.ListContentFavorites(userID, req.ContentType, req.FolderName, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, 0, fmt.Errorf("查询收藏列表失败: %w", err)
	}

	result := make([]dto.FavoriteListInfo, 0, len(favorites))
	for _, f := range favorites {
		result = append(result, dto.FavoriteListInfo{
			ID:          f.ID,
			ContentType: f.ContentType,
			ContentID:   f.ContentID,
			FolderName:  f.FolderName,
			CreatedAt:   f.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// ==================== 内容预览 ====================

// GetContentPreview 获取内容预览
func GetContentPreview(contentType string, contentID uint) (*dto.ContentPreviewInfo, error) {
	// 根据内容类型获取内容
	switch contentType {
	case "question":
		return getQuestionPreview(contentID)
	case "exam":
		return getExamPreview(contentID)
	default:
		return nil, errors.New("不支持的内容类型")
	}
}

// getQuestionPreview 获取题目预览
func getQuestionPreview(questionID uint) (*dto.ContentPreviewInfo, error) {
	// 这里需要调用question模块的方法，暂时返回示例数据
	return &dto.ContentPreviewInfo{
		ContentType: "question",
		ContentID:   questionID,
		Title:       "题目预览",
		Summary:     "这是题目的摘要信息",
		Content:     nil,
		CreatorID:   0,
		CreatorName: "",
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// getExamPreview 获取考试预览
func getExamPreview(examID uint) (*dto.ContentPreviewInfo, error) {
	// 这里需要调用exam模块的方法，暂时返回示例数据
	return &dto.ContentPreviewInfo{
		ContentType: "exam",
		ContentID:   examID,
		Title:       "考试预览",
		Summary:     "这是考试的摘要信息",
		Content:     nil,
		CreatorID:   0,
		CreatorName: "",
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// ==================== 综合搜索 ====================

// Search 综合搜索
func Search(req *dto.SearchRequest) ([]dto.SearchResultInfo, int64, error) {
	// 记录搜索日志
	userID := uint(0) // 从上下文获取
	searchLog := &model.ContentSearchLog{
		UserID:    userID,
		Keyword:   req.Keyword,
		SearchedAt: time.Now(),
	}
	contentRepo.CreateSearchLog(searchLog)

	// 根据内容类型搜索
	switch req.ContentType {
	case "question":
		return searchQuestions(req)
	case "exam":
		return searchExams(req)
	default:
		return searchAll(req)
	}
}

// searchQuestions 搜索题目
func searchQuestions(req *dto.SearchRequest) ([]dto.SearchResultInfo, int64, error) {
	// 这里需要调用question模块的搜索方法，暂时返回空结果
	return []dto.SearchResultInfo{}, 0, nil
}

// searchExams 搜索考试
func searchExams(req *dto.SearchRequest) ([]dto.SearchResultInfo, int64, error) {
	// 这里需要调用exam模块的搜索方法，暂时返回空结果
	return []dto.SearchResultInfo{}, 0, nil
}

// searchAll 搜索所有内容
func searchAll(req *dto.SearchRequest) ([]dto.SearchResultInfo, int64, error) {
	// 综合搜索，暂时返回空结果
	return []dto.SearchResultInfo{}, 0, nil
}

// GetSearchSuggestions 获取搜索建议
func GetSearchSuggestions(keyword string, limit int) ([]dto.SearchSuggestionInfo, error) {
	if limit <= 0 {
		limit = 10
	}

	keywords, err := contentRepo.SearchContentSuggestions(keyword, limit)
	if err != nil {
		return nil, fmt.Errorf("获取搜索建议失败: %w", err)
	}

	result := make([]dto.SearchSuggestionInfo, 0, len(keywords))
	for _, kw := range keywords {
		result = append(result, dto.SearchSuggestionInfo{
			Keyword: kw,
			Type:    "history",
		})
	}

	return result, nil
}

// GetHotSearches 获取热门搜索
func GetHotSearches(limit int) ([]dto.HotSearchInfo, error) {
	if limit <= 0 {
		limit = 10
	}

	hotSearches, err := contentRepo.ListHotSearches(limit)
	if err != nil {
		return nil, fmt.Errorf("获取热门搜索失败: %w", err)
	}

	result := make([]dto.HotSearchInfo, 0, len(hotSearches))
	for i, hs := range hotSearches {
		result = append(result, dto.HotSearchInfo{
			Keyword: hs.Keyword,
			Count:   hs.Count,
			Rank:    i + 1,
		})
	}

	return result, nil
}

// GetSearchHistory 获取搜索历史
func GetSearchHistory(userID uint, limit int) ([]dto.SearchHistoryInfo, error) {
	if limit <= 0 {
		limit = 20
	}

	logs, err := contentRepo.ListSearchHistory(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取搜索历史失败: %w", err)
	}

	result := make([]dto.SearchHistoryInfo, 0, len(logs))
	for _, log := range logs {
		result = append(result, dto.SearchHistoryInfo{
			ID:        log.ID,
			Keyword:   log.Keyword,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// ==================== 审核流程 ====================

// ListPendingReviews 获取待审核列表
func ListPendingReviews(req *dto.ReviewListRequest) ([]dto.ReviewInstanceInfo, int64, error) {
	instances, total, err := contentRepo.ListReviewInstances(req.Status, req.ContentType, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, 0, fmt.Errorf("查询审核列表失败: %w", err)
	}

	result := make([]dto.ReviewInstanceInfo, 0, len(instances))
	for _, instance := range instances {
		// 获取创建者名称
		creatorName := ""
		user, err := userRepo.FindByID(instance.CreatorID)
		if err == nil {
			creatorName = user.Nickname
		}

		// 获取审核步骤
		steps, _ := contentRepo.ListReviewStepLogs(instance.ID)
		stepInfos := make([]dto.ReviewStepInfo, 0, len(steps))
		for _, step := range steps {
			stepInfos = append(stepInfos, dto.ReviewStepInfo{
				ID:         step.ID,
				StepIndex:  step.StepIndex,
				StepName:   step.StepName,
				ReviewerID: step.ReviewerID,
				Action:     step.Action,
				Opinion:    step.Opinion,
				OperatedAt: step.OperatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		result = append(result, dto.ReviewInstanceInfo{
			ID:           instance.ID,
			WorkflowID:   instance.WorkflowID,
			ContentType:  instance.ContentType,
			ContentID:    instance.ContentID,
			CurrentStep:  instance.CurrentStep,
			Status:       instance.Status,
			CreatorID:    instance.CreatorID,
			CreatorName:  creatorName,
			ContentTitle: "", // 需要根据内容类型获取标题
			Steps:        stepInfos,
			CreatedAt:    instance.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    instance.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// GetReviewDetail 获取审核详情
func GetReviewDetail(id uint) (*dto.ReviewInstanceInfo, error) {
	instance, err := contentRepo.FindReviewInstanceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("审核记录不存在")
		}
		return nil, fmt.Errorf("查询审核记录失败: %w", err)
	}

	// 获取创建者名称
	creatorName := ""
	user, err := userRepo.FindByID(instance.CreatorID)
	if err == nil {
		creatorName = user.Nickname
	}

	// 获取审核步骤
	steps, _ := contentRepo.ListReviewStepLogs(instance.ID)
	stepInfos := make([]dto.ReviewStepInfo, 0, len(steps))
	for _, step := range steps {
		stepInfos = append(stepInfos, dto.ReviewStepInfo{
			ID:         step.ID,
			StepIndex:  step.StepIndex,
			StepName:   step.StepName,
			ReviewerID: step.ReviewerID,
			Action:     step.Action,
			Opinion:    step.Opinion,
			OperatedAt: step.OperatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.ReviewInstanceInfo{
		ID:           instance.ID,
		WorkflowID:   instance.WorkflowID,
		ContentType:  instance.ContentType,
		ContentID:    instance.ContentID,
		CurrentStep:  instance.CurrentStep,
		Status:       instance.Status,
		CreatorID:    instance.CreatorID,
		CreatorName:  creatorName,
		ContentTitle: "", // 需要根据内容类型获取标题
		Steps:        stepInfos,
		CreatedAt:    instance.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    instance.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ApproveReview 通过审核
func ApproveReview(reviewerID uint, id uint, req *dto.ApproveReviewRequest) error {
	instance, err := contentRepo.FindReviewInstanceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("审核记录不存在")
		}
		return fmt.Errorf("查询审核记录失败: %w", err)
	}

	// 检查状态
	if instance.Status != 0 {
		return errors.New("审核记录状态异常")
	}

	// 更新状态
	instance.Status = 1
	if err := contentRepo.UpdateReviewInstance(instance); err != nil {
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 记录审核步骤
	stepLog := &model.ReviewStepLog{
		InstanceID: instance.ID,
		StepIndex:  instance.CurrentStep,
		StepName:   "审核通过",
		ReviewerID: reviewerID,
		Action:     "pass",
		Opinion:    req.Opinion,
		OperatedAt: time.Now(),
	}
	contentRepo.CreateReviewStepLog(stepLog)

	// 发送通知
	notification := &model.ReviewNotification{
		InstanceID: instance.ID,
		UserID:     instance.CreatorID,
		Type:       "review_approved",
		Message:    "您的内容已通过审核",
	}
	contentRepo.CreateReviewNotification(notification)

	logger.Infof("审核人 %d 通过审核 %d", reviewerID, id)
	return nil
}

// RejectReview 拒绝审核
func RejectReview(reviewerID uint, id uint, req *dto.RejectReviewRequest) error {
	instance, err := contentRepo.FindReviewInstanceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("审核记录不存在")
		}
		return fmt.Errorf("查询审核记录失败: %w", err)
	}

	// 检查状态
	if instance.Status != 0 {
		return errors.New("审核记录状态异常")
	}

	// 更新状态
	instance.Status = 2
	if err := contentRepo.UpdateReviewInstance(instance); err != nil {
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 记录审核步骤
	stepLog := &model.ReviewStepLog{
		InstanceID: instance.ID,
		StepIndex:  instance.CurrentStep,
		StepName:   "审核拒绝",
		ReviewerID: reviewerID,
		Action:     "reject",
		Opinion:    req.Opinion,
		OperatedAt: time.Now(),
	}
	contentRepo.CreateReviewStepLog(stepLog)

	// 发送通知
	notification := &model.ReviewNotification{
		InstanceID: instance.ID,
		UserID:     instance.CreatorID,
		Type:       "review_rejected",
		Message:    fmt.Sprintf("您的内容未通过审核: %s", req.Opinion),
	}
	contentRepo.CreateReviewNotification(notification)

	logger.Infof("审核人 %d 拒绝审核 %d", reviewerID, id)
	return nil
}

// AddReviewComment 添加审核意见
func AddReviewComment(userID uint, id uint, req *dto.AddReviewCommentRequest) error {
	instance, err := contentRepo.FindReviewInstanceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("审核记录不存在")
		}
		return fmt.Errorf("查询审核记录失败: %w", err)
	}

	// 创建回复
	reply := &model.ReviewReply{
		InstanceID: instance.ID,
		StepLogID:  0, // 可以关联到具体的步骤
		UserID:     userID,
		Content:    req.Content,
		ReplyAt:    time.Now(),
	}

	if err := contentRepo.CreateReviewReply(reply); err != nil {
		return fmt.Errorf("添加审核意见失败: %w", err)
	}

	logger.Infof("用户 %d 为审核 %d 添加意见", userID, id)
	return nil
}
