package config

import "github.com/ilyakaznacheev/cleanenv"

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
	Topic   string `env:"KAFKA_TOPIC"`
	Brokers string `env:"KAFKA_BROKERS"`
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
	Addr string `env:"HTTP_ADDRESS"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
