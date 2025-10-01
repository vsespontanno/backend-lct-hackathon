package task

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"
	"database/sql"
	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"

	"fmt"
)

type TaskRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TaskRepo) GetTaskByID(ctx context.Context, taskID int64) (*entity.Task, error) {
	sqlStr := `SELECT id, title, content, options, correct_answer FROM tasks WHERE id = $1`

	var task entity.Task
	var options pq.StringArray

	row := r.db.QueryRowContext(ctx, sqlStr, taskID)
	err := row.Scan(&task.ID, &task.Title, &task.Content, &options, &task.CorrectAnswer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	task.Options = []string(options)
	return &task, nil
}

func (r *TaskRepo) InsertTask(ctx context.Context, task *entity.Task) error {
	sqlStr := `INSERT INTO tasks (title, content, options, correct_answer) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, sqlStr, task.Title, task.Content, pq.Array(task.Options), task.CorrectAnswer)
	return err
}
