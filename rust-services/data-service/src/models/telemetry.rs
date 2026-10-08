use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "telemetry")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub id: i64,
    pub device_id: Option<i32>,
    pub identify: Option<String>,
    pub value_type: Option<String>,
    pub value_float: Option<f64>,
    pub value_int: Option<i64>,
    pub value_string: Option<String>,
    pub value_bool: Option<bool>,
    pub timestamp: Option<DateTime>,
    #[sea_orm(created_at, updated_at)]
    pub created_at: DateTime,
    pub updated_at: DateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
