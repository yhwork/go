# Go 全栈开发 30 天学习计划

> 针对有前端开发经验、了解 Go 基础语法的开发者
> 技术栈：Go + Gin + GORM + MySQL + Redis + Docker + gRPC

---

## 目录结构

```
go-learning-plan/
├── README.md                          # 本文件 - 学习计划总览
│
├── phase1-go-advanced/                # 第一阶段：Go 进阶核心（Day 1-6）
│   ├── day01-struct-methods.md        # 结构体与方法
│   ├── day02-interface-polymorphism.md # 接口与多态
│   ├── day03-error-concurrency.md     # 错误处理与并发基础
│   ├── day04-channel-patterns.md      # Channel 与并发模式
│   ├── day05-modules-project-structure.md # 包管理与项目结构
│   └── day06-testing-benchmark.md     # 单元测试与基准测试
│
├── phase2-gin-database/               # 第二阶段：Gin + 数据库 CRUD（Day 7-14）
│   ├── day07-gin-routing.md           # Gin 入门与路由
│   ├── day08-middleware.md            # 中间件机制
│   ├── day09-validation-binding.md    # 请求验证与绑定
│   ├── day10-gorm-basics.md           # GORM 入门与模型定义
│   ├── day11-gorm-advanced.md         # GORM 进阶查询
│   ├── day12-layered-architecture.md  # 三层架构整合
│   ├── day13-jwt-auth.md              # JWT 认证
│   └── day14-redis-cache.md           # Redis 缓存
│
├── phase3-advanced-features/          # 第三阶段：进阶功能（Day 15-20）
│   ├── day15-file-upload.md           # 文件上传与静态资源
│   ├── day16-swagger-logging.md       # Swagger 文档与日志
│   ├── day17-websocket.md             # WebSocket 实时通信
│   ├── day18-cron-background-tasks.md # 定时任务与后台任务
│   ├── day19-docker.md                # Docker 容器化
│   └── day20-grpc.md                  # gRPC 入门
│
├── phase4-fullstack-project/          # 第四阶段：全栈实战项目（Day 21-28）
│   ├── day21-project-init.md          # 项目初始化与数据库设计
│   ├── day22-user-module.md           # 用户模块完整实现
│   ├── day23-article-module.md        # 文章模块完整实现
│   ├── day24-category-tag-comment.md  # 分类、标签、评论模块
│   ├── day25-like-favorite-notification.md # 点赞、收藏与通知
│   ├── day26-frontend-integration.md  # 前后端联调
│   ├── day27-optimization-security.md # 性能优化与安全加固
│   └── day28-deployment.md            # 部署上线
│
└── phase5-microservices/              # 第五阶段：微服务入门（Day 29-30）
    ├── day29-microservice-split.md    # 微服务拆分
    └── day30-mq-summary.md            # 消息队列与总结
```

---

## 学习进度跟踪

| Day | 主题 | 状态 |
|-----|------|------|
| **Phase 1: Go 进阶核心** | | |
| Day 01 | 结构体与方法 | ⬜ |
| Day 02 | 接口与多态 | ⬜ |
| Day 03 | 错误处理与并发基础 | ⬜ |
| Day 04 | Channel 与并发模式 | ⬜ |
| Day 05 | 包管理与项目结构 | ⬜ |
| Day 06 | 单元测试与基准测试 | ⬜ |
| **Phase 2: Gin + 数据库** | | |
| Day 07 | Gin 入门与路由 | ⬜ |
| Day 08 | 中间件机制 | ⬜ |
| Day 09 | 请求验证与绑定 | ⬜ |
| Day 10 | GORM 入门与模型定义 | ⬜ |
| Day 11 | GORM 进阶查询 | ⬜ |
| Day 12 | 三层架构整合 | ⬜ |
| Day 13 | JWT 认证 | ⬜ |
| Day 14 | Redis 缓存 | ⬜ |
| **Phase 3: 进阶功能** | | |
| Day 15 | 文件上传与静态资源 | ⬜ |
| Day 16 | Swagger 文档与日志 | ⬜ |
| Day 17 | WebSocket 实时通信 | ⬜ |
| Day 18 | 定时任务与后台任务 | ⬜ |
| Day 19 | Docker 容器化 | ⬜ |
| Day 20 | gRPC 入门 | ⬜ |
| **Phase 4: 全栈实战项目** | | |
| Day 21 | 项目初始化与数据库设计 | ⬜ |
| Day 22 | 用户模块完整实现 | ⬜ |
| Day 23 | 文章模块完整实现 | ⬜ |
| Day 24 | 分类、标签、评论模块 | ⬜ |
| Day 25 | 点赞、收藏与通知 | ⬜ |
| Day 26 | 前后端联调 | ⬜ |
| Day 27 | 性能优化与安全加固 | ⬜ |
| Day 28 | 部署上线 | ⬜ |
| **Phase 5: 微服务入门** | | |
| Day 29 | 微服务拆分 | ⬜ |
| Day 30 | 消息队列与总结 | ⬜ |

> 完成一天后将 ⬜ 改为 ✅

---

## 每日学习节奏

| 时间段 | 内容 | 时长 |
|--------|------|------|
| 前 30 分钟 | 复习昨天的知识点，review 昨天的代码 | 30min |
| 中间 1.5 小时 | 学习新知识点 + 编写练习代码 | 90min |
| 最后 30 分钟 | 写学习笔记 + 提交代码到 Git | 30min |

---

## 技术栈一览

| 类别 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.22+ | 后端开发 |
| Web 框架 | Gin | HTTP 路由、中间件 |
| ORM | GORM | 数据库操作 |
| 数据库 | MySQL / PostgreSQL | 数据持久化 |
| 缓存 | Redis | 缓存、限流、队列 |
| 认证 | JWT (golang-jwt) | 用户认证 |
| 日志 | Zap | 结构化日志 |
| 配置 | Viper | 配置管理 |
| 文档 | Swag | Swagger API 文档 |
| 测试 | testify | 测试断言 |
| RPC | gRPC + Protobuf | 微服务通信 |
| 消息队列 | RabbitMQ | 异步任务 |
| 容器 | Docker + Compose | 部署 |
| 代理 | Nginx | 反向代理 |

---

## 关键原则

1. **每天必须写代码**，不要只看不练
2. **所有代码提交 Git**，养成版本管理习惯
3. **遇到问题先查文档**，实在解决不了再问
4. **每个阶段结束做一次 code review**，重构不优雅的部分
5. **先跑通再优化**，不要一开始就追求完美
