# Day 02 - 接口与多态

## 今日目标

掌握 Go 接口的隐式实现机制、空接口、类型断言、类型 switch，理解 Go 中实现多态的方式。

---

## 知识点

### 1. 接口定义与隐式实现

Go 的接口是隐式实现的，不需要 `implements` 关键字：

```go
// 定义接口
type PaymentMethod interface {
    Pay(amount float64) error
    Refund(amount float64) error
    Name() string
}

// 只要实现了接口中所有方法，就自动满足该接口
type CreditCard struct {
    CardNumber string
    ExpireDate string
    CVV        string
}

func (c CreditCard) Pay(amount float64) error {
    fmt.Printf("信用卡 %s 支付 ¥%.2f\n", c.CardNumber[len(c.CardNumber)-4:], amount)
    return nil
}

func (c CreditCard) Refund(amount float64) error {
    fmt.Printf("信用卡 %s 退款 ¥%.2f\n", c.CardNumber[len(c.CardNumber)-4:], amount)
    return nil
}

func (c CreditCard) Name() string {
    return "信用卡"
}

// CreditCard 自动实现了 PaymentMethod 接口
var _ PaymentMethod = CreditCard{} // 编译期检查
```

### 2. 接口组合

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 接口组合：ReadWriter 同时要求实现 Reader 和 Writer
type ReadWriter interface {
    Reader
    Writer
}
```

### 3. 空接口与类型断言

```go
// 空接口可以接受任何类型的值
var anything interface{} // 或 any（Go 1.18+）
anything = 42
anything = "hello"
anything = CreditCard{}

// 类型断言
value, ok := anything.(CreditCard)
if ok {
    fmt.Println("是信用卡:", value.CardNumber)
}

// 类型 switch
switch v := anything.(type) {
case int:
    fmt.Println("整数:", v)
case string:
    fmt.Println("字符串:", v)
case CreditCard:
    fmt.Println("信用卡:", v.Name())
default:
    fmt.Println("未知类型")
}
```

### 4. 接口最佳实践

```go
// 1. 接口应该小而精，通常只包含 1-3 个方法
// 好的接口
type Saver interface {
    Save(data []byte) error
}

// 不好的接口（太大了）
type DoEverything interface {
    Save(data []byte) error
    Load(id string) ([]byte, error)
    Delete(id string) error
    List() ([]string, error)
    // ...更多方法
}

// 2. 接口由使用方定义，而非实现方
// 3. 接受接口，返回结构体
func ProcessPayment(method PaymentMethod, amount float64) error {
    return method.Pay(amount)
}
```

---

## 练习任务

### 任务 1：定义支付接口及多种实现（难度：基础）

```go
// 定义 PaymentMethod 接口
type PaymentMethod interface {
    Pay(amount float64) error
    Refund(amount float64) error
    Name() string
    // 返回支付手续费率（如 0.006 表示 0.6%）
    FeeRate() float64
}

// 要求实现以下三种支付方式：

// 1. CreditCard（信用卡）
//    - 字段：CardNumber, HolderName, ExpireDate, CVV
//    - 手续费率：0.6%
//    - Pay 时打印: "信用卡[****1234]支付 ¥100.00，手续费 ¥0.60"

// 2. Alipay（支付宝）
//    - 字段：Account（手机号或邮箱）
//    - 手续费率：0.1%
//    - Pay 时打印: "支付宝[138****8888]支付 ¥100.00，手续费 ¥0.10"

// 3. WechatPay（微信支付）
//    - 字段：OpenID
//    - 手续费率：0.6%
//    - Pay 时打印: "微信支付[oXxx...xxxx]支付 ¥100.00，手续费 ¥0.60"
```

### 任务 2：实现支付处理函数（难度：中等）

```go
// 1. 计算手续费的通用函数
func CalculateFee(method PaymentMethod, amount float64) float64 {
    // 手续费 = 金额 * 手续费率，保留两位小数
}

// 2. 处理支付的函数
func ProcessPayment(method PaymentMethod, amount float64) error {
    // - 检查金额是否大于 0
    // - 计算手续费
    // - 调用 Pay 方法
    // - 打印支付摘要
}

// 3. 批量退款函数
func BatchRefund(payments []struct {
    Method PaymentMethod
    Amount float64
}) []error {
    // 对每笔支付执行退款，收集所有错误
}
```

### 任务 3：基于类型断言的差异化处理（难度：中等）

```go
// 不同支付方式有不同的额外操作，使用类型断言处理

