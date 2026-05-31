package class

import (
	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Homework Submission ====================

func FindSubmission(id uint) (*model.HomeworkSubmission, error) {
	var sub model.HomeworkSubmission
	err := database.DB.First(&sub, id).Error
	return &sub, err
}

func CreateSubmission(sub *model.HomeworkSubmission) error {
	return database.DB.Create(sub).Error
}

func UpdateSubmission(sub *model.HomeworkSubmission) error {
	return database.DB.Save(sub).Error
}

func ListSubmissions(homeworkID uint, req *dto.PageRequest) ([]model.HomeworkSubmission, int64, error) {
	var subs []model.HomeworkSubmission
	var total int64

	db := database.DB.Model(&model.HomeworkSubmission{}).Where("homework_id = ?", homeworkID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("submitted_at DESC").Find(&subs).Error
	return subs, total, err
}

func FindSubmissionByUser(homeworkID, userID uint) (*model.HomeworkSubmission, error) {
	var sub model.HomeworkSubmission
	err := database.DB.Where("homework_id = ? AND user_id = ?", homeworkID, userID).
		Order("id DESC").First(&sub).Error
	return &sub, err
}

// ==================== Peer Review ====================

func CreatePeerReview(review *model.ClassPeerReview) error {
	return database.DB.Create(review).Error
}

func FindPeerReview(id uint) (*model.ClassPeerReview, error) {
	var r model.ClassPeerReview
	err := database.DB.First(&r, id).Error
	return &r, err
}

func UpdatePeerReview(r *model.ClassPeerReview) error {
	return database.DB.Save(r).Error
}

func ListPeerReviews(targetType int8, targetID uint, req *dto.PageRequest) ([]model.ClassPeerReview, int64, error) {
	var reviews []model.ClassPeerReview
	var total int64

	db := database.DB.Model(&model.ClassPeerReview{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&reviews).Error
	return reviews, total, err
}

func ListPeerReviewsByReviewer(targetType int8, targetID, reviewerID uint) ([]model.ClassPeerReview, error) {
	var reviews []model.ClassPeerReview
	err := database.DB.Where("target_type = ? AND target_id = ? AND reviewer_id = ?",
		targetType, targetID, reviewerID).Find(&reviews).Error
	return reviews, err
}

func ListPeerReviewsByReviewee(targetType int8, targetID, revieweeID uint) ([]model.ClassPeerReview, error) {
	var reviews []model.ClassPeerReview
	err := database.DB.Where("target_type = ? AND target_id = ? AND reviewee_id = ?",
		targetType, targetID, revieweeID).Find(&reviews).Error
	return reviews, err
}

func BatchCreatePeerReviews(reviews []model.ClassPeerReview) error {
	return database.DB.CreateInBatches(reviews, 100).Error
}

// ==================== Class Group ====================

func FindGroup(id uint) (*model.ClassGroup, error) {
	var g model.ClassGroup
	err := database.DB.First(&g, id).Error
	return &g, err
}

func CreateGroup(g *model.ClassGroup) error {
	return database.DB.Create(g).Error
}

func UpdateGroup(g *model.ClassGroup) error {
	return database.DB.Save(g).Error
}

func DeleteGroup(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", id).Delete(&model.ClassGroupMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClassGroup{}, id).Error
	})
}

func ListGroups(classID uint) ([]model.ClassGroup, error) {
	var groups []model.ClassGroup
	err := database.DB.Where("class_id = ?", classID).Order("sort_order ASC, id ASC").Find(&groups).Error
	return groups, err
}

func ListGroupMembers(groupID uint) ([]model.ClassGroupMember, error) {
	var members []model.ClassGroupMember
	err := database.DB.Where("group_id = ?", groupID).Order("joined_at ASC").Find(&members).Error
	return members, err
}

func FindGroupMember(groupID, userID uint) (*model.ClassGroupMember, error) {
	var m model.ClassGroupMember
	err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&m).Error
	return &m, err
}

func CreateGroupMember(m *model.ClassGroupMember) error {
	return database.DB.Create(m).Error
}

