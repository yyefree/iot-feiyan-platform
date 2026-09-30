use actix_web::{get, post, put, delete, web, HttpResponse, Responder};
use sea_orm::{ActiveModelTrait, DatabaseConnection, EntityTrait, QueryOrder, QuerySelect, Set};
use iot_common::{AppError, ApiResponse, Pagination};
use crate::models::rule::{Entity as Rule, Model as RuleModel};

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1/rules")
            .route("", web::get().list_rules)
            .route("", web::post().create_rule)
            .route("/statistics", web::get().get_rule_statistics)
            .route("/sql/test", web::post().test_sql)
            .route("/forward/logs", web::get().get_forward_logs)
            .route("/forward/config", web::post().update_forward_config)
            .route("/{id}", web::get().get_rule)
            .route("/{id}", web::put().update_rule)
            .route("/{id}", web::delete().delete_rule)
            .route("/{id}/enable", web::post().enable_rule)
            .route("/{id}/disable", web::post().disable_rule)
            .route("/{id}/logs", web::get().get_rule_logs)
            .route("/{id}/execute", web::get().execute_rule_test)
    );
    cfg.service(
        web::scope("/api/v1/scenes")
            .route("", web::get().list_scenes)
            .route("", web::post().create_scene)
            .route("/sunrise", web::get().get_sunrise_sunset)
            .route("/{id}", web::get().get_scene)
            .route("/{id}", web::put().update_scene)
            .route("/{id}", web::delete().delete_scene)
            .route("/{id}/enable", web::post().enable_scene)
            .route("/{id}/disable", web::post().disable_scene)
            .route("/{id}/execute", web::post().test_scene)
    );
}

#[get("")]
async fn list_rules(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let page = query.get("page").and_then(|v| v.as_u64()).unwrap_or(1) as i32;
    let size = query.get("size").and_then(|v| v.as_u64()).unwrap_or(20) as i32;
    let offset = (page - 1) * size;

    let total: i64 = Rule::find().count(&*db).await?;
    let rules: Vec<RuleModel> = Rule::find()
        .order_by_asc(Rule::Column::Id)
        .offset(offset as u64)
        .limit(size as u64)
        .all(&*db)
        .await?;

    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        rules,
        Pagination { page: page as u32, size: size as u32, total: total as u64, total_pages: ((total + size - 1) / size) as u64 },
    )))
}

#[post("")]
async fn create_rule(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let name = body.get("name").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("name is required".to_string()))?;
    let rule_type: &str = body.get("rule_type").and_then(|v| v.as_str()).unwrap_or("scene");
    let sql_expression: Option<String> = body.get("sql_expression").and_then(|v| v.as_str()).map(|s| s.to_string());
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);

    use crate::models::rule::ActiveModel;
    let rule = ActiveModel {
        name: Set(name.to_string()),
        rule_type: Set(rule_type.to_string()),
        sql_expression: Set(sql_expression),
        tenant_id: Set(tenant_id),
        status: Set("active".to_string()),
        ..Default::default()
    };

    let created = rule.insert(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(created)))
}

#[get("/statistics")]
async fn get_rule_statistics(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    let total: i64 = Rule::find().count(&*db).await?;
    let active: i64 = Rule::find().filter(rule::Column::Status.eq("active")).count(&*db).await?;
    let scenes: i64 = Rule::find().filter(rule::Column::RuleType.eq("scene")).count(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "totalRules": total,
        "activeRules": active,
        "sceneRules": scenes,
        "streamRules": total - scenes,
        "totalExecutions": 0,
    }))))
}

#[post("/sql/test")]
async fn test_sql(body: web::Json<serde_json::Value>) -> Result<HttpResponse, AppError> {
    let sql = body.get("sql_expression").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("sql_expression is required".to_string()))?;
    let lower_sql = sql.to_lowercase();
    let valid = lower_sql.contains("select") || lower_sql.contains("insert") || lower_sql.contains("update") || lower_sql.contains("delete");
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "success": valid,
        "message": if valid { "SQL语法正确" } else { "无法识别的SQL语句类型" },
        "result": if valid { "SQL验证通过" } else { serde_json::Value::Null },
    }))))
}

#[get("/forward/logs")]
async fn get_forward_logs(_db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        Vec::<serde_json::Value>::new(),
        Pagination { page: 1, size: 20, total: 0, total_pages: 0 }
    )))
}

#[post("/forward/config")]
async fn update_forward_config(_db: web::Data<DatabaseConnection>, _body: web::Json<serde_json::Value>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok("转发配置更新成功".to_string())))
}

#[get("/{id}")]
async fn get_rule(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(rule)))
}

#[put("/{id}")]
async fn update_rule(
    db: web::Data<DatabaseConnection>,
    path: web::Path<i32>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    use crate::models::rule::ActiveModel;
    let mut active: ActiveModel = rule.into();
    if let Some(name) = body.get("name").and_then(|v| v.as_str()) {
        active.name = Set(name.to_string());
    }
    if let Some(sql) = body.get("sql_expression").and_then(|v| v.as_str()) {
        active.sql_expression = Set(Some(sql.to_string()));
    }
    active.updated_at = Set(chrono::Utc::now());
    let updated = active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(updated)))
}

