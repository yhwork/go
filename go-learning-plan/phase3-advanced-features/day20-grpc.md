# Day 20 - gRPC 入门

## 今日目标

了解 Protocol Buffers 和 gRPC 的基本概念，定义 .proto 文件，实现 gRPC 服务端和客户端，在 Gin 中创建 HTTP 网关。

---

## 知识点

### Protocol Buffers

```protobuf
// proto/user.proto
syntax = "proto3";
package user;
option go_package = "github.com/yourname/go-blog/proto/user";

message User {
    uint64 id = 1;
    string username = 2;
    string email = 3;
    string role = 4;
    string created_at = 5;
}

message GetUserRequest {
    uint64 id = 1;
}

message GetUserResponse {
    User user = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
    string keyword = 3;
}

message ListUsersResponse {
    repeated User users = 1;
    int64 total = 2;
}

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}
```

---

## 练习任务

### 任务 1：环境搭建与 Proto 定义

```bash
# 安装 protoc 编译器和 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 编译 proto 文件
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

定义完整的 User 和 Article 两个服务的 proto 文件。

### 任务 2：实现 gRPC 服务端

```go
// internal/grpc/user_server.go

type userServer struct {
    pb.UnimplementedUserServiceServer
    userService service.UserService
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.userService.GetProfile(uint(req.Id))
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "用户不存在")
    }
    return &pb.GetUserResponse{
        User: toProtoUser(user),
    }, nil
}

// 启动 gRPC 服务器
func StartGRPCServer(port int, userSvc service.UserService) error {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return err
    }
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &userServer{userService: userSvc})
    reflection.Register(s) // 用于 grpcurl 调试
    return s.Serve(lis)
}
```

### 任务 3：实现 gRPC 客户端

```go
// 连接 gRPC 服务并调用
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
defer conn.Close()

client := pb.NewUserServiceClient(conn)
resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: 1})
```

### 任务 4：Gin HTTP 网关

```go
// 在 Gin 中转发 HTTP 请求到 gRPC 服务

func (h *GatewayHandler) GetUser(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    resp, err := h.userClient.GetUser(c.Request.Context(), &pb.GetUserRequest{Id: id})
    if err != nil {
        st, ok := status.FromError(err)
        if ok && st.Code() == codes.NotFound {
            response.NotFound(c, st.Message())
            return
        }
        response.ServerError(c, "服务调用失败")
        return
    }

    response.Success(c, resp.User)
}
```

### 任务 5：测试

```bash
# 使用 grpcurl 测试
grpcurl -plaintext -d '{"id": 1}' localhost:50051 user.UserService/GetUser

# 通过 HTTP 网关测试
curl http://localhost:8080/api/v1/users/1
```

---

## 自检清单

- [ ] Proto 文件定义完整，编译成功
- [ ] gRPC 服务端正常启动并响应请求
- [ ] gRPC 客户端能正确调用服务
- [ ] HTTP 网关正确转发请求和错误
- [ ] gRPC 错误码能正确映射到 HTTP 状态码
- [ ] 理解 gRPC 相比 REST 的优缺点
