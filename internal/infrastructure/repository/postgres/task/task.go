package task

import (
	"database/sql"

	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"

	sq "github.com/Masterminds/squirrel"

	"go.uber.org/zap"
)

type TaskRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
	logger  *zap.SugaredLogger
}

func NewTaskRepo(db *sql.DB, logger *zap.SugaredLogger) *TaskRepo {
	return &TaskRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		logger:  logger,
	}
}

func (r *TaskRepo) GetTasks(ctx context.Context) (*[]entity.Task, error) {
	builder := r.builder.
		Select("id", "title").
		From("tasks")

	query, args, err := builder.ToSql()
	if err != nil {
		r.logger.Errorw("failed to build query", "error", err, "stage", "GetTasks.ToSql")
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorw("failed to execute query", "error", err, "stage", "GetTasks.QueryContext")
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var t entity.Task
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			r.logger.Errorw("failed to scan row", "error", err, "stage", "GetTasks.Scan")
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return &tasks, nil
}
