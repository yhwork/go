# Day 05 - 包管理与项目结构

## 今日目标

掌握 Go Modules 依赖管理、标准项目目录结构、包的可见性规则，以及使用 Viper 进行配置管理。

---

## 知识点

### 1. Go Modules

```bash
# 初始化模块
go mod init github.com/yourname/myapp

# 添加依赖
go get github.com/spf13/viper

# 整理依赖（移除未使用的）
go mod tidy

# 查看依赖树
go mod graph

# 下载依赖到本地缓存
go mod download
```

### 2. go.mod 文件

```go
module github.com/yourname/myapp

go 1.22

require (
    github.com/spf13/viper v1.18.2
    github.com/gin-gonic/gin v1.9.1
)

// indirect 表示间接依赖
require (
    github.com/fsnotify/fsnotify v1.7.0 // indirect
)
```

### 3. 标准项目结构

```
myapp/
├── cmd/                    # 应用入口
│   └── server/
│       └── main.go         # main 函数
├── internal/               # 私有代码（其他模块不可导入）
│   ├── config/             # 配置管理
│   │   └── config.go
│   ├── handler/            # HTTP 处理器（Controller 层）
│   │   ├── user_handler.go
│   │   └── product_handler.go
│   ├── service/            # 业务逻辑（Service 层）
│   │   ├── user_service.go
│   │   └── product_service.go
│   ├── repository/         # 数据访问（DAO 层）
│   │   ├── user_repo.go
│   │   └── product_repo.go
│   ├── model/              # 数据模型
│   │   ├── user.go
│   │   └── product.go
│   ├── middleware/          # HTTP 中间件
│   │   └── auth.go
│   └── router/             # 路由注册
│       └── router.go
├── pkg/                    # 可被外部导入的公共包
│   ├── response/           # 统一响应格式
│   │   └── response.go
│   └── utils/              # 工具函数
│       └── hash.go
├── config/                 # 配置文件
│   ├── config.yaml
│   └── config.example.yaml
├── migrations/             # 数据库迁移文件
├── docs/                   # 文档
├── scripts/                # 脚本
├── go.mod
├── go.sum
├── Makefile
└── .gitignore
```

### 4. 包的可见性

```go
package user

// 首字母大写 -> 导出（public）
type User struct {
    ID       uint
    Username string
}

func NewUser(name string) *User {
    return &User{Username: name}
}

// 首字母小写 -> 未导出（private，包内可见）
type userCache struct {
    data map[uint]*User
}

func validate(u *User) error {
    // 只能在 user 包内调用
    return nil
}
```

### 5. internal 目录的特殊性

```
myapp/
├── internal/
│   └── service/
│       └── user_service.go  # 只有 myapp 模块内的代码可以导入
```

`internal` 目录下的包只能被其父目录及父目录的子目录导入，外部模块无法导入。这是 Go 工具链强制执行的规则。

---

## 练习任务

### 任务 1：创建项目骨架（难度：基础）

```bash
# 1. 创建项目目录并初始化
mkdir -p go-blog && cd go-blog
go mod init github.com/yourname/go-blog

# 2. 创建完整的目录结构
mkdir -p cmd/server
mkdir -p internal/{config,handler,service,repository,model,middleware,router}
mkdir -p pkg/{response,utils}
mkdir -p config
mkdir -p migrations
```

创建以下文件并填充基础代码：

```go
// cmd/server/main.go
package main

import (
    "fmt"
    "github.com/yourname/go-blog/internal/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }
    fmt.Printf("服务启动: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
```

### 任务 2：实现配置管理（难度：中等）

使用 Viper 实现灵活的配置管理：

```yaml
# config/config.yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release, test

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "123456"
  dbname: "go_blog"
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expire: 24  # 小时

log:
  level: "info"    # debug, info, warn, error
  format: "console" # console, json
  output: "stdout"  # stdout, file
  file_path: "logs/app.log"
```

