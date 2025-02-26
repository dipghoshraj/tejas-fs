package monitor

import (
	"github.com/IBM/sarama"
)

// Producer is a wrapper around sarama.SyncProducer to provide additional functionality
type Producer struct {
	sarama.SyncProducer
}

// NewProducer creates a new Producer
func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer}, nil
}
