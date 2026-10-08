@echo off
echo ==========================================
echo 编译并启动OTA和语音服务
echo ==========================================

echo 编译OTA服务...
cd /d D:\AI\test\services\ota-service
go build -o simple_service.exe simple_ota.go
if %ERRORLEVEL% NEQ 0 (
    echo 编译失败
    exit /b 1
)
echo OTA服务编译成功

echo 编译语音服务...
cd /d D:\AI\test\services\voice-service
go build -o simple_service.exe simple_voice.go
if %ERRORLEVEL% NEQ 0 (
    echo 编译失败
    exit /b 1
)
echo 语音服务编译成功

echo.
echo 启动OTA服务...
cd /d D:\AI\test\services\ota-service
start "OTA Service" /b simple_service.exe
timeout /t 2 /nobreak > nul

echo 启动语音服务...
cd /d D:\AI\test\services\voice-service
start "Voice Service" /b simple_service.exe
timeout /t 2 /nobreak > nul

echo.
echo ==========================================
echo 服务启动完成
echo ==========================================
echo OTA服务: http://localhost:8088/health
echo 语音服务: http://localhost:8089/health
echo.
echo 测试健康检查...
curl -s http://localhost:8088/health
echo.
curl -s http://localhost:8089/health
echo.
