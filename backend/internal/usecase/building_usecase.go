package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

type BuildingUsecase struct {
	repo repository.BuildingRepository
}

func NewBuildingUsecase(repo repository.BuildingRepository) *BuildingUsecase {
	return &BuildingUsecase{repo: repo}
}

func (u *BuildingUsecase) Create(
	name string,
	description *string,
	createdBy string,
) error {

	b := &domain.Building{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
	}

	return u.repo.Create(b, createdBy)
}

func (u *BuildingUsecase) List() ([]domain.Building, error) {
	return u.repo.ListAll()
}

func (u *BuildingUsecase) Update(
	id string,
	name string,
	description *string,
) error {

	b := &domain.Building{
		ID:          id,
		Name:        name,
		Description: description,
	}

	return u.repo.Update(b)
}

func (u *BuildingUsecase) Deactivate(id string) error {
	return u.repo.Deactivate(id)
}
