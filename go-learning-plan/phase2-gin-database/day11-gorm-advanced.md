# Day 11 - GORM 进阶查询

## 今日目标

掌握 GORM 的预加载（Preload）、关联查询、事务处理、Scope（查询作用域）、以及通用分页封装。

---

## 知识点

### 1. 预加载（Preload）

```go
// 查询用户时同时加载其订单
var user User
db.Preload("Orders").First(&user, 1)

// 嵌套预加载：用户 -> 订单 -> 订单项 -> 商品
db.Preload("Orders.OrderItems.Product").First(&user, 1)

// 带条件的预加载
db.Preload("Orders", "status = ?", "completed").First(&user, 1)

// 自定义预加载
db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(5)
}).First(&user, 1)
```

### 2. 关联查询（Joins）

```go
// Join 查询（比 Preload 性能更好，但只适合一对一/多对一）
var articles []Article
db.Joins("Author").Joins("Category").
    Where("articles.status = ?", "published").
    Find(&articles)

// 原生 SQL Join
db.Table("articles").
    Select("articles.*, users.username as author_name").
    Joins("LEFT JOIN users ON users.id = articles.author_id").
    Find(&results)
```

### 3. 事务

```go
// 自动事务（推荐）
err := db.Transaction(func(tx *gorm.DB) error {
    // 在事务中进行所有操作，使用 tx 而不是 db
    if err := tx.Create(&order).Error; err != nil {
        return err // 返回任何错误都会回滚
    }
    if err := tx.Create(&orderItems).Error; err != nil {
        return err
    }
    // 扣减库存
    if err := tx.Model(&Product{}).Where("id = ?", productID).
        Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
        return err
    }
    return nil // 返回 nil 提交事务
})
```

### 4. Scope（查询作用域）

```go
// Scope 是可复用的查询条件

func Published(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "published")
}

func ByCategory(categoryID uint) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("category_id = ?", categoryID)
    }
}

func OrderByLatest(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC")
}

// 使用
db.Scopes(Published, ByCategory(1), OrderByLatest).Find(&articles)
```

---

## 练习任务

### 任务 1：通用分页封装（难度：中等）

```go
// pkg/pagination/pagination.go

type Pagination struct {
    Page     int   `json:"page" form:"page"`
    PageSize int   `json:"page_size" form:"page_size"`
    Total    int64 `json:"total"`
}

func (p *Pagination) GetOffset() int {
    return (p.GetPage() - 1) * p.GetPageSize()
}

func (p *Pagination) GetPage() int {
    if p.Page <= 0 {
        p.Page = 1
    }
    return p.Page
}

func (p *Pagination) GetPageSize() int {
    if p.PageSize <= 0 {
        p.PageSize = 10
    }
    if p.PageSize > 100 {
        p.PageSize = 100
    }
    return p.PageSize
}

func (p *Pagination) GetTotalPages() int {
    totalPages := int(p.Total) / p.GetPageSize()
    if int(p.Total)%p.GetPageSize() > 0 {
        totalPages++
    }
    return totalPages
}

// Paginate 是一个 GORM Scope
func Paginate(p *Pagination) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        // 先统计总数
        var total int64
        db.Count(&total)
        p.Total = total

        return db.Offset(p.GetOffset()).Limit(p.GetPageSize())
    }
}

// 使用示例：
// var articles []Article
// page := Pagination{Page: 1, PageSize: 10}
// db.Scopes(Paginate(&page)).Find(&articles)
// 此时 page.Total 已经被填充
```

### 任务 2：多条件筛选查询（难度：中等偏上）

```go
// 实现文章的复杂查询

type ArticleQuery struct {
    Pagination
    Keyword    string `form:"keyword"`
    CategoryID uint   `form:"category_id"`
    TagID      uint   `form:"tag_id"`
    AuthorID   uint   `form:"author_id"`
    Status     string `form:"status"`
    SortBy     string `form:"sort_by"`    // created_at, views, likes
    SortOrder  string `form:"sort_order"` // asc, desc
}

func (r *articleRepository) List(query ArticleQuery) ([]model.Article, *Pagination, error) {
    var articles []model.Article
    page := &query.Pagination

    db := r.db.Model(&model.Article{})

    // 条件筛选（只在有值时添加条件）
    if query.Keyword != "" {
        db = db.Where("title LIKE ? OR content LIKE ?",
            "%"+query.Keyword+"%", "%"+query.Keyword+"%")
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
        db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
            Where("article_tags.tag_id = ?", query.TagID)
    }

    // 排序
    sortBy := "created_at"
    if query.SortBy != "" {
        sortBy = query.SortBy
    }
    sortOrder := "DESC"
    if query.SortOrder == "asc" {
        sortOrder = "ASC"
    }
    db = db.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

    // 统计总数
    var total int64
    if err := db.Count(&total).Error; err != nil {
        return nil, nil, err
    }
    page.Total = total

    // 分页查询 + 预加载
    err := db.Offset(page.GetOffset()).Limit(page.GetPageSize()).
        Preload("Author", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, username, avatar") // 只查需要的字段
        }).
        Preload("Category").
        Preload("Tags").
        Find(&articles).Error

    return articles, page, err
}
```

