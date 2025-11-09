package users

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo"
)

var _ Users = (*UsersService)(nil)

type UsersService struct {
	userRepo  repo.Users
	pointRepo repo.Points
}

func NewUsersService(userRepo repo.Users, pointRepo repo.Points) *UsersService {
	return &UsersService{userRepo: userRepo, pointRepo: pointRepo}
}

func (s *UsersService) GetLeaderboard(ctx context.Context, input UsersGetLeaderboardInput) ([]entity.Point, error) {
	return s.pointRepo.GetLeaderboard(ctx, input.Limit)
}

func (s *UsersService) GetInfo(ctx context.Context, input UsersGetInfoInput) (entity.User, error) {
	return s.userRepo.GetUserById(ctx, input.UserId)
}

func (s *UsersService) SetEmail(ctx context.Context, input UsersSetEmailInput) error {
	return s.userRepo.SetUserEmail(ctx, input.UserId, input.Email)
}

func (s *UsersService) SetReferrer(ctx context.Context, input UsersSetReferrerInput) error {
	return s.userRepo.SetUserReferrer(ctx, input.UserId, input.Referrer)
}

func (s *UsersService) CompleteTask(ctx context.Context, input UsersCompleteTaskInput) error {
	return s.pointRepo.AddPointsByUserId(ctx, input.UserId, input.TaskId, input.Points)
}

func (s *UsersService) GetHistory(ctx context.Context, input UsersGetHistoryInput) ([]entity.Point, error) {
	return s.pointRepo.GetHistoryByUserId(ctx, input.UserId)
}

func (s *UsersService) GetPoints(ctx context.Context, input UsersGetPointsInput) (int, error) {
	return s.pointRepo.GetPointsByUserId(ctx, input.UserId)
}
