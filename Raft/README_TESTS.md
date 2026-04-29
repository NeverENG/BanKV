# Raft 测试指南

## 📁 文件说明

- `raft_test.go` - 完整的测试套件（14个测试用例）
- `run-tests.bat` - 完整测试脚本（逐个运行所有测试）
- `quick-test.bat` - 快速测试脚本（一次性运行所有测试）
- `TEST_README.md` - 详细的测试说明文档
- `TEST_REPORT.md` - 测试报告和功能验证清单

## 🚀 快速开始

### 方法 1: 一键测试（最简单）
```bash
cd E:\Code\BanKv\Raft
quick-test.bat
```

### 方法 2: 详细测试（推荐）
```bash
cd E:\Code\BanKv\Raft
run-tests.bat
```

### 方法 3: 命令行测试
```bash
# 运行所有测试
go test -v ./Raft

# 运行特定测试
go test -v -run TestPersistenceLog

# 运行持久化相关测试
go test -v -run TestPersistence

# 运行快照相关测试
go test -v -run TestSnapshot
```

## ✅ 测试覆盖

### 基础功能 (7个测试)
- ✓ Raft 实例创建
- ✓ 状态查询
- ✓ 日志操作
- ✓ 选举机制
- ✓ Leader 功能

### 持久化功能 (7个测试) ⭐
- ✓ Term 和 votedFor 持久化
- ✓ 日志持久化
- ✓ 快照创建
- ✓ 快照持久化和恢复
- ✓ InstallSnapshot RPC
- ✓ 选举后持久化
- ✓ AppendEntry 后持久化

## 📊 预期输出

```
========================================
Running Raft Persistence Tests
========================================

[1/8] Testing Term and VotedFor Persistence...
=== RUN   TestPersistenceTermAndVotedFor
✓ Term and votedFor persistence test passed
--- PASS: TestPersistenceTermAndVotedFor (0.01s)
PASS

[2/8] Testing Log Persistence...
=== RUN   TestPersistenceLog
✓ Log persistence test passed
--- PASS: TestPersistenceLog (0.01s)
PASS

...

========================================
All tests completed!
========================================
```

## 🔍 测试详解

查看以下文档了解每个测试的详细信息：
- **TEST_README.md** - 测试目标、流程、预期结果
- **TEST_REPORT.md** - 完整测试报告和验证清单

## ⚙️ 测试环境

- Go 1.18+
- Windows/Linux/macOS
- 至少 10MB 磁盘空间

## 🐛 问题排查

如果测试失败：

1. **清理旧数据**
   ```bash
   rm -rf raft_data  # Linux/Mac
   rmdir /s /q raft_data  # Windows
   ```

2. **检查 Go 版本**
   ```bash
   go version
   ```

3. **查看详细错误**
   ```bash
   go test -v -run TestName
   ```

## 📖 更多信息

- Raft 论文: https://raft.github.io/raft.pdf
- MIT 6.824: https://pdos.csail.mit.edu/6.824/

---

**祝测试顺利！** 🎉
