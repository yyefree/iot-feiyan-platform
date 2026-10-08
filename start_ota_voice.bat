@echo off
echo ==========================================
echo 启动OTA和语音服务
echo ==========================================

echo 启动OTA服务...
cd /d D:\AI\test\services\ota-service
start /b simple_service.exe
timeout /t 2 /nobreak > nul
echo OTA服务已启动 (端口: 8086)

echo 启动语音服务...
cd /d D:\AI\test\services\voice-service
start /b simple_service.exe
timeout /t 2 /nobreak > nul
echo 语音服务已启动 (端口: 8087)

echo.
echo ==========================================
echo 服务启动完成
echo ==========================================
echo OTA服务: http://localhost:8086/health
echo 语音服务: http://localhost:8087/health
echo.
echo 等待服务启动...
timeout /t 3 /nobreak > nul

echo.
echo 测试健康检查...
curl -s http://localhost:8086/health
echo.
curl -s http://localhost:8087/health
echo.
