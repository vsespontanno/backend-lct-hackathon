package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/entity"
	"context"
)

type Repository interface {
	GetTaskByID(ctx context.Context, taskID int64) (*entity.Task, error)
	GetProgressByID(ctx context.Context, userID int64) (*entity.Progress, error)
	UpsertProgress(ctx context.Context, progress *entity.Progress) error
}
