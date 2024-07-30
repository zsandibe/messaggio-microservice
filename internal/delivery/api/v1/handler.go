package v1

import (
	"github.com/segmentio/kafka-go"
	"github.com/zsandibe/messaggio-microservice/internal/service"
)

type Handler struct {
	service *service.Service
	w       *kafka.Writer
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
