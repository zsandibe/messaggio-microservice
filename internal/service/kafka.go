package service

import (
	"context"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/zsandibe/messaggio-microservice/internal/storage"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

type kafkaService struct {
	messageService messageService
	storage        *storage.KafkaStorage
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
		if err := ks.CommitUpdatedMessages(ctx, msg); err != nil {
			return
		}

	}
}

func (ks *kafkaService) CommitUpdatedMessages(ctx context.Context, msg kafka.Message) error {
	id, err := strconv.Atoi(msg.Headers[0].Key)
	if err != nil {
		return err
	}

	if err := ks.messageService.UpdateStatus(ctx, id); err != nil {
		logger.Error("Failed to update message status: %v", err)
		return err
	}

	if err := ks.storage.Reader.CommitMessages(ctx, msg); err != nil {
		logger.Error("Failed to commit messages: %v", err)
		return err
	}
	return nil
}
