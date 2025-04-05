use std::result;
use std::io::{Error as IoError, ErrorKind};

use log::{info, error, debug};
use tokio::time::{sleep, Duration};
use kube8s::deploy_app;
use serde_json::Value;

pub async fn process_message(message: String) -> Result<(), Box<dyn std::error::Error + Send + Sync + 'static>> {
    debug!("Processing message: {}", message);

    // Parse the message as JSON
    let parsed: Value = match serde_json::from_str(&message) {
        Ok(value) => value,
        Err(e) => {
            error!("Failed to parse message as JSON: {}", e);
            return Err(Box::new(e));
        }
    };

    let app_name = parsed.get("app_name")
        .and_then(|v| v.as_str())
        .unwrap_or("default_app_name")
        .to_string();

    let image = parsed.get("image")
        .and_then(|v| v.as_str())
        .unwrap_or("default_image")
        .to_string();

    debug!("Parsed app_name: {}", app_name);
    debug!("Parsed image: {}", image);

    // Deploy the application
    match deploy_app(&app_name, &image).await {
        Ok(_) => {
            info!("Deployment successful for app: {}", app_name);
        },
        Err(e) => {
            error!("Error deploying app {}: {}", app_name, e);
            return Err(e as Box<dyn std::error::Error + Send + Sync + 'static>);
        }
    }
    info!("Message processed successfully.");
    Ok(())
}