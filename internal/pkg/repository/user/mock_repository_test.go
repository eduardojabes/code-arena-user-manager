package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	userService "github.com/eduardojabes/CodeArena/internal/pkg/service/user"
	"github.com/google/uuid"
)

func TestAddUserMemoryRepository(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		want := errors.New("error")

		repository := NewMockUserRepositoryWithError(want)

		user := userService.User{ID: uuid.New(), Username: "lenhador", Password: "password"}

		err := repository.AddUser(context.Background(), user)
		if err == nil {
			t.Errorf("got %v want %v", err, want)
		}
	})

	t.Run("test adding user", func(t *testing.T) {
		repository := NewMockUserRepository()

		user := userService.User{ID: uuid.New(), Username: "lenhador", Password: "password"}

		err := repository.AddUser(context.Background(), user)
		if err != nil {
			t.Errorf("got %v want nil", err)
		}

		got := repository.UserMemoryDB[user.ID]
		if !reflect.DeepEqual(got, user) {
			t.Errorf("got %v want %v", got, user)
		}
	})
}

func TestGetUserMemoryRepository(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		want := errors.New("error")

		repository := NewMockUserRepositoryWithError(want)

		_, err := repository.GetUser(context.Background(), uuid.New())
		if err == nil {
			t.Errorf("got %v want %v", err, want)
		}
	})

	t.Run("test getting user if user not exists", func(t *testing.T) {
		repository := NewMockUserRepository()

		got, err := repository.GetUser(context.Background(), uuid.New())
		if err != nil {
			t.Errorf("Get User got %v want nil", err)
		}

		if got != nil {
			t.Errorf("got %v want nil", got)
		}
	})

	t.Run("test getting user if user exists", func(t *testing.T) {
		repository := NewMockUserRepository()

		user := &userService.User{ID: uuid.New(), Username: "lenhador", Password: "password"}

		err := repository.AddUser(context.Background(), *user)
		if err != nil {
			t.Errorf("Add User got %v want nil", err)
		}

		got, err := repository.GetUser(context.Background(), user.ID)
		if err != nil {
			t.Errorf("Get User got %v want nil", err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("got %v want %v", got, user)
		}
	})

}

func TestSearchUserByUsernameMemoryRepository(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		want := errors.New("error")

		repository := NewMockUserRepositoryWithError(want)

		_, err := repository.SearchUserByUsername(context.Background(), "lenhador")
		if err == nil {
			t.Errorf("got %v want %v", err, want)
		}
	})

	t.Run("test getting user if user not exists", func(t *testing.T) {
		repository := NewMockUserRepository()

		got, err := repository.SearchUserByUsername(context.Background(), "lenhador")
		if err != nil {
			t.Errorf("Get User got %v want nil", err)
		}

		if got != nil {
			t.Errorf("got %v want nil", got)
		}
	})

	t.Run("test getting user if user exists", func(t *testing.T) {
		repository := NewMockUserRepository()

		user := &userService.User{ID: uuid.New(), Username: "lenhador", Password: "password"}

		err := repository.AddUser(context.Background(), *user)
		if err != nil {
			t.Errorf("Add User got %v want nil", err)
		}

		got, err := repository.SearchUserByUsername(context.Background(), user.Username)
		if err != nil {
			t.Errorf("Get User got %v want nil", err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("got %v want %v", got, user)
		}
	})

}
