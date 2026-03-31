# Day 09 - 请求验证与绑定

## 今日目标

掌握 Gin 的请求数据绑定（ShouldBind 系列）、validator 标签的使用、自定义验证器、以及统一的验证错误响应处理。

---

## 知识点

### 1. Gin 绑定方法

```go
// ShouldBindJSON - 绑定 JSON Body
c.ShouldBindJSON(&req)

// ShouldBindQuery - 绑定 URL 查询参数
c.ShouldBindQuery(&query)

// ShouldBindUri - 绑定路径参数
c.ShouldBindUri(&uri)

// ShouldBind - 根据 Content-Type 自动选择
c.ShouldBind(&req)

// Bind vs ShouldBind：
// Bind 失败会自动返回 400 并设置 Content-Type
// ShouldBind 失败只返回 error，由你控制响应（推荐）
```

### 2. Validator 标签

```go
type CreateUserReq struct {
    Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=32"`
    Age      int    `json:"age"      binding:"omitempty,gte=0,lte=150"`
    Phone    string `json:"phone"    binding:"omitempty,len=11"`
}

// 常用验证标签：
// required     - 必填
// omitempty    - 空值时跳过验证
// min/max      - 字符串长度或数值范围
// len          - 精确长度
// gte/lte      - 大于等于/小于等于（数值）
// email        - 邮箱格式
// url          - URL 格式
// alphanum     - 字母数字
// oneof        - 枚举值 oneof=male female
// eqfield      - 等于另一个字段 eqfield=Password（确认密码）
```

### 3. 自定义验证器

```go
import "github.com/go-playground/validator/v10"

// 自定义手机号验证
func validPhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
    return matched
}

// 注册自定义验证器
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("phone", validPhone)
}

// 使用
type Req struct {
    Phone string `json:"phone" binding:"required,phone"`
}
```

### 4. 验证错误的中文翻译

```go
import (
    "github.com/go-playground/validator/v10"
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    zh_translations "github.com/go-playground/validator/v10/translations/zh"
)
```

---

## 练习任务

### 任务 1：定义完整的请求验证结构体（难度：基础）

```go
// ===== 用户相关 =====

