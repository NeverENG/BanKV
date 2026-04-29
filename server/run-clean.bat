@echo off
REM BanKV 服务端启动脚本（清理模式）

echo ========================================
echo    BanKV Server (Clean Start)
echo ========================================
echo.

REM 清理旧的 WAL 日志
if exist ..\..\log\wal.log (
    echo 清理旧的 WAL 日志...
    del ..\..\log\wal.log
)

echo.
echo 启动服务端...
echo.

cd ..\..
go run cmd/server/server.go
