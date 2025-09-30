package service

import (
	"black-pearl/backend-hackathon/internal/domain/prize/entity"
	"black-pearl/backend-hackathon/internal/domain/prize/interfaces"
	"context"
)

type PrizeService struct {
	repo interfaces.PrizeRepoInterface
}

func NewPrizeService(repo interfaces.PrizeRepoInterface) *PrizeService {
	return &PrizeService{repo: repo}
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
