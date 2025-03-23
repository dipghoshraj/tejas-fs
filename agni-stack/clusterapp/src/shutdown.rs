use tokio::sync::oneshot;
use tokio::signal;
use log::info;


pub async fn shutdown_signal() {
    let (tx, rx) = oneshot::channel();
    
    tokio::spawn(async move {
        signal::ctrl_c().await.expect("Failed to capture shutdown signal");
        let _ = tx.send(());
    });

    rx.await.expect("Shutdown signal failed to receive");
    info!("Received shutdown signal. Exiting...");
}