package repository

import "ahu-backend/internal/domain"

type AHURepository interface {
	Create(ahu *domain.AHU) error
	ListAll() ([]domain.AHU, error)
	ListAllOrderedByCreatedAt() ([]domain.AHU, error)
	GetByID(id string) (*domain.AHU, error)
	GetByNFCUID(nfcUID string) (*domain.AHU, error)
	Update(ahu *domain.AHU) error
	Deactivate(id string) error
}
