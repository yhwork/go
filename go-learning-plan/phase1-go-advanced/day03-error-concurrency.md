# Day 03 - 错误处理与并发基础

## 今日目标

掌握 Go 的错误处理最佳实践（自定义错误、错误包装、errors.Is/As），以及并发编程基础（goroutine、WaitGroup、context）。

---

## 知识点

### 1. error 接口

```go
// Go 的 error 是一个简单的接口
type error interface {
    Error() string
}

// 创建错误的几种方式
err1 := errors.New("something went wrong")
err2 := fmt.Errorf("user %d not found", 42)
```

### 2. 自定义错误类型

```go
// 自定义错误类型可以携带更多上下文信息
type NotFoundError struct {
    Resource string
    ID       interface{}
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s (id=%v) not found", e.Resource, e.ID)
}

// 使用
func FindUser(id uint) (*User, error) {
    // ... 查找逻辑
    return nil, &NotFoundError{Resource: "User", ID: id}
}
```

### 3. 错误包装与解包（Go 1.13+）

```go
// 包装错误：保留原始错误的同时添加上下文
func GetUser(id uint) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("获取用户失败: %w", err) // %w 包装错误
    }
    return user, nil
}

// 解包：判断错误链中是否包含特定错误
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // 处理 not found 情况
}

// 提取特定类型的错误
var notFoundErr *NotFoundError
if errors.As(err, &notFoundErr) {
    fmt.Printf("资源 %s 未找到\n", notFoundErr.Resource)
}
```

### 4. Goroutine

```go
// goroutine 是 Go 的轻量级线程
go func() {
    fmt.Println("我在另一个 goroutine 中运行")
}()

// 注意：main 函数结束时，所有 goroutine 都会被终止
// 需要用同步机制等待 goroutine 完成
```

### 5. sync.WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("任务 %d 完成\n", id)
    }(i)
}

wg.Wait() // 等待所有 goroutine 完成
```

### 6. context 包

```go
// context 用于控制 goroutine 的生命周期

// 带超时的 context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 在 goroutine 中检查 context
select {
case <-ctx.Done():
    fmt.Println("超时或被取消:", ctx.Err())
case result := <-doWork(ctx):
    fmt.Println("完成:", result)
}
```

---

## 练习任务

### 任务 1：定义业务错误体系（难度：基础）

```go
// 要求：
// 1. 定义以下错误类型，每个都实现 error 接口

// NotFoundError - 资源未找到
type NotFoundError struct {
    Resource string // 资源类型，如 "User", "Product"
    ID       interface{}
}

// ValidationError - 参数验证失败
type ValidationError struct {
    Field   string // 字段名
    Message string // 错误信息
}

// AuthError - 认证/授权错误
type AuthError struct {
    Reason string // 原因，如 "token expired", "invalid credentials"
}

// 2. 定义哨兵错误（sentinel errors）
var (
    ErrNotFound       = errors.New("resource not found")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
    ErrInternalServer = errors.New("internal server error")
)

// 3. 实现一个 HTTPStatusCode() 方法，返回对应的 HTTP 状态码
//    NotFoundError -> 404
//    ValidationError -> 400
//    AuthError -> 401
```

### 任务 2：错误包装与处理链（难度：中等）

```go
// 模拟三层架构的错误传递：Repository -> Service -> Handler

// Repository 层
func (r *UserRepo) FindByID(id uint) (*User, error) {
    // 模拟：id > 100 时返回 NotFoundError
    // 返回: &NotFoundError{Resource: "User", ID: id}
}

// Service 层 - 包装 Repository 的错误
func (s *UserService) GetUser(id uint) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        // 使用 %w 包装，添加业务上下文
        return nil, fmt.Errorf("UserService.GetUser: %w", err)
    }
    return user, nil
}

