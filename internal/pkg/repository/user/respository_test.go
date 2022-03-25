package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"

	userService "github.com/eduardojabes/CodeArena/internal/pkg/service/user"
)

func TestUserRepository(t *testing.T) {
	t.Run("no_rows", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()
		mock.ExpectQuery("SELECT (.+) FROM arena_user WHERE (.+)").
			WillReturnRows(mock.NewRows([]string{"u_id", "u_username", "u_password"}))

		repository := NewPostgreUserRepository(mock)

		user, err := repository.SearchUserByUsername(context.Background(), "lenhador")

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}

		if user != nil {
			t.Errorf("got %v want nil", user)
		}
	})

	t.Run("with_user", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		user := &userService.User{
			ID:       uuid.New(),
			Username: "lenhador",
			Password: "12345",
		}

		mock.ExpectQuery("SELECT (.+) FROM arena_user WHERE (.+)").
			WillReturnRows(mock.NewRows([]string{"u_id", "u_username", "u_password"}).
				AddRow(user.ID, user.Username, user.Password))

		repository := NewPostgreUserRepository(mock)

		got, err := repository.SearchUserByUsername(context.Background(), "lenhador")

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}

		if !reflect.DeepEqual(user, got) {
			t.Errorf("got %v want %v", got, user)
		}
	})

	t.Run("with_error", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		mock.ExpectQuery("SELECT (.+) FROM arena_user WHERE (.+)").
			WillReturnError(errors.New("error"))

		repository := NewPostgreUserRepository(mock)

		_, err := repository.SearchUserByUsername(context.Background(), "lenhador")

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}
