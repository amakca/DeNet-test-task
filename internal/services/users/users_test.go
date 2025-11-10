package users

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUsersRepo struct {
	usersByID      map[int]entity.User
	getByIDErr     error
	setRefUserID   int
	setRefReferrer int
	setReferrerErr error
	setEmailUserID int
	setEmailEmail  string
	setEmailErr    error
}

func (m *mockUsersRepo) CreateUser(_ context.Context, _ entity.User) (int, error) {
	return 0, nil
}
func (m *mockUsersRepo) GetUserByUsernameAndPassword(_ context.Context, _, _ string) (entity.User, error) {
	return entity.User{}, nil
}
func (m *mockUsersRepo) GetUserById(_ context.Context, id int) (entity.User, error) {
	if m.getByIDErr != nil {
		return entity.User{}, m.getByIDErr
	}
	if u, ok := m.usersByID[id]; ok {
		return u, nil
	}
	return entity.User{}, errors.New("not found")
}
func (m *mockUsersRepo) GetUserByUsername(_ context.Context, _ string) (entity.User, error) {
	return entity.User{}, nil
}
func (m *mockUsersRepo) SetUserReferrer(_ context.Context, id int, referrer int) error {
	m.setRefUserID = id
	m.setRefReferrer = referrer
	return m.setReferrerErr
}
func (m *mockUsersRepo) SetUserEmail(_ context.Context, id int, email string) error {
	m.setEmailUserID = id
	m.setEmailEmail = email
	return m.setEmailErr
}

var _ repo.Users = (*mockUsersRepo)(nil)

type mockPointsRepo struct {
	addCalls           []struct{ UserID, TaskID, Points int }
	addErr             error
	checkCompletedResp bool
	checkCompletedErr  error
	leaderboardResp    []entity.LeaderboardItem
	leaderboardErr     error
	historyResp        []entity.Point
	historyErr         error
	pointsByUserResp   int
	pointsByUserErr    error
}

func (m *mockPointsRepo) AddPointsByUserId(_ context.Context, userId int, taskId int, points int) error {
	m.addCalls = append(m.addCalls, struct{ UserID, TaskID, Points int }{userId, taskId, points})
	return m.addErr
}
func (m *mockPointsRepo) GetPointsByUserId(_ context.Context, _ int) (int, error) {
	return m.pointsByUserResp, m.pointsByUserErr
}
func (m *mockPointsRepo) GetHistoryByUserId(_ context.Context, _ int) ([]entity.Point, error) {
	return m.historyResp, m.historyErr
}
func (m *mockPointsRepo) CheckCompletedTask(_ context.Context, _ int, _ int) (bool, error) {
	return m.checkCompletedResp, m.checkCompletedErr
}
func (m *mockPointsRepo) GetLeaderboard(_ context.Context, _ int) ([]entity.LeaderboardItem, error) {
	return m.leaderboardResp, m.leaderboardErr
}

var _ repo.Points = (*mockPointsRepo)(nil)

type mockTasksRepo struct {
	allTasks []entity.Task
	err      error
}

func (m *mockTasksRepo) GetTaskById(_ context.Context, _ int) (entity.Task, error) {
	if len(m.allTasks) > 0 {
		return m.allTasks[0], m.err
	}
	return entity.Task{}, m.err
}
func (m *mockTasksRepo) GetTaskByName(_ context.Context, _ string) (entity.Task, error) {
	if len(m.allTasks) > 0 {
		return m.allTasks[0], m.err
	}
	return entity.Task{}, m.err
}
func (m *mockTasksRepo) GetAllTasks(_ context.Context) ([]entity.Task, error) {
	return m.allTasks, m.err
}

var _ repo.Tasks = (*mockTasksRepo)(nil)

func TestNewUsersService_ErrorOnTasksFetch(t *testing.T) {
	_, err := NewUsersService(context.Background(), &mockUsersRepo{}, &mockPointsRepo{}, &mockTasksRepo{err: errors.New("boom")})
	assert.ErrorIs(t, err, ErrCannotGetTasks)
}