func DeleteGroupMember(groupID, userID uint) error {
	return database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.ClassGroupMember{}).Error
}

func CountGroupMembers(groupID uint) int64 {
	var count int64
	database.DB.Model(&model.ClassGroupMember{}).Where("group_id = ?", groupID).Count(&count)
	return count
}

// ==================== Join Application ====================

func FindApplication(id uint) (*model.ClassJoinApplication, error) {
	var app model.ClassJoinApplication
	err := database.DB.First(&app, id).Error
	return &app, err
}

func CreateApplication(app *model.ClassJoinApplication) error {
	return database.DB.Create(app).Error
}

func UpdateApplication(app *model.ClassJoinApplication) error {
	return database.DB.Save(app).Error
}

func ListApplications(classID uint, req *dto.PageRequest) ([]model.ClassJoinApplication, int64, error) {
	var apps []model.ClassJoinApplication
	var total int64

	db := database.DB.Model(&model.ClassJoinApplication{}).Where("class_id = ?", classID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&apps).Error
	return apps, total, err
}

func FindPendingApplication(classID, userID uint) (*model.ClassJoinApplication, error) {
	var app model.ClassJoinApplication
	err := database.DB.Where("class_id = ? AND user_id = ? AND status = 0", classID, userID).
		First(&app).Error
	return &app, err
}

// ==================== Attendance ====================

func FindAttendance(id uint) (*model.ClassAttendance, error) {
	var att model.ClassAttendance
	err := database.DB.First(&att, id).Error
	return &att, err
}

func CreateAttendance(att *model.ClassAttendance) error {
	return database.DB.Create(att).Error
}

func UpdateAttendance(att *model.ClassAttendance) error {
	return database.DB.Save(att).Error
}

func DeleteAttendance(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("attendance_id = ?", id).Delete(&model.ClassAttendanceRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClassAttendance{}, id).Error
	})
}

func ListAttendances(classID uint, req *dto.PageRequest) ([]model.ClassAttendance, int64, error) {
	var atts []model.ClassAttendance
	var total int64

	db := database.DB.Model(&model.ClassAttendance{}).Where("class_id = ?", classID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&atts).Error
	return atts, total, err
}

func FindAttendanceRecord(attendanceID, userID uint) (*model.ClassAttendanceRecord, error) {
	var rec model.ClassAttendanceRecord
	err := database.DB.Where("attendance_id = ? AND user_id = ?", attendanceID, userID).
		First(&rec).Error
	return &rec, err
}

func CreateAttendanceRecord(rec *model.ClassAttendanceRecord) error {
	return database.DB.Create(rec).Error
}

func UpdateAttendanceRecord(rec *model.ClassAttendanceRecord) error {
	return database.DB.Save(rec).Error
}

func ListAttendanceRecords(attendanceID uint) ([]model.ClassAttendanceRecord, error) {
	var recs []model.ClassAttendanceRecord
	err := database.DB.Where("attendance_id = ?", attendanceID).
		Order("signed_at ASC").Find(&recs).Error
	return recs, err
}

func ListAttendanceRecordsPaged(attendanceID uint, req *dto.PageRequest) ([]model.ClassAttendanceRecord, int64, error) {
	var recs []model.ClassAttendanceRecord
	var total int64

	db := database.DB.Model(&model.ClassAttendanceRecord{}).Where("attendance_id = ?", attendanceID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("signed_at ASC").Find(&recs).Error
	return recs, total, err
}

// ==================== Study Plan ====================

func FindStudyPlan(id uint) (*model.ClassStudyPlan, error) {
	var plan model.ClassStudyPlan
	err := database.DB.First(&plan, id).Error
	return &plan, err
}

func CreateStudyPlan(plan *model.ClassStudyPlan) error {
	return database.DB.Create(plan).Error
}

func UpdateStudyPlan(plan *model.ClassStudyPlan) error {
	return database.DB.Save(plan).Error
}

func DeleteStudyPlan(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id = ?", id).Delete(&model.ClassStudyPlanProgress{}).Error; err != nil {
			return err
		}
		if err := tx.Where("plan_id = ?", id).Delete(&model.ClassStudyPlanItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClassStudyPlan{}, id).Error
	})
}

