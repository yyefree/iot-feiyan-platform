-- 修复数据库schema
ALTER TABLE thing_model_defines ALTER COLUMN enum_values TYPE TEXT;
ALTER TABLE thing_model_defines ALTER COLUMN input_params TYPE TEXT;
ALTER TABLE thing_model_defines ALTER COLUMN output_params TYPE TEXT;
ALTER TABLE rules ALTER COLUMN trigger_config TYPE TEXT;
ALTER TABLE rules ALTER COLUMN action_config TYPE TEXT;
ALTER TABLE scenes ALTER COLUMN trigger_rules TYPE TEXT;
ALTER TABLE scenes ALTER COLUMN action_rules TYPE TEXT;
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS env_trigger TEXT;

-- 创建telemetry表
CREATE TABLE IF NOT EXISTS telemetry (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL,
    tenant_id BIGINT,
    product_id BIGINT,
    identify VARCHAR(128),
    value_type VARCHAR(32),
    value_float DOUBLE PRECISION,
    value_int BIGINT,
    value_string TEXT,
    value_bool BOOLEAN,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_telemetry_device_id ON telemetry(device_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_tenant_id ON telemetry(tenant_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_product_id ON telemetry(product_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_identify ON telemetry(identify);
CREATE INDEX IF NOT EXISTS idx_telemetry_timestamp ON telemetry(timestamp);

-- 创建数据看板视图
CREATE OR REPLACE VIEW v_dashboard_overview AS
SELECT 
    COUNT(DISTINCT d.id) as total_devices,
    COUNT(DISTINCT CASE WHEN d.online = true THEN d.id END) as online_devices,
    COUNT(DISTINCT CASE WHEN d.online = false THEN d.id END) as offline_devices,
    COUNT(DISTINCT t.device_id) as devices_with_telemetry,
    COUNT(t.id) as total_data_points,
    COUNT(CASE WHEN t.timestamp >= NOW() - INTERVAL '24 hours' THEN 1 END) as last_24h_points,
    COUNT(DISTINCT t.identify) as total_identifiers
FROM devices d
LEFT JOIN telemetry t ON d.id = t.device_id;
