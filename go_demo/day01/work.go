package main

import "fmt"

// BaseModel 基础模型结构体
// 用于嵌入到其他结构体中，提供通用的字段（ID、创建时间、更新时间）
// 使用匿名嵌入方式，让其他结构体自动继承这些字段
type BaseModel struct {
	ID        int `json:"id"`         // 唯一标识符
	CreatedAt int `json:"created_at"` // 创建时间戳（Unix 时间戳）
	UpdatedAt int `json:"updated_at"` // 最后更新时间戳（Unix 时间戳）
}

// Users 用户信息结构体
// 嵌入 BaseModel，自动获得 ID、CreatedAt、UpdatedAt 字段
// 用于表示系统中的用户实体
type Users struct {
	UserName  string `json:"username"` // 用户名，用户的登录名称
	Email     string `json:"email"`    // 电子邮箱地址
	Phone     string `json:"phone"`    // 手机号码
	Password  string `json:"password"` // 密码（建议存储加密后的哈希值）
	BaseModel        // 匿名嵌入 BaseModel，继承其所有字段
}

// Product 产品信息结构体
// 用于表示电商系统中的商品或产品实体
type Product struct {
	Name        string `json:"name"`        // 产品名称，商品的标题或名称
	Description string `json:"description"` // 产品描述，详细说明商品信息
	Price       int    `json:"price"`       // 价格（单位：分或元，根据业务定义）
	Stock       int    `json:"stock"`       // 库存数量，可售商品总数
	Category    string `json:"category"`    // 产品分类，如“电子产品”、“图书”等
}

// OrderItem 订单项结构体
// 表示订单中的单个商品条目，包含商品信息和购买数量
type OrderItem struct {
	OrderID   int `json:"order_id"`   // 所属订单 ID，关联到 Order 表
	ProductID int `json:"product_id"` // 产品 ID，关联到 Product 表
	Quantity  int `json:"quantity"`   // 购买数量，该商品的购买件数
	Price     int `json:"price"`      // 单价，购买时的商品价格（可能因促销而变化）
}

// Order 订单结构体
// 表示用户的完整订单信息，包含多个 OrderItem
type Order struct {
	UserID      int         `json:"user_id"`      // 用户 ID，下单用户的标识
	TotalAmount int         `json:"total_amount"` // 订单总金额，所有商品价格的总和
	Status      int         `json:"status"`       // 订单状态（如：0=待支付，1=已支付，2=已发货，3=已完成，-1=已取消）
	OrderItems  []OrderItem `json:"order_items"`  // 订单项列表，包含订单中的所有商品
}

// 为每个结构体实现，满足 fmt.Stringer 接口
func (u *Users) String() string {
	return fmt.Sprintf("用户信息：%s %s %s %s", u.UserName, u.Email, u.Phone, u.Password)
}

// 产品 信息
func (p *Product) String() string {
	return fmt.Sprintf("产品信息：%s %s %d %d %s", p.Name, p.Description, p.Price, p.Stock, p.Category)
}

// 订单信息
func (o *Order) String() string {
	return fmt.Sprintf("订单信息：%d %d %d %v", o.UserID, o.TotalAmount, o.Status, o.OrderItems)
}
