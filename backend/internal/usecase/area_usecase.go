package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type AreaUsecase struct {
	repo repository.AreaRepository
}

func NewAreaUsecase(repo repository.AreaRepository) *AreaUsecase {
	return &AreaUsecase{repo: repo}
}

func (u *AreaUsecase) Create(area *domain.Area) error {
	if area.Name == "" || area.BuildingID == "" {
		return errors.New("nama area dan gedung wajib diisi")
	}
	return u.repo.Create(area)
}

func (u *AreaUsecase) ListAll() ([]domain.Area, error) {
	return u.repo.ListAll()
}

func (u *AreaUsecase) Update(area *domain.Area) error {
	return u.repo.Update(area)
}

func (u *AreaUsecase) Deactivate(id string) error {
	return u.repo.Deactivate(id)
}
