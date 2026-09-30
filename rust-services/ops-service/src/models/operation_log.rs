use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "operation_logs")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = true)]
    pub id: i32,
    pub tenant_id: Option<i32>,
    pub user_id: Option<i32>,
    #[sea_orm(not_null)]
    pub operation: String,
    #[sea_orm(nullable)]
    pub module: Option<String>,
    pub target_id: Option<i32>,
    #[sea_orm(nullable)]
    pub target_type: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub request_data: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub response_data: Option<String>,
    #[sea_orm(nullable)]
    pub ip_address: Option<String>,
    #[sea_orm(nullable)]
    pub user_agent: Option<String>,
    #[sea_orm(created_at)]
    pub created_at: DateTime<Utc>,
    #[sea_orm(updated_at)]
    pub updated_at: DateTime<Utc>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
