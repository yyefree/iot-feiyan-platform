use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Database error: {0}")]
    DbError(#[from] sea_orm::DbErr),
    #[error("Validation error: {0}")]
    ValidationError(String),
    #[error("Resource not found")]
    NotFound,
    #[error("Internal server error: {0}")]
    InternalError(String),
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ApiResponse<T> {
    pub code: i32,
    pub message: String,
    pub data: Option<T>,
    pub pagination: Option<Pagination>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Pagination {
    pub page: u32,
    pub size: u32,
    pub total: u64,
    pub total_pages: u64,
}

impl<T> ApiResponse<T> {
    pub fn ok(data: T) -> Self {
        Self {
            code: 0,
            message: "success".to_string(),
            data: Some(data),
            pagination: None,
        }
    }

    pub fn ok_paginated(data: T, pagination: Pagination) -> Self {
        Self {
            code: 0,
            message: "success".to_string(),
            data: Some(data),
            pagination: Some(pagination),
        }
    }

    pub fn error(code: i32, message: &str) -> Self {
        Self {
            code,
            message: message.to_string(),
            data: None,
            pagination: None,
        }
    }
}
