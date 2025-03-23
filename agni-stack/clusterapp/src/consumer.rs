use log::info;
use crate::config::ConsumerConfig;
use crate::shutdown::shutdown_signal;
use crate::process::process_message;
use rdkafka::consumer::StreamConsumer;
use tokio::time::{sleep, Duration};
use tokio::sync::{mpsc, Mutex};
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
        .create()
        .expect("Failed to create Kafka consumer");

    consumer
}


fn spawn_consumer_thread(rx : Arc<Mutex<mpsc::Receiver<String>>>) {
    for _ in 0..MAX_CONCURRENT_TASKS {
        let rx = rx.clone();
        tokio::spawn(async move {
            loop {
                let mut locked_rx = rx.try_lock().unwrap();
                if let Some(message) = locked_rx.recv().await {
                    drop(locked_rx); // Release lock early
                    process_message(&message).await;
                } else {
                    break;
                }
            }
        });
    }
}


async fn consume_message(consumer: StreamConsumer,
    tx: mpsc::Sender<String>,
    shutdown: tokio::task::JoinHandle<()>,){


        let mut stream = consumer.stream();
        tokio::select! {
            message = stream.next() => {
                match message {
                    Some(Ok(msg)) => {
                        let payload = match msg.payload_view::<str>() {
                            Some(Ok(s)) => s.to_string(),
                            Some(Err(e)) => {
                                eprintln!("Error while deserializing message payload: {:?}", e);
                                return;
                            }
                            None => {
                                eprintln!("Error while deserializing message payload");
                                return;
                            }
                        };
                        tx.send(payload).await.unwrap();
                    }
                    Some(Err(e)) => {
                        eprintln!("Error while receiving from Kafka: {:?}", e);
                        sleep(BACKOFF_DELAY).await; // Backoff on error
                    }
                    None => {
                        eprintln!("Stream terminated");
                        return;
                    }
                }
            }
            _ = shutdown => {
                info!("Shutting down consumer");
                return;
            }
        }
}


pub async fn start_consumer(config: ConsumerConfig) {

    let consumer = create_consumer(&config);
    let (tx, rx) = mpsc::channel::<String>(MAX_CONCURRENT_TASKS * 2);
    let rx = Arc::new(Mutex::new(rx));

    let shutdown = tokio::spawn(shutdown_signal());
    spawn_consumer_thread(rx.clone());

    consume_message(consumer, tx, shutdown).await;

}
