use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "devices")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = true)]
    pub id: i32,
    #[sea_orm(unique, not_null)]
    pub device_name: String,
    #[sea_orm(unique)]
    pub device_key: Option<String>,
    #[sea_orm(not_null)]
    pub device_secret: String,
    pub product_id: Option<i32>,
    pub product_key: Option<String>,
    pub tenant_id: Option<i32>,
    #[sea_orm(default = r#""offline""#.to_string())]
    pub status: String,
    #[sea_orm(default = false)]
    pub online: bool,
    pub last_online_at: Option<DateTime<Utc>>,
    pub last_offline_at: Option<DateTime<Utc>>,
    pub ip_address: Option<String>,
    pub firmware_version: Option<String>,
    pub region: Option<String>,
    pub group_id: Option<i32>,
    #[sea_orm(column_type = "Text", nullable)]
    pub tags: Option<String>,
    #[sea_orm(nullable)]
    pub description: Option<String>,
    #[sea_orm(default = false)]
    pub is_virtual: bool,
    #[sea_orm(default = 60)]
    pub virtual_interval: i32,
    #[sea_orm(created_at)]
    pub created_at: DateTime<Utc>,
    #[sea_orm(updated_at)]
    pub updated_at: DateTime<Utc>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
