package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/prize/entity"
	"context"
)

type PrizeRepoInterface interface {
	GetMyPrizes(ctx context.Context, userID int) (*[]entity.Prize, error)
	GetAvailablePrizes(ctx context.Context, userID int) (*[]entity.Prize, error)
}
