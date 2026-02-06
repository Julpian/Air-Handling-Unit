package repository

import "ahu-backend/internal/domain"

type BuildingRepository interface {
	Create(building *domain.Building, createdBy string) error
	ListAll() ([]domain.Building, error)
	GetByID(id string) (*domain.Building, error) // 👈 INI
	Update(building *domain.Building) error
	Deactivate(id string) error
}
