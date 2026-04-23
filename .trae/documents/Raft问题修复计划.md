# Raft 算法问题修复计划

## 核心原则：每层完全独立，不复用任何接口

**禁止事项**：
- ❌ Raft 层不能 import `storage` 包
- ❌ Storage 层不能 import `raft` 包
- ❌ 不能复用 `istorage.IWal` 等其他层接口
- ❌ 不能共享 `LogEntry` 等跨层结构

**正确的做法**：
- ✅ 每层有自己的完整实现
- ✅ 通过 main.go 聚合各层（单向引用）
- ✅ 通过 Go 通道通信（不是接口）

---

## 架构设计

```
cmd/server/main.go  ←  聚合各层
       ↓
       ├──→ Raft (internal/Raft/)      ← 独立实现 WAL、LogEntry
       ├──→ Storage (internal/storage/) ← 独立实现 WAL、Engine
       └──→ Service (internal/service/)  ← 独立实现 Command、FSM

通信方式：
Raft     → applyCh (chan raft.LogEntry) → Service
Service  → applyCh (chan service.Command) → Storage
```

---

## Phase 1: Raft 层完善

### 步骤 1.1: Raft WAL（独立实现）

**文件**: `internal/Raft/raft_wal.go`

**不能复用 storage 层的 WAL，必须自己实现**：

```go
package Raft

// RaftWAL 自己实现的 WAL，不依赖任何其他包
type RaftWAL struct {
    file     *os.File
    logPath  string
    metaPath string
}

// 只依赖 Go 标准库
import (
    "encoding/binary"
    "os"
)
```

### 步骤 1.2: Raft LogEntry（独立定义）

```go
package Raft

// LogEntry Raft 自己的日志结构
type LogEntry struct {
    Index   int
    Term    int
    Command []byte  // 序列化的命令字节
}
```

### 步骤 1.3: 修复选举逻辑竞态

移除 `time.Sleep`，使用通道超时机制。

---

## Phase 2: Storage 层完善

### 步骤 2.1: 添加 applyCh

**文件**: `internal/storage/engine.go`

**不能引用 raft 包，独立实现**：

```go
package storage

// StorageCommand Storage 自己的命令
type StorageCommand struct {
    Key   []byte
    Value []byte
}

type Engine struct {
    memTable IMemTable
    mu       sync.RWMutex
    applyCh  chan StorageCommand  // 自己独立的通道
}
```

---

## Phase 3: Service 层（胶水层）

### 步骤 3.1: FSM 实现

**文件**: `internal/service/fsm.go`

**可以单向引用 Raft 和 Storage**：

```go
package service

// Command Service 自己的命令
type Command struct {
    Type  string  // "Put" or "Delete"
    Key   []byte
    Value []byte
}

type FSM struct {
    // 单向引用，不共享接口
    raft    *Raft.Raft      // 只能调用 Raft 暴露的方法
    storage *storage.Engine  // 只能调用 Storage 暴露的方法
}
```

---

## 文件结构

```
internal/Raft/
├── raft.go           # 核心算法，自己实现 WAL，自己定义 LogEntry
├── rpc.go            # RPC 定义
├── raft_wal.go       # 独立的 WAL 实现（不引用 storage）
└── raft_test.go

internal/storage/
├── engine.go         # 自己的 applyCh，自己的命令格式
├── istorage/
│   └── interfaces.go # 自己内部的接口
└── zstorage/
    ├── memtable.go   # 自己实现
    └── wal.go        # 自己实现（不引用 raft）

internal/service/
└── fsm.go           # 胶水层，引用 Raft 和 Storage

cmd/server/
└── main.go          # 聚合各组件
```

---

## 实施顺序

1. **Phase 1.1** - Raft WAL 独立实现
2. **Phase 1.3** - 修复选举逻辑
3. **Phase 2.1** - Storage 添加 applyCh
4. **Phase 3.1** - Service FSM 胶水层
