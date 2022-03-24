package user

import (
	"context"

	userService "github.com/eduardojabes/CodeArena/internal/pkg/service/user"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUser(ctx context.Context, ID uuid.UUID) (*userService.User, error)
	AddUser(ctx context.Context, user userService.User) (*userService.User, error)
	SearchUserByUsername(ctx context.Context, name string) (*userService.User, error)
}

type PostgreUserRepository struct {
}

func (r *PostgreUserRepository) GetUser(ctx context.Context, ID uuid.UUID) (*userService.User, error) {
	return nil, nil
}
func (r *PostgreUserRepository) AddUser(ctx context.Context, user userService.User) (*userService.User, error) {
	return nil, nil
}
func (r *PostgreUserRepository) SearchUserByUsername(ctx context.Context, name string) (*userService.User, error) {
	return nil, nil
}
