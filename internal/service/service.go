package service

import (
	"context"

	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/internal/entity"
	"github.com/zsandibe/messaggio-microservice/internal/repository"
	"github.com/zsandibe/messaggio-microservice/internal/storage"
)

type Message interface {
	CreateMessage(ctx context.Context, msg domain.CreateMessageRequest) (*entity.Message, error)
}

type Statistic interface{}

type Kafka interface {
	PublishMessage(ctx context.Context, key, value []byte) error
	ConsumeMessages(ctx context.Context)
}

type Service struct {
	Message
	Statistic
	Kafka
}

func NewService(repo repository.Repository, storage *storage.KafkaStorage) *Service {
	return &Service{
		Message:   NewMessageService(repo.Message),
		Statistic: NewStatisticService(repo.Statistic),
		Kafka:     NewKafkaService(storage),
	}
}
