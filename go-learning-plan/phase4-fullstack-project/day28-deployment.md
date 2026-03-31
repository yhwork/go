# Day 28 - 部署上线

## 今日目标

使用 Docker Compose 编排完整服务、Nginx 反向代理配置、健康检查、优雅关闭、Makefile 统一管理。

---

## 练习任务

### 任务 1：完整的 docker-compose 编排

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - frontend-dist:/usr/share/nginx/html
    depends_on:
      - app
    restart: always

  app:
    build:
      context: .
      dockerfile: Dockerfile
    expose:
      - "8080"
    environment:
      - APP_SERVER_MODE=release
      - APP_DATABASE_HOST=mysql
      - APP_REDIS_HOST=redis
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: always
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD_FILE: /run/secrets/db_password
      MYSQL_DATABASE: go_blog
    volumes:
      - mysql-data:/var/lib/mysql
    secrets:
      - db_password
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: always

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: always

volumes:
  mysql-data:
  redis-data:
  frontend-dist:

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

### 任务 2：Nginx 反向代理

```nginx
# nginx/conf.d/default.conf

upstream go_blog_api {
    server app:8080;
}

server {
    listen 80;
    server_name localhost;

    # 前端静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://go_blog_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 代理
    location /ws {
        proxy_pass http://go_blog_api;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }

    # Swagger 文档（开发环境可开放）
    location /swagger/ {
        proxy_pass http://go_blog_api;
    }

    # 静态上传文件
    location /static/uploads/ {
        alias /app/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    # Gzip 压缩
    gzip on;
    gzip_types text/plain application/json application/javascript text/css;
    gzip_min_length 1024;
}
```

### 任务 3：健康检查接口

```go
func HealthCheck(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        health := gin.H{"status": "healthy"}
        httpStatus := 200

        // 检查数据库
        sqlDB, _ := db.DB()
        if err := sqlDB.Ping(); err != nil {
            health["database"] = err.Error()
            health["status"] = "unhealthy"
            httpStatus = 503
        } else {
            health["database"] = "connected"
        }

        // 检查 Redis
        if err := rdb.Ping(context.Background()).Err(); err != nil {
            health["redis"] = err.Error()
            health["status"] = "unhealthy"
            httpStatus = 503
        } else {
            health["redis"] = "connected"
        }

        c.JSON(httpStatus, health)
    }
}
```

### 任务 4：优雅关闭（整合 Day 18）

确保 main.go 中完整实现优雅关闭逻辑，关闭顺序：
1. 停止接收新 HTTP 请求
2. 等待进行中的请求处理完成（最多 30 秒）
3. 停止 Cron 定时任务
4. 停止队列 Worker
5. 关闭 WebSocket Hub
6. 关闭数据库连接
7. 关闭 Redis 连接

### 任务 5：完整的 Makefile

```makefile
.PHONY: all build run test lint clean docker-up docker-down

# 变量
APP_NAME=go-blog
MAIN=./cmd/server

# 开发
run:
	go run $(MAIN)

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(APP_NAME) $(MAIN)

test:
	go test ./... -v -cover -coverprofile=coverage.out

test-report:
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	gofmt -w .

tidy:
	go mod tidy

# Swagger
swagger:
	swag init -g cmd/server/main.go -o docs

# Docker
docker-dev:
	docker-compose up -d --build

docker-prod:
	docker-compose -f docker-compose.prod.yml up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-clean:
	docker-compose down -v

# 数据库
migrate:
	go run $(MAIN) migrate

seed:
	go run $(MAIN) seed

# 一键构建
all: tidy fmt lint test build
```

### 任务 6：部署验证

```bash
# 启动生产环境
make docker-prod

# 验证健康检查
curl http://localhost/health

# 验证 API
curl http://localhost/api/v1/articles

# 验证前端页面
# 浏览器访问 http://localhost

# 验证 WebSocket
# 用浏览器开发工具检查 ws://localhost/ws

# 查看日志
make docker-logs
```

---

## 自检清单

- [ ] docker-compose 四个服务全部正常启动
- [ ] Nginx 反向代理配置正确（API + 前端 + WebSocket）
- [ ] 健康检查接口返回各组件状态
- [ ] Ctrl+C 优雅关闭，无请求丢失
- [ ] Makefile 命令全部可用
- [ ] 通过 Nginx 访问前端页面和 API 正常
- [ ] 生产环境不暴露 Swagger 文档（或设置密码保护）
