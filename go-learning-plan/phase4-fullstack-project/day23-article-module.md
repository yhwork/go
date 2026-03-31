# Day 23 - 文章模块完整实现

## 今日目标

完成文章模块：CRUD、Markdown/富文本内容存储、分页+筛选+搜索、阅读量统计、草稿/发布状态管理。

---

## 练习任务

### 任务 1：文章 DTO 定义

```go
// internal/dto/article_dto.go

type CreateArticleReq struct {
    Title      string   `json:"title" binding:"required,min=1,max=200"`
    Content    string   `json:"content" binding:"required,min=10"`
    Summary    string   `json:"summary" binding:"omitempty,max=500"`
    Cover      string   `json:"cover" binding:"omitempty,url"`
    CategoryID uint     `json:"category_id" binding:"required"`
    TagIDs     []uint   `json:"tag_ids" binding:"omitempty"`
    Status     string   `json:"status" binding:"required,oneof=draft published"`
}

type UpdateArticleReq struct {
    Title      *string  `json:"title" binding:"omitempty,min=1,max=200"`
    Content    *string  `json:"content" binding:"omitempty,min=10"`
    Summary    *string  `json:"summary"`
    Cover      *string  `json:"cover"`
    CategoryID *uint    `json:"category_id"`
    TagIDs     *[]uint  `json:"tag_ids"`
    Status     *string  `json:"status" binding:"omitempty,oneof=draft published"`
}

type ArticleQuery struct {
    Page       int    `form:"page"`
    PageSize   int    `form:"page_size"`
    Keyword    string `form:"keyword"`
    CategoryID uint   `form:"category_id"`
    TagID      uint   `form:"tag_id"`
    AuthorID   uint   `form:"author_id"`
    Status     string `form:"status"`
    SortBy     string `form:"sort_by"`    // created_at, views, likes
    SortOrder  string `form:"sort_order"` // asc, desc
}

type ArticleResponse struct {
    ID         uint             `json:"id"`
    Title      string           `json:"title"`
    Summary    string           `json:"summary"`
    Cover      string           `json:"cover"`
    Status     string           `json:"status"`
    Views      int              `json:"views"`
    LikeCount  int              `json:"like_count"`
    Author     UserBriefResp    `json:"author"`
    Category   CategoryResp     `json:"category"`
    Tags       []TagResp        `json:"tags"`
    CreatedAt  time.Time        `json:"created_at"`
    UpdatedAt  time.Time        `json:"updated_at"`
}

type ArticleDetailResponse struct {
    ArticleResponse
    Content string `json:"content"` // 详情才返回完整内容
}
```

### 任务 2：文章 Repository

```go
type ArticleRepository interface {
    Create(article *model.Article) error
    FindByID(id uint) (*model.Article, error)
    Update(article *model.Article) error
    Delete(id uint) error
    List(query dto.ArticleQuery) ([]model.Article, int64, error)
    IncrementViews(id uint) error
    ReplaceTags(articleID uint, tagIDs []uint) error
}

// 关键实现：复杂列表查询
func (r *articleRepo) List(query dto.ArticleQuery) ([]model.Article, int64, error) {
    var articles []model.Article
    var total int64

    db := r.db.Model(&model.Article{})

    // 动态条件构建
    if query.Keyword != "" {
        db = db.Where("title LIKE ? OR summary LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
    }
    if query.CategoryID > 0 {
        db = db.Where("category_id = ?", query.CategoryID)
    }
    if query.AuthorID > 0 {
        db = db.Where("author_id = ?", query.AuthorID)
    }
    if query.Status != "" {
        db = db.Where("status = ?", query.Status)
    }
    if query.TagID > 0 {
        db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id AND article_tags.tag_id = ?", query.TagID)
    }

    db.Count(&total)

    // 排序 + 分页 + 预加载
    db.Scopes(orderScope(query.SortBy, query.SortOrder), paginateScope(query.Page, query.PageSize)).
        Preload("Author", selectFields("id", "username", "avatar")).
        Preload("Category").
        Preload("Tags").
        Find(&articles)

    return articles, total, db.Error
}
```

### 任务 3：文章 Service

```go
type ArticleService interface {
    Create(userID uint, req *dto.CreateArticleReq) (*dto.ArticleResponse, error)
    GetByID(id uint) (*dto.ArticleDetailResponse, error) // 同时增加阅读量
    Update(userID uint, id uint, req *dto.UpdateArticleReq) (*dto.ArticleResponse, error)
    Delete(userID uint, id uint) error                    // 只能删除自己的文章（admin 除外）
    List(query *dto.ArticleQuery) (*dto.PageResult, error)
    Publish(userID uint, id uint) error                   // 草稿 -> 发布
    Unpublish(userID uint, id uint) error                 // 发布 -> 草稿
}

// GetByID 获取文章详情 + 增加阅读量
func (s *articleService) GetByID(id uint) (*dto.ArticleDetailResponse, error) {
    article, err := s.articleRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // 异步增加阅读量（不影响响应速度）
    go s.articleRepo.IncrementViews(id)

    return dto.ToArticleDetailResponse(article), nil
}
```

### 任务 4：文章 Handler 和路由

```go
articles := v1.Group("/articles")
{
    articles.GET("", articleHandler.List)           // 公开：文章列表
    articles.GET("/:id", articleHandler.GetByID)    // 公开：文章详情
}

// 需要认证的操作
authArticles := v1.Group("/articles")
authArticles.Use(middleware.JWTAuth())
{
    authArticles.POST("", articleHandler.Create)
    authArticles.PUT("/:id", articleHandler.Update)
    authArticles.DELETE("/:id", articleHandler.Delete)
    authArticles.POST("/:id/publish", articleHandler.Publish)
    authArticles.POST("/:id/unpublish", articleHandler.Unpublish)
}
```

### 任务 5：测试所有接口

验证：
- 创建文章（带标签）
- 获取文章列表（分页、搜索、分类筛选、标签筛选）
- 获取文章详情（阅读量是否增加）
- 更新文章（含修改标签）
- 删除文章（只能删自己的）
- 发布/取消发布

---

## 自检清单

- [ ] 文章 CRUD 完整实现
- [ ] 文章列表支持分页、搜索、多条件筛选
- [ ] 文章-标签多对多关系正确维护
- [ ] 阅读量统计正确
- [ ] 权限控制正确（只能编辑/删除自己的文章）
- [ ] 草稿和发布状态切换正常
