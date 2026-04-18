package model

import "time"

type CheckIn struct {
	ID        uint64    `json:"id"`
	TaskID    uint64    `json:"task_id"`
	UserID    uint64    `json:"user_id"`
	Points    int       `json:"points"`
	CheckTime time.Time `json:"check_time"`
}

type CheckInRequest struct {
	TaskID uint64 `json:"task_id"`
	UserID uint64 `json:"user_id"`
}