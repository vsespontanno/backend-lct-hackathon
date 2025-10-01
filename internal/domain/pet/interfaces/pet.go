package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	"context"
)

type PetRepoInterface interface {
	GetPetByUserID(ctx context.Context, userID int) (*entity.Pet, error)
	UpdateXP(ctx context.Context, xp int, userID int) error
	SetPetName(ctx context.Context, name string, userID int) error
	CreatePet(ctx context.Context, userID int) error
}
