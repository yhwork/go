# Day 06 - 单元测试与基准测试

## 今日目标

掌握 Go testing 包、表驱动测试、testify 断言库、mock 测试、基准测试（benchmark），为之前编写的代码补充测试。

---

## 知识点

### 1. 测试文件规范

```
user_service.go       -> user_service_test.go  （同一目录）
user_service_test.go  -> package service        （同包或 service_test）
```

- 测试文件以 `_test.go` 结尾
- 测试函数以 `Test` 开头，参数为 `*testing.T`
- 运行：`go test ./... -v`

### 2. 基础测试

```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d, want 5", result)
    }
}
```

### 3. 表驱动测试

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 4. testify 断言

```go
import "github.com/stretchr/testify/assert"

func TestUser(t *testing.T) {
    user := NewUser("alice")

    assert.NotNil(t, user)
    assert.Equal(t, "alice", user.Username)
    assert.Empty(t, user.Email)
    assert.NoError(t, user.Validate())
}

// require 与 assert 的区别：require 失败后立即停止测试
import "github.com/stretchr/testify/require"

func TestCritical(t *testing.T) {
    result, err := DoSomething()
    require.NoError(t, err)     // 失败则停止
    assert.Equal(t, 42, result) // 这行只有 err == nil 才会执行
}
```

### 5. Mock 测试

```go
// 定义接口
type UserRepository interface {
    FindByID(id uint) (*User, error)
    Create(user *User) error
}

// Mock 实现
type MockUserRepo struct {
    users map[uint]*User
}

func NewMockUserRepo() *MockUserRepo {
    return &MockUserRepo{users: make(map[uint]*User)}
}

func (m *MockUserRepo) FindByID(id uint) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, &NotFoundError{Resource: "User", ID: id}
    }
    return user, nil
}

func (m *MockUserRepo) Create(user *User) error {
    m.users[user.ID] = user
    return nil
}

// 使用 Mock 测试 Service 层
func TestUserService_GetUser(t *testing.T) {
    repo := NewMockUserRepo()
    repo.users[1] = &User{ID: 1, Username: "alice"}

    service := NewUserService(repo)
    user, err := service.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, "alice", user.Username)
}
```

### 6. 基准测试

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

// 运行基准测试
// go test -bench=. -benchmem
// 输出示例:
// BenchmarkAdd-8   1000000000   0.25 ns/op   0 B/op   0 allocs/op
```

### 7. TestMain（测试套件初始化）

```go
func TestMain(m *testing.M) {
    // 测试前的全局初始化
    setup()

    code := m.Run() // 运行所有测试

    // 测试后的全局清理
    teardown()

    os.Exit(code)
}
```

---

## 练习任务

### 任务 1：为购物车编写表驱动测试（难度：中等）

为 Day 01 的购物车代码编写完整测试：

```go
// cart_test.go
package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// 测试辅助函数
func newTestProduct(id uint, name string, price float64, stock int) Product {
    return Product{
        BaseModel: BaseModel{ID: id},
        Name:      name,
        Price:     price,
        Stock:     stock,
    }
}

func TestCart_AddItem(t *testing.T) {
    tests := []struct {
        name        string
        setup       func(*Cart)        // 测试前的准备
        product     Product
        quantity    int
        wantErr     bool
        errContains string             // 错误信息应包含的关键字
        wantCount   int                // 添加后购物车商品种类数
    }{
        {
            name:      "添加新商品",
            setup:     func(c *Cart) {},
            product:   newTestProduct(1, "Go 编程", 59.9, 100),
            quantity:  2,
            wantErr:   false,
            wantCount: 1,
        },
        {
            name: "添加已有商品增加数量",
            setup: func(c *Cart) {
                c.AddItem(newTestProduct(1, "Go 编程", 59.9, 100), 1)
            },
            product:   newTestProduct(1, "Go 编程", 59.9, 100),
            quantity:  2,
            wantErr:   false,
            wantCount: 1, // 还是 1 种
        },
        {
            name:        "库存不足",
            setup:       func(c *Cart) {},
            product:     newTestProduct(2, "限量版", 999, 2),
            quantity:    5,
            wantErr:     true,
            errContains: "库存不足",
        },
        // 补充更多用例：
        // - 添加数量为 0
        // - 累加后超过库存
        // - 多种不同商品
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cart := &Cart{UserID: 1}
            tt.setup(cart)

            err := cart.AddItem(tt.product, tt.quantity)

            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errContains)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.wantCount, cart.GetItemCount())
            }
        })
    }
}

// 补充以下测试：
func TestCart_RemoveItem(t *testing.T) {
    // 测试用例：
    // 1. 移除存在的商品 -> 成功
    // 2. 移除不存在的商品 -> 返回错误
    // 3. 移除后购物车总价更新
}

func TestCart_UpdateQuantity(t *testing.T) {
    // 测试用例：
    // 1. 修改数量为正数 -> 成功
    // 2. 修改数量为 0 -> 移除商品
    // 3. 修改数量超过库存 -> 返回错误
}

func TestCart_GetTotal(t *testing.T) {
    // 测试用例：
    // 1. 空购物车 -> 0
    // 2. 单件商品 -> 正确计算
    // 3. 多件商品 -> 正确累加
    // 注意浮点数比较用 assert.InDelta
}
```

### 任务 2：Mock Repository 测试 Service（难度：中等偏上）

```go
// 定义接口和 Service

