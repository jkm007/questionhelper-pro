-- =====================================================
-- 题小助管理后台菜单种子数据
-- 与需求文档保持一致：
--   - 角色定义：11.3 预置角色（9个）
--   - 菜单结构：12.2 管理后台菜单树
--   - 权限命名：module:submodule:operation 三级格式
--   - 权限矩阵：11.5 功能操作权限矩阵
-- =====================================================

-- =====================================================
-- 系统角色种子数据（需求文档 11.3 预置角色）
-- =====================================================

INSERT IGNORE INTO roles (id, name, code, description, is_default, is_system, sort, status, created_at, updated_at) VALUES
(1,  '超级管理员',   'super_admin',     '系统最高权限，管理所有功能，审批班级创作者申请', 0, 1, 1, 1, NOW(), NOW()),
(2,  '管理员',       'admin',           '系统管理，管理用户、内容、系统配置', 0, 1, 2, 1, NOW(), NOW()),
(3,  '内容管理员',   'content_manager', '管理公共题库、公共分类、公共知识点', 0, 1, 3, 1, NOW(), NOW()),
(4,  '用户管理员',   'user_manager',    '管理用户账号、审核', 0, 1, 4, 1, NOW(), NOW()),
(5,  '数据分析师',   'data_analyst',    '查看统计数据、生成报表', 0, 1, 5, 1, NOW(), NOW()),
(6,  '运营人员',     'operator',        '管理公告、推荐、活动', 0, 1, 6, 1, NOW(), NOW()),
(7,  '教师',         'teacher',         '创建班级、邀请成员、管理班级资源、发布考试', 0, 1, 7, 1, NOW(), NOW()),
(8,  '创作者',       'creator',         '创建私有题目/分类/知识点/考试，使用公共资源', 0, 1, 8, 1, NOW(), NOW()),
(9,  '学生',         'student',         '使用公共题库、加入班级、练习考试', 1, 1, 9, 1, NOW(), NOW());

-- =====================================================
-- 一级菜单（目录）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(1,  NULL, 'Dashboard',   '/dashboard',   'Layout', '/dashboard', '首页',     'HomeFilled',  0, 1, '', 1,  1, NOW(), NOW()),
(2,  NULL, 'User',        '/user',        'Layout', '',           '用户管理', 'User',        0, 1, '', 2,  1, NOW(), NOW()),
(3,  NULL, 'Question',    '/question',    'Layout', '',           '题库管理', 'EditPen',     0, 1, '', 3,  1, NOW(), NOW()),
(4,  NULL, 'Exam',        '/exam',        'Layout', '',           '考试管理', 'Timer',       0, 1, '', 4,  1, NOW(), NOW()),
(5,  NULL, 'Class',       '/class',       'Layout', '',           '班级管理', 'UserFilled',  0, 1, '', 5,  1, NOW(), NOW()),
(6,  NULL, 'Statistics',  '/statistics',  'Layout', '',           '数据统计', 'DataAnalysis',0, 1, '', 6,  1, NOW(), NOW()),
(7,  NULL, 'System',      '/system',      'Layout', '',           '系统设置', 'Setting',     0, 1, '', 7,  1, NOW(), NOW()),
(8,  NULL, 'Profile',     '/profile',     'Layout', '',           '个人中心', 'UserFilled',  1, 1, '', 8,  1, NOW(), NOW());

