# Day 12 - 三层架构整合

## 今日目标

将 Handler、Service、Repository 三层架构完整整合，实现依赖注入、统一错误处理、以及 HTTP 状态码与业务错误码的映射。

---

## 知识点

### 三层架构职责

| 层 | 职责 | 关注点 |
|---|------|--------|
| Handler | 接收请求、参数绑定验证、调用 Service、返回响应 | HTTP 协议、请求/响应格式 |
| Service | 业务逻辑、参数校验、编排多个 Repository | 业务规则、数据组合 |
| Repository | 数据库操作、SQL 查询 | 数据持久化、查询优化 |

数据流向：`HTTP Request -> Handler -> Service -> Repository -> DB`

---

## 练习任务

### 任务 1：定义层间接口（难度：中等）

```go
// internal/repository/interfaces.go
type UserRepository interface {
    Create(user *model.User) error
    FindByID(id uint) (*model.User, error)
    FindByUsername(username string) (*model.User, error)
    FindByEmail(email string) (*model.User, error)
    Update(user *model.User) error
    Delete(id uint) error
    List(page, pageSize int, keyword string) ([]model.User, int64, error)
}

// internal/service/interfaces.go
type UserService interface {
    Register(req *dto.RegisterReq) (*dto.UserResponse, error)
    Login(req *dto.LoginReq) (*dto.TokenResponse, error)
    GetProfile(userID uint) (*dto.UserResponse, error)
    UpdateProfile(userID uint, req *dto.UpdateProfileReq) (*dto.UserResponse, error)
    ListUsers(req *dto.ListUserQuery) (*dto.PageResponse, error)
}

// internal/dto/user_dto.go（数据传输对象）
package dto

type RegisterReq struct {
    Username        string `json:"username" binding:"required,min=3,max=20"`
    Email           string `json:"email" binding:"required,email"`
    Password        string `json:"password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Avatar    string    `json:"avatar"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    ExpiresIn    int          `json:"expires_in"`
    User         UserResponse `json:"user"`
}

// 模型转 DTO 的方法
func ToUserResponse(user *model.User) *UserResponse {
    return &UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Phone:     user.Phone,
        Avatar:    user.Avatar,
        Role:      user.Role,
        CreatedAt: user.CreatedAt,
    }
}
```

### 任务 2：实现 Service 层（难度：中等偏上）

```go
// internal/service/user_service.go
package service

type userService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) Register(req *dto.RegisterReq) (*dto.UserResponse, error) {
    // 1. 检查用户名是否已存在
    existing, _ := s.userRepo.FindByUsername(req.Username)
    if existing != nil {
        return nil, &apperror.AppError{
            Code:    apperror.ErrDuplicateUsername,
            Message: "用户名已被注册",
        }
    }

    // 2. 检查邮箱是否已存在
    existing, _ = s.userRepo.FindByEmail(req.Email)
    if existing != nil {
        return nil, &apperror.AppError{
            Code:    apperror.ErrDuplicateEmail,
            Message: "邮箱已被注册",
        }
    }

    // 3. 密码加密
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("密码加密失败: %w", err)
    }

    // 4. 创建用户
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "user",
        Status:   1,
    }
    if err := s.userRepo.Create(user); err != nil {
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }

    return dto.ToUserResponse(user), nil
}

func (s *userService) Login(req *dto.LoginReq) (*dto.TokenResponse, error) {
    // 1. 查找用户
    // 2. 验证密码
    // 3. 生成 JWT Token
    // 4. 返回 TokenResponse
}

func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    return dto.ToUserResponse(user), nil
}

// ... 补充其他方法
```

### 任务 3：实现统一错误体系（难度：中等偏上）

```go
// internal/apperror/errors.go
package apperror

