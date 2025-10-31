package models

import "time"

type Usage struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"` // Developer
	APIID        string    `json:"api_id"`
	RequestCount int64     `json:"request_count"`
	TotalRevenue float64   `json:"total_revenue"`
	Date         time.Time `json:"date"`
}

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	APIID       string    `json:"api_id,omitempty"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // "charge", "payout", "refund"
	Status      string    `json:"status"` // "pending", "completed", "failed"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Subscription struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"` // Consumer
	APIID      string    `json:"api_id"`
	Plan       string    `json:"plan"` // "free", "basic", "premium"
	Status     string    `json:"status"` // "active", "cancelled", "expired"
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
