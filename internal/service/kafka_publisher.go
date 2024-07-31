package service

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/zsandibe/messaggio-microservice/internal/storage"
)

type kafkaPublisher struct {
	storage *storage.KafkaStorage
}

func NewKafkaPublisher(storage *storage.KafkaStorage) *kafkaPublisher {
	return &kafkaPublisher{storage: storage}
}

func (kp *kafkaPublisher) PublishMessage(ctx context.Context, key, value []byte, id string) error {
	err := kp.storage.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Headers: []kafka.Header{
			{Key: id, Value: value},
		},
	})
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}
	return nil
}
