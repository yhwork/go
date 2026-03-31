# Day 24 - 分类、标签、评论模块

## 今日目标

实现分类 CRUD（树形结构）、标签 CRUD、评论系统（嵌套回复+分页）。

---

## 练习任务

### 任务 1：分类模块（树形结构）

```go
// Service 接口
type CategoryService interface {
    Create(req *dto.CreateCategoryReq) (*dto.CategoryResponse, error)
    Update(id uint, req *dto.UpdateCategoryReq) (*dto.CategoryResponse, error)
    Delete(id uint) error    // 有子分类或文章时不允许删除
    GetByID(id uint) (*dto.CategoryResponse, error)
    GetTree() ([]dto.CategoryTreeNode, error) // 返回树形结构
    List() ([]dto.CategoryResponse, error)    // 返回扁平列表
}

// 树形结构 DTO
type CategoryTreeNode struct {
    ID       uint               `json:"id"`
    Name     string             `json:"name"`
    Slug     string             `json:"slug"`
    Sort     int                `json:"sort"`
    Children []CategoryTreeNode `json:"children"`
}

// 构建树的递归函数
func buildCategoryTree(categories []model.Category, parentID *uint) []dto.CategoryTreeNode {
    var tree []dto.CategoryTreeNode
    for _, cat := range categories {
        if (cat.ParentID == nil && parentID == nil) ||
           (cat.ParentID != nil && parentID != nil && *cat.ParentID == *parentID) {
            node := dto.CategoryTreeNode{
                ID:   cat.ID,
                Name: cat.Name,
                Slug: cat.Slug,
                Sort: cat.Sort,
            }
            node.Children = buildCategoryTree(categories, &cat.ID)
            tree = append(tree, node)
        }
    }
    return tree
}
```

### 任务 2：标签模块

```go
type TagService interface {
    Create(req *dto.CreateTagReq) (*dto.TagResponse, error)
    Update(id uint, req *dto.UpdateTagReq) (*dto.TagResponse, error)
    Delete(id uint) error
    List() ([]dto.TagResponse, error)
    GetPopular(limit int) ([]dto.TagWithCount, error) // 热门标签（按文章数排序）
}

// 热门标签查询
func (r *tagRepo) GetPopular(limit int) ([]dto.TagWithCount, error) {
    var results []dto.TagWithCount
    err := r.db.Model(&model.Tag{}).
        Select("tags.id, tags.name, COUNT(article_tags.article_id) as article_count").
        Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
        Group("tags.id").
        Order("article_count DESC").
        Limit(limit).
        Find(&results).Error
    return results, err
}
```

### 任务 3：评论系统（嵌套）

```go
type CommentService interface {
    Create(userID uint, req *dto.CreateCommentReq) (*dto.CommentResponse, error)
    Delete(userID uint, commentID uint) error
    ListByArticle(articleID uint, page, pageSize int) (*dto.PageResult, error)
}

type CreateCommentReq struct {
    ArticleID uint   `json:"article_id" binding:"required"`
    Content   string `json:"content" binding:"required,min=1,max=1000"`
    ParentID  *uint  `json:"parent_id"` // 回复哪条评论
}

type CommentResponse struct {
    ID        uint              `json:"id"`
    Content   string            `json:"content"`
    User      UserBriefResp     `json:"user"`
    ParentID  *uint             `json:"parent_id"`
    Replies   []CommentResponse `json:"replies,omitempty"`
    CreatedAt time.Time         `json:"created_at"`
}

// 查询评论（只查顶级评论，子评论通过 Preload 加载）
func (r *commentRepo) ListByArticle(articleID uint, page, pageSize int) ([]model.Comment, int64, error) {
    var comments []model.Comment
    var total int64

    db := r.db.Model(&model.Comment{}).
        Where("article_id = ? AND parent_id IS NULL", articleID)

    db.Count(&total)

    err := db.Offset((page-1)*pageSize).Limit(pageSize).
        Order("created_at DESC").
        Preload("User", selectFields("id", "username", "avatar")).
        Preload("Replies", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at ASC").
                Preload("User", selectFields("id", "username", "avatar"))
        }).
        Find(&comments).Error

    return comments, total, err
}
```

### 任务 4：路由注册

```go
// 分类（admin 才能增删改）
categories := v1.Group("/categories")
{
    categories.GET("", categoryHandler.List)
    categories.GET("/tree", categoryHandler.GetTree)
    categories.GET("/:id", categoryHandler.GetByID)
}
adminCategories := v1.Group("/admin/categories")
adminCategories.Use(middleware.JWTAuth(), middleware.RequireRole("admin"))
{
    adminCategories.POST("", categoryHandler.Create)
    adminCategories.PUT("/:id", categoryHandler.Update)
    adminCategories.DELETE("/:id", categoryHandler.Delete)
}

// 标签
tags := v1.Group("/tags")
{
    tags.GET("", tagHandler.List)
    tags.GET("/popular", tagHandler.GetPopular)
}

// 评论
comments := v1.Group("/comments")
{
    comments.GET("/article/:article_id", commentHandler.ListByArticle)
}
authComments := v1.Group("/comments")
authComments.Use(middleware.JWTAuth())
{
    authComments.POST("", commentHandler.Create)
    authComments.DELETE("/:id", commentHandler.Delete)
}
```

---

## 自检清单

- [ ] 分类树形结构正确展示
- [ ] 有文章的分类不能删除
- [ ] 热门标签按文章数排序
- [ ] 评论嵌套回复正确加载
- [ ] 只能删除自己的评论
- [ ] 评论分页正常工作
