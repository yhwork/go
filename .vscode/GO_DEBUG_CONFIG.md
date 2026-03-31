# ✅ Go 调试环境配置完成！

## 🎉 已完成的配置

### 1. Delve 调试器已安装 ✓

```bash
$ dlv version
Delve Debugger
Version: 1.26.1
Build: $Id: 3f95fba2a798b133eda85dd54b3b000c4f8ba68a $
```

### 2. VS Code launch.json 已配置 ✓

位置：`d:\yang\go\.vscode\launch.json`

包含以下 8 个调试配置：

| # | 配置名称 | 用途 | 模式 |
|---|---------|------|------|
| 1 | ▶️ 运行当前文件 | 调试当前打开的 `.go` 文件 | auto |
| 2 | 🚀 运行整个项目 | 调试整个 Go 模块 | auto |
| 3 | 🧪 运行当前测试文件 | 运行并调试 `_test.go` 文件 | test |
| 4 | 🧪 运行单个测试函数 | 运行指定的 Test 函数 | test |
| 5 | 🔍 调试带参数的程序 | 带命令行参数调试 | auto |
| 6 | 🌐 运行 Web 服务器 | 调试 Gin/Echo 等 Web 应用 | auto |
| 7 | 🔧 附加到进程 | 附加到已运行的 Go 进程 | attach |
| 8 | 🐛 Delve 直接调试 | 使用 Delve 详细日志调试 | debug |

### 3. 示例代码已修复 ✓

文件：`d:\yang\go\go_demo\day01\main.go`

修复内容：
- 添加 `package main` 声明
- 结构体字段首字母大写（导出）
- 语法错误修复（逗号、冒号）
- 使用 `fmt.Println` 替代 `println`

## 🚀 快速开始调试

### 方法一：使用鼠标

1. 打开 `d:\yang\go\go_demo\day01\main.go`
2. 在代码行号左侧点击，添加断点（红点）
3. 按 `F5` 或点击左侧"运行和调试"图标
4. 选择 "▶️ 运行当前文件"
5. 程序会在断点处暂停

### 方法二：使用快捷键

| 操作 | 快捷键 | 说明 |
|------|--------|------|
| 启动调试 | `F5` | 开始/继续调试 |
| 单步跳过 | `F10` | 执行下一行代码 |
| 单步进入 | `F11` | 进入函数内部 |
| 单步跳出 | `Shift+F11` | 跳出当前函数 |
| 重启调试 | `Ctrl+Shift+F5` | 重新启动 |
| 停止调试 | `Shift+F5` | 停止调试 |
| 切换断点 | `F9` | 添加/移除断点 |

### 方法三：使用命令面板

1. `Ctrl+Shift+P` 打开命令面板
2. 输入 "Debug: Start Debugging"
3. 选择调试配置

## 📝 调试配置示例

### 调试简单的 Go 程序

```json
{
    "name": "▶️ 运行当前文件",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${fileDirname}",
    "console": "integratedTerminal"
}
```

### 调试 Web 服务器（Gin）

```json
{
    "name": "🌐 运行 Web 服务器 (Gin)",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${fileDirname}",
    "env": {
        "GIN_MODE": "debug",
        "PORT": "8080"
    }
}
```

### 调试测试

```json
{
    "name": "🧪 运行当前测试文件",
    "type": "go",
    "request": "launch",
    "mode": "test",
    "program": "${fileDirname}"
}
```

## 🔍 调试技巧

### 1. 条件断点

右键点击断点 -> 输入条件：
```go
i == 10           // 当 i 等于 10 时中断
err != nil        // 当有错误时中断
name == "alice"   // 当 name 为 alice 时中断
```

### 2. 监视变量

在"运行和调试"侧边栏 -> "监视"部分：
- 点击 "+" 添加要监视的变量
- 可以输入任意表达式，如 `len(users)`

### 3. 调用堆栈

查看函数调用链：
- 显示当前函数的调用者
- 可以点击跳转到上层调用

### 4. 变量窗口

调试时会显示：
- **局部变量**: 当前作用域的变量
- **全局变量**: 包级变量
- **寄存器**: CPU 寄存器状态

### 5. 即时求值

在调试暂停时：
- 鼠标悬停在变量上查看值
- 在"监视"窗口输入表达式实时计算

## 🛠️ 常用调试场景

### 场景 1：查找 Bug

```go
func calculateSum(numbers []int) int {
    sum := 0
    for i, num := range numbers {  // ← 在这里设断点
        sum += num                  // ← F10 单步执行
        fmt.Printf("i=%d, num=%d, sum=%d\n", i, num, sum) // ← 观察输出
    }
    return sum
}
```

### 场景 2：理解代码流程

```go
// 在 main 函数入口设断点
// 然后一直 F10/F11 单步跟踪
func main() {
    data := loadData()      // F11 进入函数
    result := process(data) // F11 进入函数
    save(result)           // F11 进入函数
}
```

### 场景 3：性能问题排查

```go
// 使用调试器的性能分析功能
// 或者在关键位置记录时间
start := time.Now()
// ... 代码 ...
elapsed := time.Since(start)
fmt.Printf("耗时：%v\n", elapsed)
```

## ⚠️ 常见问题解决

### Q1: "Cannot find Delve debugger"

**解决方案：**
```bash
# 检查 dlv 是否在 PATH 中
where dlv

# 如果找不到，手动添加到 PATH
# dlv 默认安装在：C:\Users\你的用户名\go\bin\dlv.exe
```

### Q2: 调试时防火墙弹窗

**解决方案：**
- 允许 `dlv.exe` 通过 Windows 防火墙
- 或者临时关闭防火墙测试

### Q3: 端口被占用（Web 调试）

**解决方案：**
- 修改代码中的端口
- 或在 launch.json 中添加环境变量：
```json
"env": {
    "PORT": "8081"
}
```

### Q4: 无法调试测试

**确保：**
- 文件名以 `_test.go` 结尾
- 测试函数以 `Test` 开头（大写 T）
- 测试函数签名正确：`func TestXxx(t *testing.T)`

## 📚 下一步学习

1. ✅ 已完成：安装 Delve 调试器
2. ✅ 已完成：配置 launch.json
3. ✅ 已完成：修复示例代码
4. 📖 接下来：尝试调试 `go_demo/day01/main.go`
5. 📖 然后：完成学习计划 Day 01-06 的练习
6. 📖 最后：用调试器解决实际问题

## 🎯 立即尝试

```bash
# 1. 打开文件
# d:\yang\go\go_demo\day01\main.go

# 2. 在第 19 行（fmt.Println("hello world")）设置断点

# 3. 按 F5 启动调试

# 4. 观察变量和执行流程
```

---

**配置完成时间：** 2026-03-25  
**Go 版本：** go1.26.1 windows/amd64  
**Delve 版本：** 1.26.1  
**操作系统：** Windows

祝你调试愉快！🐛 → ✅
