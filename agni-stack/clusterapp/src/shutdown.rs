use tokio::signal;

pub async fn shutdown_signal() {
    let _ = signal::ctrl_c().await; // Wait for Ctrl+C or SIGTERM
    println!("Shutdown signal received. Stopping Kafka consumer...");
}