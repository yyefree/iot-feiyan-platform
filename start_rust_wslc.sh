#!/bin/bash
set -e

echo "=========================================="
echo "IoT飞燕平台 - Rust 启动脚本 (WSLC)"
echo "=========================================="
echo ""

# 检查 WSL2
if ! wsl --list -q &>/dev/null; then
    echo "警告: 未检测到 WSL2 环境"
fi

# 检查 Docker
if ! docker info &>/dev/null; then
    echo "错误: Docker Desktop 未运行"
    exit 1
fi

echo "✓ Docker Desktop 运行中"
echo ""

# 停止现有容器
echo "停止现有容器..."
docker compose -f docker-compose.wslc.yml down 2>/dev/null || true
echo ""

# 启动所有服务
echo "启动所有服务..."
docker compose -f docker-compose.wslc.yml up -d --build 2>&1
echo ""

# 等待服务就绪
echo "等待服务就绪..."
sleep 20

# 显示容器状态
echo ""
echo "=========================================="
echo "容器状态"
echo "=========================================="
docker ps --filter "name=iot-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "=========================================="
echo "服务访问地址"
echo "=========================================="
echo "前端界面:        http://localhost:3000"
echo "API网关:         http://localhost:8080/health"
echo "设备服务:        http://localhost:8081/health"
echo "产品服务:        http://localhost:8082/health"
echo "规则服务:        http://localhost:8083/health"
echo "数据服务:        http://localhost:8084/health"
echo "运营服务:        http://localhost:8085/health"
echo "EMQX Dashboard:  http://localhost:18083"
echo "=========================================="
