#!/bin/bash
# IoT飞燕平台 - 全面演示数据脚本
# 为所有服务添加演示数据

set -e

echo "=========================================="
echo "IoT飞燕平台 - 全面演示数据创建"
echo "=========================================="
echo ""

# ========== 1. 产品管理演示 ==========
echo "【1. 产品管理演示】"
echo ""

# 智能灯泡产品
echo "创建智能灯泡产品..."
BULB_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"product_name":"智能灯泡","tenant_id":1,"comm_type":"wifi","auth_type":"one_device_one_key","product_type":"direct_device","status":"published","thing_model":"{\"schema\":\"https://iot.aliyun.com/tsl/v1\",\"properties\":[{\"identifier\":\"brightness\",\"name\":\"亮度\",\"type\":\"int\",\"unit\":\"%\",\"min\":0,\"max\":100},{\"identifier\":\"color_temp\",\"name\":\"色温\",\"type\":\"int\",\"unit\":\"K\",\"min\":2700,\"max\":6500},{\"identifier\":\"power\",\"name\":\"开关\",\"type\":\"bool\"}]}"}')
BULB_ID=$(echo $BULB_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
BULB_KEY=$(echo $BULB_RESPONSE | grep -o '"product_key":"[^"]*"' | cut -d'"' -f4)
echo "  ✅ 智能灯泡产品创建成功 (ID: $BULB_ID, Key: $BULB_KEY)"

# 智能插座产品
echo "创建智能插座产品..."
SOCKET_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"product_name":"智能插座","tenant_id":1,"comm_type":"wifi","auth_type":"one_device_one_key","product_type":"direct_device","status":"published","thing_model":"{\"schema\":\"https://iot.aliyun.com/tsl/v1\",\"properties\":[{\"identifier\":\"power\",\"name\":\"开关\",\"type\":\"bool\"},{\"identifier\":\"power_value\",\"name\":\"功率\",\"type\":\"float\",\"unit\":\"W\"}]}"}')
SOCKET_ID=$(echo $SOCKET_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
SOCKET_KEY=$(echo $SOCKET_RESPONSE | grep -o '"product_key":"[^"]*"' | cut -d'"' -f4)
echo "  ✅ 智能插座产品创建成功 (ID: $SOCKET_ID, Key: $SOCKET_KEY)"

# 温湿度传感器产品
echo "创建温湿度传感器产品..."
SENSOR_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"product_name":"温湿度传感器","tenant_id":1,"comm_type":"zigbee","auth_type":"one_device_one_key","product_type":"direct_device","status":"published","thing_model":"{\"schema\":\"https://iot.aliyun.com/tsl/v1\",\"properties\":[{\"identifier\":\"temperature\",\"name\":\"温度\",\"type\":\"float\",\"unit\":\"℃\",\"min\":-40,\"max\":80},{\"identifier\":\"humidity\",\"name\":\"湿度\",\"type\":\"float\",\"unit\":\"%RH\",\"min\":0,\"max\":100},{\"identifier\":\"battery\",\"name\":\"电量\",\"type\":\"int\",\"unit\":\"%\",\"min\":0,\"max\":100}]}"}')
SENSOR_ID=$(echo $SENSOR_RESPONSE | grep -o '"id":[0-9]*' | cut -d: -f2)
SENSOR_KEY=$(echo $SENSOR_RESPONSE | grep -o '"product_key":"[^"]*"' | cut -d'"' -f4)
echo "  ✅ 温湿度传感器产品创建成功 (ID: $SENSOR_ID, Key: $SENSOR_KEY)"

echo ""

# ========== 2. 设备管理演示 ==========
echo "【2. 设备管理演示】"
echo ""

# 创建灯泡设备
echo "创建智能灯泡设备..."
curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"客厅灯泡\",\"product_id\":$BULB_ID,\"tenant_id\":1,\"device_key\":\"$BULB_KEY\"}" > /dev/null
echo "  ✅ 客厅灯泡设备创建成功"

curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"卧室灯泡\",\"product_id\":$BULB_ID,\"tenant_id\":1,\"device_key\":\"$BULB_KEY\"}" > /dev/null
echo "  ✅ 卧室灯泡设备创建成功"

curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"厨房灯泡\",\"product_id\":$BULB_ID,\"tenant_id\":1,\"device_key\":\"$BULB_KEY\"}" > /dev/null
echo "  ✅ 厨房灯泡设备创建成功"

# 创建插座设备
echo "创建智能插座设备..."
curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"电视插座\",\"product_id\":$SOCKET_ID,\"tenant_id\":1,\"device_key\":\"$SOCKET_KEY\"}" > /dev/null
echo "  ✅ 电视插座设备创建成功"

curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"空调插座\",\"product_id\":$SOCKET_ID,\"tenant_id\":1,\"device_key\":\"$SOCKET_KEY\"}" > /dev/null
echo "  ✅ 空调插座设备创建成功"

# 创建传感器设备
echo "创建温湿度传感器设备..."
curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"客厅传感器\",\"product_id\":$SENSOR_ID,\"tenant_id\":1,\"device_key\":\"$SENSOR_KEY\"}" > /dev/null
echo "  ✅ 客厅传感器设备创建成功"

curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d "{\"device_name\":\"卧室传感器\",\"product_id\":$SENSOR_ID,\"tenant_id\":1,\"device_key\":\"$SENSOR_KEY\"}" > /dev/null
echo "  ✅ 卧室传感器设备创建成功"

# 批量创建设备
echo "批量创建设备..."
for i in {1..10}; do
  curl -s -X POST http://localhost:8081/api/v1/devices \
    -H "Content-Type: application/json" \
    -d "{\"device_name\":\"演示设备_$i\",\"product_id\":$BULB_ID,\"tenant_id\":1,\"device_key\":\"$BULB_KEY\"}" > /dev/null
done
echo "  ✅ 批量创建10台演示设备"

echo ""

# ========== 3. 物模型演示 ==========
echo "【3. 物模型演示】"
echo ""
echo "已创建以下物模型属性:"
echo "  - 智能灯泡: 亮度(0-100), 色温(2700-6500K), 开关"
echo "  - 智能插座: 开关, 功率(瓦)"
echo "  - 温湿度传感器: 温度(℃), 湿度(%RH), 电量(%)"
echo "  ✅ 物模型配置完成"
echo ""

# ========== 4. 规则引擎演示 ==========
echo "【4. 规则引擎演示】"
echo ""

echo "创建温度超限告警规则..."
curl -s -X POST http://localhost:8083/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{"name":"温度超限告警","description":"当温度超过35度时发送告警","rule_type":"stream","sql_expression":"SELECT device_id, temperature FROM telemetry WHERE temperature > 35","status":"active"}' > /dev/null
echo "  ✅ 温度超限告警规则创建成功"

echo "创建设备上线通知规则..."
curl -s -X POST http://localhost:8083/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{"name":"设备上线通知","description":"设备上线时发送通知","rule_type":"stream","sql_expression":"SELECT device_id, device_name FROM devices WHERE online = true","status":"active"}' > /dev/null
echo "  ✅ 设备上线通知规则创建成功"

echo "创建低电量告警规则..."
curl -s -X POST http://localhost:8083/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{"name":"低电量告警","description":"电量低于20%时告警","rule_type":"stream","sql_expression":"SELECT device_id, battery FROM telemetry WHERE battery < 20","status":"active"}' > /dev/null
echo "  ✅ 低电量告警规则创建成功"

echo ""

# ========== 5. 场景联动演示 ==========
echo "【5. 场景联动演示】"
echo ""

echo "创建回家模式场景..."
curl -s -X POST http://localhost:8083/api/v1/scenes \
  -H "Content-Type: application/json" \
  -d '{"name":"回家模式","description":"自动打开客厅灯光和空调","trigger_rules":"manual","action_rules":"[{\"device_id\":1,\"action\":\"set_brightness\",\"value\":80},{\"device_id\":4,\"action\":\"turn_on\"}]","cron_expr":"0 18 * * *","status":"active"}' > /dev/null
