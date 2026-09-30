# 阿里云飞燕平台深度调研报告 - 完整功能复刻方案

## 一、平台定位与核心特色

### 1.1 平台定位
- **全称**：生活物联网平台（飞燕平台）
- **定位**：面向**消费级生活智能设备**的一站式开发平台
- **目标用户**：家电/家居厂商、方案商、硬件开发者
- **核心价值**：大幅降低"设备-云-App"开发成本，实现快速智能化

### 1.2 核心特色
| 特色 | 描述 |
|------|------|
| 亿级连接 | 支持亿级设备在线、百万级并发处理 |
| 快速配置 | 1分钟完成界面配置，5小时完成智能化，10天完成量产 |
| 全球化 | 覆盖200+国家和地区，设备智能选择最近数据节点 |
| 双App方案 | 公版App（免开发）+ 自有品牌App（灵活定制） |
| 语音打通 | 一键对接天猫精灵、Amazon Alexa、Google Home |
| 安全合规 | 芯片级安全、GDPR合规、SOC1/SOC2审计 |

---

## 二、功能模块详细拆解

### 2.1 项目管理

#### 功能列表
- [x] **创建项目**：输入项目名称，最多创建4个项目
- [x] **项目设置**：修改项目名称、删除项目
- [x] **成员管理**：添加成员、分配权限
- [x] **项目授权**：将项目授权给其他阿里云账号

#### 数据结构
```json
{
  "projectId": "string",
  "projectName": "string",
  "createTime": "timestamp",
  "members": [
    {
      "userId": "string",
      "role": "owner|admin|member",
      "permission": "all|partial"
    }
  ]
}
```

---

### 2.2 产品管理

#### 功能列表
- [x] **创建产品**：选择产品类型（直连设备/网关）、品类、通讯类型
- [x] **产品类型**：
  - 直连设备：设备直接连云
  - 网关设备：支持子设备接入
- [x] **通讯类型**：WiFi、蜂窝（4G/NB-IoT）、以太网
- [x] **开发方式**：标准功能/自定义功能
- [x] **产品列表**：查看产品状态（未发布/已发布/已下线）
- [x] **产品详情**：查看产品基本信息、物模型、设备数

#### 产品类型对照表
| 类型 | 说明 | 适用场景 |
|------|------|----------|
| 直连设备 | 设备直接连接云平台 | 智能灯泡、插座、传感器 |
| 网关设备 | 通过网关接入子设备 | 智能门锁、安防传感器 |

---

### 2.3 物模型定义（核心功能）

#### 物模型三要素
1. **属性（Property）**：描述设备状态
2. **事件（Event）**：设备主动上报的信息
3. **服务（Service）**：云端可调用的指令

#### 属性类型
| 类型 | 标识符 | 示例 | 取值范围 |
|------|--------|------|----------|
| 布尔型 | Bool | isOn | true/false |
| 整型 | Int | brightness | 0-255 |
| 浮点型 | Float | temperature | -50.0~100.0 |
| 字符串 | String | mode | "auto"/"manual" |
| 枚举型 | Enum | color | ["red","green","blue"] |
| 结构体 | Struct | rgb | {r:0-255, g:0-255, b:0-255} |

#### 物模型JSON格式
```json
{
  "Identify": "temperature",
  "Name": "温度",
  "Required": true,
  "AccessPolicy": "rw",
  "Desc": "当前环境温度",
  "Type": "Float",
  "Unit": "℃",
  "UnitSymbol": "℃",
  "Min": -50.0,
  "Max": 100.0,
  "Step": 0.1
}
```

#### 服务定义
```json
{
  "Identify": "setBrightness",
  "Name": "设置亮度",
  "Desc": "设置灯具亮度",
  "Method": "thing.service.property.set",
  "Params": {
    "brightness": {
      "Type": "Int",
      "Min": 0,
      "Max": 255
    }
  }
}
```

#### 事件定义
```json
{
  "Identify": "alarmEvent",
  "Name": "报警事件",
  "Desc": "设备发生报警时上报",
  "Method": "thing.event.alarm.post",
  "Params": {
    "alarmType": {
      "Type": "String"
    },
    "alarmTime": {
      "Type": "Long"
    }
  }
}
```

---

### 2.4 设备管理

