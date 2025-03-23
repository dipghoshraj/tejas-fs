use log::info;
use tokio::net::unix::pipe::Receiver;
use crate::config::ConsumerConfig;
use crate::shutdown::shutdown_signal;
use rdkafka::consumer::{Consumer, StreamConsumer};
use tokio::time::{sleep, Duration};
use tokio::sync::{mpsc, Mutex};

use std::sync::Arc;


const MAX_CONCURRENT_TASKS: usize = 5; // Controls parallelism
const BACKOFF_DELAY: Duration = Duration::from_millis(500);


pub async fn process_message(message: &str) {
    info!("Processing message: {}", message);

    // Simulate processing time
    sleep(Duration::from_millis(500)).await;

    info!("Message processed successfully.");
}


// pub async fn StartConsumer(config: ConsumerConfig) {
    

//     consumer.subscribe(&[&config.topic]).expect("Failed to subscribe to topic");
//     let mut stream = consumer.stream();
//     info!("Starting consumer with config: {:?}", config);

//     let (tx, mut rx) = mpsc::channel::<String>(MAX_CONCURRENT_TASKS * 2);
//     let shutdown = tokio::spawn(shutdown_signal());
//     let rx = Arc::new(Mutex::new(rx)); // Shared across workers


//     // Spawn worker tasks to process messages in parallel


// }