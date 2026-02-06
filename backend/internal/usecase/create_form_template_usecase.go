package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"context"
	"errors"
)

type CreateFormTemplateUsecase struct {
	formRepo repository.FormRepository
}

func NewCreateFormTemplateUsecase(
	formRepo repository.FormRepository,
) *CreateFormTemplateUsecase {
	return &CreateFormTemplateUsecase{
		formRepo: formRepo,
	}
}

func (uc *CreateFormTemplateUsecase) Execute(
	ctx context.Context,
	template *domain.FormTemplate,
) error {

	// ===============================
	// VALIDATION (BUSINESS RULE)
	// ===============================
	if template.Name == "" {
		return errors.New("nama form wajib diisi")
	}

	if template.Period == "" {
		return errors.New("periode form wajib diisi")
	}

	if len(template.Sections) == 0 {
		return errors.New("minimal harus ada 1 section")
	}

	for _, sec := range template.Sections {
		if sec.Title == "" {
			return errors.New("judul section tidak boleh kosong")
		}

		if len(sec.Items) == 0 {
			return errors.New("section harus memiliki minimal 1 item")
		}

		for _, item := range sec.Items {
			if item.Label == "" {
				return errors.New("label item tidak boleh kosong")
			}
			if item.InputType == "" {
				return errors.New("tipe input wajib diisi")
			}
		}
	}

	// ===============================
	// EXECUTE CREATE
	// ===============================
	return uc.formRepo.CreateTemplate(ctx, template)
}
