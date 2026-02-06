package repository

import "ahu-backend/internal/domain"

type AreaRepository interface {
	Create(area *domain.Area) error
	ListAll() ([]domain.Area, error)
	GetByID(id string) (*domain.Area, error)
	Update(area *domain.Area) error
	Deactivate(id string) error
}