// 业务错误码
const (
    ErrBadRequest        = 40000
    ErrValidation        = 40001
    ErrDuplicateUsername = 40002
    ErrDuplicateEmail    = 40003
    ErrUnauthorized      = 40100
    ErrTokenExpired      = 40101
    ErrTokenInvalid      = 40102
    ErrForbidden         = 40300
    ErrNotFound          = 40400
    ErrUserNotFound      = 40401
    ErrArticleNotFound   = 40402
    ErrInternalServer    = 50000
)

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"` // 内部错误，不暴露给前端
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// HTTPStatusCode 根据业务错误码返回 HTTP 状态码
func (e *AppError) HTTPStatusCode() int {
    switch {
    case e.Code >= 50000:
        return http.StatusInternalServerError
    case e.Code >= 40400:
        return http.StatusNotFound
    case e.Code >= 40300:
        return http.StatusForbidden
    case e.Code >= 40100:
        return http.StatusUnauthorized
    case e.Code >= 40000:
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}

// 便捷构造函数
func NotFound(msg string) *AppError {
    return &AppError{Code: ErrNotFound, Message: msg}
}

func Unauthorized(msg string) *AppError {
    return &AppError{Code: ErrUnauthorized, Message: msg}
}

func BadRequest(msg string) *AppError {
    return &AppError{Code: ErrBadRequest, Message: msg}
}

func Internal(err error) *AppError {
    return &AppError{Code: ErrInternalServer, Message: "服务器内部错误", Err: err}
}
```

### 任务 4：实现 Handler 层与错误处理中间件（难度：中等）

```go
// internal/handler/user_handler.go
package handler

type UserHandler struct {
    userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
    var req dto.RegisterReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err)
        return
    }

    user, err := h.userService.Register(&req)
    if err != nil {
        response.HandleError(c, err)
        return
    }

    response.Created(c, user)
}

// pkg/response/response.go - 添加错误处理

func HandleError(c *gin.Context, err error) {
    var appErr *apperror.AppError
    if errors.As(err, &appErr) {
        c.JSON(appErr.HTTPStatusCode(), Response{
            Code:    appErr.Code,
            Message: appErr.Message,
        })
        return
    }

    // 未知错误统一返回 500
    c.JSON(http.StatusInternalServerError, Response{
        Code:    50000,
        Message: "服务器内部错误",
    })
}
```

### 任务 5：依赖注入与应用启动（难度：中等）

```go
// cmd/server/main.go
package main

func main() {
    // 1. 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // 2. 初始化数据库
    db, err := database.Init(&cfg.Database)
    if err != nil {
        log.Fatal(err)
    }

    // 3. 初始化各层（手动依赖注入）
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 4. 设置路由
    r := router.Setup(userHandler)

    // 5. 启动服务
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("服务启动: %s\n", addr)
    r.Run(addr)
}

// internal/router/router.go
func Setup(userHandler *handler.UserHandler) *gin.Engine {
    r := gin.New()
    r.Use(middleware.Recovery(), middleware.Logger(), middleware.CORS())

    v1 := r.Group("/api/v1")
    {
        auth := v1.Group("/auth")
        {
            auth.POST("/register", userHandler.Register)
            auth.POST("/login", userHandler.Login)
        }

        users := v1.Group("/users")
        users.Use(middleware.JWTAuth()) // 需要认证
        {
            users.GET("/profile", userHandler.GetProfile)
            users.PUT("/profile", userHandler.UpdateProfile)
            users.GET("", userHandler.List)
        }
    }

    return r
}
```

---

## 自检清单

- [ ] 三层架构职责清晰，没有跨层调用
- [ ] Handler 层不包含业务逻辑，只做参数绑定和响应
- [ ] Service 层不直接操作 gin.Context
- [ ] Repository 层不包含业务逻辑，只做数据库操作
- [ ] 依赖注入通过构造函数实现，层间通过接口解耦
- [ ] 统一错误体系能正确映射 HTTP 状态码
- [ ] DTO 与 Model 分离，不直接暴露数据库模型给前端