#### 设备生命周期
```
创建产品 → 添加设备 → 烧录证书 → 设备上线 → 数据通信 → 设备下线 → 删除设备
```

#### 功能列表
- [x] **添加设备**：手动添加、批量导入、激活码激活
- [x] **设备列表**：搜索、筛选、批量操作
- [x] **设备详情**：查看设备状态、属性值、在线状态
- [x] **设备调试**：实时属性读写、服务调用、事件订阅
- [x] **设备分组**：按房间、区域分组管理
- [x] **设备分享**：分享给其他用户

#### 设备认证方式
| 认证方式 | 说明 | 安全性 | 适用场景 |
|----------|------|--------|----------|
| 一机一密 | 每台设备唯一证书 | 高 | 高安全要求 |
| 一型一密 | 同型号设备共享证书 | 中 | 低成本量产 |
| 动态注册 | 设备首次联网时注册 | 中 | 灵活接入 |

#### 设备证书（三元组）
```json
{
  "ProductKey": "a1XXXXXXXX",
  "DeviceName": "device_001",
  "DeviceSecret": "xxxxxxxxxxxxxxxxxxxxxxxx"
}
```

---

### 2.5 设备配网

#### 配网方案
| 配网方式 | 技术原理 | 适用场景 |
|----------|----------|----------|
| 蓝牙辅助配网 | BLE传输WiFi配置 | 智能音箱、网关 |
| AP配网 | 设备开启热点，App连接配置 | 智能家居设备 |
| 一键广播配网 | 手机热点广播配网信息 | 路由器、中继设备 |
| 智能路由器配网 | 通过路由器DHCP选项传输 | 支持路由器的设备 |
| 设备间相互配网 | 已配网设备传输给新设备 | 网关、中控屏 |

#### 配网流程
```
1. App打开配网页面
2. 设备进入配网模式（长按配网键或自动）
3. App扫描发现设备（蓝牙/广播）
4. App输入WiFi账号密码
5. 设备接收配网信息
6. 设备连接WiFi并连云
7. 云端验证设备
8. 配网成功
```

---

### 2.6 App配置化

#### 公版App（云智能App）
- [x] **15分钟拖拽搭建面板**：可视化配置设备控制面板
- [x] **标准面板模板**：提供海量免费模板
- [x] **控件库**：开关、滑块、图表、地图、视频等
- [x] **多语言支持**：自动适配设备所在区域语言
- [x] **告警消息配置**：自定义推送消息内容

#### 面板配置项
| 控件类型 | 功能 | 配置项 |
|----------|------|--------|
| 开关控件 | 控制设备开关 | 绑定属性、显示文本 |
| 滑块控件 | 调节亮度/音量 | 绑定属性、范围、步长 |
| 图表控件 | 显示历史数据 | 绑定属性、时间范围 |
| 地图控件 | 显示设备位置 | 绑定经纬度属性 |
| 视频控件 | 实时视频预览 | 绑定视频流地址 |

#### 自定义面板
- [x] 支持上传自定义面板（HTML/JS）
- [x] 提供面板开发SDK
- [x] 本地预览与云端预览

---

### 2.7 自有品牌App

#### 功能列表
- [x] **iOS SDK**：Objective-C/Swift集成
- [x] **Android SDK**：Java/Kotlin集成
- [x] **配网SDK**：蓝牙辅助、AP、一键配网
- [x] **设备控制SDK**：属性读写、服务调用
- [x] **账号系统**：注册、登录、绑定设备
- [x] **消息推送**：告警、状态变更通知
- [x] **OTA升级**：固件远程升级

#### SDK模块
| 模块 | 功能 |
|------|------|
| 基础SDK | 账号、设备管理、数据通信 |
| 配网SDK | 多种配网方式 |
| 安全SDK | 设备认证、数据加密 |
| 推送SDK | 消息推送 |
| 语音SDK | 语音控制对接 |

---

### 2.8 场景联动

#### 功能列表
- [x] **创建场景**：设置触发条件、执行动作
- [x] **触发条件**：
  - 设备属性变化（温度>30℃）
  - 设备状态变化（开关开）
  - 时间触发（每天8:00）
  - 地理围栏（回家/离家）
  - 日出日落
- [x] **执行动作**：
  - 设备控制（开灯、调亮度）
  - 发送通知
  - 触发其他场景

