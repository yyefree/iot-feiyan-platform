use actix_web::{web, HttpResponse, App, HttpServer, Result, dev::ServiceRequest, Error};
use actix_cors::Cors;
use std::env;
use std::sync::Arc;
use tokio::sync::Mutex;
use tracing_subscriber::{fmt, prelude::*};
use http_body_util::Full;
use hyper::Uri;
use hyper_util::client::legacy::Client;
use hyper_util::rt::TokioExecutor;
use bytes::Bytes;
mod routes;

/// 后端服务地址映射
#[derive(Clone)]
struct BackendMap {
    map: Mutex<Vec<(String, String)>>,
}

impl BackendMap {
    fn new() -> Self {
        Self { map: Mutex::new(Vec::new()) }
    }

    async fn insert(&self, prefix: &str, url: &str) {
        let mut m = self.map.lock().await;
        m.push((prefix.to_string(), url.to_string()));
    }

    async fn get_backend(&self, path: &str) -> Option<String> {
        let m = self.map.lock().await;
        for (prefix, url) in m.iter() {
            if path.starts_with(prefix) {
                return Some(url.clone());
            }
        }
        None
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    tracing_subscriber::registry()
        .with(fmt::layer().with_writer(std::io::stderr))
        .init();

    dotenvy::dotenv().ok();

    let backends = Arc::new(BackendMap::new());
    backends.insert("/api/v1/devices", &env::var("DEVICE_ADDR").unwrap_or_else(|_| "http://localhost:8081".to_string())).await;
    backends.insert("/api/v1/products", &env::var("PRODUCT_ADDR").unwrap_or_else(|_| "http://localhost:8082".to_string())).await;
    backends.insert("/api/v1/rules", &env::var("RULE_ADDR").unwrap_or_else(|_| "http://localhost:8093".to_string())).await;
    backends.insert("/api/v1/scenes", &env::var("RULE_ADDR").unwrap_or_else(|_| "http://localhost:8093".to_string())).await;
    backends.insert("/api/v1/data", &env::var("DATA_ADDR").unwrap_or_else(|_| "http://localhost:8084".to_string())).await;
    backends.insert("/api/v1/ops", &env::var("OPS_ADDR").unwrap_or_else(|_| "http://localhost:8085".to_string())).await;
    backends.insert("/api/v1/push", &env::var("OPS_ADDR").unwrap_or_else(|_| "http://localhost:8085".to_string())).await;

    let host = env::var("HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
    let port = env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    let addr = format!("{}:{}", host, port);

    tracing::info!("api-gateway starting on {}", addr);

    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header()
            .max_age(3600);
        App::new()
            .app_data(web::Data::new(backends.clone()))
            .wrap(cors)
            .service(health_check)
            .service(proxy)
            .configure(routes::config)
    })
    .bind(&addr)?
    .run()
    .await
}

async fn health_check() -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok()
        .insert_header(("Content-Type", "application/json"))
        .json(serde_json::json!({"code": 0, "message": "ok", "service": "api-gateway-rust"})))
}

/// 通用反向代理
async fn proxy(
    req: ServiceRequest,
    backends: web::Data<Arc<BackendMap>>,
) -> Result<HttpResponse, Error> {
    let (parts, _body) = req.parts_mut();
    let method = parts.method.clone();
    let path = parts.uri.path().to_string();
    let query = parts.uri.query().unwrap_or("");

    let backend_url = match backends.get_backend(&path).await {
        Some(url) => url,
        None => {
            return Ok(HttpResponse::NotFound()
                .insert_header(("Content-Type", "application/json"))
                .json(serde_json::json!({"code": 404, "message": format!("No backend found for: {}", path)})));
        }
    };

    // 构建新 URI
    let new_uri: Uri = format!("{}{}", backend_url, path)
        .parse()
        .map_err(|e| actix_web::error::ErrorInternalServerError(format!("Invalid URI: {}", e)))?;

    // 如果是查询参数，附加到 URI
    let new_uri = if !query.is_empty() {
        let mut uri_str = new_uri.to_string();
        if uri_str.contains('?') {
            uri_str.push('&');
        } else {
            uri_str.push('?');
        }
        uri_str.push_str(query);
        uri_str.parse().map_err(|e| actix_web::error::ErrorInternalServerError(format!("Invalid URI: {}", e)))?
    } else {
        new_uri
    };

    // 读取请求体
    let body_bytes = match req.extract::<web::Bytes>().await {
        Ok(b) => b,
        Err(_) => Bytes::new(),
    };

    // 使用 hyper 客户端转发
    let client: Client<_, Full<Bytes>> = Client::builder(TokioExecutor::new()).build_http();
    let hyper_req = hyper::Request::builder()
        .method(method)
        .uri(new_uri)
        .header("host", new_uri.host().unwrap_or("localhost"))
        .body(Full::new(body_bytes))
        .map_err(|e| actix_web::error::ErrorInternalServerError(format!("Failed to build request: {}", e)))?;

    match client.request(hyper_req).await {
        Ok(response) => {
            let status = response.status();
            let body_bytes = match response.into_body().collect().await {
                Ok(b) => b.to_bytes(),
                Err(_) => Bytes::new(),
            };
            Ok(HttpResponse::build(actix_web::http::StatusCode::from_u16(status.as_u16()).unwrap_or(502))
                .insert_header(("Content-Type", "application/json"))
                .body(body_bytes.to_vec()))
        }
        Err(e) => {
            tracing::error!("Proxy error for {}: {}", path, e);
            Ok(HttpResponse::BadGateway()
                .insert_header(("Content-Type", "application/json"))
                .json(serde_json::json!({"code": 502, "message": format!("Backend unavailable: {}", e)})))
        }
    }
}
