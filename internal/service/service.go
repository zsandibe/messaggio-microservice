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
	DeleteMessageById(ctx context.Context, id int) error
	GetMessagesList(ctx context.Context, params domain.MessagesListParams) ([]*entity.Message, error)
	GetMessageById(ctx context.Context, id int) (*entity.Message, error)
	UpdateStatus(ctx context.Context, id int, content string) error
}

type Statistic interface {
	GetStatsList(ctx context.Context) ([]*entity.Stats, error)
	GetStatById(ctx context.Context, id int) (*entity.Stats, error)
}

type Publisher interface {
	PublishMessage(ctx context.Context, key, value []byte, id string) error
}

type Consumer interface {
	ConsumeMessages(ctx context.Context)
}

type Service struct {
	Message
	Statistic
	Publisher
	Consumer
}

func NewService(repo *repository.Repository, storage *storage.KafkaStorage) *Service {
	return &Service{
		Message:   NewMessageService(repo.Message),
		Statistic: NewStatisticService(repo.Statistic),
		Publisher: NewKafkaPublisher(storage),
		Consumer:  NewKafkaConsumer(*NewMessageService(repo.Message), storage),
	}
}