#### 场景联动规则示例
```json
{
  "ruleName": "回家模式",
  "trigger": {
    "type": "geofence",
    "condition": "enter",
    "location": {
      "latitude": 30.25,
      "longitude": 120.15,
      "radius": 100
    }
  },
  "actions": [
    {
      "type": "device_control",
      "target": "light_001",
      "action": "turn_on",
      "params": {"brightness": 80}
    },
    {
      "type": "device_control",
      "target": "ac_001",
      "action": "set_mode",
      "params": {"mode": "cool", "temp": 26}
    }
  ]
}
```

---

### 2.9 规则引擎

#### 功能列表
- [x] **创建规则**：可视化规则编辑器
- [x] **数据流转**：设备数据→规则→目标
- [x] **目标类型**：
  - 设备属性设置
  - 消息推送（MQTT/HTTP）
  - 数据库写入
  - 第三方服务调用

#### 规则语法
```
FROM /topic/sys/${productKey}/${deviceName}/thing/property/post
WHERE temperature > 30
THEN action setProperty(${deviceName}, "alarm", true)
```

---

### 2.10 语音控制

#### 对接语音平台
- [x] **天猫精灵**：国内语音助手
- [x] **Amazon Alexa**：海外语音助手
- [x] **Google Assistant**：海外语音助手
- [x] **小爱同学**：小米语音助手

#### 对接流程
```
1. 在飞燕平台绑定语音平台技能
2. 配置物模型与语音指令映射
3. 用户通过语音App发现设备
4. 语音控制设备
```

---

### 2.11 数据管理

#### 功能列表
- [x] **设备属性存储**：时序数据存储
- [x] **数据查询**：按时间范围、设备查询
- [x] **数据导出**：CSV/Excel导出
- [x] **数据保留**：默认90天，可配置
- [x] **云对云同步**：同步到客户服务器

#### 数据存储格式
```json
{
  "ts": 1632456789000,
  "id": "a1XXXXXXXX.device_001",
  "v": {
    "temperature": 25.5,
    "humidity": 60.0
  }
}
```

---

### 2.12 运营中心

#### 功能列表
- [x] **数据概览**：设备激活数、在线数、活跃数
- [x] **设备列表**：设备状态、分组管理
- [x] **用户管理**：用户数、活跃度
- [x] **数据大屏**：地理分布、趋势图表
- [x] **统计分析**：多维度数据统计

#### 统计指标
| 指标 | 说明 |
|------|------|
| 设备激活数 | 累计激活设备数量 |
| 设备在线数 | 当前在线设备数量 |
| 设备活跃数 | 当日有数据上报的设备数 |
| 消息总量 | 当日消息收发总量 |
| 用户激活数 | 累计绑定设备的用户数 |

---

### 2.13 量产中心

#### 功能列表
- [x] **激活码管理**：购买、申请激活码
- [x] **批量生产**：批量生成设备证书
- [x] **证书下载**：下载设备证书文件
- [x] **烧录支持**：支持多种烧录方式

#### 量产流程
```
1. 发布产品
2. 购买激活码
3. 批量生成设备
4. 下载证书
5. 烧录到设备
6. 设备上线
```

---

### 2.14 OTA升级

#### 功能列表
- [x] **固件上传**：上传升级包
- [x] **版本管理**：版本发布、回滚
- [x] **升级策略**：全量/分批/灰度升级
- [x] **升级监控**：升级进度、成功率统计

---

### 2.15 安全体系

#### 功能列表
- [x] **设备认证**：一机一密、一型一密、动态注册
- [x] **数据加密**：TLS/DTLS传输加密
- [x] **数据隔离**：租户级数据隔离
- [x] **访问控制**：RBAC权限管理
- [x] **操作审计**：日志记录与查询
- [x] **隐私合规**：GDPR、SOC1/SOC2认证

---

### 2.16 全球化部署

#### 功能列表
- [x] **全球节点**：支持200+国家和地区
- [x] **就近接入**：设备智能选择最近节点
- [x] **多语言**：App界面多语言支持
- [x] **海外合规**：GDPR合规、数据本地化

---

### 2.17 开放API

