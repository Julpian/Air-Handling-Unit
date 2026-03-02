package domain

import "time"

type InspectionTask struct {
	ID        string
	StartDate time.Time
	EndDate   time.Time
	Status    string

	UnitCode string
	Period   string
}
