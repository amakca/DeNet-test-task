package repo

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo/pgdb"
	"denet-test-task/pkg/postgres"
)

type Users interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
	GetUserById(ctx context.Context, id int) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)

	SetUserReferrer(ctx context.Context, id int, referrer string) error
	SetUserEmail(ctx context.Context, id int, email string) error
}

type Tasks interface {
	GetTaskById(ctx context.Context, id int) (entity.Task, error)
	GetTaskByName(ctx context.Context, name string) (entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}

type Points interface {
	AddPointsByUserId(ctx context.Context, userId int, taskId int, points int) error
	GetPointsByUserId(ctx context.Context, userId int) (int, error)
	GetHistoryByUserId(ctx context.Context, userId int) ([]entity.Point, error)
	CheckCompletedTask(ctx context.Context, userId int, taskId int) (bool, error)
	GetLeaderboard(ctx context.Context, limit int) ([]entity.Point, error)
}

type Repositories struct {
	Users
	Tasks
	Points
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Users:  pgdb.NewUsersRepo(pg),
		Tasks:  pgdb.NewTasksRepo(pg),
		Points: pgdb.NewPointsRepo(pg),
	}
}
