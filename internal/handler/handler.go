package handler

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/sladonia/news-svc/internal/post"
	"go.uber.org/zap"
)

func NewHandler(
	log *zap.Logger,
	defaultNewsLimit uint,
	postService post.Service,
	serviceName string,
) *Handler {
	return &Handler{
		log:              log,
		validator:        validator.New(),
		postService:      postService,
		defaultNewsLimit: defaultNewsLimit,
		serviceName:      serviceName,
	}
}

type Handler struct {
	serviceName      string
	log              *zap.Logger
	postService      post.Service
	validator        *validator.Validate
	defaultNewsLimit uint
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/", h.identity)
	r.HandleFunc("/posts", h.createPost).Queries(
		"limit", "{limit}",
		"offset", "{offset}",
		"from", "{from}",
		"to", "{to}",
	).Name("createPost").Methods("POST")
	r.HandleFunc("/posts", h.findPosts).Name("findPosts").Methods("GET")
	r.HandleFunc("/posts/{id}", h.postByID).Name("postByID").Methods("GET")
	r.HandleFunc("/posts/{id}", h.replacePost).Name("replacePost").Methods("PUT")
	r.HandleFunc("/posts/{id}", h.deletePost).Name("deletePost").Methods("DELETE")
}

func (h *Handler) identity(w http.ResponseWriter, r *http.Request) {
	encoded, _ := jsoniter.ConfigFastest.Marshal(struct {
		ServiceName string `json:"service_name"`
	}{ServiceName: h.serviceName})

	w.WriteHeader(200)
	w.Write(encoded)
}

func (h *Handler) writeApiError(w http.ResponseWriter, status int, level Level, msg string) {
	apiErr := NewApiError(msg, level)

	h.writeResponse(w, status, apiErr)
}

func (h *Handler) writeValidationErr(w http.ResponseWriter, err error) {
	var vErr validator.ValidationErrors
	if !errors.As(err, &vErr) {
		h.writeError(w, err, err.Error())
		return
	}

	msg := "validation error:"

	for _, err := range vErr {
		msg += fmt.Sprintf(" %s %s", strings.ToLower(err.StructField()), err.ActualTag())
	}

	apiErr := NewApiError(msg, LevelUser)

	h.writeResponse(w, http.StatusBadRequest, apiErr)
}

func (h *Handler) writeError(w http.ResponseWriter, err error, msg string) {
	status, level := h.classifyError(err)

	apiErr := NewApiError(msg, level)

	h.writeResponse(w, status, apiErr)
}

func (h *Handler) writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)

	if data == nil {
		return
	}

	encoded, err := jsoniter.ConfigFastest.Marshal(data)
	if err != nil {
		h.log.Error("failed to marshal response", zap.Error(err))
		return
	}

	_, err = w.Write(encoded)
	if err != nil {
		h.log.Error("failed to write response", zap.Error(err))
	}
}

func (h *Handler) classifyError(err error) (int, Level) {
	var operationError *net.OpError

	switch {
	case errors.As(err, &operationError):
		return http.StatusInternalServerError, LevelSystem
	case errors.Is(err, post.ErrNotFound):
		return http.StatusNotFound, LevelUser
	case errors.Is(err, post.ErrorAlreadyExists):
		return http.StatusConflict, LevelUser
	default:
		return http.StatusInternalServerError, LevelSystem
	}
}
