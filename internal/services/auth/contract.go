package auth

import (
	"context"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}
