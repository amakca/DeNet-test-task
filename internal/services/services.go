package services

import (
	"context"
	"denet-test-task/internal/repo"
	"denet-test-task/internal/services/auth"
	"denet-test-task/internal/services/tasks"
	"denet-test-task/internal/services/users"
	"denet-test-task/pkg/hasher"
	"denet-test-task/pkg/logctx"
	"time"
)

type Services struct {
	Auth  auth.Auth
	User  users.Users
	Tasks tasks.Tasks
}

type ServicesDependencies struct {
	Repos *repo.Repositories
	// GDrive webapi.GDrive
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(ctx context.Context, deps ServicesDependencies) (*Services, error) {

	userService, err := users.NewUsersService(ctx, deps.Repos.Users, deps.Repos.Points, deps.Repos.Tasks)
	if err != nil {
		logctx.FromContext(ctx).Error("Services.NewServices - users.NewUsersService", "err", err)
		return nil, err
	}

	return &Services{
		Auth:  auth.NewAuthService(deps.Repos.Users, deps.Hasher, deps.SignKey, deps.TokenTTL),
		User:  userService,
		Tasks: tasks.NewTasksService(deps.Repos.Tasks),
	}, nil
}
