package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sladonia/news-svc/internal/post"
	"go.uber.org/zap"
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		config    = mustLoadConfig()
		log       = mustInitZapLogger(config.LogLevel)
	)

	defer log.Sync()
	defer stop()

	log.Info("config loaded", zap.Any("config", config))

	var (
		db          = mustCreateDatabaseConnection(config, log)
		postStorage = newPostStorage(config, db)
		postService = post.NewService(postStorage)
		handler     = newHandler(config, log, postService)
		router      = mux.NewRouter()
		server      = createHTTPServer(config, router)
	)

	registerHTTPHandlers(log, router, handler)
	run(ctx, config, log, server, stop)
}

func run(ctx context.Context, config Config, log *zap.Logger, srv *http.Server, stop func()) {
	errCh := make(chan error)

	go func() {
		errCh <- srv.ListenAndServe()
	}()

	shutdown := func(err error) {
		stop()

		timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), config.HTTP.ShutdownTimeout)
		defer cancelTimeout()

		if err := srv.Shutdown(timeoutCtx); err != nil {
			log.Error("server shutdown", zap.Error(err))
		}

		if err != nil {
			log.Error("shutdown caused by error", zap.Error(err))
			os.Exit(1)
		}
	}

	select {
	case <-ctx.Done():
		shutdown(nil)
	case err := <-errCh:
		shutdown(err)
	}
}
