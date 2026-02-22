package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"fmt"
	"os"
)

type ScheduleApprovalUsecase struct {
	repo repository.ScheduleApprovalRepository
	pdf  *SchedulePDFService
}

func NewScheduleApprovalUsecase(
	repo repository.ScheduleApprovalRepository,
	pdf *SchedulePDFService,
) *ScheduleApprovalUsecase {
	return &ScheduleApprovalUsecase{repo: repo, pdf: pdf}
}

func (u *ScheduleApprovalUsecase) Ensure(year int) {
	if _, err := u.repo.GetByYear(year); err != nil {
		_ = u.repo.Create(year)
	}
}

func (u *ScheduleApprovalUsecase) SignSVP(year int, userID, signature string) error {
	u.Ensure(year)
	return u.repo.SignSVP(year, userID, signature)
}

func (uc *ScheduleApprovalUsecase) SignAsmen(year int, userID, signature string) error {

	// 1️⃣ Simpan tanda tangan ASMEN dulu
	err := uc.repo.SignAsmen(year, userID, signature)
	if err != nil {
		return err
	}

	// 2️⃣ Ambil data approval lengkap (termasuk SVP & ASMEN signature)
	approval, err := uc.repo.GetByYear(year)
	if err != nil {
		return err
	}

	// 3️⃣ Pastikan folder files ada
	if err := os.MkdirAll("files", os.ModePerm); err != nil {
		return err
	}

	path := fmt.Sprintf("files/schedule-%d.pdf", year)

	// 4️⃣ Generate PDF dengan data approval
	if err := uc.pdf.Generate(year, approval, path); err != nil {
		return err
	}

	// 5️⃣ Simpan path PDF ke database
	return uc.repo.SetPDF(year, path)
}

func (u *ScheduleApprovalUsecase) Get(year int) (*domain.ScheduleApproval, error) {
	return u.repo.GetByYear(year)
}

func (u *ScheduleApprovalUsecase) GetByToken(token string) (*domain.ScheduleApproval, error) {
	return u.repo.GetByToken(token)
}
