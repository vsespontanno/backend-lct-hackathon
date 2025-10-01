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

func (s *TaskService) GetTasks(ctx context.Context) (*[]entity.Task, error) {
	task, err := s.repo.GetTasks(ctx)
	if err != nil {
		s.logger.Errorw("failed to get tasks", "error", err, "stage", "TaskService.GetTasks")
		return nil, err
	}

	return task, nil
}
