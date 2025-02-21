from kafka import KafkaConsumer, KafkaProducer
import logging, signal, asyncio, os
from concurrent.futures import ThreadPoolExecutor
from kafka.errors import KafkaError
import json

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("KafkaWorker")


KAFKA_BROKER_URL =os.getenv("KAFKA_BROKER_URL", "localhost:9092")
KAFKA_UPLOAD_TOPIC = os.getenv("KAFKA_TOPIC", "metadata")
KAFKA_DLQ_TOPIC = os.getenv("KAFKA_DLQ_TOPIC", "dlq")

executor = ThreadPoolExecutor(max_workers=10)


compress_consumer = KafkaConsumer(
    KAFKA_UPLOAD_TOPIC,
    bootstrap_servers=KAFKA_BROKER_URL,
    value_deserializer=lambda x: x.decode("utf-8"),
    heartbeat_interval_ms=10000
)

dlq_producer = KafkaProducer(
    bootstrap_servers=KAFKA_BROKER_URL,
    value_serializer=lambda x: x.encode("utf-8"),
    acks='all'
)


async def process_message(message):
    """
    Process the kafka message

    :params: message: the kafka message
    :return: None
    """

    try:
        logger.error(f"Process the messaeg {message}")
        await asyncio.sleep(1)
        logger.info(f"Processed message: {message}")
    except Exception as e:
        logger.error(f"Failed to process image {message}: {str(e)}")


async def send_to_dlq(message, error):
    """
    Send failed messages to a Dead Letter Queue (DLQ).

    :param message: The failed message
    :param error: The error message
    :return: None
    """
    try:
        logger.error(f"Sending message to DLQ: {message}, Error: {error}")
        dlq_producer.send(KAFKA_DLQ_TOPIC, value=message)
    except Exception as e:
        logger.error(f"Failed to send message to DLQ: {str(e)}")


async def process_messages_batch(messages, consumer):
    """
    Process a batch of messages concurrently and commit offsets.

    :param messages: List of Kafka messages
    :param consumer: Kafka consumer instance
    :return: None
    """
    tasks = [process_message(message.value) for message in messages]
    await asyncio.gather(*tasks)
    # Commit offsets after successful processing
    consumer.commit()

async def consume_messages():
    """
    Consume messages from Kafka in an asynchronous loop.
    """
    while True:
        consumer = None
        try:
            consumer = compress_consumer 
            logger.info("Kafka worker started, waiting for messages...")

            while True:
                try:
                    # Poll for messages
                    messages = consumer.poll(timeout_ms=1000, max_records=100)
                    if messages:
                        for topic_partition, msg_list in messages.items():
                            await process_messages_batch(msg_list, consumer)
                except KafkaError as e:
                    logger.error(f"Kafka error: {str(e)}")
                    await asyncio.sleep(5)  # Backoff on Kafka errors
                except Exception as e:
                    logger.error(f"Unexpected error: {str(e)}")
                    await asyncio.sleep(5)  # Backoff on critical errors
        except Exception as e:
            logger.error(f"Consumer crashed: {str(e)}. Restarting in 10 seconds...")
            if consumer:
                consumer.close()
            await asyncio.sleep(10)  # Wait before restarting
        finally:
            if consumer:
                consumer.close()


def handle_shutdown(signal, frame):
    """
    Handle shutdown signals (e.g., SIGTERM) gracefully.
    """
    logger.info("Received shutdown signal, shutting down gracefully...")
    exit(0)

def start_worker():
    """
    Start the Kafka worker to consume and process messages.
    """
    # Register shutdown handler
    signal.signal(signal.SIGTERM, handle_shutdown)
    signal.signal(signal.SIGINT, handle_shutdown)

    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(consume_messages())
    except KeyboardInterrupt:
        logger.info("Shutting down gracefully...")
    finally:
        loop.close()
        executor.shutdown(wait=True)

if __name__ == "__main__":
    start_worker()