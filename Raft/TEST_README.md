# Raft 持久化测试说明

## 测试文件位置
`E:\Code\BanKv\Raft\raft_test.go`

## 运行测试

### 方法 1: 使用批处理脚本（推荐）
```bash
cd E:\Code\BanKv\Raft
run-tests.bat
```

### 方法 2: 使用 go test 命令
```bash
# 运行所有测试
go test -v ./Raft

# 运行特定测试
go test -v -run TestPersistenceTermAndVotedFor
go test -v -run TestSnapshotCreation
```

## 测试覆盖范围

### 1. 基础功能测试
- ✅ `TestNewRaft` - Raft 实例创建
- ✅ `TestGetState` - 状态查询
- ✅ `TestGetLog` - 日志查询
- ✅ `TestAppendEntry` - 日志追加（非 Leader 模式）
- ✅ `TestElectionTimeout` - 选举超时触发
- ✅ `TestLeaderAppendsLog` - Leader 日志追加
- ✅ `TestLeaderSendsHeartbeats` - Leader 心跳发送

### 2. 持久化功能测试（新增）

#### TestPersistenceTermAndVotedFor
**测试目标**: 验证 Term 和 votedFor 的持久化
- 设置 Term=5, votedFor=0
- 调用 persistLocked() 持久化
- 重新创建 Raft 实例
- 验证从磁盘加载的状态正确

**预期结果**: 
- r2.Term == 5
- r2.votedFor == 0

---

#### TestPersistenceLog
**测试目标**: 验证日志条目的持久化
- 添加 3 条日志条目
- 调用 persistLocked() 持久化
- 重新创建 Raft 实例
- 验证日志完整恢复

**预期结果**:
- len(r2.log) == 3
- 日志内容完全一致

---

#### TestSnapshotCreation
**测试目标**: 验证快照创建功能
- 添加日志并设置 commitIndex
- 调用 TakeSnapshot(1, data)
- 验证快照元数据更新

**预期结果**:
- LastIncludedIndex == 1
- lastSnapshotIndex == 1
- 无错误返回

---

#### TestSnapshotPersistence
**测试目标**: 验证快照持久化和恢复
- 创建包含索引 0-1 的快照
- 验证内存日志被截断
- 重新创建 Raft 实例
- 验证快照元数据恢复

**预期结果**:
- LastIncludedIndex == 1
- commitIndex == 1
- lastApplied == 1
- 内存日志只剩索引 2

---

#### TestInstallSnapshotRPC
**测试目标**: 验证 InstallSnapshot RPC 处理
- 模拟接收 InstallSnapshot RPC
- 验证状态更新（Term、LastIncludedIndex 等）
- 验证日志截断

**预期结果**:
- reply.Success == true
- Term 更新为 2
- LastIncludedIndex == 1
- commitIndex == 1

---

#### TestPersistAfterElection
**测试目标**: 验证选举后的持久化
- 等待单节点成为 Leader
- 记录当前 Term
- 重新创建 Raft 实例
- 验证 Term 已持久化

**预期结果**:
- newTerm >= originalTerm

---

#### TestPersistAfterAppendEntry
**测试目标**: 验证 AppendEntry 后的持久化
- 等待成为 Leader
- 追加一条日志
- 重新创建 Raft 实例
- 验证日志已持久化

**预期结果**:
- 日志不为空
- 最后一条日志内容正确

---

## 关键实现细节

### 持久化函数
```go
// persistLocked - 持久化所有 Raft 状态
func (r *Raft) persistLocked() {
    data := PersistData{
        CurrentTerm:       r.Term,
        VotedFor:          r.votedFor,
        Log:               r.log,
        LastIncludedIndex: r.LastIncludedIndex,
        LastIncludedTerm:  r.LastIncludedTerm,
    }
    r.wal.SavePersist(data)
}

// readPersist - 从磁盘加载所有状态
func (r *Raft) readPersist() error {
    data, err := r.wal.LoadPersist()
    // ... 恢复状态
}
```

### 持久化时机
根据 Raft 论文 Figure 2，以下情况必须持久化：
1. ✅ Term 改变时（选举开始）
2. ✅ votedFor 改变时（投票）
3. ✅ 日志追加时（AppendEntry）
4. ✅ 接收日志时（AppendEntries RPC）
5. ✅ 快照创建时（TakeSnapshot）
6. ✅ 安装快照时（InstallSnapshot RPC）

---

## 注意事项

### 测试环境清理
每个持久化测试都会：
1. 删除旧的 `raft_data` 目录
2. 测试结束后自动清理

### 单节点模式
大部分测试使用单节点配置 (`peers = ["localhost:8000"]`)，这样：
- 选举会立即成功（自己一票即过半数）
- 日志追加会立即提交
- 无需网络通信

### 时间依赖
某些测试需要等待选举超时（约 400ms），这是正常的 Raft 行为。

---

## 预期输出示例

```
=== RUN   TestPersistenceTermAndVotedFor
✓ Term and votedFor persistence test passed
--- PASS: TestPersistenceTermAndVotedFor (0.01s)

=== RUN   TestPersistenceLog
✓ Log persistence test passed
--- PASS: TestPersistenceLog (0.01s)

=== RUN   TestSnapshotCreation
[RAFT] Snapshot created: Index=1, Term=1
✓ Snapshot creation test passed
--- PASS: TestSnapshotCreation (0.01s)

...

PASS
ok      github.com/NeverENG/BanKV/Raft  1.234s
```

---

## 故障排查

### 测试失败可能原因
1. **WAL 文件权限问题** - 确保有读写权限
2. **端口占用** - 虽然测试不使用网络，但地址格式需有效
3. **磁盘空间不足** - 检查可用空间
4. **Go 版本问题** - 建议使用 Go 1.18+

### 调试技巧
```bash
# 查看详细日志
go test -v -run TestPersistenceLog

# 只运行一个测试
go test -v -run "^TestPersistenceTermAndVotedFor$"

# 显示测试覆盖率
go test -cover -run TestPersistence
```

---

## 下一步改进

1. **多节点测试** - 测试集群环境下的持久化
2. **并发测试** - 测试并发写入时的持久化一致性
3. **崩溃恢复测试** - 模拟进程崩溃后的恢复
4. **性能测试** - 测试持久化的性能影响
5. **边界条件测试** - 空日志、超大日志等场景
