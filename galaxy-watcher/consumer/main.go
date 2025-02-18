package main

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"}, // Kafka broker address
		Topic:       "my_topic",                 // The topic to consume from
		GroupID:     "my-consumer-group",        // Consumer group ID
		StartOffset: kafka.FirstOffset,
	})

	for {
		msg, err := consumer.ReadMessage(nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("message %w", msg)
	}
}
