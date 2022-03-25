package user

import (
	"context"
	"fmt"

	userService "github.com/eduardojabes/CodeArena/internal/pkg/service/user"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUser(ctx context.Context, ID uuid.UUID) (*userService.User, error)
	AddUser(ctx context.Context, user userService.User) (*userService.User, error)
	SearchUserByUsername(ctx context.Context, name string) (*userService.User, error)
}

type UserModel struct {
	ID       uuid.UUID `db:"u_id"`
	Username string    `db:"u_username"`
	Password string    `db:"u_password"`
}

type PostgreUserRepository struct {
	conn connector
}

type connector interface {
	pgxscan.Querier
}

func NewPostgreUserRepository(conn connector) *PostgreUserRepository {
	return &PostgreUserRepository{conn}
}

func (r *PostgreUserRepository) GetUser(ctx context.Context, ID uuid.UUID) (*userService.User, error) {
	return nil, nil
}
func (r *PostgreUserRepository) AddUser(ctx context.Context, user userService.User) (*userService.User, error) {

	return nil, nil
}
func (r *PostgreUserRepository) SearchUserByUsername(ctx context.Context, name string) (*userService.User, error) {
	var users []*UserModel
	err := pgxscan.Select(ctx, r.conn, &users, `SELECT * FROM arena_user WHERE u_username = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &userService.User{
		ID:       users[0].ID,
		Username: users[0].Username,
		Password: users[0].Password,
	}, nil
}
