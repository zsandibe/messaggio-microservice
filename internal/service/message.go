package service

import (
	"context"
	"database/sql"
	"errors"

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

func (s *messageService) DeleteMessageById(ctx context.Context, id int) error {
	logger.Debugf("Delete user by id: %d ", id)
	if err := s.messageRepo.DeleteMessageById(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrMessageNotFound
		}
		return err
	}
	logger.Info("User was successfully updated")
	return nil
}

func (s *messageService) GetMessageById(ctx context.Context, id int) (*entity.Message, error) {
	logger.Debugf("Get message by id: %d", id)
	message, err := s.messageRepo.GetMessageById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.Message{}, domain.ErrMessageNotFound
		}
		return &entity.Message{}, err
	}
	logger.Debug(message)
	logger.Info("Message successfully got")
	return message, nil
}

func (s *messageService) GetMessagesList(ctx context.Context, params domain.MessagesListParams) ([]*entity.Message, error) {
	logger.Debugf("Get messages list: %v ", params)
	return s.messageRepo.GetMessagesList(ctx, params)
}

func (s *messageService) UpdateStatus(ctx context.Context, id int) error {
	return s.messageRepo.UpdateStatus(ctx, id)
}
