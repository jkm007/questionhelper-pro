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
(1, NULL, 'Dashboard', '/dashboard', 'Layout', '/dashboard', '首页', 'homepage', 0, 1, '', 1, 1, NOW(), NOW());

-- 2. 用户管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(2, NULL, 'User', '/user', 'Layout', '', '用户管理', 'user', 0, 1, '', 2, 1, NOW(), NOW());

-- 3. 权限管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(3, NULL, 'Permission', '/permission', 'Layout', '', '权限管理', 'lock', 0, 1, '', 3, 1, NOW(), NOW());

-- 4. 题库管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(4, NULL, 'Question', '/question', 'Layout', '', '题库管理', 'edit', 0, 1, '', 4, 1, NOW(), NOW());

-- 5. 考试管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(5, NULL, 'Exam', '/exam', 'Layout', '', '考试管理', 'exam', 0, 1, '', 5, 1, NOW(), NOW());

-- 6. 班级管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(6, NULL, 'Class', '/class', 'Layout', '', '班级管理', 'peoples', 0, 1, '', 6, 1, NOW(), NOW());

-- 7. 内容管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(7, NULL, 'Content', '/content', 'Layout', '', '内容管理', 'documentation', 0, 1, '', 7, 1, NOW(), NOW());

-- 8. 数据统计
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(8, NULL, 'Statistics', '/statistics', 'Layout', '', '数据统计', 'chart', 0, 1, '', 8, 1, NOW(), NOW());

-- 9. 消息中心
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(9, NULL, 'Message', '/message', 'Layout', '', '消息中心', 'message', 0, 1, '', 9, 1, NOW(), NOW());

-- 10. 文件管理
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(10, NULL, 'File', '/file', 'Layout', '', '文件管理', 'file', 0, 1, '', 10, 1, NOW(), NOW());

