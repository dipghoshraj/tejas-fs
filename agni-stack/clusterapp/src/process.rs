use log::info;
use tokio::time::{sleep, Duration};


pub async fn process_message(message: String) {
    // info!("Processing message: {}", message);
    println!("Processing message: {}", message);

    // Simulate processing time
    sleep(Duration::from_millis(500)).await;
    println!("Message processed successfully.");
    drop(message);

    info!("Message processed successfully.");
}