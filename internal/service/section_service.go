package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"black-pearl/backend-hackathon/internal/domain/sections/entity"
	"black-pearl/backend-hackathon/internal/domain/sections/interfaces"
	handlerInterfaces "black-pearl/backend-hackathon/internal/handler"
)

var _ handlerInterfaces.SectionServiceInterface = (*SectionService)(nil)

var (
	ErrInvalidTitle = errors.New("invalid title of section")
)

type SectionService struct {
	repo interfaces.SectionsRepoInterface
}

func NewSectionService(repo interfaces.SectionsRepoInterface) *SectionService {
	return &SectionService{repo: repo}
}

// Получить все секции
func (s *SectionService) GetSections(ctx context.Context) (*[]entity.Sections, error) {
	sections, err := s.repo.GetSections(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return &[]entity.Sections{}, nil
		}
		return nil, fmt.Errorf("failed to get sections: %w", err)
	}
	return sections, nil
}

// Создать новую секцию
func (s *SectionService) NewSection(ctx context.Context, title string) (*entity.Sections, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}

	section := &entity.Sections{
		Title: title,
	}

	if err := s.repo.CreateSection(ctx, section); err != nil {
		return nil, fmt.Errorf("failed to create section: %w", err)
	}

	return section, nil
}