type PaymentProcessor struct{}

func (p *PaymentProcessor) Process(method PaymentMethod, amount float64) error {
    // 1. 通用流程：调用 Pay
    // 2. 根据具体类型执行额外操作：
    //    - CreditCard: 检查卡号是否过期（ExpireDate）
    //    - Alipay: 打印完整的支付宝账号（不脱敏，内部日志用）
    //    - WechatPay: 发送微信模板消息通知（模拟打印即可）
    //    - 未知类型: 打印警告日志
    // 使用 type switch 实现
}
```

### 任务 4：实现一个简单的存储接口（难度：中等偏上）

```go
// 定义通用存储接口
type Storage interface {
    Save(key string, value interface{}) error
    Get(key string) (interface{}, error)
    Delete(key string) error
    List() ([]string, error)
}

// 实现两种存储：

// 1. MemoryStorage - 基于 map 的内存存储
type MemoryStorage struct {
    data map[string]interface{}
}

// 2. FileStorage - 基于文件的存储（每个 key 对应一个 JSON 文件）
type FileStorage struct {
    basePath string
}

// 实现一个函数，接受 Storage 接口，执行一系列操作来验证存储功能：
func TestStorage(s Storage) {
    // 保存几条数据
    // 读取并打印
    // 列出所有 key
    // 删除一条数据
    // 再次列出所有 key
}
```

### 任务 5：编译期接口检查（难度：基础）

```go
// 使用编译期检查确保所有类型都正确实现了接口
// 在代码中添加以下声明：

var _ PaymentMethod = (*CreditCard)(nil)
var _ PaymentMethod = (*Alipay)(nil)
var _ PaymentMethod = (*WechatPay)(nil)
var _ Storage = (*MemoryStorage)(nil)
var _ Storage = (*FileStorage)(nil)

// 思考：为什么用 (*Type)(nil) 而不是 Type{}？
// 答案：当方法使用指针接收者时，只有 *Type 实现了接口，Type 没有
```

---

## 参考答案骨架

```go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "os"
    "path/filepath"
)

// ============ 支付接口 ============

type PaymentMethod interface {
    Pay(amount float64) error
    Refund(amount float64) error
    Name() string
    FeeRate() float64
}

// ---- 信用卡 ----
type CreditCard struct {
    CardNumber string
    HolderName string
    ExpireDate string
    CVV        string
}

func (c CreditCard) Pay(amount float64) error {
    fee := math.Round(amount*c.FeeRate()*100) / 100
    // 脱敏：只显示后四位
    masked := "****" + c.CardNumber[len(c.CardNumber)-4:]
    fmt.Printf("信用卡[%s]支付 ¥%.2f，手续费 ¥%.2f\n", masked, amount, fee)
    return nil
}

func (c CreditCard) Refund(amount float64) error {
    masked := "****" + c.CardNumber[len(c.CardNumber)-4:]
    fmt.Printf("信用卡[%s]退款 ¥%.2f\n", masked, amount)
    return nil
}

func (c CreditCard) Name() string    { return "信用卡" }
func (c CreditCard) FeeRate() float64 { return 0.006 }

// ---- 支付宝（自行补充）----
// ---- 微信支付（自行补充）----

// ============ 存储接口 ============

type Storage interface {
    Save(key string, value interface{}) error
    Get(key string) (interface{}, error)
    Delete(key string) error
    List() ([]string, error)
}

type MemoryStorage struct {
    data map[string]interface{}
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{data: make(map[string]interface{})}
}

func (m *MemoryStorage) Save(key string, value interface{}) error {
    m.data[key] = value
    return nil
}

func (m *MemoryStorage) Get(key string) (interface{}, error) {
    val, ok := m.data[key]
    if !ok {
        return nil, fmt.Errorf("key %q not found", key)
    }
    return val, nil
}

// ... 补充 Delete 和 List 方法
// ... 补充 FileStorage 的完整实现
```

---

## 自检清单

- [ ] 理解 Go 接口的隐式实现机制
- [ ] 能区分值接收者和指针接收者对接口实现的影响
- [ ] 掌握类型断言和类型 switch 的使用场景
- [ ] 理解空接口 `interface{}` / `any` 的用途和局限
- [ ] 能使用编译期检查确保接口实现的正确性
- [ ] 理解"接受接口，返回结构体"的设计原则
