# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在本仓库中工作时提供指引。

## 项目概览

题小助 (QuestionHelper) 是一款面向教育领域的题库管理与在线考试平台，由三个子项目组成：Go 后端 API、Vue 3 管理后台、UniApp 移动端。`docs/` 目录下为中文文档。

## 仓库结构

```
questionhelper-server/   # Go 后端 (Gin + GORM + MySQL + Redis)
questionhelper-admin/    # Vue 3 管理后台 (Element Plus + TypeScript)
questionhelper-app/      # UniApp 移动端 (Vue 3 + TypeScript)
docs/                    # 需求与技术设计文档
```

## 常用命令

### 后端 (questionhelper-server/)

```bash
cd questionhelper-server
make run          # 启动开发服务器 (go run cmd/server/main.go)
make build        # 构建二进制到 bin/server
make test         # 运行所有测试 (go test ./... -v)
make fmt          # 格式化代码
make lint         # 代码检查 (golangci-lint)
make swagger      # 重新生成 Swagger 文档
make migrate      # 运行数据库迁移
```

运行单个测试：`go test ./internal/service/question/ -run TestCreateQuestion -v`

API 集成测试：`bash test_api.sh`（需要后端服务运行在 localhost:8080）。

初始化数据：`go run cmd/seed/main.go`（无 Makefile 目标）。

### 管理后台 (questionhelper-admin/)

> **注意**：管理后台目录目前为空占位，尚未提交任何源代码。以下命令均为规划中的命令；基础模板为 vue3-element-admin。

包管理器：**pnpm**（强制使用，不要使用 npm/yarn）。Node 版本要求：`^20.19.0 || >=22.12.0`。

```bash
cd questionhelper-admin
pnpm dev              # Vite 开发服务器
pnpm build            # 类型检查 + 生产构建
pnpm build-only       # 生产构建（跳过类型检查）
pnpm type-check       # vue-tsc --noEmit
pnpm lint             # ESLint + Prettier + Stylelint
```

### 移动端 (questionhelper-app/)

包管理器：**npm**。

```bash
cd questionhelper-app
npm run dev:h5          # H5 开发服务器
npm run dev:mp-weixin   # 微信小程序开发服务器
npm run build:h5        # H5 生产构建
npm run build:mp-weixin # 微信小程序生产构建
```

### Docker（全栈）

```bash
docker-compose up -d    # 启动 MySQL 8.0、Redis 7、Go 服务、Nginx 管理后台
```

## 后端架构

`questionhelper-server/internal/` 下的分层架构：

```
router/       → 路由定义（Gin 路由组，按领域分文件）
controller/   → HTTP 处理器（每个领域一个包）
service/      → 业务逻辑层
repository/   → 数据访问层（GORM 查询）
dto/          → 请求/响应结构体（含验证标签）
model/        → GORM 模型定义（数据库实体）
middleware/   → 认证、CORS、限流、RBAC、操作日志、敏感词过滤
ws/           → WebSocket 实时通信（hub/client/message）
job/          → 后台任务（考试提交、导入导出、通知）
task/         → 定时任务（清理、提醒、过期处理）
```

共享包位于 `questionhelper-server/pkg/`：cache、captcha、config、consts、database、encrypt、errors、excel、jwt、logger、mq、response、sensitive、upload、validator。

**入口**：`cmd/server/main.go` —— 加载配置，初始化 MySQL/Redis/日志/JWT，启动 Gin 服务器。配置路径默认为 `config/config.yaml`，可通过 `CONFIG_PATH` 环境变量覆盖。全局单例：`database.DB`（`*gorm.DB`）、`config.Cfg`、JWT 密钥。

**依赖组装**：无依赖注入框架。`router.go` 实例化所有控制器（通过 `NewXxxController()` 返回空结构体）。Service 和 Repository 层使用**包级函数**（非结构体方法），通过导入包直接互相调用。Repository 通过全局 `database.DB` 访问数据库。

**路由层级**（均在 `/api/v1` 下）：
- **公开**（`/api/v1/auth/...`）—— 登录、注册、验证码、刷新令牌
- **已认证**（`/api/v1/...`）—— 普通用户接口（需要 JWT）
- **管理员**（`/api/v1/admin/...`）—— 管理员接口（JWT + admin 角色）

**认证**：JWT access + refresh token 模式；通过 Redis 实现令牌黑名单。

**配置**：`config/config.yaml` YAML 文件，通过 `pkg/config` 加载。

## 管理后台架构

基于 Vue 3 + Element Plus 的功能型 SPA。`src/` 下关键目录：

