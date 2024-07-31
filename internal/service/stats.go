package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/internal/entity"
	"github.com/zsandibe/messaggio-microservice/internal/repository"
)

type statisticService struct {
	statisticRepo repository.Statistic
}

func NewStatisticService(statisticRepo repository.Statistic) *statisticService {
	return &statisticService{statisticRepo: statisticRepo}
}

func (s *statisticService) GetStatsList(ctx context.Context) ([]*entity.Stats, error) {
	return s.statisticRepo.GetStatsList(ctx)
}

func (s *statisticService) GetStatById(ctx context.Context, id int) (*entity.Stats, error) {
	stat, err := s.statisticRepo.GetStatById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.Stats{}, domain.ErrStatisticNotFound
		}
		return nil, err
	}
	return stat, nil
}
