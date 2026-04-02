package main

import (
	"fmt"
)

// 购物车逻辑

// CartItem 购物车结构体
type CartItem struct {
	Product  Product `json:"product"`  // 商品
	Quantity int     `json:"quantity"` // 商品数量
}

type Cart struct {
	UserId string     `json:"userId"`
	Items  []CartItem `json:"items"` // 购物车项
}

// AddItem 添加商品到购物车
func (c *Cart) AddItem(product Product, quantity int) error {
	if product.Stock < quantity {
		return fmt.Errorf(" %s 库存 %d, 需要 %d", product.Name, product.Stock, quantity)
	}
	// 查看商品是否已经存在
	for i, item := range c.Items {
		if item.Product.Name == product.Name {
			newQty := item.Quantity + quantity
			if product.Stock < newQty {
				return fmt.Errorf(" %s 库存 %d, 需要 %d", product.Name, product.Stock, newQty)
			}
			c.Items[i].Quantity = newQty
			return nil
		}
	}

	// 添加商品到购物车
	c.Items = append(c.Items, CartItem{
		Product:  product,
		Quantity: quantity,
	})
	// 减少库存
	product.Stock -= quantity
	return nil
}
