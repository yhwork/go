# Day 04 - Channel 与并发模式

## 今日目标

掌握 Channel 的各种操作（缓冲/非缓冲、方向、关闭），以及常见的并发模式（生产者-消费者、Worker Pool、Fan-in/Fan-out、Select 超时控制）。

---

## 知识点

### 1. Channel 基础

```go
// 无缓冲 channel：发送和接收必须同步
ch := make(chan int)

// 有缓冲 channel：缓冲区满之前发送不会阻塞
ch := make(chan int, 10)

// 发送和接收
ch <- 42     // 发送
value := <-ch // 接收

// 关闭 channel
close(ch)

// 用 range 遍历 channel（channel 关闭后退出）
for v := range ch {
    fmt.Println(v)
}

// 检查 channel 是否关闭
value, ok := <-ch
if !ok {
    fmt.Println("channel 已关闭")
}
```

### 2. Channel 方向

```go
// 只发送
func producer(ch chan<- int) {
    ch <- 1
}

// 只接收
func consumer(ch <-chan int) {
    v := <-ch
}
```

### 3. Select 语句

```go
select {
case msg := <-ch1:
    fmt.Println("从 ch1 收到:", msg)
case msg := <-ch2:
    fmt.Println("从 ch2 收到:", msg)
case ch3 <- 42:
    fmt.Println("发送到 ch3")
case <-time.After(5 * time.Second):
    fmt.Println("超时")
default:
    fmt.Println("没有 channel 就绪")
}
```

### 4. 常见并发模式

#### 生产者-消费者

```go
func producer(ch chan<- int) {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch) // 生产完毕后关闭
}

func consumer(id int, ch <-chan int) {
    for v := range ch {
        fmt.Printf("消费者%d处理: %d\n", id, v)
    }
}
```

#### Worker Pool

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        result := job * 2 // 处理任务
        results <- result
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // 启动 3 个 worker
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // 发送任务
    for j := 1; j <= 10; j++ {
        jobs <- j
    }
    close(jobs)

    // 收集结果
    for r := 1; r <= 10; r++ {
        <-results
    }
}
```

#### Fan-in（扇入）

```go
// 将多个 channel 合并为一个
func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    merged := make(chan int)

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                merged <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}
```

---

## 练习任务

### 任务 1：订单处理系统 - 生产者消费者（难度：中等）

```go
// 模拟电商订单处理系统

type Order struct {
    ID        int
    UserID    int
    Amount    float64
    Status    string // "pending", "processing", "completed", "failed"
    CreatedAt time.Time
}

// 1. 实现订单生产者
//    - 每隔 100-300ms（随机）生成一个订单
//    - 订单金额随机 10.00 - 999.99
//    - 总共生成 20 个订单后停止
func orderProducer(orders chan<- Order) {
    // 你的代码
}

// 2. 实现订单消费者
//    - 处理订单（模拟耗时 200-500ms）
//    - 金额 > 500 的订单有 20% 概率失败
//    - 打印处理结果
func orderConsumer(id int, orders <-chan Order, results chan<- Order) {
    // 你的代码
}

// 3. 在 main 中启动 1 个生产者和 3 个消费者
//    - 收集所有处理结果
//    - 最后打印统计：成功数、失败数、总金额、平均处理时间
```

### 任务 2：并发下载器 - Worker Pool（难度：中等偏上）

```go
// 实现一个并发下载器，限制最大并发数

type DownloadTask struct {
    URL      string
    FileName string
}

type DownloadResult struct {
    Task     DownloadTask
    Size     int           // 模拟下载大小（字节）
    Duration time.Duration
    Err      error
}

// 实现 Worker Pool 下载器
// maxWorkers: 最大并发下载数
func Download(tasks []DownloadTask, maxWorkers int) []DownloadResult {
    // 1. 创建 jobs channel 和 results channel
    // 2. 启动 maxWorkers 个 worker goroutine
    // 3. 发送所有任务到 jobs channel
    // 4. 收集所有结果
    // 5. 返回结果切片

    // 每个 worker 的工作：
    // - 从 jobs channel 接收任务
    // - 模拟下载（随机耗时 100ms-1s，随机大小 1KB-10MB）
    // - 10% 概率下载失败
    // - 将结果发送到 results channel
}

// 测试：10 个下载任务，最大 3 个并发
tasks := []DownloadTask{
    {URL: "https://example.com/file1.zip", FileName: "file1.zip"},
    {URL: "https://example.com/file2.zip", FileName: "file2.zip"},
    // ... 更多任务
}

results := Download(tasks, 3)
// 打印每个任务的结果和总体统计
```

### 任务 3：Select 超时控制器（难度：中等）

```go
// 实现一个多服务健康检查器，带超时控制

