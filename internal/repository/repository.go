package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/internal/entity"
)

type Message interface {
	CreateMessage(ctx context.Context, msg domain.CreateMessageRequest) (*entity.Message, error)
	IsMessageExist(ctx context.Context, msg domain.CreateMessageRequest) (bool, error)
	DeleteMessageById(ctx context.Context, id int) error
	GetMessagesList(ctx context.Context, inp domain.MessagesListParams) ([]*entity.Message, error)
	GetMessageById(ctx context.Context, id int) (*entity.Message, error)
}

type Statistic interface{}

type Repository struct {
	Message
	Statistic
}

func NewPostgresRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Message:   NewMessageRepo(db),
		Statistic: NewStatisticRepo(db),
	}
}
