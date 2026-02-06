package usecase

import (
	"context"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type CreateNewFormTemplateVersionUsecase struct {
	formRepo repository.FormRepository
}

func NewCreateNewFormTemplateVersionUsecase(
	formRepo repository.FormRepository,
) *CreateNewFormTemplateVersionUsecase {
	return &CreateNewFormTemplateVersionUsecase{
		formRepo: formRepo,
	}
}

func (uc *CreateNewFormTemplateVersionUsecase) Execute(
	ctx context.Context,
	oldTemplateID string,
	template *domain.FormTemplate,
) error {
	return uc.formRepo.CreateNewVersion(ctx, oldTemplateID, template)
}
