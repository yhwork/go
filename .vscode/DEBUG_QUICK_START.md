# 🚀 Go 调试快速入门

## ✅ 你的环境已就绪！

- ✓ Delve 调试器已安装 (v1.26.1)
- ✓ launch.json 已配置 (8 种调试模式)
- ✓ 示例代码已修复并可运行

## 📖 30 秒开始调试

### 第一步：打开文件
```
d:\yang\go\go_demo\day01\main.go
```

### 第二步：设置断点
在第 **19** 行左侧点击（`fmt.Println("hello world")`）

### 第三步：启动调试
按 **F5**

### 第四步：选择配置
选择 **"▶️ 运行当前文件"**

### 第五步：观察执行
程序会在断点处暂停，可以：
- 查看变量值
- 按 F10 继续下一行
- 按 F5 继续执行

## 🎹 快捷键速查

| 功能 | Windows 快捷键 |
|------|---------------|
| **启动/继续调试** | `F5` |
| **单步跳过** | `F10` |
| **单步进入** | `F11` |
| **单步跳出** | `Shift+F11` |
| **重启调试** | `Ctrl+Shift+F5` |
| **停止调试** | `Shift+F5` |
| **切换断点** | `F9` |

## 🎯 8 种调试模式说明

### 1️⃣ ▶️ 运行当前文件
**用途：** 调试当前打开的 `.go` 文件  
**场景：** 学习语法、测试小例子

### 2️⃣ 🚀 运行整个项目
**用途：** 调试整个 Go 模块  
**场景：** 有多个包的完整项目

### 3️⃣ 🧪 运行当前测试文件
**用途：** 调试 `_test.go` 文件  
**场景：** 调试单元测试

### 4️⃣ 🧪 运行单个测试函数
**用途：** 运行指定的 Test 函数  
**场景：** 只测试某个特定功能

### 5️⃣ 🔍 调试带参数的程序
**用途：** 传递命令行参数  
**场景：** `go run main.go --config config.yaml`

### 6️⃣ 🌐 运行 Web 服务器
**用途：** 调试 Gin/Echo 等 Web 应用  
**场景：** Web API 开发、热重载调试

### 7️⃣ 🔧 附加到进程
**用途：** 附加到已运行的 Go 进程  
**场景：** 调试后台服务、长时间运行的程序

### 8️⃣ 🐛 Delve 直接调试
**用途：** 使用 Delve 详细日志  
**场景：** 深入排查复杂问题

## 💡 实用技巧

### 技巧 1：打印调试
```go
fmt.Printf("Debug: i=%d, result=%v\n", i, result)
```

### 技巧 2：条件断点
右键断点 → 输入条件：
```go
i == 10
err != nil
user.Name == "alice"
```

### 技巧 3：监视表达式
在"监视"窗口添加：
```go
len(users)
total / count
data["key"]
```

### 技巧 4：查看调用栈
调试暂停时，查看：
- 当前函数
- 谁调用了这个函数
- 完整的调用链

## 🔧 常用命令

### 终端运行
```bash
# 运行 Go 文件
go run main.go

# 运行测试
go test -v

# 运行单个测试
go test -run TestMyFunction -v

# 构建可执行文件
go build -o app.exe main.go
```

### 调试器命令（Delve CLI）
```bash
# 启动调试
dlv debug

# 附加到进程
dlv attach <pid>

# 运行并调试
dlv exec ./app --arg1 arg2
```

## 📁 项目文件位置

| 文件 | 路径 |
|------|------|
| 调试配置 | `d:\yang\go\.vscode\launch.json` |
| 示例代码 | `d:\yang\go\go_demo\day01\main.go` |
| 学习计划 | `d:\yang\go\go-learning-plan\README.md` |
| 配置指南 | `d:\yang\go\.vscode\GO_DEBUG_CONFIG.md` |

## ⚠️ 故障排除

### 问题：找不到 dlv
```bash
# 检查安装
where dlv

# 重新安装
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 问题：无法启动调试
1. 确保文件保存了（Ctrl+S）
2. 确保没有语法错误
3. 重启 VS Code

### 问题：断点不生效
1. 检查是否是可执行代码（注释和空行不能设断点）
2. 确保是调试模式而不是运行模式
3. 重新编译：`Ctrl+Shift+B`

## 📚 学习路线

```
Day 1: 学会设置断点、单步执行 ← 你现在在这里
Day 2: 学会查看变量、监视表达式
Day 3: 学会条件断点、调用栈
Day 4: 学会调试测试
Day 5: 学会调试 Web 应用
Day 6: 综合运用解决实际问题
```

## 🎉 立即开始！

1. 打开 `go_demo/day01/main.go`
2. 在第 19 行设断点
3. 按 F5
4. 享受调试的乐趣！

---

**提示：** 更多详细内容请查看：
- `GO_DEBUG_CONFIG.md` - 完整配置指南
- `GO_DEBUG_SETUP.md` - 环境搭建说明

Happy Debugging! 🐛✨
