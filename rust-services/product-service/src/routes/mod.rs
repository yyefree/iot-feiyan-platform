use actix_web::{get, post, put, delete, web, HttpResponse};
use sea_orm::{ActiveModelTrait, DatabaseConnection, EntityTrait, QueryOrder, QuerySelect, Set};
use iot_common::{AppError, ApiResponse, Pagination};
use crate::models::product::{Entity, Model as ProductModel};

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1/products")
            .route("", web::get().list_products)
            .route("", web::post().create_product)
            .route("/statistics", web::get().get_product_statistics)
            .route("/{id}", web::get().get_product)
            .route("/{id}", web::put().update_product)
            .route("/{id}", web::delete().delete_product)
            .route("/{id}/publish", web::post().publish_product)
            .route("/{id}/properties", web::get().get_product_properties)
            .route("/{id}/services", web::get().get_product_services)
            .route("/{id}/events", web::get().get_product_events)
            .route("/{id}/tsl/export", web::get().export_tsl)
    );
}

#[get("")]
async fn list_products(
    db: web::Data<DatabaseConnection>,
    query: web::Query<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let page = query.get("page").and_then(|v| v.as_u64()).unwrap_or(1) as i32;
    let size = query.get("size").and_then(|v| v.as_u64()).unwrap_or(20) as i32;
    let offset = (page - 1) * size;

    let total: i64 = Entity::find().count(&*db).await?;
    let products: Vec<ProductModel> = Entity::find()
        .order_by_asc(Entity::Column::Id)
        .offset(offset as u64)
        .limit(size as u64)
        .all(&*db)
        .await?;

    Ok(HttpResponse::Ok().json(ApiResponse::ok_paginated(
        products,
        Pagination { page: page as u32, size: size as u32, total: total as u64, total_pages: ((total + size - 1) / size) as u64 },
    )))
}

#[post("")]
async fn create_product(
    db: web::Data<DatabaseConnection>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let product_name = body.get("product_name").and_then(|v| v.as_str()).ok_or(AppError::ValidationError("product_name is required".to_string()))?;
    let tenant_id: Option<i32> = body.get("tenant_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let project_id: Option<i32> = body.get("project_id").and_then(|v| v.as_i64()).map(|v| v as i32);
    let comm_type: &str = body.get("comm_type").and_then(|v| v.as_str()).unwrap_or("wifi");
    let auth_type: &str = body.get("auth_type").and_then(|v| v.as_str()).unwrap_or("one_device_one_key");
    let description: Option<String> = body.get("description").and_then(|v| v.as_str()).map(|s| s.to_string());

    let product_key = format!("p{}", uuid::Uuid::new_v4().to_string().replace("-", "").chars().take(16).collect::<String>());

    use crate::models::product::ActiveModel;
    let product = ActiveModel {
        product_name: Set(product_name.to_string()),
        product_key: Set(product_key),
        tenant_id: Set(tenant_id),
        project_id: Set(project_id),
        comm_type: Set(comm_type.to_string()),
        auth_type: Set(auth_type.to_string()),
        status: Set("draft".to_string()),
        description: Set(description),
        ..Default::default()
    };

    let created = product.insert(&*db).await?;
    Ok(HttpResponse::Created().json(ApiResponse::ok(created)))
}

#[get("/statistics")]
async fn get_product_statistics(db: web::Data<DatabaseConnection>) -> Result<HttpResponse, AppError> {
    let total: i64 = Entity::find().count(&*db).await?;
    let published: i64 = Entity::find().filter(Entity::Column::Status.eq("published")).count(&*db).await?;
    let draft: i64 = Entity::find().filter(Entity::Column::Status.eq("draft")).count(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "totalCount": total,
        "publishedCount": published,
        "draftCount": draft,
    }))))
}

#[get("/{id}")]
async fn get_product(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let product = Entity::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(product)))
}

#[put("/{id}")]
async fn update_product(
    db: web::Data<DatabaseConnection>,
    path: web::Path<i32>,
    body: web::Json<serde_json::Value>,
) -> Result<HttpResponse, AppError> {
    let product = Entity::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    use crate::models::product::ActiveModel;
    let mut active: ActiveModel = product.into();
    if let Some(name) = body.get("product_name").and_then(|v| v.as_str()) {
        active.product_name = Set(name.to_string());
    }
    if let Some(comm_type) = body.get("comm_type").and_then(|v| v.as_str()) {
        active.comm_type = Set(comm_type.to_string());
    }
    if let Some(auth_type) = body.get("auth_type").and_then(|v| v.as_str()) {
        active.auth_type = Set(auth_type.to_string());
    }
    if let Some(desc) = body.get("description").and_then(|v| v.as_str()) {
        active.description = Set(Some(desc.to_string()));
    }
    active.updated_at = Set(chrono::Utc::now());
    let updated = active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok(updated)))
}

#[delete("/{id}")]
async fn delete_product(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let product = Entity::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let _ = product.delete(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("删除成功".to_string())))
}

#[post("/{id}/publish")]
async fn publish_product(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    use crate::models::product::ActiveModel;
    let product = Entity::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let mut active: ActiveModel = product.into();
    active.status = Set("published".to_string());
    active.updated_at = Set(chrono::Utc::now());
    active.update(&*db).await?;
    Ok(HttpResponse::Ok().json(ApiResponse::ok("产品已发布".to_string())))
}

#[get("/{id}/properties")]
async fn get_product_properties(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/{id}/services")]
async fn get_product_services(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/{id}/events")]
async fn get_product_events(_db: web::Data<DatabaseConnection>, _path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(ApiResponse::ok(Vec::<serde_json::Value>::new())))
}

#[get("/{id}/tsl/export")]
async fn export_tsl(db: web::Data<DatabaseConnection>, path: web::Path<i32>) -> Result<HttpResponse, AppError> {
    let product = Entity::find_by_id(path.into_inner()).one(&*db).await?.ok_or(AppError::NotFound)?;
    let thing_model: serde_json::Value = serde_json::from_str(product.thing_model.as_deref().unwrap_or("{}"))
        .unwrap_or_else(|_| serde_json::json!({}));
    Ok(HttpResponse::Ok().json(ApiResponse::ok(serde_json::json!({
        "schema": "https://iot.aliyun.com/tsl/v1",
        "profile": {"productKey": product.product_key},
        "properties": thing_model.get("properties").cloned().unwrap_or(serde_json::json!([])),
        "services": thing_model.get("services").cloned().unwrap_or(serde_json::json!([])),
        "events": thing_model.get("events").cloned().unwrap_or(serde_json::json!([])),
    }))))
}
