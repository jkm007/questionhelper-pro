# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装 ca 证书（HTTPS 请求需要）
RUN apk --no-cache add ca-certificates tzdata

# 复制构建产物
COPY --from=builder /app/server .
COPY --from=builder /app/pkg/config/config.yaml ./config/

# 暴露端口
EXPOSE 8080

# 启动
CMD ["./server"]
