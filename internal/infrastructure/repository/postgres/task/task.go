package task

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"
	"database/sql"

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

func (r *TaskRepo) GetTaskByID(ctx context.Context, taskID int) (*entity.Task, error) {
	builder := sq.Select("id", "title", "content", "options", "correct_answer").
		From("tasks").
		Where(sq.Eq{"id": taskID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var task entity.Task
	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	if err := row.Scan(&task.ID, &task.Title, &task.Content, &task.Options, &task.CorrectAnswer); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) InsertTask(ctx context.Context, task *entity.Task) error {
	builder := sq.Insert("tasks").
		Columns("title", "content", "options", "correct_answer").
		Values(task.Title, task.Content, task.Options, task.CorrectAnswer)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}
