package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/user/entity"
	"context"
)

type UserRepoInterface interface {
	CreateUser(ctx context.Context, id int) error
	GetUserByID(ctx context.Context, id int) (entity.User, error)
}