#### API分类
| 类别 | 接口示例 |
|------|----------|
| 产品管理 | CreateProduct, UpdateProduct, DeleteProduct |
| 设备管理 | AddDevice, QueryDeviceList, UpdateDevice |
| 物模型 | GetThingModel, UpdateThingModel |
| 数据查询 | QueryPropertyValue, QueryPropertyHistory |
| 规则引擎 | CreateRule, UpdateRule, DeleteRule |
| 统计查询 | GetDeviceStat, GetOnlineStat |

---

## 三、复刻技术方案

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                         用户交互层                                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │
│  │ Web控制台 │  │ 移动APP  │  │ 小程序   │  │ 语音助手 │           │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘           │
└─────────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                         API网关层                                   │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  Kong / Nginx + API Gateway (路由、限流、鉴权、日志)         │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                         业务服务层 (微服务)                          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │项目管理 │ │产品服务 │ │设备服务 │ │物模型服务│ │App服务  │     │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │规则引擎 │ │场景联动 │ │数据统计 │ │运营服务 │ │安全服务 │     │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │配网服务 │ │OTA服务  │ │语音服务 │ │量产服务 │ │推送服务 │     │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
└─────────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                         设备接入层                                  │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  MQTT Broker集群 (EMQX 5.x)                                │   │
│  │  支持: MQTT 3.1.1/5.0, WebSocket, CoAP, HTTP               │   │
│  │  设备认证: 动态鉴权、TLS/DTLS、设备证书                       │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                         数据存储层                                  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐             │
│  │ PostgreSQL│ │ Redis    │ │ InfluxDB │ │ Elasticsearch│           │
│  │ (元数据) │ │ (缓存/会话)│ │ (时序数据)│ │ (日志检索) │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘             │
│  ┌──────────┐ ┌──────────┐                                        │
│  │ Kafka    │ │ MinIO    │                                        │
│  │ (消息队列)│ │ (文件存储)│                                        │
│  └──────────┘ └──────────┘                                        │
└─────────────────────────────────────────────────────────────────────┘
```

### 3.2 核心功能复刻对照表

| 飞燕平台功能 | 复刻方案 | 技术实现 |
|-------------|---------|---------|
| 项目管理 | 多租户项目管理系统 | PostgreSQL + RBAC |
| 物模型定义 | JSON Schema编辑器 | MongoDB存储物模型定义 |
| 设备管理 | 设备生命周期管理 | MQTT + 状态机 |
| 设备配网 | 多模式配网服务 | BLE + AP + 广播 |
| App配置化 | 低代码面板编辑器 | Vue3 + 拖拽组件库 |
| 自有App SDK | 移动端SDK | Flutter/Android/iOS |
| 场景联动 | 可视化规则编排 | Drools/自研DSL |
| 规则引擎 | 数据流转引擎 | Kafka + 规则匹配 |
| 语音控制 | 多平台对接 | 技能绑定+映射配置 |
| 数据统计 | 实时+离线分析 | InfluxDB + ES |
| 运营中心 | 数据看板 | ECharts大屏 |
| 量产中心 | 批量证书管理 | 证书生成+导出 |
| OTA升级 | 固件分发管理 | MinIO + 版本控制 |
| 安全体系 | 设备认证+数据加密 | JWT + TLS + 审计 |
| 全球化 | 多区域部署 | K8s多集群+CDN |

---

## 四、开发计划

### Phase 1: 基础框架（2周）
- [ ] 项目脚手架搭建
- [ ] 数据库设计
- [ ] 微服务框架集成
- [ ] EMQX部署

### Phase 2: 核心功能（4周）
- [ ] 产品管理模块
- [ ] 物模型定义与管理
- [ ] 设备接入与管理
- [ ] 设备认证系统

### Phase 3: App与交互（3周）
- [ ] Web控制台开发
- [ ] App配置化编辑器
- [ ] 自有App SDK基础版

### Phase 4: 高级功能（3周）
- [ ] 场景联动与规则引擎
- [ ] 数据统计与运营中心
- [ ] OTA升级服务
- [ ] 语音控制对接

### Phase 5: 优化与部署（2周）
- [ ] 性能优化
- [ ] 安全加固
- [ ] 容器化部署
- [ ] 测试与文档

---

## 五、技术栈汇总

| 层级 | 技术选型 |
|------|---------|
| **设备接入** | EMQX 5.x (MQTT Broker) |
| **后端框架** | Spring Boot 3.x + Spring Cloud |
| **数据库** | PostgreSQL 16 + MongoDB 7.x |
| **时序存储** | InfluxDB 2.x |
| **搜索引擎** | Elasticsearch 8.x |
| **缓存** | Redis 7.x Cluster |
| **消息队列** | Kafka 3.x |
| **对象存储** | MinIO |
| **容器部署** | Kubernetes 1.29+ |
| **前端框架** | Vue 3 + TypeScript + Vite |
| **UI组件库** | Element Plus |
| **图表库** | ECharts 5.x |
| **移动端** | Flutter 3.x |
| **API网关** | Spring Cloud Gateway |
| **服务治理** | Nacos 2.x |
| **监控运维** | Prometheus + Grafana + ELK |

---

## 六、关键数据模型

### 6.1 产品表 (product)
```sql
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    product_key VARCHAR(64) UNIQUE NOT NULL,
    product_name VARCHAR(128) NOT NULL,
    category_id INTEGER REFERENCES category(id),
    product_type VARCHAR(32) NOT NULL,  -- direct_device/gateway
    comm_type VARCHAR(32) NOT NULL,     -- wifi/cellular/ethernet
    auth_type VARCHAR(32) NOT NULL,     -- one_device_one_key/one_model_one_key
    thing_model JSONB,
    status VARCHAR(32) DEFAULT 'draft', -- draft/published/offline
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 6.2 设备表 (device)
```sql
CREATE TABLE device (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES product(id),
    device_name VARCHAR(128) NOT NULL,
    device_secret VARCHAR(256),
    status VARCHAR(32) DEFAULT 'offline', -- offline/online
    last_online_at TIMESTAMP,
    group_id INTEGER REFERENCES device_group(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, device_name)
);
```

