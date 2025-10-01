package service

import (
	"black-pearl/backend-hackathon/internal/domain/theory/entity"
	"black-pearl/backend-hackathon/internal/domain/theory/interfaces"
	"context"
)

type TheoryService struct {
	repo interfaces.TheoryRepoInterface
}

func NewTheoryService(repo interfaces.TheoryRepoInterface) *TheoryService {
	return &TheoryService{repo: repo}
}

// Получить теорию по ID
func (s *TheoryService) GetTheoryByID(ctx context.Context, id int64) (*entity.Theory, error) {
	return s.repo.GetTheoryByID(ctx, id)
}

// Создать новую теорию
func (s *TheoryService) NewTheory(ctx context.Context, title, content string) (*entity.Theory, error) {
	theory := &entity.Theory{
		Title:   title,
		Content: content,
	}

	if err := s.repo.CreateTheory(ctx, theory); err != nil {
		return nil, err
	}

	return theory, nil
}
