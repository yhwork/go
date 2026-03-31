# 🔍 调试器问题排查指南

## ❓ 你看到的错误是什么？

如果看到类似 `"xxx" is not found. Please make sure it is installed and available` 的错误，请按以下步骤排查：

---

## ✅ 检查清单

### 1️⃣ 检查 Go 是否安装

```bash
go version
```

**期望输出：**
```
go version go1.26.1 windows/amd64
```

✅ **已确认：** 你的 Go 已安装（go1.26.1）

---

### 2️⃣ 检查 Delve 调试器是否安装

```bash
dlv version
```

**期望输出：**
```
Delve Debugger
Version: 1.26.1
```

✅ **已确认：** Delve 已安装（v1.26.1）

---

### 3️⃣ 检查 go.mod 是否存在

```bash
ls go.mod
```

✅ **已确认：** `d:\yang\go\go_demo\day01\go.mod` 存在

内容：
```
module day01

go 1.26.1
```

---

### 4️⃣ 检查代码是否可以运行

```bash
go run main.go
```

✅ **已确认：** 程序可以正常运行

输出：
```
hello world
1607
hangman
18
0 hangman
0 18
```

---

### 5️⃣ 检查 VS Code 配置

✅ **已创建：** `.vscode/settings.json` 包含正确的 Go 配置

---

## 🛠️ 常见错误及解决方案

### 错误 1: "command 'go' not found"

**原因：** Go 不在 PATH 环境变量中

**解决：**
```bash
# Windows: 添加 Go 安装目录到 PATH
# 通常是 C:\Program Files\Go\bin 或 D:\go\bin

# 临时测试
$env:PATH += ";D:\go\bin"
```

---

### 错误 2: "dlv command not found"

**原因：** Delve 调试器未安装或不在 PATH 中

**解决：**
```bash
# 重新安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 验证安装
dlv version

# 如果还是找不到，检查 GOPATH/bin 是否在 PATH 中
echo $GOPATH/bin
# 添加到 PATH
```

---

### 错误 3: "file not found in current directory"

**原因：** 缺少 go.mod 文件

**解决：**
```bash
cd d:\yang\go\go_demo\day01
go mod init day01
```

---

### 错误 4: "could not launch process: could not fork/exec"

**原因：** 防火墙阻止了 Delve 调试器

**解决：**
1. 允许 `dlv.exe` 通过 Windows 防火墙
2. 或者临时关闭防火墙测试
3. 以管理员身份运行 VS Code

---

### 错误 5: "building target failed"

**原因：** 代码有语法错误或依赖问题

**解决：**
```bash
# 检查代码
go build -v

# 整理依赖
go mod tidy

# 清理缓存
go clean -cache
```

---

### 错误 6: VS Code 提示 "Go extension is not installed"

**原因：** 未安装 Go 扩展插件

**解决：**
1. 按 `Ctrl+Shift+X` 打开扩展面板
2. 搜索 "Go"
3. 安装 "Go" by Go Team at Google
4. 重启 VS Code

---

### 错误 7: "launch: program does not exist"

**原因：** launch.json 中的路径配置错误

**解决：**
检查 `.vscode/launch.json` 中的 `program` 字段：
```json
{
    "program": "${fileDirname}"  // 使用变量，不要硬编码路径
}
```

---

## 🔧 完整重置步骤

如果以上方法都不行，尝试完全重置：

### Step 1: 清理所有配置

```bash
cd d:\yang\go
rm -rf .vscode
rm go_demo/day01/go.mod
```

### Step 2: 重新安装工具

```bash
# 卸载 Go 扩展
# 在 VS Code 扩展面板右键 Go -> Uninstall

# 重新安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 重新安装 Go 扩展
# 在 VS Code 扩展面板搜索 Go -> Install
```

### Step 3: 重新配置

```bash
# 创建 go.mod
cd go_demo/day01
go mod init day01

# 重新创建 .vscode 配置
# 复制之前的配置文件
```

### Step 4: 重启 VS Code

```bash
# 完全关闭 VS Code
# 重新打开项目
```

---

## 📝 当前配置状态

| 组件 | 状态 | 位置 |
|------|------|------|
| Go | ✅ 已安装 (v1.26.1) | D:\go |
| Delve | ✅ 已安装 (v1.26.1) | C:\Users\...\go\bin |
| go.mod | ✅ 已创建 | d:\yang\go\go_demo\day01 |
| launch.json | ✅ 已配置 | d:\yang\go\.vscode |
| settings.json | ✅ 已配置 | d:\yang\go\.vscode |
| 代码 | ✅ 可运行 | main.go |

---

## 🎯 正确的调试步骤

1. **打开文件**：`d:\yang\go\go_demo\day01\main.go`
2. **设置断点**：在第 18 行左侧点击
3. **启动调试**：按 `F5`
4. **选择配置**：选择 "▶️ 运行当前文件"
5. **观察执行**：程序会暂停在断点处

---

## 📞 如果还是不行

请提供完整的错误信息，包括：

1. **错误提示的完整文本**（截图或复制）
2. **在哪个步骤出现**（运行时/调试时）
3. **VS Code 的输出面板内容**（查看 "DEBUG CONSOLE"）

这样可以帮助更准确地定位问题！

---

**最后更新：** 2026-03-25  
**环境：** Windows + Go 1.26.1 + VS Code
