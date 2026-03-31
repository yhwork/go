# Day 15 - 文件上传与静态资源

## 今日目标

实现文件上传接口（单文件/多文件）、文件验证（大小、类型）、本地存储（按日期分目录）、头像上传与缩略图生成。

---

## 练习任务

### 任务 1：文件上传工具封装

```go
// pkg/upload/upload.go

type UploadConfig struct {
    MaxSize      int64    // 最大文件大小（字节）
    AllowedExts  []string // 允许的扩展名
    BasePath     string   // 存储根目录
    URLPrefix    string   // 访问 URL 前缀
}

type UploadResult struct {
    FileName    string `json:"file_name"`
    FilePath    string `json:"file_path"`
    FileSize    int64  `json:"file_size"`
    ContentType string `json:"content_type"`
    URL         string `json:"url"`
}

var DefaultImageConfig = UploadConfig{
    MaxSize:     5 * 1024 * 1024, // 5MB
    AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
    BasePath:    "uploads/images",
    URLPrefix:   "/static/uploads/images",
}

// SaveFile 保存上传的文件
func SaveFile(file *multipart.FileHeader, config UploadConfig) (*UploadResult, error) {
    // 1. 校验文件大小
    // 2. 校验文件扩展名
    // 3. 生成存储路径（按日期分目录：uploads/images/2024/01/15/）
    // 4. 生成唯一文件名（UUID + 原扩展名）
    // 5. 创建目录
    // 6. 保存文件
    // 7. 返回 UploadResult
}
```

### 任务 2：上传接口实现

```go
// POST /api/v1/upload/image     - 单文件上传
// POST /api/v1/upload/images    - 多文件上传
// POST /api/v1/upload/avatar    - 头像上传（自动生成缩略图）

func (h *UploadHandler) UploadImage(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        response.BadRequest(c, "请选择上传文件")
        return
    }
    result, err := upload.SaveFile(file, upload.DefaultImageConfig)
    if err != nil {
        response.HandleError(c, err)
        return
    }
    response.Success(c, result)
}

func (h *UploadHandler) UploadAvatar(c *gin.Context) {
    // 上传头像 + 用 disintegration/imaging 库生成 200x200 缩略图
}
```

### 任务 3：静态文件服务

```go
// 配置 Gin 静态文件服务
r.Static("/static", "./uploads")
```

---

## 自检清单

- [ ] 文件大小和类型校验生效
- [ ] 文件按日期分目录存储
- [ ] 文件名使用 UUID 防止冲突
- [ ] 多文件上传正常工作
- [ ] 头像缩略图生成正确
- [ ] 静态文件可通过 URL 访问
