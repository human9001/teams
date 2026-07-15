package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/human9001/teams/internal/application/service/task/input"
	domainErrors "github.com/human9001/teams/internal/domain/errors"
	"github.com/human9001/teams/internal/domain/task/model"
)

type TaskRepository struct {
	db *sqlx.DB
	tx TxManager
}

func NewTaskRepository(db *sqlx.DB, tx TxManager) *TaskRepository {
	return &TaskRepository{db: db, tx: tx}
}

type taskRow struct {
	ID          int64     `db:"id"`
	TeamID      int64     `db:"team_id"`
	CreatedBy   int64     `db:"created_by"`
	AssigneeID  *int64    `db:"assignee_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	Priority    string    `db:"priority"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *TaskRepository) Save(ctx context.Context, t *model.Task) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (
			team_id, created_by, assignee_id, title, description,
			status, priority, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		t.TeamID(),
		t.CreatedBy(),
		t.AssigneeID(),
		t.Title(),
		t.Description(),
		string(t.Status()),
		string(t.Priority()),
		t.CreatedAt(),
		t.UpdatedAt(),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func (r *TaskRepository) ByID(ctx context.Context, id int64) (model.Task, error) {
	var row taskRow
	err := r.db.GetContext(ctx, &row, `
		SELECT
			id, team_id, created_by, assignee_id, title, description,
			status, priority, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, domainErrors.ErrNotFound
		}
		return model.Task{}, err
	}

	return model.NewFromPersistence(
		row.ID,
		row.TeamID,
		row.CreatedBy,
		row.AssigneeID,
		row.Title,
		row.Description,
		model.Status(row.Status),
		model.Priority(row.Priority),
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func (r *TaskRepository) List(ctx context.Context, filter input.ListTasksInput) ([]model.Task, int64, error) {
	where := []string{"team_id = ?"}
	args := []any{filter.TeamID}

	if filter.Status != nil {
		where = append(where, "status = ?")
		args = append(args, string(*filter.Status))
	}
	if filter.AssigneeID != nil {
		where = append(where, "assignee_id = ?")
		args = append(args, *filter.AssigneeID)
	}

	whereSQL := strings.Join(where, " AND ")
	offset := (filter.Page - 1) * filter.Limit

	countQuery := "SELECT COUNT(*) FROM tasks WHERE " + whereSQL
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return []model.Task{}, 0, err
	}

	query := `
		SELECT id, team_id, created_by, assignee_id, title, description, status, priority, created_at, updated_at
		FROM tasks
		WHERE ` + whereSQL + `
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	args2 := append([]any{}, args, filter.Limit, offset)

	rows, err := r.db.QueryxContext(ctx, query, args2...)
	if err != nil {
		return []model.Task{}, 0, err
	}
	defer rows.Close()

	items := make([]model.Task, 0)
	for rows.Next() {
		var row taskRow
		if err := rows.StructScan(&row); err != nil {
			return []model.Task{}, 0, err
		}
		items = append(items, model.NewFromPersistence(
			row.ID, row.TeamID, row.CreatedBy, row.AssigneeID,
			row.Title, row.Description, model.Status(row.Status),
			model.Priority(row.Priority), row.CreatedAt, row.UpdatedAt,
		))
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, 0, err
	}

	return items, total, nil
}

func (r *TaskRepository) Update(ctx context.Context, t model.Task) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?, assigned_to = ?, due_date = ?, updated_at = NOW()
		WHERE id = ?
	`, t.Title(), t.Description(), string(t.Status()), string(t.Priority()), t.AssigneeID(), t.ID())
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) AddComment(ctx context.Context, in model.TaskComment) error {
	const q = `
		INSERT INTO task_comments (task_id, author_id, body)
		VALUES (?, ?, ?)
	`

	var rowItem model.TaskComment
	err := r.db.QueryRowContext(ctx, q, in.TaskID, in.UserID, in.Body).Scan(
		&rowItem.ID,
		&rowItem.TaskID,
		&rowItem.UserID,
		&rowItem.Body,
		&rowItem.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetComments(ctx context.Context, taskID int64) ([]model.TaskComment, int64, error) {
	const countq = `SELECT COUNT(*) FROM task_comments WHERE task_id = ?`
	var total int64
	if err := r.db.QueryRowContext(ctx, countq, taskID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `
		SELECT 
			id, 
			task_id, 
			author_id, 
			body, 
			created_at
		FROM 
			task_comments
		WHERE 
			task_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, q, taskID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.TaskComment, 0)
	for rows.Next() {
		var rowItem model.TaskComment
		if err := rows.Scan(&rowItem.ID, &rowItem.TaskID, &rowItem.UserID, &rowItem.Body, &rowItem.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, rowItem)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
