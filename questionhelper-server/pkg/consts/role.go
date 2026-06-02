package consts

// 系统角色编码（与需求文档 11.3 预置角色保持一致）
const (
	RoleSuperAdmin     = "super_admin"     // 超级管理员：系统最高权限
	RoleAdmin          = "admin"           // 管理员：系统管理
	RoleContentManager = "content_manager" // 内容管理员：管理公共题库
	RoleUserManager    = "user_manager"    // 用户管理员：管理用户账号
	RoleDataAnalyst    = "data_analyst"    // 数据分析师：查看统计
	RoleOperator       = "operator"        // 运营人员：管理公告推荐
	RoleTeacher        = "teacher"         // 教师：创建班级、发布考试
	RoleCreator        = "creator"         // 创作者：创建私有题目
	RoleStudent        = "student"         // 学生：使用公共题库、加入班级
)

// 班级成员角色（班级维度，非系统角色）
const (
	ClassRoleStudent = 1 // 学生
	ClassRoleTeacher = 2 // 教师
	ClassRoleAdmin   = 3 // 管理员
)
