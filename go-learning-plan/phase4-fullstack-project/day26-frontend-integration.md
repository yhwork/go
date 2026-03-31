# Day 26 - 前后端联调

## 今日目标

确保所有 API 文档完整，配置 CORS，使用你熟悉的前端框架搭建博客前端，实现核心页面与后端联调。

---

## 练习任务

### 任务 1：API 文档完善

确保所有接口都有 Swagger 注解，运行 `swag init` 生成最新文档。

```bash
swag init -g cmd/server/main.go -o docs
```

访问 `http://localhost:8080/swagger/index.html` 验证所有接口文档完整。

### 任务 2：CORS 配置确认

```go
// 开发环境允许前端端口（如 3000、5173）跨域
corsConfig := middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
    ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
    AllowCredentials: true,
    MaxAge:           86400,
}
```

### 任务 3：前端项目搭建

使用你熟悉的前端框架（React / Vue / Next.js 等）：

```bash
# React 示例
npx create-react-app go-blog-frontend --template typescript
# 或 Vue
npm create vue@latest go-blog-frontend
# 或 Next.js
npx create-next-app@latest go-blog-frontend
```

### 任务 4：核心页面实现

需要实现的页面（作为前端开发者，这部分应该是你的强项）：

| 页面 | 路由 | 说明 |
|------|------|------|
| 登录页 | /login | 登录表单、Token 存储 |
| 注册页 | /register | 注册表单、表单验证 |
| 文章列表 | / | 分页、搜索、分类筛选 |
| 文章详情 | /article/:id | 内容渲染、评论、点赞 |
| 写文章 | /editor | Markdown 编辑器 |
| 个人中心 | /profile | 个人信息、我的文章 |

### 任务 5：API 请求封装

```typescript
// 封装 axios 实例
const api = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 10000,
});

// 请求拦截器：自动携带 Token
api.interceptors.request.use(config => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// 响应拦截器：处理 401 和 Token 刷新
api.interceptors.response.use(
    response => response.data,
    async error => {
        if (error.response?.status === 401) {
            // 尝试刷新 Token
            // 如果刷新失败，跳转登录页
        }
        return Promise.reject(error);
    }
);
```

### 联调检查清单

- [ ] 注册 -> 登录 -> 自动存储 Token
- [ ] Token 过期 -> 自动刷新
- [ ] 文章列表正常加载和分页
- [ ] 搜索和分类筛选正常
- [ ] 文章详情页渲染正常
- [ ] Markdown 内容正确渲染
- [ ] 评论功能正常
- [ ] 点赞/收藏功能正常
- [ ] 文件上传（头像、文章封面）正常

---

## 自检清单

- [ ] Swagger 文档完整，所有接口都可以通过文档测试
- [ ] CORS 配置正确，前端跨域请求无报错
- [ ] 前端项目结构清晰
- [ ] 至少完成登录、文章列表、文章详情三个核心页面的联调
- [ ] Token 自动携带和刷新机制工作正常
