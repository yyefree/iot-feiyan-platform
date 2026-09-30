#!/bin/bash

echo "=========================================="
echo "IoT飞燕平台 - 服务启动与功能测试"
echo "=========================================="
echo ""

# 检查Docker是否运行
echo "1. 检查Docker状态..."
if ! docker info > /dev/null 2>&1; then
    echo -e "\033[0;31m错误: Docker未运行，请先启动Docker Desktop\033[0m"
    exit 1
fi
echo -e "\033[0;32m✓ Docker运行正常\033[0m"
echo ""

# 停止并清理旧容器
echo "2. 清理旧容器..."
docker-compose down 2>/dev/null
echo -e "\033[0;32m✓ 清理完成\033[0m"
echo ""

# 启动基础设施
echo "3. 启动基础设施服务..."
docker-compose up -d postgres redis influxdb emqx
sleep 10
echo -e "\033[0;32m✓ 基础设施启动完成\033[0m"
echo ""

# 检查基础设施状态
echo "4. 检查基础设施状态..."
docker ps --filter "name=iot-" --format "table {{.Names}}\t{{.Status}}" | grep -E "postgres|redis|influxdb|emqx"
echo ""

# 初始化数据库
echo "5. 初始化数据库..."
docker exec -i iot-postgres psql -U iot_admin -d iot_platform -c "SELECT 1" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✓ 数据库连接正常\033[0m"
else
    echo -e "\033[0;31m✗ 数据库初始化失败\033[0m"
    exit 1
fi
echo ""

# 启动后端服务
echo "6. 启动后端服务..."
cd gateway && go run main.go &
cd ../services/device-service && go run main.go &
cd ../product-service && go run main.go &
cd ../rule-service && go run main.go &
cd ../data-service && go run main.go &
cd ../ops-service && go run main.go &
sleep 5
echo -e "\033[0;32m✓ 后端服务启动完成\033[0m"
echo ""

# 检查后端服务状态
echo "7. 检查后端服务状态..."
for port in 8080 8081 8082 8083 8084 8085; do
    if curl -s http://localhost:$port/health > /dev/null 2>&1; then
        echo -e "  端口 $port: \033[0;32m运行中\033[0m"
    else
        echo -e "  端口 $port: \033[0;31m未启动\033[0m"
    fi
done
echo ""

# 启动前端
echo "8. 启动前端服务..."
cd ../frontend
npm install --silent 2>/dev/null
npm run dev &
sleep 5
echo -e "\033[0;32m✓ 前端服务启动完成\033[0m"
echo ""

# 运行测试
echo "9. 运行功能测试..."
echo ""
bash test_all.sh
echo ""

echo "=========================================="
echo "测试完成！"
echo "=========================================="
echo ""
echo "访问地址:"
echo "  - 前端: http://localhost:5173"
echo "  - 网关: http://localhost:8080"
echo "  - 设备服务: http://localhost:8081"
echo "  - 产品服务: http://localhost:8082"
echo "  - 规则服务: http://localhost:8083"
echo "  - 数据服务: http://localhost:8084"
echo "  - 运营服务: http://localhost:8085"
echo ""
echo "停止服务:"
echo "  - 按 Ctrl+C 停止所有服务"
echo "  - 或运行: docker-compose down"
echo ""
