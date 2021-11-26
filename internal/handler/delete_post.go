package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := h.postService.DeletePost(id)
	if err != nil {
		h.log.Error("failed to delete post", zap.Error(err))
		h.writeError(w, err, err.Error())

		return
	}

	h.writeResponse(w, http.StatusNoContent, nil)
}
