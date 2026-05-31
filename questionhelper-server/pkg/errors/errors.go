package errors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code string, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// 通用错误码
var (
	ErrParam       = New("A0001", "参数错误")
	ErrUnauthorized = New("A0200", "未授权")
	ErrForbidden   = New("A0301", "无权限")
	ErrNotFound    = New("A0002", "资源不存在")
	ErrInternal    = New("B0001", "系统内部错误")
)

// 用户相关错误码
var (
	ErrUserNotFound    = New("A0201", "用户不存在")
	ErrUserExists      = New("A0202", "用户已存在")
	ErrPasswordWrong   = New("A0203", "密码错误")
	ErrAccountDisabled = New("A0204", "账号已禁用")
)

// Token 相关错误码
var (
	ErrAccessTokenInvalid  = New("A0230", "访问令牌无效或过期")
	ErrRefreshTokenInvalid = New("A0231", "刷新令牌无效或过期")
)

// 题目相关错误码
var (
	ErrQuestionNotFound = New("A0401", "题目不存在")
	ErrCategoryNotFound = New("A0402", "分类不存在")
)

// 考试相关错误码
var (
	ErrExamNotFound     = New("A0501", "考试不存在")
	ErrExamNotStarted   = New("A0502", "考试未开始")
	ErrExamEnded        = New("A0503", "考试已结束")
	ErrAlreadySubmitted = New("A0504", "已经提交过")
)
