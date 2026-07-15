package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/human9001/teams/internal/interfaces/http/api/v1"
	apiErrs "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/internal/interfaces/http/api/v1/mocks"
)

func TestLoginFailed(t *testing.T) {
	teamSrv := mocks.NewITeamService(t)
	userSrv := mocks.NewIUserService(t)
	taskSrv := mocks.NewITaskService(t)
	historySrv := mocks.NewITaskHistoryService(t)
	commentSrv := mocks.NewICommentService(t)
	api := apiv1.NewAPI(teamSrv, userSrv, taskSrv, historySrv, commentSrv)

	userSrv.EXPECT().
		Login(mock.Anything, mock.Anything).
		Return(nil, apiErrs.ErrUnauthorized)

	body := []byte(`{"email":"user@example.com","password":"secret"}`)
	ctx := context.Background()
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Login(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLoginInvalidJson(t *testing.T) {
	teamSrv := mocks.NewITeamService(t)
	userSrv := mocks.NewIUserService(t)
	taskSrv := mocks.NewITaskService(t)
	historySrv := mocks.NewITaskHistoryService(t)
	commentSrv := mocks.NewICommentService(t)
	api := apiv1.NewAPI(teamSrv, userSrv, taskSrv, historySrv, commentSrv)

	body := []byte(`{bad json}`)
	ctx := context.Background()
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Login(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestLoginMissingPassword(t *testing.T) {
	teamSrv := mocks.NewITeamService(t)
	userSrv := mocks.NewIUserService(t)
	taskSrv := mocks.NewITaskService(t)
	historySrv := mocks.NewITaskHistoryService(t)
	commentSrv := mocks.NewICommentService(t)
	api := apiv1.NewAPI(teamSrv, userSrv, taskSrv, historySrv, commentSrv)

	body := []byte(`{"email":"user@example.com"}`)
	ctx := context.Background()
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	api.Login(rr, req)
	var resp ErrorResponse
	_ = json.NewDecoder(rr.Body).Decode(&resp)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, "missing required fields", resp.Error)
}
