@echo off
REM BanKV 客户端启动脚本

echo ========================================
echo    BanKV Interactive Client
echo ========================================
echo.

REM 检查是否提供了服务器地址
if "%1"=="" (
    echo 连接到默认服务器: localhost:8080
    echo.
    go run .
) else (
    echo 连接到服务器: %1
    echo.
    go run . %1
)
