-- =====================================================
-- 题小助管理后台菜单种子数据
-- 与管理后台功能设计文档保持一致
-- =====================================================

-- 清空现有菜单数据（可选，谨慎使用）
-- TRUNCATE TABLE menus;

-- =====================================================
-- 一级菜单（目录）
-- =====================================================

-- 1. 首页
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(1, NULL, 'Dashboard', '/dashboard', 'Layout', '/dashboard', '首页', 'HomeFilled', 0, 1, '', 1, 1, NOW(), NOW());

-- 2. 用户管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(2, NULL, 'User', '/user', 'Layout', '', '用户管理', 'User', 0, 1, '', 2, 1, NOW(), NOW());

-- 3. 权限管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(3, NULL, 'Permission', '/permission', 'Layout', '', '权限管理', 'Lock', 0, 1, '', 3, 1, NOW(), NOW());

-- 4. 题库管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(4, NULL, 'Question', '/question', 'Layout', '', '题库管理', 'EditPen', 0, 1, '', 4, 1, NOW(), NOW());

-- 5. 考试管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(5, NULL, 'Exam', '/exam', 'Layout', '', '考试管理', 'Timer', 0, 1, '', 5, 1, NOW(), NOW());

-- 6. 班级管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(6, NULL, 'Class', '/class', 'Layout', '', '班级管理', 'UserFilled', 0, 1, '', 6, 1, NOW(), NOW());

-- 7. 内容管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(7, NULL, 'Content', '/content', 'Layout', '', '内容管理', 'Document', 0, 1, '', 7, 1, NOW(), NOW());

-- 8. 数据统计
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(8, NULL, 'Statistics', '/statistics', 'Layout', '', '数据统计', 'DataAnalysis', 0, 1, '', 8, 1, NOW(), NOW());

-- 9. 消息中心
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(9, NULL, 'Message', '/message', 'Layout', '', '消息中心', 'Message', 0, 1, '', 9, 1, NOW(), NOW());

-- 10. 文件管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(10, NULL, 'File', '/file', 'Layout', '', '文件管理', 'Folder', 0, 1, '', 10, 1, NOW(), NOW());

