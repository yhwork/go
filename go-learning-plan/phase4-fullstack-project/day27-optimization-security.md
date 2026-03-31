# Day 27 - 性能优化与安全加固

## 今日目标

对博客系统进行性能优化（Redis 缓存、数据库索引）和安全加固（SQL 注入防护、XSS 防护、CSRF、限流）。

---

## 练习任务

### 任务 1：Redis 缓存优化

```go
// 1. 热门文章列表缓存
const hotArticlesCacheKey = "articles:hot"
const hotArticlesCacheTTL = 5 * time.Minute

func (s *articleService) GetHotArticles(limit int) ([]dto.ArticleResponse, error) {
    ctx := context.Background()

    // 先查缓存
    var cached []dto.ArticleResponse
    if err := s.cache.Get(ctx, hotArticlesCacheKey, &cached); err == nil {
        return cached, nil
    }

    // 查数据库
    articles, err := s.articleRepo.GetHotArticles(limit)
    if err != nil {
        return nil, err
    }

    result := dto.ToArticleResponseList(articles)
    s.cache.Set(ctx, hotArticlesCacheKey, result, hotArticlesCacheTTL)
    return result, nil
}

// 2. 文章详情缓存
// 3. 分类列表缓存（变更少，TTL 可以设长一些）
// 4. 用户信息缓存
```

### 任务 2：数据库优化

```sql
-- 检查慢查询
-- MySQL 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1; -- 超过 1 秒记录

-- 分析查询执行计划
EXPLAIN SELECT * FROM articles WHERE category_id = 1 AND status = 'published' ORDER BY created_at DESC LIMIT 10;

-- 添加复合索引
CREATE INDEX idx_articles_category_status_created ON articles(category_id, status, created_at DESC);

-- 检查索引使用情况
SHOW INDEX FROM articles;
```

```go
// 在 GORM 模型中添加复合索引
type Article struct {
    // ...
    CategoryID uint   `gorm:"index:idx_cat_status_created"`
    Status     string `gorm:"index:idx_cat_status_created"`
    CreatedAt  time.Time `gorm:"index:idx_cat_status_created"`
}
```

### 任务 3：安全加固

```go
// 1. SQL 注入防护 - 确认所有查询使用参数化
// 正确 ✓
db.Where("username = ?", username).First(&user)
// 危险 ✗ - 永远不要这样做
db.Where(fmt.Sprintf("username = '%s'", username)).First(&user)

// 2. XSS 防护 - 使用 bluemonday 过滤 HTML
import "github.com/microcosm-cc/bluemonday"

func SanitizeHTML(input string) string {
    p := bluemonday.UGCPolicy() // 允许常见的用户生成内容标签
    return p.Sanitize(input)
}

// 在创建/更新文章和评论时过滤：
func (s *articleService) Create(userID uint, req *dto.CreateArticleReq) (*dto.ArticleResponse, error) {
    req.Content = SanitizeHTML(req.Content)
    req.Title = bluemonday.StrictPolicy().Sanitize(req.Title)
    // ...
}

// 3. 接口限流（分级别）
// 全局：100 QPS
// 登录接口：5 次/分钟（防暴力破解）
// 注册接口：3 次/分钟
// 上传接口：10 次/分钟

loginGroup := auth.Group("")
loginGroup.Use(middleware.RedisRateLimiter(cache, 5, time.Minute))
{
    loginGroup.POST("/login", userHandler.Login)
}

// 4. 请求体大小限制
r.Use(gin.Recovery())
r.MaxMultipartMemory = 8 << 20 // 8MB

// 5. 安全响应头
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

### 任务 4：安全审计检查清单

逐项检查：

- [ ] 所有数据库查询使用参数化（无 SQL 注入风险）
- [ ] 用户输入的 HTML 经过 bluemonday 过滤（无 XSS 风险）
- [ ] 密码使用 bcrypt 加密存储
- [ ] JWT Secret 不硬编码在代码中
- [ ] 配置文件不包含生产环境密码
- [ ] 敏感接口有限流保护
- [ ] 登录失败不提示具体原因
- [ ] 文件上传有类型和大小限制
- [ ] 响应中不返回密码等敏感字段（json:"-"）
- [ ] 安全响应头配置正确

---

## 自检清单

- [ ] 热门文章缓存减少数据库查询
- [ ] 数据库添加合适的索引
- [ ] 慢查询优化到 < 100ms
- [ ] 安全审计清单全部通过
- [ ] 限流配置正确
- [ ] XSS 防护生效
