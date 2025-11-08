package repo

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo/pgdb"
	"denet-test-task/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
	GetUserById(ctx context.Context, id int) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)

	SetUserReferrer(ctx context.Context, id int, referrer string) error
	SetUserEmail(ctx context.Context, id int, email string) error
}

type Repositories struct {
	User
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(pg),
	}
}
