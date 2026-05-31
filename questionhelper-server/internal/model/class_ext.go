package model

import (
	"time"

	"gorm.io/gorm"
)

// HomeworkSubmission 作业提交表，记录学生提交作业的情况
type HomeworkSubmission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	HomeworkID  uint           `gorm:"index;not null;comment:作业ID" json:"homework_id"`
	Homework    Homework       `gorm:"foreignKey:HomeworkID" json:"homework,omitempty"`
	UserID      uint           `gorm:"index;not null;comment:提交学生ID" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content     string         `gorm:"type:text;comment:提交内容" json:"content"`
	Attachments string         `gorm:"type:text;comment:附件路径(JSON数组)" json:"attachments"`
	Score       *float64       `gorm:"comment:得分" json:"score"`
	Feedback    string         `gorm:"type:text;comment:教师评语" json:"feedback"`
	GradedBy    *uint          `gorm:"index;comment:批改人ID" json:"graded_by"`
	Status      int8           `gorm:"default:0;comment:状态:0=未批改,1=已批改,2=退回重做" json:"status"`
	SubmittedAt time.Time      `gorm:"index;not null;comment:提交时间" json:"submitted_at"`
	GradedAt    *time.Time     `gorm:"comment:批改时间" json:"graded_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HomeworkSubmission) TableName() string {
	return "homework_submissions"
}

// ClassTag 班级标签表，用于对班级进行分类标记
type ClassTag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null;comment:标签名称" json:"name"`
	CreatorID uint           `gorm:"index;comment:创建人ID" json:"creator_id"`
	SortOrder int            `gorm:"default:0;comment:排序值" json:"sort_order"`
	Status    int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassTag) TableName() string {
	return "class_tags"
}

// ClassTagRelation 班级标签关联表，记录班级与标签的多对多关系
type ClassTagRelation struct {
	ID      uint      `gorm:"primarykey" json:"id"`
	ClassID uint      `gorm:"uniqueIndex:idx_class_tag;not null;comment:班级ID" json:"class_id"`
	TagID   uint      `gorm:"uniqueIndex:idx_class_tag;not null;comment:标签ID" json:"tag_id"`
}

func (ClassTagRelation) TableName() string {
	return "class_tag_relations"
}

// ClassGroup 班级分组表，支持班级内的小组协作学习
type ClassGroup struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ClassID     uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	Name        string         `gorm:"size:100;not null;comment:分组名称" json:"name"`
	Description string         `gorm:"size:500;comment:分组描述" json:"description"`
	LeaderID    *uint          `gorm:"index;comment:组长用户ID" json:"leader_id"`
	MaxMembers  int            `gorm:"default:0;comment:最大成员数(0=不限制)" json:"max_members"`
	SortOrder   int            `gorm:"default:0;comment:排序值" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassGroup) TableName() string {
	return "class_groups"
}

// ClassGroupMember 班级分组成员表，记录分组内的学生
type ClassGroupMember struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	GroupID  uint      `gorm:"uniqueIndex:idx_group_user;not null;comment:分组ID" json:"group_id"`
	UserID   uint      `gorm:"uniqueIndex:idx_group_user;not null;comment:用户ID" json:"user_id"`
	Role     int8      `gorm:"default:1;comment:角色:1=普通成员,2=组长" json:"role"`
	JoinedAt time.Time `gorm:"not null;comment:加入时间" json:"joined_at"`
}

func (ClassGroupMember) TableName() string {
	return "class_group_members"
}

// ClassJoinApplication 班级加入申请表，记录学生申请加入班级的请求
type ClassJoinApplication struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ClassID   uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	UserID    uint           `gorm:"index;not null;comment:申请人ID" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reason    string         `gorm:"size:500;comment:申请理由" json:"reason"`
	Status    int8           `gorm:"default:0;comment:状态:0=待审批,1=已通过,2=已拒绝" json:"status"`
	Remark    string         `gorm:"size:200;comment:审批备注" json:"remark"`
	ReviewBy  *uint          `gorm:"index;comment:审批人ID" json:"review_by"`
	ReviewAt  *time.Time     `gorm:"comment:审批时间" json:"review_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassJoinApplication) TableName() string {
	return "class_join_applications"
}

// ClassAttendance 考勤表，定义一次考勤活动
type ClassAttendance struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ClassID     uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	Title       string         `gorm:"size:200;not null;comment:考勤标题" json:"title"`
	Description string         `gorm:"size:500;comment:考勤说明" json:"description"`
	Type        int8           `gorm:"default:1;comment:类型:1=普通签到,2=限时签到,3=位置签到,4=手势签到" json:"type"`
	Deadline    *time.Time     `gorm:"comment:签到截止时间" json:"deadline"`
	CreatorID   uint           `gorm:"index;comment:发起人ID" json:"creator_id"`
	Status      int8           `gorm:"default:1;comment:状态:0=已结束,1=进行中" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassAttendance) TableName() string {
	return "class_attendances"
}

