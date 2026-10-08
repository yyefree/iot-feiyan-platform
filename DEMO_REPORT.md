# IoT飞燕平台 - 演示数据报告

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
| EMQX | iot-emqx | 1883/18083 | ✅ Running |

---

## 演示数据统计

### 1. 产品管理 (Product Service - :8082)

**统计:**
- 草稿产品: 31个
- 已发布: 1个
- 总计: 32个

**演示产品:**
| 产品名称 | 产品ID | 通信类型 | 状态 | 物模型属性 |
|----------|--------|----------|------|------------|
| Smart Bulb | 1 | WiFi | 已发布 | 亮度、色温、开关 |
| Smart Socket | 新增 | WiFi | 已发布 | 开关、功率 |
| 温湿度传感器 | 新增 | ZigBee | 已发布 | 温度、湿度、电量 |

---

### 2. 设备管理 (Device Service - :8081)

**统计:**
- 总设备数: 97台
- 在线设备: 19台
- 离线设备: 78台

**演示设备:**
| 设备名称 | 产品类型 | 位置 | 状态 |
|----------|----------|------|------|
| 客厅灯泡 | 智能灯泡 | 上海市 | 在线 |
| 卧室灯泡 | 智能灯泡 | 上海市 | 在线 |
| 厨房灯泡 | 智能灯泡 | 上海市 | 在线 |
| 电视插座 | 智能插座 | 广州市 | 在线 |
| 空调插座 | 智能插座 | 广州市 | 在线 |
| 客厅传感器 | 温湿度传感器 | 深圳市 | 在线 |
| 卧室传感器 | 温湿度传感器 | 上海市 | 在线 |
| 演示设备_1~10 | 智能灯泡 | 多地 | 离线 |

---

### 3. 规则引擎 (Rule Service - :8083)

**统计:**
- 总规则数: 22条
- 活跃规则: 21条
- 场景规则: 16条
- 流规则: 6条
- 总执行次数: 8次

**演示规则:**
| 规则名称 | 规则类型 | SQL表达式 | 状态 |
|----------|----------|-----------|------|
| 温度超限告警 | 流规则 | SELECT device_id, temperature FROM telemetry WHERE temperature > 35 | ✅ 活跃 |
| 设备上线通知 | 流规则 | SELECT device_id, device_name FROM devices WHERE online = true | ✅ 活跃 |
| 低电量告警 | 流规则 | SELECT device_id, battery FROM telemetry WHERE battery < 20 | ✅ 活跃 |

---

### 4. 场景联动 (Rule Service - :8083)

**统计:**
- 总场景数: 16个
- 活跃场景: 16个

**演示场景:**
| 场景名称 | 触发方式 | Cron表达式 | 执行动作 |
|----------|----------|------------|----------|
| 回家模式 | 手动 | 0 18 * * * | 打开客厅灯、打开电视 |
| 离家模式 | 手动 | 0 19 * * * | 关闭所有灯光和电器 |
| 睡眠模式 | 手动 | 0 22 * * * | 调暗灯光、关闭电器 |
| 起床模式 | 手动 | 0 7 * * * | 缓慢打开灯光 |

---

### 5. 数据统计 (Data Service - :8084)

**遥测数据:**
- 已上报: 60条
- 存储: InfluxDB时序数据库

**数据类型:**
- 灯泡数据: 亮度、色温、开关状态
- 插座数据: 开关状态、功率
- 传感器数据: 温度(℃)、湿度(%RH)、电量(%)

---

### 6. 运营中心 (API Gateway - :8080)

**运营大屏数据:**
```json
{
  "deviceStats": {
    "offline": 78,
    "online": 19,
    "total": 97
  },
  "alertStats": {
    "pending": 0,
    "total": 0
  }
}
```

---

### 7. OTA升级 (OTA Service - :8088)

**统计:**
- 固件版本: 5个
- 升级任务: 9个
- 草稿: 5个
- 已完成: 0个
- 失败: 0个

**演示固件:**
| 固件名称 | 版本号 | 产品 | 大小 | 状态 |
|----------|--------|------|------|------|
| 固件v1.0.0 | 1.0.0 | 智能灯泡 | 1MB | 已发布 |
| 固件v1.1.0 | 1.1.0 | 智能灯泡 | 1MB | 已发布 |
| 固件v1.2.0 | 1.2.0 | 温湿度传感器 | 512KB | 已发布 |

**演示任务:**
- 全量升级v1.1.0: 目标100%，状态运行中
- 传感器升级v1.2.0: 目标50%，状态待执行

---

### 8. 语音控制 (Voice Service - :8089)

**统计:**
- 连接设备: 5台
- 成功指令: 15条
- 总指令: 15条

**演示语音设备:**
| 设备名 | 语音平台 | 状态 |
|--------|----------|------|
| 客厅灯 | 天猫精灵 | 已连接 |
| 卧室灯 | 天猫精灵 | 已连接 |
| 厨房灯 | 天猫精灵 | 已连接 |
| 电视 | 天猫精灵 | 已连接 |
| 空调 | Alexa | 已连接 |

**演示指令记录:**
- "打开客厅灯" ✅
- "关闭卧室灯" ✅
- "打开电视" ✅
- "打开空调" ✅
- "把客厅灯调暗一点" ✅

---

### 9. 设备地理位置 (Rule Extensions - :8094)

**设备分布:**
| 城市 | 设备数 | 区域 |
|------|--------|------|
| 上海市 | 3台 | 华东 |
| 广州市 | 2台 | 华南 |
| 深圳市 | 1台 | 华南 |

---

## 测试覆盖率

```
总测试数: 77
通过:    77 (100%)
失败:     0
跳过:     0
```

---

## 访问地址

| 服务 | 地址 |
|------|------|
| **API网关** | http://localhost:8080 |
| **前端页面** | http://localhost:3000 |
| **EMQX Dashboard** | http://localhost:18083 |
| **InfluxDB UI** | http://localhost:8086 |
| **GitHub仓库** | https://github.com/yyefree/iot-feiyan-platform |

---

## 技术栈

### 后端
- **Rust (Actix-web + SeaORM)**: 6个核心微服务
- **Go (Gin + GORM)**: 3个附加服务

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

## 项目亮点

1. **完整的功能覆盖** - 产品、设备、物模型、规则、场景、数据、OTA、语音、地理
2. **多语言架构** - Rust核心服务 + Go附加服务，发挥各自优势
3. **容器化部署** - 一键启动14个服务，适合演示和开发环境
4. **完善的测试** - 77个测试用例全部通过
5. **真实演示数据** - 覆盖所有功能模块的演示数据

---

*报告生成时间: 2026-10-08*
