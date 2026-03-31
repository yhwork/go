# 🔧 调试错误解决方案

## ❌ 错误信息

```
Build Error: go build -o D:\yang\go\go_demo\day01\debug\bin.exe ...
file not found in current directory or any parent directory; see 'go help modules'
(exit status 1)
```

## ✅ 解决方案

### 原因
缺少 `go.mod` 文件。Go 1.16+ 版本要求使用 Go Modules 管理项目。

### 解决方法

#### 方法一：初始化 go.mod（推荐）

在包含 `main.go` 的目录中执行：

```bash
cd d:\yang\go\go_demo\day01
go mod init day01
```

然后重新运行或调试即可。

#### 方法二：在父目录创建 go.mod

如果你希望整个 `go_demo` 是一个模块：

```bash
cd d:\yang\go\go_demo
go mod init go-demo
```

然后在每个子目录中都可以正常运行。

#### 方法三：临时禁用模块（不推荐）

```bash
$env:GO111MODULE="off"
go run main.go
```

但这会导致无法使用现代 Go 的最佳实践。

## 📝 已执行的修复

✅ 已在 `d:\yang\go\go_demo\day01` 目录创建 `go.mod`

验证：
```bash
cd d:\yang\go\go_demo\day01
go run main.go  # ✓ 成功输出
```

## 🎯 项目结构建议

### 推荐结构 A：每个练习独立模块

```
go_demo/
├── day01/
│   ├── main.go
│   └── go.mod      # module day01
├── day02/
│   ├── main.go
│   └── go.mod      # module day02
└── ...
```

**优点：** 每个练习独立，互不影响  
**缺点：** 多个 go.mod 文件

### 推荐结构 B：统一模块

```
go_demo/
├── go.mod          # module go-demo
├── day01/
│   └── main.go
├── day02/
│   └── main.go
└── ...
```

**优点：** 统一管理，代码共享方便  
**缺点：** 所有代码在同一个模块中

### 当前选择：结构 A（每个练习独立）

这样每个练习都是独立的，符合学习计划的设计。

## 🔍 验证步骤

### 1. 检查 go.mod 是否存在

```bash
ls go.mod
```

### 2. 检查 go.mod 内容

```bash
cat go.mod
```

应该看到类似：
```
module day01

go 1.26
```

### 3. 测试运行

```bash
go run main.go
```

### 4. 测试调试

在 VS Code 中按 F5，选择 "▶️ 运行当前文件"

## ⚠️ 其他可能的错误

### 错误：package is not in stdlib

**现象：**
```
package xxx is not in stdlib
```

**解决：**
确保使用了正确的模块路径，或者运行 `go mod tidy`。

### 错误：cannot find package

**现象：**
```
cannot find package "xxx"
```

**解决：**
```bash
go mod tidy
```

### 错误：build constraints exclude all Go files

**现象：**
通常发生在跨平台编译时。

**解决：**
确保文件名没有特殊构建标签，如 `file_windows.go` 在非 Windows 系统。

## 🛠️ 常用命令

```bash
# 初始化模块
go mod init <module-name>

# 整理依赖
go mod tidy

# 查看依赖
go list -m all

# 清理缓存
go clean -modcache

# 更新依赖
go get -u ./...
```

## 📚 参考资料

- [Go Modules 官方文档](https://go.dev/ref/mod)
- [go mod 命令详解](https://go.dev/cmd/go/#hdr-Module_maintenance)

---

**提示：** 如果遇到其他问题，请查看错误信息中的具体提示，通常 Go 的错误信息都很友好且包含解决方案。
