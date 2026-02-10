package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type GetFormByInspectionUsecase struct {
	inspectionRepo repository.InspectionRepository
	formRepo       repository.FormRepository
}

func NewGetFormByInspectionUsecase(
	inspectionRepo repository.InspectionRepository,
	formRepo repository.FormRepository,
) *GetFormByInspectionUsecase {
	return &GetFormByInspectionUsecase{
		inspectionRepo: inspectionRepo,
		formRepo:       formRepo,
	}
}

func (u *GetFormByInspectionUsecase) Execute(
	inspectionID string,
) (*domain.FormTemplate, error) {

	inspection, err := u.inspectionRepo.GetByID(inspectionID)
	if err != nil || inspection == nil {
		return nil, errors.New("inspection tidak ditemukan")
	}

	// 🔥 AMBIL FORM BERDASARKAN ID YANG SUDAH TERKUNCI
	return u.formRepo.GetTemplateByID(inspection.FormTemplateID)
}
