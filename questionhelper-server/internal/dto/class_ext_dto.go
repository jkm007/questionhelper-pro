package dto

import "time"

// ==================== Homework Extension ====================

// UpdateHomeworkRequest 更新作业请求
type UpdateHomeworkRequest struct {
	Title       string    `json:"title" binding:"omitempty,max=200"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

// SubmitHomeworkRequest 提交作业请求
type SubmitHomeworkRequest struct {
	Content     string   `json:"content" binding:"required"`
	Attachments []string `json:"attachments"`
}

// GradeHomeworkRequest 批改作业请求
type GradeHomeworkRequest struct {
	Score    float64 `json:"score" binding:"required,min=0"`
	Feedback string  `json:"feedback"`
	Status   int8    `json:"status" binding:"required,oneof=1 2"` // 1=已批改,2=退回重做
}

// HomeworkSubmissionInfo 作业提交信息
type HomeworkSubmissionInfo struct {
	ID          uint       `json:"id"`
	HomeworkID  uint       `json:"homework_id"`
	UserID      uint       `json:"user_id"`
	UserName    string     `json:"user_name"`
	Content     string     `json:"content"`
	Attachments string     `json:"attachments"`
	Score       *float64   `json:"score"`
	Feedback    string     `json:"feedback"`
	Status      int8       `json:"status"`
	SubmittedAt time.Time  `json:"submitted_at"`
	GradedAt    *time.Time `json:"graded_at"`
}

// AssignPeerReviewRequest 分配互评请求
type AssignPeerReviewRequest struct {
	TargetType int8   `json:"target_type" binding:"required"` // 1=作业
	TargetID   uint   `json:"target_id" binding:"required"`
	Strategy   string `json:"strategy"` // random=随机,assign=指定
	Pairs      []PeerReviewPair `json:"pairs"` // 指定配对
}

// PeerReviewPair 互评配对
type PeerReviewPair struct {
	ReviewerID uint `json:"reviewer_id" binding:"required"`
	RevieweeID uint `json:"reviewee_id" binding:"required"`
}

// SubmitPeerReviewRequest 提交互评请求
type SubmitPeerReviewRequest struct {
	Score   float64 `json:"score" binding:"required,min=0,max=100"`
	Content string  `json:"content" binding:"required"`
}

// PeerReviewInfo 互评信息
type PeerReviewInfo struct {
	ID         uint       `json:"id"`
	TargetType int8       `json:"target_type"`
	TargetID   uint       `json:"target_id"`
	ReviewerID uint       `json:"reviewer_id"`
	ReviewerName string   `json:"reviewer_name"`
	RevieweeID uint       `json:"reviewee_id"`
	RevieweeName string   `json:"reviewee_name"`
	Score      *float64   `json:"score"`
	Content    string     `json:"content"`
	Status     int8       `json:"status"`
	ReviewedAt *time.Time `json:"reviewed_at"`
}

// PeerReviewResult 互评结果
type PeerReviewResult struct {
	UserID       uint    `json:"user_id"`
	UserName     string  `json:"user_name"`
	AvgScore     float64 `json:"avg_score"`
	ReviewCount  int     `json:"review_count"`
	Reviews      []PeerReviewInfo `json:"reviews"`
}

// ==================== Group Extension ====================

// CreateGroupRequest 创建分组请求
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	LeaderID    *uint  `json:"leader_id"`
	MaxMembers  int    `json:"max_members"`
}

// UpdateGroupRequest 更新分组请求
type UpdateGroupRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	LeaderID    *uint  `json:"leader_id"`
	MaxMembers  int    `json:"max_members"`
	SortOrder   int    `json:"sort_order"`
}

// GroupInfo 分组信息
type GroupInfo struct {
	ID          uint             `json:"id"`
	ClassID     uint             `json:"class_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	LeaderID    *uint            `json:"leader_id"`
	MaxMembers  int              `json:"max_members"`
	MemberCount int              `json:"member_count"`
	Members     []GroupMemberInfo `json:"members,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

// GroupMemberInfo 分组成员信息
type GroupMemberInfo struct {
	ID       uint      `json:"id"`
	GroupID  uint      `json:"group_id"`
	UserID   uint      `json:"user_id"`
	UserName string    `json:"user_name"`
	Role     int8      `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// AddGroupMemberRequest 添加分组成员请求
type AddGroupMemberRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required"`
}

// ==================== Attendance Extension ====================

