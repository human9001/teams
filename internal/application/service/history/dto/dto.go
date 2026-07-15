package dto

import "github.com/human9001/teams/internal/domain/history/model"

type ListHistoryResult struct {
	Items []model.HistoryItem `json:"items"`
}
