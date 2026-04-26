# BanKV - 分布式键值数据库

基于 Raft 共识算法和 LSM-Tree 存储引擎构建的高可用分布式 KV 数据库。

## 🌟 项目亮点

- **高可用性**: 基于 Raft 共识算法实现数据一致性，支持多节点集群
- **高性能存储**: 采用 LSM-Tree 架构，包含 MemTable (跳表) 和 SSTable，优化写入性能
- **预写日志(WAL)**: 确保数据持久性和崩溃恢复能力
- **自定义网络框架**: 内置 TCP 网络通信框架，支持消息路由和连接管理
- **交互式客户端**: 提供命令行交互模式和批量操作模式
- **纯 Go 实现**: 无第三方依赖，使用 Go 标准库构建

## 🏗️ 系统架构

```
BanKV/
├── cmd/
│   ├── server/          # 服务端入口，负责启动 Raft、存储引擎和网络服务
│   │   ├── server.go    # 服务端主程序
│   │   ├── run.bat      # Windows 启动脚本
│   │   └── run-clean.bat# Windows 清理数据后启动脚本
│   └── client/          # 客户端入口，支持交互模式和命令行模式
│       ├── main.go      # 客户端主程序
│       ├── client.go    # 客户端核心逻辑
│       ├── interactive.go# 交互式客户端实现
│       ├── README.md    # 客户端使用说明
│       └── run.bat      # Windows 启动脚本
├── internal/
│   ├── Raft/            # Raft 一致性协议实现
│   │   ├── raft.go      # 选主、日志复制、状态机应用核心逻辑
│   │   ├── raft_test.go # Raft 算法测试
│   │   ├── raft_wal.go  # Raft 日志持久化
│   │   ├── rpc.go       # RPC 通信结构定义
│   │   └── rpc_test.go  # RPC 通信测试
│   ├── storage/         # 存储引擎核心 (LSM-Tree)
│   │   ├── engine.go    # 存储引擎封装，对外提供 Put/Get/Delete 接口
│   │   ├── engine_test.go# 存储引擎测试
│   │   ├── istorage/    # 存储接口定义
│   │   │   └── interfaces.go# 存储相关接口
│   │   └── zstorage/    # 存储引擎具体实现
│   │       ├── memtable.go  # 内存表 (跳表实现)
│   │       ├── memtable_test.go# MemTable 测试
│   │       ├── WAL.go       # 预写日志实现
│   │       └── wal_test.go  # WAL 测试
│   ├── network/         # 网络通信框架
│   │   ├── ziface/      # 网络接口定义
│   │   │   ├── IConnManager.go# 连接管理器接口
│   │   │   ├── IDataPack.go   # 数据打包接口
│   │   │   ├── IMsgHandle.go  # 消息处理接口
│   │   │   ├── iRequest.go    # 请求接口
│   │   │   ├── iRouter.go     # 路由接口
│   │   │   ├── iconnect.go    # 连接接口
│   │   │   ├── imessage.go    # 消息接口
│   │   │   └── isever.go      # 服务器接口
│   │   └── znet/        # 网络框架实现
│   │       ├── server.go    # TCP 服务器实现
│   │       ├── server_test.go# 服务器测试
│   │       ├── connection.go# 连接管理实现
│   │       ├── ConnManager.go# 连接管理器
│   │       ├── msgHandle.go # 消息处理器
│   │       ├── DataPack.go  # 数据打包器
│   │       ├── request.go   # 请求实现
│   │       ├── router.go    # 路由实现
│   │       └── message.go   # 消息实现
│   └── service/         # 业务逻辑层
│       ├── fsm.go       # 状态机应用逻辑 (Raft -> Storage)
│       ├── fsm_test.go  # FSM 测试
│       ├── router.go    # 请求路由处理 (PUT/GET/DELETE)
│       ├── ha.go        # 高可用监控
│       └── test_service_wal.log# 测试用 WAL 日志
├── config/              # 配置管理
│   ├── config.json      # JSON 配置文件
│   └── global.go        # 全局配置加载和管理
├── pkg/utils/           # 工具函数
│   ├── byteBuilder.go   # 字节构建工具
│   ├── datapack.go      # 数据打包工具
│   └── message.go       # 消息处理工具
├── data/                # 数据存储目录 (SSTable)
├── log/                 # 日志文件目录
│   ├── SSTable          # SSTable 日志
│   └── wal.log          # WAL 日志
├── .idea/               # IDE 配置目录
├── .trae/               # Trae 配置目录
├── .ignore              # 忽略文件配置
├── go.mod               # Go 模块定义
├── README.md            # 项目说明文档
└── QUICKSTART.md        # 快速启动指南
```

## 🚀 快速开始

### 前置要求

- Go 1.26.1 或更高版本

### 启动服务端

```bash
# 方法1: 使用启动脚本 (推荐)
cd cmd/server
./run.bat

# 方法2: 手动启动
cd BanKV
go run cmd/server/server.go
```

### 启动客户端

```bash
# 方法1: 交互模式 (推荐)
cd cmd/client
./run.bat

# 方法2: 命令行模式
go run cmd/client put name Alice
go run cmd/client get name
go run cmd/client delete name
```

### 基本操作示例

在交互模式中:

```
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

> quit
再见！
```

## 📋 功能特性

### 核心功能
- ✅ PUT: 插入或更新键值对
- ✅ GET: 根据键查询值
- ✅ DELETE: 删除键值对
- ✅ Raft 共识: 保证多节点数据一致性
- ✅ WAL: 预写日志保证数据持久性
- ✅ LSM-Tree: 高效的存储引擎架构

### 高级特性
- 🔧 跳表实现的 MemTable，支持 O(log n) 查找
- 🔧 SSTable 持久化存储
- 🔧 自动 Flush 机制 (MemTable -> SSTable)
- 🔧 单节点和多节点集群模式
- 🔧 心跳检测和领导者选举

## 🛠️ 技术栈

- **语言**: Go 1.26.1
- **共识算法**: Raft
- **存储引擎**: LSM-Tree (MemTable + SSTable + WAL)
- **数据结构**: 跳表 (SkipList)
- **网络**: TCP 自定义协议
- **依赖**: 无第三方依赖，纯标准库实现

## 📊 性能特点

- **写入优化**: LSM-Tree 架构将随机写转换为顺序写
- **读取加速**: MemTable 内存缓存 + 多层索引
- **崩溃恢复**: WAL 日志确保数据不丢失
- **一致性保证**: Raft 算法保证强一致性

## 📖 详细文档

- [快速启动指南](QUICKSTART.md) - 详细的启动和问题排查指南

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。