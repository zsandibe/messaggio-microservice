package main

import (
	"github.com/zsandibe/messaggio-microservice/internal/app"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

func main() {
	if err := app.Start(); err != nil {
		logger.Error(err)
		return
	}
}
