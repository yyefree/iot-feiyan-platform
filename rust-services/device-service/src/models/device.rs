use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "devices")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub id: i32,
    #[sea_orm(unique)]
    pub device_name: String,
    pub device_key: Option<String>,
    pub device_secret: String,
    pub product_id: Option<i32>,
    pub tenant_id: Option<i32>,
    pub status: String,
    pub online: bool,
    pub is_virtual: bool,
    pub description: Option<String>,
    #[sea_orm(json)]
    pub extra: Option<serde_json::Value>,
    pub last_online_at: Option<DateTime>,
    pub last_offline_at: Option<DateTime>,
    #[sea_orm(created_at, updated_at)]
    pub created_at: DateTime,
    pub updated_at: DateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
