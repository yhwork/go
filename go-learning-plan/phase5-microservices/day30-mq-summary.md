# Day 30 - 消息队列与总结

## 今日目标

集成消息队列（RabbitMQ）实现异步通知，回顾 30 天所学，梳理知识图谱。

---

## 上午：消息队列实践

### 任务 1：RabbitMQ 集成

```go
// 添加到 docker-compose.yml
// rabbitmq:
//   image: rabbitmq:3-management-alpine
//   ports:
//     - "5672:5672"   # AMQP
//     - "15672:15672" # 管理界面
//   environment:
//     RABBITMQ_DEFAULT_USER: guest
//     RABBITMQ_DEFAULT_PASS: guest

// pkg/mq/rabbitmq.go
import "github.com/rabbitmq/amqp091-go"

type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(url) // "amqp://guest:guest@localhost:5672/"
    if err != nil {
        return nil, err
    }
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    return &RabbitMQ{conn: conn, channel: ch}, nil
}

// 发布消息
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
    return r.channel.PublishWithContext(
        context.Background(),
        exchange,
        routingKey,
        false, false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}

// 消费消息
func (r *RabbitMQ) Consume(queue string, handler func([]byte) error) error {
    msgs, err := r.channel.Consume(queue, "", false, false, false, false, nil)
    if err != nil {
        return err
    }
    for msg := range msgs {
        if err := handler(msg.Body); err != nil {
            msg.Nack(false, true) // 处理失败，重新入队
        } else {
            msg.Ack(false)
        }
    }
    return nil
}

func (r *RabbitMQ) Close() {
    r.channel.Close()
    r.conn.Close()
}
```

### 任务 2：异步通知流程

```go
// 发布文章 -> 通过 MQ 异步通知关注者

// 1. Article Service 发布文章后发送消息
func (s *articleService) Publish(userID, articleID uint) error {
    // 更新文章状态
    // ...

    // 发送消息到 MQ
    event := map[string]interface{}{
        "type":       "article_published",
        "article_id": articleID,
        "author_id":  userID,
        "timestamp":  time.Now(),
    }
    data, _ := json.Marshal(event)
    return s.mq.Publish("blog_events", "article.published", data)
}

// 2. Notification Service 消费消息并发送通知
func handleArticlePublished(body []byte) error {
    var event struct {
        ArticleID uint `json:"article_id"`
        AuthorID  uint `json:"author_id"`
    }
    json.Unmarshal(body, &event)

    // 查询作者的所有关注者
    // 为每个关注者创建通知
    // 通过 WebSocket 推送实时通知
    return nil
}
```

### 任务 3：声明交换机和队列

```go
func setupQueues(ch *amqp.Channel) error {
    // 声明交换机
    err := ch.ExchangeDeclare("blog_events", "topic", true, false, false, false, nil)
    if err != nil {
        return err
    }

    // 声明队列
    q, err := ch.QueueDeclare("notification_queue", true, false, false, false, nil)
    if err != nil {
        return err
    }

    // 绑定队列到交换机
    return ch.QueueBind(q.Name, "article.*", "blog_events", false, nil)
}
```

---

## 下午：30 天总结

### 知识图谱

```
Go 全栈开发知识体系
│
├── Go 语言核心
│   ├── 结构体与方法
│   ├── 接口与多态
│   ├── 错误处理（自定义错误、包装、errors.Is/As）
│   ├── 并发编程（goroutine、channel、select、context）
│   ├── 并发模式（生产者消费者、Worker Pool、Fan-in/out）
│   ├── 包管理（Go Modules）
│   └── 测试（表驱动、Mock、Benchmark）
│
├── Web 开发（Gin）
│   ├── 路由与路由组
│   ├── 中间件（Logger、Recovery、CORS、限流）
│   ├── 请求绑定与验证
│   └── 统一响应格式
│
├── 数据库（GORM）
│   ├── 模型定义与关联
│   ├── CRUD 操作
│   ├── 预加载与复杂查询
│   ├── 事务处理
│   └── Scope 封装
│
├── 缓存（Redis）
│   ├── 基本数据类型操作
│   ├── Cache-Aside 模式
│   ├── 限流
│   └── JWT 黑名单
│
├── 认证与安全
│   ├── JWT（双 Token 机制）
│   ├── bcrypt 密码加密
│   ├── XSS/SQL 注入防护
│   └── 接口限流
│
├── 实时通信
│   ├── WebSocket
│   └── 消息推送
│
├── 异步处理
│   ├── Cron 定时任务
│   ├── Redis 任务队列
│   └── RabbitMQ 消息队列
│
├── 微服务
│   ├── Protocol Buffers
│   ├── gRPC
│   ├── API Gateway
│   └── 服务拆分
│
├── 工程化
│   ├── 项目结构（三层架构）
│   ├── 依赖注入
│   ├── 配置管理（Viper）
│   ├── 结构化日志（Zap）
│   ├── Swagger 文档
│   └── Makefile
│
└── 部署
    ├── Docker（多阶段构建）
    ├── Docker Compose
    ├── Nginx 反向代理
    └── 优雅关闭
```

### 回顾与自评

逐个阶段回顾，给自己打分（1-5）：

| 阶段 | 内容 | 掌握度 | 需要加强 |
|------|------|--------|----------|
| Phase 1 | Go 进阶 | /5 | |
| Phase 2 | Gin + DB | /5 | |
| Phase 3 | 进阶功能 | /5 | |
| Phase 4 | 全栈项目 | /5 | |
| Phase 5 | 微服务 | /5 | |

### 后续学习方向

完成这 30 天后，可以继续深入的方向：

1. **服务治理**：服务发现（Consul/etcd）、配置中心、链路追踪（Jaeger）
2. **高可用**：负载均衡、熔断器（Hystrix）、服务降级
3. **性能**：pprof 性能分析、内存优化、连接池调优
4. **监控**：Prometheus + Grafana 监控系统
5. **CI/CD**：GitHub Actions、自动化测试、自动部署
6. **云原生**：Kubernetes 部署、Helm Chart、Service Mesh

---

## 自检清单

- [ ] RabbitMQ 消息发布和消费正常
- [ ] 异步通知流程完整可用
- [ ] 知识图谱梳理完成
- [ ] 各阶段自评诚实客观
- [ ] 确定了后续学习方向
- [ ] 所有代码已提交 Git