// ClassAttendanceRecord 考勤记录表，记录每个学生的签到情况
type ClassAttendanceRecord struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	AttendanceID uint       `gorm:"uniqueIndex:idx_attendance_user;not null;comment:考勤ID" json:"attendance_id"`
	UserID       uint       `gorm:"uniqueIndex:idx_attendance_user;not null;comment:用户ID" json:"user_id"`
	Status       int8       `gorm:"default:1;comment:状态:1=正常,2=迟到,3=缺勤,4=请假" json:"status"`
	Remark       string     `gorm:"size:200;comment:备注" json:"remark"`
	IP           string     `gorm:"size:45;comment:签到IP" json:"ip"`
	Location     string     `gorm:"size:200;comment:签到位置信息" json:"location"`
	SignedAt     time.Time  `gorm:"not null;comment:签到时间" json:"signed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (ClassAttendanceRecord) TableName() string {
	return "class_attendance_records"
}

// ClassStudyPlan 学习计划表，教师为班级制定的学习计划
type ClassStudyPlan struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ClassID     uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	Title       string         `gorm:"size:200;not null;comment:计划标题" json:"title"`
	Description string         `gorm:"type:text;comment:计划描述" json:"description"`
	StartDate   time.Time      `gorm:"not null;comment:开始日期" json:"start_date"`
	EndDate     time.Time      `gorm:"not null;comment:结束日期" json:"end_date"`
	CreatorID   uint           `gorm:"index;comment:创建人ID" json:"creator_id"`
	Status      int8           `gorm:"default:1;comment:状态:0=已归档,1=进行中,2=未开始" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassStudyPlan) TableName() string {
	return "class_study_plans"
}

// ClassStudyPlanItem 学习计划任务表，定义计划中的具体学习任务
type ClassStudyPlanItem struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	PlanID        uint           `gorm:"index;not null;comment:学习计划ID" json:"plan_id"`
	Plan          ClassStudyPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Title         string         `gorm:"size:200;not null;comment:任务标题" json:"title"`
	Description   string         `gorm:"type:text;comment:任务描述" json:"description"`
	ItemType      int8           `gorm:"default:1;comment:任务类型:1=做题,2=考试,3=阅读资料,4=观看视频,5=其他" json:"item_type"`
	ResourceType  string         `gorm:"size:50;comment:资源类型" json:"resource_type"`
	ResourceID    uint           `gorm:"index;comment:资源ID" json:"resource_id"`
	TargetCount   int            `gorm:"default:0;comment:目标数量(如题目数)" json:"target_count"`
	Required      bool           `gorm:"default:true;comment:是否必做" json:"required"`
	SortOrder     int            `gorm:"default:0;comment:排序值" json:"sort_order"`
	DueDate       *time.Time     `gorm:"comment:截止日期" json:"due_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassStudyPlanItem) TableName() string {
	return "class_study_plan_items"
}

// ClassStudyPlanProgress 学习计划进度表，记录学生完成计划任务的进度
type ClassStudyPlanProgress struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	PlanID      uint      `gorm:"index;not null;comment:学习计划ID" json:"plan_id"`
	ItemID      uint      `gorm:"uniqueIndex:idx_item_user;not null;comment:任务ID" json:"item_id"`
	UserID      uint      `gorm:"uniqueIndex:idx_item_user;not null;comment:用户ID" json:"user_id"`
	Completed   int       `gorm:"default:0;comment:已完成数量" json:"completed"`
	Status      int8      `gorm:"default:0;comment:状态:0=未开始,1=进行中,2=已完成" json:"status"`
	CompletedAt *time.Time `gorm:"comment:完成时间" json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ClassStudyPlanProgress) TableName() string {
	return "class_study_plan_progress"
}

