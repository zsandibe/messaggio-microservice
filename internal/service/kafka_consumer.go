package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/internal/storage"
)

type kafkaConsumer struct {
	messageService messageService
	storage        *storage.KafkaStorage
}

func NewKafkaConsumer(messageService messageService, storage *storage.KafkaStorage) *kafkaConsumer {
	return &kafkaConsumer{messageService: messageService, storage: storage}
}

func (kc *kafkaConsumer) ConsumeMessages(ctx context.Context) {
	go func() {
		for {
			msg, err := kc.storage.Reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Failed to read message from Kafka: %v", err)
				continue
			}

			log.Printf("Received message: key=%s, value=%s", string(msg.Key), string(msg.Value))
			// msg.Headers[0].Key = id
			if err = kc.CommitUpdatedMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit updated messages: %v", err)
				continue
			}
		}
	}()
}

func (kc *kafkaConsumer) CommitUpdatedMessages(ctx context.Context, msg kafka.Message) error {
	if kc == nil {
		return errors.New("messageService is nil")
	}

	if len(msg.Key) == 0 {
		return errors.New("message headers are empty")
	}

	id, err := strconv.Atoi(msg.Headers[0].Key)
	if err != nil {
		return err
	}

	var message domain.CreateMessageRequest
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message value: %w", err)
	}

	if err := kc.messageService.UpdateStatus(ctx, id, message.Content); err != nil {
		log.Printf("Failed to update message status: %v", err)
		return err
	}

	if err := kc.storage.Reader.CommitMessages(ctx, msg); err != nil {
		log.Printf("Failed to commit messages: %v", err)
		return err
	}
	return nil
}
