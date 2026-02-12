package usecase

import (
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

// InspectionUsecase berisi logika pemeriksaan AHU
type InspectionUsecase struct {
	ahuRepo        repository.AHURepository
	scheduleRepo   repository.ScheduleRepository
	inspectionRepo repository.InspectionRepository
	formRepo       repository.FormRepository
	auditRepo      repository.AuditTrailRepository
}

func NewInspectionUsecase(
	ahuRepo repository.AHURepository,
	scheduleRepo repository.ScheduleRepository,
	inspectionRepo repository.InspectionRepository,
	formRepo repository.FormRepository,
	auditRepo repository.AuditTrailRepository,
) *InspectionUsecase {
	return &InspectionUsecase{
		ahuRepo:        ahuRepo,
		scheduleRepo:   scheduleRepo,
		inspectionRepo: inspectionRepo,
		formRepo:       formRepo,
		auditRepo:      auditRepo,
	}
}

// SubmitInspection menyelesaikan pemeriksaan
func (u *InspectionUsecase) SubmitInspection(
	inspectionID string,
	results []domain.InspectionResult,
	inspectorID string,
) error {

	inspection, err := u.inspectionRepo.GetByID(inspectionID)
	if err != nil || inspection == nil {
		return errors.New("inspection tidak ditemukan")
	}

	if inspection.Status != domain.InspectionStatusSedangDiisi {
		return errors.New("inspection tidak dalam status yang valid")
	}

	for _, r := range results {
		r.ID = uuid.NewString()
		r.InspectionID = inspectionID

		if err := u.inspectionRepo.SaveResult(&r); err != nil {
			return err
		}
	}

	if err := u.inspectionRepo.UpdateStatus(
		inspectionID,
		domain.InspectionStatusTerkirim,
		nil,
	); err != nil {
		return err
	}

	if err := u.scheduleRepo.UpdateStatus(
		inspection.ScheduleID,
		domain.ScheduleStatusSelesai,
	); err != nil {
		return err
	}

	return nil
}

func (u *InspectionUsecase) GetByID(id string) (*domain.Inspection, error) {
	return u.inspectionRepo.GetByID(id)
}

func (u *InspectionUsecase) ListByStatus(status string) ([]domain.Inspection, error) {
	return u.inspectionRepo.ListByStatus(status)
}