func ListStudyPlans(classID uint, req *dto.PageRequest) ([]model.ClassStudyPlan, int64, error) {
	var plans []model.ClassStudyPlan
	var total int64

	db := database.DB.Model(&model.ClassStudyPlan{}).Where("class_id = ?", classID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&plans).Error
	return plans, total, err
}

func FindStudyPlanItem(id uint) (*model.ClassStudyPlanItem, error) {
	var item model.ClassStudyPlanItem
	err := database.DB.First(&item, id).Error
	return &item, err
}

func CreateStudyPlanItem(item *model.ClassStudyPlanItem) error {
	return database.DB.Create(item).Error
}

func UpdateStudyPlanItem(item *model.ClassStudyPlanItem) error {
	return database.DB.Save(item).Error
}

func DeleteStudyPlanItem(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("item_id = ?", id).Delete(&model.ClassStudyPlanProgress{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClassStudyPlanItem{}, id).Error
	})
}

func ListStudyPlanItems(planID uint) ([]model.ClassStudyPlanItem, error) {
	var items []model.ClassStudyPlanItem
	err := database.DB.Where("plan_id = ?", planID).Order("sort_order ASC, id ASC").Find(&items).Error
	return items, err
}

func FindStudyPlanProgress(planID, itemID, userID uint) (*model.ClassStudyPlanProgress, error) {
	var p model.ClassStudyPlanProgress
	err := database.DB.Where("plan_id = ? AND item_id = ? AND user_id = ?", planID, itemID, userID).
		First(&p).Error
	return &p, err
}

func CreateStudyPlanProgress(p *model.ClassStudyPlanProgress) error {
	return database.DB.Create(p).Error
}

func UpdateStudyPlanProgress(p *model.ClassStudyPlanProgress) error {
	return database.DB.Save(p).Error
}

func ListStudyPlanProgressByPlan(planID uint) ([]model.ClassStudyPlanProgress, error) {
	var progress []model.ClassStudyPlanProgress
	err := database.DB.Where("plan_id = ?", planID).Find(&progress).Error
	return progress, err
}

func CountStudyPlanItems(planID uint) int64 {
	var count int64
	database.DB.Model(&model.ClassStudyPlanItem{}).Where("plan_id = ?", planID).Count(&count)
	return count
}

// ==================== Class File ====================

func FindClassFile(id uint) (*model.ClassFile, error) {
	var f model.ClassFile
	err := database.DB.First(&f, id).Error
	return &f, err
}

func CreateClassFile(f *model.ClassFile) error {
	return database.DB.Create(f).Error
}

func DeleteClassFile(id uint) error {
	return database.DB.Delete(&model.ClassFile{}, id).Error
}

func ListClassFiles(classID uint, req *dto.PageRequest) ([]model.ClassFile, int64, error) {
	var files []model.ClassFile
	var total int64

	db := database.DB.Model(&model.ClassFile{}).Where("class_id = ?", classID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&files).Error
	return files, total, err
}

func IncrementFileDownloads(id uint) error {
	return database.DB.Model(&model.ClassFile{}).Where("id = ?", id).
		UpdateColumn("downloads", database.DB.Raw("downloads + 1")).Error
}

// ==================== Ranking ====================

func ListRankings(classID uint, rankingType string) ([]model.ClassRanking, error) {
	var rankings []model.ClassRanking
	err := database.DB.Where("class_id = ? AND ranking_type = ?", classID, rankingType).
		Order("`rank` ASC").Find(&rankings).Error
	return rankings, err
}

func UpsertRanking(r *model.ClassRanking) error {
	var existing model.ClassRanking
	err := database.DB.Where("class_id = ? AND user_id = ? AND ranking_type = ?",
		r.ClassID, r.UserID, r.RankingType).First(&existing).Error
	if err != nil {
		return database.DB.Create(r).Error
	}
	existing.Score = r.Score
	existing.Rank = r.Rank
	return database.DB.Save(&existing).Error
}

