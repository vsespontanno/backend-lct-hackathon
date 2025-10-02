package service

import (
	"black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	"black-pearl/backend-hackathon/internal/domain/sectionItems/interfaces"
	"context"

	"go.uber.org/zap"
)

type SectionItemsService struct {
	repo   interfaces.SectionItemsRepoInterface
	logger *zap.SugaredLogger
}

func NewSectionItemsService(repo interfaces.SectionItemsRepoInterface, logger *zap.SugaredLogger) *SectionItemsService {
	return &SectionItemsService{repo: repo, logger: logger}
}

// Получить все айтемы секции
func (s *SectionItemsService) GetSectionItemsBySectionID(ctx context.Context, sectionID int) (*[]entity.SectionItem, error) {
	items, err := s.repo.GetSectionItemsBySectionId(ctx, sectionID)
	if err != nil {
		s.logger.Errorw("failed to get section items", "error", err, "stage", "GetSectionItemsBySectionID.GetSectionItemsBySectionId")
		return nil, err
	}
	return items, nil
}

// Создать новый айтем в секции
func (s *SectionItemsService) NewSectionItem(ctx context.Context, sectionID int, title string, isTest bool, itemID int) (*entity.SectionItem, error) {
	sectionItem := &entity.SectionItem{
		SectionID: sectionID,
		Title:     title,
		IsTest:    isTest,
		ItemID:    itemID,
	}

	if err := s.repo.CreateSectionItem(ctx, sectionItem); err != nil {
		s.logger.Errorw("failed to create section item", "error", err, "stage", "NewSectionItem.CreateSectionItem")
		return nil, err
	}

	return sectionItem, nil
}