### 任务 3：事务处理 - 创建订单（难度：中等偏上）

```go
// 创建订单的完整事务流程

type CreateOrderReq struct {
    Items []CreateOrderItemReq `json:"items" binding:"required,min=1"`
}

type CreateOrderItemReq struct {
    ProductID uint `json:"product_id" binding:"required"`
    Quantity  int  `json:"quantity"   binding:"required,min=1"`
}

func (s *OrderService) CreateOrder(userID uint, req CreateOrderReq) (*model.Order, error) {
    var order model.Order

    err := s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 生成订单号
        orderNo := generateOrderNo() // 如: "ORD20240101120000001"

        // 2. 查询所有商品信息并校验库存
        var totalAmount float64
        var orderItems []model.OrderItem

        for _, item := range req.Items {
            var product model.Product
            // 使用 FOR UPDATE 锁定行，防止超卖
            if err := tx.Set("gorm:query_option", "FOR UPDATE").
                First(&product, item.ProductID).Error; err != nil {
                return fmt.Errorf("商品 %d 不存在", item.ProductID)
            }

            if product.Stock < item.Quantity {
                return fmt.Errorf("商品 %s 库存不足（剩余 %d）", product.Name, product.Stock)
            }

            // 3. 扣减库存
            result := tx.Model(&product).
                Where("stock >= ?", item.Quantity).
                Update("stock", gorm.Expr("stock - ?", item.Quantity))
            if result.RowsAffected == 0 {
                return fmt.Errorf("商品 %s 库存扣减失败", product.Name)
            }

            orderItems = append(orderItems, model.OrderItem{
                ProductID: product.ID,
                Quantity:  item.Quantity,
                Price:     product.Price,
            })
            totalAmount += product.Price * float64(item.Quantity)
        }

        // 4. 创建订单
        order = model.Order{
            OrderNo:     orderNo,
            UserID:      userID,
            TotalAmount: totalAmount,
            Status:      "pending",
            OrderItems:  orderItems,
        }
        if err := tx.Create(&order).Error; err != nil {
            return fmt.Errorf("创建订单失败: %w", err)
        }

        return nil // 提交事务
    })

    if err != nil {
        return nil, err
    }

    // 5. 重新查询订单（含关联数据）
    s.db.Preload("OrderItems.Product").First(&order, order.ID)
    return &order, nil
}
```

### 任务 4：常用 Scope 集合（难度：中等）

```go
// internal/repository/scopes.go

// 状态过滤
func WithStatus(status string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if status == "" {
            return db
        }
        return db.Where("status = ?", status)
    }
}

// 时间范围
func CreatedBetween(start, end time.Time) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("created_at BETWEEN ? AND ?", start, end)
    }
}

// 排序
func OrderBy(field, order string) func(*gorm.DB) *gorm.DB {
    allowedFields := map[string]bool{
        "created_at": true, "updated_at": true,
        "views": true, "likes": true, "price": true,
    }
    return func(db *gorm.DB) *gorm.DB {
        if !allowedFields[field] {
            field = "created_at"
        }
        if order != "asc" {
            order = "desc"
        }
        return db.Order(fmt.Sprintf("%s %s", field, order))
    }
}

// 模糊搜索
func Search(keyword string, fields ...string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if keyword == "" {
            return db
        }
        conditions := make([]string, len(fields))
        args := make([]interface{}, len(fields))
        for i, f := range fields {
            conditions[i] = fmt.Sprintf("%s LIKE ?", f)
            args[i] = "%" + keyword + "%"
        }
        return db.Where(strings.Join(conditions, " OR "), args...)
    }
}

// 使用示例：
db.Scopes(
    WithStatus("published"),
    Search("Go", "title", "content"),
    OrderBy("views", "desc"),
    Paginate(&page),
).Preload("Author").Find(&articles)
```

---

## 自检清单

- [ ] 掌握 Preload 和嵌套 Preload 的使用
- [ ] 能正确使用 GORM 事务，包括错误回滚
- [ ] 分页封装通用且可复用
- [ ] 多条件筛选查询正确处理空值（不添加条件）
- [ ] Scope 封装可复用且安全（防止 SQL 注入）
- [ ] 创建订单事务中正确处理库存校验和扣减
