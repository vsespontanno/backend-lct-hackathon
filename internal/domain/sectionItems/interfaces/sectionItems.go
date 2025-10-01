package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	"context"
)

type SectionItemsRepoInterface interface {
	GetSectionItemsBySectionId(ctx context.Context, sectionId int64) (*[]entity.SectionItem, error)
	CreateSectionItem(ctx context.Context, item *entity.SectionItem) error
}
