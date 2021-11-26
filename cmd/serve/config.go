package main

import (
	"time"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type httpConfig struct {
	Listen          string        `env:"HTTP_LISTEN" default:"8080" json:"listen"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" default:"5s" json:"read_timeout"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" default:"5s" json:"write_timeout"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" default:"5s" json:"shutdown_timeout"`
}

type Config struct {
	HTTP             httpConfig
	ServiceName      string `env:"SERVICE_NAME" default:"news-svc" json:"service_name"`
	LogLevel         string `env:"LOG_LEVEL" default:"info" json:"log_level"`
	PostgresDSN      string `env:"POSTGRES_DSN" default:"postgresql://user:password@localhost:5432/news?sslmode=disable" json:"-"`
	PostTableName    string `env:"POST_TABLE_NAME" default:"post" json:"post_table_name"`
	DefaultNewsLimit uint   `env:"DEFAULT_NEWS_LIMIT" default:"100" json:"default_news_limit"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	var config Config

	err := configor.Load(&config)

	return config, err
}
