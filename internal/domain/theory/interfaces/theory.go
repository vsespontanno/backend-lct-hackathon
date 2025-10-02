package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/theory/entity"
	"context"
)

type TheoryRepoInterface interface {
	GetTheoryByID(ctx context.Context, id int) (*entity.Theory, error)
	CreateTheory(ctx context.Context, t *entity.Theory) error
}
