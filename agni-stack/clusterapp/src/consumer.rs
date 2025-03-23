use crate::config::ConsumerConfig;
use crate::shutdown::shutdown_signal;
use crate::process::process_message;
use rdkafka::consumer::{Consumer, StreamConsumer};
use tokio::time::{sleep, Duration};
use tokio::sync::{mpsc, Mutex, watch};
use std::sync::Arc;
use tokio_stream::StreamExt;
use rdkafka::Message;

const MAX_CONCURRENT_TASKS: usize = 5; // Controls parallelism
const BACKOFF_DELAY: Duration = Duration::from_millis(500);


fn create_consumer(config: &ConsumerConfig) -> StreamConsumer {

    let consumer: StreamConsumer = rdkafka::ClientConfig::new()
        .set("group.id", &config.group_id)
        .set("bootstrap.servers", &config.broker)
        .set("enable.auto.commit", "true")
        .set("auto.commit.interval.ms", "1000")
        .set("session.timeout.ms", "30000")
        .set("enable.partition.eof", "false")
        .set("auto.offset.reset", "earliest")
        .set("queued.max.messages.kbytes", "1024") // Reduce max buffer size
        .set("fetch.message.max.bytes", "1024000") // Limit per-message size
        .create()
        .expect("Failed to create Kafka consumer");

    consumer
}


fn spawn_consumer_thread(rx : Arc<Mutex<mpsc::Receiver<String>>>) {
    for _ in 0..MAX_CONCURRENT_TASKS {
        let rx = rx.clone();
        tokio::spawn(async move {
            loop {
                if let Ok(mut locked_rx) = rx.try_lock() {
                    if let Some(message) = locked_rx.recv().await {
                        drop(locked_rx); // Release lock early
                        process_message(message).await;

                    }
                } else {
                    sleep(Duration::from_millis(10)).await; // Prevent busy looping
                }
            }
        });
    }
}


async fn consume_message(consumer: StreamConsumer,
    tx: mpsc::Sender<String>,
    mut shutdown: watch::Receiver<bool>,){


        let mut stream = consumer.stream();
        loop {
            tokio::select! {
                message = stream.next() => {
                    match message {
                        Some(Ok(msg)) => {
                            let payload = match msg.payload_view::<str>() {
                                Some(Ok(s)) => s.to_string(),
                                Some(Err(e)) => {
                                    println!("Error while deserializing message payload: {:?}", e);
                                    continue;
                                }
                                None => {
                                    println!("Error while deserializing message payload");
                                    continue;
                                }
                            };
                            if let Err(_) = tx.send(payload).await {
                                println!("Receiver dropped. Stopping consumer.");
                                return;
                            }
                        }
                        Some(Err(e)) => {
                            println!("Error while receiving from Kafka: {:?}", e);
                            sleep(BACKOFF_DELAY).await; // Backoff on error
                        }
                        None => {
                            sleep(Duration::from_millis(500)).await; // Short sleep to avoid high CPU usage
                        }
                    }
                }

                _ = shutdown.changed() => {
                    println!("Shutdown detected, exiting consumer.");
                    return;
                }
    
            }
        }
}


pub async fn start_consumer(config: ConsumerConfig) {

    let consumer = create_consumer(&config);
    consumer
        .subscribe(&[&config.topic])
        .expect("Failed to subscribe to Kafka topic");


    let (tx, rx) = mpsc::channel::<String>(MAX_CONCURRENT_TASKS * 2);
    let rx = Arc::new(Mutex::new(rx));
    let (shutdown_tx, shutdown_rx) = watch::channel(false);


    _ = tokio::spawn(shutdown_signal());
    spawn_consumer_thread(rx.clone());

    // consume_message(consumer, tx, shutdown).await;

    let consumer_task = tokio::spawn(consume_message(consumer, tx, shutdown_rx));


    shutdown_signal().await;
    shutdown_tx.send(true).unwrap(); // Notify all tasks to shut down
    consumer_task.await.unwrap(); // Ensure consumer task exits cleanly
    println!("Kafka consumer stopped.");

}
