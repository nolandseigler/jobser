// adapted from: https://github.com/tokio-rs/axum/blob/main/examples/graceful-shutdown/src/main.rs

use axum::{
    extract::Query, http::StatusCode, response::Html, response::IntoResponse, routing::get, Json,
    Router,
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::net::SocketAddr;
use tokio::signal;

#[tokio::main]
async fn main() {
    // build our application with a route
    let app = Router::new()
        .route("/", get(handler_root))
        .route("/api/v1/synonyms", get(handler_get_synonyms));

    // run it
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    println!("listening on {addr}");
    hyper::Server::bind(&addr)
        .serve(app.into_make_service())
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

async fn handler_root() -> Html<&'static str> {
    Html("<h1>Hello, World!</h1>")
}

#[derive(Debug, Deserialize)]
struct GetSynonymsReq {
    word: String,
}

#[derive(Debug, Serialize)]
struct GetSynonymsResp {
    synonymns: Vec<String>,
}

async fn handler_get_synonyms(params: Query<GetSynonymsReq>) -> impl IntoResponse {
    println!("received text from wodrserweb service: {}", params.word);
    let resp = reqwest::blocking::get(format!(
        "https://www.dictionaryapi.com/api/v3/references/thesaurus/json/{}?key={}",
        params.word, "not-real-api-key",
    ));

    let resp = match resp {
        Ok(x) => x.json::<serde_json::Value>(),
        Err(e) => {
            println!("error: {e}");
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(GetSynonymsResp {
                    synonymns: Vec::new(),
                }),
            );
        }
    };
    println!("{:#?}", resp);
    let resp = GetSynonymsResp {
        synonymns: vec![String::from("a")],
    };

    (StatusCode::OK, Json(resp))
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }

    println!("signal received, starting graceful shutdown");
}
