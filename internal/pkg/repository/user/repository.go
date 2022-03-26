package user

import (
	"context"
	"fmt"

	"github.com/eduardojabes/code-arena-user-manager/internal/pkg/entity"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

type UserRepository interface {
	GetUser(ctx context.Context, ID uuid.UUID) (*entity.User, error)
	AddUser(ctx context.Context, user entity.User) error
	SearchUserByUsername(ctx context.Context, name string) (*entity.User, error)
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
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func NewPostgreUserRepository(conn connector) *PostgreUserRepository {
	return &PostgreUserRepository{conn}
}

func (r *PostgreUserRepository) GetUser(ctx context.Context, ID uuid.UUID) (*entity.User, error) {
	var users []*UserModel
	err := pgxscan.Select(ctx, r.conn, &users, `SELECT * FROM arena_user WHERE u_id = $1`, ID)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &entity.User{
		ID:       users[0].ID,
		Username: users[0].Username,
		Password: users[0].Password,
	}, nil
}

func (r *PostgreUserRepository) AddUser(ctx context.Context, user entity.User) error {
	_, err := r.conn.Exec(ctx, `INSERT INTO arena_user(u_id, u_username, u_password) values($1, $2, $3)`, user.ID, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreUserRepository) SearchUserByUsername(ctx context.Context, name string) (*entity.User, error) {
	var users []*UserModel
	err := pgxscan.Select(ctx, r.conn, &users, `SELECT * FROM arena_user WHERE u_username = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &entity.User{
		ID:       users[0].ID,
		Username: users[0].Username,
		Password: users[0].Password,
	}, nil
}
