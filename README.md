# 题小助 (QuestionHelper)

一款面向教育领域的题库管理与在线考试平台，支持多角色权限、班级管理、刷题练习、错题本、评论社区等完整功能。

## 项目概览

| 子项目 | 技术栈 | 说明 |
|--------|--------|------|
| `questionhelper-server/` | Go + Gin + GORM + MySQL + Redis | 后端 API 服务 |
| `questionhelper-admin/` | Vue 3 + Element Plus + TypeScript | 管理后台前端 |
| `questionhelper-app/` | UniApp (Vue 3 + TypeScript) | 移动端 H5 / 微信小程序 / 原生 App |

## 技术栈

### 后端

- **语言**: Go 1.21+
- **Web 框架**: Gin v1.9
- **ORM**: GORM v1.25 (MySQL 8.0)
- **缓存**: Redis 7 (go-redis v9)
- **认证**: JWT (golang-jwt v5) + OAuth 2.0
- **实时通信**: WebSocket (gorilla/websocket)
- **文件存储**: 本地文件系统 / MinIO (S3 兼容)
- **Excel**: excelize v2
- **日志**: Zap (结构化日志)
- **验证码**: base64Captcha

### 前端 (管理后台)

- **框架**: Vue 3.5 + TypeScript 5.9
- **构建**: Vite 8 (Rolldown)
- **UI 库**: Element Plus 2.13
- **状态管理**: Pinia 3
- **图表**: ECharts 6
- **富文本**: @wangeditor-next/editor
- **样式**: SCSS + UnoCSS
- **包管理**: pnpm (强制)

### 移动端

- **框架**: UniApp (DCloud) + Vue 3.4 + TypeScript
- **状态管理**: Pinia 2.1
- **目标平台**: H5、微信小程序、iOS / Android 原生

## 快速启动

### 前置要求

| 工具 | 版本 |
|------|------|
| Go | 1.21+ |
| Node.js | ^20.19.0 或 >=22.12.0 |
| pnpm | 最新版 (管理后台专用) |
| MySQL | 8.0 |
| Redis | 7+ |
| Docker & Docker Compose | 可选，用于容器化部署 |

### 方式一：Docker Compose (推荐)

```bash
# 克隆项目
git clone <repo-url> questionhelper-pro
cd questionhelper-pro

# 一键启动所有服务 (MySQL + Redis + 后端 + 前端)
docker-compose up -d

# 查看日志
docker-compose logs -f server
```

启动后访问：
- 管理后台: http://localhost:3000
- API 服务: http://localhost:8080
- 默认管理员: `admin / admin123`

### 方式二：本地开发

#### 1. 启动 MySQL 和 Redis

```bash
# 确保 MySQL 8.0 和 Redis 7 已运行
# 创建数据库
mysql -u root -p -e "CREATE DATABASE questionhelper DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
```

#### 2. 启动后端

```bash
cd questionhelper-server

# 编辑配置文件（数据库连接、Redis 地址等）
vim config/config.yaml

# 运行数据库迁移
make migrate

# 启动开发服务器
make run
```

后端将在 http://localhost:8080 启动。

#### 3. 启动管理后台

```bash
cd questionhelper-admin

# 安装依赖（必须使用 pnpm）
pnpm install

# 启动开发服务器
pnpm dev
```

#### 4. 启动移动端

```bash
cd questionhelper-app

# 安装依赖
npm install

# H5 开发模式
npm run dev:h5

# 微信小程序开发模式
npm run dev:mp-weixin
```

## 常用命令

### 后端

```bash
make run          # 启动开发服务器
make build        # 构建二进制到 bin/server
make test         # 运行所有测试
make fmt          # 格式化代码
make lint         # 代码检查 (golangci-lint)
make swagger      # 重新生成 Swagger 文档
make migrate      # 运行数据库迁移
```

### 管理后台

```bash
pnpm dev              # Vite 开发服务器
pnpm build            # 类型检查 + 生产构建
pnpm build-only       # 生产构建 (跳过类型检查)
pnpm type-check       # 仅类型检查
pnpm lint             # ESLint + Prettier + Stylelint
pnpm commit           # 交互式提交 (Commitizen)
```

### 移动端

