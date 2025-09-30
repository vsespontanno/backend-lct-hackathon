package progress

import (
	"black-pearl/backend-hackathon/internal/domain/progress/entity"
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

type ProgressRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewProgressRepo(db *sql.DB) *ProgressRepo {
	return &ProgressRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ProgressRepo) GetProgressByID(ctx context.Context, userID int64) (*entity.Progress, error) {
	builder := squirrel.Select("user_id", "points", "progress").
		From("progress").
		Where(squirrel.Eq{"user_id": userID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var progress entity.Progress
	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	if err := row.Scan(&progress.UserID, &progress.Points, &progress.Progress); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("progress not found")
		}
		return nil, err
	}
	return &progress, nil
}

// апдейтит прогресс, если нет прогресса, то создаст его с дефолт значениями
func (r *ProgressRepo) UpsertProgress(ctx context.Context, progress *entity.Progress) error {
	updateBuilder := squirrel.Update("progress").
		Set("points", progress.Points).
		Set("progress", progress.Progress).
		Where(squirrel.Eq{"user_id": progress.UserID})

	updateSQL, updateArgs, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, updateSQL, updateArgs...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		insertBuilder := squirrel.Insert("progress").
			Columns("user_id", "points", "progress").
			Values(progress.UserID, progress.Points, progress.Progress).
			PlaceholderFormat(squirrel.Dollar)

		insertSQL, insertArgs, err := insertBuilder.ToSql()
		if err != nil {
			return err
		}

		_, err = r.db.ExecContext(ctx, insertSQL, insertArgs...)
		if err != nil {
			return err
		}
	}

	return nil
}
