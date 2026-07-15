package converter

import (
	"github.com/human9001/teams/internal/application/service/task/input"
	"github.com/human9001/teams/internal/domain/task/model"
)

func FromModelToListTasksItem(items []model.Task) []input.ListTasksItem {
	res := make([]input.ListTasksItem, 0, len(items))
	for _, v := range items {
		item := input.ListTasksItem{
			ID:          v.ID(),
			TeamID:      v.TeamID(),
			CreatedBy:   v.CreatedBy(),
			AssigneeID:  *v.AssigneeID(),
			Title:       v.Title(),
			Description: v.Description(),
			Status:      v.Status().String(),
			Priority:    v.Priority().String(),
			CreatedAt:   v.CreatedAt(),
			UpdatedAt:   v.UpdatedAt(),
		}

		res = append(res, item)

	}
	return res
}
