package user

import (
	"context"

	userService "github.com/eduardojabes/CodeArena/internal/pkg/service/user"
	"github.com/google/uuid"
)

type MockUserRepositorys struct {
	UserMemoryDB map[uuid.UUID]userService.User
	err          error
}

func NewMockUserRepository() *MockUserRepositorys {
	return &MockUserRepositorys{UserMemoryDB: make(map[uuid.UUID]userService.User), err: nil}
}

func NewMockUserRepositoryWithError(err error) *MockUserRepositorys {
	return &MockUserRepositorys{UserMemoryDB: make(map[uuid.UUID]userService.User), err: err}
}

func (mur *MockUserRepositorys) GetUser(ctx context.Context, ID uuid.UUID) (*userService.User, error) {
	if mur.err != nil {
		return nil, mur.err
	}

	user, ok := mur.UserMemoryDB[ID]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (mur *MockUserRepositorys) AddUser(ctx context.Context, user userService.User) error {
	if mur.err != nil {
		return mur.err
	}

	mur.UserMemoryDB[user.ID] = user
	return nil
}

func (mur *MockUserRepositorys) SearchUserByUsername(ctx context.Context, name string) (*userService.User, error) {
	if mur.err != nil {
		return nil, mur.err
	}

	for _, user := range mur.UserMemoryDB {
		if user.Username == name {
			return &user, nil
		}
	}
	return nil, nil
}