```go
// internal/config/config.go
package config

import "github.com/spf13/viper"

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
    Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
    Driver       string `mapstructure:"driver"`
    Host         string `mapstructure:"host"`
    Port         int    `mapstructure:"port"`
    Username     string `mapstructure:"username"`
    Password     string `mapstructure:"password"`
    DBName       string `mapstructure:"dbname"`
    MaxIdleConns int    `mapstructure:"max_idle_conns"`
    MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// ... 补充 RedisConfig, JWTConfig, LogConfig

// Load 加载配置
// 优先级：环境变量 > config.yaml > 默认值
func Load() (*Config, error) {
    // 1. 设置默认值
    // 2. 设置配置文件路径
    // 3. 绑定环境变量（如 SERVER_PORT -> server.port）
    // 4. 读取配置文件
    // 5. 反序列化到 Config 结构体
}

// DatabaseConfig 的 DSN 方法
func (c DatabaseConfig) DSN() string {
    // 返回数据库连接字符串
    // MySQL: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
}
```

### 任务 3：实现统一响应格式（难度：基础）

```go
// pkg/response/response.go
package response

// Response 统一响应格式
type Response struct {
    Code    int         `json:"code"`    // 业务状态码
    Message string      `json:"message"` // 提示信息
    Data    interface{} `json:"data"`    // 数据
}

// PageData 分页数据
type PageData struct {
    List     interface{} `json:"list"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"page_size"`
}

// 业务状态码定义
const (
    CodeSuccess       = 0
    CodeBadRequest    = 400
    CodeUnauthorized  = 401
    CodeForbidden     = 403
    CodeNotFound      = 404
    CodeInternalError = 500
)

// 预定义消息
var codeMessages = map[int]string{
    CodeSuccess:       "success",
    CodeBadRequest:    "请求参数错误",
    CodeUnauthorized:  "未授权",
    CodeForbidden:     "禁止访问",
    CodeNotFound:      "资源不存在",
    CodeInternalError: "服务器内部错误",
}

func Success(data interface{}) Response {
    // 返回成功响应
}

func Error(code int, msg ...string) Response {
    // 返回错误响应，msg 可选，不传则使用 codeMessages 中的默认消息
}

func PageSuccess(list interface{}, total int64, page, pageSize int) Response {
    // 返回分页成功响应
}
```

### 任务 4：实现工具函数包（难度：基础）

```go
// pkg/utils/hash.go
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行 bcrypt 加密
func HashPassword(password string) (string, error) {}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {}

// pkg/utils/string.go

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {}

// MaskPhone 手机号脱敏: 13812345678 -> 138****5678
func MaskPhone(phone string) string {}

// MaskEmail 邮箱脱敏: alice@example.com -> ali***@example.com
func MaskEmail(email string) string {}
```

### 任务 5：编写 Makefile（难度：基础）

```makefile
# Makefile

.PHONY: build run test lint clean

# 构建
build:
	go build -o bin/server ./cmd/server

# 运行
run:
	go run ./cmd/server

# 测试
test:
	go test ./... -v -cover

# 代码检查
lint:
	golangci-lint run

# 清理
clean:
	rm -rf bin/

# 格式化
fmt:
	gofmt -w .

# 依赖整理
tidy:
	go mod tidy
```

---

## 参考答案骨架

```go
// internal/config/config.go
package config

import (
    "fmt"
    "strings"

    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

// ... 结构体定义

func Load() (*Config, error) {
    v := viper.New()

    // 设置默认值
    v.SetDefault("server.host", "0.0.0.0")
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.mode", "debug")

    // 配置文件
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath("./config")
    v.AddConfigPath(".")

    // 环境变量
    v.SetEnvPrefix("APP")
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    v.AutomaticEnv()

    // 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("解析配置失败: %w", err)
    }

    return &cfg, nil
}

func (c DatabaseConfig) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        c.Username, c.Password, c.Host, c.Port, c.DBName)
}
```

---

## 自检清单

- [ ] 理解 Go Modules 的工作原理
- [ ] 项目目录结构清晰，职责分明
- [ ] 配置管理支持文件 + 环境变量，优先级正确
- [ ] 统一响应格式定义完整
- [ ] internal 和 pkg 目录的区别理解清楚
- [ ] 包的可见性规则（大小写）掌握牢固
- [ ] Makefile 可用于日常开发操作
