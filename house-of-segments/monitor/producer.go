package monitor

import (
	"github.com/IBM/sarama"
)

const TOPIC = "metadata"

// SendMessage sends a message to the specified topic
func (p *Producer) SendMessage() error {
	_, _, err := p.SendMessageWithKey(TOPIC, "", "NODEID:123,STORAGE:1000")
	return err
}

// SendMessageWithKey sends a message to the specified topic with the specified key
func (p *Producer) SendMessageWithKey(topic string, key string, message string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.SyncProducer.SendMessage(msg)
	return partition, offset, err
}
