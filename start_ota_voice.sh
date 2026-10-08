#!/bin/bash

echo "=========================================="
echo "启动OTA和语音服务"
echo "=========================================="

# 启动OTA服务
echo "启动OTA服务..."
cd D:\AI\test\services\ota-service
export DATABASE_URL="postgres://iot_admin:iot_admin_2024@localhost:5432/iot_platform?sslmode=disable"
export HOST="0.0.0.0"
export PORT="8086"
./simple_service &
OTA_PID=$!
echo "OTA服务PID: $OTA_PID"

# 启动语音服务
echo "启动语音服务..."
cd D:\AI\test\services\voice-service
export DATABASE_URL="postgres://iot_admin:iot_admin_2024@localhost:5432/iot_platform?sslmode=disable"
export HOST="0.0.0.0"
export PORT="8087"
./simple_service &
VOICE_PID=$!
echo "语音服务PID: $VOICE_PID"

echo ""
echo "=========================================="
echo "服务启动完成"
echo "=========================================="
echo "OTA服务: http://localhost:8086/health (PID: $OTA_PID)"
echo "语音服务: http://localhost:8087/health (PID: $VOICE_PID)"
echo ""
echo "等待服务启动..."
sleep 3

# 测试健康检查
echo ""
echo "测试健康检查..."
curl -s http://localhost:8086/health
echo ""
curl -s http://localhost:8087/health
echo ""