### 6.3 物模型表 (thing_model)
```sql
CREATE TABLE thing_model (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES product(id),
    define_type VARCHAR(32) NOT NULL,  -- property/service/event
    identify VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    desc TEXT,
    access_policy VARCHAR(32),         -- r/rw
    type VARCHAR(32) NOT NULL,         -- bool/int/float/string/enum/struct
    extra JSONB,
    UNIQUE(product_id, identify)
);
```

### 6.4 设备属性表 (device_property)
```sql
CREATE TABLE device_property (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT REFERENCES device(id),
    identify VARCHAR(128) NOT NULL,
    value TEXT,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id, identify)
);
```

### 6.5 时序数据表 (telemetry)
```sql
CREATE TABLE telemetry (
    time TIMESTAMPTZ NOT NULL,
    device_id BIGINT REFERENCES device(id),
    identify VARCHAR(128) NOT NULL,
    value DOUBLE PRECISION,
    PRIMARY KEY (time, device_id, identify)
);
```

---

## 七、API接口设计

### 7.1 设备接入API
```
POST /api/v1/device/connect       # 设备连接
POST /api/v1/device/online        # 设备上线
POST /api/v1/device/offline       # 设备下线
POST /api/v1/device/property/set  # 设置属性
POST /api/v1/device/property/get  # 获取属性
POST /api/v1/device/event/post    # 上报事件
POST /api/v1/device/service/call  # 调用服务
```

### 7.2 设备管理API
```
GET    /api/v1/devices           # 设备列表
POST   /api/v1/devices           # 创建设备
GET    /api/v1/devices/{id}      # 设备详情
PUT    /api/v1/devices/{id}      # 更新设备
DELETE /api/v1/devices/{id}      # 删除设备
POST   /api/v1/devices/{id}/share # 分享设备
```

### 7.3 产品管理API
```
GET    /api/v1/products          # 产品列表
POST   /api/v1/products          # 创建产品
GET    /api/v1/products/{id}     # 产品详情
PUT    /api/v1/products/{id}     # 更新产品
DELETE /api/v1/products/{id}     # 删除产品
POST   /api/v1/products/{id}/publish # 发布产品
```

### 7.4 数据统计API
```
GET    /api/v1/statistics/overview   # 数据概览
GET    /api/v1/statistics/devices    # 设备统计
GET    /api/v1/statistics/online     # 在线统计
GET    /api/v1/statistics/geography  # 地理分布
GET    /api/v1/statistics/history    # 历史数据
```

---

*调研完成时间: 2026-09-22*
*复刻方案版本: v1.0*
