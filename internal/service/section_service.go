package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"black-pearl/backend-hackathon/internal/domain/sections/entity"
	"black-pearl/backend-hackathon/internal/domain/sections/interfaces"

	"go.uber.org/zap"
)

var (
	ErrInvalidTitle = errors.New("invalid title of section")
)

type SectionService struct {
	repo   interfaces.SectionsRepoInterface
	logger *zap.SugaredLogger
}

func NewSectionService(repo interfaces.SectionsRepoInterface, logger *zap.SugaredLogger) *SectionService {
	return &SectionService{repo: repo, logger: logger}
}

// Получить все секции
func (s *SectionService) GetSections(ctx context.Context) (*[]entity.Sections, error) {
	sections, err := s.repo.GetSections(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return &[]entity.Sections{}, nil
		}
		s.logger.Errorw("failed to get sections", "error", err, "stage", "GetSections.GetSections")
		return nil, fmt.Errorf("failed to get sections: %w", err)
	}
	return sections, nil
}

// Создать новую секцию
func (s *SectionService) NewSection(ctx context.Context, title string) (*entity.Sections, error) {
	if title == "" {
		s.logger.Errorw("failed to create section", "error", ErrInvalidTitle, "stage", "NewSection.NewSection")
		return nil, ErrInvalidTitle
	}

	section := &entity.Sections{
		Title: title,
	}

	if err := s.repo.CreateSection(ctx, section); err != nil {
		s.logger.Errorw("failed to create section", "error", err, "stage", "NewSection.CreateSection")
		return nil, err
	}

	return section, nil
}
