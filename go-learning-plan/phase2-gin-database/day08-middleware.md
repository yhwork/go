# Day 08 - 中间件机制

## 今日目标

理解 Gin 中间件的运行原理（洋葱模型）、学会编写自定义中间件（Logger、Recovery、CORS、限流），以及中间件链的组合使用。

---

## 知识点

### 1. 中间件原理

```go
// Gin 中间件本质就是 gin.HandlerFunc
type HandlerFunc func(*Context)

// 中间件执行流程（洋葱模型）：
// Request -> Logger -> Auth -> Handler -> Auth(后续) -> Logger(后续) -> Response

func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ===== 请求前（Before）=====
        start := time.Now()

        c.Next() // 调用后续的中间件和 Handler

        // ===== 请求后（After）=====
        latency := time.Since(start)
        fmt.Printf("耗时: %v\n", latency)
    }
}
```

### 2. 中间件注册

```go
// 全局中间件
r := gin.New() // 不含默认中间件
r.Use(Logger(), Recovery())

// 路由组中间件
auth := r.Group("/api", AuthMiddleware())
{
    auth.GET("/profile", getProfile)
}

// 单个路由中间件
r.GET("/admin", AdminOnly(), adminHandler)
```

### 3. c.Next() 与 c.Abort()

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "未提供 Token"})
            return // Abort 后必须 return，不会调用后续 handler
        }
        c.Set("userID", 42)
        c.Next() // 继续执行后续 handler
    }
}
```

---

## 练习任务

### 任务 1：Logger 中间件（难度：中等）

```go
// 实现一个详细的请求日志中间件

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ===== Before =====
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery
        method := c.Request.Method
        clientIP := c.ClientIP()
        requestID := uuid.New().String() // 或自己生成

        // 将 requestID 设置到 context 中，后续可用于链路追踪
        c.Set("requestID", requestID)
        c.Header("X-Request-ID", requestID)

        c.Next()

        // ===== After =====
        latency := time.Since(start)
        statusCode := c.Writer.Status()
        bodySize := c.Writer.Size()

        // 输出格式：
        // [2024-01-01 12:00:00] [req-id] 200 | 12.34ms | 127.0.0.1 | GET /api/v1/users?page=1 | 1234 bytes
        fmt.Printf("[%s] [%s] %d | %13v | %15s | %-7s %s%s | %d bytes\n",
            time.Now().Format("2006-01-02 15:04:05"),
            requestID[:8],
            statusCode,
            latency,
            clientIP,
            method,
            path,
            formatQuery(query),
            bodySize,
        )

        // 额外要求：
        // - 4xx 状态码用黄色输出（可选，用 ANSI 颜色码）
        // - 5xx 状态码用红色输出
        // - 慢请求（> 500ms）额外标记 [SLOW]
    }
}

func formatQuery(query string) string {
    if query != "" {
        return "?" + query
    }
    return ""
}
```

### 任务 2：Recovery 中间件（难度：中等）

```go
// 捕获 panic，返回 500 错误而不是让服务崩溃

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 1. 获取堆栈信息
                stack := debug.Stack()

                // 2. 记录错误日志
                fmt.Printf("[PANIC RECOVERY] %v\n%s\n", err, string(stack))

                // 3. 获取 requestID（从 Logger 中间件设置的）
                requestID, _ := c.Get("requestID")

                // 4. 返回统一错误响应
                c.AbortWithStatusJSON(500, gin.H{
                    "code":       500,
                    "message":    "服务器内部错误",
                    "request_id": requestID,
                })
            }
        }()

        c.Next()
    }
}

// 测试：创建一个会 panic 的路由
r.GET("/panic", func(c *gin.Context) {
    panic("something went wrong!")
})
```

### 任务 3：CORS 中间件（难度：中等）

```go
// 实现跨域资源共享中间件

