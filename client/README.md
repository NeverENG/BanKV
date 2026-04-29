# BanKV 客户端使用指南

## 📖 概述

BanKV 客户端支持两种模式：
1. **交互模式**：类似 Redis CLI，保持长连接，持续输入命令
2. **命令行模式**：单次执行命令后退出

## 🚀 快速开始

### 1. 启动服务端

```powershell
cd E:\Code\BanKv
go run cmd/server/server.go
```

### 2. 使用客户端

#### 方式一：交互模式（推荐）

**方法1：使用启动脚本（最简单）**
```powershell
cd E:\Code\BanKv\cmd\client
.\run.bat

# 或指定服务器地址
.\run.bat localhost:8080
```

**方法2：使用 go run**
```powershell
# 从项目根目录
cd E:\Code\BanKv
go run cmd/client

# 或从 client 目录
cd E:\Code\BanKv\cmd\client
go run .

# 或指定服务器地址
go run . localhost:8080
```

**⚠️ 注意：** 不要直接运行 `go run main.go`，需要运行整个包！

**交互示例：**
```
已连接到 localhost:8080
输入命令进行操作，输入 'quit' 或 'exit' 退出
支持命令: put <key> <value>, get <key>, delete <key>

> put name Alice
✅ OK

> put age 25
✅ OK

> get name
"Alice"

> get age
"25"

> delete age
✅ OK

> get age
❌ 错误: key not found or server error

> help

=== BanKV 客户端帮助 ===
命令:
  put <key> <value>  - 存储键值对
  get <key>          - 获取值
  delete <key>       - 删除键
  help               - 显示此帮助信息
  quit/exit          - 退出客户端

> quit
再见！
```

#### 方式二：命令行模式

```powershell
# 从项目根目录
cd E:\Code\BanKv

# PUT 操作
go run cmd/client put name Alice

# GET 操作
go run cmd/client get name

# DELETE 操作
go run cmd/client delete name
```

## 📋 支持的命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `put <key> <value>` | 存储键值对 | `put name Alice` |
| `get <key>` | 获取值 | `get name` |
| `delete <key>` | 删除键 | `delete name` |
| `help` | 显示帮助 | `help` |
| `quit` / `exit` | 退出客户端 | `quit` |

## ✨ 特性

- ✅ **长连接**：交互模式下复用同一个 TCP 连接
- ✅ **友好提示**：清晰的 ✅/❌ 状态标识
- ✅ **智能解析**：支持 value 中有空格
- ✅ **帮助系统**：输入 `help` 查看使用说明
- ✅ **超时保护**：5秒读写超时，防止卡死

## 🔧 技术实现

### 数据包格式

**客户端 → 服务端：**
```
[dataLen(4字节)] [msgID(4字节)] [data]
```

其中 data 的内容：
- **PUT (msgID=1)**: `keylen(4) + key + valuelen(4) + value`
- **GET (msgID=2)**: `keylen(4) + key`
- **DELETE (msgID=3)**: `keylen(4) + key`

**服务端 → 客户端：**
```
[dataLen(4字节)] [msgID(4字节)] [status(1字节)] [可选的value]
```

### 架构设计

```
cmd/client/
├── main.go          # 入口文件，支持两种模式
├── client.go        # 客户端核心实现（TCP连接、消息发送）
└── interactive.go   # 交互式客户端（命令循环、用户输入）

pkg/utils/
├── message.go       # Message 结构体和 KV 数据构建
├── byteBuilder.go   # 字节切片拼接工具
└── datapack.go      # 数据打包/解包

internal/network/znet/
└── DataPack.go      # 网络层数据包处理
```

## 💡 使用技巧

1. **批量操作**：在交互模式下可以连续执行多个命令，无需重新连接
2. **Value 含空格**：`put greeting hello world` 会自动合并为 "hello world"
3. **快速退出**：输入 `q` 不行，必须输入完整的 `quit` 或 `exit`
4. **查看帮助**：随时输入 `help` 查看命令列表

## ⚠️ 注意事项

1. 确保服务端已启动后再运行客户端
2. 默认连接地址是 `localhost:8080`
3. 如果连接失败，检查服务端是否正常运行
4. 交互模式下按 Ctrl+C 也可以强制退出

## 🎯 下一步

可以考虑增强的功能：
- 命令历史记录（上下键翻阅）
- 自动补全
- 批量导入/导出
- 性能测试模式
