package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/sections/entity"
	"context"
)

type SectionsRepoInterface interface {
	GetSections(ctx context.Context) (*[]entity.Sections, error)
	CreateSection(ctx context.Context, s *entity.Sections) error
}
