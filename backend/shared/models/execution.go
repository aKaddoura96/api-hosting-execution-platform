package models

import "time"

type Execution struct {
	ID             string        `json:"id"`
	APIID          string        `json:"api_id"`
	UserID         string        `json:"user_id"` // Consumer who invoked
	StatusCode     int           `json:"status_code"`
	Duration       time.Duration `json:"duration"` // Execution time in ms
	RequestSize    int64         `json:"request_size"` // Bytes
	ResponseSize   int64         `json:"response_size"` // Bytes
	Error          string        `json:"error,omitempty"`
	ExecutedAt     time.Time     `json:"executed_at"`
}
