# QuestionHelper 部署指南

## 项目结构

```
questionhelper-pro/
├── questionhelper-server/    # 后端服务 (Go + Gin + GORM)
├── questionhelper-admin/     # 前端管理后台 (Vue3 + Element Plus)
├── questionhelper-app/       # 移动端应用 (UniApp)
├── docs/                     # 项目文档
├── docker-compose.yml        # Docker 编排配置
└── DEPLOYMENT.md            # 部署文档
```

## 环境要求

### 开发环境

- **Go**: 1.21+
- **Node.js**: 18+
- **pnpm**: 8+
- **MySQL**: 8.0
- **Redis**: 7.x

### 生产环境

- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **服务器内存**: 2GB+
- **磁盘空间**: 10GB+

## 快速部署（Docker）

### 1. 克隆代码

```bash
git clone <repository-url>
cd questionhelper-pro
```

### 2. 配置环境变量

创建 `.env` 文件：

```bash
# 数据库配置
MYSQL_ROOT_PASSWORD=your_secure_password
MYSQL_DATABASE=questionhelper

# JWT 密钥
JWT_SECRET=your_jwt_secret_key

# 服务器配置
SERVER_PORT=8080
GIN_MODE=release
```

### 3. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 访问服务

- **前端管理后台**: http://localhost:3000
- **后端 API**: http://localhost:8080
- **MySQL**: localhost:3306
- **Redis**: localhost:6379

### 5. 默认账号

- **管理员**: admin / admin123

## 手动部署

### 1. 部署后端服务

#### 1.1 安装依赖

```bash
cd questionhelper-server
go mod download
```

#### 1.2 配置数据库

编辑 `config/config.yaml`：

```yaml
server:
  port: 8080
  mode: release

mysql:
  host: your_mysql_host
  port: 3306
  user: root
  password: your_password
  dbname: questionhelper

redis:
  host: your_redis_host
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your_jwt_secret
  expire: 7200
  refresh_expire: 604800
```

#### 1.3 编译运行

```bash
# 编译
go build -o bin/server ./cmd/server/main.go

# 运行
./bin/server
```

### 2. 部署前端服务

#### 2.1 安装依赖

```bash
cd questionhelper-admin
pnpm install
```

#### 2.2 配置环境变量

编辑 `.env.production`：

```bash
VITE_APP_BASE_API=/prod-api
VITE_APP_API_URL=http://your_server_ip:8080
```

#### 2.3 构建

```bash
pnpm build
```

#### 2.4 部署到 Nginx

```nginx
server {
    listen       80;
    server_name  your_domain.com;

    location / {
        root   /path/to/dist;
        index  index.html index.htm;
        try_files $uri $uri/ /index.html;
    }

    location /prod-api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 数据库初始化

### 1. 创建数据库

```sql
CREATE DATABASE questionhelper DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 运行迁移

后端服务启动时会自动创建表结构。

### 3. 插入种子数据

```bash
# 连接数据库
mysql -u root -p questionhelper

# 插入默认角色
INSERT INTO roles (name, code, description, is_default, sort, status) VALUES
('超级管理员', 'admin', '系统超级管理员', 0, 1, 1),
('普通用户', 'user', '普通用户', 1, 2, 1),
('教师', 'teacher', '教师角色', 0, 3, 1);

# 插入默认菜单
INSERT INTO menus (parent_id, name, path, component, title, icon, type, sort, status) VALUES
(NULL, 'Dashboard', '/dashboard', 'Layout', '仪表盘', 'homepage', 1, 1, 1),
(NULL, 'User', '/user', 'Layout', '用户管理', 'user', 1, 2, 1),
(NULL, 'Question', '/question', 'Layout', '题库管理', 'question', 1, 3, 1),
(NULL, 'Exam', '/exam', 'Layout', '考试管理', 'exam', 1, 4, 1),
(NULL, 'Class', '/class', 'Layout', '班级管理', 'class', 1, 5, 1),
(NULL, 'Practice', '/practice', 'Layout', '练习管理', 'practice', 1, 6, 1),
(NULL, 'Statistics', '/statistics', 'Layout', '统计分析', 'statistics', 1, 7, 1),
(NULL, 'System', '/system', 'Layout', '系统管理', 'system', 1, 10, 1);

# 给管理员角色分配菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus;
```

## 生产环境优化

### 1. 安全配置

- 修改默认密码
- 配置 HTTPS
- 设置防火墙规则
- 配置 CORS 白名单

### 2. 性能优化

- 配置 Redis 缓存
- 启用 Gzip 压缩
- 配置 CDN 加速
- 数据库索引优化

### 3. 监控配置

- 配置日志收集
- 设置监控告警
- 配置备份策略

## 常见问题

### 1. 数据库连接失败

检查 MySQL 服务是否启动，配置是否正确。

### 2. Redis 连接失败

检查 Redis 服务是否启动，密码是否正确。

### 3. 前端无法访问后端 API

检查 Nginx 配置，确保 API 代理正确。

### 4. 文件上传失败

检查上传目录权限，确保目录可写。

## 维护命令

```bash
# 查看服务状态
docker-compose ps

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f server
docker-compose logs -f admin

# 停止服务
docker-compose down

# 清理数据
docker-compose down -v
```

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并部署
docker-compose up -d --build

# 清理旧镜像
docker image prune -f
```

## 备份恢复

### 备份数据库

```bash
docker exec questionhelper-mysql mysqldump -u root -p questionhelper > backup.sql
```

### 恢复数据库

```bash
docker exec -i questionhelper-mysql mysql -u root -p questionhelper < backup.sql
```

## 技术支持

如有问题，请查看：
- 项目文档：docs/
- 日志文件：docker-compose logs
- 配置文件：config/

## License

MIT License
