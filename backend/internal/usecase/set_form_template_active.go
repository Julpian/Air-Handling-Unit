package usecase

import "ahu-backend/internal/repository"

type SetFormTemplateActiveUsecase struct {
	formRepo repository.FormRepository
}

func NewSetFormTemplateActiveUsecase(
	formRepo repository.FormRepository,
) *SetFormTemplateActiveUsecase {
	return &SetFormTemplateActiveUsecase{formRepo}
}

func (uc *SetFormTemplateActiveUsecase) Execute(
	id string,
	active bool,
) error {
	return uc.formRepo.SetActive(id, active)
}
