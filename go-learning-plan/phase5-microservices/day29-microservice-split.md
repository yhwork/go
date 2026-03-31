# Day 29 - 微服务拆分

## 今日目标

将博客系统拆分为独立的微服务（User Service + Article Service），通过 gRPC 通信，统一 API Gateway 入口。

---

## 练习任务

### 任务 1：服务拆分规划

```
go-blog-microservices/
├── api-gateway/          # HTTP API 网关（Gin）
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── handler/      # 转发请求到各微服务
│   │   ├── middleware/    # JWT 认证、限流等
│   │   └── router/
│   └── go.mod
│
├── user-service/          # 用户微服务（gRPC）
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── proto/
│   │   └── user.proto
│   └── go.mod
│
├── article-service/       # 文章微服务（gRPC）
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── proto/
│   │   └── article.proto
│   └── go.mod
│
├── proto/                 # 共享 Proto 定义
│   ├── user/
│   └── article/
│
└── docker-compose.yml
```

### 任务 2：定义 Proto 文件

```protobuf
// proto/user/user.proto
syntax = "proto3";
package user;
option go_package = "proto/user";

service UserService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}

message RegisterRequest {
    string username = 1;
    string email = 2;
    string password = 3;
}

message LoginRequest {
    string username = 1;
    string password = 2;
}

message LoginResponse {
    string access_token = 1;
    string refresh_token = 2;
    UserInfo user = 3;
}

message ValidateTokenRequest {
    string token = 1;
}

message ValidateTokenResponse {
    bool valid = 1;
    uint64 user_id = 2;
    string username = 3;
    string role = 4;
}

message UserInfo {
    uint64 id = 1;
    string username = 2;
    string email = 3;
    string avatar = 4;
    string role = 5;
}

// proto/article/article.proto - 自行定义
```

### 任务 3：实现 User Service

```go
// user-service/cmd/main.go
func main() {
    // 初始化数据库
    // 启动 gRPC 服务
    lis, _ := net.Listen("tcp", ":50051")
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &userGRPCServer{...})
    reflection.Register(s)
    s.Serve(lis)
}

// user-service/internal/grpc_server.go
type userGRPCServer struct {
    pb.UnimplementedUserServiceServer
    userService service.UserService
}

func (s *userGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    result, err := s.userService.Login(&dto.LoginReq{
        Username: req.Username,
        Password: req.Password,
    })
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, err.Error())
    }
    return &pb.LoginResponse{
        AccessToken:  result.AccessToken,
        RefreshToken: result.RefreshToken,
        User: &pb.UserInfo{
            Id:       uint64(result.User.ID),
            Username: result.User.Username,
        },
    }, nil
}

// ValidateToken - 供 API Gateway 调用，验证 Token 有效性
func (s *userGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    claims, err := s.jwtManager.ParseToken(req.Token)
    if err != nil {
        return &pb.ValidateTokenResponse{Valid: false}, nil
    }
    return &pb.ValidateTokenResponse{
        Valid:    true,
        UserId:   uint64(claims.UserID),
        Username: claims.Username,
        Role:     claims.Role,
    }, nil
}
```

### 任务 4：实现 API Gateway

```go
// api-gateway/cmd/main.go
func main() {
    // 连接各微服务
    userConn, _ := grpc.Dial("user-service:50051", grpc.WithInsecure())
    articleConn, _ := grpc.Dial("article-service:50052", grpc.WithInsecure())

    userClient := userpb.NewUserServiceClient(userConn)
    articleClient := articlepb.NewArticleServiceClient(articleConn)

    // 启动 Gin HTTP 服务
    r := gin.Default()

    // 认证中间件改为调用 User Service 验证 Token
    authMiddleware := func(c *gin.Context) {
        token := extractToken(c)
        resp, err := userClient.ValidateToken(c.Request.Context(), &userpb.ValidateTokenRequest{Token: token})
        if err != nil || !resp.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("userID", resp.UserId)
        c.Set("username", resp.Username)
        c.Set("role", resp.Role)
        c.Next()
    }

    // 路由转发
    v1 := r.Group("/api/v1")
    {
        // 认证相关 -> User Service
        v1.POST("/auth/login", gatewayHandler.Login)
        v1.POST("/auth/register", gatewayHandler.Register)

        // 文章相关 -> Article Service
        v1.GET("/articles", gatewayHandler.ListArticles)
        v1.GET("/articles/:id", gatewayHandler.GetArticle)

        protected := v1.Group("")
        protected.Use(authMiddleware)
        {
            protected.POST("/articles", gatewayHandler.CreateArticle)
        }
    }

    r.Run(":8080")
}
```

### 任务 5：Docker Compose 编排微服务

```yaml
services:
  api-gateway:
    build: ./api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - user-service
      - article-service

  user-service:
    build: ./user-service
    expose:
      - "50051"
    depends_on:
      - mysql
      - redis

  article-service:
    build: ./article-service
    expose:
      - "50052"
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    # ...

  redis:
    image: redis:7-alpine
    # ...
```

---

## 自检清单

- [ ] User Service 和 Article Service 独立运行
- [ ] gRPC 通信正常
- [ ] API Gateway 正确转发请求
- [ ] Token 验证通过跨服务调用 User Service 完成
- [ ] Docker Compose 编排所有服务正常启动
- [ ] 通过 API Gateway 的 HTTP 接口能完成所有操作