```bash
npm run dev:h5            # H5 开发服务器
npm run dev:mp-weixin     # 微信小程序开发服务器
npm run build:h5          # H5 生产构建
npm run build:mp-weixin   # 微信小程序生产构建
```

### API 测试

```bash
# 需要后端服务运行在 localhost:8080
bash test_api.sh
```

## 项目结构

```
questionhelper-pro/
├── questionhelper-server/         # Go 后端
│   ├── cmd/server/                # 入口 (main.go)
│   ├── config/                    # 配置文件 (config.yaml)
│   ├── internal/
│   │   ├── router/                # 路由定义
│   │   ├── controller/            # HTTP 处理器
│   │   ├── service/               # 业务逻辑层
│   │   ├── repository/            # 数据访问层
│   │   ├── dto/                   # 请求/响应结构体
│   │   ├── model/                 # GORM 模型 (数据库实体)
│   │   ├── middleware/            # 中间件 (认证、RBAC、日志等)
│   │   ├── ws/                    # WebSocket 实时通信
│   │   ├── job/                   # 后台任务
│   │   └── task/                  # 定时任务
│   ├── pkg/                       # 共享包 (cache, jwt, logger 等)
│   └── Makefile
├── questionhelper-admin/          # Vue 3 管理后台
│   ├── src/
│   │   ├── api/                   # API 模块
│   │   ├── views/                 # 页面组件
│   │   ├── stores/                # Pinia 状态管理
│   │   ├── router/                # 路由 (动态权限路由)
│   │   ├── components/            # 共享组件 (含 CURD 框架)
│   │   ├── composables/           # Vue 组合式函数
│   │   └── layouts/               # 布局组件
│   └── package.json
├── questionhelper-app/            # UniApp 移动端
│   ├── pages/                     # 页面 (40+ 页面)
│   ├── components/                # 组件
│   ├── api/                       # API 模块
│   ├── store/                     # Pinia 状态管理
│   └── package.json
├── docs/                          # 项目文档
│   ├── 技术设计文档/               # 14 个模块的详细设计
│   ├── 品牌设计/                   # Favicon、品牌风格、图标说明
│   ├── 需求与规划/                 # 功能需求、技术栈、目录结构
│   └── ...
├── docker-compose.yml             # Docker 编排
├── DEPLOYMENT.md                  # 部署指南
├── CLAUDE.md                      # AI 辅助开发指南
└── README.md                      # 本文件
```

## 功能模块

| 模块 | 说明 |
|------|------|
| 用户认证 | 注册、登录、JWT、OAuth、验证码、手机号/邮箱绑定 |
| 用户管理 | 个人资料、隐私设置、实名认证、设备管理、安全日志 |
| 题库管理 | 题目 CRUD、多题型支持、分类/知识点、版本管理、审核、分享 |
| 考试管理 | 组卷、考试发布、在线答题、自动/手动阅卷、防作弊、成绩分析 |
| 班级管理 | 创建/加入班级、成员管理、作业、通知公告 |
| 刷题练习 | 专项练习、随机练习、练习统计 |
| 错题本 | 自动收集错题、复习、掌握标记、错题分析 |
| 评论社区 | 评论、回复、点赞、举报 |
| 数据统计 | 仪表盘、用户统计、题目统计、考试统计、排行榜 |
| 系统设置 | 角色权限 (RBAC)、菜单管理、操作日志、敏感词管理 |

## 默认账号

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin123 |

## 文档

- [功能需求文档](docs/需求与规划/题小助V2.0-功能需求文档.md)
- [技术栈文档](docs/需求与规划/题小助V2.0-技术栈文档.md)
- [目录结构文档](docs/需求与规划/题小助V2.0-目录结构文档.md)
- [技术设计文档](docs/技术设计文档/) (14 个模块)
- [品牌设计](docs/品牌设计/) (Favicon、品牌风格、图标说明)
- [部署指南](DEPLOYMENT.md)
- [API 接口文档](docs/API接口文档.md)
- [数据库设计文档](docs/数据库设计文档.md)
- [环境配置说明](docs/环境配置说明.md)
- [开发规范文档](docs/开发规范文档.md)

## 许可证

本项目仅供学习和内部使用。
