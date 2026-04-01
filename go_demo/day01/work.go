package main

// BaseModel 创建 BaseModel 结构体
type BaseModel struct {
	ID        int `json:"id"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

// Users User 创建 user 结构体
type Users struct {
	UserName string `json:"username"` // 用户名
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"` // 密码
	BaseModel
}

// Product 创建 product 结构体
type Product struct {
	Name        string `json:"name"`        // 产品名称
	Description string `json:"description"` // 描述
	Price       int    `json:"price"`       // 价格
	Stock       int    `json:"stock"`       // 库存
	Category    string `json:"category"`    // 分类
}

type OrderItem struct {
	OrderID   int `json:"order_id"`   // 订单ID
	ProductID int `json:"product_id"` // 产品ID
	Quantity  int `json:"quantity"`   // 数量
	Price     int `json:"price"`      //  价格
}

type Order struct {
	UserID      int         `json:"user_id"`      // 用户ID
	TotalAmount int         `json:"total_amount"` // 总金额
	Status      int         `json:"status"`
	OrderItems  []OrderItem `json:"order_items"` // 订单项列表
}
