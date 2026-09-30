use actix_web::{web, HttpResponse, Responder};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex;
use actix_web::http::header::{HeaderValue, CONTENT_TYPE};

pub fn config(cfg: &mut web::ServiceConfig, _backends: Arc<Mutex<HashMap<String, String>>>) {
    cfg.service(
        web::resource("/health").to(|| async {
            HttpResponse::Ok()
                .insert_header((CONTENT_TYPE, HeaderValue::from_static("application/json")))
                .json(serde_json::json!({"code": 0, "message": "ok", "service": "api-gateway"}))
        })
    );
}