// CreateAttendanceRequest 创建考勤请求
type CreateAttendanceRequest struct {
	Title       string     `json:"title" binding:"required,max=200"`
	Description string     `json:"description" binding:"omitempty,max=500"`
	Type        int8       `json:"type" binding:"required,oneof=1 2 3 4"`
	Deadline    *time.Time `json:"deadline"`
}

// UpdateAttendanceRequest 更新考勤请求
type UpdateAttendanceRequest struct {
	Title       string     `json:"title" binding:"omitempty,max=200"`
	Description string     `json:"description" binding:"omitempty,max=500"`
	Type        int8       `json:"type" binding:"omitempty,oneof=1 2 3 4"`
	Deadline    *time.Time `json:"deadline"`
	Status      int8       `json:"status" binding:"omitempty,oneof=0 1"`
}

// AttendanceInfo 考勤信息
type AttendanceInfo struct {
	ID          uint       `json:"id"`
	ClassID     uint       `json:"class_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        int8       `json:"type"`
	Deadline    *time.Time `json:"deadline"`
	CreatorID   uint       `json:"creator_id"`
	Status      int8       `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

// AttendanceRecordInfo 考勤记录信息
type AttendanceRecordInfo struct {
	ID           uint      `json:"id"`
	AttendanceID uint      `json:"attendance_id"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	Status       int8      `json:"status"`
	Remark       string    `json:"remark"`
	IP           string    `json:"ip"`
	Location     string    `json:"location"`
	SignedAt     time.Time `json:"signed_at"`
}

// ==================== Study Plan Extension ====================

// CreateStudyPlanRequest 创建学习计划请求
type CreateStudyPlanRequest struct {
	Title       string    `json:"title" binding:"required,max=200"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

// UpdateStudyPlanRequest 更新学习计划请求
type UpdateStudyPlanRequest struct {
	Title       string    `json:"title" binding:"omitempty,max=200"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      int8      `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// StudyPlanInfo 学习计划信息
type StudyPlanInfo struct {
	ID          uint               `json:"id"`
	ClassID     uint               `json:"class_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	StartDate   time.Time          `json:"start_date"`
	EndDate     time.Time          `json:"end_date"`
	CreatorID   uint               `json:"creator_id"`
	Status      int8               `json:"status"`
	Items       []StudyPlanItemInfo `json:"items,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

// CreateStudyPlanItemRequest 创建学习计划任务请求
type CreateStudyPlanItemRequest struct {
	Title        string     `json:"title" binding:"required,max=200"`
	Description  string     `json:"description"`
	ItemType     int8       `json:"item_type" binding:"required"`
	ResourceType string     `json:"resource_type"`
	ResourceID   uint       `json:"resource_id"`
	TargetCount  int        `json:"target_count"`
	Required     bool       `json:"required"`
	DueDate      *time.Time `json:"due_date"`
}

// UpdateStudyPlanItemRequest 更新学习计划任务请求
type UpdateStudyPlanItemRequest struct {
	Title        string     `json:"title" binding:"omitempty,max=200"`
	Description  string     `json:"description"`
	ItemType     int8       `json:"item_type"`
	ResourceType string     `json:"resource_type"`
	ResourceID   uint       `json:"resource_id"`
	TargetCount  int        `json:"target_count"`
	Required     *bool      `json:"required"`
	DueDate      *time.Time `json:"due_date"`
}

// StudyPlanItemInfo 学习计划任务信息
type StudyPlanItemInfo struct {
	ID           uint       `json:"id"`
	PlanID       uint       `json:"plan_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ItemType     int8       `json:"item_type"`
	ResourceType string     `json:"resource_type"`
	ResourceID   uint       `json:"resource_id"`
	TargetCount  int        `json:"target_count"`
	Required     bool       `json:"required"`
	SortOrder    int        `json:"sort_order"`
	DueDate      *time.Time `json:"due_date"`
	CreatedAt    time.Time  `json:"created_at"`
}

// StudyPlanProgressInfo 学习计划进度信息
type StudyPlanProgressInfo struct {
	TotalItems    int                      `json:"total_items"`
	CompletedItems int                    `json:"completed_items"`
	Progress      float64                  `json:"progress"` // 百分比
	UserProgress  []UserStudyPlanProgress  `json:"user_progress,omitempty"`
}

// UserStudyPlanProgress 用户学习计划进度
type UserStudyPlanProgress struct {
	UserID      uint    `json:"user_id"`
	UserName    string  `json:"user_name"`
	TotalItems  int     `json:"total_items"`
	Completed   int     `json:"completed"`
	Progress    float64 `json:"progress"`
}

// ==================== File Extension ====================

// ClassFileInfo 班级文件信息
type ClassFileInfo struct {
	ID        uint      `json:"id"`
	ClassID   uint      `json:"class_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	FolderID  *uint     `json:"folder_id"`
	CreatorID uint      `json:"creator_id"`
	Downloads int       `json:"downloads"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== Ranking Extension ====================

// ClassRankingInfo 班级排名信息
type ClassRankingInfo struct {
	UserID      uint    `json:"user_id"`
	UserName    string  `json:"user_name"`
	RankingType string  `json:"ranking_type"`
	Score       float64 `json:"score"`
	Rank        int     `json:"rank"`
}

// CalculateRankingRequest 计算排名请求
type CalculateRankingRequest struct {
	RankingType string `json:"ranking_type" binding:"required,oneof=practice exam homework overall"`
}

// ==================== Tag Extension ====================

// AddClassTagRequest 为班级添加标签请求
type AddClassTagRequest struct {
	TagIDs []uint `json:"tag_ids" binding:"required"`
}

// ==================== Template Extension ====================

// CreateClassFromTemplateRequest 从模板创建班级请求
type CreateClassFromTemplateRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

// ==================== Discussion Extension ====================

// CreateDiscussionRequest 创建讨论请求
type CreateDiscussionRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
}

// UpdateDiscussionRequest 更新讨论请求
type UpdateDiscussionRequest struct {
	Title   string `json:"title" binding:"omitempty,max=200"`
	Content string `json:"content"`
}

// DiscussionInfo 讨论信息
type DiscussionInfo struct {
	ID          uint       `json:"id"`
	ClassID     uint       `json:"class_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	CreatorID   uint       `json:"creator_id"`
	CreatorName string     `json:"creator_name"`
	ReplyCount  int        `json:"reply_count"`
	ViewCount   int        `json:"view_count"`
	IsTop       bool       `json:"is_top"`
	IsClosed    bool       `json:"is_closed"`
	LastReplyAt *time.Time `json:"last_reply_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ==================== Resource Extension ====================

// ResourceInfo 资源信息
type ResourceInfo struct {
	ID           uint      `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	Title        string    `json:"title"`
	Version      int       `json:"version"`
	CreatorID    uint      `json:"creator_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// ResourceStatistics 资源统计
type ResourceStatistics struct {
	FileCount      int64 `json:"file_count"`
	HomeworkCount  int64 `json:"homework_count"`
	NoticeCount    int64 `json:"notice_count"`
	TotalSize      int64 `json:"total_size"`
}

// ImportResourceRequest 导入资源请求
type ImportResourceRequest struct {
	ResourceType string `json:"resource_type" binding:"required"`
	Data         string `json:"data" binding:"required"`
}

// ==================== Creator Extension ====================

// CreatorApplyRequest 申请创作者请求
type CreatorApplyRequest struct {
	Reason     string `json:"reason" binding:"required,max=200"`
	MaxClasses int    `json:"max_classes"`
}

// CreatorApplicationInfo 创作者申请信息
type CreatorApplicationInfo struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	UserName   string     `json:"user_name"`
	Reason     string     `json:"reason"`
	MaxClasses int        `json:"max_classes"`
	Status     int8       `json:"status"`
	Remark     string     `json:"remark"`
	ReviewBy   *uint      `json:"review_by"`
	ReviewAt   *time.Time `json:"review_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreatorInfo 创作者信息
type CreatorInfo struct {
	UserID     uint       `json:"user_id"`
	UserName   string     `json:"user_name"`
	MaxClasses int        `json:"max_classes"`
	CanCreate  bool       `json:"can_create"`
	ExpiresAt  *time.Time `json:"expires_at"`
	GrantedBy  uint       `json:"granted_by"`
	Reason     string     `json:"reason"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ==================== Management Enhancement ====================

// SetExpireRequest 设置有效期请求
type SetExpireRequest struct {
	ExpireAt time.Time `json:"expire_at" binding:"required"`
}

// QRCodeInfo 二维码信息
type QRCodeInfo struct {
	Code    string `json:"code"`
	URL     string `json:"url"`
	Image   string `json:"image"` // base64
}

// ClassExamInfo 班级考试信息
type ClassExamInfo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    int8      `json:"status"`
}

// ClassNoticeDetailInfo 班级公告详情
type ClassNoticeDetailInfo struct {
	ID          uint      `json:"id"`
	ClassID     uint      `json:"class_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatorID   uint      `json:"creator_id"`
	CreatorName string    `json:"creator_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// SearchClassRequest 搜索班级请求
type SearchClassRequest struct {
	PageRequest
	Keyword string `form:"keyword" binding:"required"`
}
