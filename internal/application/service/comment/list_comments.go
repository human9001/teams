package comment

import (
	"context"

	"github.com/human9001/teams/internal/application/converter"
	"github.com/human9001/teams/internal/application/service/comment/dto"
	domainErrors "github.com/human9001/teams/internal/domain/errors"
)

func (s CommentService) ListComments(ctx context.Context, in dto.CommentListInput) (dto.CommentListResult, error) {
	if in.TaskID == 0 {
		return dto.CommentListResult{}, domainErrors.ErrNotFound
	}
	if in.UserID == 0 {
		return dto.CommentListResult{}, domainErrors.ErrForbidden
	}
	if in.Page == 0 {
		in.Page = 1
	}
	if in.Limit == 0 {
		in.Limit = 20
	}
	if in.Limit > 500 {
		return dto.CommentListResult{}, domainErrors.ErrInvalidLimit
	}

	ok, err := s.repo.UserHasAccess(ctx, in.TaskID, in.UserID)
	if err != nil {
		return dto.CommentListResult{}, err
	}
	if !ok {
		return dto.CommentListResult{}, domainErrors.ErrForbidden
	}

	offset := (in.Page - 1) * in.Limit
	items, total, err := s.repo.ListByTaskID(ctx, in.TaskID, in.Limit, offset)
	if err != nil {
		return dto.CommentListResult{}, err
	}

	return dto.CommentListResult{
		Items: converter.FromModelToCommentListItem(items),
		Page:  in.Page,
		Limit: in.Limit,
		Total: total,
	}, nil
}
