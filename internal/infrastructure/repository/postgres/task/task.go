package task

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"
	"database/sql"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) GetTaskByID(ctx context.Context, taskID int64) (*entity.Task, error) {
	var task entity.Task
	row := r.db.QueryRowContext(ctx, "SELECT * FROM tasks WHERE id = $1", taskID)
	if err := row.Scan(&task.ID, &task.Title, &task.Content, &task.Options, &task.CorrectAnswer); err != nil {
		return nil, err
	}
	return &task, nil
}
