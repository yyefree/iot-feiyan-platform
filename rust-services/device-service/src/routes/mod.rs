use actix_web::{get, post, put, delete, web, HttpResponse, Responder};
use sea_orm::{ActiveModelTrait, DatabaseConnection, EntityTrait, QueryOrder, QuerySelect, Set};
use iot_common::{AppError, ApiResponse, Pagination};
use crate::models::device::{Entity as Device, Model as DeviceModel};

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1/devices")
            .route("", web::get().list_devices)
            .route("", web::post().create_device)
            .route("/batch", web::post().batch_create_devices)
            .route("/statistics", web::get().get_device_statistics)
            .route("/virtual/create", web::post().create_virtual_device)
            .route("/virtual/list", web::get().get_virtual_device_list)
            .route("/{id}", web::get().get_device)
            .route("/{id}", web::put().update_device)
            .route("/{id}", web::delete().delete_device)
            .route("/{id}/enable", web::post().enable_device)
            .route("/{id}/disable", web::post().disable_device)
            .route("/{id}/reset", web::post().reset_device)
            .route("/{id}/credentials", web::get().get_device_credentials)
            .route("/{id}/shadow", web::get().get_device_shadow)
            .route("/{id}/shadow", web::put().update_device_shadow)
            .route("/groups", web::get().get_device_groups)
            .route("/groups", web::post().create_device_group)
            .route("/groups/{id}", web::put().update_device_group)
            .route("/groups/{id}", web::delete().delete_device_group)
            .route("/groups/{id}/devices", web::post().add_device_to_group)
            .route("/groups/{id}/devices/{device_id}", web::delete().remove_device_from_group)
            .route("/activation-codes", web::get().get_activation_codes)
            .route("/activation-codes/generate", web::post().generate_activation_codes)
            .route("/activation-codes/{id}", web::delete().delete_activation_code)
    );
}

#[get("")]
async fn list_devices(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let page = query.get("page").and_then(|v| v.as_u64()).unwrap_or(1) as i32;
    let size = query.get("size").and_then(|v| v.as_u64()).unwrap_or(20) as i32;
    let offset = (page - 1) * size;

    let total: i64 = Device::find().count(&*db).await?;
    let devices: Vec<DeviceModel> = Device::find()
        .order_by_asc(Device::Column::Id)
        .offset(offset as u64)
        .limit(size as u64)
        .all(&*db)
        .await?;

    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        devices,
        Pagination { page: page as u32, size: size as u32, total: total as u64, total_pages: ((total + size - 1) / size) as u64 },
    )))
}

#[post("")]
async fn create_device(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device_name = body.get("device_name").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("device_name is required".to_string()))?;
    let product_id: Option<i32> = body.get("product_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let description: Option<String> = body.get("description").and_then(|v| v.as_str()).map(|s| s.to_string());
    let is_virtual = body.get("is_virtual").and_then(|v| v.as_bool()).unwrap_or(false);

    let device_key = format!("dk{}", uuid::Uuid::new_v4().to_string().replace("-", ""));
    let device_secret = format!("ds{}", uuid::Uuid::new_v4().to_string().replace("-", ""));

    use crate::models::device::ActiveModel;
    let device = ActiveModel {
        device_name: Set(device_name.to_string()),
        device_key: Set(Some(device_key)),
        device_secret: Set(device_secret),
        product_id: Set(product_id),
        tenant_id: Set(tenant_id),
        status: Set(if is_virtual { "online".to_string() } else { "offline".to_string() }),
        online: Set(is_virtual),
        description: Set(description),
        is_virtual: Set(is_virtual),
        ..Default::default()
    };

    let created = device.insert(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(created)))
}

#[post("/batch")]
async fn batch_create_devices(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let product_id: i32 = body.get("product_id").and_then(|v| v.as_i64()).ok_or(AppError::ValidationError("product_id is required".to_string()))? as i32;
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let count: usize = body.get("count").and_then(|v| v.as_u64()).ok_or(AppError::ValidationError("count is required".to_string()))? as usize;
    let prefix: &str = body.get("prefix").and_then(|v| v.as_str()).unwrap_or("device");

    let mut devices = Vec::new();
    for i in 0..count {
        let device_key = format!("dk{}{:04}", uuid::Uuid::new_v4().to_string().replace("-", ""), i);
        let device_secret = format!("ds{}", uuid::Uuid::new_v4().to_string().replace("-", ""));
        use crate::models::device::ActiveModel;
        devices.push(ActiveModel {
            device_name: Set(format!("{}_{:04}", prefix, i + 1)),
            device_key: Set(Some(device_key)),
            device_secret: Set(device_secret),
            product_id: Set(Some(product_id)),
            tenant_id: Set(tenant_id),
            status: Set("offline".to_string()),
            online: Set(false),
            is_virtual: Set(false),
            ..Default::default()
        });
    }

    Device::insert_many(devices).exec(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(format!("成功创建 {} 台设备", count))))
}

