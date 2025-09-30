package service

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	petInterface "black-pearl/backend-hackathon/internal/domain/pet/interfaces"
	userInterface "black-pearl/backend-hackathon/internal/domain/user/interfaces"
	"context"
	"database/sql"
	"fmt"
)

var (
	ErrInvalidName = fmt.Errorf("invalid name of pet")
)

type PetService struct {
	petRepo  petInterface.PetRepoInterface
	userRepo userInterface.UserRepoInterface
}

func NewPetService(petRepo petInterface.PetRepoInterface, userRepo userInterface.UserRepoInterface) *PetService {
	return &PetService{
		petRepo:  petRepo,
		userRepo: userRepo,
	}
}

func (s *PetService) SetName(ctx context.Context, name string, userID int) error {
	if name == "" {
		return ErrInvalidName
	}
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		if err == sql.ErrNoRows {
			err = s.userRepo.CreateUser(ctx, userID)
			if err != nil {
				return err
			}
			err := s.petRepo.CreatePet(ctx, userID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err := s.petRepo.SetPetName(ctx, name, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PetService) GetPetByUserID(ctx context.Context, userID int) (*entity.Pet, error) {
	pet, err := s.petRepo.GetPetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return pet, nil
}

func (s *PetService) UpdateXP(ctx context.Context, xp int, userID int) error {
	err := s.petRepo.UpdateXP(ctx, xp, userID)
	if err != nil {
		return err
	}
	return nil
}
