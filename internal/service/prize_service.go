package service

import (
	"black-pearl/backend-hackathon/internal/domain/prize/entity"
	"black-pearl/backend-hackathon/internal/domain/prize/interfaces"
	"context"

	"go.uber.org/zap"
)

type PrizeService struct {
	repo   interfaces.PrizeRepoInterface
	logger *zap.SugaredLogger
}

func NewPrizeService(repo interfaces.PrizeRepoInterface, logger *zap.SugaredLogger) *PrizeService {
	return &PrizeService{repo: repo, logger: logger}
}

func (s *PrizeService) MyPrizes(ctx context.Context, userID int) (*[]entity.Prize, error) {
	prizes, err := s.repo.GetMyPrizes(ctx, userID)
	if err != nil {
		return nil, err
	}
	return prizes, nil
}

func (s *PrizeService) AvailablePrizes(ctx context.Context, userID int) (*[]entity.Prize, error) {
	prizes, err := s.repo.GetAvailablePrizes(ctx, userID)
	if err != nil {
		return nil, err
	}
	return prizes, nil
}
