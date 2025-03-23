// Code: src/main.rs
// Cargo.toml
// [package]
// name = "clusterApp"
// version = "0.1.0"
// edition = "2025"
//


mod config;
mod process;


use env_logger;


fn main() {
    env_logger::init();
    
    let config = config::ConsumerConfig::load();
    println!("Starting consumer with config: {:?}", config);
}
