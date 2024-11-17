package domain

import (
	"database/sql"
	"time"
)

type Task struct {
	ID        uint         `json:"ID"`
	Name      string       `json:"name"`
	FaName    string       `json:"fa_name"`
	Done      bool         `json:"done"`
	CreatedAt time.Time    `json:"created_at"` // Timestamp of record creation
	UpdatedAt time.Time    `json:"updated_at"` // Timestamp of the last record update
	DeletedAt sql.NullTime `json:"deleted_at"` // Nullable timestamp of record deletion
}
type TaskResult struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	FaName string `json:"fa_name"`
	Done   bool   `json:"done"`
}

type TaskRequest struct {
	Name string `json:"name" binding:"required"`
	Done bool   `json:"done" binding:"omitempty,oneof=false"`
}
