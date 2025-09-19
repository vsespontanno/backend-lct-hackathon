package service

import (
	"black-pearl/backend-hackathon/internal/domain/entity"
	"black-pearl/backend-hackathon/internal/domain/interfaces"
	"context"
)

type TaskService struct {
	repo interfaces.Repository
}

func NewTaskService(repo interfaces.Repository) *TaskService {
	return &TaskService{repo: repo}
}

// Временные затычки
func (s *TaskService) Task(ctx context.Context, taskID int64) (*entity.Task, error) {
	return nil, nil
}

// Временные затычки
func (s *TaskService) Progress(ctx context.Context, taskID int64) (*entity.Progress, error) {
	return nil, nil
}
