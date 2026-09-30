#!/bin/bash
set -e

echo "=========================================="
echo "IoT飞燕平台 - Rust版本功能测试"
echo "=========================================="
echo ""

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

PASS=0
FAIL=0
SKIP=0

test_api() {
    local name=$1
    local url=$2
    local method=${3:-GET}
    local data=${4:-''}
    
    echo -n "测试 $name ... "
    
    local response http_code body
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" --max-time 10 "$url" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" --max-time 10 "$url" 2>/dev/null)
    fi
    
    http_code=$(echo "$response" | tail -1)
    
    if [ "$http_code" = "200" ] || [ "$http_code" = "201" ] || [ "$http_code" = "204" ]; then
        echo -e "${GREEN}✓ 通过${NC}"
        PASS=$((PASS+1))
    elif [ "$http_code" = "000" ]; then
        echo -e "${YELLOW}○ 跳过（服务未启动）${NC}"
        SKIP=$((SKIP+1))
    else
        echo -e "${RED}✗ 失败 (HTTP $http_code)${NC}"
        FAIL=$((FAIL+1))
    fi
}

echo "1. 健康检查测试"
echo "----------------"
test_api "网关健康检查" "http://localhost:8080/health"
test_api "设备服务健康检查" "http://localhost:8081/health"
test_api "产品服务健康检查" "http://localhost:8082/health"
test_api "规则服务健康检查" "http://localhost:8083/health"
test_api "数据服务健康检查" "http://localhost:8084/health"
test_api "运营服务健康检查" "http://localhost:8085/health"
echo ""

echo "2. 产品管理测试"
echo "----------------"
test_api "创建产品" "http://localhost:8082/api/v1/products" "POST" "{\"product_name\":\"Test Product $(date +%s)\",\"tenant_id\":1,\"comm_type\":\"wifi\",\"auth_type\":\"one_device_one_key\"}"
test_api "获取产品列表" "http://localhost:8082/api/v1/products"
test_api "获取产品统计" "http://localhost:8082/api/v1/products/statistics"
echo ""

echo "3. 设备管理测试"
echo "----------------"
test_api "创建设备" "http://localhost:8081/api/v1/devices" "POST" "{\"device_name\":\"Test Device $(date +%s)\",\"product_id\":1,\"tenant_id\":1}"
test_api "获取设备列表" "http://localhost:8081/api/v1/devices"
test_api "获取设备统计" "http://localhost:8081/api/v1/devices/statistics"
test_api "创建虚拟设备" "http://localhost:8081/api/v1/devices/virtual/create" "POST" "{\"device_name\":\"Virtual Device $(date +%s)\",\"product_id\":1,\"tenant_id\":1,\"is_virtual\":true}"
test_api "获取虚拟设备列表" "http://localhost:8081/api/v1/devices/virtual/list"
echo ""

echo "4. 物模型测试"
echo "----------------"
test_api "添加属性" "http://localhost:8082/api/v1/products/1/properties" "POST" "{\"identifier\":\"test_temp\",\"name\":\"温度\",\"type\":\"float\"}"
test_api "获取属性列表" "http://localhost:8082/api/v1/products/1/properties"
test_api "添加服务" "http://localhost:8082/api/v1/products/1/services" "POST" "{\"identifier\":\"test_svc\",\"name\":\"测试服务\"}"
test_api "获取服务列表" "http://localhost:8082/api/v1/products/1/services"
test_api "添加事件" "http://localhost:8082/api/v1/products/1/events" "POST" "{\"identifier\":\"test_evt\",\"name\":\"测试事件\",\"level\":\"info\"}"
test_api "获取事件列表" "http://localhost:8082/api/v1/products/1/events"
test_api "导出TSL" "http://localhost:8082/api/v1/products/1/tsl/export"
echo ""

