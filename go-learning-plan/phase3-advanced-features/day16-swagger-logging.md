# Day 16 - Swagger 文档与日志

## 今日目标

集成 Swagger 自动生成 API 文档，集成 Zap 结构化日志，实现按级别输出和请求链路追踪。

---

## 练习任务

### 任务 1：Swagger 文档集成

```bash
# 安装 swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# 安装 Gin Swagger 中间件
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
```

```go
// cmd/server/main.go 顶部添加全局注解
// @title           Go Blog API
// @version         1.0
// @description     Go 全栈博客系统 API 文档
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// 为每个接口添加注解，例如：
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterReq true "注册信息"
// @Success      201  {object} response.Response{data=dto.UserResponse}
// @Failure      400  {object} response.Response
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {}

// 生成文档
// swag init -g cmd/server/main.go -o docs

// 注册 Swagger 路由
import swaggerFiles "github.com/swaggo/files"
import ginSwagger "github.com/swaggo/gin-swagger"
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

为所有已实现的接口添加 Swagger 注解。

### 任务 2：Zap 日志集成

```go
// pkg/logger/logger.go

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Init(level string, format string, output string) error {
    // 1. 根据 level 参数设置日志级别
    // 2. 根据 format 选择 console 或 json 编码
    // 3. 根据 output 选择 stdout 或文件输出
    // 4. 初始化全局 Log 实例
}

// 在请求日志中间件中使用：
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        requestID := uuid.New().String()[:8]
        c.Set("requestID", requestID)

        c.Next()

        logger.Log.Infow("HTTP请求",
            "request_id", requestID,
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "latency", time.Since(start).String(),
            "ip", c.ClientIP(),
        )
    }
}
```

### 任务 3：不同环境的日志配置

```go
// 开发环境：Console 格式，Debug 级别，输出到 stdout
// 生产环境：JSON 格式，Info 级别，输出到文件 + stdout

// 日志文件轮转（使用 lumberjack）
import "gopkg.in/natefinished/lumberjack.v2"

writer := &lumberjack.Logger{
    Filename:   "logs/app.log",
    MaxSize:    100, // MB
    MaxBackups: 30,
    MaxAge:     7,   // 天
    Compress:   true,
}
```

---

## 自检清单

- [ ] 所有接口都有 Swagger 注解
- [ ] 访问 /swagger/index.html 能看到文档
- [ ] Swagger UI 可以直接测试接口（Try it out）
- [ ] Zap 日志正确输出，包含 request_id
- [ ] 开发和生产环境使用不同的日志配置
- [ ] 日志文件轮转配置正确