// ClassFile 班级文件表，管理班级共享的教学资源文件
type ClassFile struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ClassID   uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	Name      string         `gorm:"size:255;not null;comment:文件名称" json:"name"`
	Path      string         `gorm:"size:500;not null;comment:文件路径" json:"path"`
	Size      int64          `gorm:"default:0;comment:文件大小(字节)" json:"size"`
	MimeType  string         `gorm:"size:100;comment:MIME类型" json:"mime_type"`
	FolderID  *uint          `gorm:"index;comment:所属文件夹ID" json:"folder_id"`
	CreatorID uint           `gorm:"index;comment:上传人ID" json:"creator_id"`
	Downloads int            `gorm:"default:0;comment:下载次数" json:"downloads"`
	SortOrder int            `gorm:"default:0;comment:排序值" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassFile) TableName() string {
	return "class_files"
}

// ClassPeerReview 学生互评表，支持学生之间互相评价作业或表现
type ClassPeerReview struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	TargetType int8           `gorm:"not null;comment:评价目标类型:1=作业,2=讨论,3=其他" json:"target_type"`
	TargetID   uint           `gorm:"index;not null;comment:评价目标ID" json:"target_id"`
	ReviewerID uint           `gorm:"uniqueIndex:idx_review_pair;not null;comment:评价人ID" json:"reviewer_id"`
	RevieweeID uint           `gorm:"uniqueIndex:idx_review_pair;not null;comment:被评价人ID" json:"reviewee_id"`
	Score      *float64       `gorm:"comment:评分" json:"score"`
	Content    string         `gorm:"type:text;comment:评价内容" json:"content"`
	Status     int8           `gorm:"default:0;comment:状态:0=待评,1=已评" json:"status"`
	ReviewedAt *time.Time     `gorm:"comment:评价时间" json:"reviewed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassPeerReview) TableName() string {
	return "class_peer_reviews"
}

// ClassRanking 班级排名缓存表，缓存学生的班级内排名数据
type ClassRanking struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ClassID     uint      `gorm:"uniqueIndex:idx_class_user_type;not null;comment:班级ID" json:"class_id"`
	UserID      uint      `gorm:"uniqueIndex:idx_class_user_type;not null;comment:用户ID" json:"user_id"`
	RankingType string    `gorm:"uniqueIndex:idx_class_user_type;size:30;not null;comment:排名类型:practice=练习,exam=考试,homework=作业,overall=综合" json:"ranking_type"`
	Score       float64   `gorm:"default:0;comment:综合得分" json:"score"`
	Rank        int       `gorm:"default:0;comment:排名" json:"rank"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ClassRanking) TableName() string {
	return "class_rankings"
}

// ClassImportRecord 批量导入记录表，记录批量导入学生等操作的日志
type ClassImportRecord struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ClassID      uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	ImportType   string         `gorm:"size:30;not null;comment:导入类型:student=学生,question=题目,homework=作业" json:"import_type"`
	FileName     string         `gorm:"size:255;comment:导入文件名" json:"file_name"`
	TotalCount   int            `gorm:"default:0;comment:总记录数" json:"total_count"`
	SuccessCount int            `gorm:"default:0;comment:成功数" json:"success_count"`
	FailCount    int            `gorm:"default:0;comment:失败数" json:"fail_count"`
	ErrorDetail  string         `gorm:"type:text;comment:失败详情(JSON)" json:"error_detail"`
	OperatorID   uint           `gorm:"index;comment:操作人ID" json:"operator_id"`
	Status       int8           `gorm:"default:0;comment:状态:0=处理中,1=已完成,2=失败" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassImportRecord) TableName() string {
	return "class_import_records"
}

// ClassTemplate 班级模板表，保存班级配置模板以便快速创建班级
type ClassTemplate struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null;comment:模板名称" json:"name"`
	Description string         `gorm:"size:500;comment:模板描述" json:"description"`
	Config      string         `gorm:"type:text;comment:模板配置(JSON)" json:"config"`
	CreatorID   uint           `gorm:"index;comment:创建人ID" json:"creator_id"`
	IsPublic    bool           `gorm:"default:false;comment:是否公开" json:"is_public"`
	UsedCount   int            `gorm:"default:0;comment:使用次数" json:"used_count"`
	Status      int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassTemplate) TableName() string {
	return "class_templates"
}

