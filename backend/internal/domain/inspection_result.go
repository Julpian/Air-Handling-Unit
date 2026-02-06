package domain

import "time"

type InspectionResult struct {
	ID           string
	InspectionID string
	FormItemID   string

	ValueText   *string
	ValueNumber *float64
	ValueBool   *bool

	Result    string // pass | fail
	CreatedAt time.Time
}
