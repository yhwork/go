# Go Demo 练习项目

这是 Go 全栈学习计划的配套练习代码。

## 📁 目录结构

```
go_demo/
├── README.md              # 本文件
├── day01/                 # Day 01 练习
│   ├── main.go           # 主程序
│   └── go.mod            # Go 模块定义
├── day02/                # Day 02 练习（待创建）
└── ...
```

## 🚀 快速开始

### 运行单个练习

```bash
cd day01
go run main.go
```

### 调试

在 VS Code 中打开 `day01/main.go`，按 F5 选择 "▶️ 运行当前文件"。

## 📋 已完成修复

### ✅ Day 01

**问题：** 缺少 go.mod 导致无法调试  
**解决：** 创建 `go.mod` 文件  
**状态：** 可以正常运行和调试

**测试输出：**
```
hello world
1607
zhangsan
18
```

## 🎯 学习计划关联

每个 `dayXX` 目录对应学习计划中的一天的练习内容：

| 目录 | 对应文档 | 主题 |
|------|---------|------|
| day01 | `../go-learning-plan/phase1-go-advanced/day01-struct-methods.md` | 结构体与方法 |
| day02 | `../go-learning-plan/phase1-go-advanced/day02-interface-polymorphism.md` | 接口与多态 |
| day03 | `../go-learning-plan/phase1-go-advanced/day03-error-concurrency.md` | 错误处理与并发 |
| ... | ... | ... |

## 🛠️ 常见问题

### Q: 运行时提示 "file not found"
A: 确保目录中有 `go.mod` 文件。如果没有，运行：
```bash
go mod init <day-name>
```

### Q: 如何添加新的一天练习？
A: 
1. 创建新目录 `mkdir day02`
2. 创建 `main.go` 文件
3. 初始化模块 `go mod init day02`
4. 编写代码并测试

### Q: 如何在多个练习间切换？
A: 每个练习是独立的模块，直接 `cd` 到对应目录即可。

## 📖 相关资源

- [Go 学习主计划](../go-learning-plan/README.md)
- [调试配置指南](../.vscode/GO_DEBUG_CONFIG.md)
- [快速入门指南](../.vscode/DEBUG_QUICK_START.md)
- [故障排除](../.vscode/TROUBLESHOOTING.md)

---

**提示：** 按照 `../go-learning-plan/README.md` 中的计划每天完成一个练习！
