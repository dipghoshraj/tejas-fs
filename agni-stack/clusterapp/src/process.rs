use log::info;
use tokio::time::{sleep, Duration};


pub async fn process_message(message: &str) {
    info!("Processing message: {}", message);
    println!("Processing message: {}", message);

    // Simulate processing time
    sleep(Duration::from_millis(500)).await;
    println!("Message processed successfully.");

    info!("Message processed successfully.");
}