package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"
)

type TaskRepoInterface interface {
	GetTaskByID(ctx context.Context, taskID int) (*entity.Task, error)
}