type UserRepository interface {
    FindByID(id uint) (*User, error)
    FindByUsername(username string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(id uint) (*User, error) {
    if id == 0 {
        return nil, &ValidationError{Field: "id", Message: "用户 ID 不能为 0"}
    }
    return s.repo.FindByID(id)
}

func (s *UserService) Register(username, email, password string) (*User, error) {
    // 1. 验证用户名长度 3-20
    // 2. 检查用户名是否已存在
    // 3. 密码加密
    // 4. 创建用户
    // 5. 返回用户（不含密码）
}

// 编写 Mock 和测试：
type MockUserRepo struct {
    users    map[uint]*User
    byName   map[string]*User
    nextID   uint
    // 可选：记录调用次数
    calls    map[string]int
}

func TestUserService_GetUser(t *testing.T) {
    // 测试用例：
    // 1. 正常获取用户
    // 2. 用户不存在 -> NotFoundError
    // 3. ID 为 0 -> ValidationError
}

func TestUserService_Register(t *testing.T) {
    // 测试用例：
    // 1. 正常注册
    // 2. 用户名太短 -> ValidationError
    // 3. 用户名已存在 -> 错误
    // 4. 返回的用户密码字段为空
}
```

### 任务 3：基准测试对比（难度：中等）

```go
// 比较不同并发 worker 数量对 Worker Pool 性能的影响

func benchmarkWorkerPool(b *testing.B, workerCount int) {
    tasks := make([]int, 100)
    for i := range tasks {
        tasks[i] = i
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 用指定数量的 worker 处理所有任务
        runWorkerPool(tasks, workerCount)
    }
}

func BenchmarkWorkerPool_1(b *testing.B)  { benchmarkWorkerPool(b, 1) }
func BenchmarkWorkerPool_2(b *testing.B)  { benchmarkWorkerPool(b, 2) }
func BenchmarkWorkerPool_4(b *testing.B)  { benchmarkWorkerPool(b, 4) }
func BenchmarkWorkerPool_8(b *testing.B)  { benchmarkWorkerPool(b, 8) }
func BenchmarkWorkerPool_16(b *testing.B) { benchmarkWorkerPool(b, 16) }

// 同时比较 Mutex vs Channel 在并发计数场景的性能

func BenchmarkMutexCounter(b *testing.B) {
    var mu sync.Mutex
    counter := 0

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    })
}

func BenchmarkChannelCounter(b *testing.B) {
    ch := make(chan int, 1)
    ch <- 0

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            val := <-ch
            ch <- val + 1
        }
    })
}

func BenchmarkAtomicCounter(b *testing.B) {
    var counter int64

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            atomic.AddInt64(&counter, 1)
        }
    })
}
```

### 任务 4：测试覆盖率分析（难度：基础）

```bash
# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out

# 查看覆盖率摘要
go tool cover -func=coverage.out

# 生成 HTML 报告（在浏览器中查看）
go tool cover -html=coverage.out -o coverage.html

# 目标：核心业务代码覆盖率 > 80%
```

---

## 参考答案骨架

```go
// mock_user_repo.go
package service_test

type MockUserRepo struct {
    users  map[uint]*User
    byName map[string]*User
    nextID uint
}

func NewMockUserRepo() *MockUserRepo {
    return &MockUserRepo{
        users:  make(map[uint]*User),
        byName: make(map[string]*User),
        nextID: 1,
    }
}

func (m *MockUserRepo) FindByID(id uint) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, &NotFoundError{Resource: "User", ID: id}
    }
    return user, nil
}

func (m *MockUserRepo) FindByUsername(username string) (*User, error) {
    user, ok := m.byName[username]
    if !ok {
        return nil, &NotFoundError{Resource: "User", ID: username}
    }
    return user, nil
}

func (m *MockUserRepo) Create(user *User) error {
    user.ID = m.nextID
    m.nextID++
    m.users[user.ID] = user
    m.byName[user.Username] = user
    return nil
}

// ... 补充 Update, Delete

// user_service_test.go
func TestUserService_GetUser(t *testing.T) {
    repo := NewMockUserRepo()
    repo.Create(&User{Username: "alice", Email: "alice@example.com"})

    svc := NewUserService(repo)

    t.Run("existing user", func(t *testing.T) {
        user, err := svc.GetUser(1)
        require.NoError(t, err)
        assert.Equal(t, "alice", user.Username)
    })

    t.Run("not found", func(t *testing.T) {
        _, err := svc.GetUser(999)
        require.Error(t, err)

        var notFoundErr *NotFoundError
        assert.True(t, errors.As(err, &notFoundErr))
    })

    t.Run("invalid id", func(t *testing.T) {
        _, err := svc.GetUser(0)
        require.Error(t, err)

        var validErr *ValidationError
        assert.True(t, errors.As(err, &validErr))
    })
}
```

---

## 自检清单

- [ ] 掌握表驱动测试的写法
- [ ] 会使用 testify 的 assert 和 require
- [ ] 理解 Mock 的作用和实现方式
- [ ] 能编写基准测试并分析结果
- [ ] 能生成和查看测试覆盖率报告
- [ ] 核心代码测试覆盖率达到 80% 以上