func TestUsersService_CompleteTask_Restricted(t *testing.T) {
	svc, err := NewUsersService(context.Background(),
		&mockUsersRepo{},
		&mockPointsRepo{},
		&mockTasksRepo{allTasks: []entity.Task{}},
	)
	assert.NoError(t, err)
	for _, restricted := range []int{TaskCompleteEmail, TaskGetReferral, TaskGiveReferral} {
		err := svc.CompleteTask(context.Background(), UsersCompleteTaskInput{UserId: 1, TaskId: restricted})
		assert.ErrorIs(t, err, ErrTaskNotAllowedToComplete)
	}
}

func TestUsersService_CompleteTask_NotFound(t *testing.T) {
	svc, err := NewUsersService(context.Background(),
		&mockUsersRepo{},
		&mockPointsRepo{},
		&mockTasksRepo{allTasks: []entity.Task{}},
	)
	assert.NoError(t, err)
	err = svc.CompleteTask(context.Background(), UsersCompleteTaskInput{UserId: 1, TaskId: 999})
	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestUsersService_CompleteTask_CheckError(t *testing.T) {
	points := &mockPointsRepo{checkCompletedErr: errors.New("db")}
	svc, err := NewUsersService(context.Background(),
		&mockUsersRepo{},
		points,
		&mockTasksRepo{allTasks: []entity.Task{{Id: 100, Points: 15}}},
	)
	assert.NoError(t, err)
	err = svc.CompleteTask(context.Background(), UsersCompleteTaskInput{UserId: 1, TaskId: 100})
	assert.ErrorIs(t, err, ErrCannotCheckCompletedTask)
}

func TestUsersService_CompleteTask_AlreadyCompleted(t *testing.T) {
	points := &mockPointsRepo{checkCompletedResp: true}
	svc, err := NewUsersService(context.Background(),
		&mockUsersRepo{},
		points,
		&mockTasksRepo{allTasks: []entity.Task{{Id: 101, Points: 7}}},
	)
	assert.NoError(t, err)
	err = svc.CompleteTask(context.Background(), UsersCompleteTaskInput{UserId: 1, TaskId: 101})
	assert.ErrorIs(t, err, ErrTaskAlreadyCompleted)
}

func TestUsersService_CompleteTask_AddPointsError(t *testing.T) {
	points := &mockPointsRepo{checkCompletedResp: false, addErr: errors.New("db")}
	svc, err := NewUsersService(context.Background(),
		&mockUsersRepo{},
		points,
		&mockTasksRepo{allTasks: []entity.Task{{Id: 102, Points: 13}}},
	)
	assert.NoError(t, err)
	err = svc.CompleteTask(context.Background(), UsersCompleteTaskInput{UserId: 1, TaskId: 102})
	assert.ErrorIs(t, err, ErrCannotAddPoints)
}

func TestUsersService_CompleteTask_Success(t *testing.T) {
	points := &mockPointsRepo{checkCompletedResp: false}
	svc, err := NewUsersService(context.Background(),
		&mockUsersRepo{},
		points,
		&mockTasksRepo{allTasks: []entity.Task{{Id: 103, Points: 33}}},
	)
	assert.NoError(t, err)
	err = svc.CompleteTask(context.Background(), UsersCompleteTaskInput{UserId: 10, TaskId: 103})
	assert.NoError(t, err)
	assert.Len(t, points.addCalls, 1)
	call := points.addCalls[0]
	assert.Equal(t, 10, call.UserID)
	assert.Equal(t, 103, call.TaskID)
	assert.Equal(t, 33, call.Points)
}

func TestUsersService_SetReferrer_SelfReferrerError(t *testing.T) {
	uRepo := &mockUsersRepo{
		usersByID: map[int]entity.User{
			1: {Id: 1},
			2: {Id: 2, Referrer: strPtr(strconv.Itoa(1))},
		},
	}
	svc, err := NewUsersService(context.Background(), uRepo, &mockPointsRepo{}, &mockTasksRepo{allTasks: []entity.Task{}})
	assert.NoError(t, err)
	err = svc.SetReferrer(context.Background(), UsersSetReferrerInput{UserId: 1, Referrer: 2})
	assert.ErrorIs(t, err, ErrReferrerCannotBeTheSameAsUser)
}

func TestUsersService_SetReferrer_UserAlreadyHasReferrer(t *testing.T) {
	referrer := 2
	refStr := strconv.Itoa(referrer)
	uRepo := &mockUsersRepo{
		usersByID: map[int]entity.User{
			1: {Id: 1, Referrer: &refStr},
			2: {Id: 2},
		},
	}
	svc, err := NewUsersService(context.Background(), uRepo, &mockPointsRepo{}, &mockTasksRepo{allTasks: []entity.Task{}})
	assert.NoError(t, err)
	err = svc.SetReferrer(context.Background(), UsersSetReferrerInput{UserId: 1, Referrer: referrer})
	assert.ErrorIs(t, err, ErrUserAlreadySetReferrer)
}

func TestUsersService_SetReferrer_TaskMappingMissing(t *testing.T) {
	uRepo := &mockUsersRepo{
		usersByID: map[int]entity.User{
			1: {Id: 1},
			2: {Id: 2},
		},
	}
	// No tasks provided -> mapping missing
	svc, err := NewUsersService(context.Background(), uRepo, &mockPointsRepo{}, &mockTasksRepo{allTasks: []entity.Task{}})
	assert.NoError(t, err)
	err = svc.SetReferrer(context.Background(), UsersSetReferrerInput{UserId: 1, Referrer: 2})
	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestUsersService_SetReferrer_Success(t *testing.T) {
	uRepo := &mockUsersRepo{
		usersByID: map[int]entity.User{
			1: {Id: 1},
			2: {Id: 2},
		},
	}
	points := &mockPointsRepo{}
	// Provide both referral tasks
	tasks := []entity.Task{
		{Id: TaskGiveReferral, Points: 5},
		{Id: TaskGetReferral, Points: 7},
	}
	svc, err := NewUsersService(context.Background(), uRepo, points, &mockTasksRepo{allTasks: tasks})
	assert.NoError(t, err)
	err = svc.SetReferrer(context.Background(), UsersSetReferrerInput{UserId: 1, Referrer: 2})
	assert.NoError(t, err)
	assert.Len(t, points.addCalls, 2)
	assert.Equal(t, 1, uRepo.setRefUserID)
	assert.Equal(t, 2, uRepo.setRefReferrer)
}

func TestUsersService_SetEmail(t *testing.T) {
	points := &mockPointsRepo{addErr: nil}
	uRepo := &mockUsersRepo{}
	tasks := []entity.Task{{Id: TaskCompleteEmail, Points: 11}}
	svc, err := NewUsersService(context.Background(), uRepo, points, &mockTasksRepo{allTasks: tasks})
	assert.NoError(t, err)
	err = svc.SetEmail(context.Background(), UsersSetEmailInput{UserId: 99, Email: "x@y.z"})
	assert.NoError(t, err)
	assert.Len(t, points.addCalls, 1)
	assert.Equal(t, 99, points.addCalls[0].UserID)
	assert.Equal(t, TaskCompleteEmail, points.addCalls[0].TaskID)
	assert.Equal(t, 11, points.addCalls[0].Points)
	assert.Equal(t, 99, uRepo.setEmailUserID)
	assert.Equal(t, "x@y.z", uRepo.setEmailEmail)
}

func TestUsersService_SetEmail_AddPointsError(t *testing.T) {
	points := &mockPointsRepo{addErr: errors.New("db")}
	uRepo := &mockUsersRepo{}
	tasks := []entity.Task{{Id: TaskCompleteEmail, Points: 11}}
	svc, err := NewUsersService(context.Background(), uRepo, points, &mockTasksRepo{allTasks: tasks})
	assert.NoError(t, err)
	err = svc.SetEmail(context.Background(), UsersSetEmailInput{UserId: 1, Email: "a@b.c"})
	assert.ErrorIs(t, err, ErrCannotAddPoints)
}

func strPtr(s string) *string { return &s }
