# Go 调试环境配置指南

## ✅ 已完成的配置

你的 `launch.json` 已经配置了以下调试模式：

| 配置名称 | 用途 | 快捷键 |
|---------|------|--------|
| ▶️ 运行当前文件 | 调试当前打开的 Go 文件 | F5 选择 |
| 🚀 运行整个项目 | 调试整个 Go 模块项目 | F5 选择 |
| 🧪 运行当前测试文件 | 运行并调试测试文件 | F5 选择 |
| 🧪 运行单个测试函数 | 运行指定的测试函数 | F5 选择 |
| 🔍 调试带参数的程序 | 带命令行参数调试 | F5 选择 |
| 🌐 运行 Web 服务器 (Gin) | 调试 Gin Web 应用 | F5 选择 |
| 🔧 附加到进程 | 附加到已运行的 Go 进程 | F5 选择 |
| 🐛 Delve 直接调试 | 使用 Delve 详细模式调试 | F5 选择 |

## ⚠️ 需要安装的调试工具

### 方法一：使用 Go 命令安装（推荐）

```bash
# 安装 Delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest

# 验证安装
dlv version
```

如果下载失败（网络问题），可以尝试设置 GOPROXY：

```bash
# 设置国内代理（中国用户）
$env:GOPROXY="https://goproxy.cn,direct"
go install github.com/go-delve/delve/cmd/dlv@latest

# 或使用全局代理
$env:GOPROXY="https://goproxy.io,direct"
```

### 方法二：手动下载

1. 访问 GitHub Releases: https://github.com/go-delve/delve/releases
2. 下载对应版本的 `dlv-windows-amd64.exe`
3. 重命名为 `dlv.exe`
4. 放到 `$GOPATH/bin` 目录（通常是 `C:\Users\你的用户名\go\bin`）

### 方法三：使用 Chocolatey（Windows）

```powershell
choco install delve
```

## 📝 使用示例

### 1. 调试简单的 Go 程序

打开 `go_demo/day01/main.go`，按 F5 选择 "▶️ 运行当前文件"

### 2. 设置断点

在代码行号左侧点击，会出现红点，表示断点。

### 3. 调试控制台

启动调试后，可以使用以下命令：
- **继续 (F5)**: 继续执行直到下一个断点
- **单步跳过 (F10)**: 执行下一行代码
- **单步进入 (F11)**: 进入函数内部
- **单步跳出 (Shift+F11)**: 跳出当前函数
- **重启 (Ctrl+Shift+F5)**: 重新启动调试会话
- **停止 (Shift+F5)**: 停止调试

### 4. 查看变量

调试时，鼠标悬停在变量上可以查看当前值，或在"运行和调试"侧边栏查看。

### 5. 条件断点

右键点击断点，输入条件表达式，例如：
```
i == 10
err != nil
```

## 🔧 常见问题

### Q: 提示 "Cannot find Delve debugger"
A: 确保 dlv.exe 已安装并且在 PATH 环境变量中。检查：
```bash
where dlv
```

### Q: 调试时卡住不动
A: 可能是防火墙阻止了 Delve。尝试：
1. 关闭防火墙临时测试
2. 允许 dlv.exe 通过防火墙

### Q: 无法调试测试
A: 确保测试文件名以 `_test.go` 结尾，测试函数以 `Test` 开头。

### Q: 调试 Web 服务器时端口被占用
A: 在 launch.json 中修改 PORT 环境变量，或在代码中修改默认端口。

## 📦 VS Code 扩展推荐

确保安装了以下扩展：
- **Go** (by Go Team at Google) - 官方 Go 支持
- **CodeLLDB** - 可选的调试器替代

## 🎯 下一步

1. 安装 Delve 调试器（参考上方方法）
2. 打开一个 Go 文件（如 `go_demo/day01/main.go`）
3. 设置断点
4. 按 F5 开始调试
5. 使用调试工具栏控制程序执行

---

配置完成后，你就可以享受完整的 Go 调试体验了！🚀
