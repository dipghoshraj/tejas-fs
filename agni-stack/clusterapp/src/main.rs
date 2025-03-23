/*
Code: src/main.rs
Cargo.toml
[package]
name = "clusterApp"
version = "0.1.0"
edition = "2025"
*/

mod config;
mod process;
mod shutdown;
mod consumer;


use env_logger;
use consumer::start_consumer;


#[tokio::main]
async fn main() {
    env_logger::init();
    
    let config = config::ConsumerConfig::load();
    println!("Starting consumer with config: {:?}", config);
    start_consumer(config).await;
}
