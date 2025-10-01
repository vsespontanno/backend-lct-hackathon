package service

import (
	"black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	"black-pearl/backend-hackathon/internal/domain/sectionItems/interfaces"
	"context"
)

type SectionItemsService struct {
	repo interfaces.SectionItemsRepoInterface
}

func NewSectionItemsService(repo interfaces.SectionItemsRepoInterface) *SectionItemsService {
	return &SectionItemsService{repo: repo}
}

// Получить все айтемы секции
func (s *SectionItemsService) GetSectionItemsBySectionID(ctx context.Context, sectionID int64) (*[]entity.SectionItem, error) {
	items, err := s.repo.GetSectionItemsBySectionId(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Создать новый айтем в секции
func (s *SectionItemsService) NewSectionItem(ctx context.Context, sectionID int64, title string, isTest bool, itemID int64) (*entity.SectionItem, error) {
	sectionItem := &entity.SectionItem{
		SectionID: sectionID,
		Title:     title,
		IsTest:    isTest,
		ItemID:    itemID,
	}

	if err := s.repo.CreateSectionItem(ctx, sectionItem); err != nil {
		return nil, err
	}

	return sectionItem, nil
}
