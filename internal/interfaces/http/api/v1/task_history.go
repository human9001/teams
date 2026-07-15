package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	domainErrs "github.com/human9001/teams/internal/domain/errors"
)

func (a *API) TaskHistory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || taskID == 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	res, err := a.taskHistoryService.GetTaskHistory(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, domainErrs.ErrNotFound) {
			http.Error(w, "history not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("task history", "error", err)
	}
}