-- =====================================================
-- 二级菜单（用户管理）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(201, 2, 'UserList',       '/user/list',       'system/user/index',       '', '用户列表',       'User',      0, 2, 'user:list',       1, 1, NOW(), NOW()),
(202, 2, 'ClassCreator',   '/user/creator',     'system/creator/index',    '', '班级创作者管理', 'Star',      0, 2, 'class:creator:view', 2, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（题库管理）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(301, 3, 'QuestionPublic',    '/question/public',    'question/public/index',    '', '公共题目',   'Document',  0, 2, 'question:public:view',   1, 1, NOW(), NOW()),
(302, 3, 'QuestionPrivate',   '/question/private',   'question/private/index',   '', '私有题目',   'Lock',      0, 2, 'question:private:view',  2, 1, NOW(), NOW()),
(303, 3, 'QuestionClass',     '/question/class',     'question/class/index',     '', '班级题目',   'School',    0, 2, 'question:class:view',    3, 1, NOW(), NOW()),
(304, 3, 'QuestionCategory',  '/question/category',  'question/category/index',  '', '题目分类',   'Menu',      0, 2, 'category:view',          4, 1, NOW(), NOW()),
(305, 3, 'QuestionKnowledge', '/question/knowledge', 'question/knowledge/index', '', '知识点管理', 'Collection',0, 2, 'knowledge:view',         5, 1, NOW(), NOW()),
(306, 3, 'QuestionImport',    '/question/import',    'question/import/index',    '', '批量导入',   'Upload',    0, 2, 'question:import',        6, 1, NOW(), NOW()),
(307, 3, 'QuestionExport',    '/question/export',    'question/export/index',    '', '批量导出',   'Download',  0, 2, 'question:export',        7, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（考试管理）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(401, 4, 'ExamPaper',    '/exam/paper',    'exam/paper/index',    '', '试卷列表', 'Document',  0, 2, 'exam:view',      1, 1, NOW(), NOW()),
(402, 4, 'ExamComposer', '/exam/composer', 'exam/composer/index', '', '组卷管理', 'EditPen',   0, 2, 'composer:view',  2, 1, NOW(), NOW()),
(403, 4, 'ExamPublish',  '/exam/publish',  'exam/publish/index',  '', '考试发布', 'Promotion', 0, 2, 'exam:publish',   3, 1, NOW(), NOW()),
(404, 4, 'ExamMonitor',  '/exam/monitor',  'exam/monitor/index',  '', '考试监控', 'Monitor',   0, 2, 'exam:monitor',   4, 1, NOW(), NOW()),
(405, 4, 'ExamScore',    '/exam/score',    'exam/score/index',    '', '成绩管理', 'Trophy',    0, 2, 'score:view',     5, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（班级管理）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(501, 5, 'ClassList',      '/class/list',      'class/list/index',      '', '班级列表',       'School',    0, 2, 'class:view',            1, 1, NOW(), NOW()),
(502, 5, 'ClassDetail',    '/class/detail',    'class/detail/index',    '', '班级详情',       'View',      0, 2, 'class:detail',          2, 1, NOW(), NOW()),
(503, 5, 'ClassApproval',  '/class/approval',  'class/approval/index',  '', '班级创作者审批', 'Stamp',     0, 2, 'class:creator:approve', 3, 1, NOW(), NOW()),
(504, 5, 'ClassResource',  '/class/resource',  'class/resource/index',  '', '班级资源管理',   'Resources', 0, 2, 'class:resource:view',   4, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（数据统计）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(601, 6, 'StatisticsOverview', '/statistics/overview', 'statistics/overview/index', '', '数据概览', 'Odometer',    0, 2, 'statistics:overview', 1, 1, NOW(), NOW()),
(602, 6, 'StatisticsUser',     '/statistics/user',     'statistics/user/index',     '', '用户统计', 'Histogram',   0, 2, 'statistics:user',     2, 1, NOW(), NOW()),
(603, 6, 'StatisticsAnswer',   '/statistics/answer',   'statistics/answer/index',   '', '答题统计', 'PieChart',    0, 2, 'statistics:answer',   3, 1, NOW(), NOW()),
(604, 6, 'StatisticsScore',    '/statistics/score',    'statistics/score/index',    '', '成绩分析', 'TrendCharts', 0, 2, 'statistics:score',    4, 1, NOW(), NOW()),
(605, 6, 'StatisticsSystem',   '/statistics/system',   'statistics/system/index',   '', '系统监控', 'Monitor',     0, 2, 'statistics:system',   5, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（系统设置）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(701, 7, 'SystemSetting',    '/system/setting',    'setting/basic/index',          '', '基础设置', 'Setting',    0, 2, 'system:settings',    1, 1, NOW(), NOW()),
(702, 7, 'SystemSecurity',   '/system/security',   'system/security/index',        '', '安全设置', 'Lock',       0, 2, 'system:security',    2, 1, NOW(), NOW()),
(703, 7, 'SystemRole',       '/system/role',       'system/role/index',            '', '角色管理', 'UserFilled', 0, 2, 'role:view',          3, 1, NOW(), NOW()),
(704, 7, 'SystemPermission', '/system/permission', 'system/permission/index',      '', '权限管理', 'Key',        0, 2, 'permission:view',    4, 1, NOW(), NOW()),
(705, 7, 'SystemAudit',      '/system/audit',      'system/audit/index',           '', '操作审计', 'Document',   0, 2, 'system:audit',       5, 1, NOW(), NOW()),
(706, 7, 'SystemLog',        '/system/log',        'log/system/index',             '', '日志管理', 'Notebook',   0, 2, 'system:logs',        6, 1, NOW(), NOW()),
(707, 7, 'SystemNotify',     '/system/notify',     'system/notify/index',          '', '通知设置', 'Bell',       0, 2, 'system:notification',7, 1, NOW(), NOW()),
(708, 7, 'SystemBackup',     '/system/backup',     'backup/list/index',            '', '数据备份', 'Upload',     0, 2, 'system:backup',      8, 1, NOW(), NOW());

-- =====================================================
-- 二级菜单（个人中心）— 需求文档 12.2
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(801, 8, 'ProfileInfo',     '/profile/info',     'profile/info/index',     '', '个人信息', 'User',    0, 2, 'profile:view',    1, 1, NOW(), NOW()),
(802, 8, 'ProfilePassword', '/profile/password', 'profile/password/index', '', '修改密码', 'Lock',    0, 2, 'profile:password',2, 1, NOW(), NOW()),
(803, 8, 'ProfileMessage',  '/profile/message',  'profile/message/index',  '', '我的消息', 'Message', 0, 2, 'profile:message', 3, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（用户管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(20101, 201, 'UserCreate',    '', '', '', '新增用户', '', 0, 3, 'user:list:create',    1, 1, NOW(), NOW()),
(20102, 201, 'UserUpdate',    '', '', '', '编辑用户', '', 0, 3, 'user:list:edit',      2, 1, NOW(), NOW()),
(20103, 201, 'UserDelete',    '', '', '', '删除用户', '', 0, 3, 'user:list:delete',    3, 1, NOW(), NOW()),
(20104, 201, 'UserExport',    '', '', '', '导出用户', '', 0, 3, 'user:list:export',    4, 1, NOW(), NOW()),
(20105, 201, 'UserResetPwd',  '', '', '', '重置密码', '', 0, 3, 'user:list:reset-pwd', 5, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（题库管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
-- 公共题目按钮
(30101, 301, 'PublicQCreate',  '', '', '', '新增公共题目', '', 0, 3, 'question:public:create',  1, 1, NOW(), NOW()),
(30102, 301, 'PublicQUpdate',  '', '', '', '编辑公共题目', '', 0, 3, 'question:public:edit',    2, 1, NOW(), NOW()),
(30103, 301, 'PublicQDelete',  '', '', '', '删除公共题目', '', 0, 3, 'question:public:delete',  3, 1, NOW(), NOW()),
-- 私有题目按钮
(30201, 302, 'PrivateQCreate', '', '', '', '新增私有题目', '', 0, 3, 'question:private:create', 1, 1, NOW(), NOW()),
(30202, 302, 'PrivateQUpdate', '', '', '', '编辑私有题目', '', 0, 3, 'question:private:edit',   2, 1, NOW(), NOW()),
(30203, 302, 'PrivateQDelete', '', '', '', '删除私有题目', '', 0, 3, 'question:private:delete', 3, 1, NOW(), NOW()),
-- 班级题目按钮
(30301, 303, 'ClassQCreate',   '', '', '', '新增班级题目', '', 0, 3, 'question:class:create',   1, 1, NOW(), NOW()),
(30302, 303, 'ClassQUpdate',   '', '', '', '编辑班级题目', '', 0, 3, 'question:class:edit',     2, 1, NOW(), NOW()),
(30303, 303, 'ClassQDelete',   '', '', '', '删除班级题目', '', 0, 3, 'question:class:delete',   3, 1, NOW(), NOW()),
-- 分类按钮
(30401, 304, 'CategoryCreate', '', '', '', '新增分类', '', 0, 3, 'category:create', 1, 1, NOW(), NOW()),
(30402, 304, 'CategoryUpdate', '', '', '', '编辑分类', '', 0, 3, 'category:edit',   2, 1, NOW(), NOW()),
(30403, 304, 'CategoryDelete', '', '', '', '删除分类', '', 0, 3, 'category:delete', 3, 1, NOW(), NOW()),
-- 知识点按钮
(30501, 305, 'KnowledgeCreate','', '', '', '新增知识点', '', 0, 3, 'knowledge:create', 1, 1, NOW(), NOW()),
(30502, 305, 'KnowledgeUpdate','', '', '', '编辑知识点', '', 0, 3, 'knowledge:edit',   2, 1, NOW(), NOW()),
(30503, 305, 'KnowledgeDelete','', '', '', '删除知识点', '', 0, 3, 'knowledge:delete', 3, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（考试管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
-- 试卷按钮
(40101, 401, 'PaperCreate',   '', '', '', '新增试卷', '', 0, 3, 'exam:paper:create', 1, 1, NOW(), NOW()),
(40102, 401, 'PaperUpdate',   '', '', '', '编辑试卷', '', 0, 3, 'exam:paper:edit',   2, 1, NOW(), NOW()),
(40103, 401, 'PaperDelete',   '', '', '', '删除试卷', '', 0, 3, 'exam:paper:delete', 3, 1, NOW(), NOW()),
-- 考试按钮
(40301, 403, 'ExamPublish',   '', '', '', '发布考试', '', 0, 3, 'exam:publish:create', 1, 1, NOW(), NOW()),
(40302, 403, 'ExamUpdate',    '', '', '', '编辑考试', '', 0, 3, 'exam:publish:edit',   2, 1, NOW(), NOW()),
(40303, 403, 'ExamDelete',    '', '', '', '删除考试', '', 0, 3, 'exam:publish:delete', 3, 1, NOW(), NOW()),
-- 成绩按钮
(40501, 405, 'ScoreExport',   '', '', '', '导出成绩', '', 0, 3, 'score:export', 1, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（班级管理）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(50101, 501, 'ClassCreate',      '', '', '', '新增班级',       '', 0, 3, 'class:create',          1, 1, NOW(), NOW()),
(50102, 501, 'ClassUpdate',      '', '', '', '编辑班级',       '', 0, 3, 'class:edit',            2, 1, NOW(), NOW()),
(50103, 501, 'ClassDelete',      '', '', '', '删除班级',       '', 0, 3, 'class:delete',          3, 1, NOW(), NOW()),
(50104, 501, 'ClassMemberManage','', '', '', '管理班级成员',   '', 0, 3, 'class:member:manage',   4, 1, NOW(), NOW()),
(50301, 503, 'CreatorApprove',   '', '', '', '审批班级创作者', '', 0, 3, 'class:creator:approve', 1, 1, NOW(), NOW()),
(50302, 503, 'CreatorReject',    '', '', '', '拒绝班级创作者', '', 0, 3, 'class:creator:reject',  2, 1, NOW(), NOW());

-- =====================================================
-- 按钮权限（系统设置）
-- =====================================================

INSERT INTO menus (id, parent_id, name, path, component, redirect, title, icon, hidden, type, permission, sort, status, created_at, updated_at) VALUES
(70101, 701, 'SettingEdit',     '', '', '', '修改设置', '', 0, 3, 'system:settings:edit',     1, 1, NOW(), NOW()),
(70301, 703, 'RoleCreate',      '', '', '', '新增角色', '', 0, 3, 'role:create',              1, 1, NOW(), NOW()),
(70302, 703, 'RoleUpdate',      '', '', '', '编辑角色', '', 0, 3, 'role:edit',                2, 1, NOW(), NOW()),
(70303, 703, 'RoleDelete',      '', '', '', '删除角色', '', 0, 3, 'role:delete',              3, 1, NOW(), NOW()),
(70304, 703, 'RoleAssign',      '', '', '', '分配权限', '', 0, 3, 'role:assign',              4, 1, NOW(), NOW()),
(70801, 708, 'BackupCreate',    '', '', '', '创建备份', '', 0, 3, 'system:backup:create',     1, 1, NOW(), NOW()),
(70802, 708, 'BackupRestore',   '', '', '', '恢复备份', '', 0, 3, 'system:backup:restore',    2, 1, NOW(), NOW());

-- =====================================================
-- 超级管理员（ID=1）拥有所有菜单权限
-- =====================================================

INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus;

-- =====================================================
-- 管理员（ID=2）菜单权限
-- 需求文档 11.5.2：用户管理、公共题目管理、私有题目管理、班级题目管理、
--   题目审核、公共试卷管理、私有试卷管理、班级试卷管理、考试发布、
--   班级创建、班级成员管理、班级创作者审批、成绩管理、数据统计
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(2, 1), (2, 8), (2, 801), (2, 802), (2, 803),
-- 用户管理（全部）
(2, 2), (2, 201), (2, 202),
(2, 20101), (2, 20102), (2, 20103), (2, 20104), (2, 20105),
-- 题库管理（全部：公共+私有+班级+分类+知识点+导入导出）
(2, 3), (2, 301), (2, 302), (2, 303), (2, 304), (2, 305), (2, 306), (2, 307),
(2, 30101), (2, 30102), (2, 30103),
(2, 30201), (2, 30202), (2, 30203),
(2, 30301), (2, 30302), (2, 30303),
(2, 30401), (2, 30402), (2, 30403),
(2, 30501), (2, 30502), (2, 30503),
-- 考试管理（全部）
(2, 4), (2, 401), (2, 402), (2, 403), (2, 404), (2, 405),
(2, 40101), (2, 40102), (2, 40103),
(2, 40301), (2, 40302), (2, 40303),
(2, 40501),
-- 班级管理（全部）
(2, 5), (2, 501), (2, 502), (2, 503), (2, 504),
(2, 50101), (2, 50102), (2, 50103), (2, 50104),
(2, 50301), (2, 50302),
-- 数据统计（全部）
(2, 6), (2, 601), (2, 602), (2, 603), (2, 604), (2, 605),
-- 系统设置（全部）
(2, 7), (2, 701), (2, 702), (2, 703), (2, 704), (2, 705), (2, 706), (2, 707), (2, 708),
(2, 70101), (2, 70301), (2, 70302), (2, 70303), (2, 70304),
(2, 70801), (2, 70802);

-- =====================================================
-- 内容管理员（ID=3）菜单权限
-- 需求文档 11.5.2：公共题目管理、题目审核、公共试卷管理
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(3, 1), (3, 8), (3, 801), (3, 802), (3, 803),
-- 题库管理（公共题目+分类+知识点）
(3, 3), (3, 301), (3, 304), (3, 305),
(3, 30101), (3, 30102), (3, 30103),
(3, 30401), (3, 30402), (3, 30403),
(3, 30501), (3, 30502), (3, 30503),
-- 考试管理（公共试卷）
(3, 4), (3, 401), (3, 40101), (3, 40102), (3, 40103);

-- =====================================================
-- 用户管理员（ID=4）菜单权限
-- 需求文档 11.5.2：用户管理
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(4, 1), (4, 8), (4, 801), (4, 802), (4, 803),
-- 用户管理（全部）
(4, 2), (4, 201), (4, 202),
(4, 20101), (4, 20102), (4, 20103), (4, 20104), (4, 20105);

-- =====================================================
-- 数据分析师（ID=5）菜单权限
-- 需求文档 11.5.2：数据统计
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(5, 1), (5, 8), (5, 801), (5, 802), (5, 803),
-- 数据统计（全部，只读）
(5, 6), (5, 601), (5, 602), (5, 603), (5, 604), (5, 605);

-- =====================================================
-- 运营人员（ID=6）菜单权限
-- 需求文档 11.5.2：公告、推荐、活动（暂用内容管理+通知）
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(6, 1), (6, 8), (6, 801), (6, 802), (6, 803),
-- 数据统计（查看）
(6, 6), (6, 601);

-- =====================================================
-- 教师（ID=7）菜单权限
-- 需求文档 11.5.2：私有题目管理、班级题目管理、私有试卷管理、
--   班级试卷管理、考试发布(班级)、班级创建、班级成员管理、
--   班级创作者审批、成绩管理(自己班级)、数据统计(自己班级)
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(7, 1), (7, 8), (7, 801), (7, 802), (7, 803),
-- 用户管理（查看班级创作者）
(7, 2), (7, 202),
-- 题库管理（私有题目+班级题目+分类+知识点）
(7, 3), (7, 302), (7, 303), (7, 304), (7, 305),
(7, 30201), (7, 30202), (7, 30203),
(7, 30301), (7, 30302), (7, 30303),
(7, 30401), (7, 30402), (7, 30403),
(7, 30501), (7, 30502), (7, 30503),
-- 考试管理（私有试卷+班级试卷+发布+成绩）
(7, 4), (7, 401), (7, 402), (7, 403), (7, 405),
(7, 40101), (7, 40102), (7, 40103),
(7, 40301), (7, 40302), (7, 40303),
(7, 40501),
-- 班级管理（全部）
(7, 5), (7, 501), (7, 502), (7, 503), (7, 504),
(7, 50101), (7, 50102), (7, 50103), (7, 50104),
(7, 50301), (7, 50302),
-- 数据统计（自己班级）
(7, 6), (7, 601), (7, 602), (7, 603), (7, 604);

-- =====================================================
-- 创作者（ID=8）菜单权限
-- 需求文档 11.5.2：私有题目管理、私有试卷管理
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(8, 1), (8, 8), (8, 801), (8, 802), (8, 803),
-- 题库管理（私有题目+分类+知识点）
(8, 3), (8, 302), (8, 304), (8, 305),
(8, 30201), (8, 30202), (8, 30203),
(8, 30401), (8, 30402), (8, 30403),
(8, 30501), (8, 30502), (8, 30503),
-- 考试管理（私有试卷）
(8, 4), (8, 401), (8, 40101), (8, 40102), (8, 40103);

-- =====================================================
-- 学生（ID=9）菜单权限
-- 需求文档 11.5.2：答题、错题本、收藏、评论
-- 学生通过移动端使用，不访问管理后台，仅保留最小菜单
-- =====================================================

INSERT INTO role_menus (role_id, menu_id) VALUES
-- 首页 + 个人中心
(9, 1), (9, 8), (9, 801), (9, 802), (9, 803);

-- =====================================================
-- 完成
-- =====================================================

SELECT '菜单种子数据初始化完成！（9个角色，对齐需求文档）' AS message;
