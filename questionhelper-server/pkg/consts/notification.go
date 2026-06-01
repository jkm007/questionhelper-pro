package consts

// 通知类型
const (
	NotificationTypeSystem   int8 = 1 // 系统通知
	NotificationTypeExam     int8 = 2 // 考试通知
	NotificationTypeHomework int8 = 3 // 作业通知
	NotificationTypeClass    int8 = 4 // 班级通知
	NotificationTypeComment  int8 = 5 // 评论通知
)

// 通知渠道
const (
	NotificationChannelApp    = "app"    // 应用内通知
	NotificationChannelEmail  = "email"  // 邮件通知
	NotificationChannelSms    = "sms"    // 短信通知
	NotificationChannelWechat = "wechat" // 微信通知
)

// 定时通知状态
const (
	ScheduledStatusPending  int8 = 0 // 待发送
	ScheduledStatusSent     int8 = 1 // 已发送
	ScheduledStatusFailed   int8 = 2 // 发送失败
)

// 通知已读状态
const (
	NotificationUnread = false // 未读
	NotificationRead   = true  // 已读
)

// 群发通知目标类型
const (
	BatchTargetAll   = "all"   // 全部用户
	BatchTargetRole  = "role"  // 按角色
	BatchTargetClass = "class" // 按班级
	BatchTargetGroup = "group" // 按用户组
)
