use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "products")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = true)]
    pub id: i32,
    #[sea_orm(not_null)]
    pub product_name: String,
    #[sea_orm(unique, not_null)]
    pub product_key: String,
    pub tenant_id: Option<i32>,
    pub project_id: Option<i32>,
    pub category_id: Option<i32>,
    #[sea_orm(default = r#""direct_device""#.to_string())]
    pub product_type: String,
    #[sea_orm(default = r#""wifi""#.to_string())]
    pub comm_type: String,
    #[sea_orm(default = r#""one_device_one_key""#.to_string())]
    pub auth_type: String,
    #[sea_orm(default = r#""draft""#.to_string())]
    pub status: String,
    #[sea_orm(nullable)]
    pub icon: Option<String>,
    #[sea_orm(nullable)]
    pub version: Option<String>,
    #[sea_orm(nullable)]
    pub description: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub thing_model: Option<String>,
    #[sea_orm(created_at)]
    pub created_at: DateTime<Utc>,
    #[sea_orm(updated_at)]
    pub updated_at: DateTime<Utc>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
