package domain

import "time"

type Building struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` // ← pointer
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}
