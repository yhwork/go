# Day 07 - Gin 入门与路由

## 今日目标

掌握 Gin 框架的基础使用：引擎创建、路由注册、路由分组、参数获取（路径参数、查询参数、JSON Body），以及统一响应格式的封装。

---

## 知识点

### 1. Gin 快速上手

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default() // 包含 Logger 和 Recovery 中间件

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080") // 默认监听 0.0.0.0:8080
}
```

### 2. 路由参数

```go
// 路径参数
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id") // string 类型
})

// 查询参数
// GET /users?page=1&size=10
r.GET("/users", func(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    size := c.DefaultQuery("size", "10")
})

// JSON Body
r.POST("/users", func(c *gin.Context) {
    var req CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
})
```

### 3. 路由分组

```go
v1 := r.Group("/api/v1")
{
    users := v1.Group("/users")
    {
        users.GET("", listUsers)
        users.GET("/:id", getUser)
        users.POST("", createUser)
        users.PUT("/:id", updateUser)
        users.DELETE("/:id", deleteUser)
    }

    products := v1.Group("/products")
    {
        products.GET("", listProducts)
        products.GET("/:id", getProduct)
    }
}
```

### 4. gin.Context 常用方法

```go
func handler(c *gin.Context) {
    // 获取参数
    c.Param("id")                  // 路径参数
    c.Query("page")                // 查询参数
    c.DefaultQuery("size", "10")   // 带默认值的查询参数
    c.PostForm("username")         // 表单参数

    // 绑定请求体
    c.ShouldBindJSON(&req)         // JSON
    c.ShouldBindQuery(&req)        // Query 参数
    c.ShouldBind(&req)             // 自动判断

    // 返回响应
    c.JSON(200, data)              // JSON 响应
    c.String(200, "hello")         // 纯文本
    c.File("path/to/file")        // 文件下载

    // 设置/获取上下文值
    c.Set("userID", 42)
    userID, _ := c.Get("userID")

    // 中止请求
    c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
}
```

---

## 练习任务

### 任务 1：搭建 User CRUD API（难度：中等）

创建一个完整的用户 CRUD API，暂时使用内存存储（map）：

```go
// model/user.go
type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 请求结构体
type CreateUserReq struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
}

type UpdateUserReq struct {
    Username string `json:"username,omitempty"`
    Email    string `json:"email,omitempty"`
    Phone    string `json:"phone,omitempty"`
}

// 列表查询参数
type ListUserQuery struct {
    Page     int    `form:"page"`
    PageSize int    `form:"page_size"`
    Keyword  string `form:"keyword"`  // 搜索用户名或邮箱
}
```

实现以下接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/users | 用户列表（分页+搜索） |
| GET | /api/v1/users/:id | 用户详情 |
| POST | /api/v1/users | 创建用户 |
| PUT | /api/v1/users/:id | 更新用户 |
| DELETE | /api/v1/users/:id | 删除用户 |

### 任务 2：统一响应格式封装（难度：中等）

```go
// pkg/response/response.go

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PageResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    PageData    `json:"data"`
}

type PageData struct {
    List     interface{} `json:"list"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"page_size"`
}

// 封装便捷函数，直接操作 gin.Context

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Created(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, Response{
        Code:    0,
        Message: "created",
        Data:    data,
    })
}

func PageOK(c *gin.Context, list interface{}, total int64, page, pageSize int) {
    c.JSON(http.StatusOK, PageResponse{
        Code:    0,
        Message: "success",
        Data: PageData{
            List:     list,
            Total:    total,
            Page:     page,
            PageSize: pageSize,
        },
    })
}

func BadRequest(c *gin.Context, msg string) {
    c.JSON(http.StatusBadRequest, Response{
        Code:    400,
        Message: msg,
    })
}

func NotFound(c *gin.Context, msg string) {
    c.JSON(http.StatusNotFound, Response{
        Code:    404,
        Message: msg,
    })
}

func ServerError(c *gin.Context, msg string) {
    c.JSON(http.StatusInternalServerError, Response{
        Code:    500,
        Message: msg,
    })
}
```

### 任务 3：实现分页逻辑（难度：基础）

```go
// 在 handler 中实现分页

func (h *UserHandler) List(c *gin.Context) {
    var query ListUserQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }

    // 设置默认值
    if query.Page <= 0 {
        query.Page = 1
    }
    if query.PageSize <= 0 {
        query.PageSize = 10
    }
    if query.PageSize > 100 {
        query.PageSize = 100 // 最大 100 条
    }

    // 从内存存储中获取数据，支持 keyword 搜索
    // 计算分页偏移量
    // 返回分页响应
}
```

### 任务 4：用 curl/Postman 测试所有接口

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","phone":"13800138000"}'

# 获取用户列表
curl http://localhost:8080/api/v1/users?page=1&page_size=10

# 搜索用户
curl http://localhost:8080/api/v1/users?keyword=alice

# 获取用户详情
curl http://localhost:8080/api/v1/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"email":"newalice@example.com"}'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1
```

逐一测试并验证：
- 正常请求返回正确的响应格式
- ID 不存在返回 404
- 参数错误返回 400
- 分页参数生效
- 搜索功能正常

---

## 参考答案骨架

```go
package main

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// ============ 内存存储 ============

type MemoryStore struct {
    mu     sync.RWMutex
    users  map[uint]*User
    nextID uint
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        users:  make(map[uint]*User),
        nextID: 1,
    }
}

func (s *MemoryStore) Create(user *User) {
    s.mu.Lock()
    defer s.mu.Unlock()
    user.ID = s.nextID
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    s.nextID++
    s.users[user.ID] = user
}

func (s *MemoryStore) FindByID(id uint) (*User, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    user, ok := s.users[id]
    return user, ok
}

// ... 补充 Update, Delete, List 方法

// ============ Handler ============

type UserHandler struct {
    store *MemoryStore
}

func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
        return
    }

    user := &User{
        Username: req.Username,
        Email:    req.Email,
        Phone:    req.Phone,
    }
    h.store.Create(user)

    c.JSON(http.StatusCreated, Response{Code: 0, Message: "created", Data: user})
}

// ... 补充其他 handler

// ============ 路由 ============

func setupRouter(h *UserHandler) *gin.Engine {
    r := gin.Default()

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

func main() {
    store := NewMemoryStore()
    handler := &UserHandler{store: store}
    r := setupRouter(handler)
    r.Run(":8080")
}
```

---

## 自检清单

- [ ] Gin 服务能正常启动并响应请求
- [ ] 所有 5 个 CRUD 接口实现完整
- [ ] 统一响应格式正确（code/message/data）
- [ ] 分页查询正常工作
- [ ] 搜索功能正常工作
- [ ] 错误情况返回正确的 HTTP 状态码
- [ ] 使用 curl 或 Postman 测试通过所有接口
