#!/bin/bash
# IoT飞燕平台 - 本地 Rust 服务启动脚本
# 不需要 Docker，直接运行 Rust 服务

set -e

echo "=========================================="
echo "IoT飞燕平台 - 本地 Rust 服务启动"
echo "=========================================="
echo ""

# 检查 Rust
if ! command -v cargo &> /dev/null; then
    echo "错误: 未检测到 Rust 工具链"
    echo "请先安装 Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
    exit 1
fi

# 检查 PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "警告: 未检测到本地 PostgreSQL，将使用 Docker 启动"
    docker compose -f docker-compose.wslc.yml up -d postgres redis influxdb emqx
    sleep 10
fi

# 设置环境变量
export DATABASE_URL="postgres://iot_admin:iot_admin_2024@localhost:5432/iot_platform"
export DEVICE_ADDR="http://localhost:8081"
export PRODUCT_ADDR="http://localhost:8082"
export RULE_ADDR="http://localhost:8083"
export DATA_ADDR="http://localhost:8084"
export OPS_ADDR="http://localhost:8085"

cd D:\AI\test\rust-services

echo "构建 Rust 服务..."
cargo build --release 2>&1
echo ""

echo "启动服务..."
echo ""

# 启动各服务（后台运行）
echo "启动 device-service (8081)..."
DATABASE_URL=$DATABASE_URL PORT=8081 cargo run --package iot-device-service --release &
echo "启动 product-service (8082)..."
DATABASE_URL=$DATABASE_URL PORT=8082 cargo run --package iot-product-service --release &
echo "启动 rule-service (8083)..."
DATABASE_URL=$DATABASE_URL PORT=8083 cargo run --package iot-rule-service --release &
echo "启动 data-service (8084)..."
DATABASE_URL=$DATABASE_URL PORT=8084 cargo run --package iot-data-service --release &
echo "启动 ops-service (8085)..."
DATABASE_URL=$DATABASE_URL PORT=8085 cargo run --package iot-ops-service --release &
echo "启动 api-gateway (8080)..."
DEVICE_ADDR=$DEVICE_ADDR PRODUCT_ADDR=$PRODUCT_ADDR RULE_ADDR=$RULE_ADDR DATA_ADDR=$DATA_ADDR OPS_ADDR=$OPS_ADDR cargo run --package iot-gateway --release &

echo ""
echo "等待服务启动..."
sleep 15

echo ""
echo "=========================================="
echo "服务状态检查"
echo "=========================================="
for port in 8080 8081 8082 8083 8084 8085; do
    if curl -s --max-time 3 "http://localhost:$port/health" | grep -q '"code":0'; then
        echo "✓ 服务 :$port 正常"
    else
        echo "✗ 服务 :$port 异常"
    fi
done

echo ""
echo "=========================================="
echo "服务访问地址"
echo "=========================================="
echo "前端界面:        http://localhost:3000 (需要单独启动)"
echo "API网关:         http://localhost:8080/health"
echo "设备服务:        http://localhost:8081/health"
echo "产品服务:        http://localhost:8082/health"
echo "规则服务:        http://localhost:8083/health"
echo "数据服务:        http://localhost:8084/health"
echo "运营服务:        http://localhost:8085/health"
echo "=========================================="
echo ""
echo "提示: 按 Ctrl+C 停止所有服务"
echo ""

# 等待所有后台进程
wait
