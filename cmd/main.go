package main

import (
	"github.com/zsandibe/messaggio-microservice/internal/app"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

// @title Messaggio test task
// @version 1.0
// @description This is basic server for a message sending
// @host 0.0.0.0:7777
// @BasePath /api/v1
func main() {
	if err := app.Start(); err != nil {
		logger.Error(err)
		return
	}
}
