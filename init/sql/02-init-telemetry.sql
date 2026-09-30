-- =====================================================
-- 时序数据表（用于 InfluxDB 或 PostgreSQL 时序表）
-- =====================================================

-- 方法一: PostgreSQL 分区表（时序数据）
CREATE TABLE IF NOT EXISTS telemetry (
    id BIGSERIAL,
    device_id BIGINT NOT NULL REFERENCES device(id) ON DELETE CASCADE,
    tenant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    identify VARCHAR(128) NOT NULL,
    value_type VARCHAR(32) DEFAULT 'float',
    value_float DOUBLE PRECISION,
    value_int BIGINT,
    value_string TEXT,
    value_bool BOOLEAN,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

-- 创建月度分区
CREATE TABLE telemetry_2024_01 PARTITION OF telemetry
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE telemetry_2024_02 PARTITION OF telemetry
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
CREATE TABLE telemetry_2024_03 PARTITION OF telemetry
    FOR VALUES FROM ('2024-03-01') TO ('2024-04-01');

-- 创建索引
CREATE INDEX idx_telemetry_device_time ON telemetry(device_id, timestamp DESC);
CREATE INDEX idx_telemetry_timestamp ON telemetry(timestamp DESC);
CREATE INDEX idx_telemetry_tenant ON telemetry(tenant_id, timestamp DESC);

-- =====================================================
-- 设备在线状态表
-- =====================================================
CREATE TABLE IF NOT EXISTS device_online_status (
    device_id BIGINT PRIMARY KEY REFERENCES device(id) ON DELETE CASCADE,
    tenant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    online_at TIMESTAMPTZ,
    offline_at TIMESTAMPTZ,
    last_heartbeat_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- =====================================================
-- 设备事件记录表
-- =====================================================
CREATE TABLE IF NOT EXISTS device_event (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES device(id) ON DELETE CASCADE,
    tenant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    event_identify VARCHAR(128) NOT NULL,
    event_type VARCHAR(32) NOT NULL,
    event_params JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_event_device ON device_event(device_id, timestamp DESC);
CREATE INDEX idx_device_event_tenant ON device_event(tenant_id, timestamp DESC);
CREATE INDEX idx_device_event_type ON device_event(event_type, timestamp DESC);
