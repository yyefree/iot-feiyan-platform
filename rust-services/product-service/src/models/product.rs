use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "products")]
pub struct Model {
    #[sea_orm(primary_key)]
    #[serde(rename = "id")]
    pub id: i32,
    #[sea_orm(unique)]
    #[serde(rename = "product_name")]
    pub product_name: String,
    #[sea_orm(unique)]
    #[serde(rename = "product_key")]
    pub product_key: String,
    #[serde(rename = "tenant_id")]
    pub tenant_id: Option<i32>,
    #[serde(rename = "project_id")]
    pub project_id: Option<i32>,
    #[serde(rename = "comm_type")]
    pub comm_type: String,
    #[serde(rename = "auth_type")]
    pub auth_type: String,
    #[serde(rename = "status")]
    pub status: String,
    #[sea_orm(json)]
    #[serde(rename = "thing_model")]
    pub thing_model: Option<serde_json::Value>,
    #[serde(rename = "description")]
    pub description: Option<String>,
    #[sea_orm(created_at)]
    #[serde(rename = "created_at")]
    pub created_at: DateTime,
    #[sea_orm(updated_at)]
    #[serde(rename = "updated_at")]
    pub updated_at: DateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}
