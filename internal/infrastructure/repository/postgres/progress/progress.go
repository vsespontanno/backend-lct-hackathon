package progress

import (
	"black-pearl/backend-hackathon/internal/domain/progress/entity"
	"context"
	"database/sql"
	"fmt"
)

type ProgressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) *ProgressRepo {
	return &ProgressRepo{
		db: db,
	}
}

func (r *ProgressRepo) GetProgressByID(ctx context.Context, userID int64) (*entity.Progress, error) {
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

func (r *ProgressRepo) UpsertProgress(ctx context.Context, progress *entity.Progress) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO progress (user_id, points, progress) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET points = $2, progress = $3", progress.UserID, progress.Points, progress.Progress)
	return err
}
