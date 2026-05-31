# 题小助 V2.0 API 接口文档

**项目名称**：题小助 (QuestionHelper)
**版本**：V2.0
**基础路径**：`/api/v1`
**认证方式**：JWT Bearer Token (Authorization Header)
**日期**：2026-05-31

---

## 目录

- [一、概述](#一概述)
- [二、认证与授权](#二认证与授权)
- [三、公共接口 (无需认证)](#三公共接口-无需认证)
- [四、用户端接口 (需登录)](#四用户端接口-需登录)
- [五、管理端接口 (需管理员权限)](#五管理端接口-需管理员权限)
- [六、通用约定](#六通用约定)

---

## 一、概述

### 接口统计

| 权限级别 | 数量 |
|----------|------|
| 公共接口 (无认证) | 12 |
| 用户端接口 (JWT) | ~313 |
| 管理端接口 (JWT + Admin) | ~424 |
| **合计** | **~749** |

### 请求格式

- Content-Type: `application/json` (除文件上传外)
- 文件上传: `multipart/form-data`
- 所有接口均返回 JSON 格式

### 通用响应结构

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 分页请求参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 10 | 每页数量 |

### 分页响应结构

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 二、认证与授权

### 认证方式

1. **JWT Token**：登录后获取 access_token 和 refresh_token
2. **请求头**：`Authorization: Bearer <access_token>`
3. **Token 刷新**：access_token 过期后使用 refresh_token 获取新 token

### 权限级别

| 级别 | 说明 |
|------|------|
| Public | 无需认证 |
| Authorized | 需要有效 JWT Token |
| Admin | 需要 JWT Token + 管理员角色 |

### Token 配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| access_token 有效期 | 7200s (2小时) | JWT 配置 |
| refresh_token 有效期 | 604800s (7天) | JWT 配置 |

---

## 三、公共接口 (无需认证)

### 3.1 健康检查

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 服务健康检查 |

### 3.2 认证 (auth)

#### 登录注册

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |

**POST /api/v1/auth/login**

```json
// 请求
{
  "username": "admin",
  "password": "admin123",
  "captcha_id": "xxx",
  "captcha_code": "xxxx"
}

// 响应
{
  "code": 200,
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 7200,
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "roles": ["admin"]
    }
  }
}
```

#### 验证码

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/captcha` | 获取验证码 |
| POST | `/api/v1/auth/captcha/verify` | 验证验证码 |

#### 密码重置

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/password/reset-request` | 申请密码重置 |
| POST | `/api/v1/auth/password/reset` | 执行密码重置 |

#### 短信 / 邮箱验证码

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/sms/send` | 发送短信验证码 |
| POST | `/api/v1/auth/email/send` | 发送邮箱验证码 |

#### OAuth 第三方登录

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/oauth/:provider/url` | 获取 OAuth 授权 URL |
| POST | `/api/v1/auth/oauth/:provider` | OAuth 登录回调 |

`provider` 可选值：`wechat`、`github`、`google` 等

---

## 四、用户端接口 (需登录)

> 以下所有接口需要在请求头中携带 `Authorization: Bearer <access_token>`

### 4.1 用户资料 (user)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/profile` | 获取当前用户资料 |
| GET | `/api/v1/users/me` | 获取当前用户资料 (兼容路径) |
| PUT | `/api/v1/user/profile` | 更新个人资料 |
| POST | `/api/v1/user/avatar` | 上传头像 |

#### 隐私设置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/privacy` | 获取隐私设置 |
| PUT | `/api/v1/user/privacy` | 更新隐私设置 |

隐私可见性枚举：`1`=全部可见, `2`=仅班级成员, `3`=仅自己

#### 绑定管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/bind-phone` | 绑定手机号 |
| POST | `/api/v1/user/bind-email` | 绑定邮箱 |

#### 实名认证

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/realname` | 获取实名认证状态 |
| POST | `/api/v1/user/realname` | 提交实名认证 |

#### OAuth 绑定

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/oauth` | 获取 OAuth 绑定列表 |
| DELETE | `/api/v1/user/oauth/:provider` | 解绑 OAuth |

#### 收藏

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/favorites` | 获取收藏列表 |
| POST | `/api/v1/user/favorites` | 添加收藏 |
| DELETE | `/api/v1/user/favorites/:id` | 取消收藏 |
| GET | `/api/v1/favorites` | 收藏列表（内容创作） |
| POST | `/api/v1/favorites` | 添加收藏（内容创作） |
| DELETE | `/api/v1/favorites/:id` | 取消收藏（内容创作） |
| PUT | `/api/v1/favorites/:id` | 更新收藏（修改收藏夹/备注） |
| POST | `/api/v1/favorites/batch-delete` | 批量取消收藏 |

### 4.2 安全与设备 (auth)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/logout` | 退出登录 |
| POST | `/api/v1/user/logout-all` | 退出所有设备 |
| PUT | `/api/v1/user/password` | 修改密码 |
| GET | `/api/v1/user/devices` | 获取登录设备列表 |
| DELETE | `/api/v1/user/devices/:id` | 踢出指定设备 |
| DELETE | `/api/v1/user/devices` | 踢出所有设备 |
| GET | `/api/v1/user/security/logs` | 获取安全日志 |
| PUT | `/api/v1/user/security/settings` | 更新安全设置 |

#### OAuth 绑定 (安全)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/oauth/status` | 获取 OAuth 绑定状态 |
| POST | `/api/v1/user/oauth/bind/:provider` | 绑定 OAuth 账号 |
| DELETE | `/api/v1/user/oauth/unbind/:provider` | 解绑 OAuth 账号 |

#### 账号注销

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/account/deactivate` | 申请注销账号 |
| POST | `/api/v1/user/account/deactivate/confirm` | 确认注销 |
| POST | `/api/v1/user/account/cancel-deactivate` | 取消注销 |
| GET | `/api/v1/user/account/export` | 导出个人数据 |

### 4.3 标签 (tag)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/tags` | 获取所有标签 |

### 4.4 菜单 (menu)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/menu/user` | 获取当前用户菜单 |
| GET | `/api/v1/menu/buttons` | 获取当前用户按钮权限 |
| GET | `/api/v1/menus/routes` | 获取前端动态路由 |

### 4.5 题库 (question)

#### 题目浏览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/questions` | 题目列表 (支持分页、筛选) |
| GET | `/api/v1/questions/search` | 题目搜索 |
| GET | `/api/v1/questions/:id` | 题目详情 |
| POST | `/api/v1/questions/:id/favorite` | 收藏/取消收藏题目 |
| POST | `/api/v1/questions/:id/like` | 点赞/取消点赞题目 |
| GET | `/api/v1/questions/:id/versions` | 题目版本历史 |

#### 分享

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/shares/:code` | 通过分享码查看题目 |

#### 分类与知识点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/categories` | 分类列表 |
| GET | `/api/v1/categories/tree` | 分类树形结构 |
| GET | `/api/v1/knowledge-points` | 知识点列表 |

#### 收藏夹

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/favorites/folders` | 收藏夹列表 |
| POST | `/api/v1/favorites/folders` | 创建收藏夹 |
| PUT | `/api/v1/favorites/folders/:id` | 更新收藏夹 |
| DELETE | `/api/v1/favorites/folders/:id` | 删除收藏夹 |

### 4.6 考试 (exam)

#### 考试列表与详情

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/exam` | 考试列表 |
| GET | `/api/v1/exam/:id` | 考试详情 |
| GET | `/api/v1/exam/history` | 考试历史 |

#### 参加考试

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/exam/:id/start` | 开始考试 |
| POST | `/api/v1/exam/:recordId/submit` | 提交考试 |
| GET | `/api/v1/exam/:id/result` | 查看考试结果 |
| GET | `/api/v1/exam/:id/standard-answers` | 查看标准答案 |
| GET | `/api/v1/exam/:id/guide` | 获取考试指南 |
| POST | `/api/v1/exam/:id/feedback` | 提交考试反馈 |

#### 答题操作

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/exam-records/:recordId/save-answer` | 保存单题答案 |
| POST | `/api/v1/exam-records/:recordId/save-answers` | 批量保存答案 |
| POST | `/api/v1/exam-records/:recordId/mark` | 标记/取消标记题目 |
| GET | `/api/v1/exam-records/:recordId/marked` | 获取标记题目列表 |
| POST | `/api/v1/exam-records/:recordId/warning` | 上报考试异常行为 |

#### 考试公告与排名

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/exam/:id/notices` | 考试公告列表 |
| GET | `/api/v1/exam/:id/ranking` | 考试排名 |
| GET | `/api/v1/exam/:id/ranking/me` | 我的排名 |
| POST | `/api/v1/exam/:id/verify-password` | 验证考试密码 |

### 4.7 班级 (class)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class` | 班级列表 |
| GET | `/api/v1/class/:id` | 班级详情 |
| POST | `/api/v1/class` | 创建班级 |
| PUT | `/api/v1/class/:id` | 更新班级信息 |
| DELETE | `/api/v1/class/:id` | 删除班级 |
| POST | `/api/v1/class/:id/join` | 加入班级 |
| POST | `/api/v1/class/:id/leave` | 退出班级 |
| GET | `/api/v1/class/:id/members` | 班级成员列表 |
| GET | `/api/v1/class/:id/notices` | 班级通知列表 |
| GET | `/api/v1/class/:id/homework` | 班级作业列表 |

### 4.8 刷题练习 (practice)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice` | 练习历史 |
| GET | `/api/v1/practice/:id` | 练习详情 |
| GET | `/api/v1/practice/stats` | 练习统计 |
| POST | `/api/v1/practice/start` | 开始练习 |
| POST | `/api/v1/practice/submit` | 提交练习 |

**POST /api/v1/practice/start**

```json
// 请求
{
  "category_id": 1,
  "question_count": 20,
  "difficulty": 2,
  "question_type": 1
}

// 响应
{
  "code": 200,
  "data": {
    "session_id": 1,
    "questions": [...]
  }
}
```

### 4.9 错题本 (wrong)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/wrong-questions` | 错题列表 |
| GET | `/api/v1/wrong-questions/:id` | 错题详情 |
| POST | `/api/v1/wrong-questions/:id/review` | 复习错题 |
| DELETE | `/api/v1/wrong-questions/:id` | 移除错题 |
| GET | `/api/v1/wrong-questions/analysis` | 错题分析 |

### 4.10 评论 (comment)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/comments` | 评论列表 (通过 target_type + target_id 查询) |
| POST | `/api/v1/comments` | 发表评论 |
| DELETE | `/api/v1/comments/:id` | 删除评论 |
| POST | `/api/v1/comments/:id/like` | 点赞评论 |
| POST | `/api/v1/comments/:id/report` | 举报评论 |

评论目标类型：`1`=题目, `2`=考试, `3`=班级

### 4.11 通知 (notification)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/notifications` | 通知列表 |
| GET | `/api/v1/notifications/unread-count` | 未读通知数量 |
| PUT | `/api/v1/notifications/:id/read` | 标记单条已读 |
| PUT | `/api/v1/notifications/read-all` | 标记全部已读 |
| DELETE | `/api/v1/notifications/:id` | 删除通知 |

通知类型：`1`=系统, `2`=考试, `3`=作业, `4`=班级, `5`=评论

### 4.12 数据统计 (statistics)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/statistics/overview` | 仪表盘概览 |
| GET | `/api/v1/statistics/users` | 用户统计 |
| GET | `/api/v1/statistics/questions` | 练习统计 |
| GET | `/api/v1/statistics/exams` | 考试统计 |
| GET | `/api/v1/statistics/classes` | 班级统计 |
| GET | `/api/v1/statistics/ranking` | 排行榜 |

### 4.13 文件 (file)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/files` | 上传文件 |
| DELETE | `/api/v1/files` | 删除文件 |

### 4.14 角色申请 (application)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/applications` | 提交角色申请 |
| GET | `/api/v1/applications` | 我的申请列表 |
| GET | `/api/v1/applications/:id` | 申请详情 |
| POST | `/api/v1/user/apply/creator` | 申请成为创作者 |
| POST | `/api/v1/user/apply/teacher` | 申请成为教师 |
| GET | `/api/v1/user/apply/status` | 查询申请状态 |


### 内容创作

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/creator/level` | 获取当前创作者等级信息 |
| GET | `/api/v1/creator/points` | 获取创作者积分明细 |
| GET | `/api/v1/creator/points/logs` | 获取积分变动记录 |
| GET | `/api/v1/creator/levels` | 获取所有创作者等级列表 |
| GET | `/api/v1/creator/agreement/latest` | 获取最新创作者协议 |
| POST | `/api/v1/creator/agreement/sign` | 签署创作者协议 |
| GET | `/api/v1/creator/agreement/signed` | 查询协议签署状态 |
| GET | `/api/v1/creator/portfolios` | 获取我的作品集列表 |
| POST | `/api/v1/creator/portfolios` | 创建作品集 |
| PUT | `/api/v1/creator/portfolios/:id` | 更新作品集 |
| DELETE | `/api/v1/creator/portfolios/:id` | 删除作品集 |
| POST | `/api/v1/creator/portfolios/:id/items` | 添加作品到作品集 |
| DELETE | `/api/v1/creator/portfolios/:id/items/:itemId` | 从作品集移除作品 |
| GET | `/api/v1/creator/portfolios/:id` | 获取作品集详情 |
| GET | `/api/v1/users/:id/portfolios` | 查看他人作品集 |
| GET | `/api/v1/content/versions` | 获取内容版本列表 |
| GET | `/api/v1/content/versions/:id` | 获取内容版本详情 |
| POST | `/api/v1/content/versions/:id/rollback` | 回滚到指定版本 |
| POST | `/api/v1/content/tags` | 给内容打标签 |
| DELETE | `/api/v1/content/tags` | 移除内容标签 |
| GET | `/api/v1/content/:type/:id/tags` | 获取内容标签列表 |
| POST | `/api/v1/tags` | 创建标签 |
| PUT | `/api/v1/tags/:id` | 更新标签 |
| DELETE | `/api/v1/tags/:id` | 删除标签 |
| GET | `/api/v1/content/:type/:id/preview` | 内容预览 |
| POST | `/api/v1/content/preview` | 预览草稿内容 |
| PUT | `/api/v1/creator/contents/:id` | 编辑创作者内容 |
| DELETE | `/api/v1/creator/contents/:id` | 删除创作者内容 |
| POST | `/api/v1/creator/contents/:id/copy` | 复制创作者内容 |
| GET | `/api/v1/creator/contents/:id/export` | 导出创作者内容（Word/PDF） |
| POST | `/api/v1/creator/contents/:id/share` | 生成内容分享链接/二维码 |
| GET | `/api/v1/creator/contents/:id/share/:token` | 通过分享链接访问内容 |

### 内容搜索

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/search` | 综合搜索（题目、分类、知识点） |
| GET | `/api/v1/search/suggestions` | 搜索建议（自动补全） |
| GET | `/api/v1/search/history` | 获取搜索历史 |
| DELETE | `/api/v1/search/history` | 清空搜索历史 |

### 刷题练习

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/practice/:id/answer` | 提交练习答案（即时评判） |
| POST | `/api/v1/practice/:id/pause` | 暂停练习 |
| POST | `/api/v1/practice/:id/resume` | 继续练习 |
| POST | `/api/v1/practice/:id/finish` | 完成练习 |
| GET | `/api/v1/practice/:id/result` | 获取练习结果 |
| GET | `/api/v1/practice/resume-list` | 获取可继续的练习列表 |
| GET | `/api/v1/practice/history` | 获取练习历史记录 |

### 审核流程

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/reviews/submit` | 提交内容审核 |
| POST | `/api/v1/reviews/:id/withdraw` | 撤回审核申请 |
| POST | `/api/v1/reviews/:id/approve` | 审核通过 |
| POST | `/api/v1/reviews/:id/reject` | 审核拒绝 |
| GET | `/api/v1/reviews/:id` | 获取审核详情 |
| GET | `/api/v1/reviews/:id/steps` | 获取审核步骤记录 |
| GET | `/api/v1/reviews/:id/replies` | 获取审核意见回复列表 |
| POST | `/api/v1/reviews/:id/replies` | 回复审核意见 |
| DELETE | `/api/v1/replies/:replyId` | 删除自己的审核回复 |
| GET | `/api/v1/review-notifications` | 获取审核通知列表 |
| GET | `/api/v1/review-notifications/unread-count` | 获取未读审核通知数量 |
| POST | `/api/v1/review-notifications/:id/read` | 标记审核通知已读 |
| POST | `/api/v1/review-notifications/read-all` | 全部审核通知标记已读 |
| DELETE | `/api/v1/review-notifications/:id` | 删除审核通知 |

### 成绩管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/exam/:id/score-review` | 学生申请成绩复核 |
| GET | `/api/v1/admin/score-reviews` | 获取成绩复核申请列表 |
| PUT | `/api/v1/admin/score-reviews/:id` | 处理成绩复核（通过/驳回并给出复核结果） |
| POST | `/api/v1/admin/exams/:id/notice` | 发布考试公告 |
| GET | `/api/v1/admin/exams/:id/notices` | 获取考试公告列表 |

### 数据统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/statistics/dashboard` | 获取仪表盘数据（总用户数、今日活跃、数据趋势等） |
| POST | `/api/v1/admin/statistics/export` | 创建统计导出任务 |
| GET | `/api/v1/admin/statistics/export/list` | 获取导出记录列表 |
| GET | `/api/v1/admin/statistics/export/:id/download` | 下载导出文件 |
| DELETE | `/api/v1/admin/statistics/export/:id` | 删除导出记录 |
| GET | `/api/v1/admin/statistics/compare` | 数据对比（不同时间段） |
| GET | `/api/v1/admin/statistics/compare/metrics` | 获取对比指标列表 |
| GET | `/api/v1/admin/statistics/retention` | 获取用户留存率统计 |
| GET | `/api/v1/admin/statistics/retention/trend` | 获取留存率趋势 |
| GET | `/api/v1/admin/statistics/churn` | 获取用户流失统计 |
| GET | `/api/v1/admin/statistics/churn/trend` | 获取流失趋势 |
| POST | `/api/v1/admin/statistics/events` | 上报用户行为事件 |
| GET | `/api/v1/admin/statistics/events/analysis` | 获取用户行为分析 |
| GET | `/api/v1/admin/statistics/events/overview` | 获取行为概览 |
| GET | `/api/v1/admin/statistics/paths` | 获取用户访问路径分析 |
| GET | `/api/v1/admin/statistics/paths/flow` | 获取页面流转分析 |
| GET | `/api/v1/admin/statistics/segments` | 获取用户分群列表 |
| POST | `/api/v1/admin/statistics/segments` | 创建用户分群 |
| PUT | `/api/v1/admin/statistics/segments/:id` | 更新用户分群 |
| DELETE | `/api/v1/admin/statistics/segments/:id` | 删除用户分群 |
| POST | `/api/v1/admin/statistics/segments/:id/sync` | 同步分群成员 |
| GET | `/api/v1/admin/statistics/segments/:id/members` | 获取分群成员列表 |
| GET | `/api/v1/admin/statistics/funnels` | 获取转化漏斗列表 |
| POST | `/api/v1/admin/statistics/funnels` | 创建转化漏斗 |
| GET | `/api/v1/admin/statistics/funnels/:id/stats` | 获取漏斗转化统计 |
| GET | `/api/v1/admin/statistics/alerts/rules` | 获取预警规则列表 |
| POST | `/api/v1/admin/statistics/alerts/rules` | 创建预警规则 |
| PUT | `/api/v1/admin/statistics/alerts/rules/:id` | 更新预警规则 |
| DELETE | `/api/v1/admin/statistics/alerts/rules/:id` | 删除预警规则 |
| GET | `/api/v1/admin/statistics/alerts/records` | 获取预警记录列表 |
| PUT | `/api/v1/admin/statistics/alerts/records/:id/handle` | 处理预警记录 |
| GET | `/api/v1/statistics/subscriptions` | 获取我的数据订阅列表 |
| POST | `/api/v1/statistics/subscriptions` | 创建数据订阅 |
| PUT | `/api/v1/statistics/subscriptions/:id` | 更新数据订阅 |
| DELETE | `/api/v1/statistics/subscriptions/:id` | 删除数据订阅 |
| POST | `/api/v1/statistics/refresh` | 手动刷新统计数据 |
| POST | `/api/v1/statistics/share` | 分享数据报告 |
| GET | `/api/v1/statistics/share/:token` | 获取分享的数据报告 |
| GET | `/api/v1/statistics/practice/dimensions` | 获取答题统计详细维度 |
| GET | `/api/v1/statistics/score/dimensions` | 获取成绩统计详细维度 |
| GET | `/api/v1/statistics/question/difficulty` | 获取题目难度分析 |
| GET | `/api/v1/statistics/question/discrimination` | 获取题目区分度分析 |
| GET | `/api/v1/statistics/score/prediction` | 获取成绩预测 |
| GET | `/api/v1/statistics/score/alert/config` | 获取成绩预警配置 |
| POST | `/api/v1/statistics/score/alert/config` | 设置成绩预警 |
| GET | `/api/v1/statistics/class/overview` | 获取班级概览统计（教师视角） |
| GET | `/api/v1/statistics/class/:id/progress` | 获取班级成员学习进度（教师视角） |
| GET | `/api/v1/statistics/class/:id/analysis` | 获取班级成绩分析（教师视角） |
| GET | `/api/v1/statistics/class/:id/creators` | 获取班级创作者统计（教师视角） |
| GET | `/api/v1/statistics/class/:id/activity` | 获取班级活跃度统计（教师视角） |
| GET | `/api/v1/statistics/class/:id/ranking` | 获取班级排名统计（教师视角） |
| GET | `/api/v1/statistics/class/:id/report` | 生成班级报告（教师视角） |
| GET | `/api/v1/statistics/class/compare` | 班级间对比统计（教师视角） |
| GET | `/api/v1/statistics/class/:id/warning` | 获取班级预警（教师视角） |
| GET | `/api/v1/statistics/class/:id/trend` | 获取班级趋势分析（教师视角） |
| GET | `/api/v1/statistics/mobile/overview` | 获取移动端个人统计概览 |
| GET | `/api/v1/statistics/mobile/practice` | 获取移动端练习统计 |
| GET | `/api/v1/statistics/mobile/accuracy-trend` | 获取移动端正确率趋势 |
| GET | `/api/v1/statistics/mobile/category-analysis` | 获取移动端分类分析 |

### 文件管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/file/upload/image` | 上传图片（带压缩+缩略图+水印） |
| POST | `/api/v1/file/upload/batch` | 批量上传文件（最多20个） |
| GET | `/api/v1/file/:id` | 获取文件信息 |
| GET | `/api/v1/file/:id/download` | 下载文件 |
| GET | `/api/v1/file/:id/thumbnail/:size` | 获取图片缩略图（100x100/300x300/800x800） |
| DELETE | `/api/v1/file/:id` | 删除文件（检查引用计数） |
| POST | `/api/v1/admin/files/batch/delete` | 批量删除文件 |
| GET | `/api/v1/file/list` | 获取文件列表（管理端，支持按类型/大小/时间筛选） |
| POST | `/api/v1/file/:id/reference` | 添加文件引用 |
| DELETE | `/api/v1/file/:id/reference` | 移除文件引用 |
| GET | `/api/v1/file/:id/references` | 查看文件引用列表 |
| GET | `/api/v1/file/access-logs` | 查询文件访问日志 |
| GET | `/api/v1/file/access-logs/statistics` | 获取文件访问日志统计 |
| GET | `/api/v1/file/hotlink-rules` | 获取防盗链规则列表 |
| POST | `/api/v1/file/hotlink-rules` | 创建防盗链规则 |
| PUT | `/api/v1/file/hotlink-rules/:id` | 更新防盗链规则 |
| DELETE | `/api/v1/file/hotlink-rules/:id` | 删除防盗链规则 |
| GET | `/api/v1/file/watermark-configs` | 获取水印配置列表 |
| POST | `/api/v1/file/watermark-configs` | 创建水印配置 |
| PUT | `/api/v1/file/watermark-configs/:id` | 更新水印配置 |
| DELETE | `/api/v1/file/watermark-configs/:id` | 删除水印配置 |
| POST | `/api/v1/file/watermark-configs/:id/default` | 设置默认水印配置 |
| POST | `/api/v1/file/watermark/preview` | 预览水印效果 |
| POST | `/api/v1/file/cleanup/orphan` | 手动触发孤立文件清理 |
| GET | `/api/v1/file/cleanup/logs` | 查看文件清理记录 |
| GET | `/api/v1/file/storage/statistics` | 获取存储使用统计（总量、各类型占比、增长趋势） |
| POST | `/api/v1/file/:id/virus-scan` | 手动触发文件病毒扫描 |
| POST | `/api/v1/file/:id/content-check` | 手动触发文件敏感内容检测 |

### 最近浏览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/recent-categories` | 获取最近浏览的分类列表 |
| POST | `/api/v1/practice/recent-categories` | 记录分类浏览 |
| DELETE | `/api/v1/practice/recent-categories/:id` | 删除最近浏览记录 |

### 练习提醒

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/reminder/config` | 获取用户练习提醒配置 |
| PUT | `/api/v1/practice/reminder/config` | 更新用户练习提醒配置（enabled, reminder_time, reminder_days, channel） |

### 模拟考试

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/practice/mock/start` | 开始模拟考试（限时、计分） |
| POST | `/api/v1/practice/mock/:id/answer` | 模拟考试答题 |
| POST | `/api/v1/practice/mock/:id/submit` | 模拟考试交卷 |
| GET | `/api/v1/practice/mock/:id/result` | 获取模拟考试结果 |
| GET | `/api/v1/practice/mock/configs` | 获取模拟考试配置列表 |
| POST | `/api/v1/practice/mock/configs` | 保存模拟考试配置 |

### 每日练习

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/daily` | 获取今日推荐题目 |
| POST | `/api/v1/practice/daily/start` | 开始今日练习 |

### 班级作业管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/class/:id/homework/:homeworkId/peer-reviews/assign` | 分配学生互评（教师操作） |
| GET | `/api/v1/class/:id/homework/:homeworkId/peer-reviews` | 获取互评列表 |
| GET | `/api/v1/class/:id/homework/:homeworkId/peer-reviews/mine` | 获取我的互评任务 |
| POST | `/api/v1/class/:id/homework/:homeworkId/peer-reviews/:reviewId` | 提交互评评分 |
| GET | `/api/v1/class/:id/homework/:homeworkId/peer-reviews/result` | 查看互评结果 |

### 班级分组管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/groups` | 获取班级分组列表 |
| POST | `/api/v1/class/:id/groups` | 创建班级分组（学习小组） |
| PUT | `/api/v1/class/:id/groups/:groupId` | 编辑班级分组 |
| DELETE | `/api/v1/class/:id/groups/:groupId` | 删除班级分组 |
| POST | `/api/v1/class/:id/groups/:groupId/members` | 添加分组成员（支持批量） |
| DELETE | `/api/v1/class/:id/groups/:groupId/members/:userId` | 移除分组成员 |

### 班级创作者管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/creator-applications` | 获取班级创作者申请列表（教师视角） |
| POST | `/api/v1/class/:id/creator-apply` | 申请成为班级创作者 |
| POST | `/api/v1/class/:id/creator-applications/:appId/approve` | 审批通过班级创作者申请 |
| POST | `/api/v1/class/:id/creator-applications/:appId/reject` | 审批驳回班级创作者申请 |
| DELETE | `/api/v1/class/:id/creators/:userId` | 撤销班级创作者权限 |
| GET | `/api/v1/class/:id/creators` | 获取班级创作者列表 |

### 班级学习计划

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/study-plans` | 获取班级学习计划列表 |
| POST | `/api/v1/class/:id/study-plans` | 创建学习计划 |
| GET | `/api/v1/class/:id/study-plans/:planId` | 获取学习计划详情（含任务列表和进度） |
| PUT | `/api/v1/class/:id/study-plans/:planId` | 编辑学习计划 |
| DELETE | `/api/v1/class/:id/study-plans/:planId` | 删除学习计划 |
| POST | `/api/v1/class/:id/study-plans/:planId/items` | 添加学习计划任务 |
| PUT | `/api/v1/class/:id/study-plans/:planId/items/:itemId` | 编辑学习计划任务 |
| DELETE | `/api/v1/class/:id/study-plans/:planId/items/:itemId` | 删除学习计划任务 |
| POST | `/api/v1/class/:id/study-plans/:planId/items/:itemId/complete` | 标记学习计划任务完成 |
| GET | `/api/v1/class/:id/study-plans/:planId/progress` | 查看学习计划完成进度 |

### 班级成员管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/class/:id/members/import` | Excel批量导入班级成员 |
| POST | `/api/v1/class/:id/members/:userId/disable` | 禁用班级成员（不影响其他班级） |
| POST | `/api/v1/class/:id/members/:userId/enable` | 启用班级成员 |
| GET | `/api/v1/class/:id/members/export` | 导出成员列表为Excel |

### 班级排名

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/ranking` | 获取班级排名（支持单次/总排名） |
| POST | `/api/v1/class/:id/ranking/calculate` | 触发排名计算（教师操作） |

### 班级文件管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/files` | 获取班级共享文件列表 |
| POST | `/api/v1/class/:id/files` | 上传班级共享文件 |
| DELETE | `/api/v1/class/:id/files/:fileId` | 删除班级共享文件 |
| GET | `/api/v1/class/:id/files/:fileId/download` | 下载班级共享文件 |

### 班级标签管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/tags` | 获取所有班级标签列表 |
| POST | `/api/v1/class/tags` | 创建班级标签 |
| PUT | `/api/v1/class/tags/:tagId` | 编辑班级标签 |
| DELETE | `/api/v1/class/tags/:tagId` | 删除班级标签 |
| POST | `/api/v1/class/:id/tags` | 为班级添加标签 |
| DELETE | `/api/v1/class/:id/tags/:tagId` | 移除班级标签 |

### 班级模板

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/templates` | 获取班级模板列表 |
| POST | `/api/v1/class/templates` | 创建班级模板 |
| GET | `/api/v1/class/templates/:templateId` | 获取班级模板详情 |
| PUT | `/api/v1/class/templates/:templateId` | 编辑班级模板 |
| DELETE | `/api/v1/class/templates/:templateId` | 删除班级模板 |
| POST | `/api/v1/class/templates/:templateId/create` | 从模板创建班级 |

### 班级申请管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/class/:id/apply` | 提交加入班级申请 |
| GET | `/api/v1/class/:id/applications` | 获取班级申请列表（教师视角） |
| POST | `/api/v1/class/:id/applications/:appId/approve` | 审批通过加入申请 |
| POST | `/api/v1/class/:id/applications/:appId/reject` | 审批驳回加入申请 |

### 班级管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/class/:id/archive` | 归档班级（只读，禁止新增成员/发布考试） |
| POST | `/api/v1/class/:id/unarchive` | 取消归档班级 |
| POST | `/api/v1/class/:id/pin` | 置顶班级 |
| POST | `/api/v1/class/:id/unpin` | 取消置顶班级 |
| GET | `/api/v1/class/search` | 按关键词搜索班级 |
| GET | `/api/v1/class/:id/qrcode` | 生成班级二维码（返回Base64图片或图片URL） |
| PUT | `/api/v1/class/:id/expire` | 设置/修改班级有效期 |
| GET | `/api/v1/class/:id/exams` | 获取班级考试列表 |
| POST | `/api/v1/class/:id/exams` | 在班级内发布考试 |
| GET | `/api/v1/class/:id/notice` | 获取班级公告列表 |

### 班级考勤管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/attendances` | 获取班级考勤列表 |
| POST | `/api/v1/class/:id/attendances` | 创建考勤（教师发起签到/签退） |
| PUT | `/api/v1/class/:id/attendances/:attId` | 编辑考勤记录 |
| DELETE | `/api/v1/class/:id/attendances/:attId` | 删除考勤记录 |
| POST | `/api/v1/class/:id/attendances/:attId/checkin` | 学生签到 |
| POST | `/api/v1/class/:id/attendances/:attId/checkout` | 学生签退 |
| GET | `/api/v1/class/:id/attendances/:attId/records` | 获取考勤记录详情 |
| GET | `/api/v1/class/:id/attendances/:attId/export` | 导出考勤记录 |
| POST | `/api/v1/class/:id/attendances/:attId/records/:recordId` | 教师修改考勤状态 |

### 班级讨论

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/discussions` | 获取班级讨论列表（置顶优先） |
| POST | `/api/v1/class/:id/discussions` | 发布班级讨论帖 |
| GET | `/api/v1/class/:id/discussions/:discussionId` | 获取讨论详情（含回复列表） |
| PUT | `/api/v1/class/:id/discussions/:discussionId` | 编辑讨论帖 |
| DELETE | `/api/v1/class/:id/discussions/:discussionId` | 删除讨论帖 |
| POST | `/api/v1/class/:id/discussions/:discussionId/pin` | 置顶讨论帖 |
| POST | `/api/v1/class/:id/discussions/:discussionId/unpin` | 取消置顶讨论帖 |

#### 讨论回复管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/class/:id/discussions/:discussionId/replies` | 回复讨论帖（支持嵌套回复） |
| PUT | `/api/v1/class/:id/discussions/:discussionId/replies/:replyId` | 编辑回复 |
| DELETE | `/api/v1/class/:id/discussions/:discussionId/replies/:replyId` | 删除回复 |

### 班级资源管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/resources` | 获取班级资源列表（题目、分类、知识点） |
| GET | `/api/v1/class/:id/resources/statistics` | 获取班级资源统计（题目数、分类数、知识点数） |
| POST | `/api/v1/class/:id/resources/import` | 从公共题库导入题目到班级 |
| GET | `/api/v1/class/:id/resources/export` | 导出班级资源 |

#### 班级资源搜索

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/resources/search` | 在班级资源库中搜索（支持按类型、标签、知识点筛选） |

#### 班级资源版本管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/resources/:resourceType/:resourceId/versions` | 资源版本历史列表 |
| GET | `/api/v1/class/:id/resources/:resourceType/:resourceId/versions/:versionId` | 版本详情（含内容快照） |
| POST | `/api/v1/class/:id/resources/:resourceType/:resourceId/versions/:versionId/rollback` | 回滚到指定版本 |

#### 班级资源标签

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/resource-tags` | 获取班级资源标签列表 |
| POST | `/api/v1/class/:id/resource-tags` | 创建资源标签 |
| PUT | `/api/v1/class/:id/resource-tags/:tagId` | 编辑资源标签 |
| DELETE | `/api/v1/class/:id/resource-tags/:tagId` | 删除资源标签 |
| POST | `/api/v1/class/:id/resources/:resourceType/:resourceId/tags` | 为资源打标签 |
| DELETE | `/api/v1/class/:id/resources/:resourceType/:resourceId/tags/:tagId` | 移除资源标签 |

#### 班级资源审核

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/resource-reviews` | 获取待审核资源列表 |
| POST | `/api/v1/class/:id/resource-reviews/:reviewId/approve` | 审核通过 |
| POST | `/api/v1/class/:id/resource-reviews/:reviewId/reject` | 审核驳回 |
| PUT | `/api/v1/class/:id/resource-review-config` | 配置是否需要审核（针对创作者角色） |
| GET | `/api/v1/class/:id/resource-review-config` | 获取审核配置 |

### 班级作业模板管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/homework-templates` | 获取班级作业模板列表 |
| POST | `/api/v1/class/:id/homework-templates` | 创建作业模板 |
| GET | `/api/v1/class/:id/homework-templates/:templateId` | 获取作业模板详情 |
| PUT | `/api/v1/class/:id/homework-templates/:templateId` | 编辑作业模板 |
| DELETE | `/api/v1/class/:id/homework-templates/:templateId` | 删除作业模板 |
| POST | `/api/v1/class/:id/homework-templates/:templateId/create` | 从模板创建作业 |

### 班级成绩预警

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/score-alerts` | 获取成绩预警列表（支持按状态筛选） |
| POST | `/api/v1/class/:id/score-alerts/rules` | 配置预警规则（阈值、提醒方式） |
| GET | `/api/v1/class/:id/score-alerts/rules` | 获取预警规则 |
| POST | `/api/v1/class/:id/score-alerts/:alertId/confirm` | 确认预警（教师操作） |
| GET | `/api/v1/class/:id/score-alerts/statistics` | 获取预警统计（各科目预警人数） |

### 班级学习报告

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/students/:studentId/report` | 生成学生学习报告（成绩、进度、参与度） |
| GET | `/api/v1/class/:id/students/report/batch` | 批量生成学习报告 |
| GET | `/api/v1/class/:id/report/template` | 获取报告模板配置 |
| PUT | `/api/v1/class/:id/report/template` | 配置报告模板 |

### 班级日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/logs` | 获取班级日志列表（支持按操作类型、时间、操作人筛选） |
| GET | `/api/v1/class/:id/logs/actions` | 获取日志操作类型列表 |

### 班级成员统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/members/statistics` | 获取成员整体统计（总人数、活跃人数、新增人数等） |
| GET | `/api/v1/class/:id/members/:userId/activity` | 获取单个成员活跃度详情（登录、参与考试、提交作业等） |
| GET | `/api/v1/class/:id/members/participation` | 获取成员参与度排行 |

### 班级创作者权限配置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/creator-permissions` | 获取创作者权限列表 |
| PUT | `/api/v1/class/:id/creator-permissions` | 批量配置创作者权限 |
| PUT | `/api/v1/class/:id/creator-permissions/:permissionKey` | 更新单个权限项 |

### 班级创作者统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/class/:id/creators/statistics` | 获取创作者整体统计 |
| GET | `/api/v1/class/:id/creators/:userId/statistics` | 获取单个创作者统计（创建题目数、活跃度等） |
| GET | `/api/v1/class/:id/creators/ranking` | 获取创作者贡献排行 |

### 移动端

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/user/stats/overview` | 获取用户个人统计概览（总题量、今日练习数、正确率、连续打卡天数） |
| GET | `/api/v1/categories/hot` | 获取热门分类列表 |
| GET | `/api/v1/questions/recommend` | 获取推荐题目列表 |
| GET | `/api/v1/classes/:id` | 获取班级详情（移动端路径） |
| GET | `/api/v1/creator/questions` | 获取创作者自己创建的题目列表 |
| GET | `/api/v1/creator/papers` | 获取创作者自己创建的试卷列表 |
| GET | `/api/v1/user/class-creator-status` | 获取班级创作者申请状态 |
| GET | `/api/v1/study-plans` | 获取学习计划列表 |
| GET | `/api/v1/study-plans/today` | 获取今日学习计划 |
| POST | `/api/v1/study-plans` | 创建学习计划 |
| PUT | `/api/v1/study-plans/:id` | 更新学习计划 |
| DELETE | `/api/v1/study-plans/:id` | 删除学习计划 |
| GET | `/api/v1/study-reports` | 获取学习报告 |
| GET | `/api/v1/checkin` | 获取打卡记录 |
| POST | `/api/v1/checkin` | 执行每日打卡 |
| POST | `/api/v1/checkin/makeup` | 补打卡 |
| GET | `/api/v1/leaderboard` | 获取排行榜（支持type和scope参数） |
| GET | `/api/v1/daily-recommend` | 获取每日推荐题目 |
| POST | `/api/v1/daily-recommend/refresh` | 刷新每日推荐题目 |
| GET | `/api/v1/notifications/unread-counts` | 获取未读通知数量（按类型分组） |
| POST | `/api/v1/practice/records/sync` | 同步练习记录（离线数据上传） |
| POST | `/api/v1/exam/answers/sync` | 同步考试答案（离线数据上传） |
| GET | `/api/v1/practice/records` | 获取练习记录列表 |

### 练习打卡

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/practice/checkin` | 手动打卡 |
| GET | `/api/v1/practice/checkin/status` | 获取今日打卡状态 |
| GET | `/api/v1/practice/checkin/calendar` | 获取打卡日历 |
| GET | `/api/v1/practice/checkin/streak` | 获取连续打卡信息 |

### 练习报告

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/report` | 获取练习报告（含分析和建议） |
| GET | `/api/v1/practice/report/daily` | 获取每日练习报告 |
| GET | `/api/v1/practice/report/weekly` | 获取每周练习报告 |

### 练习排行榜

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/leaderboard` | 获取练习排行榜列表（按题数、正确率等） |
| GET | `/api/v1/practice/leaderboard/rank` | 获取我的排名 |

### 练习统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/statistics` | 获取练习统计概览（总题数、正确率等） |
| GET | `/api/v1/practice/statistics/trend` | 获取练习趋势数据 |
| GET | `/api/v1/practice/statistics/category` | 获取各分类练习统计 |

### 练习计划

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/plans` | 获取练习计划列表 |
| POST | `/api/v1/practice/plans` | 创建练习计划 |
| PUT | `/api/v1/practice/plans/:id` | 更新练习计划 |
| DELETE | `/api/v1/practice/plans/:id` | 删除练习计划 |
| PUT | `/api/v1/practice/plans/:id/pause` | 暂停练习计划 |
| PUT | `/api/v1/practice/plans/:id/resume` | 恢复练习计划 |
| GET | `/api/v1/practice/plans/:id/progress` | 获取练习计划进度 |

### 考试管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/exams/:id/scores` | 获取考试成绩列表 |
| GET | `/api/v1/admin/exams/:id/scores/export` | 导出考试成绩为Excel |
| GET | `/api/v1/admin/exams/:id/statistics` | 获取考试数据统计（参与率、完成率、平均用时） |
| POST | `/api/v1/admin/exams/:id/extend` | 延长考试截止时间（记录延期原因） |
| POST | `/api/v1/admin/exams/:id/pause` | 暂停考试（已开始考生可继续，新考生无法进入） |
| POST | `/api/v1/admin/exams/:id/resume` | 恢复考试 |
| PUT | `/api/v1/admin/exams/:id/publish` | 发布考试 |
| PUT | `/api/v1/admin/exams/:id/close` | 结束考试 |
| POST | `/api/v1/exam/:id/switch-screen` | 上报考试切屏事件（防作弊） |
| GET | `/api/v1/exam/:id/resume` | 断线续考恢复（读取Redis答题进度） |
| POST | `/api/v1/exam/:recordId/save-answer` | 保存单题答案（自动保存/防丢失） |
| POST | `/api/v1/exam/:recordId/save-answers` | 批量保存答案 |
| POST | `/api/v1/exam/:recordId/mark` | 标记/取消标记题目 |
| GET | `/api/v1/exam/:recordId/marked` | 获取已标记题目列表 |
| POST | `/api/v1/exam/:recordId/warning` | 上报考试异常行为 |
| POST | `/api/v1/exam/:id/feedback` | 提交考后反馈（难度评价、建议、满意度） |
| GET | `/api/v1/exam/:id/guide` | 获取答题指引/操作说明 |
| GET | `/api/v1/exams/upcoming` | 获取即将开始的考试列表（公共+班级） |
| GET | `/api/v1/exams/:id/rankings` | 获取考试成绩排名 |

### 计时模式

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/practice/timed/start` | 开始计时练习 |
| POST | `/api/v1/practice/timed/:id/answer` | 计时练习答题 |
| POST | `/api/v1/practice/timed/:id/submit` | 提交计时练习（超时自动交卷） |
| GET | `/api/v1/practice/timed/:id/result` | 获取计时练习结果 |

### 评论系统

| 方法 | 路径 | 说明 |
|------|------|------|
| PUT | `/api/v1/comment/:id` | 编辑评论 |
| POST | `/api/v1/comment/:id/pin` | 置顶/取消置顶评论（管理员/教师） |
| POST | `/api/v1/comment/:id/featured` | 精选/取消精选评论（管理员/教师） |
| POST | `/api/v1/comment/:id/official` | 标记/取消官方解答 |
| POST | `/api/v1/comment/upload-image` | 上传评论图片 |
| GET | `/api/v1/sticker/list` | 获取表情包列表 |
| GET | `/api/v1/user/search` | 搜索用户（@功能用） |

### 通知系统

| 方法 | 路径 | 说明 |
|------|------|------|
| DELETE | `/api/v1/notification/batch-delete` | 批量删除通知 |
| GET | `/api/v1/notification/settings` | 获取通知偏好设置 |
| PUT | `/api/v1/notification/settings` | 更新通知偏好设置 |

### 错题本分享

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/:id/share` | 生成错题分享链接 |
| GET | `/api/v1/wrong/share/:code` | 通过分享码查看错题 |

### 错题本分析

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/wrong/search` | 错题搜索 |
| POST | `/api/v1/wrong` | 手动添加错题 |
| GET | `/api/v1/wrong/analysis/trend` | 获取错题趋势（近7/30/90天） |
| GET | `/api/v1/wrong/analysis/weak` | 获取薄弱知识点（各分类错误率排行） |
| GET | `/api/v1/wrong/analysis/heatmap` | 获取错题时间分布热力图 |
| GET | `/api/v1/wrong/analysis/graph` | 获取错题知识点关联图谱 |
| GET | `/api/v1/wrong/analysis/predict` | 获取错题预测（可能再次犯错的题目） |
| GET | `/api/v1/wrong/analysis/report` | 生成错题分析报告 |
| GET | `/api/v1/wrong/compare` | 错题对比（不同时间段） |
| GET | `/api/v1/wrong/compare/summary` | 错题对比摘要（错题数量变化、掌握率变化） |
| GET | `/api/v1/wrong/ranking` | 错题排行榜（按数量/掌握率/进步） |

### 错题本同步

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/sync` | 同步错题（多设备同步） |
| GET | `/api/v1/wrong/sync/status` | 获取同步状态 |

### 错题本备注

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/wrong/:id/notes` | 获取错题备注列表 |
| POST | `/api/v1/wrong/:id/notes` | 添加错题备注 |
| PUT | `/api/v1/wrong/:id/notes/:noteId` | 更新错题备注 |
| DELETE | `/api/v1/wrong/:id/notes/:noteId` | 删除错题备注 |

### 错题本导入

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/import` | 导入错题 |

### 错题本导出

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/export` | 导出错题（Word/PDF格式） |

### 错题本打印

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/print` | 生成错题打印内容 |

### 错题本批量操作

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/batch/remove` | 批量移除错题 |
| POST | `/api/v1/wrong/batch/favorite` | 批量收藏错题 |
| POST | `/api/v1/wrong/batch/unfavorite` | 批量取消收藏错题 |
| POST | `/api/v1/wrong/batch/tag` | 批量给错题打标签 |
| POST | `/api/v1/wrong/batch/export` | 批量导出错题 |

### 错题本收藏

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/:id/favorite` | 收藏错题 |
| DELETE | `/api/v1/wrong/:id/favorite` | 取消收藏错题 |
| GET | `/api/v1/wrong/favorites` | 获取收藏错题列表 |

### 错题本标签

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/wrong/tags` | 获取用户错题标签列表 |
| POST | `/api/v1/wrong/tags` | 创建错题标签（如：粗心、不会、混淆） |
| PUT | `/api/v1/wrong/tags/:id` | 更新错题标签 |
| DELETE | `/api/v1/wrong/tags/:id` | 删除错题标签 |
| POST | `/api/v1/wrong/:id/tags` | 给错题打标签 |
| DELETE | `/api/v1/wrong/:id/tags/:tagId` | 移除错题标签 |

### 错题本附件

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/wrong/:id/images` | 上传错题图片（如手写解题过程） |
| DELETE | `/api/v1/wrong/:id/images/:imageId` | 删除错题图片 |
| POST | `/api/v1/wrong/:id/voice` | 上传错题语音备注 |
| DELETE | `/api/v1/wrong/:id/voice/:voiceId` | 删除错题语音备注 |

### 闯关模式

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/practice/challenge/levels` | 获取关卡列表（按分类） |
| GET | `/api/v1/practice/challenge/levels/:id` | 获取关卡详情 |
| POST | `/api/v1/practice/challenge/:id/start` | 开始闯关 |
| POST | `/api/v1/practice/challenge/:id/answer` | 闯关答题 |
| POST | `/api/v1/practice/challenge/:id/submit` | 提交闯关结果 |
| GET | `/api/v1/practice/challenge/progress` | 获取用户闯关进度 |
| GET | `/api/v1/practice/challenge/progress/:categoryId` | 获取指定分类闯关进度 |

### 题目互动

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/questions/:id/notes` | 获取题目笔记（当前用户） |
| POST | `/api/v1/questions/:id/notes` | 创建/更新题目笔记 |
| DELETE | `/api/v1/questions/:id/notes` | 删除题目笔记 |
| GET | `/api/v1/practice/notes` | 获取用户所有笔记列表（分页） |
| POST | `/api/v1/questions/:id/difficulty-rating` | 提交题目难度评价（1-5） |
| GET | `/api/v1/questions/:id/difficulty-rating` | 获取题目难度评价统计（平均分、各分值人数） |
| GET | `/api/v1/questions/:id/difficulty-rating/my` | 获取当前用户的难度评价 |
| POST | `/api/v1/questions/:id/quality-rating` | 提交题目质量评分（1-5星） |
| GET | `/api/v1/questions/:id/quality-rating` | 获取题目质量评分统计 |
| GET | `/api/v1/questions/:id/quality-rating/my` | 获取当前用户的质量评分 |
| POST | `/api/v1/questions/:id/corrections` | 提交题目纠错 |
| GET | `/api/v1/questions/:id/corrections` | 获取题目纠错列表（当前用户） |
| GET | `/api/v1/practice/corrections/mine` | 获取我的纠错记录 |
| GET | `/api/v1/questions/:id/related` | 获取相关题目推荐（基于知识点关联） |
| POST | `/api/v1/questions/:id/share` | 生成题目分享链接 |
| GET | `/api/v1/questions/share/:code` | 通过分享码获取题目详情（无需登录） |

### 题目管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/questions/:id/copy` | 复制题目创建副本 |
| GET | `/api/v1/admin/questions/:id/preview` | 题目预览（支持PC端和移动端渲染） |
| POST | `/api/v1/admin/questions/check-duplicate` | 题目查重检测（基于内容MD5） |
| POST | `/api/v1/admin/questions/import-word` | 从Word文档导入题目 |
| GET | `/api/v1/admin/questions/export-word` | 导出题目为Word文档 |
| POST | `/api/v1/admin/questions/batch-import-attachments` | 批量导入题目附件（图片/音频/视频） |
| GET | `/api/v1/admin/questions/import/template` | 下载Excel导入模板 |
| GET | `/api/v1/admin/questions/import/progress/:task_id` | 查询导入进度 |
| POST | `/api/v1/admin/questions/import/retry/:task_id` | 导入失败重试 |
| GET | `/api/v1/admin/questions/import/report/:task_id` | 获取导入报告 |
| GET | `/api/v1/admin/questions/statistics` | 题目统计概览（总数、各状态数量、今日创建） |
| GET | `/api/v1/admin/questions/statistics/by-type` | 按题型统计题目数量 |
| GET | `/api/v1/admin/questions/statistics/by-category` | 按分类统计题目数量 |
| GET | `/api/v1/admin/questions/statistics/by-difficulty` | 按难度统计题目数量 |
| GET | `/api/v1/admin/questions/statistics/trend` | 题目创建趋势图表数据 |
| GET | `/api/v1/question/:id/favorite/status` | 获取题目收藏状态 |
| GET | `/api/v1/question/:id/preview` | 用户端题目预览 |
| POST | `/api/v1/favorite-folders/:id/move` | 移动收藏到指定收藏夹 |


---

## 五、管理端接口 (需管理员权限)

> 以下所有接口需要 JWT Token + 管理员角色，路径前缀为 `/api/v1/admin`

### 5.1 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/users` | 用户列表 |
| GET | `/api/v1/admin/users/:id` | 用户详情 |
| POST | `/api/v1/admin/users` | 创建用户 |
| PUT | `/api/v1/admin/users/:id` | 更新用户 |
| DELETE | `/api/v1/admin/users/:id` | 删除用户 |
| PUT | `/api/v1/admin/users/:id/status` | 更新用户状态 |
| POST | `/api/v1/admin/users/:id/reset-password` | 重置用户密码 |
| PUT | `/api/v1/admin/users/:id/roles` | 分配用户角色 |
| POST | `/api/v1/admin/users/batch-status` | 批量更新状态 |
| POST | `/api/v1/admin/users/batch-delete` | 批量删除用户 |
| POST | `/api/v1/admin/users/batch-roles` | 批量分配角色 |
| GET | `/api/v1/admin/users/export` | 导出用户 |

### 5.2 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/roles` | 角色列表 |
| GET | `/api/v1/admin/roles/:id` | 角色详情 |
| POST | `/api/v1/admin/roles` | 创建角色 |
| PUT | `/api/v1/admin/roles/:id` | 更新角色 |
| DELETE | `/api/v1/admin/roles/:id` | 删除角色 |
| GET | `/api/v1/admin/roles/:id/menus` | 获取角色菜单 |
| PUT | `/api/v1/admin/roles/:id/menus` | 分配角色菜单 |
| GET | `/api/v1/admin/roles/:id/permissions` | 获取角色权限 |
| PUT | `/api/v1/admin/roles/:id/permissions` | 分配角色权限 |

### 5.3 实名认证审核

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/realnames` | 实名认证列表 |
| PUT | `/api/v1/admin/realnames/:id/review` | 审核实名认证 |

### 5.4 标签管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/tags` | 标签列表 |
| GET | `/api/v1/admin/tags/:id` | 标签详情 |
| POST | `/api/v1/admin/tags` | 创建标签 |
| PUT | `/api/v1/admin/tags/:id` | 更新标签 |
| DELETE | `/api/v1/admin/tags/:id` | 删除标签 |

### 5.5 申请审核

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/applications` | 申请列表 |
| GET | `/api/v1/admin/applications/:id` | 申请详情 |
| PUT | `/api/v1/admin/applications/:id/review` | 审核申请 |
| POST | `/api/v1/admin/applications/:id/approve` | 审批通过（创作者/教师申请） |
| POST | `/api/v1/admin/applications/:id/reject` | 审批拒绝（创作者/教师申请） |

### 5.6 操作日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/logs/operation` | 操作日志列表 |
| GET | `/api/v1/admin/logs/operation/:id` | 日志详情 |
| GET | `/api/v1/admin/logs/operation/export` | 导出操作日志 |
| POST | `/api/v1/admin/logs/operation/clean` | 清空操作日志 |
| GET | `/api/v1/admin/logs/export` | 导出日志（支持指定日志类型） |

### 5.6b 练习管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/practice/statistics` | 全局练习统计概览 |
| GET | `/api/v1/admin/practice/users` | 用户练习统计列表（支持分页、排序） |
| GET | `/api/v1/admin/practice/users/:id` | 指定用户练习详情 |
| GET | `/api/v1/admin/practice/hot-questions` | 热门题目排行（练习次数最多） |
| GET | `/api/v1/admin/practice/accuracy-analysis` | 正确率分析（按分类/难度/题型） |
| GET | `/api/v1/admin/practice/difficulty-distribution` | 难度分布统计 |

### 5.7 登录日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/logs/login` | 登录日志列表 |
| GET | `/api/v1/admin/logs/login/export` | 导出登录日志 |
| POST | `/api/v1/admin/logs/login/clean` | 清空登录日志 |

### 5.8 菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/menus` | 菜单列表 |
| GET | `/api/v1/admin/menus/tree` | 菜单树形结构 |
| GET | `/api/v1/admin/menus/:id` | 菜单详情 |
| POST | `/api/v1/admin/menus` | 创建菜单 |
| PUT | `/api/v1/admin/menus/:id` | 更新菜单 |
| DELETE | `/api/v1/admin/menus/:id` | 删除菜单 |
| PUT | `/api/v1/admin/menus/sort` | 菜单排序（拖拽排序） |

### 5.9 题目管理

#### 基本 CRUD

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/questions` | 题目列表 (管理视图) |
| GET | `/api/v1/admin/questions/:id` | 题目详情 |
| POST | `/api/v1/admin/questions` | 创建题目 |
| PUT | `/api/v1/admin/questions/:id` | 更新题目 |
| DELETE | `/api/v1/admin/questions/:id` | 删除题目 |
| PUT | `/api/v1/admin/questions/:id/status` | 更新题目状态 |

#### 版本管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/questions/:id/versions` | 版本历史 |
| GET | `/api/v1/admin/questions/:id/versions/:versionId` | 版本详情 |
| POST | `/api/v1/admin/questions/:id/versions/:version/rollback` | 回滚版本 |

#### 分享

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/questions/:id/share` | 创建分享链接 |
| GET | `/api/v1/admin/questions/:id/shares` | 我的分享列表 |
| GET | `/api/v1/admin/shares` | 所有分享列表 |
| DELETE | `/api/v1/admin/shares/:id` | 撤销分享 |

#### 导入导出

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/questions/import` | 导入题目 (Excel) |
| GET | `/api/v1/admin/questions/export` | 导出题目 (Excel) |

#### 批量操作

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/questions/batch/publish` | 批量发布 |
| POST | `/api/v1/admin/questions/batch/archive` | 批量归档 |
| POST | `/api/v1/admin/questions/batch/delete` | 批量删除 |
| POST | `/api/v1/admin/questions/batch/move` | 批量移动分类 |

#### 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/questions/stats` | 题目统计 |

### 5.10 试卷管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/papers` | 试卷列表 |
| GET | `/api/v1/admin/papers/:id` | 试卷详情 |
| POST | `/api/v1/admin/papers` | 创建试卷 |
| PUT | `/api/v1/admin/papers/:id` | 更新试卷 |
| DELETE | `/api/v1/admin/papers/:id` | 删除试卷 |
| GET | `/api/v1/admin/papers/:id/preview` | 预览试卷 |
| POST | `/api/v1/admin/papers/:id/copy` | 复制试卷 |
| PUT | `/api/v1/admin/papers/:id/publish` | 发布试卷 |
| POST | `/api/v1/admin/papers/:id/save-template` | 保存为模板 |
| GET | `/api/v1/admin/papers/:id/export` | 导出试卷 |
| GET | `/api/v1/admin/papers/:id/stats` | 试卷统计 |

### 5.11 模板管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/templates` | 模板列表 |
| POST | `/api/v1/admin/templates/create` | 从模板创建试卷 |

### 5.12 考试管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/exams` | 考试列表 |
| GET | `/api/v1/admin/exams/:id` | 考试详情 |
| POST | `/api/v1/admin/exams` | 创建考试 |
| PUT | `/api/v1/admin/exams/:id` | 更新考试 |
| DELETE | `/api/v1/admin/exams/:id` | 删除考试 |
| PUT | `/api/v1/admin/exams/:id/publish` | 发布考试 |
| PUT | `/api/v1/admin/exams/:id/close` | 结束考试 |
| GET | `/api/v1/admin/exams/:id/monitor` | 考试监控 |
| POST | `/api/v1/admin/exams/:id/review` | 阅卷 |
| GET | `/api/v1/admin/exams/:id/analysis` | 考试分析 |
| GET | `/api/v1/admin/exams/:id/review-progress` | 阅卷进度 |
| GET | `/api/v1/admin/exams/:id/warnings` | 异常行为列表 |
| GET | `/api/v1/admin/exams/:id/extensions` | 延期记录 |
| GET | `/api/v1/admin/exams/:id/pauses` | 暂停记录 |

### 5.13 成绩管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/scores` | 成绩列表 |
| GET | `/api/v1/admin/scores/:id` | 成绩详情 |
| GET | `/api/v1/admin/scores/analysis` | 成绩分析 |
| GET | `/api/v1/admin/scores/:id/export` | 导出成绩 |

### 5.14 班级管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/class` | 班级列表 |
| GET | `/api/v1/admin/class/:id` | 班级详情 |
| POST | `/api/v1/admin/class` | 创建班级 |
| PUT | `/api/v1/admin/class/:id` | 更新班级 |
| DELETE | `/api/v1/admin/class/:id` | 删除班级 |
| GET | `/api/v1/admin/class/:id/members` | 成员列表 |
| DELETE | `/api/v1/admin/class/:id/members/:uid` | 移除成员 |

### 5.15 系统设置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/settings` | 获取系统设置 |
| PUT | `/api/v1/admin/settings` | 更新系统设置 |


### 内容审核

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/reviews` | 获取审核列表（支持状态筛选、分页） |
| GET | `/api/v1/admin/reviews/:id` | 获取审核详情 |
| POST | `/api/v1/admin/reviews/:id/approve` | 审核通过（支持填写审核意见） |
| POST | `/api/v1/admin/reviews/:id/reject` | 审核驳回（需填写驳回原因） |
| POST | `/api/v1/admin/reviews/:id/revision` | 标记需修改（提出修改意见） |
| POST | `/api/v1/admin/reviews/:id/reply` | 回复审核意见（创作者与审核人沟通） |
| POST | `/api/v1/admin/reviews/batch/approve` | 批量审核通过 |
| POST | `/api/v1/admin/reviews/batch/reject` | 批量审核驳回 |
| POST | `/api/v1/admin/reviews/batch-approve` | 批量审核通过（别名） |
| POST | `/api/v1/admin/reviews/batch-reject` | 批量审核拒绝（别名） |
| GET | `/api/v1/admin/reviews/statistics` | 获取审核统计（通过率、驳回率、平均审核时长） |
| GET | `/api/v1/admin/reviews/timeout` | 获取超时审核提醒列表（超过24小时未处理） |

### 创作者等级管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/creator-levels` | 获取创作者等级列表 |
| POST | `/api/v1/admin/creator-levels` | 创建创作者等级 |
| PUT | `/api/v1/admin/creator-levels/:id` | 更新创作者等级 |
| DELETE | `/api/v1/admin/creator-levels/:id` | 删除创作者等级 |

### 创作者协议管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/agreements` | 获取创作者协议列表 |
| POST | `/api/v1/admin/agreements` | 创建创作者协议 |
| PUT | `/api/v1/admin/agreements/:id` | 更新创作者协议 |

### 审核流程管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/review-workflows` | 获取审核流程列表 |
| POST | `/api/v1/admin/review-workflows` | 创建审核流程 |
| PUT | `/api/v1/admin/review-workflows/:id` | 更新审核流程 |
| DELETE | `/api/v1/admin/review-workflows/:id` | 删除审核流程 |

### 分类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| PUT | `/api/v1/admin/categories/sort` | 分类拖拽排序 |
| POST | `/api/v1/admin/categories/merge` | 分类合并（关联题目自动归属目标分类） |
| GET | `/api/v1/admin/categories/statistics` | 分类统计（各分类题目数量） |

### 审批流程

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/approvals` | 获取审批列表 |
| GET | `/api/v1/admin/approvals/:id` | 获取审批详情 |
| PUT | `/api/v1/admin/approvals/:id/approve` | 审批通过 |
| PUT | `/api/v1/admin/approvals/:id/reject` | 审批驳回 |
| GET | `/api/v1/admin/approvals/pending-count` | 获取待审批数量 |

### 敏感词管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/sensitive-words` | 获取敏感词列表 |
| POST | `/api/v1/admin/sensitive-words` | 添加敏感词 |
| PUT | `/api/v1/admin/sensitive-words/:id` | 编辑敏感词 |
| DELETE | `/api/v1/admin/sensitive-words/:id` | 删除敏感词 |
| POST | `/api/v1/admin/sensitive-words/import` | 批量导入敏感词 |
| POST | `/api/v1/admin/sensitive-words/reload` | 重新加载敏感词库到内存 |

### 权限管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/permissions/buttons` | 获取按钮权限列表 |
| POST | `/api/v1/admin/permissions/buttons` | 创建按钮权限 |
| PUT | `/api/v1/admin/permissions/buttons/:id` | 更新按钮权限 |
| DELETE | `/api/v1/admin/permissions/buttons/:id` | 删除按钮权限 |
| GET | `/api/v1/admin/roles/:id/permissions` | 获取角色权限（菜单权限、按钮权限） |
| PUT | `/api/v1/admin/roles/:id/permissions` | 分配角色权限 |

### 知识点管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/knowledge-points` | 获取知识点列表（支持按分类筛选和关键词搜索） |
| GET | `/api/v1/admin/knowledge-points/tree` | 获取知识点树形结构 |
| POST | `/api/v1/admin/knowledge-points` | 创建知识点（设置名称、所属分类、权重等） |
| PUT | `/api/v1/admin/knowledge-points/:id` | 更新知识点信息 |
| DELETE | `/api/v1/admin/knowledge-points/:id` | 删除知识点（有关联题目时禁止删除） |
| POST | `/api/v1/admin/knowledge-points/batch/delete` | 批量删除知识点 |
| POST | `/api/v1/admin/knowledge-points/batch/move` | 批量移动知识点到其他分类 |
| GET | `/api/v1/admin/knowledge-points/:id/statistics` | 获取知识点统计信息（使用频次、正确率、关联题目数） |

### 管理后台

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/dashboard/export` | 导出仪表盘数据（Excel/CSV） |
| GET | `/api/v1/admin/dashboard/export-pdf` | 导出仪表盘数据为PDF |

### 系统日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/logs/system` | 获取系统日志列表 |
| GET | `/api/v1/admin/logs/error` | 获取错误日志列表 |
| GET | `/api/v1/admin/logs/apply` | 获取班级创作者申请日志列表 |
| GET | `/api/v1/admin/logs/stats` | 获取日志统计 |
| GET | `/api/v1/admin/logs/search` | 日志全文搜索 |
| POST | `/api/v1/admin/logs/archive` | 手动触发日志归档 |
| GET | `/api/v1/admin/logs/archives` | 获取日志归档记录列表 |

### 系统监控

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/system/metrics` | 获取系统运行指标（CPU/内存/磁盘/网络） |
| GET | `/api/v1/admin/system/status` | 获取系统健康状态 |
| GET | `/api/v1/admin/system/alerts` | 获取系统告警列表 |
| PUT | `/api/v1/admin/system/alerts/:id/handle` | 处理系统告警 |

### 系统设置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/features` | 获取功能开关配置 |
| PUT | `/api/v1/admin/features/:key` | 更新功能开关 |
| GET | `/api/v1/admin/security` | 获取安全配置 |
| PUT | `/api/v1/admin/security` | 更新安全配置 |
| GET | `/api/v1/admin/storage` | 获取存储配置 |
| PUT | `/api/v1/admin/storage` | 更新存储配置 |
| GET | `/api/v1/admin/settings/class` | 获取班级相关设置 |
| PUT | `/api/v1/admin/settings/class` | 更新班级相关设置 |
| GET | `/api/v1/admin/settings/resource` | 获取资源相关设置 |
| PUT | `/api/v1/admin/settings/resource` | 更新资源相关设置 |
| GET | `/api/v1/admin/email/config` | 获取邮件服务器配置 |
| PUT | `/api/v1/admin/email/config` | 更新邮件服务器配置 |
| GET | `/api/v1/admin/email/templates` | 获取邮件模板列表 |
| PUT | `/api/v1/admin/email/templates/:id` | 更新邮件模板 |
| POST | `/api/v1/admin/email/test` | 测试邮件发送 |
| GET | `/api/v1/admin/sms/config` | 获取短信服务商配置 |
| PUT | `/api/v1/admin/sms/config` | 更新短信服务商配置 |
| GET | `/api/v1/admin/sms/templates` | 获取短信模板列表 |
| POST | `/api/v1/admin/sms/test` | 测试短信发送 |
| GET | `/api/v1/admin/notifications/channels` | 获取通知渠道列表（站内信/邮件/短信/微信） |
| PUT | `/api/v1/admin/notifications/channels/:channel` | 更新通知渠道配置 |
| GET | `/api/v1/admin/notifications/rate-limits` | 获取通知频率限制配置 |
| PUT | `/api/v1/admin/notifications/rate-limits/:id` | 更新通知频率限制配置 |
| POST | `/api/v1/admin/notifications/test` | 测试通知发送 |
| GET | `/api/v1/admin/cache/stats` | 获取缓存统计信息 |
| POST | `/api/v1/admin/cache/clear` | 清除缓存 |
| GET | `/api/v1/admin/i18n/languages` | 获取语言列表 |
| PUT | `/api/v1/admin/i18n/languages/:id` | 更新语言配置 |
| GET | `/api/v1/admin/theme` | 获取主题配置 |
| PUT | `/api/v1/admin/theme` | 更新主题配置 |
| GET | `/api/v1/admin/alerts/rules` | 获取系统告警规则列表 |
| POST | `/api/v1/admin/alerts/rules` | 创建系统告警规则 |
| PUT | `/api/v1/admin/alerts/rules/:id` | 更新系统告警规则 |
| DELETE | `/api/v1/admin/alerts/rules/:id` | 删除系统告警规则 |
| GET | `/api/v1/admin/alerts/records` | 获取系统告警记录 |
| GET | `/api/v1/admin/backup/config` | 获取备份配置 |
| PUT | `/api/v1/admin/backup/config` | 更新备份配置 |
| POST | `/api/v1/admin/backup` | 创建数据库备份 |
| GET | `/api/v1/admin/backup/list` | 获取备份列表 |
| POST | `/api/v1/admin/backup/:id/restore` | 恢复数据库备份 |
| DELETE | `/api/v1/admin/backup/:id` | 删除备份 |
| GET | `/api/v1/admin/backup/status` | 获取备份监控状态 |
| POST | `/api/v1/admin/backup/drill` | 手动触发备份恢复演练 |
| GET | `/api/v1/admin/backup/drill/records` | 获取恢复演练记录 |

### 评论管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/comments` | 获取全部评论列表（支持状态筛选） |
| PUT | `/api/v1/admin/comments/:id/status` | 审核评论（通过/隐藏） |
| DELETE | `/api/v1/admin/comments/:id` | 管理员删除评论 |
| POST | `/api/v1/admin/comments/batch-delete` | 批量删除评论 |
| GET | `/api/v1/admin/comments/reports` | 获取评论举报列表 |
| PUT | `/api/v1/admin/comments/reports/:id` | 处理评论举报 |
| GET | `/api/v1/admin/comments/blacklist` | 获取评论黑名单列表 |
| POST | `/api/v1/admin/comments/blacklist` | 添加评论黑名单（禁止特定用户评论） |
| DELETE | `/api/v1/admin/comments/blacklist/:id` | 移除评论黑名单 |
| GET | `/api/v1/admin/comments/stats` | 获取评论统计 |
| GET | `/api/v1/admin/comments/export` | 导出评论为Excel |
| GET | `/api/v1/admin/comments/audit-rules` | 获取评论审核规则列表 |
| POST | `/api/v1/admin/comments/audit-rules` | 创建评论审核规则 |
| PUT | `/api/v1/admin/comments/audit-rules/:id` | 更新评论审核规则 |
| DELETE | `/api/v1/admin/comments/audit-rules/:id` | 删除评论审核规则 |
| PUT | `/api/v1/admin/comments/batch-audit` | 批量审核评论 |

### 试卷管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/papers/:id` | 获取试卷详情 |
| POST | `/api/v1/admin/papers/:id/share` | 试卷共享（指定教师/教研组） |
| POST | `/api/v1/admin/papers/import` | 从Word/PDF导入试卷 |
| POST | `/api/v1/admin/papers/:id/favorite` | 收藏/取消收藏试卷 |
| GET | `/api/v1/admin/templates` | 获取试卷模板列表 |
| POST | `/api/v1/admin/papers/from-template` | 从模板创建试卷 |

### 通知管理

| 方法 | 路径 | 说明 |
|------|------|------|
| PUT | `/api/v1/notification/:id/recall` | 撤回已发送通知 |
| POST | `/api/v1/notification/batch-send` | 群发通知给指定用户 |
| GET | `/api/v1/notification/stats` | 获取通知统计 |
| POST | `/api/v1/notification/scheduled` | 创建定时通知 |
| GET | `/api/v1/notification/scheduled` | 获取定时通知列表 |
| DELETE | `/api/v1/notification/scheduled/:id` | 取消定时通知 |
| GET | `/api/v1/admin/notification/templates` | 获取通知模板列表 |
| POST | `/api/v1/admin/notification/templates` | 创建通知模板 |
| PUT | `/api/v1/admin/notification/templates/:id` | 更新通知模板 |
| DELETE | `/api/v1/admin/notification/templates/:id` | 删除通知模板 |

### 题目纠错管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/corrections` | 管理员获取所有纠错列表（支持状态筛选） |
| PUT | `/api/v1/admin/corrections/:id` | 管理员审核纠错（采纳/驳回） |
| PUT | `/api/v1/admin/corrections/:id/status` | 管理员更新纠错状态（已修复） |

### 评语模板管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/comment-templates` | 评语模板列表 |
| POST | `/api/v1/admin/comment-templates` | 创建评语模板 |
| PUT | `/api/v1/admin/comment-templates/:id` | 更新评语模板 |
| DELETE | `/api/v1/admin/comment-templates/:id` | 删除评语模板 |

### 公告管理

| 方法 | 路径 | 说明 |
|------|------|------|
| PUT | `/api/v1/admin/notices/:id` | 编辑公告 |
| DELETE | `/api/v1/admin/notices/:id` | 撤回公告 |


---

## 六、通用约定

### 6.1 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 / Token 过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突 (如用户名已存在) |
| 422 | 请求格式正确但语义错误 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

### 6.2 题目类型枚举

| 值 | 类型 |
|----|------|
| 1 | 单选题 |
| 2 | 多选题 |
| 3 | 判断题 |
| 4 | 填空题 |
| 5 | 简答题 |

### 6.3 难度枚举

| 值 | 难度 |
|----|------|
| 1 | 简单 |
| 2 | 中等 |
| 3 | 困难 |

### 6.4 题目状态枚举

| 值 | 状态 |
|----|------|
| 0 | 草稿 |
| 1 | 已发布 |
| 2 | 已归档 |
| 3 | 待审核 |
| 4 | 需修改 |
| 5 | 审核超时 |

### 6.5 可见性枚举

| 值 | 可见性 |
|----|--------|
| 1 | 公开 |
| 2 | 私有 |
| 3 | 班级可见 |

### 6.6 排序参数

列表接口支持 `sort_by` 和 `sort_order` 参数：

```
GET /api/v1/questions?sort_by=created_at&sort_order=desc
```

### 6.7 搜索参数

部分列表接口支持关键词搜索：

```
GET /api/v1/questions?keyword=数学
```