-- 11. 系统设置
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(11, NULL, 'System', '/system', 'Layout', '', '系统设置', 'Setting', 0, 1, '', 11, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（用户管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(101, 2, 'UserList', '/user/list', 'system/user/index', '', '用户列表', 'User', 0, 2, 'user:list', 1, 1, NOW(), NOW()),
(102, 2, 'UserTag', '/user/tag', 'system/tag/index', '', '用户标签', 'PriceTag', 0, 2, 'user:tag', 2, 1, NOW(), NOW()),
(103, 2, 'UserLog', '/user/log', 'log/operation/index', '', '操作日志', 'Notebook', 0, 2, 'user:log', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（权限管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(201, 3, 'Role', '/permission/role', 'system/role/index', '', '角色管理', 'UserFilled', 0, 2, 'permission:role', 1, 1, NOW(), NOW()),
(202, 3, 'Menu', '/permission/menu', 'system/menu/index', '', '菜单管理', 'Menu', 0, 2, 'permission:menu', 2, 1, NOW(), NOW()),
(203, 3, 'Button', '/permission/button', 'permission/button/index', '', '按钮权限', 'Key', 0, 2, 'permission:button', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（题库管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(401, 4, 'QuestionList', '/question/list', 'question/list/index', '', '题目列表', 'List', 0, 2, 'question:list', 1, 1, NOW(), NOW()),
(402, 4, 'QuestionCreate', '/question/create', 'question/create/index', '', '题目创建', 'Plus', 0, 2, 'question:create', 2, 1, NOW(), NOW()),
(403, 4, 'QuestionCategory', '/question/category', 'question/category/index', '', '分类管理', 'Menu', 0, 2, 'question:category', 3, 1, NOW(), NOW()),
(404, 4, 'QuestionKnowledge', '/question/knowledge', 'question/knowledge/index', '', '知识点管理', 'Collection', 0, 2, 'question:knowledge', 4, 1, NOW(), NOW()),
(405, 4, 'QuestionReview', '/question/review', 'question/review/index', '', '题目审核', 'View', 0, 2, 'question:review', 5, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（考试管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(501, 5, 'ExamPaper', '/exam/paper', 'exam/paper/index', '', '试卷管理', 'Document', 0, 2, 'exam:paper', 1, 1, NOW(), NOW()),
(502, 5, 'ExamList', '/exam/list', 'exam/list/index', '', '考试列表', 'Tickets', 0, 2, 'exam:list', 2, 1, NOW(), NOW()),
(503, 5, 'ExamScore', '/exam/score', 'exam/score/index', '', '成绩管理', 'Trophy', 0, 2, 'exam:score', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（班级管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(601, 6, 'ClassList', '/class/list', 'class/list/index', '', '班级列表', 'School', 0, 2, 'class:list', 1, 1, NOW(), NOW()),
(602, 6, 'ClassAudit', '/class/audit', 'class/audit/index', '', '班级审核', 'Checked', 0, 2, 'class:audit', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（内容管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(701, 7, 'ContentComment', '/content/comment', 'content/comment/index', '', '评论管理', 'ChatDotRound', 0, 2, 'content:comment', 1, 1, NOW(), NOW()),
(702, 7, 'ContentReview', '/content/review', 'content/review/index', '', '资源审核', 'View', 0, 2, 'content:review', 2, 1, NOW(), NOW()),
(703, 7, 'ContentCreator', '/content/creator', 'content/creator/index', '', '创作者管理', 'Star', 0, 2, 'content:creator', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（数据统计）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(801, 8, 'StatisticsOverview', '/statistics/overview', 'statistics/overview/index', '', '数据概览', 'Odometer', 0, 2, 'statistics:overview', 1, 1, NOW(), NOW()),
(802, 8, 'StatisticsUser', '/statistics/user', 'statistics/user/index', '', '用户统计', 'Histogram', 0, 2, 'statistics:user', 2, 1, NOW(), NOW()),
(803, 8, 'StatisticsAnswer', '/statistics/answer', 'statistics/answer/index', '', '答题统计', 'PieChart', 0, 2, 'statistics:answer', 3, 1, NOW(), NOW()),
(804, 8, 'StatisticsScore', '/statistics/score', 'statistics/score/index', '', '成绩统计', 'TrendCharts', 0, 2, 'statistics:score', 4, 1, NOW(), NOW()),
(805, 8, 'StatisticsClass', '/statistics/class', 'statistics/class/index', '', '班级统计', 'DataBoard', 0, 2, 'statistics:class', 5, 1, NOW(), NOW()),
(806, 8, 'StatisticsAlert', '/statistics/alert', 'statistics/alert/index', '', '数据预警', 'Warning', 0, 2, 'statistics:alert', 6, 1, NOW(), NOW()),
(807, 8, 'StatisticsAdvanced', '/statistics/advanced', 'Layout', '', '高级分析', 'MagicStick', 0, 1, 'statistics:advanced', 7, 1, NOW(), NOW());

-- =====================================================
-- 三级菜单（高级分析）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(80701, 807, 'AdvancedRetention', '/statistics/advanced/retention', 'statistics/retention/index', '', '用户留存', 'DataLine', 0, 2, 'statistics:advanced:retention', 1, 1, NOW(), NOW()),
(80702, 807, 'AdvancedChurn', '/statistics/advanced/churn', 'statistics/advanced/churn/index', '', '用户流失', 'Remove', 0, 2, 'statistics:advanced:churn', 2, 1, NOW(), NOW()),
(80703, 807, 'AdvancedBehavior', '/statistics/advanced/behavior', 'statistics/advanced/behavior/index', '', '用户行为', 'Monitor', 0, 2, 'statistics:advanced:behavior', 3, 1, NOW(), NOW()),
(80704, 807, 'AdvancedSegment', '/statistics/advanced/segment', 'statistics/segment/index', '', '用户分群', 'PieChart', 0, 2, 'statistics:advanced:segment', 4, 1, NOW(), NOW()),
(80705, 807, 'AdvancedPath', '/statistics/advanced/path', 'statistics/path/index', '', '用户路径', 'Guide', 0, 2, 'statistics:advanced:path', 5, 1, NOW(), NOW()),
(80706, 807, 'AdvancedConversion', '/statistics/advanced/conversion', 'statistics/advanced/conversion/index', '', '用户转化', 'TrendCharts', 0, 2, 'statistics:advanced:conversion', 6, 1, NOW(), NOW()),
(80707, 807, 'AdvancedQuestion', '/statistics/advanced/question', 'statistics/advanced/question/index', '', '题目分析', 'DataAnalysis', 0, 2, 'statistics:advanced:question', 7, 1, NOW(), NOW()),
(80708, 807, 'AdvancedSubscribe', '/statistics/advanced/subscribe', 'statistics/advanced/subscribe/index', '', '数据订阅', 'Bell', 0, 2, 'statistics:advanced:subscribe', 8, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（消息中心）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(901, 9, 'MessageNotification', '/message/notification', 'message/notification/index', '', '通知管理', 'Bell', 0, 2, 'message:notification', 1, 1, NOW(), NOW()),
(902, 9, 'MessageTemplate', '/message/template', 'notification/template/index', '', '通知模板', 'Stamp', 0, 2, 'message:template', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（文件管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(1001, 10, 'FileList', '/file/list', 'file/list/index', '', '文件列表', 'Files', 0, 2, 'file:list', 1, 1, NOW(), NOW()),
(1002, 10, 'FileStorage', '/file/storage', 'file/storage/index', '', '存储配置', 'Box', 0, 2, 'file:storage', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（系统设置）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(1101, 11, 'SystemSetting', '/system/setting', 'setting/basic/index', '', '基础设置', 'Setting', 0, 2, 'system:setting', 1, 1, NOW(), NOW()),
(1102, 11, 'SystemSettingClass', '/system/setting/class', 'system/setting/class/index', '', '班级设置', 'School', 0, 2, 'system:setting:class', 2, 1, NOW(), NOW()),
(1103, 11, 'SystemSettingResource', '/system/setting/resource', 'system/setting/resource/index', '', '资源设置', 'Resources', 0, 2, 'system:setting:resource', 3, 1, NOW(), NOW()),
(1104, 11, 'SystemLog', '/system/log', 'log/system/index', '', '系统日志', 'Notebook', 0, 2, 'system:log', 4, 1, NOW(), NOW()),
(1105, 11, 'SystemBackup', '/system/backup', 'backup/list/index', '', '数据备份', 'Upload', 0, 2, 'system:backup', 5, 1, NOW(), NOW()),
(1106, 11, 'SystemApproval', '/system/approval', 'system/approval/index', '', '审批管理', 'Stamp', 0, 2, 'system:approval', 6, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（用户管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(10101, 101, 'UserCreate', '', '', '', '新增用户', '', 0, 3, 'user:create', 1, 1, NOW(), NOW()),
(10102, 101, 'UserUpdate', '', '', '', '编辑用户', '', 0, 3, 'user:update', 2, 1, NOW(), NOW()),
(10103, 101, 'UserDelete', '', '', '', '删除用户', '', 0, 3, 'user:delete', 3, 1, NOW(), NOW()),
(10104, 101, 'UserExport', '', '', '', '导出用户', '', 0, 3, 'user:export', 4, 1, NOW(), NOW()),
(10105, 101, 'UserResetPwd', '', '', '', '重置密码', '', 0, 3, 'user:reset-password', 5, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（权限管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(20101, 201, 'RoleCreate', '', '', '', '新增角色', '', 0, 3, 'permission:role:create', 1, 1, NOW(), NOW()),
(20102, 201, 'RoleUpdate', '', '', '', '编辑角色', '', 0, 3, 'permission:role:update', 2, 1, NOW(), NOW()),
(20103, 201, 'RoleDelete', '', '', '', '删除角色', '', 0, 3, 'permission:role:delete', 3, 1, NOW(), NOW()),
(20104, 201, 'RoleAssign', '', '', '', '分配权限', '', 0, 3, 'permission:role:assign', 4, 1, NOW(), NOW()),
(20201, 202, 'MenuCreate', '', '', '', '新增菜单', '', 0, 3, 'permission:menu:create', 1, 1, NOW(), NOW()),
(20202, 202, 'MenuUpdate', '', '', '', '编辑菜单', '', 0, 3, 'permission:menu:update', 2, 1, NOW(), NOW()),
(20203, 202, 'MenuDelete', '', '', '', '删除菜单', '', 0, 3, 'permission:menu:delete', 3, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（题库管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(40101, 401, 'QuestionCreateBtn', '', '', '', '新增题目', '', 0, 3, 'question:create', 1, 1, NOW(), NOW()),
(40102, 401, 'QuestionUpdateBtn', '', '', '', '编辑题目', '', 0, 3, 'question:update', 2, 1, NOW(), NOW()),
(40103, 401, 'QuestionDeleteBtn', '', '', '', '删除题目', '', 0, 3, 'question:delete', 3, 1, NOW(), NOW()),
(40104, 401, 'QuestionImport', '', '', '', '导入题目', '', 0, 3, 'question:import', 4, 1, NOW(), NOW()),
(40105, 401, 'QuestionExport', '', '', '', '导出题目', '', 0, 3, 'question:export', 5, 1, NOW(), NOW()),
(40301, 403, 'CategoryCreate', '', '', '', '新增分类', '', 0, 3, 'question:category:create', 1, 1, NOW(), NOW()),
(40302, 403, 'CategoryUpdate', '', '', '', '编辑分类', '', 0, 3, 'question:category:update', 2, 1, NOW(), NOW()),
(40303, 403, 'CategoryDelete', '', '', '', '删除分类', '', 0, 3, 'question:category:delete', 3, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（考试管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(50101, 501, 'PaperCreate', '', '', '', '新增试卷', '', 0, 3, 'exam:paper:create', 1, 1, NOW(), NOW()),
(50102, 501, 'PaperUpdate', '', '', '', '编辑试卷', '', 0, 3, 'exam:paper:update', 2, 1, NOW(), NOW()),
(50103, 501, 'PaperDelete', '', '', '', '删除试卷', '', 0, 3, 'exam:paper:delete', 3, 1, NOW(), NOW()),
(50201, 502, 'ExamCreate', '', '', '', '发布考试', '', 0, 3, 'exam:create', 1, 1, NOW(), NOW()),
(50202, 502, 'ExamUpdate', '', '', '', '编辑考试', '', 0, 3, 'exam:update', 2, 1, NOW(), NOW()),
(50203, 502, 'ExamDelete', '', '', '', '删除考试', '', 0, 3, 'exam:delete', 3, 1, NOW(), NOW()),
(50301, 503, 'ScoreExport', '', '', '', '导出成绩', '', 0, 3, 'exam:score:export', 1, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（班级管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(60101, 601, 'ClassCreate', '', '', '', '新增班级', '', 0, 3, 'class:create', 1, 1, NOW(), NOW()),
(60102, 601, 'ClassUpdate', '', '', '', '编辑班级', '', 0, 3, 'class:update', 2, 1, NOW(), NOW()),
(60103, 601, 'ClassDelete', '', '', '', '删除班级', '', 0, 3, 'class:delete', 3, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（内容管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(70101, 701, 'CommentDelete', '', '', '', '删除评论', '', 0, 3, 'content:comment:delete', 1, 1, NOW(), NOW()),
(70102, 701, 'CommentAudit', '', '', '', '审核评论', '', 0, 3, 'content:comment:audit', 2, 1, NOW(), NOW()),
(70201, 702, 'ReviewApprove', '', '', '', '审核通过', '', 0, 3, 'content:review:approve', 1, 1, NOW(), NOW()),
(70202, 702, 'ReviewReject', '', '', '', '审核拒绝', '', 0, 3, 'content:review:reject', 2, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（系统设置）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(110101, 1101, 'SettingUpdate', '', '', '', '修改设置', '', 0, 3, 'system:setting:update', 1, 1, NOW(), NOW()),
(110501, 1105, 'BackupCreate', '', '', '', '创建备份', '', 0, 3, 'system:backup:create', 1, 1, NOW(), NOW()),
(110502, 1105, 'BackupRestore', '', '', '', '恢复备份', '', 0, 3, 'system:backup:restore', 2, 1, NOW(), NOW()),
(110601, 1106, 'ApprovalApprove', '', '', '', '审批通过', '', 0, 3, 'system:approval:approve', 1, 1, NOW(), NOW()),
(110602, 1106, 'ApprovalReject', '', '', '', '审批拒绝', '', 0, 3, 'system:approval:reject', 2, 1, NOW(), NOW());

-- =====================================================
-- 系统角色种子数据（确保角色存在）
-- =====================================================

INSERT IGNORE INTO roles (id, name, code, description, is_default, is_system, sort, status, created_at, updated_at) VALUES
(1, '超级管理员', 'admin', '系统超级管理员，拥有全部权限', 0, 1, 1, 1, NOW(), NOW()),
(2, '教师', 'teacher', '教师角色，可管理题库、考试、班级', 0, 1, 2, 1, NOW(), NOW()),
(3, '创作者', 'creator', '内容创作者，可管理题目和内容', 0, 1, 3, 1, NOW(), NOW()),
(4, '学生', 'student', '学生角色，可参加考试和查看班级', 1, 1, 4, 1, NOW(), NOW()),
(5, '普通用户', 'user', '普通用户，仅可访问首页', 0, 1, 5, 1, NOW(), NOW());

-- =====================================================
-- 超级管理员角色（ID=1）拥有所有菜单权限
-- =====================================================

INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus;

-- =====================================================
-- 教师角色（ID=2）菜单权限
-- 范围：首页、用户列表(只读)、题库、考试、班级、统计、文件
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页
(2, 1),
-- 用户管理（仅查看用户列表）
(2, 2), (2, 101),
-- 题库管理（全部）
(2, 4), (2, 401), (2, 402), (2, 403), (2, 404), (2, 405),
(2, 40101), (2, 40102), (2, 40103), (2, 40104), (2, 40105),
(2, 40301), (2, 40302), (2, 40303),
-- 考试管理（全部）
(2, 5), (2, 501), (2, 502), (2, 503),
(2, 50101), (2, 50102), (2, 50103), (2, 50201), (2, 50202), (2, 50203), (2, 50301),
-- 班级管理（全部）
(2, 6), (2, 601), (2, 602), (2, 60101), (2, 60102), (2, 60103),
-- 数据统计（基础）
(2, 8), (2, 801), (2, 802), (2, 803), (2, 804), (2, 805), (2, 806),
-- 文件管理（查看）
(2, 10), (2, 1001);

-- =====================================================
-- 创作者角色（ID=3）菜单权限
-- 范围：首页、题库(含分类)、内容管理、文件
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页
(3, 1),
-- 题库管理（全部）
(3, 4), (3, 401), (3, 402), (3, 403), (3, 404),
(3, 40101), (3, 40102), (3, 40103), (3, 40104), (3, 40105),
(3, 40301), (3, 40302), (3, 40303),
-- 内容管理（全部）
(3, 7), (3, 701), (3, 702), (3, 703), (3, 70101), (3, 70102), (3, 70201), (3, 70202),
-- 文件管理（查看）
(3, 10), (3, 1001);

-- =====================================================
-- 学生角色（ID=4）菜单权限
-- 范围：首页、考试(参加)、班级(查看)
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页
(4, 1),
-- 考试管理（仅查看考试列表和成绩）
(4, 5), (4, 502), (4, 503),
-- 班级管理（仅查看）
(4, 6), (4, 601);

-- =====================================================
-- 普通用户角色（ID=5）菜单权限
-- 范围：首页
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
(5, 1);

-- =====================================================
-- 完成
-- =====================================================

SELECT '菜单种子数据初始化完成！' AS message;
