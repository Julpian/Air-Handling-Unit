package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"ahu-backend/internal/usecase/dto"
)

type InspectionQueryUsecase struct {
	inspectionRepo repository.InspectionRepository
	userRepo       repository.UserRepository
}

func NewInspectionQueryUsecase(
	inspectionRepo repository.InspectionRepository,
	userRepo repository.UserRepository,
) *InspectionQueryUsecase {
	return &InspectionQueryUsecase{
		inspectionRepo: inspectionRepo,
		userRepo:       userRepo,
	}
}

// List inspection berdasarkan status
func (u *InspectionQueryUsecase) ListByStatus(
	status string,
) ([]domain.Inspection, error) {

	// kalau status kosong, bisa kamu atur mau ambil semua / error
	if status == "" {
		return u.inspectionRepo.ListByStatus("sedang_diisi")
	}

	return u.inspectionRepo.ListByStatus(status)
}

func (uc *InspectionQueryUsecase) ListDropdown() ([]*dto.InspectorDropdownDTO, error) {
	users, err := uc.userRepo.ListInspectors()
	if err != nil {
		return nil, err
	}

	result := make([]*dto.InspectorDropdownDTO, 0, len(users))

	for _, user := range users {
		result = append(result, &dto.InspectorDropdownDTO{
			ID:   user.ID,
			Name: user.Name,
		})
	}

	return result, nil
}
