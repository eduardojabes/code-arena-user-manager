package test

import (
	"context"
	"errors"

	"github.com/eduardojabes/CodeArena/internal/pkg/entity"
	"github.com/google/uuid"
)

type MockUserRepository struct {
	GetUserMock              func(ctx context.Context, ID uuid.UUID) (*entity.User, error)
	AddUserMock              func(ctx context.Context, user entity.User) error
	SearchUserByUsernameMock func(ctx context.Context, name string) (*entity.User, error)
}

func (mur *MockUserRepository) GetUser(ctx context.Context, ID uuid.UUID) (*entity.User, error) {
	if mur.GetUserMock != nil {
		return mur.GetUserMock(ctx, ID)
	}

	return nil, errors.New("GetUserMock must be set")
}

func (mur *MockUserRepository) AddUser(ctx context.Context, user entity.User) error {
	if mur.AddUserMock != nil {
		return mur.AddUserMock(ctx, user)
	}

	return errors.New("AddUserMock must be set")
}

func (mur *MockUserRepository) SearchUserByUsername(ctx context.Context, name string) (*entity.User, error) {
	if mur.SearchUserByUsernameMock != nil {
		return mur.SearchUserByUsernameMock(ctx, name)
	}

	return nil, errors.New("SearchUserByUsernameMock must be set")
}
