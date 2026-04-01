package main

import (
	"fmt"
)

// User 创建 user 结构体
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 创建 classInfo 结构体
/**
First         User   `json:"first,omitempty"`  // 确定字段总会有合理的默认值,不关心"未设置"和"零值"的区别
First         *User   `json:"first,omitempty"` // 需要明确表示"未设置"状态 数据库 ORM 映射（区分 NULL 和默认值）

*/
type classInfo struct {
	First         User   `json:"first,omitempty"` // 首席用户
	ClassName     string `json:"class_name"`      // 班级名称
	ClassUserList []User `json:"class_user"`      // 用户列表
}

func main() {
	fmt.Println("————————————结构体定义——————————————————")
	u1 := User{
		Name: "hangman",
		Age:  18,
	}
	class2 := classInfo{
		ClassName: "1班",
		First:     User{},
		ClassUserList: []User{
			{
				Name: "张三",
				Age:  18,
			},
		},
	}
	fmt.Println("————————————结构体实例化——————————————————")
	fmt.Println(u1.Name)
	fmt.Println("————————————结构体字段访问——————————————————")
	fmt.Println(class2.ClassUserList[0].Name)
	fmt.Println("————————————结构体字段修改——————————————————")
	class2.ClassUserList[0].Name = "张三 1"
	fmt.Println(class2.ClassUserList[0].Name)

	fmt.Println("————————————方法调用示例——————————————————")
	// 使用 GetName() 方法 - 值接收者
	fmt.Println("使用 GetName() 方法获取名称:", u1.GetName())
	// 使用 SetName() 方法 - 指针接收者
	u1.SetName("李四")
	fmt.Println("使用 SetName() 方法修改后的名称:", u1.Name)
	fmt.Println("再次使用 GetName() 方法:", u1.GetName())

	// 对指针类型使用方法
	userPtr := &User{Name: "王五", Age: 20}
	fmt.Println("指针类型的 GetName():", userPtr.GetName())
	userPtr.SetName("赵六")
	fmt.Println("指针类型的 SetName() 后:", userPtr.GetName())

	fmt.Println("---------------------------------------父类---------------------------------------")
	parent := ParentClass{Name: "父类"}
	parent.ChildClass.Name = "这是子类"
	fmt.Println("直接调用优先显示父类名称:", parent.Name)
	parent.SetParentName("我是你爸爸")
	fmt.Println("通过 SetParentName 设置父类", parent.Name)
	fmt.Println("调用字类名称", parent.ChildClass.Name)
	fmt.Println("通过 GetParentName() 方法获取父类名称:", parent.GetParentName())
	fmt.Println("通过 GetChildName() 方法获取子类名称:", parent.GetChildName())
	fmt.Println("综合数据：", parent.GetAll())
}
