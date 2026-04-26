# BanKV 快速启动指南

## 🚀 快速开始（3步）

### 第一步：启动服务端

**⚠️ 重要：** 如果是第一次启动或需要清空数据，使用清理模式启动

**方法1：使用启动脚本（推荐）**
```powershell
cd E:\Code\BanKv\cmd\server
.\run.bat

# 或清理旧数据后启动
.\run-clean.bat
```

**方法2：手动启动**
```powershell
cd E:\Code\BanKv
go run cmd/server/server.go
```

看到以下输出表示启动成功：
```
Starting server...
HA initialized, initial health status: true
```

---

### 第二步：启动客户端（新开一个终端）

**方法1：使用启动脚本（推荐）**
```powershell
cd E:\Code\BanKv\cmd\client
.\run.bat
```

**方法2：手动启动**
```powershell
cd E:\Code\BanKv
go run cmd/client
```

看到以下输出表示连接成功：
```
已连接到 localhost:8080
输入命令进行操作，输入 'quit' 或 'exit' 退出
支持命令: put <key> <value>, get <key>, delete <key>

> 
```

---

### 第三步：测试操作

在客户端中输入：

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

---

## ⚠️ 常见错误

### 错误1：`undefined: NewClient`

**原因：** 直接运行了 `go run main.go`

**解决：** 运行整个包
```powershell
# ❌ 错误
go run main.go

# ✅ 正确
go run .
# 或
go run main.go client.go interactive.go
```

---

### 错误2：配置文件读取失败

**原因：** 不在正确的目录运行

**解决：** 从项目根目录运行
```powershell
cd E:\Code\BanKv
go run cmd/server/server.go
```

现在代码已经修复，会自动查找配置文件。

---

### 错误3：连接失败

**原因：** 服务端未启动或端口被占用

**解决：**
1. 确保服务端已启动
2. 检查端口 8080 是否被占用
3. 查看服务端日志是否有错误

---

## 📝 完整命令参考

### 服务端

```powershell
# 启动服务端
cd E:\Code\BanKv\cmd\server
.\run.bat

# 或
cd E:\Code\BanKv
go run cmd/server/server.go
```

### 客户端 - 交互模式

```powershell
# 启动客户端
cd E:\Code\BanKv\cmd\client
.\run.bat

# 或指定服务器地址
.\run.bat localhost:8080

# 或使用 go run
cd E:\Code\BanKv
go run cmd/client
```

### 客户端 - 命令行模式

```powershell
cd E:\Code\BanKv

# PUT
go run cmd/client put name Alice

# GET
go run cmd/client get name

# DELETE
go run cmd/client delete name
```

---

## 🎯 下一步

- 查看 [客户端详细文档](README.md)
- 了解 [项目架构](../../README.md)
- 探索更多功能

---

## 💡 提示

1. **保持两个终端**：一个运行服务端，一个运行客户端
2. **使用启动脚本**：最简单，不容易出错
3. **输入 help**：在客户端中随时查看帮助
4. **Ctrl+C**：可以快速退出客户端

祝您使用愉快！🎉
