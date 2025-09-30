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
func (s *TaskService) Task(ctx context.Context, taskID int) (*entity.Task, error) {
	return nil, nil
}
