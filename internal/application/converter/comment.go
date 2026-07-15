package converter

import (
	"github.com/human9001/teams/internal/application/service/comment/dto"
	"github.com/human9001/teams/internal/domain/comment/model"
)

func FromModelToCommentListItem(items []model.Comment) []dto.CommentItem {
	res := make([]dto.CommentItem, 0, len(items))
	for _, v := range items {
		item := dto.CommentItem{
			ID:        v.ID,
			TaskID:    v.TaskID,
			AuthorID:  v.AuthorID,
			Body:      v.Body,
			CreatedAt: v.CreatedAt,
		}
		res = append(res, item)
	}
	return res
}
