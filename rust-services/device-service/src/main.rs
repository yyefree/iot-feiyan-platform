use actix_web::{web, App, HttpServer, Error, HttpResponse};
use actix_cors::Cors;
use sea_orm::{Database, DatabaseConnection};
use std::env;
use tracing_subscriber::{fmt, prelude::*};

mod models;
mod routes;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    tracing_subscriber::registry()
        .with(fmt::layer().with_writer(std::io::stderr))
        .init();

    dotenvy::dotenv().ok();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let db: DatabaseConnection = Database::connect(&database_url)
        .await
        .expect("Failed to connect to database");

    let host = env::var("HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
    let port = env::var("PORT").unwrap_or_else(|_| "8081".to_string());
    let addr = format!("{}:{}", host, port);

    tracing::info!("device-service starting on {}", addr);

    HttpServer::new(move || {
        let cors = Cors::default().allow_any_origin().allow_any_method().allow_any_header();
        App::new()
            .app_data(web::Data::new(db.clone()))
            .wrap(cors)
            .route("/health", web::get().to(health_check))
            .configure(routes::config)
    })
    .bind(&addr)?
    .run()
    .await
}

async fn health_check() -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok()
        .insert_header(("Content-Type", "application/json"))
        .json(iot_common::ApiResponse::ok("ok".to_string())))
}
