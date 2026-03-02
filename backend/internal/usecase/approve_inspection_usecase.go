package usecase

import "ahu-backend/internal/repository"

type ApproveInspectionUsecase struct {
	repo repository.InspectionRepository
}

func NewApproveInspectionUsecase(repo repository.InspectionRepository) *ApproveInspectionUsecase {
	return &ApproveInspectionUsecase{repo: repo}
}

func (uc *ApproveInspectionUsecase) Execute(
	inspectionID string,
	spvID string,
	signature string,
) error {

	return uc.repo.ApproveInspection(inspectionID, spvID, signature)
}
