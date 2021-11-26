package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *Handler) postByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	p, err := h.postService.GetPost(id)
	if err != nil {
		h.log.Error("failed to get post", zap.Error(err))
		h.writeError(w, err, err.Error())

		return
	}

	h.writeResponse(w, http.StatusOK, p)
}
