package models

import "time"

type API struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Runtime     string    `json:"runtime"` // "python", "go", "nodejs"
	Visibility  string    `json:"visibility"` // "public", "private", "paid"
	Status      string    `json:"status"` // "pending", "deployed", "failed"
	Endpoint    string    `json:"endpoint"` // Generated endpoint URL
	CodePath    string    `json:"code_path"` // Path to uploaded code
	ContainerID string    `json:"container_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type APIKey struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	APIID     string    `json:"api_id,omitempty"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
