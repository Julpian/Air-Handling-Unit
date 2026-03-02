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

// 🔥 FIX: Hanya gunakan versi ini (2 parameter)
func (uc *InspectionQueryUsecase) ListByStatus(status string, inspectorID string) ([]domain.Inspection, error) {
	// FIX TYPO: Gunakan uc.inspectionRepo (sesuai nama di struct)
	return uc.inspectionRepo.ListByStatus(status, inspectorID)
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
