package service

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	petInterface "black-pearl/backend-hackathon/internal/domain/pet/interfaces"
	userInterface "black-pearl/backend-hackathon/internal/domain/user/interfaces"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

var (
	ErrInvalidName = fmt.Errorf("invalid name of pet")
)

type PetService struct {
	petRepo  petInterface.PetRepoInterface
	userRepo userInterface.UserRepoInterface
	logger   *zap.SugaredLogger
}

func NewPetService(petRepo petInterface.PetRepoInterface, userRepo userInterface.UserRepoInterface, logger *zap.SugaredLogger) *PetService {
	return &PetService{
		petRepo:  petRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *PetService) SetName(ctx context.Context, name string, userID int) error {
	if name == "" {
		s.logger.Errorw("invalid name of pet", "error", ErrInvalidName, "stage", "SetName")
		return ErrInvalidName
	}
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		if err == sql.ErrNoRows {
			err = s.userRepo.CreateUser(ctx, userID)
			if err != nil {
				s.logger.Errorw("failed to create user", "error", err, "stage", "SetName.CreateUser")
				return err
			}
			err := s.petRepo.CreatePet(ctx, userID)
			if err != nil {
				s.logger.Errorw("failed to create pet", "error", err, "stage", "SetName.CreatePet")
				return err
			}
		} else {
			s.logger.Errorw("failed to get user", "error", err, "stage", "SetName.GetUserByID")
			return err
		}
	}

	err := s.petRepo.SetPetName(ctx, name, userID)
	if err != nil {
		s.logger.Errorw("failed to set pet name", "error", err, "stage", "SetName.SetPetName")
		return err
	}
	return nil
}

func (s *PetService) GetPetByUserID(ctx context.Context, userID int) (*entity.Pet, error) {
	pet, err := s.petRepo.GetPetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorw("failed to get pet", "error", err, "stage", "GetPetByUserID.GetPetByUserID")
		return nil, err
	}
	return pet, nil
}

func (s *PetService) UpdateXP(ctx context.Context, xp int, userID int) error {
	err := s.petRepo.UpdateXP(ctx, xp, userID)
	if err != nil {
		s.logger.Errorw("failed to update xp", "error", err, "stage", "UpdateXP.UpdateXP")
		return err
	}
	return nil
}
