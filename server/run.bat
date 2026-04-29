@echo off
REM BanKV 服务端启动脚本

echo ========================================
echo    BanKV Server
echo ========================================
echo.

cd ..\..
go run cmd/server/server.go