// ClassDiscussion 班级讨论表，支持班级内的主题讨论
type ClassDiscussion struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ClassID     uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	Title       string         `gorm:"size:200;not null;comment:讨论标题" json:"title"`
	Content     string         `gorm:"type:text;not null;comment:讨论内容" json:"content"`
	CreatorID   uint           `gorm:"index;comment:发起人ID" json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	ReplyCount  int            `gorm:"default:0;comment:回复数" json:"reply_count"`
	ViewCount   int            `gorm:"default:0;comment:浏览数" json:"view_count"`
	IsTop       bool           `gorm:"default:false;comment:是否置顶" json:"is_top"`
	IsClosed    bool           `gorm:"default:false;comment:是否已关闭" json:"is_closed"`
	LastReplyAt *time.Time     `gorm:"index;comment:最后回复时间" json:"last_reply_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassDiscussion) TableName() string {
	return "class_discussions"
}

// ClassDiscussionReply 班级讨论回复表，记录讨论下的回复内容
type ClassDiscussionReply struct {
	ID           uint             `gorm:"primarykey" json:"id"`
	DiscussionID uint             `gorm:"index;not null;comment:讨论ID" json:"discussion_id"`
	Discussion   ClassDiscussion  `gorm:"foreignKey:DiscussionID" json:"discussion,omitempty"`
	ParentID     *uint            `gorm:"index;comment:父回复ID(楼中楼)" json:"parent_id"`
	Content      string           `gorm:"type:text;not null;comment:回复内容" json:"content"`
	UserID       uint             `gorm:"index;not null;comment:回复人ID" json:"user_id"`
	User         User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LikeCount    int              `gorm:"default:0;comment:点赞数" json:"like_count"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (ClassDiscussionReply) TableName() string {
	return "class_discussion_replies"
}

// HomeworkTemplate 作业模板表，保存常用作业配置以便快速创建作业
type HomeworkTemplate struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"size:200;not null;comment:模板标题" json:"title"`
	Description string         `gorm:"type:text;comment:模板描述" json:"description"`
	Content     string         `gorm:"type:text;comment:作业内容模板" json:"content"`
	Subject     string         `gorm:"size:50;index;comment:学科" json:"subject"`
	GradeLevel  string         `gorm:"size:20;comment:年级" json:"grade_level"`
	CreatorID   uint           `gorm:"index;comment:创建人ID" json:"creator_id"`
	IsPublic    bool           `gorm:"default:false;comment:是否公开" json:"is_public"`
	UsedCount   int            `gorm:"default:0;comment:使用次数" json:"used_count"`
	Status      int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HomeworkTemplate) TableName() string {
	return "homework_templates"
}

// ScoreAlert 成绩预警表，当学生成绩低于阈值时触发预警
type ScoreAlert struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ClassID     uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	UserID      uint           `gorm:"index;not null;comment:学生ID" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AlertType   string         `gorm:"size:30;not null;comment:预警类型:score=单次成绩,average=平均分,decline=成绩下滑,absent=缺考" json:"alert_type"`
	SourceType  string         `gorm:"size:30;not null;comment:来源类型:homework=作业,exam=考试,practice=练习" json:"source_type"`
	SourceID    uint           `gorm:"index;comment:来源ID" json:"source_id"`
	Threshold   float64        `gorm:"comment:预警阈值" json:"threshold"`
	ActualValue float64        `gorm:"comment:实际值" json:"actual_value"`
	Message     string         `gorm:"size:500;comment:预警信息" json:"message"`
	IsRead      bool           `gorm:"default:false;comment:是否已读" json:"is_read"`
	HandledBy   *uint          `gorm:"index;comment:处理人ID" json:"handled_by"`
	HandledAt   *time.Time     `gorm:"comment:处理时间" json:"handled_at"`
	HandleNote  string         `gorm:"size:500;comment:处理备注" json:"handle_note"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ScoreAlert) TableName() string {
	return "score_alerts"
}

