package model

import "time"

type LongTermItemStatus string

const (
	LongTermItemActive   LongTermItemStatus = "active"
	LongTermItemScrapped LongTermItemStatus = "scrapped"
)

type LongTermItem struct {
	ID              uint64             `json:"id"`
	UserID          uint64             `json:"user_id"`
	Name            string             `json:"name"`
	Price           float64            `json:"price"`
	PurchaseDate    time.Time          `json:"purchase_date"`
	ScrapDate       *time.Time         `json:"scrap_date"`
	FrozenDailyCost *float64           `json:"frozen_daily_cost"`
	Status          LongTermItemStatus `json:"status"`
	DailyCost       float64            `json:"daily_cost"`
	OwnedDays       int                `json:"owned_days"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type CreateLongTermItemRequest struct {
	UserID       uint64    `json:"user_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	PurchaseDate time.Time `json:"purchase_date"`
}

type UpdateLongTermItemRequest struct {
	UserID       uint64    `json:"user_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	PurchaseDate time.Time `json:"purchase_date"`
}

type ScrapLongTermItemRequest struct {
	UserID    uint64    `json:"user_id"`
	ScrapDate time.Time `json:"scrap_date"`
}

type DeleteLongTermItemRequest struct {
	UserID uint64 `json:"user_id"`
}

type LongTermSummary struct {
	ActiveDailyCost float64 `json:"active_daily_cost"`
	ActiveCount     int     `json:"active_count"`
	ScrappedCount   int     `json:"scrapped_count"`
}

type LongTermItemListResponse struct {
	Items   []*LongTermItem  `json:"items"`
	Summary *LongTermSummary `json:"summary"`
}
