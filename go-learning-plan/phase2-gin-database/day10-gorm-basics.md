# Day 10 - GORM 入门与模型定义

## 今日目标

掌握 GORM 的数据库连接、模型定义（字段标签、关联关系）、自动迁移、基础 CRUD 操作。

---

## 知识点

### 1. 连接数据库

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

dsn := "root:123456@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // 打印 SQL
})

// 获取底层 *sql.DB 设置连接池
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. 模型定义

```go
// gorm.Model 提供了 ID, CreatedAt, UpdatedAt, DeletedAt（软删除）
type User struct {
    gorm.Model                          // ID, CreatedAt, UpdatedAt, DeletedAt
    Username string `gorm:"uniqueIndex;size:50;not null"`
    Email    string `gorm:"uniqueIndex;size:100;not null"`
    Password string `gorm:"size:100;not null"`
    Phone    string `gorm:"size:20"`
    Avatar   string `gorm:"size:255"`
    Role     string `gorm:"size:20;default:user"` // user, admin
    Status   int    `gorm:"default:1"`             // 1=active, 0=disabled
}

// 自定义表名
func (User) TableName() string {
    return "users"
}
```

### 3. GORM 字段标签

```go
type Product struct {
    gorm.Model
    Name        string  `gorm:"size:100;not null;index"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"type:decimal(10,2);not null"`
    Stock       int     `gorm:"default:0;not null"`
    CategoryID  uint    `gorm:"index;not null"`
    IsActive    bool    `gorm:"default:true"`
}

// 常用标签：
// size:100          - VARCHAR(100)
// type:text         - TEXT 类型
// type:decimal(10,2)- DECIMAL(10,2)
// not null          - NOT NULL
// default:value     - 默认值
// uniqueIndex       - 唯一索引
// index             - 普通索引
// primaryKey        - 主键
// column:col_name   - 自定义列名
// -                 - 忽略该字段
```

### 4. 自动迁移

```go
db.AutoMigrate(&User{}, &Product{}, &Order{}, &OrderItem{})
// 自动创建表、添加缺失的列、创建索引
// 不会删除列或更改列类型（安全）
```

### 5. 基础 CRUD

```go
// Create
user := User{Username: "alice", Email: "alice@example.com"}
db.Create(&user) // user.ID 会自动填充

// Read
var user User
db.First(&user, 1)                          // 按主键查找
db.First(&user, "username = ?", "alice")     // 条件查找
db.Where("age > ?", 18).Find(&users)        // 多条记录

// Update
db.Model(&user).Update("email", "new@example.com")
db.Model(&user).Updates(User{Email: "new@example.com", Phone: "13800138000"})
db.Model(&user).Updates(map[string]interface{}{"email": "new@example.com"})

// Delete（软删除）
db.Delete(&user, 1) // 设置 DeletedAt，不会物理删除
db.Unscoped().Delete(&user, 1) // 物理删除
```

---

## 练习任务

### 任务 1：定义完整的数据模型（难度：中等）

```go
// internal/model/base.go
package model

import (
    "gorm.io/gorm"
    "time"
)

// 所有模型都嵌入 gorm.Model

// internal/model/user.go
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;size:50;not null"  json:"username"`
    Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password string `gorm:"size:100;not null"             json:"-"` // json:"-" 不返回密码
    Phone    string `gorm:"size:20"                       json:"phone"`
    Avatar   string `gorm:"size:255"                      json:"avatar"`
    Role     string `gorm:"size:20;default:user"          json:"role"`
    Status   int    `gorm:"default:1"                     json:"status"`
    // 关联
    Orders   []Order `json:"orders,omitempty"`
}

