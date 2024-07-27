package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

type Kafka struct {
	Broker  []string `envconfig:"KAFKA_BROKER" required:"true"`
	Topic   string   `envconfig:"KAFKA_TOPIC" required:"true"`
	GroupId string   `envconfig:"KAFKA_GROUP_ID" required:"true"`
}

type Server struct {
	Host string `envconfig:"SERVER_HOST" required:"true"`
	Port string `envconfig:"SERVER_PORT" required:"true"`
}

type Postgres struct {
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

type Config struct {
	Postgres Postgres
	Server   Server
	Kafka    Kafka
}

func NewConfig(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		logger.Error("godotenv.Load():%v", err)
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	var config Config

	if err := envconfig.Process("", &config); err != nil {
		logger.Error("envconfig.Process() :%v", err)
		return nil, fmt.Errorf("error processing .env file: %v", err)
	}
	return &config, nil
}
