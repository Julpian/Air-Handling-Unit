package domain

import "time"

type InspectionNote struct {
	ID           string
	InspectionID string
	Note         string
	CreatedAt    time.Time
}
