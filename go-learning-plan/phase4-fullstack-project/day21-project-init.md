# Day 21 - 博客项目初始化与数据库设计

## 今日目标

从零开始搭建个人博客系统项目，完成需求分析、数据库表设计、项目骨架初始化、数据库迁移。

---

## 需求分析

### 核心功能模块

| 模块 | 功能 |
|------|------|
| 用户 | 注册、登录、个人信息管理、角色权限（admin/user） |
| 文章 | CRUD、Markdown 内容、草稿/发布状态、阅读量统计 |
| 分类 | CRUD、树形结构（支持子分类） |
| 标签 | CRUD、文章-标签多对多关联 |
| 评论 | 嵌套评论、分页加载 |
| 点赞/收藏 | 文章点赞、取消点赞、收藏夹 |

### 数据库 ER 图（文字描述）

```
users
  ├── 1:N articles（一个用户可以发布多篇文章）
  ├── 1:N comments（一个用户可以发表多条评论）
  ├── M:N articles（点赞关系，通过 article_likes 中间表）
  └── M:N articles（收藏关系，通过 article_favorites 中间表）

articles
  ├── N:1 categories（一篇文章属于一个分类）
  ├── M:N tags（通过 article_tags 中间表）
  └── 1:N comments

categories
  └── 自引用（parent_id 实现树形结构）
```

---

## 练习任务

### 任务 1：创建项目并初始化

```bash
mkdir go-blog && cd go-blog
go mod init github.com/yourname/go-blog
```

按照 Day 05 学到的标准项目结构创建所有目录。

### 任务 2：定义完整数据模型

```go
// internal/model/user.go
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;size:50;not null"`
    Email    string `gorm:"uniqueIndex;size:100;not null"`
    Password string `gorm:"size:100;not null"`
    Nickname string `gorm:"size:50"`
    Avatar   string `gorm:"size:255"`
    Bio      string `gorm:"size:500"`
    Role     string `gorm:"size:20;default:user;index"`
    Status   int8   `gorm:"default:1"`
}

// internal/model/article.go
type Article struct {
    gorm.Model
    Title      string `gorm:"size:200;not null;index"`
    Content    string `gorm:"type:longtext;not null"`
    Summary    string `gorm:"size:500"`
    Cover      string `gorm:"size:255"`
    AuthorID   uint   `gorm:"index;not null"`
    CategoryID uint   `gorm:"index"`
    Status     string `gorm:"size:20;default:draft;index"` // draft, published
    Views      int    `gorm:"default:0"`
    LikeCount  int    `gorm:"default:0"`

    Author   User     `gorm:"foreignKey:AuthorID"`
    Category Category `gorm:"foreignKey:CategoryID"`
    Tags     []Tag    `gorm:"many2many:article_tags"`
    Comments []Comment
}

// internal/model/category.go
type Category struct {
    gorm.Model
    Name        string     `gorm:"size:50;not null"`
    Slug        string     `gorm:"size:50;uniqueIndex"`
    Description string     `gorm:"size:200"`
    ParentID    *uint      `gorm:"index"`
    Sort        int        `gorm:"default:0"`
    Parent      *Category  `gorm:"foreignKey:ParentID"`
    Children    []Category `gorm:"foreignKey:ParentID"`
}

// internal/model/tag.go
type Tag struct {
    gorm.Model
    Name     string    `gorm:"size:30;not null;uniqueIndex"`
    Articles []Article `gorm:"many2many:article_tags"`
}

// internal/model/comment.go
type Comment struct {
    gorm.Model
    Content   string `gorm:"type:text;not null"`
    ArticleID uint   `gorm:"index;not null"`
    UserID    uint   `gorm:"index;not null"`
    ParentID  *uint  `gorm:"index"` // 父评论 ID，支持嵌套

    User     User      `gorm:"foreignKey:UserID"`
    Article  Article   `gorm:"foreignKey:ArticleID"`
    Parent   *Comment  `gorm:"foreignKey:ParentID"`
    Replies  []Comment `gorm:"foreignKey:ParentID"`
}

// internal/model/article_like.go
type ArticleLike struct {
    gorm.Model
    ArticleID uint `gorm:"uniqueIndex:idx_article_user;not null"`
    UserID    uint `gorm:"uniqueIndex:idx_article_user;not null"`
}

// internal/model/article_favorite.go
type ArticleFavorite struct {
    gorm.Model
    ArticleID uint `gorm:"uniqueIndex:idx_article_user_fav;not null"`
    UserID    uint `gorm:"uniqueIndex:idx_article_user_fav;not null"`
}
```

### 任务 3：初始化数据库连接和迁移

将之前所有学到的内容整合：配置加载、数据库连接、自动迁移、种子数据。

### 任务 4：编写种子数据

```go
func Seed(db *gorm.DB) error {
    // 创建管理员
    // 创建 3-5 个分类
    // 创建 10 个标签
    // 创建 5 个测试用户
    // 创建 20 篇测试文章（随机分配作者、分类、标签）
    // 创建一些测试评论
}
```

### 任务 5：验证

启动项目，确认：
- 数据库表全部正确创建
- 种子数据成功插入
- 关联关系正确（外键、中间表）

---

## 自检清单

- [ ] 项目目录结构完整
- [ ] 所有模型定义正确，GORM 标签无误
- [ ] 数据库迁移成功，表结构符合预期
- [ ] 种子数据成功插入
- [ ] 关联关系测试通过（查询用户的文章、文章的标签等）
