# 题小助 V2.0 后端服务

基于 Go + Gin + GORM 的后端服务。

## 技术栈

- **语言**: Go 1.21
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.x
- **消息队列**: Redis Stream
- **WebSocket**: Gorilla WebSocket
- **对象存储**: MinIO / 阿里云 OSS
- **日志**: Zap
- **JWT**: golang-jwt

## 项目结构

```
questionhelper-server/
├── cmd/server/          # 程序入口
├── internal/            # 内部代码
│   ├── controller/      # 控制器层
│   ├── service/         # 业务逻辑层
│   ├── repository/      # 数据访问层
│   ├── model/           # 数据模型
│   ├── middleware/       # 中间件
│   ├── router/          # 路由定义
│   ├── dto/             # 数据传输对象
│   ├── task/            # 定时任务
│   ├── job/             # 异步任务
│   └── ws/              # WebSocket
├── pkg/                 # 公共工具包
├── migrations/          # 数据库迁移
├── scripts/             # 脚本
├── docs/                # 文档
├── config/              # 配置文件
├── Dockerfile           # Docker 构建文件
├── docker-compose.yml   # Docker Compose
├── Makefile             # 构建命令
└── go.mod               # Go 模块定义
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 7.x+

### 安装依赖

```bash
go mod tidy
```

### 配置

复制并修改配置文件：

```bash
cp pkg/config/config.yaml config/config.yaml
```

编辑 `config/config.yaml`，配置数据库、Redis 等连接信息。

### 运行

```bash
# 开发模式
make run

# 或直接运行
go run cmd/server/main.go
```

### 构建

```bash
make build
```

### Docker

```bash
# 构建镜像
docker build -t questionhelper-server .

# 运行
docker-compose up -d
```

## API 文档

启动服务后访问：`http://localhost:8080/swagger/index.html`

## 目录说明

| 目录 | 说明 |
|------|------|
| `cmd/` | 程序入口 |
| `internal/controller/` | 控制器层，处理 HTTP 请求 |
| `internal/service/` | 业务逻辑层，核心业务处理 |
| `internal/repository/` | 数据访问层，数据库操作 |
| `internal/model/` | 数据模型，GORM 模型定义 |
| `internal/middleware/` | 中间件，JWT认证、权限校验等 |
| `internal/router/` | 路由定义 |
| `internal/dto/` | 数据传输对象 |
| `internal/task/` | 定时任务 |
| `internal/job/` | 异步任务 |
| `internal/ws/` | WebSocket |
| `pkg/` | 公共工具包 |
| `migrations/` | 数据库迁移文件 |
| `scripts/` | 脚本文件 |
| `docs/` | 文档 |

## License

MIT
