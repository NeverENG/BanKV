# Storage 层问题修复计划

## 问题概述

经过分析，storage 层存在以下问题：

### 问题 1: istorage 包不存在（严重）
- **位置**: [engine.go](file:///e:/Code/BanKv/internal/storage/engine.go#L7-L8) 和 [WAL.go](file:///e:/Code/BanKv/internal/storage/zstorage/WAL.go#L9-L10)
- **描述**: 两个文件都引用了 `github.com/NeverENG/BanKV/internal/storage/istorage` 包，但该包不存在
- **影响**: 代码无法编译

### 问题 2: WAL.go 中的接口使用不一致
- **位置**: [WAL.go](file:///e:/Code/BanKv/internal/storage/zstorage/WAL.go#L31-L31)
- **描述**: `Write` 方法接收 `istorage.LogEntry` 类型，但 `zstorage` 包内已有自己的 `LogEntry` 类型定义在 [model.go](file:///e:/Code/BanKv/internal/storage/zstorage/model.go)

### 问题 3: MemTable 未实现任何接口
- **描述**: `MemTable` 结构体有 `FlushChan`、`wal` 等字段，但 WAL 引用的是接口 `istorage.IWal`，而 MemTable 使用的是具体的 `*WAL` 类型

---

## 修复方案（采用创建 istorage 接口包）

### 步骤 1: 创建 istorage 接口包
**文件**: `internal/storage/istorage/interfaces.go`

创建以下接口和结构体：
- `LogEntry` 结构体
- `IWal` 接口：Write, Read, Close, Sync, Clear
- `IMemTable` 接口：Get, Put, Delete, Size, StartFlush

### 步骤 2: 修复 engine.go
- 保持 import `istorage` 包
- 确保 Engine 使用 `istorage.IMemTable` 接口

### 步骤 3: 修复 WAL.go
- 使用 `istorage.LogEntry` 替代 `zstorage.LogEntry`
- 确保 WAL 实现 `istorage.IWal` 接口
- 删除 `zstorage/model.go` 中的 LogEntry（如果不再需要）

### 步骤 4: 修复 memtable.go
- 让 MemTable 实现 `istorage.IMemTable` 接口

### 步骤 5: 创建黑盒测试
**文件**:
- `internal/storage/zstorage/memtable_test.go` - 已存在，需检查
- `internal/storage/zstorage/wal_test.go` - 新建
- `internal/storage/engine_test.go` - 新建

测试要求：
- 黑盒测试，不直接访问内部实现
- 通过公共接口进行测试
- 覆盖 Put/Get/Delete 基本操作
- 覆盖 Flush 逻辑

### 步骤 6: 验证编译和测试
- 运行 `go build ./...` 确保编译通过
- 运行 `go test ./...` 确保测试通过

---

## 修复后预期
- 所有 storage 层代码可以正常编译
- 接口定义清晰，模块间解耦
- 测试覆盖完整
- 目录结构正确
