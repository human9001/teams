package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/human9001/teams/internal/application/service/comment/dto"
	domainErrors "github.com/human9001/teams/internal/domain/errors"
	apiErrors "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/pkg/helpers"
)

func (a *API) ListComments(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || taskID == 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	page, err := strconv.ParseInt(q.Get("page"), 10, 64)
	if err != nil {
		slog.Error("parse page param", "error", err)
	}
	limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)
	if err != nil {
		slog.Error("parse limit param", "error", err)
	}

	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrUnauthorized)
		return
	}
	res, err := a.commentService.ListComments(r.Context(), dto.CommentListInput{
		TaskID: taskID,
		UserID: claims.UserID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrForbidden):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, domainErrors.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domainErrors.ErrInvalidLimit):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("list comments encode ", "error", err)
	}
}
