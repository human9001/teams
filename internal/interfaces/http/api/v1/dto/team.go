package dto

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type TeamResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	OwnerID int64  `json:"owner_id"`
}

type ListTeamsRequest struct {
	OwnerID int64 `json:"owner_id"`
}

type ListTeamsResponse struct {
	Name []string `json:"names"`
}