echo "  ✅ 回家模式场景创建成功"

echo "创建离家模式场景..."
curl -s -X POST http://localhost:8083/api/v1/scenes \
  -H "Content-Type: application/json" \
  -d '{"name":"离家模式","description":"关闭所有灯光和电器","trigger_rules":"manual","action_rules":"[{\"device_id\":1,\"action\":\"turn_off\"},{\"device_id\":2,\"action\":\"turn_off\"},{\"device_id\":3,\"action\":\"turn_off\"},{\"device_id\":4,\"action\":\"turn_off\"},{\"device_id\":5,\"action\":\"turn_off\"}]","cron_expr":"0 19 * * *","status":"active"}' > /dev/null
echo "  ✅ 离家模式场景创建成功"

echo "创建睡眠模式场景..."
curl -s -X POST http://localhost:8083/api/v1/scenes \
  -H "Content-Type: application/json" \
  -d '{"name":"睡眠模式","description":"调暗灯光，关闭电器","trigger_rules":"manual","action_rules":"[{\"device_id\":1,\"action\":\"set_brightness\",\"value\":10},{\"device_id\":2,\"action\":\"set_brightness\",\"value\":5},{\"device_id\":3,\"action\":\"set_brightness\",\"value\":5},{\"device_id\":4,\"action\":\"turn_off\"},{\"device_id\":5,\"action\":\"turn_off\"}]","cron_expr":"0 22 * * *","status":"active"}' > /dev/null
echo "  ✅ 睡眠模式场景创建成功"

echo "创建起床模式场景..."
curl -s -X POST http://localhost:8083/api/v1/scenes \
  -H "Content-Type: application/json" \
  -d '{"name":"起床模式","description":"缓慢打开灯光，模拟日出","trigger_rules":"manual","action_rules":"[{\"device_id\":1,\"action\":\"set_brightness\",\"value\":30},{\"device_id\":2,\"action\":\"set_brightness\",\"value\":20},{\"device_id\":3,\"action\":\"set_brightness\",\"value\":20}]","cron_expr":"0 7 * * *","status":"active"}' > /dev/null
echo "  ✅ 起床模式场景创建成功"

echo ""

# ========== 6. 遥测数据演示 ==========
echo "【6. 遥测数据演示】"
echo ""

echo "上报遥测数据..."
for i in {1..20}; do
  # 灯泡设备数据
  curl -s -X POST http://localhost:8084/api/v1/telemetry \
    -H "Content-Type: application/json" \
    -d "{\"device_id\":$i,\"data\":{\"brightness\":$((i * 5)),\"color_temp\":$((3000 + i * 100)),\"power\":$((i % 2))}}" > /dev/null
  
  # 插座设备数据
  curl -s -X POST http://localhost:8084/api/v1/telemetry \
    -H "Content-Type: application/json" \
    -d "{\"device_id\":$(($i + 10)),\"data\":{\"power\":$((i % 2)),\"power_value\":$((i * 10))}}" > /dev/null
  
  # 传感器设备数据
  curl -s -X POST http://localhost:8084/api/v1/telemetry \
    -H "Content-Type: application/json" \
    -d "{\"device_id\":$(($i + 20)),\"data\":{\"temperature\":$((20 + i)),\"humidity\":$((50 + i)),\"battery\":$((100 - i))}}" > /dev/null
done
echo "  ✅ 已上报60条遥测数据"

echo ""

# ========== 7. OTA升级演示 ==========
echo "【7. OTA升级演示】"
echo ""

echo "创建固件版本..."
# 固件v1.0.0
curl -s -X POST http://localhost:8088/api/v1/ota/firmware \
  -H "Content-Type: application/json" \
  -d "{\"firmware_name\":\"固件v1.0.0\",\"product_id\":$BULB_ID,\"version\":\"1.0.0\",\"download_url\":\"https://example.com/firmware/v1.0.0.bin\",\"size\":1024000,\"status\":\"published\",\"release_notes\":\"初始版本\"}" > /dev/null
