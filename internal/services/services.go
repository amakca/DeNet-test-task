package services

import (
	"denet-test-task/internal/repo"
	"denet-test-task/internal/services/contracts"
	"denet-test-task/internal/services/users"
	"denet-test-task/pkg/hasher"
	"time"
)

type Services struct {
	Auth contracts.Auth
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
		Auth: users.NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
	}
}
