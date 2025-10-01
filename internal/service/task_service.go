package service

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/domain/task/interfaces"
	"context"
)

type TaskService struct {
	repo interfaces.TaskRepoInterface
}

func NewTaskService(repo interfaces.TaskRepoInterface) *TaskService {
	return &TaskService{repo: repo}
}

// Временные затычки
func (s *TaskService) GetTask(ctx context.Context, taskID int64) (*entity.Quiz, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return task, err
}
