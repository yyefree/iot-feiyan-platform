@echo off
echo ==========================================
echo IoT飞燕平台 - 服务启动脚本
echo ==========================================
echo.

echo 启动规则服务扩展...
cd /d D:\AI\test\services\rule-service-extensions
start /b go run main.go
timeout /t 3 /nobreak > nul

echo.
echo ==========================================
echo 服务启动完成
echo ==========================================
echo 规则服务扩展: http://localhost:8094/health
echo.
