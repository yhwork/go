package main

import (
	"fmt"
	"reflect"
)

// 创建 user 结构体
type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 创建 classInfo 结构体
type classInfo struct {
	ClassName string `json:"class_name"` // 班级名称
	ClassUser []user `json:"class_user"` // 用户列表
}

func main() {
	fmt.Println("hello world")
	u1 := user{
		Name: "hangman",
		Age:  18,
	}
	// 创建 user 实例
	class := classInfo{
		ClassName: "1607",
		ClassTeacher: map[string]user{
			"teacher1": {
				Name: "hangman",
				Age:  18,
			},
			"teacher2": {
				Name: "hangman",
				Age:  18,
			},
		}
		// 创建 user 实例列表 普通
		ClassUser: []user{
			u1,
			{
				Name: "hangman2",
				Age:  19,
			},
		},
	}
	// 输出 user 实例
	fmt.Println(class.ClassName)        // 输出班级名称
	// 循环输出 user 实例
	for i := 0; i < len(class.ClassUser); i++ {
		fmt.Println(class.ClassUser[i].Name)
		fmt.Println(class.ClassUser[i].Age)
	}
	// 循环输出 user 实例
	for key, v := range class.ClassUser {
		fmt.Println(key, v.Name)
		fmt.Println(key, v.Age)
	}

	// 遍历 user 对象的字段（key）
	fmt.Println("\n=== 使用反射遍历 user 对象字段 ===")
	userExample := class.ClassUser[0]
	val := reflect.ValueOf(userExample)
	typ := reflect.TypeOf(userExample)
	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name           // 字段名
		fieldValue := val.Field(i).Interface()   // 字段值
		jsonTag := typ.Field(i).Tag.Get("json")  // json 标签
		fmt.Printf("字段：%s, JSON Key: %s, 值：%v\n", fieldName, jsonTag, fieldValue)
	}

}