echo "  ✅ 固件v1.0.0创建成功"

# 固件v1.1.0
curl -s -X POST http://localhost:8088/api/v1/ota/firmware \
  -H "Content-Type: application/json" \
  -d "{\"firmware_name\":\"固件v1.1.0\",\"product_id\":$BULB_ID,\"version\":\"1.1.0\",\"download_url\":\"https://example.com/firmware/v1.1.0.bin\",\"size\":1048576,\"status\":\"published\",\"release_notes\":\"新增色温调节功能\"}" > /dev/null
echo "  ✅ 固件v1.1.0创建成功"

# 固件v1.2.0
curl -s -X POST http://localhost:8088/api/v1/ota/firmware \
  -H "Content-Type: application/json" \
  -d "{\"firmware_name\":\"固件v1.2.0\",\"product_id\":$SENSOR_ID,\"version\":\"1.2.0\",\"download_url\":\"https://example.com/firmware/v1.2.0.bin\",\"size\":524288,\"status\":\"published\",\"release_notes\":\"优化功耗\"}" > /dev/null
echo "  ✅ 固件v1.2.0创建成功"

echo ""
echo "创建升级任务..."
# 升级任务1
curl -s -X POST http://localhost:8088/api/v1/ota/tasks \
  -H "Content-Type: application/json" \
  -d '{"task_name":"全量升级v1.1.0","firmware_id":2,"target_percent":100,"status":"running","description":"所有设备升级到v1.1.0"}' > /dev/null
echo "  ✅ 升级任务1创建成功"

# 升级任务2
curl -s -X POST http://localhost:8088/api/v1/ota/tasks \
  -H "Content-Type: application/json" \
  -d '{"task_name":"传感器升级v1.2.0","firmware_id":3,"target_percent":50,"status":"pending","description":"传感器设备升级到v1.2.0"}' > /dev/null
echo "  ✅ 升级任务2创建成功"

echo ""

# ========== 8. 语音控制演示 ==========
echo "【8. 语音控制演示】"
echo ""

echo "绑定语音设备..."
curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":1,"voice_name":"客厅灯","platform":"tmall","status":"connected"}' > /dev/null
echo "  ✅ 绑定天猫精灵 - 客厅灯"

curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":2,"voice_name":"卧室灯","platform":"tmall","status":"connected"}' > /dev/null
echo "  ✅ 绑定天猫精灵 - 卧室灯"

curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":3,"voice_name":"厨房灯","platform":"tmall","status":"connected"}' > /dev/null
echo "  ✅ 绑定天猫精灵 - 厨房灯"

curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":4,"voice_name":"电视","platform":"tmall","status":"connected"}' > /dev/null
echo "  ✅ 绑定天猫精灵 - 电视"

curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":5,"voice_name":"空调","platform":"alexa","status":"connected"}' > /dev/null
echo "  ✅ 绑定Alexa - 空调"

curl -s -X POST http://localhost:8089/api/v1/voice/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":6,"voice_name":"客厅传感器","platform":"tmall","status":"connected"}' > /dev/null
echo "  ✅ 绑定天猫精灵 - 客厅传感器"

echo ""
echo "执行语音指令..."
curl -s -X POST http://localhost:8089/api/v1/voice/commands \
  -H "Content-Type: application/json" \
  -d '{"device_id":1,"command":"打开客厅灯","result":"success"}' > /dev/null
echo "  ✅ 执行指令: 打开客厅灯"

curl -s -X POST http://localhost:8089/api/v1/voice/commands \
  -H "Content-Type: application/json" \
  -d '{"device_id":2,"command":"关闭卧室灯","result":"success"}' > /dev/null
echo "  ✅ 执行指令: 关闭卧室灯"

curl -s -X POST http://localhost:8089/api/v1/voice/commands \
  -H "Content-Type: application/json" \
  -d '{"device_id":4,"command":"打开电视","result":"success"}' > /dev/null
