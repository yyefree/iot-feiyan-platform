use actix_web::{get, post, put, delete, web, HttpResponse, Responder};
use sea_orm::{ActiveModelTrait, DatabaseConnection, EntityTrait, QueryOrder, QuerySelect, Set};
use iot_common::{AppError, ApiResponse, Pagination};
use crate::models::telemetry::{Entity as Telemetry, Model as TelemetryModel};

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1/data")
            .route("/telemetry", web::post().report_telemetry)
            .route("/telemetry/batch", web::post().batch_report_telemetry)
            .route("/telemetry", web::get().query_telemetry)
            .route("/telemetry/device/{id}", web::get().get_device_telemetry)
            .route("/telemetry/device/{id}/latest", web::get().get_latest_telemetry)
            .route("/telemetry/device/{id}/stats", web::get().get_device_stats)
            .route("/properties", web::get().get_latest_properties)
            .route("/properties/device/{id}", web::get().get_device_latest_properties)
            .route("/statistics", web::get().get_data_statistics)
            .route("/statistics/overview", web::get().get_statistics_overview)
            .route("/charts/line", web::get().get_line_chart)
            .route("/charts/bar", web::get().get_bar_chart)
            .route("/charts/pie", web::get().get_pie_chart)
            .route("/charts/radar", web::get().get_radar_chart)
            .route("/export", web::post().export_data)
    );
}

#[post("/telemetry")]
async fn report_telemetry(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device_id: Option<i32> = body.get("device_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let identify: Option<String> = body.get("identify").and_then(|v| v.as_str()).map(|s| s.to_string());
    let value_float: Option<f64> = body.get("value_float").and_then(|v| v.as_f64());
    let value_int: Option<i64> = body.get("value_int").and_then(|v| v.as_i64());
    let value_string: Option<String> = body.get("value_string").and_then(|v| v.as_str()).map(|s| s.to_string());
    let value_bool: Option<bool> = body.get("value_bool").and_then(|v| v.as_bool());
    let timestamp_str: Option<String> = body.get("timestamp").and_then(|v| v.as_str()).map(|s| s.to_string());

    let timestamp = if let Some(ts_str) = &timestamp_str {
        chrono::DateTime::parse_from_rfc3339(ts_str).ok().map(|dt| dt.with_timezone(&chrono::Utc))
    } else {
        None
    }.unwrap_or_else(chrono::Utc::now);

    // 确定值类型
    let value_type = if value_bool.is_some() { "bool" }
        else if value_int.is_some() { "int" }
        else if value_string.is_some() { "string" }
        else { "float" };

    use crate::models::telemetry::ActiveModel;
    let telemetry = ActiveModel {
        device_id: Set(device_id),
        identify: Set(identify),
        value_type: Set(Some(value_type.to_string())),
        value_float: Set(value_float),
        value_int: Set(value_int),
        value_string: Set(value_string),
        value_bool: Set(value_bool),
        timestamp: Set(Some(timestamp)),
        ..Default::default()
    };

    let created = telemetry.insert(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(created)))
}

#[post("/telemetry/batch")]
async fn batch_report_telemetry(
    db: web::Data<DatabaseConnection>,
    body: web::Json<Vec<serde_json::Value>>,
) -> Result<HttpResponse, AppError> {
    use crate::models::telemetry::ActiveModel;
    let mut records = Vec::new();
    for item in body.iter() {
        let device_id: Option<i32> = item.get("device_id").and_then(|v| v.as_i64()).map(|v| v as i32);
        let identify: Option<String> = item.get("identify").and_then(|v| v.as_str()).map(|s| s.to_string());
        let value_float: Option<f64> = item.get("value_float").and_then(|v| v.as_f64());
        let value_int: Option<i64> = item.get("value_int").and_then(|v| v.as_i64());
        let value_string: Option<String> = item.get("value_string").and_then(|v| v.as_str()).map(|s| s.to_string());
        let value_bool: Option<bool> = item.get("value_bool").and_then(|v| v.as_bool());
        let timestamp = chrono::Utc::now();
        let value_type = if value_bool.is_some() { "bool" } else if value_int.is_some() { "int" }
            else if value_string.is_some() { "string" } else { "float" };

        records.push(ActiveModel {
            device_id: Set(device_id),
            identify: Set(identify),
            value_type: Set(Some(value_type.to_string())),
            value_float: Set(value_float),
            value_int: Set(value_int),
            value_string: Set(value_string),
            value_bool: Set(value_bool),
            timestamp: Set(Some(timestamp)),
            ..Default::default()
        });
    }

    Telemetry::insert_many(records).exec(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(format!("成功上报 {} 条数据", body.len()))))
}

#[get("/telemetry")]
async fn query_telemetry(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device_id: Option<i32> = query.get("device_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let limit: u64 = query.get("limit").and_then(|v| v.as_u64()).unwrap_or(1000);

    let mut db_query = Telemetry::find();
    if let Some(id) = device_id {
        db_query = db_query.filter(telemetry::Column::DeviceId.eq(id));
    }
    let telemetries: Vec<TelemetryModel> = db_query
        .order_by_asc(telemetry::Column::Timestamp)
        .limit(limit)
        .all(&*db)
        .await?;

    Ok(HttpResponse::Ok().json(ApiResponse::ok(telemetries)))
}

#[get("/telemetry/device/{id}")]
async fn get_device_telemetry(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let telemetries: Vec<TelemetryModel> = Telemetry::find()
        .filter(telemetry::Column::DeviceId.eq(path.into_inner()))
        .order_by_desc(telemetry::Column::Timestamp)
        .limit(100)
        .all(&*db)
        .await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(telemetries)))
}

