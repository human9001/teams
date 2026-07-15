package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/human9001/teams/internal/domain/comment/model"
)

type CommentRepository struct {
	db *sqlx.DB
	tx TxManager
}

func NewCommentRepository(db *sqlx.DB, tx TxManager) *CommentRepository {
	return &CommentRepository{db: db, tx: tx}
}

func (r *CommentRepository) ListByTaskID(ctx context.Context, taskID, limit, offset int64) ([]model.Comment, int64, error) {
	const countq = `SELECT COUNT(*) FROM task_comments WHERE task_id = ?`
	var total int64
	if err := r.db.QueryRowContext(ctx, countq, taskID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, task_id, author_id, body, created_at
		FROM task_comments
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, taskID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.Comment, 0)
	for rows.Next() {
		var rowItem model.Comment
		if err := rows.Scan(&rowItem.ID, &rowItem.TaskID, &rowItem.AuthorID, &rowItem.Body, &rowItem.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, rowItem)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *CommentRepository) UserHasAccess(ctx context.Context, taskID, userID int64) (bool, error) {
	const query = `
				SELECT EXISTS (
				SELECT * from tasks 
				INNER join team_members on team_members.team_id = tasks.team_id
				WHERE tasks.id = ? and team_members.user_id = ? ;  
			);
	`
	var ok bool
	if err := r.db.QueryRowContext(ctx, query, taskID, userID).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}
