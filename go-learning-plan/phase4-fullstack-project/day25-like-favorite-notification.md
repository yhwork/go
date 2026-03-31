# Day 25 - 点赞、收藏与通知

## 今日目标

实现文章点赞/取消点赞（Redis 去重）、收藏功能、WebSocket 实时通知（新评论通知）、通知列表管理。

---

## 练习任务

### 任务 1：文章点赞（Redis + DB）

```go
type LikeService interface {
    Like(userID, articleID uint) error       // 点赞
    Unlike(userID, articleID uint) error     // 取消点赞
    IsLiked(userID, articleID uint) (bool, error)
    GetLikeCount(articleID uint) (int64, error)
    GetUserLikedArticles(userID uint, page, pageSize int) (*dto.PageResult, error)
}

// 使用 Redis Set 实现高性能点赞
// Key: "article:likes:{articleID}"
// Value: Set 中存储 userID

func (s *likeService) Like(userID, articleID uint) error {
    ctx := context.Background()
    key := fmt.Sprintf("article:likes:%d", articleID)

    // 1. Redis 中添加
    added, err := s.rdb.SAdd(ctx, key, userID).Result()
    if err != nil {
        return err
    }
    if added == 0 {
        return nil // 已经点过赞了
    }

    // 2. 数据库记录（异步或同步）
    s.db.Create(&model.ArticleLike{ArticleID: articleID, UserID: userID})

    // 3. 更新文章点赞计数
    s.db.Model(&model.Article{}).Where("id = ?", articleID).
        Update("like_count", gorm.Expr("like_count + 1"))

    return nil
}
```

### 任务 2：文章收藏

```go
type FavoriteService interface {
    Favorite(userID, articleID uint) error
    Unfavorite(userID, articleID uint) error
    IsFavorited(userID, articleID uint) (bool, error)
    GetUserFavorites(userID uint, page, pageSize int) (*dto.PageResult, error)
}

// 实现类似点赞，使用 DB 唯一索引防重复
```

### 任务 3：通知系统

```go
// internal/model/notification.go
type Notification struct {
    gorm.Model
    UserID    uint   `gorm:"index;not null"`        // 接收者
    Type      string `gorm:"size:20;not null;index"` // comment, like, system
    Title     string `gorm:"size:200;not null"`
    Content   string `gorm:"size:500"`
    RelatedID uint   `gorm:"index"`                  // 关联的文章/评论 ID
    IsRead    bool   `gorm:"default:false;index"`
}

type NotificationService interface {
    Create(userID uint, nType, title, content string, relatedID uint) error
    List(userID uint, page, pageSize int, onlyUnread bool) (*dto.PageResult, error)
    MarkAsRead(userID uint, notificationID uint) error
    MarkAllAsRead(userID uint) error
    UnreadCount(userID uint) (int64, error)
}

// 创建评论时发送通知
func (s *commentService) Create(userID uint, req *dto.CreateCommentReq) (*dto.CommentResponse, error) {
    // ... 创建评论 ...

    // 通知文章作者
    article, _ := s.articleRepo.FindByID(req.ArticleID)
    if article.AuthorID != userID {
        s.notificationService.Create(
            article.AuthorID,
            "comment",
            fmt.Sprintf("%s 评论了你的文章", username),
            comment.Content[:min(len(comment.Content), 100)],
            article.ID,
        )
        // 通过 WebSocket 推送实时通知
        s.wsHub.SendToUser(article.AuthorID, wsMessage)
    }

    return resp, nil
}
```

### 任务 4：WebSocket 实时推送通知

```go
// 扩展 Day 17 的 Hub，添加定向推送能力

func (h *Hub) SendToUser(userID uint, msg *Message) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    if client, ok := h.clients[userID]; ok {
        data, _ := json.Marshal(msg)
        select {
        case client.Send <- data:
        default:
            // channel 满了，丢弃
        }
    }
}

// 通知消息类型
type NotificationMessage struct {
    Type    string `json:"type"`    // "notification"
    Title   string `json:"title"`
    Content string `json:"content"`
    Unread  int64  `json:"unread"`  // 当前未读数
}
```

### 任务 5：路由注册

```go
// 点赞
likes := v1.Group("/likes")
likes.Use(middleware.JWTAuth())
{
    likes.POST("/article/:id", likeHandler.Like)
    likes.DELETE("/article/:id", likeHandler.Unlike)
    likes.GET("/article/:id", likeHandler.IsLiked) // 当前用户是否已点赞
}

// 收藏
favorites := v1.Group("/favorites")
favorites.Use(middleware.JWTAuth())
{
    favorites.POST("/article/:id", favoriteHandler.Favorite)
    favorites.DELETE("/article/:id", favoriteHandler.Unfavorite)
    favorites.GET("", favoriteHandler.List) // 我的收藏列表
}

// 通知
notifications := v1.Group("/notifications")
notifications.Use(middleware.JWTAuth())
{
    notifications.GET("", notificationHandler.List)
    notifications.GET("/unread-count", notificationHandler.UnreadCount)
    notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
    notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
}
```

---

## 自检清单

- [ ] 点赞/取消点赞正常，不能重复点赞
- [ ] 点赞数与实际数一致
- [ ] 收藏功能正常工作
- [ ] 评论时作者收到通知
- [ ] WebSocket 实时推送通知正常
- [ ] 通知已读/未读状态管理正确
- [ ] 未读通知计数准确
