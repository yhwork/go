package main

import "fmt"

// ChildClass 子类
type ChildClass struct {
	Name  string
	Group []string
	Info  map[string]string
}

// ParentClass 父类
type ParentClass struct {
	Name       string
	ChildClass // 继承子类 	ChildClass
}

// InitParentName 初始化父类名称
func (p *ParentClass) InitParentName() {
	p.Name = "parent"
}

// GetParentName 获取父类名称
func (p *ParentClass) GetParentName() string {
	return p.Name
}

// SetParentName 设置父类名称
func (p *ParentClass) SetParentName(name string) {
	p.Name = name
}

// SetChildName 设置子类名称
func (p *ParentClass) SetChildName(name string) {
	p.ChildClass.Name = name
}

// GetChildName 获取子类名称
func (p *ParentClass) GetChildName() string {
	return p.ChildClass.Name
}

// GetAll 综合数据
func (p *ParentClass) GetAll() string {
	p.Info = map[string]string{"age": "18岁"}
	p.Group = []string{"hangman", "tom", "jack", "lucy", "lily", "lucy", "lily", "lucy", "lily", "lucy", "lily", "lucy", "lily", "lucy", "lily"}
	return fmt.Sprintf("ChildClass{Name: %s, Group: %s, Info: %s}", p.ChildClass.Name, p.ChildClass.Group, p.ChildClass.Info)
}
