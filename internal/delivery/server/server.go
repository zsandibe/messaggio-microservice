package server

import (
	"context"
	"net/http"
	"time"

	"github.com/zsandibe/messaggio-microservice/config"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

type Server struct {
	httpServer http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:           cfg.Server.Host + ":" + cfg.Server.Port,
			Handler:        handler,
			MaxHeaderBytes: 1024 * 1024,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	logger.Info("Starting server on: ", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Close(); err != nil {
		return err
	}
	return nil
}
