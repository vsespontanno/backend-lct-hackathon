package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"context"
)

type TaskRepoInterface interface {
	GetTasks(ctx context.Context) (*[]entity.Task, error)
}
