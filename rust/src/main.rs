mod tickets;

use axum::Router;
use mongodb::Client;

#[derive(Clone)]
pub struct AppState {
    pub db: mongodb::Database,
}

#[tokio::main]
async fn main() {
    let uri = std::env::var("GK_MONGO_SERVER").expect("mongo not set");

    let client = Client::with_uri_str(&uri)
        .await
        .expect("failed to connect to mongo");

    let state = AppState {
        db: client.database("gatekeeper"),
    };

    println!("reached the mongodb finals");

    let app = Router::new()
        .nest("/tickets", tickets::router())
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("127.0.0.1:3001")
        .await
        .unwrap();

    println!("gatekeeper is reaching the finals");

    axum::serve(listener, app).await.unwrap();
}
