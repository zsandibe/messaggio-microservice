package service

import (
	"context"

	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/internal/entity"
	"github.com/zsandibe/messaggio-microservice/internal/repository"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

type messageService struct {
	messageRepo repository.Message
}

func NewMessageService(messageRepo repository.Message) *messageService {
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) CreateMessage(ctx context.Context, inp domain.CreateMessageRequest) (*entity.Message, error) {
	logger.Debugf("Creating message: %+v", inp)

	exists, _ := s.messageRepo.IsMessageExist(ctx, inp)
	if !exists {
	}

	message, err := s.messageRepo.CreateMessage(ctx, inp)
	if err != nil {
		return &entity.Message{}, domain.ErrCreatingMessage
	}
	logger.Debug(message)
	logger.Info("Message successfully created")
	return message, nil
}