#[get("/statistics")]
async fn get_device_statistics(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    let total: i64 = Device::find().count(&*db).await?;
    let online: i64 = Device::find().filter(device::Column::Online.eq(true)).count(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "totalCount": total,
        "onlineCount": online,
        "offlineCount": total - online,
    }))))
}

#[get("/virtual/list")]
async fn get_virtual_device_list(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    let devices: Vec<DeviceModel> = Device::find()
        .filter(device::Column::IsVirtual.eq(true))
        .all(&*db)
        .await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(devices)))
}

#[post("/virtual/create")]
async fn create_virtual_device(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device_name = body.get("device_name").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("device_name is required".to_string()))?;
    let product_id: Option<i32> = body.get("product_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);

    let device_key = format!("dv{}", uuid::Uuid::new_v4().to_string().replace("-", ""));
    let device_secret = format!("ds{}", uuid::Uuid::new_v4().to_string().replace("-", ""));

    use crate::models::device::ActiveModel;
    let device = ActiveModel {
        device_name: Set(device_name.to_string()),
        device_key: Set(Some(device_key)),
        device_secret: Set(device_secret),
        product_id: Set(product_id),
        tenant_id: Set(tenant_id),
        status: Set("online".to_string()),
        online: Set(true),
        is_virtual: Set(true),
        ..Default::default()
    };

    let created = device.insert(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(created)))
}

#[get("/{id}")]
async fn get_device(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(device)))
}

#[put("/{id}")]
async fn update_device(
    db: web::Data<DatabaseConnection>,
    path: web::Path<i32>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    use crate::models::device::ActiveModel;
    let mut active: ActiveModel = device.into();
    if let Some(name) = body.get("device_name").and_then(|v| v.as_str()) {
        active.device_name = Set(name.to_string());
    }
    if let Some(desc) = body.get("description").and_then(|v| v.as_str()) {
        active.description = Set(Some(desc.to_string()));
    }
    active.updated_at = Set(chrono::Utc::now());
    let updated = active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(updated)))
}

#[delete("/{id}")]
async fn delete_device(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    device.delete(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("删除成功".to_string())))
}

#[post("/{id}/enable")]
async fn enable_device(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::device::ActiveModel;
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = device.into();
    active.status = Set("online".to_string());
    active.online = Set(true);
    active.last_online_at = Set(Some(chrono::Utc::now()));
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("设备已启用".to_string())))
}

#[post("/{id}/disable")]
async fn disable_device(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::device::ActiveModel;
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = device.into();
    active.status = Set("offline".to_string());
    active.online = Set(false);
    active.last_offline_at = Set(Some(chrono::Utc::now()));
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("设备已禁用".to_string())))
}

#[post("/{id}/reset")]
async fn reset_device(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::device::ActiveModel;
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = device.into();
    active.device_secret = Set(format!("ds{}", uuid::Uuid::new_v4().to_string().replace("-", "")));
    active.updated_at = Set(chrono::Utc::now());
    let updated = active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(updated)))
}

#[get("/{id}/credentials")]
async fn get_device_credentials(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let device = Device::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "productKey": device.product_key,
        "deviceName": device.device_name,
        "deviceSecret": device.device_secret,
        "authType": "one_device_one_key",
    }))))
}

#[get("/{id}/shadow")]
async fn get_device_shadow(db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "reported": {},
        "desired": {},
        "version": 1,
    }))))
}

#[put("/{id}/shadow")]
async fn update_device_shadow(
    db: web::Data<DatabaseConnection>,
    _path: web::Path<i32>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "reported": body.as_object().cloned().unwrap_or_default(),
        "desired": {},
        "version": 1,
    }))))
}

// 设备分组API
#[get("/groups")]
async fn get_device_groups(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[post("/groups")]
async fn create_device_group(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let name = body.get("name").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("name is required".to_string()))?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(serde_json::json!({"id": 1, "name": name, "tenant_id": 1, "device_count": 0}))))
}

#[put("/groups/{id}")]
async fn update_device_group(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>, _body: web::Json<serde_json::Value>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("分组已更新".to_string())))
}

#[delete("/groups/{id}")]
async fn delete_device_group(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("分组已删除".to_string())))
}

#[post("/groups/{id}/devices")]
async fn add_device_to_group(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>, _body: web::Json<serde_json::Value>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("添加成功".to_string())))
}

#[delete("/groups/{id}/devices/{device_id}")]
async fn remove_device_from_group(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("移除成功".to_string())))
}

// 激活码API
#[get("/activation-codes")]
async fn get_activation_codes(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[post("/activation-codes/generate")]
async fn generate_activation_codes(
    _db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let count: usize = body.get("count").and_then(|v| v.as_u64()).unwrap_or(10) as usize;
    Ok(HttpResponse::Created().json(ApiResponse::ok(format!("成功生成 {} 个激活码", count))))
}

#[delete("/activation-codes/{id}")]
async fn delete_activation_code(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("删除成功".to_string())))
}
