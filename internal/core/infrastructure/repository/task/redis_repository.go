package task

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"hexagonal/internal/core/domain"
	"hexagonal/internal/core/infrastructure/db"
	"time"
)

type RedisTaskRepository struct {
	db *db.Registry
}

func NewRedisTaskRepository(db *db.Registry) *RedisTaskRepository {
	return &RedisTaskRepository{db: db}
}

// Save stores a taskPO in Redis with a context that handles timeouts or cancellations.
func (r *RedisTaskRepository) Save(ctx context.Context, task domain.TaskRequest) error {
	// Marshal the taskPO to JSON format
	data, err := json.Marshal(task)
	if err != nil {
		logrus.Errorf("Failed to marshal taskPO: %v", err)
		return err
	}

	// Set the taskPO data in Redis using the taskPO ID as the key
	err = r.db.GetRedisClient().Set(ctx, "ehsan", data, 10*time.Minute).Err()
	if err != nil {
		logrus.Errorf("Failed to save taskPO in Redis: %v", err)
		return err
	}

	logrus.Infof("Task saved successfully with key: %v", task)
	return nil
}

func (r *RedisTaskRepository) FindAll(ctx context.Context) ([]domain.Task, error) {
	var tasks []domain.Task

	// Example key, adjust as needed
	data, err := r.db.GetRedisClient().Get(ctx, "ehsan").Result()
	if err != nil {
		logrus.Errorf("Failed to get data from Redis: %v", err)
		return nil, err
	}

	logrus.Infof("Data fetched from Redis: %s", data)

	// Unmarshal the Redis JSON data into the tasks slice
	err = json.Unmarshal([]byte(data), &tasks)
	if err != nil {
		logrus.Errorf("Failed to unmarshal Redis data: %v", err)
		return nil, err
	}

	return tasks, nil
}
func (r *RedisTaskRepository) Update(ctx context.Context, taskID string, updatedTask domain.TaskRequest) error {
	// Marshal the updated task to JSON format
	data, err := json.Marshal(updatedTask)
	if err != nil {
		logrus.Errorf("Failed to marshal updated task: %v", err)
		return err
	}

	// Update the task data in Redis using the task ID as the key
	err = r.db.GetRedisClient().Set(ctx, taskID, data, 10*time.Minute).Err()
	if err != nil {
		logrus.Errorf("Failed to update task in Redis: %v", err)
		return err
	}

	logrus.Infof("Task updated successfully with key: %v", taskID)
	return nil
}
func (r *RedisTaskRepository) Delete(ctx context.Context, taskID string) error {
	// Delete the task from Redis by its key (task ID)
	err := r.db.GetRedisClient().Del(ctx, taskID).Err()
	if err != nil {
		logrus.Errorf("Failed to delete task from Redis: %v", err)
		return err
	}

	logrus.Infof("Task with ID %v deleted successfully", taskID)
	return nil
}
