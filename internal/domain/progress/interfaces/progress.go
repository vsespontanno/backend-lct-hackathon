package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/progress/entity"
	"context"
)

type ProgressRepoInterface interface {
	GetProgressByID(ctx context.Context, userID int64) (*entity.Progress, error)
	UpsertProgress(ctx context.Context, progress *entity.Progress) error
}
