package auth

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo"
	"denet-test-task/internal/repo/repoerrs"
	"denet-test-task/pkg/hasher"
	"denet-test-task/pkg/logctx"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var _ Auth = (*AuthService)(nil)

var (
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthService struct {
	usersRepo      repo.Users
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAuthService(usersRepo repo.Users, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		usersRepo:      usersRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error) {
	user := entity.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userId, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		logctx.FromContext(ctx).Error("AuthService.CreateUser - userRepo.CreateUser", "err", err)
		return 0, ErrCannotCreateUser
	}
	return userId, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	// get user from DB
	user, err := s.usersRepo.GetUserByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrUserNotFound
		}
		logctx.FromContext(ctx).Error("AuthService.GenerateToken: cannot get user", "err", err)
		return "", ErrCannotGetUser
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		logctx.FromContext(ctx).Error("AuthService.GenerateToken: cannot sign token", "err", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrCannotParseToken
	}

	return claims.UserId, nil
}
