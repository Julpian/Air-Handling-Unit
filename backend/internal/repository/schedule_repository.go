package repository

import "ahu-backend/internal/domain"

type ScheduleRepository interface {
	Create(*domain.Schedule) error
	Update(*domain.Schedule) error
	UpdateStatus(scheduleID string, status string) error // ⬅️ TAMBAHKAN
	DeleteByPlan(planID string) error

	GetByID(id string) (*domain.Schedule, error)

	AssignInspector(scheduleID, inspectorID string) error
	ListAll() ([]domain.Schedule, error)
	ListWithDetail() ([]*domain.ScheduleWithDetail, error)
	GetActiveByAHU(ahuID string) (*domain.Schedule, error)
	GetByYear(year int) ([]domain.Schedule, error)
	ListWithDetailByYear(year int) ([]*domain.ScheduleWithDetail, error)
	SetPDFHash(year int, hash string) error
	ListByYear(year int) ([]*domain.ScheduleWithDetail, error)
	ListTasksByInspector(inspectorID string) ([]domain.InspectionTask, error)
}
