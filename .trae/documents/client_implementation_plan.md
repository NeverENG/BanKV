# 客户端连接服务端实现计划

## 项目现状分析

### 当前结构
- **服务端**：已实现基于 Network 层的 TCP 服务，支持 PUT(1)、GET(2)、DELETE(3) 操作
- **协议**：基于消息 ID 的 TLV 格式，消息结构包含 Id、DataLen、Data
- **客户端**：当前为空，需要实现

### 客户端需求
- 连接到服务端
- 发送 PUT、GET、DELETE 请求
- 接收和处理响应
- 错误处理
- 命令行界面

## 实现方案

### 1. 客户端核心功能

**文件**：`cmd/client/client.go`

- **功能**：
  - 连接管理
  - 消息编解码
  - 请求发送
  - 响应处理
  - 命令行界面

**关键代码**：
```go
// Client 客户端结构
type Client struct {
	conn net.Conn
	addr string
}

// NewClient 创建客户端实例
func NewClient(addr string) *Client

// Connect 连接到服务端
func (c *Client) Connect() error

// SendPut 发送 PUT 请求
func (c *Client) SendPut(key, value []byte) error

// SendGet 发送 GET 请求
func (c *Client) SendGet(key []byte) ([]byte, error)

// SendDelete 发送 DELETE 请求
func (c *Client) SendDelete(key []byte) error

// Close 关闭连接
func (c *Client) Close() error
```

### 2. 消息编解码

**文件**：`cmd/client/packer.go`

- **功能**：
  - 编码消息为二进制格式
  - 解码二进制数据为消息

**关键代码**：
```go
// Pack 打包消息
func Pack(msgID uint32, data []byte) ([]byte, error)

// Unpack 解包消息
func Unpack(data []byte) (uint32, []byte, error)
```

### 3. 命令行界面

**文件**：`cmd/client/main.go`

- **功能**：
  - 解析命令行参数
  - 执行相应的操作
  - 显示结果

**关键代码**：
```go
func main() {
    // 解析命令行参数
    // 执行 PUT/GET/DELETE 操作
    // 显示结果
}
```

## 实现步骤

### 阶段 1：核心客户端实现
1. 创建 `cmd/client/client.go`
2. 实现连接管理和基本操作

### 阶段 2：消息编解码
1. 创建 `cmd/client/packer.go`
2. 实现消息打包和解包

### 阶段 3：命令行界面
1. 创建 `cmd/client/main.go`
2. 实现命令行参数解析和执行

### 阶段 4：测试
1. 测试 PUT 操作
2. 测试 GET 操作
3. 测试 DELETE 操作
4. 测试错误处理

## 技术要点

### 1. 协议格式
- **消息结构**：消息 ID (4字节) + 数据长度 (4字节) + 数据内容
- **操作类型**：PUT(1)、GET(2)、DELETE(3)
- **数据格式**：
  - PUT：key_len(4字节) + key + value_len(4字节) + value
  - GET：key_len(4字节) + key
  - DELETE：key_len(4字节) + key

### 2. 错误处理
- 连接错误
- 发送错误
- 接收错误
- 响应错误

### 3. 性能考虑
- 连接复用
- 超时设置
- 缓冲区大小

## 预期成果

- ✅ 完整的客户端实现
- ✅ 支持 PUT、GET、DELETE 操作
- ✅ 命令行界面
- ✅ 错误处理机制
- ✅ 与服务端的正确通信

## 示例使用

```bash
# PUT 操作
./client put key1 value1

# GET 操作
./client get key1

# DELETE 操作
./client delete key1
```
