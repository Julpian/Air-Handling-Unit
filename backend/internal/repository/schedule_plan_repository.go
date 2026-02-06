package repository

import "ahu-backend/internal/domain"

type SchedulePlanRepository interface {
	Create(plan *domain.SchedulePlan) error
	Update(plan *domain.SchedulePlan) error
	Delete(id string) error

	ListAllWithAHU() ([]domain.SchedulePlanWithAHU, error)
	GetActiveByYear(year int) ([]domain.SchedulePlan, error)
	MarkGenerated(planID string) error
}