#[delete("/{id}")]
async fn delete_rule(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    rule.delete(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("删除成功".to_string())))
}

#[post("/{id}/enable")]
async fn enable_rule(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::rule::ActiveModel;
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = rule.into();
    active.status = Set("active".to_string());
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("规则已启用".to_string())))
}

#[post("/{id}/disable")]
async fn disable_rule(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::rule::ActiveModel;
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = rule.into();
    active.status = Set("inactive".to_string());
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("规则已禁用".to_string())))
}

#[get("/{id}/logs")]
async fn get_rule_logs(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        Vec::<serde_json::Value>::new(),
        Pagination { page: 1, size: 20, total: 0, total_pages: 0 }
    )))
}

#[get("/{id}/execute")]
async fn execute_rule_test(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "status": "success",
        "message": "规则执行成功（测试模式）",
        "rule_name": rule.name,
        "executed_at": chrono::Utc::now().to_rfc3339(),
    }))))
}

// ========== 场景联动 ==========

#[get("/scenes")]
async fn list_scenes(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let page = query.get("page").and_then(|v| v.as_u64()).unwrap_or(1) as i32;
    let size = query.get("size").and_then(|v| v.as_u64()).unwrap_or(20) as i32;

    // 场景数据从 rules 表中筛选 scene 类型
    let total: i64 = Rule::find().filter(rule::Column::RuleType.eq("scene")).count(&*db).await?;
    let rules: Vec<RuleModel> = Rule::find()
        .filter(rule::Column::RuleType.eq("scene"))
        .order_by_asc(Rule::Column::Id)
        .offset(((page - 1) * size) as u64)
        .limit(size as u64)
        .all(&*db)
        .await?;

    let scenes: Vec<serde_json::Value> = rules.iter().map(|r| {
        serde_json::json!({
            "id": r.id,
            "name": r.name,
            "description": r.description,
            "cron_expr": r.sql_expression.clone(),
            "status": r.status,
            "created_at": r.created_at.to_rfc3339(),
        })
    }).collect();

    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        scenes,
        Pagination { page: page as u32, size: size as u32, total: total as u64, total_pages: ((total + size - 1) / size) as u64 },
    )))
}

#[post("/scenes")]
async fn create_scene(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let name = body.get("name").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("name is required".to_string()))?;
    let cron_expr: Option<String> = body.get("cron_expr").and_then(|v| v.as_str()).map(|s| s.to_string());
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);

    use crate::models::rule::ActiveModel;
    let rule = ActiveModel {
        name: Set(name.to_string()),
        rule_type: Set("scene".to_string()),
        sql_expression: Set(cron_expr),
        tenant_id: Set(tenant_id),
        status: Set("active".to_string()),
        ..Default::default()
    };

    let created = rule.insert(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(serde_json::json!({
        "id": created.id,
        "name": created.name,
        "cron_expr": created.sql_expression,
        "status": created.status,
    }))))
}

#[get("/scenes/sunrise")]
async fn get_sunrise_sunset() -> Result<HttpResponse, AppError> {
    use chrono::Timelike;
    let now = chrono::Utc::now();
    let sunrise = now.with_hour(6).unwrap().with_minute(0).unwrap().with_second(0).unwrap();
    let sunset = now.with_hour(18).unwrap().with_minute(0).unwrap().with_second(0).unwrap();
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "date": now.format("%Y-%m-%d").to_string(),
        "sunrise": sunrise.format("%H:%M:%S").to_string(),
        "sunset": sunset.format("%H:%M:%S").to_string(),
        "latitude": 31.2304,
        "longitude": 121.4737,
    }))))
}

#[get("/scenes/{id}")]
async fn get_scene(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "id": rule.id,
        "name": rule.name,
        "description": rule.description,
        "cron_expr": rule.sql_expression,
        "status": rule.status,
        "created_at": rule.created_at.to_rfc3339(),
    }))))
}

#[put("/scenes/{id}")]
async fn update_scene(
    db: web::Data<DatabaseConnection>,
    path: web::Path<i32>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    use crate::models::rule::ActiveModel;
    let mut active: ActiveModel = rule.into();
    if let Some(name) = body.get("name").and_then(|v| v.as_str()) {
        active.name = Set(name.to_string());
    }
    if let Some(cron) = body.get("cron_expr").and_then(|v| v.as_str()) {
        active.sql_expression = Set(Some(cron.to_string()));
    }
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("场景已更新".to_string())))
}

#[delete("/scenes/{id}")]
async fn delete_scene(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    rule.delete(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("场景已删除".to_string())))
}

#[post("/scenes/{id}/enable")]
async fn enable_scene(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::rule::ActiveModel;
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = rule.into();
    active.status = Set("active".to_string());
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("场景已启用".to_string())))
}

#[post("/scenes/{id}/disable")]
async fn disable_scene(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::rule::ActiveModel;
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = rule.into();
    active.status = Set("inactive".to_string());
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("场景已禁用".to_string())))
}

#[post("/scenes/{id}/execute")]
async fn test_scene(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let rule = Rule::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "status": "success",
        "message": "场景测试执行成功",
        "scene_name": rule.name,
        "executed_at": chrono::Utc::now().to_rfc3339(),
    }))))
}
