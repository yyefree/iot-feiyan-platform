use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "rules")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub id: i32,
    pub name: String,
    pub description: String,
    pub rule_type: String,
    pub sql_expression: Option<String>,
    pub trigger_config: Option<String>,
    pub action_config: Option<String>,
    pub tenant_id: Option<i32>,
    pub status: String,
    pub created_by: i32,
    pub last_triggered: Option<DateTime>,
    pub trigger_count: i64,
    #[sea_orm(created_at, updated_at)]
    pub created_at: DateTime,
    pub updated_at: DateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
