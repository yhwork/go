# Day 13 - JWT 认证

## 今日目标

实现完整的 JWT 认证流程：用户注册（密码加密）、登录（生成 Token）、认证中间件（Token 校验）、Token 刷新。

---

## 知识点

### JWT 结构

```
Header.Payload.Signature

Header:  {"alg": "HS256", "typ": "JWT"}
Payload: {"user_id": 1, "username": "alice", "exp": 1700000000}
Signature: HMACSHA256(base64(header) + "." + base64(payload), secret)
```

### 双 Token 机制

- **Access Token**：短期（如 2 小时），用于接口认证
- **Refresh Token**：长期（如 7 天），用于刷新 Access Token
- Access Token 过期后，用 Refresh Token 获取新的 Access Token，无需重新登录

---

## 练习任务

### 任务 1：JWT 工具封装（难度：中等）

```go
// pkg/jwt/jwt.go
package jwt

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type JWTManager struct {
    secretKey       []byte
    accessExpiry    time.Duration
    refreshExpiry   time.Duration
    issuer          string
}

func NewJWTManager(secret string, accessHours, refreshDays int) *JWTManager {
    return &JWTManager{
        secretKey:     []byte(secret),
        accessExpiry:  time.Duration(accessHours) * time.Hour,
        refreshExpiry: time.Duration(refreshDays) * 24 * time.Hour,
        issuer:        "go-blog",
    }
}

// GenerateAccessToken 生成访问令牌
func (j *JWTManager) GenerateAccessToken(userID uint, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    j.issuer,
            Subject:   fmt.Sprintf("%d", userID),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
    // 刷新令牌只包含 userID 和过期时间
}

// ParseToken 解析并验证令牌
func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}

// GenerateTokenPair 同时生成 access_token 和 refresh_token
func (j *JWTManager) GenerateTokenPair(userID uint, username, role string) (accessToken, refreshToken string, err error) {
    accessToken, err = j.GenerateAccessToken(userID, username, role)
    if err != nil {
        return
    }
    refreshToken, err = j.GenerateRefreshToken(userID)
    return
}
```

### 任务 2：认证中间件（难度：中等）

```go
// internal/middleware/auth.go

func JWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取 Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, response.Error(40100, "请先登录"))
            return
        }

        // 2. 检查 Bearer 前缀
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(401, response.Error(40102, "Token 格式错误"))
            return
        }

        // 3. 解析 Token
        claims, err := jwtManager.ParseToken(parts[1])
        if err != nil {
            if errors.Is(err, jwt.ErrTokenExpired) {
                c.AbortWithStatusJSON(401, response.Error(40101, "Token 已过期"))
            } else {
                c.AbortWithStatusJSON(401, response.Error(40102, "Token 无效"))
            }
            return
        }

        // 4. 将用户信息写入 Context
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}

// 权限中间件
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(401, response.Error(40100, "请先登录"))
            return
        }

        for _, role := range roles {
            if userRole.(string) == role {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(403, response.Error(40300, "没有权限执行此操作"))
    }
}

// 获取当前用户 ID 的辅助函数
func GetCurrentUserID(c *gin.Context) (uint, error) {
    userID, exists := c.Get("userID")
    if !exists {
        return 0, errors.New("未找到用户信息")
    }
    return userID.(uint), nil
}
```

### 任务 3：注册与登录接口（难度：中等）

```go
// 完善 UserService 的 Register 和 Login 方法

func (s *userService) Login(req *dto.LoginReq) (*dto.TokenResponse, error) {
    // 1. 根据用户名查找用户
    user, err := s.userRepo.FindByUsername(req.Username)
    if err != nil {
        return nil, &apperror.AppError{
            Code:    apperror.ErrUnauthorized,
            Message: "用户名或密码错误", // 不提示具体哪个错
        }
    }

    // 2. 验证密码
    if !utils.CheckPassword(req.Password, user.Password) {
        return nil, &apperror.AppError{
            Code:    apperror.ErrUnauthorized,
            Message: "用户名或密码错误",
        }
    }

    // 3. 检查用户状态
    if user.Status != 1 {
        return nil, &apperror.AppError{
            Code:    apperror.ErrForbidden,
            Message: "账号已被禁用",
        }
    }

    // 4. 生成 Token
    accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
        user.ID, user.Username, user.Role,
    )
    if err != nil {
        return nil, apperror.Internal(err)
    }

    return &dto.TokenResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    7200, // 2 小时
        User:         *dto.ToUserResponse(user),
    }, nil
}
```

### 任务 4：Token 刷新接口（难度：中等）

```go
// POST /api/v1/auth/refresh
type RefreshTokenReq struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
    var req dto.RefreshTokenReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err)
        return
    }

    // 1. 解析 refresh_token
    // 2. 检查是否过期
    // 3. 查询用户（确保用户还存在且未被禁用）
    // 4. 生成新的 token pair
    // 5. 返回新的 access_token 和 refresh_token
}
```

### 任务 5：完整路由整合测试

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"Pass1234","confirm_password":"Pass1234"}'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"Pass1234"}'
# 记录返回的 access_token

# 访问需要认证的接口
curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <access_token>"

# 不带 Token 访问 -> 401
curl http://localhost:8080/api/v1/users/profile

# 带过期 Token 访问 -> 401 (token expired)

# 刷新 Token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

---

## 自检清单

- [ ] 密码使用 bcrypt 加密存储
- [ ] JWT 生成和解析正确
- [ ] 认证中间件正确解析 Bearer Token
- [ ] Token 过期返回 401 和明确的错误信息
- [ ] 登录失败不提示"用户不存在"或"密码错误"的具体信息（安全）
- [ ] Refresh Token 机制可以正常刷新 Access Token
- [ ] 权限中间件能区分 admin 和普通用户
