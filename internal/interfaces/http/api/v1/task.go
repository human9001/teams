package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel"

	"github.com/human9001/teams/internal/application/service/task/input"
	domainErrs "github.com/human9001/teams/internal/domain/errors"
	"github.com/human9001/teams/internal/domain/task/model"
	"github.com/human9001/teams/internal/interfaces/http/api/v1/dto"
	apiErrors "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/pkg/helpers"
)

func (a *API) CreateTask(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateTask handler")
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrUnauthorized)
		return
	}
	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return
	}

	priority, err := model.NewPriority(strings.ToUpper(req.Priority))
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return

	}
	in := input.CreateTaskInput{
		TeamID:      req.TeamID,
		AssigneeID:  req.AssigneeID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
	}
	t, err := a.taskService.CreateTask(r.Context(), claims.UserID, in)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, dto.TaskResponse{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Priority:    string(t.Priority()),
	})
}

func (a *API) ListTasks(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("teams-service")
	ctx, span := tracer.Start(r.Context(), "health-check")
	defer span.End()

	slog.InfoContext(ctx, "health endpoint called")
	q := r.URL.Query()

	teamID, err := strconv.ParseInt(q.Get("team_id"), 10, 64)
	if err != nil || teamID == 0 {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrUnknownTeamID)
		return
	}

	var status *model.Status
	if s := q.Get("status"); s != "" {
		v := model.Status(strings.ToUpper(s))
		status = &v
	}

	var assigneeID *int64
	if v := q.Get("assignee_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidAssigneeID)
			return
		}
		assigneeID = &id
	}

	page, err := strconv.ParseInt(q.Get("page"), 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "parse page param", "error", err)
	}
	limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "parse limit param", "error", err)
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}

	res, err := a.taskService.ListTasks(r.Context(), input.ListTasksInput{
		TeamID:     teamID,
		Status:     status,
		AssigneeID: assigneeID,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	items := make([]dto.TaskItem, 0, len(res.Items))
	for _, t := range res.Items {
		items = append(items, dto.TaskItem{
			ID:          t.ID,
			TeamID:      t.TeamID,
			AssigneeID:  &t.AssigneeID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			Priority:    t.Priority,
		})
	}

	helpers.WriteJSON(w, http.StatusOK, dto.ListTasksResponse{Items: items, Total: res.Total, Page: res.Page, Limit: res.Limit})
}

func (a *API) UpdateTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidTaskID)
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return
	}

	in := input.UpdateTaskInput{
		TaskID:      id,
		UserID:      claims.UserID,
		Description: helpers.PtrOrNIl(req.Description),
		AssigneeID:  helpers.PtrOrNIl(req.AssigneeID),
		Status:      helpers.PtrOrNIl(new(model.Status(strings.ToUpper(*req.Status)))),
		Priority:    helpers.PtrOrNIl(new(model.Priority(strings.ToUpper(*req.Priority)))),
	}

	if req.Title != nil {
		v := strings.TrimSpace(*req.Title)
		if v == "" {
			helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrMissingRequiredFields)
		}
		in.Title = &v
	}

	err = a.taskService.UpdateTask(r.Context(), in)
	if err != nil {
		switch {
		case errors.Is(err, domainErrs.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domainErrs.ErrForbidden):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusOK, nil)
}
