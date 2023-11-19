// adapted from: https://github.com/tokio-rs/axum/blob/main/examples/graceful-shutdown/src/main.rs

use axum::{extract::Query, http::StatusCode, response::IntoResponse, routing::get, Json, Router};
use dotenv::dotenv;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, env, net::SocketAddr};
use tokio::signal;

#[tokio::main]
async fn main() {
    dotenv().unwrap();
    // TODO: we would want to pass this into a handler but that proved to be a time suck due
    // to my lack of rust syntax experience. we want this to explode on start not when endpoint is called
    env::var("WEBSTER_THESAURUS_API_KEY").expect("WEBSTER_THESAURUS_API_KEY not in env");
    // build our application with a route
    let app = Router::new().route("/api/v1/synonyms", get(handler_get_synonyms));

    // run it
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    println!("listening on {addr}");
    hyper::Server::bind(&addr)
        .serve(app.into_make_service())
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

#[derive(Debug, Deserialize)]
struct GetSynonymsReq {
    word: String,
}

#[derive(Debug, Serialize)]
struct GetSynonymsResp {
    synonymns: Vec<String>,
}

#[derive(Deserialize, Debug)]
struct ThesaurusPartialResp {
    meta: HashMap<String, serde_json::Value>,
}

async fn handler_get_synonyms(params: Query<GetSynonymsReq>) -> impl IntoResponse {
    println!("received text from wodrserweb service: {}", params.word);
    // TODO: I dont want to read from env everytime and I tried to use a closure to inject
    // the api-key during handler creation but it was too much of a pain
    // to try to do something that simple.
    let resp = reqwest::blocking::get(format!(
        "https://www.dictionaryapi.com/api/v3/references/thesaurus/json/{}?key={}",
        params.word,
        env::var("WEBSTER_THESAURUS_API_KEY").expect("WEBSTER_THESAURUS_API_KEY not in env"),
    ));

    if resp.is_err() {
        println!("resp: {:#?}", resp);
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(GetSynonymsResp {
                synonymns: Vec::new(),
            }),
        );
    }
    let resp = resp.unwrap();
    // let txt = resp.text();
    // println!("resp: {:#?}", txt);
    let partial: Result<Vec<ThesaurusPartialResp>, reqwest::Error> = resp.json();
    if partial.is_err() {
        println!("partial: {:#?}", partial);
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(GetSynonymsResp {
                synonymns: Vec::new(),
            }),
        );
    }
    let partial: Vec<ThesaurusPartialResp> = partial.unwrap();
    let meta_inner: &Vec<serde_json::Value> = match partial[0].meta.get("syns") {
        Some(m) => match m.as_array() {
            Some(s) => s,
            None => {
                println!("got None instead of syns");
                return (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(GetSynonymsResp {
                        synonymns: Vec::new(),
                    }),
                );
            }
        },
        None => {
            println!("got None instead of syns");
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(GetSynonymsResp {
                    synonymns: Vec::new(),
                }),
            );
        }
    };

    let mut syns: Vec<String> = Vec::new();
    for vect in meta_inner.iter() {
        let inner_vect = match vect.as_array() {
            Some(s) => s,
            None => {
                println!("got None instead of syns");
                return (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(GetSynonymsResp {
                        synonymns: Vec::new(),
                    }),
                );
            }
        };
        for inner_v in inner_vect.iter() {
            let inner_str = match inner_v.as_str() {
                Some(s) => s,
                None => {
                    println!("got None instead of syns");
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Json(GetSynonymsResp {
                            synonymns: Vec::new(),
                        }),
                    );
                }
            };
            syns.push(inner_str.to_string())
        }
    }
    println!("syns: {:?}", syns);
    let resp = GetSynonymsResp { synonymns: syns };

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
