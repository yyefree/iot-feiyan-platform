@echo off
cd /d D:\AI\test\services\ota-service
set DATABASE_URL=postgres://iot_admin:iot_admin_2024@localhost:5432/iot_platform?sslmode=disable
set HOST=0.0.0.0
set PORT=8086
start /b ota-service.exe
timeout /t 3 /nobreak > nul

cd /d D:\AI\test\services\voice-service
set DATABASE_URL=postgres://iot_admin:iot_admin_2024@localhost:5432/iot_platform?sslmode=disable
set HOST=0.0.0.0
set PORT=8087
start /b voice-service.exe
timeout /t 3 /nobreak > nul

echo OTA服务启动: http://localhost:8086/health
echo 语音服务启动: http://localhost:8087/health
