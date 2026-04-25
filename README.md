基于 Raft 共识算法开发的 KV数据库
```
gokv/
├── cmd/
│   └── server/          # 程序的入口，负责解析 flag、初始化并启动服务
│       └── main.go
├── internal/# 私有逻辑，不希望被作为第三方库引用
│   ├── storage/         # 存储引擎核心 (LSM-Tree)
│   │   ├── memtable.go  # 内存表 (跳表实现)
│   │   ├── sstable.go   # 磁盘有序表 (SSTable 读写逻辑)
│   │   ├── wal.go       # 预写日志
│   │   └── engine.go    # 封装上面三者，对外提供 Put/Get 接口
│   ├── raft/            # Raft 一致性协议实现
│   │   ├── raft.go      # 选主、日志分发核心
│   │   └── rpc.go       # RPC 请求/响应结构定义
│   └── service/         # 胶水层：把 Raft 和 Storage 缝合在一起
│       └── fsm.go       # 状态机应用逻辑 (从 Raft 取日志写入 Storage)
├── api/                 # 协议定义
│   └── kv.proto         # 客户端与服务端通信的 Protobuf 定义
├── pkg/                 # 可以暴露给外部的公共工具
│   └── utils/           # 比如自定义的编码工具、位运算等
├── go.mod
└── Makefile             # 定义编译、测试命令
```