package domain

import "time"

type AuditTrailView struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`

	Action   string `json:"action"`
	Entity   string `json:"entity"`
	EntityID string `json:"entity_id"`

	Metadata map[string]interface{} `json:"metadata"`
}