echo "5. 规则引擎测试"
echo "----------------"
test_api "创建规则" "http://localhost:8083/api/v1/rules" "POST" "{\"name\":\"Test Rule $(date +%s)\",\"sql_expression\":\"SELECT device_id FROM devices WHERE online = true\",\"tenant_id\":1}"
test_api "获取规则列表" "http://localhost:8083/api/v1/rules"
test_api "测试SQL" "http://localhost:8083/api/v1/rules/sql/test" "POST" "{\"sql_expression\":\"SELECT * FROM devices WHERE online = true\"}"
test_api "获取规则统计" "http://localhost:8083/api/v1/rules/statistics"
echo ""

echo "6. 场景联动测试"
echo "----------------"
test_api "创建场景" "http://localhost:8083/api/v1/scenes" "POST" "{\"name\":\"Test Scene $(date +%s)\",\"cron_expr\":\"0 18 * * *\",\"tenant_id\":1}"
test_api "获取场景列表" "http://localhost:8083/api/v1/scenes"
test_api "获取日出日落" "http://localhost:8083/api/v1/scenes/sunrise"
test_api "测试场景" "http://localhost:8083/api/v1/scenes/1/execute" "POST"
echo ""

echo "7. 数据统计测试"
echo "----------------"
test_api "上报遥测数据" "http://localhost:8084/api/v1/data/telemetry" "POST" "{\"device_id\":1,\"identify\":\"temperature\",\"value_float\":25.5}"
test_api "批量上报遥测" "http://localhost:8084/api/v1/data/telemetry/batch" "POST" "[{\"device_id\":1,\"identify\":\"temp\",\"value_float\":26.0},{\"device_id\":1,\"identify\":\"humid\",\"value_float\":60.0}]"
test_api "查询遥测数据" "http://localhost:8084/api/v1/data/telemetry?device_id=1"
test_api "获取最新遥测" "http://localhost:8084/api/v1/data/telemetry/device/1/latest"
test_api "获取统计概览" "http://localhost:8084/api/v1/data/statistics/overview"
test_api "获取折线图" "http://localhost:8084/api/v1/data/charts/line?device_id=1"
test_api "获取柱状图" "http://localhost:8084/api/v1/data/charts/bar?device_id=1"
test_api "获取饼图" "http://localhost:8084/api/v1/data/charts/pie"
test_api "获取雷达图" "http://localhost:8084/api/v1/data/charts/radar"
echo ""

echo "8. 运营中心测试"
echo "----------------"
test_api "获取运营大屏" "http://localhost:8085/api/v1/ops/dashboard"
test_api "获取设备统计" "http://localhost:8085/api/v1/ops/devices"
test_api "获取用户统计" "http://localhost:8085/api/v1/ops/users"
test_api "获取告警列表" "http://localhost:8085/api/v1/ops/alerts"
test_api "获取操作日志" "http://localhost:8085/api/v1/ops/logs"
test_api "获取大屏配置" "http://localhost:8085/api/v1/ops/dashboard/config"
echo ""

echo "9. 消息推送测试"
echo "----------------"
test_api "发送App推送" "http://localhost:8085/api/v1/push/app" "POST" "{\"tenant_id\":1,\"user_id\":1,\"msg_type\":\"alert\",\"title\":\"测试\",\"content\":\"测试内容\"}"
test_api "发送微信推送" "http://localhost:8085/api/v1/push/wechat" "POST" "{\"tenant_id\":1,\"user_id\":1,\"msg_type\":\"alert\",\"title\":\"测试\",\"content\":\"测试内容\"}"
test_api "获取推送模板" "http://localhost:8085/api/v1/push/templates"
echo ""

echo "=========================================="
echo "测试结果汇总"
echo "=========================================="
echo -e "通过: ${GREEN}$PASS${NC}"
echo -e "失败: ${RED}$FAIL${NC}"
echo -e "跳过: ${YELLOW}$SKIP${NC}"
echo "总计: $((PASS+FAIL+SKIP))"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}✗ 有 $FAIL 个测试失败${NC}"
    exit 1
fi
