package users

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo"
	"denet-test-task/pkg/logctx"
	"fmt"
	"strconv"
)

var _ Users = (*UsersService)(nil)

var (
	ErrTaskNotFound                  = fmt.Errorf("task not found")
	ErrTaskAlreadyCompleted          = fmt.Errorf("task already completed")
	ErrCannotCheckCompletedTask      = fmt.Errorf("cannot check if task is completed")
	ErrCannotAddPoints               = fmt.Errorf("cannot add points")
	ErrCannotGetTasks                = fmt.Errorf("cannot get tasks")
	ErrUserAlreadySetReferrer        = fmt.Errorf("user already has a referrer")
	ErrTaskNotAllowedToComplete      = fmt.Errorf("task not allowed to complete")
	ErrReferrerCannotBeTheSameAsUser = fmt.Errorf("referrer cannot be the same as user")
)

const (
	TaskGiveReferral = iota + 1
	TaskGetReferral
	TaskSubscribeTelegram
	TaskSubscribeTwitter
	TaskCompleteEmail
)

type UsersService struct {
	usersRepo  repo.Users
	pointsRepo repo.Points
	tasksRepo  repo.Tasks

	tasksList map[int]int // map[task_id]points
}

func NewUsersService(ctx context.Context, userRepo repo.Users, pointRepo repo.Points, tasksRepo repo.Tasks) (*UsersService, error) {

	service := &UsersService{usersRepo: userRepo, pointsRepo: pointRepo, tasksRepo: tasksRepo}

	tasks, err := tasksRepo.GetAllTasks(ctx)
	if err != nil {
		logctx.FromContext(ctx).Error("UsersService.NewUsersService - tasksRepo.GetAllTasks", "err", err)
		return nil, ErrCannotGetTasks
	}

	service.tasksList = make(map[int]int, len(tasks))
	for _, task := range tasks {
		service.tasksList[task.Id] = task.Points
	}
	return service, nil
}

func (s *UsersService) GetLeaderboard(ctx context.Context, input UsersGetLeaderboardInput) ([]entity.LeaderboardItem, error) {
	return s.pointsRepo.GetLeaderboard(ctx, input.Limit)
}

func (s *UsersService) GetInfo(ctx context.Context, input UsersGetInfoInput) (entity.User, error) {
	return s.usersRepo.GetUserById(ctx, input.UserId)
}

func (s *UsersService) SetEmail(ctx context.Context, input UsersSetEmailInput) error {

	pointsForEmail, ok := s.tasksList[TaskCompleteEmail]
	if !ok {
		logctx.FromContext(ctx).Error("UsersService.SetEmail - task not found")
	}

	err := s.pointsRepo.AddPointsByUserId(ctx, input.UserId, TaskCompleteEmail, pointsForEmail)
	if err != nil {
		logctx.FromContext(ctx).Error("UsersService.SetEmail - pointsRepo.AddPointsByUserId", "err", err)
		return ErrCannotAddPoints
	}

	return s.usersRepo.SetUserEmail(ctx, input.UserId, input.Email)
}

func (s *UsersService) SetReferrer(ctx context.Context, input UsersSetReferrerInput) error {

	referrer, err := s.usersRepo.GetUserById(ctx, input.Referrer)
	if err != nil {
		logctx.FromContext(ctx).Error("UsersService.SetReferrer - usersRepo.GetUserById", "err", err)
		return err
	}
	if referrer.Referrer != nil && *referrer.Referrer == strconv.Itoa(input.UserId) {
		logctx.FromContext(ctx).Error("UsersService.SetReferrer - referrer cannot be the same as user")
		return ErrReferrerCannotBeTheSameAsUser
	}

	user, err := s.usersRepo.GetUserById(ctx, input.UserId)
	if err != nil {
		logctx.FromContext(ctx).Error("UsersService.SetReferrer - usersRepo.GetUserById", "err", err)
		return err
	}
	if user.Referrer != nil {
		logctx.FromContext(ctx).Error("UsersService.SetReferrer - user already has a referrer")
		return ErrUserAlreadySetReferrer
	}

	pointsForReferrer, ok := s.tasksList[TaskGiveReferral]
	if !ok {
		logctx.FromContext(ctx).Error("UsersService.SetReferrer - task not found")
		return ErrTaskNotFound
	}
	pointsForUser, ok := s.tasksList[TaskGetReferral]
	if !ok {
		logctx.FromContext(ctx).Error("UsersService.SetReferrer - task not found")
		return ErrTaskNotFound
	}

	s.pointsRepo.AddPointsByUserId(ctx, input.Referrer, TaskGiveReferral, pointsForReferrer)
	s.pointsRepo.AddPointsByUserId(ctx, input.UserId, TaskGetReferral, pointsForUser)
	return s.usersRepo.SetUserReferrer(ctx, input.UserId, input.Referrer)
}

func (s *UsersService) CompleteTask(ctx context.Context, input UsersCompleteTaskInput) error {

	if input.TaskId == TaskCompleteEmail || input.TaskId == TaskGetReferral || input.TaskId == TaskGiveReferral {
		logctx.FromContext(ctx).Error("UsersService.CompleteTask - task not allowed to complete")
		return ErrTaskNotAllowedToComplete
	}

	points, ok := s.tasksList[input.TaskId]
	if !ok {
		logctx.FromContext(ctx).Error("UsersService.CompleteTask - task not found")
		return ErrTaskNotFound
	}

	completed, err := s.pointsRepo.CheckCompletedTask(ctx, input.UserId, input.TaskId)
	if err != nil {
		logctx.FromContext(ctx).Error("UsersService.CompleteTask - pointsRepo.CheckCompletedTask", "err", err)
		return ErrCannotCheckCompletedTask
	}
	if completed {
		logctx.FromContext(ctx).Error("UsersService.CompleteTask - task already completed")
		return ErrTaskAlreadyCompleted
	}

	err = s.pointsRepo.AddPointsByUserId(ctx, input.UserId, input.TaskId, points)
	if err != nil {
		logctx.FromContext(ctx).Error("UsersService.CompleteTask - pointsRepo.AddPointsByUserId", "err", err)
		return ErrCannotAddPoints
	}

	return nil
}

func (s *UsersService) GetHistory(ctx context.Context, input UsersGetHistoryInput) ([]entity.Point, error) {
	return s.pointsRepo.GetHistoryByUserId(ctx, input.UserId)
}

func (s *UsersService) GetPoints(ctx context.Context, input UsersGetPointsInput) (int, error) {
	return s.pointsRepo.GetPointsByUserId(ctx, input.UserId)
}
