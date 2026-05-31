.PHONY: build run clean test migrate

# 构建
build:
	go build -o bin/server cmd/server/main.go

# 运行
run:
	go run cmd/server/main.go

# 清理
clean:
	rm -rf bin/

# 测试
test:
	go test ./... -v

# 数据库迁移
migrate:
	go run cmd/migrate/main.go

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...

# 生成 swagger 文档
swagger:
	swag init -g cmd/server/main.go -o docs/swagger