type ServiceStatus struct {
    Name     string
    Healthy  bool
    Latency  time.Duration
    ErrorMsg string
}

// 模拟检查单个服务的健康状态
func checkService(name string) ServiceStatus {
    // 模拟不同服务的响应时间：
    // "auth-service": 100-200ms
    // "user-service": 50-150ms
    // "order-service": 200-800ms（偶尔很慢）
    // "payment-service": 300-1500ms（经常超时）
}

// 实现带超时的健康检查
// timeout: 单个服务的超时时间
func HealthCheck(services []string, timeout time.Duration) []ServiceStatus {
    // 1. 并发检查所有服务
    // 2. 每个服务有独立的超时控制（用 select + time.After）
    // 3. 超时的服务标记为 Healthy: false, ErrorMsg: "timeout"
    // 4. 返回所有服务的状态
}

// 测试：
services := []string{"auth-service", "user-service", "order-service", "payment-service"}
results := HealthCheck(services, 500*time.Millisecond)
// 打印结果表格
```

### 任务 4：Fan-in 数据聚合（难度：中等偏上）

```go
// 实现数据流聚合：从 3 个不同的数据源收集数据，合并到一个 channel

type DataPoint struct {
    Source    string
    Value     float64
    Timestamp time.Time
}

// 模拟数据源：每个数据源以不同频率产生数据
func dataSource(name string, interval time.Duration, count int) <-chan DataPoint {
    ch := make(chan DataPoint)
    go func() {
        defer close(ch)
        for i := 0; i < count; i++ {
            time.Sleep(interval)
            ch <- DataPoint{
                Source:    name,
                Value:     rand.Float64() * 100,
                Timestamp: time.Now(),
            }
        }
    }()
    return ch
}

// 实现 Fan-in 合并
func mergeDataSources(sources ...<-chan DataPoint) <-chan DataPoint {
    // 将多个数据源合并为一个 channel
    // 所有数据源关闭后，合并的 channel 也要关闭
}

// 实现数据消费：实时打印 + 统计
func processData(merged <-chan DataPoint) {
    // 1. 实时打印每条数据
    // 2. 统计每个数据源的：数据条数、平均值、最大值、最小值
    // 3. 所有数据处理完后打印统计表格
}

// main 中：
// 数据源 A：每 100ms 产生一条，共 10 条
// 数据源 B：每 200ms 产生一条，共 8 条
// 数据源 C：每 300ms 产生一条，共 5 条
```

---

## 参考答案骨架

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// ============ 订单处理 ============

type Order struct {
    ID        int
    UserID    int
    Amount    float64
    Status    string
    CreatedAt time.Time
}

func orderProducer(orders chan<- Order) {
    defer close(orders)
    for i := 1; i <= 20; i++ {
        time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
        orders <- Order{
            ID:        i,
            UserID:    rand.Intn(100) + 1,
            Amount:    float64(rand.Intn(99000)+1000) / 100.0,
            Status:    "pending",
            CreatedAt: time.Now(),
        }
        fmt.Printf("[生产者] 生成订单 #%d, 金额 ¥%.2f\n", i, float64(rand.Intn(99000)+1000)/100.0)
    }
}

func orderConsumer(id int, orders <-chan Order, results chan<- Order, wg *sync.WaitGroup) {
    defer wg.Done()
    for order := range orders {
        // 模拟处理
        processTime := time.Duration(200+rand.Intn(300)) * time.Millisecond
        time.Sleep(processTime)

        // 大额订单有概率失败
        if order.Amount > 500 && rand.Float64() < 0.2 {
            order.Status = "failed"
        } else {
            order.Status = "completed"
        }

        fmt.Printf("[消费者%d] 订单 #%d %s (耗时 %v)\n", id, order.ID, order.Status, processTime)
        results <- order
    }
}

// ... 补充 main 函数和其他模式的实现
```

---

## 常见陷阱

1. **向已关闭的 channel 发送数据会 panic**
2. **从已关闭的 channel 接收会立即返回零值**
3. **忘记关闭 channel 导致 goroutine 泄漏**
4. **并发读写 map 需要加锁**（用 `sync.Mutex` 或 `sync.Map`）
5. **goroutine 中的闭包变量捕获问题**（循环变量要作为参数传入）

---

## 自检清单

- [ ] 理解缓冲 channel 和非缓冲 channel 的区别
- [ ] 能正确实现生产者-消费者模式
- [ ] 能用 Worker Pool 限制并发数
- [ ] 理解 Fan-in 模式的应用场景
- [ ] 会用 select 实现超时控制
- [ ] 代码中没有 goroutine 泄漏
- [ ] 代码中没有向已关闭 channel 发送数据的情况
