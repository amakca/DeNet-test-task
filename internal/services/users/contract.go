package users

import (
	"context"
	"denet-test-task/internal/entity"
)

type UsersGetInfoInput struct {
	UserId int
}

type UsersSetReferrerInput struct {
	UserId   int
	Referrer string
}

type UsersSetEmailInput struct {
	UserId int
	Email  string
}

type UsersCompleteTaskInput struct {
	UserId int
	TaskId int
	Points int
}

type UsersGetHistoryInput struct {
	UserId int
	Limit  int
}

type UsersGetPointsInput struct {
	UserId int
}

type UsersGetLeaderboardInput struct {
	Limit int
}

type Users interface {
	GetInfo(ctx context.Context, input UsersGetInfoInput) (entity.User, error)
	SetReferrer(ctx context.Context, input UsersSetReferrerInput) error
	SetEmail(ctx context.Context, input UsersSetEmailInput) error
	CompleteTask(ctx context.Context, input UsersCompleteTaskInput) error
	GetHistory(ctx context.Context, input UsersGetHistoryInput) ([]entity.Point, error)
	GetPoints(ctx context.Context, input UsersGetPointsInput) (int, error)
	GetLeaderboard(ctx context.Context, input UsersGetLeaderboardInput) ([]entity.Point, error)
}
