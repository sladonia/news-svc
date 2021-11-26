package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sladonia/news-svc/internal/handler"
	"github.com/sladonia/news-svc/internal/handler/middlewares"
	"github.com/sladonia/news-svc/internal/logger"
	"github.com/sladonia/news-svc/internal/post"
	"github.com/sladonia/news-svc/internal/poststorage"
	"go.uber.org/zap"
)

func mustInitZapLogger(logLevel string) *zap.Logger {
	log, err := logger.NewZapLogger(logLevel)
	if err != nil {
		panic(fmt.Sprintf("init logger: %s", err))
	}

	return log
}

func mustLoadConfig() Config {
	config, err := LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("load config: %s", err))
	}

	return config
}

func mustCreateDatabaseConnection(config Config, log *zap.Logger) *goqu.Database {
	postgresClient, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		log.Panic("create database", zap.Error(err))
	}

	return goqu.New("postgres", postgresClient)
}

func newPostStorage(config Config, db *goqu.Database) post.Storage {
	return poststorage.New(db, config.PostTableName)
}

func newHandler(config Config, log *zap.Logger, postService post.Service) *handler.Handler {
	return handler.NewHandler(log, config.DefaultNewsLimit, postService, config.ServiceName)
}

func createHTTPServer(config Config, router http.Handler) *http.Server {
	return &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", config.HTTP.Listen),
		ReadTimeout:  config.HTTP.ReadTimeout,
		WriteTimeout: config.HTTP.WriteTimeout,
	}
}

func registerHTTPHandlers(log *zap.Logger, r *mux.Router, handler *handler.Handler) {
	middlewares.NewHandlerLogger(log).Register(r)
	middlewares.NewJsonResponse().Register(r)

	handler.Register(r)
}
