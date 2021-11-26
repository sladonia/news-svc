package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sladonia/news-svc/internal/post"
	"go.uber.org/zap"
)

func (h *Handler) findPosts(w http.ResponseWriter, r *http.Request) {
	limit, err := h.parseUint(r.FormValue("limit"))
	if err != nil {
		h.log.Info("atoi error. limit", zap.Error(err))
		h.writeApiError(w, http.StatusBadRequest, LevelUser, "limit query parameter should be integer")

		return
	}

	offset, err := h.parseUint(r.FormValue("offset"))
	if err != nil {
		h.log.Info("atoi error. offset", zap.Error(err))
		h.writeApiError(w, http.StatusBadRequest, LevelUser, "offset query parameter should be integer")

		return
	}

	from, err := h.parseTime(r.FormValue("from"))
	if err != nil {
		h.log.Info("time parse error. offset", zap.Error(err))
		h.writeApiError(w, http.StatusBadRequest, LevelUser, "from query parameter should be RFC3339 formatted")

		return
	}

	to, err := h.parseTime(r.FormValue("to"))
	if err != nil {
		h.log.Info("time parse error. offset", zap.Error(err))
		h.writeApiError(w, http.StatusBadRequest, LevelUser, "to query parameter should be RFC3339 formatted")

		return
	}

	if limit == 0 {
		limit = h.defaultNewsLimit
	}

	f := post.Filter{
		From:   from,
		To:     to,
		Limit:  limit,
		Offset: offset,
	}

	posts, err := h.postService.FindPosts(f)
	if err != nil {
		h.log.Error("failed to delete post", zap.Error(err))
		h.writeError(w, err, err.Error())

		return
	}

	h.writeResponse(w, http.StatusOK, posts)
}

func (h *Handler) parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, timeStr)
}

func (h *Handler) parseUint(intStr string) (uint, error) {
	if intStr == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(intStr)
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}
