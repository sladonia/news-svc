package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func (h *Handler) replacePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

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

	err = h.postService.UpsertPost(id, request.Title, request.Content)
	if err != nil {
		h.log.Error("failed to upsert post", zap.Error(err))
		h.writeError(w, err, err.Error())

		return
	}

	h.writeResponse(w, http.StatusOK, nil)
}
