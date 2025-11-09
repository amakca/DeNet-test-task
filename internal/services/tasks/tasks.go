package tasks

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo"
)

var _ Tasks = (*TasksService)(nil)

type TasksService struct {
	tasksRepo repo.Tasks
}

func NewTasksService(tasksRepo repo.Tasks) *TasksService {
	return &TasksService{tasksRepo: tasksRepo}
}

func (s *TasksService) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	return s.tasksRepo.GetAllTasks(ctx)
}
