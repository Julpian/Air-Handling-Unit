package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type ListFormTemplateUsecase struct {
	formRepo repository.FormRepository
}

func NewListFormTemplateUsecase(
	formRepo repository.FormRepository,
) *ListFormTemplateUsecase {
	return &ListFormTemplateUsecase{
		formRepo: formRepo,
	}
}

func (uc *ListFormTemplateUsecase) Execute() ([]domain.FormTemplate, error) {
	return uc.formRepo.ListTemplates()
}
