package v1

import (
	"encoding/json"
	"net/http"

	"github.com/human9001/teams/internal/application/service/user"
	"github.com/human9001/teams/internal/interfaces/http/api/v1/dto"
	apiErrors "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/pkg/helpers"
)

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return
	}

	u, err := a.authService.Register(r.Context(), user.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, dto.AuthResponse{UserID: u.ID()})
}
