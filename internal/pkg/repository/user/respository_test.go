package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardojabes/code-arena-user-manager/internal/pkg/entity"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func TestGetUser(t *testing.T) {
	t.Run("no_rows", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()
		mock.ExpectQuery("SELECT (.+) FROM arena_user WHERE (.+)").
			WillReturnRows(mock.NewRows([]string{"u_id", "u_username", "u_password"}))

		repository := NewPostgreUserRepository(mock)

		user, err := repository.GetUser(context.Background(), uuid.New())

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}

		if user != nil {
			t.Errorf("got %v want nil", user)
		}
	})

	t.Run("with_user", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		user := &entity.User{
			ID:       uuid.New(),
			Username: "lenhador",
			Password: "12345",
		}

		mock.ExpectQuery("SELECT (.+) FROM arena_user WHERE (.+)").
			WillReturnRows(mock.NewRows([]string{"u_id", "u_username", "u_password"}).
				AddRow(user.ID, user.Username, user.Password))

		repository := NewPostgreUserRepository(mock)

		got, err := repository.GetUser(context.Background(), user.ID)

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

		_, err := repository.GetUser(context.Background(), uuid.New())

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}
func TestAddUser(t *testing.T) {
	t.Run("Adding User", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		user := &entity.User{
			ID:       uuid.New(),
			Username: "lenhador",
			Password: "12345",
		}

		mock.ExpectExec("INSERT INTO arena_user").
			WithArgs(user.ID, user.Username, user.Password).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		repository := NewPostgreUserRepository(mock)
		err := repository.AddUser(context.Background(), *user)

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}
	})

	t.Run("with_error", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		user := &entity.User{
			ID:       uuid.New(),
			Username: "lenhador",
			Password: "12345",
		}

		mock.ExpectQuery("SELECT (.+) FROM arena_user WHERE (.+)").
			WillReturnError(errors.New("error"))

		repository := NewPostgreUserRepository(mock)
		err := repository.AddUser(context.Background(), *user)

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}

func TestSearchUserByUsername(t *testing.T) {
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

		user := &entity.User{
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