#[get("/telemetry/device/{id}/latest")]
async fn get_latest_telemetry(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let telemetries: Vec<TelemetryModel> = Telemetry::find()
        .filter(telemetry::Column::DeviceId.eq(path.into_inner()))
        .order_by_desc(telemetry::Column::Timestamp)
        .limit(10)
        .all(&*db)
        .await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(telemetries)))
}

#[get("/telemetry/device/{id}/stats")]
async fn get_device_stats(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "avg": 0.0, "min": 0.0, "max": 0.0, "count": 0,
        "sum": 0.0, "last_24h_count": 0,
    }))))
}

#[get("/properties")]
async fn get_latest_properties(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/properties/device/{id}")]
async fn get_device_latest_properties(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/statistics")]
async fn get_data_statistics(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/statistics/overview")]
async fn get_statistics_overview(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    let total_data: i64 = Telemetry::find().count(&*db).await?;
    let last_24h: i64 = Telemetry::find()
        .filter(telemetry::Column::Timestamp.gte(chrono::Utc::now() - chrono::Duration::hours(24)))
        .count(&*db)
        .await?;
    let total_devices: i64 = 0; // Would query devices table
    let online_devices: i64 = 0;
    let total_identifiers: i64 = Telemetry::find().group_by([telemetry::Column::Identify]).count(&*db).await?;

    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "totalDevices": total_devices,
        "onlineDevices": online_devices,
        "offlineDevices": total_devices - online_devices,
        "totalDataPoints": total_data,
        "last24hDataPoints": last_24h,
        "totalIdentifiers": total_identifiers,
    }))))
}

#[get("/charts/line")]
async fn get_line_chart(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device_id: Option<i32> = query.get("device_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let identify: Option<String> = query.get("identify").and_then(|v| v.as_str()).map(|s| s.to_string());

    let mut db_query = Telemetry::find();
    if let Some(id) = device_id {
        db_query = db_query.filter(telemetry::Column::DeviceId.eq(id));
    }
    if let Some(ident) = identify {
        db_query = db_query.filter(telemetry::Column::Identify.eq(ident));
    }
    let telemetries: Vec<TelemetryModel> = db_query
        .order_by_asc(telemetry::Column::Timestamp)
        .limit(1000)
        .all(&*db)
        .await?;

    let points: Vec<serde_json::Value> = telemetries.iter().map(|t| {
        serde_json::json!({
            "time": t.timestamp.map(|ts| ts.format("%Y-%m-%d %H:%M:%S").to_string()).unwrap_or_default(),
            "value": t.value_float.unwrap_or(0.0),
        })
    }).collect();

    Ok(HttpResponse::Ok().json(ApiResponse::ok(points)))
}

#[get("/charts/bar")]
async fn get_bar_chart(_db: web::Data<DatabaseConnection>, _query: web::Query<serde_json::Value>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/charts/pie")]
async fn get_pie_chart(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    // 设备状态分布 - 从 device 表查询
    let online: i64 = 0;
    let offline: i64 = 0;
    let points = vec![
        serde_json::json!({"name": "在线", "value": online}),
        serde_json::json!({"name": "离线", "value": offline}),
    ];
    Ok(HttpResponse::Ok().json(ApiResponse::ok(points)))
}

#[get("/charts/radar")]
async fn get_radar_chart() -> Result<HttpResponse, AppError> {
    let points = vec![
        serde_json::json!({"name": "在线率", "value": 85.0}),
        serde_json::json!({"name": "数据上报率", "value": 90.0}),
        serde_json::json!({"name": "指令响应率", "value": 92.0}),
        serde_json::json!({"name": "故障率", "value": 5.0}),
        serde_json::json!({"name": "用户满意度", "value": 95.0}),
    ];
    Ok(HttpResponse::Ok().json(ApiResponse::ok(points)))
}

#[post("/export")]
async fn export_data(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let device_id: i32 = body.get("device_id").and_then(|v| v.as_i64()).ok_or(AppError::ValidationError("device_id is required".to_string()))? as i32;
    let telemetries: Vec<TelemetryModel> = Telemetry::find()
        .filter(telemetry::Column::DeviceId.eq(device_id))
        .order_by_asc(telemetry::Column::Timestamp)
        .limit(10000)
        .all(&*db)
        .await?;

    // 生成 CSV
    let mut csv = String::from("time,device_id,identify,value_type,value_float,value_int,value_string\n");
    for t in &telemetries {
        csv.push_str(&format!("{},{},{},{},{},{},{}\n",
            t.timestamp.map(|ts| ts.to_rfc3339()).unwrap_or_default(),
            t.device_id.unwrap_or(0),
            t.identify.as_deref().unwrap_or(""),
            t.value_type.as_deref().unwrap_or(""),
            t.value_float.unwrap_or(0.0),
            t.value_int.unwrap_or(0),
            t.value_string.as_deref().unwrap_or(""),
        ));
    }

    // 返回 CSV 内容（实际应该用 StreamingResponseBody）
    Ok(HttpResponse::Ok()
        .insert_header(("Content-Type", "text/csv"))
        .insert_header(("Content-Disposition", "attachment; filename=telemetry_export.csv"))
        .body(csv))
}
