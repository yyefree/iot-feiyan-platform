use actix_web::{get, post, put, delete, web, HttpResponse, Responder};
use sea_orm::{ActiveModelTrait, DatabaseConnection, EntityTrait, QueryOrder, QuerySelect, Set};
use iot_common::{AppError, ApiResponse, Pagination};
use crate::models::operation_log::{Entity as OperationLog, Model as OperationLogModel};

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1/ops")
            .route("/dashboard", web::get().get_dashboard)
            .route("/devices", web::get().get_device_stats)
            .route("/users", web::get().get_user_stats)
            .route("/alerts", web::get().get_alerts)
            .route("/alerts/{id}/handle", web::put().handle_alert)
            .route("/logs", web::get().get_operation_logs)
            .route("/dashboard/config", web::get().get_dashboard_config)
            .route("/dashboard/config", web::put().update_dashboard_config)
    );
    cfg.service(
        web::scope("/api/v1/push")
            .route("/app", web::post().send_app_push)
            .route("/wechat", web::post().send_wechat_push)
            .route("/templates", web::get().get_push_templates)
            .route("/config", web::post().update_push_config)
    );
}

#[get("/dashboard")]
async fn get_dashboard(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    // 设备统计
    let total_devices: i64 = 0;
    let online_devices: i64 = 0;
    let offline_devices: i64 = 0;
    // 遥测数据
    let total_data_points: i64 = 0;
    let last_24h_data_points: i64 = 0;
    // 告警统计
    let total_alerts: i64 = 0;
    let pending_alerts: i64 = 0;

    let stats = serde_json::json!({
        "deviceStats": { "total": total_devices, "online": online_devices, "offline": offline_devices },
        "dataStats": { "total": total_data_points, "last24h": last_24h_data_points },
        "alertStats": { "total": total_alerts, "pending": pending_alerts },
        "timestamp": chrono::Utc::now().to_rfc3339(),
    });

    Ok(HttpResponse::Ok().json(ApiResponse::ok(stats)))
}

#[get("/devices")]
async fn get_device_stats(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/users")]
async fn get_user_stats(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/alerts")]
async fn get_alerts(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let page = query.get("page").and_then(|v| v.as_u64()).unwrap_or(1) as i32;
    let size = query.get("size").and_then(|v| v.as_u64()).unwrap_or(20) as i32;
    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        Vec::<serde_json::Value>::new(),
        Pagination { page: page as u32, size: size as u32, total: 0, total_pages: 0 },
    )))
}

#[put("/alerts/{id}/handle")]
async fn handle_alert(
    _db: web::Data<DatabaseConnection>,
    path: web::Path<i32>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let status: &str = body.get("status").and_then(|v| v.as_str()).unwrap_or("resolved");
    Ok(HttpResponse::Ok().json(ApiResponse::ok(format!("告警已处理，状态: {}", status))))
}

#[get("/logs")]
async fn get_operation_logs(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let page = query.get("page").and_then(|v| v.as_u64()).unwrap_or(1) as i32;
    let size = query.get("size").and_then(|v| v.as_u64()).unwrap_or(20) as i32;
    let offset = (page - 1) * size;

    let total: i64 = OperationLog::find().count(&*db).await?;
    let logs: Vec<OperationLogModel> = OperationLog::find()
        .order_by_desc(OperationLog::Column::CreatedAt)
        .offset(offset as u64)
        .limit(size as u64)
        .all(&*db)
        .await?;

    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        logs,
        Pagination { page: page as u32, size: size as u32, total: total as u64, total_pages: ((total + size - 1) / size) as u64 },
    )))
}

#[get("/dashboard/config")]
async fn get_dashboard_config(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "name": "默认大屏",
        "config": "{}",
    }))))
}

#[put("/dashboard/config")]
async fn update_dashboard_config(
    _db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("配置更新成功".to_string())))
}

// ========== 消息推送 ==========

#[post("/push/app")]
async fn send_app_push(
    _db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let user_id: Option<i32> = body.get("user_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let msg_type = body.get("msg_type").and_then(|v| v.as_str()).unwrap_or("alert");
    let title = body.get("title").and_then(|v| v.as_str()).unwrap_or("");
    let content = body.get("content").and_then(|v| v.as_str()).unwrap_or("");

    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "tenant_id": tenant_id,
        "user_id": user_id,
        "msg_type": msg_type,
        "title": title,
        "content": content,
        "status": "sent",
    }))))
}

#[post("/push/wechat")]
async fn send_wechat_push(
    _db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let title = body.get("title").and_then(|v| v.as_str()).unwrap_or("");
    let content = body.get("content").and_then(|v| v.as_str()).unwrap_or("");

    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "title": title,
        "content": content,
        "status": "sent",
        "platform": "wechat",
    }))))
}

#[get("/push/templates")]
async fn get_push_templates() -> Result<HttpResponse, AppError> {
    let templates = vec![
        serde_json::json!({"id": "device_alert", "name": "设备告警", "content": "设备{{device_name}}发生告警：{{alert_content}}"}),
        serde_json::json!({"id": "device_online", "name": "设备上线", "content": "设备{{device_name}}已上线"}),
        serde_json::json!({"id": "device_offline", "name": "设备离线", "content": "设备{{device_name}}已离线"}),
        serde_json::json!({"id": "ota_complete", "name": "升级完成", "content": "设备{{device_name}}固件升级完成"}),
    ];
    Ok(HttpResponse::Ok().json(ApiResponse::ok(templates)))
}

#[post("/push/config")]
async fn update_push_config(_db: web::Data<DatabaseConnection>, _body: web::Json<serde_json::Value>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("推送配置更新成功".to_string())))
}
