package tasks

import (
	"context"
	"denet-test-task/internal/entity"
)

type Tasks interface {
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}
