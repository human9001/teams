package comment

import "github.com/human9001/teams/internal/domain/comment/repository"

type CommentService struct {
	repo repository.IRepository
}

func NewCommentService(repo repository.IRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}
