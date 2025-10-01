package service

import (
	"black-pearl/backend-hackathon/internal/domain/theory/entity"
	"black-pearl/backend-hackathon/internal/domain/theory/interfaces"
	"context"

	"go.uber.org/zap"
)

type TheoryService struct {
	repo   interfaces.TheoryRepoInterface
	logger *zap.SugaredLogger
}

func NewTheoryService(repo interfaces.TheoryRepoInterface, logger *zap.SugaredLogger) *TheoryService {
	return &TheoryService{repo: repo, logger: logger}
}

// Получить теорию по ID
func (s *TheoryService) GetTheoryByID(ctx context.Context, id int64) (*entity.Theory, error) {
	theory, err := s.repo.GetTheoryByID(ctx, id)
	if err != nil {
		s.logger.Errorw("failed to get theory", "error", err, "stage", "GetTheoryByID.GetTheoryByID")
		return nil, err
	}
	return theory, nil
}

// Создать новую теорию
func (s *TheoryService) NewTheory(ctx context.Context, title, content string) (*entity.Theory, error) {
	theory := &entity.Theory{
		Title:   title,
		Content: content,
	}

	if err := s.repo.CreateTheory(ctx, theory); err != nil {
		s.logger.Errorw("failed to create theory", "error", err, "stage", "NewTheory.CreateTheory")
		return nil, err
	}

	return theory, nil
}
