package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/human9001/teams/internal/domain/user"
)

type UserRepository struct {
	db *sqlx.DB
	tx TxManager
}

type userRow struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUserRepository(db *sqlx.DB, tx TxManager) *UserRepository {
	return &UserRepository{db: db, tx: tx}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`, u.Name(), u.Email(), u.PasswordHash(), time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (r *UserRepository) ByEmail(ctx context.Context, email string) (*user.User, error) {
	var row userRow
	err := r.db.GetContext(ctx, &row, `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return user.NewFromPersistence(row.ID, row.Name, row.Email, row.Password, row.CreatedAt, row.UpdatedAt), nil
}