// Handler 层 - 解析错误并返回合适的响应
func (h *UserHandler) GetUser(id uint) {
    user, err := h.service.GetUser(id)
    if err != nil {
        // 使用 errors.As 提取具体错误类型
        // 根据不同错误类型返回不同的 HTTP 状态码和消息
        var notFoundErr *NotFoundError
        if errors.As(err, &notFoundErr) {
            fmt.Printf("[404] %s\n", notFoundErr.Error())
            return
        }
        // ... 处理其他错误类型
        fmt.Printf("[500] 内部错误: %s\n", err.Error())
        return
    }
    fmt.Printf("[200] %s\n", user)
}

// 在 main 中测试完整的错误传递链
```

### 任务 3：并发请求模拟器（难度：中等）

```go
// 模拟同时请求 5 个 URL 并汇总结果

type FetchResult struct {
    URL      string
    Status   int
    Duration time.Duration
    Err      error
}

func FetchURL(url string) FetchResult {
    start := time.Now()
    // 模拟 HTTP 请求（用 time.Sleep 模拟随机耗时 100-500ms）
    // 随机返回成功(200)或失败(500)
    duration := time.Since(start)
    return FetchResult{URL: url, Status: 200, Duration: duration}
}

func FetchAll(urls []string) []FetchResult {
    // 使用 goroutine + WaitGroup 并发请求所有 URL
    // 收集所有结果到切片中
    // 注意：并发写入切片需要用 mutex 或预分配索引
}

// 测试 URL 列表
urls := []string{
    "https://api.example.com/users",
    "https://api.example.com/products",
    "https://api.example.com/orders",
    "https://api.example.com/categories",
    "https://api.example.com/reviews",
}
```

### 任务 4：带超时的任务执行器（难度：中等偏上）

```go
// 实现一个通用的带超时任务执行器

type TaskResult struct {
    Value interface{}
    Err   error
}

// ExecuteWithTimeout 在指定超时时间内执行任务
// 如果任务在超时前完成，返回结果
// 如果超时，返回 context.DeadlineExceeded 错误
func ExecuteWithTimeout(ctx context.Context, timeout time.Duration, task func(ctx context.Context) (interface{}, error)) TaskResult {
    // 1. 创建带超时的 context
    // 2. 在 goroutine 中执行任务
    // 3. 用 select 等待结果或超时
}

// 测试用例：
// 1. 正常任务（100ms），超时 1s -> 应该成功
// 2. 慢任务（2s），超时 500ms -> 应该超时
// 3. 任务出错 -> 应该返回错误
// 4. 外部 context 被取消 -> 应该返回取消错误
```

---

## 参考答案骨架

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// ============ 错误定义 ============

type NotFoundError struct {
    Resource string
    ID       interface{}
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s (id=%v) not found", e.Resource, e.ID)
}

func (e *NotFoundError) HTTPStatusCode() int { return 404 }

// ... 补充 ValidationError, AuthError

// ============ 并发请求 ============

type FetchResult struct {
    URL      string
    Status   int
    Duration time.Duration
    Err      error
}

func FetchAll(urls []string) []FetchResult {
    results := make([]FetchResult, len(urls))
    var wg sync.WaitGroup

    for i, url := range urls {
        wg.Add(1)
        go func(idx int, u string) {
            defer wg.Done()
            results[idx] = FetchURL(u) // 用索引避免并发问题
        }(i, url)
    }

    wg.Wait()
    return results
}

// ============ 超时执行器 ============

func ExecuteWithTimeout(ctx context.Context, timeout time.Duration, task func(ctx context.Context) (interface{}, error)) TaskResult {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    ch := make(chan TaskResult, 1)
    go func() {
        val, err := task(ctx)
        ch <- TaskResult{Value: val, Err: err}
    }()

    select {
    case result := <-ch:
        return result
    case <-ctx.Done():
        return TaskResult{Err: ctx.Err()}
    }
}
```

---

## 自检清单

- [ ] 能定义自定义错误类型并实现 error 接口
- [ ] 理解 `%w` 包装错误和 `errors.Is/As` 解包错误
- [ ] 能使用 goroutine + WaitGroup 实现并发任务
- [ ] 理解 context 的超时控制和取消机制
- [ ] 并发代码没有数据竞争（race condition）
