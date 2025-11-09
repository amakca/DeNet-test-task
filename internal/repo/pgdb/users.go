package pgdb

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo/repoerrs"
	"denet-test-task/pkg/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

type UsersRepo struct {
	*postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("username", "password").
		Values(user.Username, user.Password).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("UsersRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *UsersRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("id, username, password, created_at, referrer, email").
		From("users").
		Where("username = ? AND password = ?", username, password).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.Referrer,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UsersRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UsersRepo) GetUserById(ctx context.Context, id int) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("id, username, password, created_at, referrer, email").
		From("users").
		Where("id = ?", id).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.Referrer,
		&user.Email,
		&user.Referrer,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UsersRepo.GetUserById - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UsersRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("id, username, password, created_at, referrer, email").
		From("users").
		Where("username = ?", username).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.Referrer,
		&user.Email,
		&user.Referrer,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UsersRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UsersRepo) SetUserReferrer(ctx context.Context, id int, referrer int) error {
	sql, args, _ := r.Builder.
		Update("users").
		Set("referrer", referrer).
		Where("id = ? AND referrer IS NULL", id).
		ToSql()

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UsersRepo.UpdateUserReferrer - r.Pool.Exec: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return repoerrs.ErrNotFound
	}

	return nil
}

func (r *UsersRepo) SetUserEmail(ctx context.Context, id int, email string) error {
	sql, args, _ := r.Builder.
		Update("users").
		Set("email", email).
		Where("id = ?", id).
		ToSql()

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UsersRepo.SetUserEmail - r.Pool.Exec: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return repoerrs.ErrNotFound
	}

	return nil
}
