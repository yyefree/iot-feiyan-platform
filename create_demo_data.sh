#!/bin/bash
# IoT飞燕平台 - 演示数据脚本
# 运行方式: bash create_demo_data.sh

set -e

echo "=========================================="
echo "IoT飞燕平台 - 演示数据创建"
echo "=========================================="
echo ""

# 1. 创建租户
echo "【1. 创建租户】"
TENANT_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{"tenant_name":"演示企业","contact":"管理员","phone":"13800138000"}')
echo "租户响应: $TENANT_RESPONSE"
TENANT_ID=$(echo $TENANT_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
echo "租户ID: $TENANT_ID"
echo ""

# 2. 创建产品
echo "【2. 创建产品】"
PRODUCT_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/products \
  -H "Content-Type: application/json" \
  -d "{\"product_name\":\"智能温湿度传感器\",\"tenant_id\":$TENANT_ID,\"comm_type\":\"wifi\",\"auth_type\":\"one_device_one_key\"}")
echo "产品响应: $PRODUCT_RESPONSE"
PRODUCT_ID=$(echo $PRODUCT_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
PRODUCT_KEY=$(echo $PRODUCT_RESPONSE | grep -o '"product_key":"[^"]*"' | cut -d'"' -f4)
echo "产品ID: $PRODUCT_ID, 产品密钥: $PRODUCT_KEY"
echo ""

# 3. 创建设备
echo "【3. 创建设备】"
DEVICE_RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"温度传感器-001\",\"product_id\":$PRODUCT_ID,\"tenant_id\":$TENANT_ID,\"device_key\":\"$PRODUCT_KEY\"}")
echo "设备响应: $DEVICE_RESPONSE"
DEVICE_ID=$(echo $DEVICE_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
echo "设备ID: $DEVICE_ID"
echo ""

# 4. 添加物模型属性
echo "【4. 添加物模型属性】"
curl -s -X POST http://localhost:8082/api/v1/tsl/properties \
  -H "Content-Type: application/json" \
  -d "{\"product_id\":$PRODUCT_ID,\"identifier\":\"temperature\",\"name\":\"温度\",\"type\":\"float\",\"unit\":\"℃\",\"min\":-40,\"max\":80,\"step\":0.1}" > /dev/null

curl -s -X POST http://localhost:8082/api/v1/tsl/properties \
  -H "Content-Type: application/json" \
  -d "{\"product_id\":$PRODUCT_ID,\"identifier\":\"humidity\",\"name\":\"湿度\",\"type\":\"float\",\"unit\":\"%RH\",\"min\":0,\"max\":100,\"step\":0.1}" > /dev/null

curl -s -X POST http://localhost:8082/api/v1/tsl/properties \
  -H "Content-Type: application/json" \
  -d "{\"product_id\":$PRODUCT_ID,\"identifier\":\"battery\",\"name\":\"电量\",\"type\":\"int\",\"unit\":\"%\",\"min\":0,\"max\":100,\"step\":1}" > /dev/null
echo "物模型属性已添加"
echo ""

# 5. 上报遥测数据
echo "【5. 上报遥测数据】"
for i in {1..10}; do
  curl -s -X POST http://localhost:8084/api/v1/telemetry \
    -H "Content-Type: application/json" \
    -d "{\"device_id\":$DEVICE_ID,\"data\":{\"temperature\":$((20 + i)),\"humidity\":$((50 + i)),\"battery\":$((90 - i))}}" > /dev/null
done
echo "已上报10条遥测数据"
echo ""

# 6. 创建规则引擎
echo "【6. 创建规则引擎】"
RULE_RESPONSE=$(curl -s -X POST http://localhost:8083/api/v1/rules \
  -H "Content-Type: application/json" \
  -d "{\"rule_name\":\"温度超限告警\",\"sql\":\"SELECT * FROM telemetry WHERE temperature > 35\",\"enabled\":true}")
echo "规则响应: $RULE_RESPONSE"
RULE_ID=$(echo $RULE_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
echo "规则ID: $RULE_ID"
echo ""

# 7. 创建场景联动
echo "【7. 创建场景联动】"
SCENE_RESPONSE=$(curl -s -X POST http://localhost:8083/api/v1/scenes \
  -H "Content-Type: application/json" \
  -d "{\"scene_name\":\"回家模式\",\"trigger_type\":\"manual\",\"actions\":[{\"device_id\":$DEVICE_ID,\"action\":\"set_temperature\",\"value\":25}]}")
echo "场景响应: $SCENE_RESPONSE"
SCENE_ID=$(echo $SCENE_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
echo "场景ID: $SCENE_ID"
echo ""

# 8. 创建固件
echo "【8. 创建OTA固件】"
FIRMWARE_RESPONSE=$(curl -s -X POST http://localhost:8088/api/v1/ota/firmware \
  -H "Content-Type: application/json" \
  -d "{\"firmware_name\":\"v1.2.0\",\"product_id\":$PRODUCT_ID,\"version\":\"1.2.0\",\"download_url\":\"https://example.com/firmware/v1.2.0.bin\",\"size\":1024000}")
echo "固件响应: $FIRMWARE_RESPONSE"
FIRMWARE_ID=$(echo $FIRMWARE_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
echo "固件ID: $FIRMWARE_ID"
echo ""

# 9. 绑定语音设备
echo "【9. 绑定语音设备】"
curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_id\":$DEVICE_ID,\"voice_name\":\"客厅温度传感器\",\"platform\":\"tmall\"}" > /dev/null
echo "语音设备已绑定"
echo ""

# 10. 更新设备位置
echo "【10. 更新设备位置】"
curl -s -X POST http://localhost:8094/api/v1/devices/$DEVICE_ID/location \
  -H "Content-Type: application/json" \
  -d "{\"latitude\":31.2304,\"longitude\":121.4737,\"location\":\"上海市浦东新区\",\"region\":\"华东\"}" > /dev/null
echo "设备位置已更新"
echo ""

# 11. 查看统计信息
echo "【11. 查看统计信息】"
echo "=== 产品统计 ==="
curl -s http://localhost:8082/api/v1/products/statistics
echo ""
echo ""
echo "=== 数据统计 ==="
curl -s http://localhost:8084/api/v1/statistics/overview
echo ""
echo ""
echo "=== 运营大屏 ==="
curl -s http://localhost:8085/api/v1/dashboard
echo ""
echo ""
echo "=== OTA统计 ==="
curl -s http://localhost:8088/api/v1/ota/statistics
echo ""
echo ""
echo "=== 语音统计 ==="
curl -s http://localhost:8089/api/v1/voice/statistics
echo ""
echo ""

echo "=========================================="
echo "演示数据创建完成！"
echo "=========================================="
echo ""
echo "演示数据汇总:"
echo "  - 租户: $TENANT_ID (演示企业)"
echo "  - 产品: $PRODUCT_ID (智能温湿度传感器)"
echo "  - 设备: $DEVICE_ID (温度传感器-001)"
echo "  - 规则: $RULE_ID (温度超限告警)"
echo "  - 场景: $SCENE_ID (回家模式)"
echo "  - 固件: $FIRMWARE_ID (v1.2.0)"
echo ""
echo "访问地址:"
echo "  - API网关: http://localhost:8080"
echo "  - 前端页面: http://localhost:3000"
echo "  - EMQX Dashboard: http://localhost:18083"