echo "  ✅ 执行指令: 打开电视"

curl -s -X POST http://localhost:8089/api/v1/voice/commands \
  -H "Content-Type: application/json" \
  -d '{"device_id":5,"command":"打开空调","result":"success"}' > /dev/null
echo "  ✅ 执行指令: 打开空调"

curl -s -X POST http://localhost:8089/api/v1/voice/commands \
  -H "Content-Type: application/json" \
  -d '{"device_id":1,"command":"把客厅灯调暗一点","result":"success"}' > /dev/null
echo "  ✅ 执行指令: 把客厅灯调暗一点"

echo ""

# ========== 9. 设备地理位置演示 ==========
echo "【9. 设备地理位置演示】"
echo ""

echo "更新设备位置..."
curl -s -X POST http://localhost:8094/api/v1/devices/1/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":31.2304,"longitude":121.4737,"location":"上海市浦东新区","region":"华东"}' > /dev/null
echo "  ✅ 更新设备1位置: 上海市"

curl -s -X POST http://localhost:8094/api/v1/devices/2/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":31.2304,"longitude":121.4737,"location":"上海市浦东新区","region":"华东"}' > /dev/null
echo "  ✅ 更新设备2位置: 上海市"

curl -s -X POST http://localhost:8094/api/v1/devices/3/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":31.2304,"longitude":121.4737,"location":"上海市浦东新区","region":"华东"}' > /dev/null
echo "  ✅ 更新设备3位置: 上海市"

curl -s -X POST http://localhost:8094/api/v1/devices/4/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":23.1291,"longitude":113.2644,"location":"广州市天河区","region":"华南"}' > /dev/null
echo "  ✅ 更新设备4位置: 广州市"

curl -s -X POST http://localhost:8094/api/v1/devices/5/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":23.1291,"longitude":113.2644,"location":"广州市天河区","region":"华南"}' > /dev/null
echo "  ✅ 更新设备5位置: 广州市"

curl -s -X POST http://localhost:8094/api/v1/devices/6/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":22.5431,"longitude":114.0579,"location":"深圳市南山区","region":"华南"}' > /dev/null
echo "  ✅ 更新设备6位置: 深圳市"

echo ""

# ========== 10. 统计信息汇总 ==========
echo "【10. 统计信息汇总】"
echo ""

echo "=== 产品统计 ==="
curl -s http://localhost:8082/api/v1/products/statistics | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8082/api/v1/products/statistics
echo ""

echo "=== 设备统计 ==="
curl -s http://localhost:8081/api/v1/devices/statistics | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8081/api/v1/devices/statistics
echo ""

echo "=== 规则统计 ==="
curl -s http://localhost:8083/api/v1/rules/statistics | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8083/api/v1/rules/statistics
echo ""

echo "=== 运营大屏 ==="
curl -s http://localhost:8080/api/v1/ops/dashboard | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8080/api/v1/ops/dashboard
echo ""

echo "=== OTA统计 ==="
curl -s http://localhost:8088/api/v1/ota/statistics | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8088/api/v1/ota/statistics
echo ""

echo "=== 语音统计 ==="
curl -s http://localhost:8089/api/v1/voice/statistics | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8089/api/v1/voice/statistics
echo ""

echo "=========================================="
echo "演示数据创建完成！"
echo "=========================================="
echo ""
echo "数据统计:"
echo "  - 产品: 5个 (已发布)"
echo "  - 设备: 16台 (3灯泡 + 2插座 + 2传感器 + 10演示)"
echo "  - 规则: 3条"
echo "  - 场景: 4个"
echo "  - 固件: 3个版本"
echo "  - 升级任务: 2个"
echo "  - 语音设备: 6台"
echo "  - 语音指令: 5条"
echo "  - 地理位置: 6台设备"
echo ""
echo "访问地址:"
echo "  - API网关: http://localhost:8080"
echo "  - 前端页面: http://localhost:3000"
echo "  - EMQX Dashboard: http://localhost:18083"
