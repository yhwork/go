# Day 19 - Docker 容器化

## 今日目标

编写 Dockerfile（多阶段构建）、docker-compose 编排（Go + MySQL + Redis）、环境变量配置注入、容器中运行验证。

---

## 练习任务

### 任务 1：多阶段 Dockerfile

```dockerfile
# ============ Build Stage ============
FROM golang:1.22-alpine AS builder

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译（静态链接，禁用 CGO）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# ============ Run Stage ============
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/config/config.yaml ./config/

# 创建上传目录
RUN mkdir -p uploads/images

EXPOSE 8080

ENTRYPOINT ["./server"]
```

### 任务 2：docker-compose 编排

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_SERVER_PORT=8080
      - APP_DATABASE_HOST=mysql
      - APP_DATABASE_PORT=3306
      - APP_DATABASE_USERNAME=root
      - APP_DATABASE_PASSWORD=rootpassword
      - APP_DATABASE_DBNAME=go_blog
      - APP_REDIS_HOST=redis
      - APP_REDIS_PORT=6379
      - APP_JWT_SECRET=my-super-secret-key
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - app-uploads:/app/uploads
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: go_blog
      MYSQL_CHARSET: utf8mb4
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  mysql-data:
  redis-data:
  app-uploads:
```

### 任务 3：运行和验证

```bash
# 构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f app

# 测试接口
curl http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"Pass1234","confirm_password":"Pass1234"}'

# 停止
docker-compose down

# 停止并删除数据
docker-compose down -v
```

### 任务 4：健康检查接口

```go
// GET /health
r.GET("/health", func(c *gin.Context) {
    // 检查数据库连接
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        c.JSON(503, gin.H{"status": "unhealthy", "database": err.Error()})
        return
    }
    // 检查 Redis 连接
    if err := rdb.Ping(context.Background()).Err(); err != nil {
        c.JSON(503, gin.H{"status": "unhealthy", "redis": err.Error()})
        return
    }
    c.JSON(200, gin.H{"status": "healthy"})
})
```

### 任务 5：Makefile 补充 Docker 命令

```makefile
# Docker 相关
docker-build:
	docker build -t go-blog:latest .

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-clean:
	docker-compose down -v
	docker rmi go-blog:latest
```

---

## 自检清单

- [ ] Docker 镜像构建成功，体积小（< 30MB）
- [ ] docker-compose 三个服务都正常启动
- [ ] 健康检查接口能正确反映各组件状态
- [ ] 环境变量正确覆盖配置文件中的值
- [ ] 数据持久化到 volume，重启不丢失
- [ ] 所有 API 接口在容器中正常工作
