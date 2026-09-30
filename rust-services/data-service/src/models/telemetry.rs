use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "telemetry")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = true)]
    pub id: i32,
    pub device_id: Option<i32>,
    pub tenant_id: Option<i32>,
    pub product_id: Option<i32>,
    pub identify: Option<String>,
    #[sea_orm(nullable)]
    pub value_type: Option<String>,
    pub value_float: Option<f64>,
    pub value_int: Option<i64>,
    #[sea_orm(nullable)]
    pub value_string: Option<String>,
    #[sea_orm(nullable)]
    pub value_bool: Option<bool>,
    #[sea_orm(index)]
    pub timestamp: Option<DateTime<Utc>>,
    #[sea_orm(created_at)]
    pub created_at: DateTime<Utc>,
    #[sea_orm(updated_at)]
    pub updated_at: DateTime<Utc>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
