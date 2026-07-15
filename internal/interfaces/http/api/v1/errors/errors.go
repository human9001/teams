package apiErrors

import "errors"

var (
	ErrInvalidJSON           = errors.New("invalid JSON")
	ErrNotFound              = errors.New("not found")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrBearerToken           = errors.New("missing bearer token")
	ErrInvalidToken          = errors.New("invalid token")
	ErrInvalidTokenClaims    = errors.New("invalid token claims")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrUnknownTeamID         = errors.New("team_id is required")
	ErrInvalidTaskID         = errors.New("invalid task id")
	ErrInvalidAssigneeID     = errors.New("invalid assignee_id")
)
