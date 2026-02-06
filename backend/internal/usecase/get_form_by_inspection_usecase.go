package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"errors"
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

func (uc *GetFormByInspectionUsecase) Execute(
	inspectionID string,
) (*domain.FormTemplate, error) {

	inspection, err := uc.inspectionRepo.GetByID(inspectionID)
	if err != nil || inspection == nil {
		return nil, errors.New("inspection tidak ditemukan")
	}

	if inspection.FormTemplateID == "" {
		return nil, errors.New("form template belum terpasang")
	}

	return uc.formRepo.GetTemplateByID(inspection.FormTemplateID)
}
