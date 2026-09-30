#!/bin/bash
# IoT飞燕平台 - Rust 服务一键启动脚本

echo "=========================================="
echo "IoT飞燕平台 - Rust 服务启动"
echo "=========================================="
echo ""

# 检查 Docker
if ! docker info &>/dev/null; then
    echo "错误: Docker Desktop 未运行"
    echo "请先启动 Docker Desktop"
    exit 1
fi

echo "✓ Docker Desktop 运行中"
echo ""

# 停止现有容器
echo "停止现有容器..."
docker compose -f docker-compose.wslc.yml down 2>/dev/null || true
echo ""

# 启动所有服务
echo "启动所有 Rust 服务..."
docker compose -f docker-compose.wslc.yml up -d --build 2>&1
echo ""

# 等待服务就绪
echo "等待服务就绪..."
sleep 25

# 显示容器状态
echo ""
echo "=========================================="
echo "容器状态"
echo "=========================================="
docker ps --filter "name=iot-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

# 健康检查
echo "=========================================="
echo "服务健康检查"
echo "=========================================="
for svc in "网关:8080" "设备:8081" "产品:8082" "规则:8083" "数据:8084" "运营:8085"; do
    name="${svc%%:*}"
    port="${svc##*:}"
    if curl -s --max-time 5 "http://localhost:$port/health" | grep -q '"code":0'; then
        echo "✓ $name 服务正常 (:$port)"
    else
        echo "✗ $name 服务异常 (:$port)"
    fi
done

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
