package taskSRV

import (
	"context"
	"github.com/sirupsen/logrus"
	"hexagonal/internal/core/domain"
	"hexagonal/internal/core/port/taskPO"
)

type TaskService struct {
	repo taskPO.TaskRepository
}

func NewTaskService(repo taskPO.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, task domain.TaskRequest) error {
	// Delegate the task creation to the repository's Save method
	return s.repo.Save(ctx, task)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	// Fetch all tasks from the repository
	return s.repo.FindAll(ctx)
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID string, updatedTask domain.TaskRequest) error {

	// Call the repository's Update method to persist the changes
	err := s.repo.Update(ctx, taskID, updatedTask)
	if err != nil {
		logrus.Errorf("Failed to update task: %v", err)
		return err
	}

	logrus.Infof("Task updated successfully with ID: %v", taskID)
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	// Call the repository's Delete method to remove the task
	err := s.repo.Delete(ctx, taskID)
	if err != nil {
		logrus.Errorf("Failed to delete task: %v", err)
		return err
	}

	logrus.Infof("Task with ID %v deleted successfully", taskID)
	return nil
}
