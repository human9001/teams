package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/human9001/teams/internal/application/service/user"
	"github.com/human9001/teams/internal/interfaces/http/api/v1/dto"
	apiErrors "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/pkg/helpers"
)

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrInvalidJSON)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		helpers.WriteError(w, http.StatusBadRequest, apiErrors.ErrMissingRequiredFields)
		return
	}
	res, err := a.authService.Login(r.Context(), user.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		UserID: res.UserID,
		Token:  res.Token,
	})
}
