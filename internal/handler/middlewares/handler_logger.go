package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewHandlerLogger(logger *zap.Logger) *HandlerLogger {
	return &HandlerLogger{logger: logger}
}

type HandlerLogger struct {
	logger *zap.Logger
}

func (mw *HandlerLogger) Register(r *mux.Router) {
	r.Use(mw.log)
}

func (mw *HandlerLogger) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw.logger.Debug(
			"incoming request",
			zap.String("url", r.URL.String()),
			zap.String("method", r.Method),
		)

		next.ServeHTTP(w, r)
	})
}