// internal/model/category.go
type Category struct {
    gorm.Model
    Name     string `gorm:"size:50;not null;uniqueIndex" json:"name"`
    ParentID *uint  `gorm:"index"                        json:"parent_id"` // 支持树形
    Sort     int    `gorm:"default:0"                    json:"sort"`
    // 关联
    Articles []Article `json:"articles,omitempty"`
    Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// internal/model/article.go
type Article struct {
    gorm.Model
    Title      string `gorm:"size:200;not null;index"      json:"title"`
    Content    string `gorm:"type:longtext;not null"        json:"content"`
    Summary    string `gorm:"size:500"                      json:"summary"`
    Cover      string `gorm:"size:255"                      json:"cover"`
    AuthorID   uint   `gorm:"index;not null"                json:"author_id"`
    CategoryID uint   `gorm:"index;not null"                json:"category_id"`
    Status     string `gorm:"size:20;default:draft;index"   json:"status"` // draft, published
    Views      int    `gorm:"default:0"                     json:"views"`
    Likes      int    `gorm:"default:0"                     json:"likes"`
    // 关联
    Author   User     `gorm:"foreignKey:AuthorID"  json:"author,omitempty"`
    Category Category `json:"category,omitempty"`
    Tags     []Tag    `gorm:"many2many:article_tags" json:"tags,omitempty"`
}

// internal/model/tag.go
type Tag struct {
    gorm.Model
    Name     string    `gorm:"size:30;not null;uniqueIndex" json:"name"`
    Articles []Article `gorm:"many2many:article_tags"       json:"articles,omitempty"`
}

// internal/model/order.go
type Order struct {
    gorm.Model
    OrderNo     string      `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
    UserID      uint        `gorm:"index;not null"               json:"user_id"`
    TotalAmount float64     `gorm:"type:decimal(10,2);not null"  json:"total_amount"`
    Status      string      `gorm:"size:20;default:pending"      json:"status"` // pending, paid, shipped, completed, cancelled
    // 关联
    User       User        `json:"user,omitempty"`
    OrderItems []OrderItem `json:"order_items,omitempty"`
}

// internal/model/order_item.go
type OrderItem struct {
    gorm.Model
    OrderID   uint    `gorm:"index;not null"              json:"order_id"`
    ProductID uint    `gorm:"index;not null"              json:"product_id"`
    Quantity  int     `gorm:"not null"                    json:"quantity"`
    Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`
    // 关联
    Product Product `json:"product,omitempty"`
}
```

### 任务 2：数据库初始化与迁移（难度：基础）

```go
// internal/database/database.go
package database

import (
    "fmt"
    "your-project/internal/config"
    "your-project/internal/model"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
    var err error

    // 1. 根据 driver 选择数据库驱动（mysql/postgres）
    // 2. 设置 GORM 配置（日志级别、命名策略等）
    // 3. 连接数据库
    // 4. 设置连接池参数
    // 5. 自动迁移

    return nil
}

func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        &model.Category{},
        &model.Article{},
        &model.Tag{},
        &model.Order{},
        &model.OrderItem{},
    )
}
```

### 任务 3：实现 User Repository（难度：中等偏上）

```go
// internal/repository/user_repo.go
package repository

type UserRepository interface {
    Create(user *model.User) error
    FindByID(id uint) (*model.User, error)
    FindByUsername(username string) (*model.User, error)
    FindByEmail(email string) (*model.User, error)
    Update(user *model.User) error
    Delete(id uint) error
    List(page, pageSize int, keyword string) ([]model.User, int64, error)
    ExistsByUsername(username string) (bool, error)
    ExistsByEmail(email string) (bool, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &NotFoundError{Resource: "User", ID: id}
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
    var users []model.User
    var total int64

    query := r.db.Model(&model.User{})

    if keyword != "" {
        query = query.Where("username LIKE ? OR email LIKE ?",
            "%"+keyword+"%", "%"+keyword+"%")
    }

    // 先查总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // 再查分页数据
    offset := (page - 1) * pageSize
    if err := query.Offset(offset).Limit(pageSize).
        Order("created_at DESC").
        Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}

// 补充其他方法...
```

### 任务 4：编写种子数据（难度：基础）

```go
// internal/database/seed.go
func Seed(db *gorm.DB) error {
    // 检查是否已有数据，避免重复插入
    var count int64
    db.Model(&model.User{}).Count(&count)
    if count > 0 {
        return nil
    }

    // 创建管理员用户
    admin := model.User{
        Username: "admin",
        Email:    "admin@example.com",
        Password: hashPassword("admin123"), // bcrypt 加密
        Role:     "admin",
    }
    db.Create(&admin)

    // 创建测试用户
    // 创建分类
    // 创建标签
    // 创建测试文章

    return nil
}
```

---

## 自检清单

- [ ] 能正确连接 MySQL/PostgreSQL 数据库
- [ ] 所有模型的 GORM 标签正确（类型、索引、默认值）
- [ ] 关联关系定义正确（一对多、多对多）
- [ ] 自动迁移成功创建所有表
- [ ] Repository 的 CRUD 方法全部实现并可用
- [ ] 分页查询正确处理 offset 和 limit
- [ ] 软删除生效（Delete 后记录仍在数据库中）
