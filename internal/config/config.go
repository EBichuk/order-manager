package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Cache
	Kafka
	Db
	HttpServer
}

type Cache struct {
	Size int `env:"CACHE_SIZE" env-default:"100"`
}

type Kafka struct {
	Topic   string `env:"KAFKA_TOPIC" env-default:"order"`
	Brokers string `env:"KAFKA_BROKERS" env-default:"broker-1:19092,broker-2:19092,broker-3:19092"`
}

type Db struct {
	Name     string `env:"POSTGRES_NAME"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Ssl      string `env:"POSTGRES_SSL"`
}

type HttpServer struct {
	Addr string `env:"HTTP_ADDRESS" env-default:"localhost:8081"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return cfg, fmt.Errorf("failed to read env: %w", err)
		}
		return cfg, nil
	}
	return cfg, nil
}
