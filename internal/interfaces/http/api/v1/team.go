package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/human9001/teams/internal/interfaces/http/api/v1/dto"
	apiErrors "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/pkg/helpers"
)

func (a *API) CreateTeam(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateTeam handler")
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrUnauthorized)
		return
	}
	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return
	}
	if req.Name == "" {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrMissingRequiredFields)
		return
	}

	t, err := a.teamService.CreateTeam(r.Context(), req.Name, claims.UserID)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, dto.TeamResponse{
		ID:      t.ID(),
		Name:    t.Name(),
		OwnerID: t.OwnerID(),
	})
}

func (a *API) ListTeams(w http.ResponseWriter, r *http.Request) {
	slog.Info("ListTeams handler")
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrUnauthorized)
		return
	}
	var req dto.ListTeamsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return
	}

	if req.OwnerID == 0 {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrMissingRequiredFields)
		return
	}

	t, err := a.teamService.ListTeams(r.Context(), claims.UserID)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	names := make([]string, len(t))
	for i, team := range t {
		names[i] = team.Name()
	}
	helpers.WriteJSON(w, http.StatusCreated, dto.ListTeamsResponse{
		Name: names,
	})
}

func (a *API) InviteUser(w http.ResponseWriter, r *http.Request) {
	slog.Info("InviteUser handler")
}
