package handler

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type createPostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var request createPostRequest

	err := jsoniter.ConfigFastest.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.log.Error("failed to unmarshal request", zap.Error(err))
		h.writeApiError(w, http.StatusBadRequest, LevelUser, "failed to unmarshal json")

		return
	}

	err = h.validator.StructCtx(r.Context(), request)
	if err != nil {
		h.log.Info("validation error", zap.String("error", err.Error()))
		h.writeValidationErr(w, err)

		return
	}

	p, err := h.postService.CreatePost(request.Title, request.Content)
	if err != nil {
		h.log.Error("failed to create post", zap.Error(err))
		h.writeError(w, err, err.Error())

		return
	}

	h.writeResponse(w, http.StatusOK, p)
}
