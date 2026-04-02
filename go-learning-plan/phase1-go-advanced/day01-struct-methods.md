# Day 01 - 结构体与方法

## 今日目标

深入理解 Go 的结构体（struct）定义、方法接收者（值接收者 vs 指针接收者）、结构体组合（嵌入），为后续面向对象风格的编程打下基础。

---

## 知识点
[go.mod](..%2F..%2Fgo_demo%2Fday01%2Fgo.mod)
### 1. 结构体定义

```go
type User struct {
    ID       uint
    Username string
    Email    string
    Age      int
}

// 创建结构体实例的几种方式
u1 := User{ID: 1, Username: "alice", Email: "alice@example.com", Age: 25}
u2 := User{}              // 零值初始化
u3 := new(User)           // 返回 *User
u4 := &User{Username: "bob"} // 返回 *User，未赋值字段为零值
```

### 2. 方法接收者

```go
// 值接收者：方法内部操作的是副本，不会修改原始值
func (u User) FullInfo() string {
    return fmt.Sprintf("%s (%s)", u.Username, u.Email)
}

// 指针接收者：方法内部操作的是原始值，可以修改
func (u *User) UpdateEmail(newEmail string) {
    u.Email = newEmail
}
```

**选择原则：**
- 需要修改接收者的值 -> 用指针接收者
- 结构体较大，避免拷贝开销 -> 用指针接收者
- 需要实现某个接口且接口方法中有指针接收者 -> 所有方法统一用指针接收者
- 结构体较小且不需要修改 -> 值接收者也可以

### 3. 结构体组合（嵌入）

Go 没有继承，但通过组合实现类似效果：

```go
type BaseModel struct {
    ID        uint
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Product struct {
    BaseModel        // 嵌入，Product 自动拥有 BaseModel 的所有字段和方法
    Name     string
    Price    float64
    Stock    int
}

p := Product{
    BaseModel: BaseModel{ID: 1, CreatedAt: time.Now()},
    Name:      "Go 编程指南",
    Price:     59.9,
    Stock:     100,
}
fmt.Println(p.ID)        // 直接访问嵌入字段
fmt.Println(p.CreatedAt) // 直接访问嵌入字段
```

### 4. fmt.Stringer 接口

```go
// fmt.Stringer 是 Go 标准库中最常用的接口之一
type Stringer interface {
    String() string
}

// 实现 String() 方法后，fmt.Println 等函数会自动调用
func (u User) String() string {
    return fmt.Sprintf("User{ID: %d, Username: %s, Email: %s}", u.ID, u.Username, u.Email)
}

fmt.Println(u1) // 输出: User{ID: 1, Username: alice, Email: alice@example.com}
```

---

## 练习任务

### 任务 1：定义业务模型（难度：基础）

定义以下三个结构体，并建立它们之间的组合关系：

```go
// 要求：
// 1. 所有结构体都嵌入 BaseModel（包含 ID, CreatedAt, UpdatedAt）
// 2. User 包含字段：Username, Email, Phone, Password
// 3. Product 包含字段：Name, Description, Price, Stock, CategoryID
// 4. Order 包含字段：UserID, TotalAmount, Status, OrderItems(切片)
// 5. OrderItem 包含字段：OrderID, ProductID, Quantity, Price

type BaseModel struct {
    // 你的代码
}

type User struct {
    // 你的代码
}

type Product struct {
    // 你的代码
}

type Order struct {
    // 你的代码
}

type OrderItem struct {
    // 你的代码
}
```

### 任务 2：实现 String() 方法（难度：基础）

为每个结构体实现 `fmt.Stringer` 接口：

```go
// 要求：
// User.String() 输出格式: "User[1]: alice (alice@example.com)"
// Product.String() 输出格式: "Product[1]: Go编程指南 - ¥59.90 (库存: 100)"
// Order.String() 输出格式: "Order[1]: 用户1 - ¥299.80 (状态: pending) - 共2件商品"
```

### 任务 3：实现购物车（难度：中等）

```go
// 要求实现 Cart 结构体及以下方法：

type CartItem struct {
    Product  Product
    Quantity int
}

type Cart struct {
    UserID uint
    Items  []CartItem
}

// 1. AddItem：添加商品到购物车
//    - 如果商品已存在，增加数量
//    - 如果商品不存在，新增条目
//    - 需要检查库存是否充足，不足则返回 error
func (c *Cart) AddItem(product Product, quantity int) error

// 2. RemoveItem：从购物车移除商品
//    - 根据 ProductID 移除
//    - 如果商品不存在，返回 error
func (c *Cart) RemoveItem(productID uint) error

// 3. UpdateQuantity：修改商品数量
//    - 如果数量为 0，则移除该商品
//    - 需要检查库存
func (c *Cart) UpdateQuantity(productID uint, quantity int) error

// 4. GetTotal：计算购物车总价
func (c *Cart) GetTotal() float64

// 5. GetItemCount：获取购物车中商品种类数
func (c *Cart) GetItemCount() int

// 6. Clear：清空购物车
func (c *Cart) Clear()

// 7. String：实现 Stringer 接口，输出购物车摘要
func (c Cart) String() string
```

### 任务 4：编写 main 函数验证（难度：基础）

```go
func main() {
    // 1. 创建几个测试商品
    // 2. 创建购物车
    // 3. 添加商品、修改数量、移除商品
    // 4. 打印购物车信息和总价
    // 5. 测试库存不足的错误处理
}
```

---

## 参考答案骨架

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

// ============ 模型定义 ============

type BaseModel struct {
    ID        uint
    CreatedAt time.Time
    UpdatedAt time.Time
}

type User struct {
    BaseModel
    Username string
    Email    string
    Phone    string
    Password string
}

func (u User) String() string {
    return fmt.Sprintf("User[%d]: %s (%s)", u.ID, u.Username, u.Email)
}

// ... 补充 Product, Order, OrderItem 的定义和 String() 方法

// ============ 购物车实现 ============

var (
    ErrInsufficientStock = errors.New("库存不足")
    ErrItemNotFound      = errors.New("商品不在购物车中")
)

type CartItem struct {
    Product  Product
    Quantity int
}

type Cart struct {
    UserID uint
    Items  []CartItem
}

func (c *Cart) AddItem(product Product, quantity int) error {
    if product.Stock < quantity {
        return fmt.Errorf("%w: %s 库存 %d, 需要 %d", ErrInsufficientStock, product.Name, product.Stock, quantity)
    }

    // 检查是否已存在
    for i, item := range c.Items {
        if item.Product.ID == product.ID {
            newQty := item.Quantity + quantity
            if product.Stock < newQty {
                return fmt.Errorf("%w: %s 库存 %d, 需要 %d", ErrInsufficientStock, product.Name, product.Stock, newQty)
            }
            c.Items[i].Quantity = newQty
            return nil
        }
    }

    c.Items = append(c.Items, CartItem{Product: product, Quantity: quantity})
    return nil
}

// ... 补充其他方法的实现
```

---

## 自检清单

- [ ] 理解值接收者和指针接收者的区别
- [ ] 能正确使用结构体嵌入（组合）
- [ ] 理解 `fmt.Stringer` 接口的作用
- [ ] 购物车所有方法都有正确的错误处理
- [ ] 代码能正常编译运行，输出符合预期
