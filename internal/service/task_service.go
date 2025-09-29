package service

import (
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/domain/task/interfaces"
	"context"
)

type TaskService struct {
	repo interfaces.TaskInterface
}

func NewTaskService(repo interfaces.TaskInterface) *TaskService {
	return &TaskService{repo: repo}
}

// Временные затычки
func (s *TaskService) Task(ctx context.Context, taskID int64) (*entity.Task, error) {
	return nil, nil
}
