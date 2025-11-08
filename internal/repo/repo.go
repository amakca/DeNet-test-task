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

type Task interface {
	GetTaskById(ctx context.Context, id int) (entity.Task, error)
	GetTaskByName(ctx context.Context, name string) (entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}

type Point interface {
	AddPointsByUserId(ctx context.Context, userId int, taskId int, points int) error
	GetPointsByUserId(ctx context.Context, userId int) (int, error)
	GetHistoryByUserId(ctx context.Context, userId int) ([]entity.Point, error)
	CheckCompletedTask(ctx context.Context, userId int, taskId int) (bool, error)
}

type Repositories struct {
	User
	Task
	Point
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:  pgdb.NewUserRepo(pg),
		Task:  pgdb.NewTaskRepo(pg),
		Point: pgdb.NewPointRepo(pg),
	}
}
