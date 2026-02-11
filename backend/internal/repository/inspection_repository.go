package repository

import (
	"ahu-backend/internal/domain"
	"time"
)

type InspectionRepository interface {
	Create(i *domain.Inspection) error

	GetByID(id string) (*domain.Inspection, error)
	SetScanToken(id string, token string, expires time.Time, uid string) error
	GetLastByScheduleID(scheduleID string) (*domain.Inspection, error)

	UpdateStatus(id string, status string, note *string) error

	ExistsApproved(scheduleID string) (bool, error)
	ListByStatus(status string) ([]domain.Inspection, error)
	SaveResult(result *domain.InspectionResult) error
	ClearScanToken(id string) error
	GetByScheduleID(scheduleID string) (*domain.Inspection, error)

	Approve(
		inspectionID string,
		approverID string,
		approvedAt any,
		metadata map[string]any,
	) error
}
