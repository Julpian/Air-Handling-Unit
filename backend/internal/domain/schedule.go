package domain

import "time"

type Schedule struct {
	ID     string
	PlanID string
	AHUId  string

	StartDate time.Time
	EndDate   time.Time

	InspectorID *string
	Status      string
	NFCBypass   bool
	CreatedAt   time.Time
}
