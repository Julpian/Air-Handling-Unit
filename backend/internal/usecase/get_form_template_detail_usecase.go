package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type GetFormTemplateDetailUsecase struct {
	formRepo repository.FormRepository
}

func NewGetFormTemplateDetailUsecase(
	formRepo repository.FormRepository,
) *GetFormTemplateDetailUsecase {
	return &GetFormTemplateDetailUsecase{
		formRepo: formRepo,
	}
}

func (uc *GetFormTemplateDetailUsecase) Execute(
	templateID string,
) (*domain.FormTemplate, error) {

	if templateID == "" {
		return nil, errors.New("template id wajib diisi")
	}

	return uc.formRepo.GetTemplateByID(templateID)
}
