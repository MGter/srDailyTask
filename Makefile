.PHONY: build run clean test deps frontend pack

APP_NAME=daily_task

# 构建后端
build:
	go build -o $(APP_NAME) cmd/main.go

# 构建前端
frontend:
	cd web && npm run build
	cp pics/kita.webp web/dist/assets/kita.webp

# 完整打包（前端+后端）
pack: frontend build
	@echo "打包完成: $(APP_NAME)"

# 运行开发模式
run:
	go run cmd/main.go

# 清理
clean:
	rm -f $(APP_NAME)
	rm -rf web/dist

# 测试
test:
	go test ./...

# 下载依赖
deps:
	go mod download
	go mod tidy
	cd web && npm install