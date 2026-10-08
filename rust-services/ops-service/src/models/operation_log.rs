use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "operation_logs")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub id: i32,
    pub operator: Option<String>,
    pub operator_id: Option<i32>,
    pub module: String,
    pub action: String,
    pub target_type: String,
    pub target_id: Option<i32>,
    pub detail: Option<String>,
    pub ip: Option<String>,
    #[sea_orm(created_at, updated_at)]
    pub created_at: DateTime,
    pub updated_at: DateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
