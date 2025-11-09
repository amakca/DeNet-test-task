package pgdb

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/pkg/postgres"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type PointsRepo struct {
	*postgres.Postgres
}

func NewPointsRepo(pg *postgres.Postgres) *PointsRepo {
	return &PointsRepo{pg}
}

func (r *PointsRepo) AddPointsByUserId(ctx context.Context, userId int, taskId int, points int) error {

	subquery := r.Builder.
		Select("COALESCE(points, 0)").
		From("points").
		Where("user_id = ?", userId).
		OrderBy("upd_at DESC").
		Limit(1)

	sql, args, _ := r.Builder.
		Insert("points").
		Columns("user_id", "task_id", "points").
		Values(
			userId,
			taskId,
			squirrel.Expr("? + ?", subquery, points),
		).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PointsRepo.AddPointsByUserId - r.Pool.Exec: %v", err)
	}
	return nil
}

func (r *PointsRepo) GetHistoryByUserId(ctx context.Context, userId int) ([]entity.Point, error) {
	sql, args, _ := r.Builder.
		Select("user_id, task_id, points, upd_at").
		From("points").
		Where("user_id = ?", userId).
		OrderBy("upd_at DESC").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PointsRepo.GetHistoryByUserId - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	pointsHistory, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Point])
	if err != nil {
		return nil, fmt.Errorf("PointsRepo.GetHistoryByUserId - pgx.CollectRows: %v", err)
	}
	return pointsHistory, nil
}

func (r *PointsRepo) CheckCompletedTask(ctx context.Context, userId int, taskId int) (bool, error) {
	sql, args, _ := r.Builder.
		Select("count(1)").
		From("points").
		Where("user_id = ? AND task_id = ?", userId, taskId).
		ToSql()

	var cnt int
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&cnt); err != nil {
		return false, fmt.Errorf("PointsRepo.CheckCompletedTask - r.Pool.QueryRow: %v", err)
	}
	return cnt > 0, nil
}

func (r *PointsRepo) GetPointsByUserId(ctx context.Context, userId int) (int, error) {

	sql, args, _ := r.Builder.
		Select("COALESCE(SUM(points), 0)").
		From("points").
		Where("user_id = ?", userId).
		ToSql()

	var points int
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&points); err != nil {
		return 0, fmt.Errorf("PointsRepo.GetPointsByUserId - r.Pool.QueryRow: %v", err)
	}
	return points, nil
}

func (r *PointsRepo) GetLeaderboard(ctx context.Context, limit int) ([]entity.Point, error) {
	sql, args, _ := r.Builder.
		Select("user_id, SUM(points) as points").
		From("points").
		GroupBy("user_id").
		OrderBy("points DESC").
		Limit(uint64(limit)).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PointsRepo.GetLeaderboard - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	leaderboard, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Point])
	if err != nil {
		return nil, fmt.Errorf("PointsRepo.GetLeaderboard - pgx.CollectRows: %v", err)
	}
	return leaderboard, nil
}
