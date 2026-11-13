use axum::{
    Json, Router,
    extract::{Path, State},
    http::StatusCode,
    routing::get,
};

use futures_util::TryStreamExt;
use mongodb::bson::{Document, doc};

use crate::AppState;

pub fn router() -> Router<AppState> {
    Router::new().route("/{variant}", get(list))
}

async fn list(
    State(state): State<AppState>,
    Path(variant): Path<String>,
) -> Result<Json<Vec<Document>>, StatusCode> {
    let collection = match variant.as_str() {
        "groups" => "groupTickets",
        "users" => "userTickets",
        _ => return Err(StatusCode::NOT_FOUND),
    };

    let tickets = state
        .db
        .collection::<Document>(collection)
        .find(doc! {})
        .sort(doc! { "priority": -1, "_id": 1 })
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?
        .try_collect::<Vec<Document>>()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    
    let mut tickets = tickets;

    for ticket in &mut tickets {
        if let Ok(id) = ticket.get_object_id("_id") {
            ticket.insert("_id", id.to_hex());
        }
    }

    Ok(Json(tickets))
}
