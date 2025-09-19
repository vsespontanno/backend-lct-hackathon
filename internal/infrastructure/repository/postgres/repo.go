package postgres

import (
	"black-pearl/backend-hackathon/internal/domain/entity"
	"context"
	"database/sql"
	"fmt"
)

type PoostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PoostgresRepo {
	return &PoostgresRepo{
		db: db,
	}
}

func (r *PoostgresRepo) GetTaskByID(ctx context.Context, taskID int64) (*entity.Task, error) {
	var task entity.Task
	row := r.db.QueryRowContext(ctx, "SELECT * FROM tasks WHERE id = $1", taskID)
	if err := row.Scan(&task.ID, &task.Title, &task.Content, &task.Options, &task.CorrectAnswer); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PoostgresRepo) GetProgressByID(ctx context.Context, userID int64) (*entity.Progress, error) {
	var progress entity.Progress
	row := r.db.QueryRowContext(ctx, "SELECT * FROM progress WHERE user_id = $1", userID)
	if err := row.Scan(&progress.UserID, &progress.Points, &progress.Progress); err != nil {
		if err == sql.ErrNoRows {
			// по хорошему тут нужно инициализировать но я пока затычку сделаю
			// как с общей схемой определимся так и сделаем потом
			return nil, fmt.Errorf("progress not found")
		}
		return nil, err
	}
	return &progress, nil
}

func (r *PoostgresRepo) UpsertProgress(ctx context.Context, progress *entity.Progress) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO progress (user_id, points, progress) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET points = $2, progress = $3", progress.UserID, progress.Points, progress.Progress)
	return err
}
