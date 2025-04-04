use log::info;
use tokio::time::{sleep, Duration};
use kube8s::deploy_app;
use serde_json::Value;

pub async fn process_message(message: String) {
    println!("Processing message: {}", message);

    // Parse the message as JSON
    let mut parsed = serde_json::from_str::<Value>(&message).unwrap();
    let app_name = parsed.get_mut("app_name").as_str().unwrap_or("default_app").to_string();
    let image = parsed.get_mut("image").as_str().unwrap_or("default_image").to_string();
    // let targetport = parsed.get_mut("targetport").as_str().unwrap_or("default_targetport").to_string();

    deploy_app(&app_name, &image).await.unwrap_or_else(|e| {
        println!("Error deploying app: {}", e);
    });
    drop(message);

    info!("Message processed successfully.");
}