type RegisterReq struct {
    Username        string `json:"username"         binding:"required,min=3,max=20,alphanum"`
    Email           string `json:"email"            binding:"required,email"`
    Password        string `json:"password"         binding:"required,min=8,max=32"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
    Phone           string `json:"phone"            binding:"omitempty,phone"`
}

type LoginReq struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type UpdateProfileReq struct {
    Username string `json:"username" binding:"omitempty,min=3,max=20"`
    Email    string `json:"email"    binding:"omitempty,email"`
    Phone    string `json:"phone"    binding:"omitempty,phone"`
    Avatar   string `json:"avatar"   binding:"omitempty,url"`
}

// ===== 文章相关 =====

type CreateArticleReq struct {
    Title      string   `json:"title"       binding:"required,min=1,max=200"`
    Content    string   `json:"content"     binding:"required,min=10"`
    CategoryID uint     `json:"category_id" binding:"required,gt=0"`
    Tags       []string `json:"tags"        binding:"omitempty,max=5,dive,min=1,max=20"`
    Status     string   `json:"status"      binding:"required,oneof=draft published"`
}

type ListArticleQuery struct {
    Page       int    `form:"page"        binding:"omitempty,min=1"`
    PageSize   int    `form:"page_size"   binding:"omitempty,min=1,max=100"`
    CategoryID uint   `form:"category_id" binding:"omitempty"`
    Tag        string `form:"tag"         binding:"omitempty"`
    Keyword    string `form:"keyword"     binding:"omitempty,max=50"`
    SortBy     string `form:"sort_by"     binding:"omitempty,oneof=created_at updated_at views"`
    SortOrder  string `form:"sort_order"  binding:"omitempty,oneof=asc desc"`
}

// ===== 路径参数 =====

type IDUri struct {
    ID uint `uri:"id" binding:"required,gt=0"`
}
```

### 任务 2：自定义验证器注册（难度：中等）

```go
// validator/custom.go

// 1. 手机号验证器
func validPhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
    return matched
}

// 2. 密码强度验证器（至少包含大写、小写、数字）
func validStrongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    return hasUpper && hasLower && hasDigit
}

// 3. 非空字符串验证器（不能全是空格）
func validNotBlank(fl validator.FieldLevel) bool {
    return strings.TrimSpace(fl.Field().String()) != ""
}

// 注册所有自定义验证器
func RegisterCustomValidators(engine *validator.Validate) {
    engine.RegisterValidation("phone", validPhone)
    engine.RegisterValidation("strong_password", validStrongPassword)
    engine.RegisterValidation("notblank", validNotBlank)
}
```

### 任务 3：统一验证错误响应（难度：中等偏上）

```go
// 将 validator 的错误信息转换为用户友好的中文提示

type ValidationErrorItem struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// 错误消息映射表
var validationMessages = map[string]string{
    "required":        "不能为空",
    "min":             "长度不能小于 %s",
    "max":             "长度不能大于 %s",
    "email":           "邮箱格式不正确",
    "phone":           "手机号格式不正确",
    "alphanum":        "只能包含字母和数字",
    "eqfield":         "两次输入不一致",
    "oneof":           "值必须是 %s 之一",
    "gt":              "必须大于 %s",
    "gte":             "必须大于或等于 %s",
    "url":             "URL 格式不正确",
    "strong_password": "密码必须包含大写字母、小写字母和数字",
}

// 解析验证错误
func ParseValidationErrors(err error) []ValidationErrorItem {
    var errs validator.ValidationErrors
    if !errors.As(err, &errs) {
        return []ValidationErrorItem{{Field: "unknown", Message: err.Error()}}
    }

    items := make([]ValidationErrorItem, 0, len(errs))
    for _, e := range errs {
        items = append(items, ValidationErrorItem{
            Field:   toSnakeCase(e.Field()), // 字段名转 snake_case
            Message: getValidationMessage(e),
        })
    }
    return items
}

func getValidationMessage(e validator.FieldError) string {
    // 根据 e.Tag() 从 validationMessages 中获取模板
    // 根据 e.Param() 填充参数
    // 返回格式化后的中文消息
}

// 在 handler 中使用：
func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterReq
    if err := c.ShouldBindJSON(&req); err != nil {
        items := ParseValidationErrors(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "请求参数验证失败",
            "errors":  items,
        })
        return
    }
    // ... 业务逻辑
}

// 期望的错误响应格式：
// {
//     "code": 400,
//     "message": "请求参数验证失败",
//     "errors": [
//         {"field": "username", "message": "长度不能小于 3"},
//         {"field": "email", "message": "邮箱格式不正确"},
//         {"field": "confirm_password", "message": "两次输入不一致"}
//     ]
// }
```

### 任务 4：封装验证中间件（难度：中等）

```go
// 封装一个通用的请求验证中间件，减少 handler 中的重复代码

// BindJSON 绑定 JSON Body 并验证
func BindJSON[T any](c *gin.Context) (*T, bool) {
    var req T
    if err := c.ShouldBindJSON(&req); err != nil {
        items := ParseValidationErrors(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "请求参数验证失败",
            "errors":  items,
        })
        return nil, false
    }
    return &req, true
}

// BindQuery 绑定查询参数并验证
func BindQuery[T any](c *gin.Context) (*T, bool) {
    var req T
    if err := c.ShouldBindQuery(&req); err != nil {
        items := ParseValidationErrors(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "请求参数验证失败",
            "errors":  items,
        })
        return nil, false
    }
    return &req, true
}

// 使用方式：
func (h *UserHandler) Register(c *gin.Context) {
    req, ok := BindJSON[RegisterReq](c)
    if !ok {
        return // 错误响应已在 BindJSON 中处理
    }
    // 使用 req 处理业务逻辑
}
```

---

## 测试验证

```bash
# 测试正常注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"Pass1234","confirm_password":"Pass1234"}'

# 测试验证失败
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"ab","email":"invalid","password":"123","confirm_password":"456"}'

# 期望返回多条验证错误
```

---

## 自检清单

- [ ] 掌握 ShouldBindJSON/Query/Uri 的使用
- [ ] 验证标签（required/min/max/email/oneof 等）使用正确
- [ ] 自定义验证器（手机号、密码强度）注册并生效
- [ ] 验证错误能转换为中文友好提示
- [ ] 错误响应格式统一、字段名清晰
- [ ] 泛型绑定函数减少了 handler 中的重复代码
