package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaClient struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

func NewKafkaProducer() (*kafka.Producer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	}
	return kafka.NewProducer(&config)
}

func (kc *KafkaClient) Publish(topic string, message interface{}) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	kc.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msgBytes,
	}, nil)
	return nil
}

func NewKafkaConsumer(brokers []string, groupID string) (*kafka.Consumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}
	return kafka.NewConsumer(&config)
}

func (kc *KafkaClient) Subscribe(topic string, handler func(message kafka.Message)) error {
	err := kc.Consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	go func() {
		for {
			msg, err := kc.Consumer.ReadMessage(-1)
			if err == nil {
				handler(*msg)
			}
		}
	}()
	return nil
}
