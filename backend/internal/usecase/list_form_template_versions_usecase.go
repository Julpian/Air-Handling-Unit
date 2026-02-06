package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type ListFormTemplateVersionsUsecase struct {
	formRepo repository.FormRepository
}

func NewListFormTemplateVersionsUsecase(
	formRepo repository.FormRepository,
) *ListFormTemplateVersionsUsecase {
	return &ListFormTemplateVersionsUsecase{formRepo}
}

func (uc *ListFormTemplateVersionsUsecase) Execute(
	templateID string,
) ([]domain.FormTemplate, error) {
	return uc.formRepo.ListTemplateVersions(templateID)
}
