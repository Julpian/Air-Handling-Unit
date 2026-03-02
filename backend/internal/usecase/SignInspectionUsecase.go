package usecase

import (
	"ahu-backend/internal/repository"
	"errors"
)

type SignInspectionUsecase struct {
	repo repository.InspectionRepository
	pdf  *InspectionPDFService
}

func NewSignInspectionUsecase(
	repo repository.InspectionRepository,
	pdf *InspectionPDFService,
) *SignInspectionUsecase {
	return &SignInspectionUsecase{
		repo: repo,
		pdf:  pdf,
	}
}

func (uc *SignInspectionUsecase) Execute(id string, signature string) error {

	if signature == "" {
		return errors.New("signature kosong")
	}

	if err := uc.repo.SaveSignature(id, signature); err != nil {
		return err
	}

	if err := uc.repo.UpdateStatus(id, "waiting_spv", nil); err != nil {
		return err
	}

	// 🔥 GENERATE PDF SETELAH SIGN
	return uc.pdf.GenerateInspectionPDF(id)
}
