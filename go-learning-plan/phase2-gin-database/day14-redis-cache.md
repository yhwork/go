# Day 14 - Redis 缓存

## 今日目标

掌握 Go Redis 客户端的使用，实现缓存策略（Cache-Aside）、基于 Redis 的接口限流、JWT 黑名单（支持退出登录）。

---

## 知识点

### 1. go-redis 连接

```go
import "github.com/redis/go-redis/v9"

rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    PoolSize: 10,
})

// 测试连接
ctx := context.Background()
_, err := rdb.Ping(ctx).Result()
```

### 2. 基本操作

```go
// String
rdb.Set(ctx, "key", "value", 10*time.Minute)
val, err := rdb.Get(ctx, "key").Result()

// Hash
rdb.HSet(ctx, "user:1", "name", "alice", "email", "alice@example.com")
rdb.HGetAll(ctx, "user:1")

// List
rdb.LPush(ctx, "queue", "task1", "task2")
rdb.RPop(ctx, "queue")

// Set
rdb.SAdd(ctx, "likes:article:1", 1, 2, 3)
rdb.SIsMember(ctx, "likes:article:1", 1)

// Sorted Set
rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "alice"})
```

---

## 练习任务

### 任务 1：Redis 缓存封装（难度：中等）

```go
// internal/cache/cache.go
package cache

type Cache struct {
    rdb *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
    return &Cache{rdb: rdb}
}

// Get 获取缓存，反序列化到 dest
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return ErrCacheMiss
    }
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

// Set 设置缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Delete 删除缓存
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
    return c.rdb.Del(ctx, keys...).Err()
}

// DeleteByPrefix 按前缀删除
func (c *Cache) DeleteByPrefix(ctx context.Context, prefix string) error {
    iter := c.rdb.Scan(ctx, 0, prefix+"*", 100).Iterator()
    for iter.Next(ctx) {
        c.rdb.Del(ctx, iter.Val())
    }
    return iter.Err()
}

var ErrCacheMiss = errors.New("cache miss")
```

### 任务 2：Cache-Aside 模式（难度：中等偏上）

```go
// 在 Service 层实现 Cache-Aside

const (
    userCachePrefix = "user:"
    userCacheTTL    = 30 * time.Minute
)

func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
    ctx := context.Background()
    cacheKey := fmt.Sprintf("%s%d", userCachePrefix, userID)

    // 1. 先查缓存
    var cached dto.UserResponse
    err := s.cache.Get(ctx, cacheKey, &cached)
    if err == nil {
        return &cached, nil // 缓存命中
    }

    // 2. 缓存未命中，查数据库
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }

    resp := dto.ToUserResponse(user)

    // 3. 写入缓存
    _ = s.cache.Set(ctx, cacheKey, resp, userCacheTTL)

    return resp, nil
}

// 更新时删除缓存
func (s *userService) UpdateProfile(userID uint, req *dto.UpdateProfileReq) (*dto.UserResponse, error) {
    // 1. 更新数据库
    // 2. 删除缓存
    ctx := context.Background()
    cacheKey := fmt.Sprintf("%s%d", userCachePrefix, userID)
    s.cache.Delete(ctx, cacheKey)
    // 3. 返回新数据
}
```

### 任务 3：基于 Redis 的限流（难度：中等）

```go
// 滑动窗口限流

func (c *Cache) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    now := time.Now().UnixMilli()
    windowStart := now - window.Milliseconds()

    pipe := c.rdb.Pipeline()

    // 移除窗口外的记录
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
    // 统计窗口内的请求数
    countCmd := pipe.ZCard(ctx, key)
    // 添加当前请求
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
    // 设置过期时间
    pipe.Expire(ctx, key, window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }

    return countCmd.Val() < int64(limit), nil
}

// 中间件使用
func RedisRateLimiter(cache *cache.Cache, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("ratelimit:%s", c.ClientIP())
        allowed, err := cache.RateLimit(c.Request.Context(), key, limit, window)
        if err != nil || !allowed {
            c.AbortWithStatusJSON(429, gin.H{
                "code":    429,
                "message": "请求过于频繁，请稍后重试",
            })
            return
        }
        c.Next()
    }
}
```

### 任务 4：JWT 黑名单（退出登录）（难度：中等）

```go
// 用 Redis 存储已注销的 Token

const tokenBlacklistPrefix = "token:blacklist:"

// AddToBlacklist 将 Token 加入黑名单
func (c *Cache) AddToBlacklist(ctx context.Context, tokenID string, expiry time.Duration) error {
    key := tokenBlacklistPrefix + tokenID
    return c.rdb.Set(ctx, key, "1", expiry).Err()
}

// IsBlacklisted 检查 Token 是否在黑名单中
func (c *Cache) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
    key := tokenBlacklistPrefix + tokenID
    _, err := c.rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return true, nil
}

// 退出登录接口
func (s *userService) Logout(accessToken string) error {
    claims, err := s.jwtManager.ParseToken(accessToken)
    if err != nil {
        return nil // Token 无效也算退出成功
    }

    // 将 Token 加入黑名单，过期时间设为 Token 的剩余有效期
    remaining := time.Until(claims.ExpiresAt.Time)
    if remaining > 0 {
        return s.cache.AddToBlacklist(context.Background(), accessToken, remaining)
    }
    return nil
}

// 在认证中间件中检查黑名单
func JWTAuth(jwtManager *jwt.JWTManager, cache *cache.Cache) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... 解析 Token ...

        // 检查黑名单
        blacklisted, _ := cache.IsBlacklisted(c.Request.Context(), tokenString)
        if blacklisted {
            c.AbortWithStatusJSON(401, response.Error(40101, "Token 已失效，请重新登录"))
            return
        }

        // ... 设置用户信息 ...
        c.Next()
    }
}
```

---

## 测试验证

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"Pass1234"}'

# 用 Token 访问 -> 200
curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <token>"

# 退出登录
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <token>"

# 再用同一个 Token 访问 -> 401 (Token 已失效)
curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <token>"

# 快速请求测试限流（连续发 20 次）
for i in $(seq 1 20); do curl -s http://localhost:8080/api/v1/users; done
```

---

## 自检清单

- [ ] Redis 连接和基础操作正常
- [ ] 缓存封装支持序列化/反序列化任意类型
- [ ] Cache-Aside 模式正确实现（先查缓存、miss 后查 DB、写入缓存）
- [ ] 更新数据时正确删除缓存
- [ ] 限流功能在请求超限时返回 429
- [ ] JWT 黑名单机制生效，退出登录后 Token 不可用
- [ ] 黑名单过期时间与 Token 剩余有效期一致（不浪费内存）
