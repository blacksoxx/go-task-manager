package models

import (
    "time"
)

type TaskStatus string

const (
    StatusPending   TaskStatus = "pending"
    StatusInProgress TaskStatus = "in_progress"
    StatusCompleted TaskStatus = "completed"
)

type Task struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    UserID      string     `json:"user_id"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
    Title       string     `json:"title"`
    Description string     `json:"description"`
    UserID      string     `json:"user_id"`
    DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
    Title       string     `json:"title,omitempty"`
    Description string     `json:"description,omitempty"`
    Status      TaskStatus `json:"status,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
}

type TaskResponse struct {
    Task *Task `json:"task"`
}

type TasksResponse struct {
    Tasks []*Task `json:"tasks"`
    Total int     `json:"total"`
}
