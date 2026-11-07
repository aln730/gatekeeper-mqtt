use axum::Router;

#[tokio::main]
async fn main() {
    let app = Router::new();

    let listener = tokio::net::TcpListener::bind("127.0.0.1:3001")
        .await
        .unwrap();

    println!("gatekeeper is reaching the finals");

    axum::serve(listener, app).await.unwrap();
}
