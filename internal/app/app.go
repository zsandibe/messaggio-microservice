package app

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/zsandibe/messaggio-microservice/config"
	"github.com/zsandibe/messaggio-microservice/internal/storage"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

func Start() error {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		return fmt.Errorf("config.NewConfig: %v", err)
	}
	logger.Info("Config loaded successfully")

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		return fmt.Errorf("storage.NewPostgresDB: %v", err)
	}
	defer db.Close()
	logger.Info("Database  loaded successfully")

	if err = db.MigrateUp(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Debug(err)
	}

	kafka := storage.NewKafkaStorage(cfg)
	defer kafka.Close()

	return nil
}
