package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	"context"
)

type PetInterface interface {
	GetPetByUserID(ctx context.Context, userID int) (*entity.Pet, error)
	SetPetName(ctx context.Context, name string, userID int) (bool, error)
}
