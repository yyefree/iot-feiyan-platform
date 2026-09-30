use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "rules")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = true)]
    pub id: i32,
    #[sea_orm(not_null)]
    pub name: String,
    #[sea_orm(nullable)]
    pub description: Option<String>,
    #[sea_orm(default = r#""scene""#.to_string())]
    pub rule_type: String,
    pub tenant_id: Option<i32>,
    #[sea_orm(column_type = "Text", nullable)]
    pub sql_expression: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub trigger_config: Option<String>,
    #[sea_orm(column_type = "Text", nullable)]
    pub action_config: Option<String>,
    #[sea_orm(default = r#""active""#.to_string())]
    pub status: String,
    pub created_by: Option<i32>,
    pub last_triggered: Option<DateTime<Utc>>,
    #[sea_orm(default = 0)]
    pub trigger_count: i32,
    #[sea_orm(created_at)]
    pub created_at: DateTime<Utc>,
    #[sea_orm(updated_at)]
    pub updated_at: DateTime<Utc>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
