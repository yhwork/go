# Day 22 - 用户模块完整实现

## 今日目标

完成博客系统用户模块的完整实现：注册、登录、个人信息管理、角色权限、头像上传、单元测试。

---

## 练习任务

### 任务 1：用户 Repository

整合 Day 10-11 的内容，实现完整的 UserRepository。

### 任务 2：用户 Service

```go
// 接口定义
type UserService interface {
    Register(req *dto.RegisterReq) (*dto.UserResponse, error)
    Login(req *dto.LoginReq) (*dto.TokenResponse, error)
    Logout(token string) error
    RefreshToken(refreshToken string) (*dto.TokenResponse, error)
    GetProfile(userID uint) (*dto.UserResponse, error)
    UpdateProfile(userID uint, req *dto.UpdateProfileReq) (*dto.UserResponse, error)
    UpdatePassword(userID uint, req *dto.UpdatePasswordReq) error
    UpdateAvatar(userID uint, avatarURL string) (*dto.UserResponse, error)
    ListUsers(query *dto.ListUserQuery) (*dto.PageResult, error) // admin only
}
```

实现所有方法，整合 JWT 认证和 Redis 缓存。

### 任务 3：用户 Handler + 路由

```go
// 路由定义
auth := v1.Group("/auth")
{
    auth.POST("/register", userHandler.Register)
    auth.POST("/login", userHandler.Login)
    auth.POST("/refresh", userHandler.RefreshToken)
}

user := v1.Group("/user")
user.Use(middleware.JWTAuth())
{
    user.GET("/profile", userHandler.GetProfile)
    user.PUT("/profile", userHandler.UpdateProfile)
    user.PUT("/password", userHandler.UpdatePassword)
    user.POST("/avatar", userHandler.UploadAvatar)
    user.POST("/logout", userHandler.Logout)
}

// Admin 路由
admin := v1.Group("/admin")
admin.Use(middleware.JWTAuth(), middleware.RequireRole("admin"))
{
    admin.GET("/users", userHandler.ListUsers)
}
```

### 任务 4：单元测试

为 UserService 的关键方法编写表驱动测试：
- Register：正常注册、用户名重复、邮箱重复、密码不一致
- Login：正常登录、用户不存在、密码错误、账号被禁用
- GetProfile：正常获取、用户不存在
- UpdatePassword：正常修改、旧密码错误

### 任务 5：接口测试

用 curl 逐一测试所有用户相关接口，验证正常和异常场景。

---

## 自检清单

- [ ] 注册/登录/退出登录完整流程正常
- [ ] Token 刷新机制工作正常
- [ ] 个人信息修改和密码修改正常
- [ ] 头像上传正常
- [ ] Admin 权限验证正确（普通用户无法访问 admin 接口）
- [ ] 单元测试覆盖率 > 80%
