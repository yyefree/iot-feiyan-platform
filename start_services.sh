#!/bin/bash

echo "=========================================="
echo "IoT飞燕平台 - 服务启动脚本"
echo "=========================================="
echo ""

# 启动规则服务扩展
echo "启动规则服务扩展..."
cd D:\AI\test\services\rule-service-extensions
go run main.go &
RULE_EXT_PID=$!
echo "规则服务扩展PID: $RULE_EXT_PID (端口: 8094)"

echo ""
echo "=========================================="
echo "服务启动完成"
echo "=========================================="
echo "规则服务扩展: http://localhost:8094/health"
echo ""
