package taskPO

import (
	"context"
	"hexagonal/internal/core/domain"
)

type TaskRepository interface {
	Save(ctx context.Context, task domain.TaskRequest) error
	FindAll(ctx context.Context) ([]domain.Task, error)
	Update(ctx context.Context, taskID string, updatedTask domain.TaskRequest) error
	Delete(ctx context.Context, taskID string) error
}

type TaskService interface {
	CreateTask(ctx context.Context, task domain.TaskRequest) error
	GetAllTasks(ctx context.Context) ([]domain.Task, error)
	UpdateTask(ctx context.Context, taskID string, updatedTask domain.TaskRequest) error
	DeleteTask(ctx context.Context, taskID string) error
}