- `api/` —— 按领域划分的 API 模块（Axios 实例）
- `views/` —— 按领域划分的页面组件
- `stores/` —— Pinia 状态管理（app、user、permission、settings、dict、tags-view、tenant）
- `router/` —— 静态 + 动态路由；**路由根据用户权限从后端获取**并动态注入
- `components/` —— 共享组件（含 CURD 框架）
- `composables/` —— Vue 组合式函数（SSE、表格选择、最近菜单）
- `layouts/` —— 布局变体（Left、Top、Mix、Base）

**样式规范**：SCSS + UnoCSS；Prettier 配置为 100 字符宽度、2 空格缩进、双引号、分号、尾逗号。ESLint 要求模板中组件名使用 PascalCase，Vue 块顺序为 template → script → style。

**提交规范**：通过 commitlint + Commitizen (cz-git) 强制执行 Conventional Commits。使用 `pnpm commit` 进行交互式提交。

## 移动端架构

基于 UniApp (Vue 3 + TypeScript)，目标平台为 H5、微信小程序和原生 iOS/Android。页面路由定义在 `pages.json` 中（42 个页面，分布在 13 个功能目录下）。底部导航栏：首页、题库、练习、班级、我的。使用 `#ifdef` 条件编译处理平台差异代码。

关键目录：`pages/`、`components/`、`api/`、`store/`（Pinia，Composition API setup 风格）、`hooks/`、`utils/`。

API 层（`api/request.ts`）封装 `uni.request`，自动注入 Bearer token，基础 URL 来自 `VITE_API_BASE_URL`（默认 `/api/v1`），成功响应码为 `'00000'`。

## 关键技术细节

- **数据库**：MySQL 8.0 + GORM 自动迁移。表名使用复数蛇形命名，必含 `id`/`created_at`/`updated_at`/`deleted_at` 字段，布尔字段以 `is_` 前缀，全局使用软删除。
- **缓存**：Redis 7，用于会话/令牌黑名单及通用缓存
- **文件存储**：本地文件系统或 MinIO（S3 兼容），可配置
- **实时通信**：WebSocket 用于通知；管理后台使用 SSE（Server-Sent Events）
- **Excel**：后端使用 excelize 进行导入导出；管理后台使用 ExcelJS
- **富文本**：管理后台使用 @wangeditor-next/editor；移动端使用自定义 RichText 组件
- **RBAC**：基于角色的权限控制，支持菜单/权限分配；中间件按路由检查权限
- **敏感词过滤**：内置内容审核中间件
- **后台处理**：内存消息队列（`pkg/mq`）配合 job worker 和定时任务（`internal/task/` 含 9 个定时任务，`internal/job/` 含 6 个异步任务类型）

## 开发规范

### API 响应格式

所有端点通过 `pkg/response/` 返回统一 JSON 信封：

```json
{"code": "00000", "msg": "success", "data": {...}}
```

成功码为 `"00000"`。分页响应在 `data` 中嵌套 `{"list", "total", "page", "page_size"}`。移动端的 `request.ts` 还会检查码 `"A0003"`/`"A0004"`（令牌过期 → 跳转登录）。

### 错误处理

`pkg/errors/` 定义了 `AppError`，包含 `Code`（如 `"A0001"`、`"A201"`）和 `Message`。预定义了常见场景的哨兵错误（参数错误、未授权、禁止访问、未找到、内部错误、用户错误、令牌错误、考试错误）。控制器通常直接使用 `response.Error(c, httpCode, message)`。

### DTO 模式

请求/响应结构体位于 `internal/dto/`。Service 层接收 DTO 指针并返回 DTO info 结构体。Repository 函数操作 `internal/model/` 结构体；转换在 Service 层通过 `toXxxInfo()` 辅助函数完成。

### 测试

后端目前**没有测试文件**。添加测试时使用 `testing` + `testify`。Repository 和 Service 层未定义接口，因此 mock 需要引入接口或使用真实数据库测试。运行单个测试：`go test ./path/to/pkg/ -run TestFunctionName -v`

## Git 规范

分支模型：`main`（生产）、`develop`、`feature/*`、`hotfix/*`、`release/*`。通过 commitlint 强制执行 Conventional Commits。提交类型：`feat`、`fix`、`docs`、`style`、`refactor`、`perf`、`test`、`build`、`ci`、`chore`、`revert`、`wip`。scope 对应模块（auth、user、question、exam、class 等）。

## 基础设施

本地开发默认需要 MySQL 和 Redis。服务端配置位于 `questionhelper-server/config/config.yaml`。默认管理员账号：`admin / admin123`。
