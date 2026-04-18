package model

import "time"

type CircleMode string

const (
	CircleOnce    CircleMode = "once"
	CircleWeekly  CircleMode = "weekly"
	CircleWorkday CircleMode = "workday"
	CircleWeekend CircleMode = "weekend"
	CircleCustom  CircleMode = "custom"
)

type Task struct {
	ID          uint64     `json:"id"`
	UserID      uint64     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CircleMode  CircleMode `json:"circle_mode"`
	Points      int        `json:"points"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	IsExpired   bool       `json:"is_expired"`
}

type CreateTaskRequest struct {
	UserID      uint64     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CircleMode  CircleMode `json:"circle_mode"`
	Points      int        `json:"points"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CircleMode  CircleMode `json:"circle_mode"`
	Points      int        `json:"points"`
	IsExpired   bool       `json:"is_expired"`
}