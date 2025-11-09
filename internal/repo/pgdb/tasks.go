package pgdb

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo/repoerrs"
	"denet-test-task/pkg/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TasksRepo struct {
	*postgres.Postgres
}

func NewTasksRepo(pg *postgres.Postgres) *TasksRepo {
	return &TasksRepo{pg}
}

func (r *TasksRepo) GetTaskById(ctx context.Context, id int) (entity.Task, error) {
	sql, args, _ := r.Builder.
		Select("id, name, descr").
		From("tasks").
		Where("id = ?", id).
		ToSql()

	var task entity.Task
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&task.Id,
		&task.Name,
		&task.Descr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, repoerrs.ErrNotFound
		}
		return entity.Task{}, fmt.Errorf("TasksRepo.GetTaskById - r.Pool.QueryRow: %v", err)
	}

	return task, nil
}

func (r *TasksRepo) GetTaskByName(ctx context.Context, name string) (entity.Task, error) {
	sql, args, _ := r.Builder.
		Select("id, name, descr").
		From("tasks").
		Where("name = ?", name).
		ToSql()

	var task entity.Task
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&task.Id,
		&task.Name,
		&task.Descr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, repoerrs.ErrNotFound
		}
		return entity.Task{}, fmt.Errorf("TasksRepo.GetTaskByName - r.Pool.QueryRow: %v", err)
	}

	return task, nil
}

func (r *TasksRepo) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	sql, args, _ := r.Builder.
		Select("id, name, descr").
		From("tasks").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TasksRepo.GetAllTasks - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Task])
	if err != nil {
		return nil, fmt.Errorf("TasksRepo.GetAllTasks - pgx.CollectRows: %v", err)
	}

	return tasks, nil
}
