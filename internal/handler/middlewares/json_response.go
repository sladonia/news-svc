package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewJsonResponse() *JsonResponseMiddleware {
	return &JsonResponseMiddleware{}
}

type JsonResponseMiddleware struct{}

func (m *JsonResponseMiddleware) Register(r *mux.Router) {
	r.Use(m.jsonResponse)
}

func (m *JsonResponseMiddleware) jsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