func DeleteRankings(classID uint, rankingType string) error {
	return database.DB.Where("class_id = ? AND ranking_type = ?", classID, rankingType).
		Delete(&model.ClassRanking{}).Error
}

// ==================== Tag ====================

func FindTag(id uint) (*model.ClassTag, error) {
	var tag model.ClassTag
	err := database.DB.First(&tag, id).Error
	return &tag, err
}

func CreateTag(tag *model.ClassTag) error {
	return database.DB.Create(tag).Error
}

func UpdateTag(tag *model.ClassTag) error {
	return database.DB.Save(tag).Error
}

func DeleteTag(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tag_id = ?", id).Delete(&model.ClassTagRelation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClassTag{}, id).Error
	})
}

func ListTags() ([]model.ClassTag, error) {
	var tags []model.ClassTag
	err := database.DB.Where("status = 1").Order("sort_order ASC, id ASC").Find(&tags).Error
	return tags, err
}

func ListAllTags() ([]model.ClassTag, error) {
	var tags []model.ClassTag
	err := database.DB.Order("sort_order ASC, id ASC").Find(&tags).Error
	return tags, err
}

func AddClassTag(classID, tagID uint) error {
	rel := &model.ClassTagRelation{ClassID: classID, TagID: tagID}
	return database.DB.Create(rel).Error
}

func RemoveClassTag(classID, tagID uint) error {
	return database.DB.Where("class_id = ? AND tag_id = ?", classID, tagID).
		Delete(&model.ClassTagRelation{}).Error
}

func ListClassTags(classID uint) ([]model.ClassTag, error) {
	var tags []model.ClassTag
	err := database.DB.Joins("JOIN class_tag_relations ON class_tag_relations.tag_id = class_tags.id").
		Where("class_tag_relations.class_id = ?", classID).Find(&tags).Error
	return tags, err
}

// ==================== Template ====================

func FindTemplate(id uint) (*model.ClassTemplate, error) {
	var tpl model.ClassTemplate
	err := database.DB.First(&tpl, id).Error
	return &tpl, err
}

func CreateTemplate(tpl *model.ClassTemplate) error {
	return database.DB.Create(tpl).Error
}

func UpdateTemplate(tpl *model.ClassTemplate) error {
	return database.DB.Save(tpl).Error
}

func DeleteTemplate(id uint) error {
	return database.DB.Delete(&model.ClassTemplate{}, id).Error
}

