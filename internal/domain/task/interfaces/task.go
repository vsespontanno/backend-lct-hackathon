package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"
)

type TaskInterface interface {
	GetTaskByID(ctx context.Context, taskID int64) (*entity.Task, error)
}
