package service

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/zsandibe/messaggio-microservice/internal/storage"
)

type kafkaService struct {
	storage *storage.KafkaStorage
}

func NewKafkaService(storage *storage.KafkaStorage) *kafkaService {
	return &kafkaService{storage: storage}
}

func (ks *kafkaService) PublishMessage(ctx context.Context, key, value []byte) error {
	err := ks.storage.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}
	return nil
}

func (ks *kafkaService) ConsumeMessages(ctx context.Context) {
	for {
		msg, err := ks.storage.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message from Kafka: %v", err)
			continue
		}
		log.Printf("Received message: key=%s, value=%s", string(msg.Key), string(msg.Value))
	}
}