-- 11. 系统设置
INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(11, NULL, 'System', '/system', 'Layout', '', '系统设置', 'setting', 0, 1, '', 11, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（用户管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(101, 2, 'UserList', '/user/list', 'user/list/index', '', '用户列表', 'user', 0, 2, 'user:list', 1, 1, NOW(), NOW()),
(102, 2, 'UserTag', '/user/tag', 'user/tag/index', '', '用户标签', 'tag', 0, 2, 'user:tag', 2, 1, NOW(), NOW()),
(103, 2, 'UserLog', '/user/log', 'user/log/index', '', '操作日志', 'log', 0, 2, 'user:log', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（权限管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(201, 3, 'Role', '/permission/role', 'permission/role/index', '', '角色管理', 'role', 0, 2, 'permission:role', 1, 1, NOW(), NOW()),
(202, 3, 'Menu', '/permission/menu', 'permission/menu/index', '', '菜单管理', 'menu', 0, 2, 'permission:menu', 2, 1, NOW(), NOW()),
(203, 3, 'Button', '/permission/button', 'permission/button/index', '', '按钮权限', 'button', 0, 2, 'permission:button', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（题库管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(401, 4, 'QuestionList', '/question/list', 'question/list/index', '', '题目列表', 'list', 0, 2, 'question:list', 1, 1, NOW(), NOW()),
(402, 4, 'QuestionCreate', '/question/create', 'question/create/index', '', '题目创建', 'add', 0, 2, 'question:create', 2, 1, NOW(), NOW()),
(403, 4, 'QuestionCategory', '/question/category', 'question/category/index', '', '分类管理', 'category', 0, 2, 'question:category', 3, 1, NOW(), NOW()),
(404, 4, 'QuestionKnowledge', '/question/knowledge', 'question/knowledge/index', '', '知识点管理', 'knowledge', 0, 2, 'question:knowledge', 4, 1, NOW(), NOW()),
(405, 4, 'QuestionReview', '/question/review', 'question/review/index', '', '题目审核', 'review', 0, 2, 'question:review', 5, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（考试管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(501, 5, 'ExamPaper', '/exam/paper', 'exam/paper/index', '', '试卷管理', 'paper', 0, 2, 'exam:paper', 1, 1, NOW(), NOW()),
(502, 5, 'ExamList', '/exam/list', 'exam/list/index', '', '考试列表', 'exam-list', 0, 2, 'exam:list', 2, 1, NOW(), NOW()),
(503, 5, 'ExamScore', '/exam/score', 'exam/score/index', '', '成绩管理', 'score', 0, 2, 'exam:score', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（班级管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(601, 6, 'ClassList', '/class/list', 'class/list/index', '', '班级列表', 'class-list', 0, 2, 'class:list', 1, 1, NOW(), NOW()),
(602, 6, 'ClassAudit', '/class/audit', 'class/audit/index', '', '班级审核', 'audit', 0, 2, 'class:audit', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（内容管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(701, 7, 'ContentComment', '/content/comment', 'content/comment/index', '', '评论管理', 'comment', 0, 2, 'content:comment', 1, 1, NOW(), NOW()),
(702, 7, 'ContentReview', '/content/review', 'content/review/index', '', '资源审核', 'review', 0, 2, 'content:review', 2, 1, NOW(), NOW()),
(703, 7, 'ContentCreator', '/content/creator', 'content/creator/index', '', '创作者管理', 'creator', 0, 2, 'content:creator', 3, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（数据统计）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(801, 8, 'StatisticsOverview', '/statistics/overview', 'statistics/overview/index', '', '数据概览', 'overview', 0, 2, 'statistics:overview', 1, 1, NOW(), NOW()),
(802, 8, 'StatisticsUser', '/statistics/user', 'statistics/user/index', '', '用户统计', 'user-stats', 0, 2, 'statistics:user', 2, 1, NOW(), NOW()),
(803, 8, 'StatisticsAnswer', '/statistics/answer', 'statistics/answer/index', '', '答题统计', 'answer-stats', 0, 2, 'statistics:answer', 3, 1, NOW(), NOW()),
(804, 8, 'StatisticsScore', '/statistics/score', 'statistics/score/index', '', '成绩统计', 'score-stats', 0, 2, 'statistics:score', 4, 1, NOW(), NOW()),
(805, 8, 'StatisticsClass', '/statistics/class', 'statistics/class/index', '', '班级统计', 'class-stats', 0, 2, 'statistics:class', 5, 1, NOW(), NOW()),
(806, 8, 'StatisticsAlert', '/statistics/alert', 'statistics/alert/index', '', '数据预警', 'alert', 0, 2, 'statistics:alert', 6, 1, NOW(), NOW()),
(807, 8, 'StatisticsAdvanced', '/statistics/advanced', 'Layout', '', '高级分析', 'advanced', 0, 1, 'statistics:advanced', 7, 1, NOW(), NOW());

-- =====================================================
-- 三级菜单（高级分析）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(80701, 807, 'AdvancedRetention', '/statistics/advanced/retention', 'statistics/advanced/retention/index', '', '用户留存', 'retention', 0, 2, 'statistics:advanced:retention', 1, 1, NOW(), NOW()),
(80702, 807, 'AdvancedChurn', '/statistics/advanced/churn', 'statistics/advanced/churn/index', '', '用户流失', 'churn', 0, 2, 'statistics:advanced:churn', 2, 1, NOW(), NOW()),
(80703, 807, 'AdvancedBehavior', '/statistics/advanced/behavior', 'statistics/advanced/behavior/index', '', '用户行为', 'behavior', 0, 2, 'statistics:advanced:behavior', 3, 1, NOW(), NOW()),
(80704, 807, 'AdvancedSegment', '/statistics/advanced/segment', 'statistics/advanced/segment/index', '', '用户分群', 'segment', 0, 2, 'statistics:advanced:segment', 4, 1, NOW(), NOW()),
(80705, 807, 'AdvancedPath', '/statistics/advanced/path', 'statistics/advanced/path/index', '', '用户路径', 'path', 0, 2, 'statistics:advanced:path', 5, 1, NOW(), NOW()),
(80706, 807, 'AdvancedConversion', '/statistics/advanced/conversion', 'statistics/advanced/conversion/index', '', '用户转化', 'conversion', 0, 2, 'statistics:advanced:conversion', 6, 1, NOW(), NOW()),
(80707, 807, 'AdvancedQuestion', '/statistics/advanced/question', 'statistics/advanced/question/index', '', '题目分析', 'question-stats', 0, 2, 'statistics:advanced:question', 7, 1, NOW(), NOW()),
(80708, 807, 'AdvancedSubscribe', '/statistics/advanced/subscribe', 'statistics/advanced/subscribe/index', '', '数据订阅', 'subscribe', 0, 2, 'statistics:advanced:subscribe', 8, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（消息中心）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(901, 9, 'MessageNotification', '/message/notification', 'message/notification/index', '', '通知管理', 'notification', 0, 2, 'message:notification', 1, 1, NOW(), NOW()),
(902, 9, 'MessageTemplate', '/message/template', 'message/template/index', '', '通知模板', 'template', 0, 2, 'message:template', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（文件管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(1001, 10, 'FileList', '/file/list', 'file/list/index', '', '文件列表', 'file-list', 0, 2, 'file:list', 1, 1, NOW(), NOW()),
(1002, 10, 'FileStorage', '/file/storage', 'file/storage/index', '', '存储配置', 'storage', 0, 2, 'file:storage', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（系统设置）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(1101, 11, 'SystemSetting', '/system/setting', 'system/setting/index', '', '基础设置', 'setting', 0, 2, 'system:setting', 1, 1, NOW(), NOW()),
(1102, 11, 'SystemSettingClass', '/system/setting/class', 'system/setting/class/index', '', '班级设置', 'class-setting', 0, 2, 'system:setting:class', 2, 1, NOW(), NOW()),
(1103, 11, 'SystemSettingResource', '/system/setting/resource', 'system/setting/resource/index', '', '资源设置', 'resource-setting', 0, 2, 'system:setting:resource', 3, 1, NOW(), NOW()),
(1104, 11, 'SystemLog', '/system/log', 'system/log/index', '', '系统日志', 'log', 0, 2, 'system:log', 4, 1, NOW(), NOW()),
(1105, 11, 'SystemBackup', '/system/backup', 'system/backup/index', '', '数据备份', 'backup', 0, 2, 'system:backup', 5, 1, NOW(), NOW()),
(1106, 11, 'SystemApproval', '/system/approval', 'system/approval/index', '', '审批管理', 'approval', 0, 2, 'system:approval', 6, 1, NOW(), NOW());

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
-- 超级管理员角色（ID=1）拥有所有菜单权限
-- =====================================================

INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus;

-- =====================================================
-- 完成
-- =====================================================

SELECT '菜单种子数据初始化完成！' AS message;
