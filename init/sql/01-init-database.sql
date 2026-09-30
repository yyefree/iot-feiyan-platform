-- =====================================================
-- 物联网平台数据库初始化脚本
-- 数据库: iot_platform
-- =====================================================

-- 启用扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- =====================================================
-- 1. 租户表
-- =====================================================
CREATE TABLE tenant (
    id BIGSERIAL PRIMARY KEY,
    tenant_code VARCHAR(64) UNIQUE NOT NULL,
    tenant_name VARCHAR(128) NOT NULL,
    contact_name VARCHAR(64),
    contact_phone VARCHAR(32),
    contact_email VARCHAR(128),
    address VARCHAR(256),
    status VARCHAR(32) DEFAULT 'active',
    expire_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 2. 用户表
-- =====================================================
CREATE TABLE "user" (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(256) NOT NULL,
    nickname VARCHAR(64),
    email VARCHAR(128),
    phone VARCHAR(32),
    avatar VARCHAR(512),
    role VARCHAR(32) DEFAULT 'user',
    status VARCHAR(32) DEFAULT 'active',
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, username)
);

-- =====================================================
-- 3. 项目表
-- =====================================================
CREATE TABLE project (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    project_key VARCHAR(64) UNIQUE NOT NULL,
    project_name VARCHAR(128) NOT NULL,
    description TEXT,
    status VARCHAR(32) DEFAULT 'active',
    region VARCHAR(64) DEFAULT 'cn-shanghai',
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 4. 项目成员表
-- =====================================================
CREATE TABLE project_member (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    role VARCHAR(32) DEFAULT 'member',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(project_id, user_id)
);

-- =====================================================
-- 5. 品类表
-- =====================================================
CREATE TABLE category (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT REFERENCES category(id),
    name VARCHAR(128) NOT NULL,
    name_en VARCHAR(128),
    level INTEGER DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    status VARCHAR(32) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 6. 产品表
-- =====================================================
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    project_id BIGINT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    product_key VARCHAR(64) UNIQUE NOT NULL,
    product_name VARCHAR(128) NOT NULL,
    category_id BIGINT REFERENCES category(id),
    product_type VARCHAR(32) DEFAULT 'direct_device',
    comm_type VARCHAR(32) DEFAULT 'wifi',
    auth_type VARCHAR(32) DEFAULT 'one_device_one_key',
    desc TEXT,
    icon VARCHAR(512),
    thing_model JSONB,
    status VARCHAR(32) DEFAULT 'draft',
    version VARCHAR(32),
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 7. 设备表
-- =====================================================
CREATE TABLE device (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    device_name VARCHAR(128) NOT NULL,
    device_secret VARCHAR(256),
    device_key VARCHAR(256),
    status VARCHAR(32) DEFAULT 'offline',
    online BOOLEAN DEFAULT FALSE,
    last_online_at TIMESTAMP,
    last_offline_at TIMESTAMP,
    ip_address VARCHAR(64),
    firmware_version VARCHAR(64),
    group_id BIGINT REFERENCES device_group(id),
    region VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, device_name)
);

-- =====================================================
-- 8. 设备分组表
-- =====================================================
CREATE TABLE device_group (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    parent_id BIGINT REFERENCES device_group(id),
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 9. 物模型定义表
-- =====================================================
CREATE TABLE thing_model_define (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    define_type VARCHAR(32) NOT NULL,
    identify VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    desc TEXT,
    access_policy VARCHAR(32),
    type VARCHAR(32) NOT NULL,
    unit VARCHAR(32),
    unit_symbol VARCHAR(32),
    min DECIMAL(20, 6),
    max DECIMAL(20, 6),
    step DECIMAL(20, 6),
    enum_values JSONB,
    struct_items JSONB,
    extra JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, identify)
);

-- =====================================================
-- 10. 设备属性表
-- =====================================================
CREATE TABLE device_property (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES device(id) ON DELETE CASCADE,
    identify VARCHAR(128) NOT NULL,
    value TEXT,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id, identify)
);

-- =====================================================
-- 11. 设备影子表
-- =====================================================
CREATE TABLE device_shadow (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES device(id) ON DELETE CASCADE,
    desired_state JSONB,
    reported_state JSONB,
    reported_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id)
);

