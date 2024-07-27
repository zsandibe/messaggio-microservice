package service

import "github.com/zsandibe/messaggio-microservice/internal/repository"

type statisticService struct {
	statisticRepo repository.Statistic
}

func NewStatisticService(statisticRepo repository.Statistic) *statisticService {
	return &statisticService{statisticRepo: statisticRepo}
}