type CORSConfig struct {
    AllowOrigins     []string // 允许的域名，["*"] 表示所有
    AllowMethods     []string // 允许的 HTTP 方法
    AllowHeaders     []string // 允许的请求头
    ExposeHeaders    []string // 暴露给前端的响应头
    AllowCredentials bool     // 是否允许携带 Cookie
    MaxAge           int      // 预检请求缓存时间（秒）
}

func DefaultCORSConfig() CORSConfig {
    return CORSConfig{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
        ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
        AllowCredentials: false,
        MaxAge:           86400,
    }
}

func CORS(config ...CORSConfig) gin.HandlerFunc {
    cfg := DefaultCORSConfig()
    if len(config) > 0 {
        cfg = config[0]
    }

    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // 1. 检查 origin 是否在允许列表中
        // 2. 设置 CORS 响应头：
        //    Access-Control-Allow-Origin
        //    Access-Control-Allow-Methods
        //    Access-Control-Allow-Headers
        //    Access-Control-Expose-Headers
        //    Access-Control-Allow-Credentials
        //    Access-Control-Max-Age
        // 3. 处理 OPTIONS 预检请求（直接返回 204）
        // 4. 其他请求调用 c.Next()

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

### 任务 4：请求限流中间件（难度：中等偏上）

```go
// 使用令牌桶算法实现接口限流

import "golang.org/x/time/rate"

// 方式一：全局限流（所有请求共享一个限流器）
func RateLimiter(rps int, burst int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), burst)

    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "code":    429,
                "message": "请求过于频繁，请稍后重试",
            })
            return
        }
        c.Next()
    }
}

// 方式二：基于 IP 的限流（每个 IP 独立限流器）
func IPRateLimiter(rps int, burst int) gin.HandlerFunc {
    // 使用 sync.Map 存储每个 IP 的限流器
    var limiters sync.Map

    return func(c *gin.Context) {
        ip := c.ClientIP()

        // 获取或创建该 IP 的限流器
        val, _ := limiters.LoadOrStore(ip, rate.NewLimiter(rate.Limit(rps), burst))
        limiter := val.(*rate.Limiter)

        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "code":    429,
                "message": "请求过于频繁，请稍后重试",
            })
            return
        }

        c.Next()
    }

    // 思考：limiters map 会无限增长，如何清理过期的 limiter？
    // 提示：可以存储 limiter 和最后访问时间，定期清理
}
```

### 任务 5：将中间件整合到 Day 07 的项目中

```go
func setupRouter(h *UserHandler) *gin.Engine {
    r := gin.New() // 不用 Default，自己组装中间件

    // 全局中间件
    r.Use(Recovery())
    r.Use(Logger())
    r.Use(CORS())
    r.Use(RateLimiter(100, 200)) // 全局 100 QPS

    v1 := r.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("", h.List)
            users.GET("/:id", h.Get)
            users.POST("", h.Create)
            users.PUT("/:id", h.Update)
            users.DELETE("/:id", h.Delete)
        }
    }

    return r
}
```

---

## 参考答案骨架

```go
package middleware

import (
    "fmt"
    "net/http"
    "runtime/debug"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()

        logLine := fmt.Sprintf("[%s] %3d | %13v | %s | %s",
            time.Now().Format("2006-01-02 15:04:05"),
            status,
            latency,
            method,
            path,
        )

        if latency > 500*time.Millisecond {
            logLine = "[SLOW] " + logLine
        }

        fmt.Println(logLine)
    }
}

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("[PANIC] %v\n%s\n", err, debug.Stack())
                c.AbortWithStatusJSON(500, gin.H{
                    "code":    500,
                    "message": "服务器内部错误",
                })
            }
        }()
        c.Next()
    }
}

// ... 补充 CORS 和 RateLimiter
```

---

## 自检清单

- [ ] 理解中间件的洋葱模型执行顺序
- [ ] Logger 中间件正确记录请求方法、路径、状态码、耗时
- [ ] Recovery 中间件能捕获 panic 并返回 500
- [ ] CORS 中间件正确处理预检请求
- [ ] 限流中间件能在请求过多时返回 429
- [ ] 所有中间件与 Day 07 的 CRUD API 整合后正常工作
