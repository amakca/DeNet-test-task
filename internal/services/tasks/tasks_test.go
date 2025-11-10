package tasks

import (
	"context"
	"denet-test-task/internal/entity"
	"denet-test-task/internal/repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestTasksService_GetAllTasks(t *testing.T) {
	r := &mockTasksRepo{
		allTasks: []entity.Task{
			{Id: 1, Name: "A", Points: 10},
			{Id: 2, Name: "B", Points: 20},
		},
	}
	s := NewTasksService(r)
	got, err := s.GetAllTasks(context.Background())
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "A", got[0].Name)
	assert.Equal(t, 20, got[1].Points)
}
