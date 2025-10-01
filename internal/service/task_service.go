package service

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/domain/task/interfaces"
	"context"

	"go.uber.org/zap"
)

type TaskService struct {
	repo   interfaces.TaskRepoInterface
	logger *zap.SugaredLogger
}

func NewTaskService(repo interfaces.TaskRepoInterface, logger *zap.SugaredLogger) *TaskService {
	return &TaskService{repo: repo, logger: logger}
}

// Временные затычки
func (s *TaskService) Task(ctx context.Context, taskID int) (*entity.Task, error) {
	return nil, nil
}