-- =====================================================
-- 12. 设备规则表
-- =====================================================
CREATE TABLE device_rule (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    desc TEXT,
    rule_type VARCHAR(32) DEFAULT 'scene',
    trigger_config JSONB,
    action_config JSONB,
    status VARCHAR(32) DEFAULT 'active',
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 13. 场景联动表
-- =====================================================
CREATE TABLE scene_linkage (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    project_id BIGINT REFERENCES project(id),
    name VARCHAR(128) NOT NULL,
    desc TEXT,
    trigger_rules JSONB,
    action_rules JSONB,
    status VARCHAR(32) DEFAULT 'active',
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 14. 设备分享表
-- =====================================================
CREATE TABLE device_share (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES device(id) ON DELETE CASCADE,
    sharer_id BIGINT NOT NULL REFERENCES "user"(id),
    sharee_id BIGINT NOT NULL REFERENCES "user"(id),
    permission VARCHAR(32) DEFAULT 'read',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_id, sharee_id)
);

-- =====================================================
-- 15. 消息推送表
-- =====================================================
CREATE TABLE message_push (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    user_id BIGINT,
    device_id BIGINT,
    msg_type VARCHAR(32) NOT NULL,
    title VARCHAR(256),
    content TEXT,
    params JSONB,
    status VARCHAR(32) DEFAULT 'pending',
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 16. OTA固件表
-- =====================================================
CREATE TABLE ota_firmware (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    product_id BIGINT REFERENCES product(id),
    firmware_name VARCHAR(128) NOT NULL,
    version VARCHAR(64) NOT NULL,
    file_url VARCHAR(512) NOT NULL,
    file_size BIGINT,
    checksum VARCHAR(128),
    release_notes TEXT,
    status VARCHAR(32) DEFAULT 'draft',
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 17. 设备升级任务表
-- =====================================================
CREATE TABLE ota_task (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    firmware_id BIGINT NOT NULL REFERENCES ota_firmware(id),
    task_name VARCHAR(128) NOT NULL,
    target_devices JSONB,
    strategy VARCHAR(32) DEFAULT 'full',
    progress JSONB,
    status VARCHAR(32) DEFAULT 'pending',
    created_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 18. 激活码表
-- =====================================================
CREATE TABLE activation_code (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    code VARCHAR(64) UNIQUE NOT NULL,
    device_name VARCHAR(128),
    status VARCHAR(32) DEFAULT 'unused',
    used_at TIMESTAMP,
    used_by BIGINT REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 19. 语音平台绑定表
-- =====================================================
CREATE TABLE voice_platform_bind (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    platform VARCHAR(32) NOT NULL,
    skill_id VARCHAR(128),
    bind_status VARCHAR(32) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 20. 操作日志表
-- =====================================================
CREATE TABLE operation_log (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    user_id BIGINT,
    operation VARCHAR(128) NOT NULL,
    module VARCHAR(64),
    target_id BIGINT,
    target_type VARCHAR(64),
    request_data JSONB,
    response_data JSONB,
    ip_address VARCHAR(64),
    user_agent VARCHAR(512),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 21. 语音映射表
-- =====================================================
CREATE TABLE voice_mapping (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    voice_name VARCHAR(128) NOT NULL,
    identify VARCHAR(128) NOT NULL,
    params JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 创建索引
-- =====================================================
CREATE INDEX idx_device_product ON device(product_id);
CREATE INDEX idx_device_status ON device(status);
CREATE INDEX idx_device_tenant ON device(tenant_id);
CREATE INDEX idx_product_tenant ON product(tenant_id);
CREATE INDEX idx_product_project ON product(project_id);
CREATE INDEX idx_product_key ON product(product_key);
CREATE INDEX idx_user_tenant ON "user"(tenant_id);
CREATE INDEX idx_user_username ON "user"(username);
CREATE INDEX idx_project_tenant ON project(tenant_id);
CREATE INDEX idx_project_key ON project(project_key);
CREATE INDEX idx_thing_model_product ON thing_model_define(product_id);
CREATE INDEX idx_device_property_device ON device_property(device_id);
CREATE INDEX idx_ota_firmware_product ON ota_firmware(product_id);
CREATE INDEX idx_ota_task_firmware ON ota_task(firmware_id);
CREATE INDEX idx_activation_code_product ON activation_code(product_id);
CREATE INDEX idx_activation_code_code ON activation_code(code);
CREATE INDEX idx_voice_platform_product ON voice_platform_bind(product_id);
CREATE INDEX idx_operation_log_tenant ON operation_log(tenant_id);
CREATE INDEX idx_operation_log_time ON operation_log(created_at);

-- =====================================================
-- 插入初始数据
-- =====================================================
-- 默认租户
INSERT INTO tenant (tenant_code, tenant_name, contact_name, contact_email, status)
VALUES ('default_tenant', '默认租户', '管理员', 'admin@example.com', 'active');

-- 默认管理员用户 (密码: admin123, BCrypt加密)
INSERT INTO "user" (tenant_id, username, password, nickname, role, status)
VALUES (1, 'admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '系统管理员', 'admin', 'active');

-- 默认项目
INSERT INTO project (tenant_id, project_key, project_name, description, status)
VALUES (1, 'default-project', '默认项目', '默认项目描述', 'active');

-- 示例品类
INSERT INTO category (name, name_en, level, sort_order) VALUES
('智能照明', 'Smart Lighting', 1, 1),
('智能安防', 'Smart Security', 1, 2),
('智能家电', 'Smart Appliance', 1, 3),
('智能传感器', 'Smart Sensor', 1, 4),
('智能插座', 'Smart Socket', 1, 5);

-- =====================================================
-- 创建更新触发器函数
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表添加自动更新 updated_at 的触发器
CREATE TRIGGER update_tenant_updated_at BEFORE UPDATE ON tenant FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_project_updated_at BEFORE UPDATE ON project FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_product_updated_at BEFORE UPDATE ON product FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_device_updated_at BEFORE UPDATE ON device FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_thing_model_updated_at BEFORE UPDATE ON thing_model_define FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ota_firmware_updated_at BEFORE UPDATE ON ota_firmware FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ota_task_updated_at BEFORE UPDATE ON ota_task FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_voice_platform_bind_updated_at BEFORE UPDATE ON voice_platform_bind FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
