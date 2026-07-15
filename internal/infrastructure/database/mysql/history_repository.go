package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/human9001/teams/internal/domain/history/model"
)

type TaskHistoryRepository struct {
	db *sqlx.DB
	tx TxManager
}

func NewTaskHistoryRepository(db *sqlx.DB, tx TxManager) *TaskHistoryRepository {
	return &TaskHistoryRepository{db: db, tx: tx}
}

func (r *TaskHistoryRepository) ListByTaskID(ctx context.Context, taskID int64) ([]model.HistoryItem, error) {
	const q = `
		SELECT 
			id, 
			task_id, 
			event_type, 
			actor_id, 
			before_state, 
			after_state, 
			created_at
		FROM 
			task_history
		WHERE 
			task_id = ?
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, q, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.HistoryItem, 0)
	for rows.Next() {
		var itemRow model.HistoryItem
		if err := rows.Scan(
			&itemRow.ID,
			&itemRow.TaskID,
			&itemRow.EventType,
			&itemRow.ActorID,
			&itemRow.Before,
			&itemRow.After,
			&itemRow.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, itemRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
