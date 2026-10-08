# IoT飞燕平台 - 演示报告

## 系统状态

### 容器服务状态 (14/14 运行中)

| 服务 | 容器名 | 端口 | 状态 |
|------|--------|------|------|
| API Gateway | iot-api-gateway | 8080 | ✅ Running |
| Device Service | iot-device-service | 8081 | ✅ Running |
| Product Service | iot-product-service | 8082 | ✅ Running |
| Rule Service | iot-rule-service | 8083 | ✅ Running |
| Data Service | iot-data-service | 8084 | ✅ Running |
| Ops Service | iot-ops-service | 8085 | ✅ Running |
| OTA Service | iot-ota-service | 8088 | ✅ Running |
| Voice Service | iot-voice-service | 8089 | ✅ Running |
| Rule Extensions | iot-rule-extensions | 8094 | ✅ Running |
| Web Frontend | iot-web-frontend | 3000 | ✅ Running |
| PostgreSQL | iot-postgres | 5432 | ✅ Healthy |
| Redis | iot-redis | 6379 | ✅ Healthy |
| InfluxDB | iot-influxdb | 8086 | ✅ Healthy |
| EMQX | iot-emqx | 1883/18083 | ⚠️ Unhealthy |

---

## 功能演示

### 1. 产品管理 (Product Service - :8082)

**已有产品数据:**
- **Smart Bulb** (ID: 1) - 已发布，带TSL物模型
  - 产品密钥: p1790146096257
  - 通信类型: WiFi
  - 认证方式: one_device_one_key
  - TSL属性: test_prop (测试属性)

- **Smart Switch** (ID: 2) - 草稿状态
- **Test Product** (ID: 3-25) - 多个测试产品

**API端点:**
- `GET /api/v1/products` - 产品列表
- `GET /api/v1/products/statistics` - 产品统计
  - 草稿: 24个
  - 已发布: 1个
  - 总计: 25个

---

### 2. 设备管理 (Device Service - :8081)

**设备统计:**
- 总设备数: 92台
- 在线设备: 18台
- 离线设备: 74台

**设备类型:**
- 真实设备: Test Bulb, Test Switch, Test Device等
- 虚拟设备: 虚拟设备_1791421094, 虚拟设备_1791421932等
- 批量设备: batch_0001, batch_0002, batch_0003等

**API端点:**
- `GET /api/v1/devices` - 设备列表
- `GET /api/v1/devices/statistics` - 设备统计

---

### 3. 规则引擎 (Rule Service - :8083)

**规则统计:**
- 总规则数: 15条
- 活跃规则: 14条
- 场景规则: 15条
- 流规则: 0条
- 总执行次数: 7次

**规则示例:**
- SQL表达式: `SELECT device_id FROM devices WHERE online = true`
- 触发类型: scene (场景触发)

**API端点:**
- `GET /api/v1/rules` - 规则列表
- `GET /api/v1/rules/statistics` - 规则统计
- `POST /api/v1/rules/sql/test` - SQL测试

---

### 4. 场景联动 (Rule Service - :8083)

**场景统计:**
- 总场景数: 15个
- 所有场景状态: active (活跃)

**场景示例:**
- 场景名称: 测试场景_1791422825, 测试场景_1791423646等
- Cron表达式: `0 18 * * *` (每天18:00触发)
- 触发类型: 定时触发

**API端点:**
- `GET /api/v1/scenes` - 场景列表
- `GET /api/v1/scenes/sunrise` - 日出日落时间

---

### 5. 数据统计 (Data Service - :8084)

**遥测数据:**
- 已上报多条遥测数据
- 数据存储: InfluxDB时序数据库

**API端点:**
- `POST /api/v1/telemetry` - 上报遥测数据
- `GET /api/v1/telemetry/latest/:device_id` - 最新数据
- `GET /api/v1/telemetry/query` - 查询数据

---

### 6. 运营中心 (API Gateway - :8080)

**运营大屏数据:**
```json
{
  "deviceStats": {
    "offline": 74,
    "online": 18,
    "total": 92
  },
  "alertStats": {
    "pending": 0,
    "total": 0
  },
  "dataStats": {
    "last24h": 0,
    "total": 0
  }
}
```

**API端点:**
- `GET /api/v1/ops/dashboard` - 运营大屏

---

### 7. OTA升级 (OTA Service - :8088)

**固件列表:**
- 共4个固件版本
- 固件名称: 固件测试
- 版本号: 1.0.0
- 状态: draft (草稿)

**升级任务:**
- 共4个升级任务
- 任务名称: 测试任务
- 目标比例: 100%
- 状态: pending (待执行)

**API端点:**
- `GET /api/v1/ota/firmware` - 固件列表
- `GET /api/v1/ota/tasks` - 升级任务列表
- `GET /api/v1/ota/statistics` - OTA统计
  - 活动任务: 0
  - 已完成: 0
  - 草稿: 4
  - 失败: 0
  - 固件总数: 4
  - 任务总数: 4

---

### 8. 语音控制 (Voice Service - :8089)

**语音设备:**
- 共4个绑定设备
- 供应商: 天猫精灵
- 状态: connected (已连接)

**语音指令记录:**
- 指令: "打开灯"
- 结果: success
- 共4条指令记录

**API端点:**
- `GET /api/v1/voice/devices` - 语音设备列表
- `GET /api/v1/voice/commands/logs` - 指令日志
- `GET /api/v1/voice/statistics` - 语音统计
  - 连接设备: 4台
  - 成功指令: 4条
  - 总指令: 4条

---

### 9. 设备地理位置 (Rule Extensions - :8094)

**设备位置数据:**
- 上海设备: 纬度31.2304, 经度121.4737 (华东)
- 广州设备: 纬度23.1291, 经度113.2644 (华南)
- 深圳设备: 纬度22.5431, 经度114.0579 (华南)

**API端点:**
- `GET /api/v1/devices/locations` - 设备位置列表
- `POST /api/v1/devices/:id/location` - 更新设备位置

---

## 访问地址

| 服务 | 地址 |
|------|------|
| API网关 | http://localhost:8080 |
| 前端页面 | http://localhost:3000 |
| EMQX Dashboard | http://localhost:18083 |
| InfluxDB UI | http://localhost:8086 |

---

## 测试覆盖率

```
总测试数: 77
通过:    77 (100%)
失败:     0
跳过:     0
```

---

## 技术栈总结

### 后端
- **Rust (Actix-web + SeaORM)**: 6个核心微服务
- **Go (Gin + GORM)**: 3个附加服务 (OTA, Voice, Rule Extensions)

### 基础设施
- **PostgreSQL 16**: 关系型数据库
- **Redis 7**: 缓存层
- **InfluxDB 2.x**: 时序数据库
- **EMQX 5.7**: MQTT Broker

### 前端
- **Vue 3 + Element Plus + Vite**
- **Nginx**: 静态文件服务

### 部署
- **Docker Compose**: 容器编排
- **WSLC**: Windows Subsystem for Linux Containers

---

## 项目仓库

- **GitHub**: https://github.com/yyefree/iot-feiyan-platform
- **分支**: master
- **最新提交**: 75e6940 IoT飞燕平台 - 容器化完整版 (77/77测试通过)
