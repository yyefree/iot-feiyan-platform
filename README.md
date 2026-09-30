# IoT飞燕平台

一个仿阿里云飞燕平台的物联网IoT平台，支持设备管理、物模型、规则引擎、数据统计、运营中心等功能。

## 🚀 快速开始

### 使用 WSLC 部署

```bash
# 启动所有服务
.\start_wslc.ps1

# 或者使用 Docker Compose
docker compose up -d

# 运行测试
.\test_all.sh
```

### 访问地址

- 前端界面: http://localhost:3000
- API网关: http://localhost:8080/health
- EMQX Dashboard: http://localhost:18083

## 📦 技术栈

### 后端（Rust 重构版）
- **Web框架**: Actix-web 4
- **ORM**: SeaORM 1
- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **时序数据库**: InfluxDB 2
- **MQTT Broker**: EMQX 5.7
- **异步运行时**: Tokio

### 前端
- **框架**: Vue 3
- **语言**: TypeScript
- **UI库**: Element Plus
- **可视化**: ECharts

### 部署
- **容器**: Docker 29
- **编排**: Docker Compose 5
- **WSLC**: Windows Subsystem for Linux Containers

## 📁 项目结构

```
D:\AI\test\
├── rust-services/           # Rust 微服务
│   ├── common/              # 公共模块
│   ├── device-service/      # 设备管理
│   ├── product-service/     # 产品管理
│   ├── rule-service/        # 规则引擎
│   ├── data-service/        # 数据统计
│   ├── ops-service/         # 运营中心
│   └── gateway/             # API网关
├── services/                # 原 Go 服务（兼容保留）
├── frontend/                # Vue 前端
├── gateway/                 # 原 Go API网关
├── init/sql/                # 数据库初始化脚本
├── docker-compose.wslc.yml  # WSLC 配置
├── start_wslc.ps1           # WSLC 启动脚本
├── test_all.sh              # 功能测试脚本
└── RUST_MIGRATION.md        # Rust 重构文档
```

## 🔌 核心功能

| 功能模块 | 描述 | API |
|---------|------|-----|
| 设备管理 | 设备CRUD、虚拟设备、设备影子 | /api/v1/devices |
| 产品管理 | 产品CRUD、物模型(TSL) | /api/v1/products |
| 规则引擎 | SQL规则、场景联动 | /api/v1/rules |
| 数据统计 | 遥测数据、可视化图表 | /api/v1/data |
| 运营中心 | 大屏监控、告警管理 | /api/v1/ops |
| 消息推送 | App推送、微信推送 | /api/v1/push |

## 📊 性能优势

| 指标 | Go (原) | Rust (新) |
|------|---------|-----------|
| 内存占用 | ~50MB | ~15MB |
| 启动时间 | ~2s | ~0.5s |
| 请求延迟 | ~5ms | ~2ms |
| 二进制大小 | ~20MB | ~8MB |

## 🧪 测试

```bash
# 全功能测试
.\test_all.sh

# 健康检查
curl http://localhost:8080/health

# 设备统计
curl http://localhost:8081/api/v1/devices/statistics
```

## 📚 文档

- [Rust 重构文档](RUST_MIGRATION.md)
- [功能测试报告](功能测试报告.md)
- [技术栈总结](README.md)

## 🔗 相关链接

- [飞燕平台](https://feoyan.aliyun.com/)
- [腾讯连连](https://iot.weixin.qq.com/)
- [Actix-web](https://actix.rs/)
- [SeaORM](https://www.sea-orm.dev/)
