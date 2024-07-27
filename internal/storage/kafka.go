package storage

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zsandibe/messaggio-microservice/config"
)

type KafkaStorage struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func NewKafkaStorage(cfg *config.Config) *KafkaStorage {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Kafka.Broker,
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Broker,
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupId,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  10 * time.Second,
	})

	return &KafkaStorage{
		Writer: writer,
		Reader: reader,
	}
}

func (ks *KafkaStorage) Close() error {
	if err := ks.Writer.Close(); err != nil {
		log.Printf("Failed to close Kafka writer: %v", err)
		return err
	}

	if err := ks.Reader.Close(); err != nil {
		log.Printf("Failed to close Kafka reader: %v", err)
		return err
	}

	return nil
}
