# Day 17 - WebSocket 实时通信

## 今日目标

使用 gorilla/websocket 实现一个简单的聊天室：广播消息、私聊、在线用户列表维护、连接管理。

---

## 练习任务

### 任务 1：WebSocket Hub（连接管理中心）

```go
// internal/ws/hub.go

type Client struct {
    ID       uint
    Username string
    Conn     *websocket.Conn
    Send     chan []byte
}

type Message struct {
    Type     string `json:"type"`      // "broadcast", "private", "system"
    From     string `json:"from"`
    To       uint   `json:"to,omitempty"` // 私聊目标用户 ID
    Content  string `json:"content"`
    Time     string `json:"time"`
}

type Hub struct {
    clients    map[uint]*Client // userID -> Client
    register   chan *Client
    unregister chan *Client
    broadcast  chan *Message
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[uint]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan *Message),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.ID] = client
            h.mu.Unlock()
            // 广播系统消息：xxx 加入了聊天室

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.ID]; ok {
                close(client.Send)
                delete(h.clients, client.ID)
            }
            h.mu.Unlock()
            // 广播系统消息：xxx 离开了聊天室

        case msg := <-h.broadcast:
            if msg.To > 0 {
                // 私聊：只发给目标用户
            } else {
                // 广播：发给所有人
            }
        }
    }
}
```

### 任务 2：客户端读写 Goroutine

```go
// internal/ws/client.go

func (c *Client) ReadPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.Conn.Close()
    }()

    c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    for {
        var msg Message
        err := c.Conn.ReadJSON(&msg)
        if err != nil {
            break
        }
        msg.From = c.Username
        msg.Time = time.Now().Format("15:04:05")
        hub.broadcast <- &msg
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(30 * time.Second)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case msg, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.Conn.WriteMessage(websocket.TextMessage, msg)
        case <-ticker.C:
            c.Conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

### 任务 3：WebSocket 升级接口

```go
// internal/handler/ws_handler.go

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 开发阶段允许所有来源
    },
}

func (h *WSHandler) HandleWS(c *gin.Context) {
    // 1. 从 JWT 获取用户信息
    userID, _ := middleware.GetCurrentUserID(c)
    username, _ := c.Get("username")

    // 2. 升级为 WebSocket 连接
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }

    // 3. 创建 Client 并注册到 Hub
    client := &Client{
        ID:       userID,
        Username: username.(string),
        Conn:     conn,
        Send:     make(chan []byte, 256),
    }
    h.hub.register <- client

    go client.WritePump()
    go client.ReadPump(h.hub)
}

// 获取在线用户列表
func (h *WSHandler) OnlineUsers(c *gin.Context) {
    h.hub.mu.RLock()
    defer h.hub.mu.RUnlock()

    users := make([]gin.H, 0, len(h.hub.clients))
    for _, client := range h.hub.clients {
        users = append(users, gin.H{
            "id":       client.ID,
            "username": client.Username,
        })
    }
    response.Success(c, users)
}
```

### 任务 4：简单的 HTML 测试页面

创建一个简单的 HTML 页面用来测试 WebSocket 聊天功能。用你的前端经验快速搭建即可。

---

## 自检清单

- [ ] WebSocket 连接建立和断开正常
- [ ] 广播消息所有在线用户都能收到
- [ ] 私聊消息只有目标用户能收到
- [ ] 在线用户列表实时更新
- [ ] Ping/Pong 心跳保持连接活跃
- [ ] 连接断开后资源正确清理（无 goroutine 泄漏）