// ClassLog 班级操作日志表，记录班级管理的操作历史
type ClassLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ClassID     uint      `gorm:"index;not null;comment:班级ID" json:"class_id"`
	OperatorID  uint      `gorm:"index;not null;comment:操作人ID" json:"operator_id"`
	Action      string    `gorm:"size:50;not null;index;comment:操作类型" json:"action"`
	TargetType  string    `gorm:"size:50;comment:目标类型" json:"target_type"`
	TargetID    uint      `gorm:"comment:目标ID" json:"target_id"`
	Detail      string    `gorm:"type:text;comment:操作详情(JSON)" json:"detail"`
	IP          string    `gorm:"size:45;comment:操作IP" json:"ip"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ClassLog) TableName() string {
	return "class_logs"
}

// ClassResourceVersion 班级资源版本表，管理班级教学资源的版本历史
type ClassResourceVersion struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ResourceType string   `gorm:"size:50;not null;index;comment:资源类型:file=文件,homework=作业,notice=通知" json:"resource_type"`
	ResourceID  uint      `gorm:"not null;index;comment:资源ID" json:"resource_id"`
	Version     int       `gorm:"default:1;comment:版本号" json:"version"`
	Title       string    `gorm:"size:200;comment:版本标题" json:"title"`
	Content     string    `gorm:"type:text;comment:版本内容快照" json:"content"`
	FilePath    string    `gorm:"size:500;comment:文件路径" json:"file_path"`
	CreatorID   uint      `gorm:"index;comment:创建人ID" json:"creator_id"`
	Remark      string    `gorm:"size:200;comment:版本说明" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ClassResourceVersion) TableName() string {
	return "class_resource_versions"
}

// ClassResourceTag 班级资源标签表，用于对班级资源进行分类
type ClassResourceTag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ClassID   uint           `gorm:"index;not null;comment:班级ID" json:"class_id"`
	Name      string         `gorm:"size:50;not null;comment:标签名称" json:"name"`
	CreatorID uint           `gorm:"index;comment:创建人ID" json:"creator_id"`
	SortOrder int            `gorm:"default:0;comment:排序值" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassResourceTag) TableName() string {
	return "class_resource_tags"
}

// ClassResourceTagRelation 班级资源标签关联表，记录资源与标签的多对多关系
type ClassResourceTagRelation struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	TagID        uint      `gorm:"uniqueIndex:idx_tag_resource;not null;comment:标签ID" json:"tag_id"`
	ResourceType string    `gorm:"uniqueIndex:idx_tag_resource;size:50;not null;comment:资源类型" json:"resource_type"`
	ResourceID   uint      `gorm:"uniqueIndex:idx_tag_resource;not null;comment:资源ID" json:"resource_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (ClassResourceTagRelation) TableName() string {
	return "class_resource_tag_relations"
}

// ClassResourceReview 班级资源审核表，记录教学资源的审核流程
type ClassResourceReview struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ResourceType string         `gorm:"size:50;not null;index;comment:资源类型" json:"resource_type"`
	ResourceID   uint           `gorm:"not null;index;comment:资源ID" json:"resource_id"`
	SubmitterID  uint           `gorm:"index;not null;comment:提交人ID" json:"submitter_id"`
	ReviewerID   *uint          `gorm:"index;comment:审核人ID" json:"reviewer_id"`
	Status       int8           `gorm:"default:0;comment:状态:0=待审核,1=已通过,2=已拒绝" json:"status"`
	Remark       string         `gorm:"size:500;comment:审核备注" json:"remark"`
	ReviewedAt   *time.Time     `gorm:"comment:审核时间" json:"reviewed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassResourceReview) TableName() string {
	return "class_resource_reviews"
}

// ClassCreatorPermission 班级创作者权限表，控制用户创建班级的权限
type ClassCreatorPermission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"uniqueIndex;not null;comment:用户ID" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MaxClasses  int            `gorm:"default:5;comment:最大可创建班级数" json:"max_classes"`
	CanCreate   bool           `gorm:"default:true;comment:是否允许创建班级" json:"can_create"`
	ExpiresAt   *time.Time     `gorm:"comment:权限过期时间(空=永久)" json:"expires_at"`
	GrantedBy   uint           `gorm:"index;comment:授权人ID" json:"granted_by"`
	Reason      string         `gorm:"size:200;comment:授权原因" json:"reason"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassCreatorPermission) TableName() string {
	return "class_creator_permissions"
}
