package task

import (
	"context"
	"database/sql"
	"hexagonal/internal/core/domain"
	"hexagonal/internal/core/infrastructure/db"
	"log"
)

// PostgresTaskRepository implements port.TaskRepository for PostgreSQL
type PostgresTaskRepository struct {
	db *db.Registry
}

// NewPostgresTaskRepository creates a new PostgresTaskRepository
func NewPostgresTaskRepository(db *db.Registry) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

// Save saves a taskPO to the PostgreSQL dbPO
func (r *PostgresTaskRepository) Save(ctx context.Context, task domain.TaskRequest) error {
	_, err := r.db.GetPostgres().QueryContext(ctx, "INSERT INTO tasks (name, done) VALUES ($1, $2)", task.Name, task.Done)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all tasks from the PostgreSQL dbPO
func (r *PostgresTaskRepository) FindAll(ctx context.Context) ([]domain.Task, error) {
	var tasks []domain.Task

	// Use Query to retrieve rows
	rows, err := r.db.GetPostgres().QueryContext(ctx, "SELECT id, name, fa_name, done, created_at, updated_at, deleted_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(rows)

	// Iterate over rows and populate the tasks slice
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.FaName, &task.Done, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
