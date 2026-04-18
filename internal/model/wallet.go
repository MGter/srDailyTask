package model

import "time"

type WalletType string

const (
	WalletEarn  WalletType = "earn"
	WalletSpend WalletType = "spend"
)

type Wallet struct {
	ID          uint64     `json:"id"`
	UserID      uint64     `json:"user_id"`
	Balance     int        `json:"balance"`
	Type        WalletType `json:"type"`
	Amount      int        `json:"amount"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	RecordTime  time.Time  `json:"record_time"`
}

type SpendRequest struct {
	UserID      uint64 `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type AddRecordRequest struct {
	UserID      uint64     `json:"user_id"`
	Type        WalletType `json:"type"`
	Amount      int        `json:"amount"`
	Description string     `json:"description"`
	RecordTime  time.Time  `json:"record_time"`
}

type DeleteRecordRequest struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
}