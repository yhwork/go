# Day 18 - 定时任务与后台任务

## 今日目标

使用 robfig/cron 实现定时任务、基于 Redis 实现简单任务队列、异步任务处理、优雅关闭。

---

## 练习任务

### 任务 1：Cron 定时任务

```go
// internal/cron/cron.go
import "github.com/robfig/cron/v3"

type CronManager struct {
    c *cron.Cron
}

func NewCronManager() *CronManager {
    c := cron.New(cron.WithSeconds()) // 支持秒级精度
    return &CronManager{c: c}
}

func (m *CronManager) Start() {
    // 每天凌晨 2 点清理过期 Token
    m.c.AddFunc("0 0 2 * * *", cleanExpiredTokens)

    // 每小时统计在线用户数
    m.c.AddFunc("0 0 * * * *", logOnlineStats)

    // 每 5 分钟刷新热门文章缓存
    m.c.AddFunc("0 */5 * * * *", refreshHotArticlesCache)

    m.c.Start()
}

func (m *CronManager) Stop() {
    ctx := m.c.Stop() // 返回 context，等待正在运行的任务完成
    <-ctx.Done()
}
```

### 任务 2：Redis 任务队列

```go
// internal/queue/queue.go

type Task struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

type Queue struct {
    rdb   *redis.Client
    name  string
}

func (q *Queue) Enqueue(ctx context.Context, task Task) error {
    data, _ := json.Marshal(task)
    return q.rdb.LPush(ctx, q.name, data).Err()
}

func (q *Queue) Dequeue(ctx context.Context, timeout time.Duration) (*Task, error) {
    result, err := q.rdb.BRPop(ctx, timeout, q.name).Result()
    if err != nil {
        return nil, err
    }
    var task Task
    json.Unmarshal([]byte(result[1]), &task)
    return &task, nil
}

// Worker 消费任务
func (q *Queue) StartWorker(ctx context.Context, handler func(Task) error) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            task, err := q.Dequeue(ctx, 5*time.Second)
            if err != nil {
                continue
            }
            if err := handler(*task); err != nil {
                logger.Log.Errorw("任务处理失败", "type", task.Type, "error", err)
                // 可以将失败任务放入死信队列
            }
        }
    }
}
```

### 任务 3：异步邮件通知

```go
// 注册成功后异步发送欢迎邮件
func (s *userService) Register(req *dto.RegisterReq) (*dto.UserResponse, error) {
    // ... 创建用户 ...

    // 异步发送欢迎邮件
    s.queue.Enqueue(context.Background(), queue.Task{
        Type:    "send_email",
        Payload: json.RawMessage(fmt.Sprintf(`{"to":"%s","template":"welcome","username":"%s"}`, user.Email, user.Username)),
    })

    return dto.ToUserResponse(user), nil
}

// 邮件任务处理器
func handleSendEmail(task queue.Task) error {
    var payload struct {
        To       string `json:"to"`
        Template string `json:"template"`
        Username string `json:"username"`
    }
    json.Unmarshal(task.Payload, &payload)
    // 调用邮件发送 SDK（此处模拟打印）
    fmt.Printf("发送邮件: to=%s, template=%s\n", payload.To, payload.Template)
    return nil
}
```

### 任务 4：优雅关闭

```go
// cmd/server/main.go

func main() {
    // ... 初始化 ...

    // 启动 HTTP 服务（非阻塞）
    srv := &http.Server{Addr: addr, Handler: r}
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("服务启动失败: %v", err)
        }
    }()

    // 启动 Cron
    cronManager.Start()

    // 启动队列 Worker
    ctx, cancel := context.WithCancel(context.Background())
    go queue.StartWorker(ctx, taskHandler)

    // 等待退出信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("正在关闭服务...")

    // 1. 停止接收新请求
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    srv.Shutdown(shutdownCtx)

    // 2. 停止 Cron
    cronManager.Stop()

    // 3. 停止队列 Worker
    cancel()

    // 4. 关闭数据库连接
    sqlDB, _ := db.DB()
    sqlDB.Close()

    // 5. 关闭 Redis 连接
    rdb.Close()

    log.Println("服务已安全关闭")
}
```

---

## 自检清单

- [ ] Cron 定时任务按预期执行
- [ ] Redis 队列入队和出队正常
- [ ] 注册后异步邮件任务被正确处理
- [ ] Ctrl+C 时所有组件依次优雅关闭
- [ ] 没有 goroutine 泄漏
