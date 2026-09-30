use actix_web::{web, App, HttpServer, Result};
use actix_cors::Cors;
use sea_orm::{Database, DatabaseConnection, Schema};
use std::env;
use tracing_subscriber::{fmt, prelude::*, EnvFilter};

mod models;
mod routes;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    tracing_subscriber::registry()
        .with(fmt::layer().with_writer(std::io::stderr))
        .with(EnvFilter::try_from_default_env().unwrap_or_else(|_| "info".into()))
        .init();

    dotenvy::dotenv().ok();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let db: DatabaseConnection = Database::connect(&database_url)
        .await
        .expect("Failed to connect to database");

    let schema = Schema::new(sea_orm::DatabaseBackend::Postgres);
    use models::telemetry::Entity as Telemetry;
    schema.create_table_from_model(Telemetry::default()).exec(&db).await.expect("Failed to migrate telemetry");

    let host = env::var("HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
    let port = env::var("PORT").unwrap_or_else(|_| "8084".to_string());
    let addr = format!("{}:{}", host, port);

    tracing::info!("data-service starting on {}", addr);

    HttpServer::new(move || {
        let cors = Cors::default().allow_any_origin().allow_any_method().allow_any_header();
        App::new()
            .app_data(web::Data::new(db.clone()))
            .wrap(cors)
            .service(health_check)
            .service(web::scope("").configure(routes::config))
    })
    .bind(&addr)?
    .run()
    .await
}

async fn health_check() -> Result<actix_web::HttpResponse, actix_web::Error> {
    Ok(actix_web::HttpResponse::Ok()
        .insert_header(("Content-Type", "application/json"))
        .json(iot_common::ApiResponse::ok("ok".to_string())))
}
