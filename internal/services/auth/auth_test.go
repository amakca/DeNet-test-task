package auth

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo/repoerrs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockUsersRepo struct {
	createdUser   entity.User
	createUserID  int
	createErr     error
	getUserResp   entity.User
	getUserErr    error
	lastGetUser   struct{ Username, Password string }
	getByIDMap    map[int]entity.User
	getByIDErr    error
	setReferrerID struct {
		ID, Referrer int
	}
	setReferrerErr error
	setEmailID     struct {
		ID    int
		Email string
	}
	setEmailErr error
}

func (m *mockUsersRepo) CreateUser(_ context.Context, user entity.User) (int, error) {
	m.createdUser = user
	if m.createErr != nil {
		return 0, m.createErr
	}
	if m.createUserID == 0 {
		m.createUserID = 1
	}
	return m.createUserID, nil
}
func (m *mockUsersRepo) GetUserByUsernameAndPassword(_ context.Context, username, password string) (entity.User, error) {
	m.lastGetUser = struct{ Username, Password string }{Username: username, Password: password}
	return m.getUserResp, m.getUserErr
}
func (m *mockUsersRepo) GetUserById(_ context.Context, id int) (entity.User, error) {
	if m.getByIDErr != nil {
		return entity.User{}, m.getByIDErr
	}
	if u, ok := m.getByIDMap[id]; ok {
		return u, nil
	}
	return entity.User{}, repoerrs.ErrNotFound
}
func (m *mockUsersRepo) GetUserByUsername(_ context.Context, _ string) (entity.User, error) {
	return entity.User{}, repoerrs.ErrNotFound
}
func (m *mockUsersRepo) SetUserReferrer(_ context.Context, id int, referrer int) error {
	m.setReferrerID = struct {
		ID, Referrer int
	}{ID: id, Referrer: referrer}
	return m.setReferrerErr
}
func (m *mockUsersRepo) SetUserEmail(_ context.Context, id int, email string) error {
	m.setEmailID = struct {
		ID    int
		Email string
	}{ID: id, Email: email}
	return m.setEmailErr
}

type mockHasher struct {
	out string
}

func (m mockHasher) Hash(password string) string {
	_ = password
	return m.out
}

func TestAuthService_CreateUser_Success(t *testing.T) {
	repoMock := &mockUsersRepo{}
	h := mockHasher{out: "HPASS"}
	svc := NewAuthService(repoMock, h, "key", time.Hour)

	id, err := svc.CreateUser(context.Background(), AuthCreateUserInput{
		Username: "john", Password: "secret",
	})
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.Equal(t, "john", repoMock.createdUser.Username)
	assert.Equal(t, "HPASS", repoMock.createdUser.Password)
}

func TestAuthService_CreateUser_AlreadyExists(t *testing.T) {
	repoMock := &mockUsersRepo{createErr: repoerrs.ErrAlreadyExists}
	svc := NewAuthService(repoMock, mockHasher{out: "x"}, "key", time.Hour)
	_, err := svc.CreateUser(context.Background(), AuthCreateUserInput{Username: "a", Password: "b"})
	assert.ErrorIs(t, err, ErrUserAlreadyExists)
}

func TestAuthService_CreateUser_InternalError(t *testing.T) {
	repoMock := &mockUsersRepo{createErr: errors.New("db down")}
	svc := NewAuthService(repoMock, mockHasher{out: "x"}, "key", time.Hour)
	_, err := svc.CreateUser(context.Background(), AuthCreateUserInput{Username: "a", Password: "b"})
	assert.ErrorIs(t, err, ErrCannotCreateUser)
}

func TestAuthService_GenerateToken_Errors(t *testing.T) {
	// not found
	repoMock := &mockUsersRepo{getUserErr: repoerrs.ErrNotFound}
	svc := NewAuthService(repoMock, mockHasher{out: "HPASS"}, "key", time.Hour)
	_, err := svc.GenerateToken(context.Background(), AuthGenerateTokenInput{Username: "a", Password: "b"})
	assert.ErrorIs(t, err, ErrUserNotFound)
	// internal
	repoMock2 := &mockUsersRepo{getUserErr: errors.New("db")}
	svc2 := NewAuthService(repoMock2, mockHasher{out: "HPASS"}, "key", time.Hour)
	_, err = svc2.GenerateToken(context.Background(), AuthGenerateTokenInput{Username: "a", Password: "b"})
	assert.ErrorIs(t, err, ErrCannotGetUser)
}

func TestAuthService_GenerateAndParseToken_Success(t *testing.T) {
	repoMock := &mockUsersRepo{
		getUserResp: entity.User{Id: 42, Username: "john"},
	}
	svc := NewAuthService(repoMock, mockHasher{out: "HPASS"}, "secret-key", time.Hour)
	token, err := svc.GenerateToken(context.Background(), AuthGenerateTokenInput{Username: "john", Password: "pwd"})
	assert.NoError(t, err)
	id, err := svc.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, 42, id)
	assert.Equal(t, "john", repoMock.lastGetUser.Username)
	assert.Equal(t, "HPASS", repoMock.lastGetUser.Password)
}

func TestAuthService_ParseToken_Invalid(t *testing.T) {
	repoMock := &mockUsersRepo{}
	svc := NewAuthService(repoMock, mockHasher{out: "x"}, "secret", time.Hour)
	_, err := svc.ParseToken("not-a-token")
	assert.ErrorIs(t, err, ErrCannotParseToken)
}
