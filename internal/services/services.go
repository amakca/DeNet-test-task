package services

import (
	"denet-test-task/internal/repo"
	"denet-test-task/internal/services/auth"
	"denet-test-task/internal/services/tasks"
	"denet-test-task/internal/services/users"
	"denet-test-task/pkg/hasher"
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

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:  auth.NewAuthService(deps.Repos.Users, deps.Hasher, deps.SignKey, deps.TokenTTL),
		User:  users.NewUsersService(deps.Repos.Users, deps.Repos.Points),
		Tasks: tasks.NewTasksService(deps.Repos.Tasks),
	}
}