func ListTemplates(req *dto.PageRequest, creatorID uint, isAdmin bool) ([]model.ClassTemplate, int64, error) {
	var templates []model.ClassTemplate
	var total int64

	db := database.DB.Model(&model.ClassTemplate{}).Where("status = 1")
	if !isAdmin {
		db = db.Where("creator_id = ? OR is_public = true", creatorID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("used_count DESC, id DESC").Find(&templates).Error
	return templates, total, err
}

func IncrementTemplateUsedCount(id uint) error {
	return database.DB.Model(&model.ClassTemplate{}).Where("id = ?", id).
		UpdateColumn("used_count", database.DB.Raw("used_count + 1")).Error
}

// ==================== Discussion ====================

func FindDiscussion(id uint) (*model.ClassDiscussion, error) {
	var d model.ClassDiscussion
	err := database.DB.Preload("Creator").First(&d, id).Error
	return &d, err
}

func CreateDiscussion(d *model.ClassDiscussion) error {
	return database.DB.Create(d).Error
}

func UpdateDiscussion(d *model.ClassDiscussion) error {
	return database.DB.Save(d).Error
}

func DeleteDiscussion(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("discussion_id = ?", id).Delete(&model.ClassDiscussionReply{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClassDiscussion{}, id).Error
	})
}

func ListDiscussions(classID uint, req *dto.PageRequest) ([]model.ClassDiscussion, int64, error) {
	var discussions []model.ClassDiscussion
	var total int64

	db := database.DB.Model(&model.ClassDiscussion{}).Where("class_id = ?", classID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Creator").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("is_top DESC, last_reply_at DESC, created_at DESC").Find(&discussions).Error
	return discussions, total, err
}

func IncrementDiscussionViewCount(id uint) error {
	return database.DB.Model(&model.ClassDiscussion{}).Where("id = ?", id).
		UpdateColumn("view_count", database.DB.Raw("view_count + 1")).Error
}

// ==================== Class Resource ====================

func CountClassFiles(classID uint) int64 {
	var count int64
	database.DB.Model(&model.ClassFile{}).Where("class_id = ?", classID).Count(&count)
	return count
}

func SumClassFileSize(classID uint) int64 {
	var total int64
	database.DB.Model(&model.ClassFile{}).Where("class_id = ?", classID).
		Select("COALESCE(SUM(size), 0)").Scan(&total)
	return total
}

func CountClassHomework(classID uint) int64 {
	var count int64
	database.DB.Model(&model.Homework{}).Where("class_id = ?", classID).Count(&count)
	return count
}

func CountClassNotices(classID uint) int64 {
	var count int64
	database.DB.Model(&model.ClassNotice{}).Where("class_id = ?", classID).Count(&count)
	return count
}

func ListClassResources(classID uint, req *dto.PageRequest) ([]model.ClassResourceVersion, int64, error) {
	var resources []model.ClassResourceVersion
	var total int64

	db := database.DB.Model(&model.ClassResourceVersion{}).
		Where("creator_id IN (SELECT user_id FROM class_members WHERE class_id = ?)", classID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&resources).Error
	return resources, total, err
}

func CreateClassResource(r *model.ClassResourceVersion) error {
	return database.DB.Create(r).Error
}

// ==================== Creator Permission ====================

func FindCreatorPermission(userID uint) (*model.ClassCreatorPermission, error) {
	var p model.ClassCreatorPermission
	err := database.DB.Where("user_id = ?", userID).First(&p).Error
	return &p, err
}

func CreateCreatorPermission(p *model.ClassCreatorPermission) error {
	return database.DB.Create(p).Error
}

func UpdateCreatorPermission(p *model.ClassCreatorPermission) error {
	return database.DB.Save(p).Error
}

func DeleteCreatorPermission(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.ClassCreatorPermission{}).Error
}

func ListCreatorPermissions(classID uint) ([]model.ClassCreatorPermission, error) {
	var perms []model.ClassCreatorPermission
	// 创作者权限是全局的，这里列出所有有权限的用户
	err := database.DB.Preload("User").Where("can_create = true").Find(&perms).Error
	return perms, err
}

func ListCreatorApplications(classID uint, req *dto.PageRequest) ([]model.ClassCreatorPermission, int64, error) {
	// 使用 ClassCreatorPermission 表来管理创作者申请
	// 这里简化处理：将 ClassCreatorPermission 中 can_create=false 的视为申请中
	var apps []model.ClassCreatorPermission
	var total int64

	db := database.DB.Model(&model.ClassCreatorPermission{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&apps).Error
	return apps, total, err
}

// ==================== Class Enhancement ====================

func SearchClasses(keyword string, req *dto.PageRequest) ([]model.Class, int64, error) {
	var classes []model.Class
	var total int64

	db := database.DB.Model(&model.Class{}).
		Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Creator").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("member_count DESC").Find(&classes).Error
	return classes, total, err
}

func ListClassExams(classID uint, req *dto.PageRequest) ([]model.Exam, int64, error) {
	var exams []model.Exam
	var total int64

	// 通过班级成员ID来查找关联的考试
	db := database.DB.Model(&model.Exam{}).
		Where("id IN (SELECT exam_id FROM exam_classes WHERE class_id = ?)", classID)

	if err := db.Count(&total).Error; err != nil {
		// 如果 exam_classes 表不存在，尝试其他方式
		return nil, 0, nil
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&exams).Error
	return exams, total, err
}

func FindLatestNotice(classID uint) (*model.ClassNotice, error) {
	var notice model.ClassNotice
	err := database.DB.Where("class_id = ?", classID).
		Order("created_at DESC").First(&notice).Error
	return &notice, err
}
