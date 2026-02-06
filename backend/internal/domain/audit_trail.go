package domain

import "time"

// AuditTrail menyimpan catatan aktivitas penting di sistem
type AuditTrail struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Action    string                 `json:"action"`
	Entity    string                 `json:"entity"`
	EntityID  string                 `json:"entity_id"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}
