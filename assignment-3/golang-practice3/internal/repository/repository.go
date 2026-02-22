package repository

import (
	"context"

	"golang/internal/repository/_postgres"
	"golang/internal/repository/_postgres/users"
	"golang/pkg/modules"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]modules.User, error)
	GetUserByID(ctx context.Context, id int) (*modules.User, error)
	CreateUser(ctx context.Context, u modules.User) (*modules.User, error)
	UpdateUser(ctx context.Context, id int, u modules.User) error
	DeleteUser(ctx context.Context, id int) (int64, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
