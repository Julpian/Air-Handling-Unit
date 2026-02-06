package usecase

import (
	"errors"
	"time"

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

// ScanNFC memulai pemeriksaan dengan scan NFC
func (u *InspectionUsecase) ScanNFC(
	nfcUID, _ string,
	inspectorID string,
) (*domain.Inspection, error) {

	// 1. Cari AHU dari NFC
	ahu, err := u.ahuRepo.GetByNFCUID(nfcUID)
	if err != nil || ahu == nil {
		return nil, errors.New("NFC tidak terdaftar")
	}

	// 2. Cari schedule aktif dari AHU
	schedule, err := u.scheduleRepo.GetActiveByAHU(ahu.ID)
	if err != nil || schedule == nil {
		return nil, errors.New("tidak ada jadwal aktif")
	}

	// 3. Validasi status
	if schedule.Status != domain.ScheduleStatusSiapDiperiksa {
		return nil, errors.New("jadwal tidak dapat discan")
	}

	// 4. Ambil form template
	form, err := u.formRepo.GetTemplateBySchedule(schedule.ID)
	if err != nil || form == nil {
		return nil, errors.New("form inspeksi tidak ditemukan")
	}

	// 5. Ambil last inspection (jika ada)
	last, err := u.inspectionRepo.GetLastByScheduleID(schedule.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	inspection := &domain.Inspection{
		ID:             uuid.NewString(),
		ScheduleID:     schedule.ID,
		InspectorID:    inspectorID,
		FormTemplateID: form.ID,
		Status:         domain.InspectionStatusSedangDiisi,
		ScannedNFCUID:  &nfcUID,
		InspectedAt:    &now,
	}

	if last != nil && last.Status == "revisi" {
		inspection.ParentID = &last.ID
	}

	if err := u.inspectionRepo.Create(inspection); err != nil {
		return nil, err
	}

	if err := u.scheduleRepo.UpdateStatus(
		schedule.ID,
		domain.ScheduleStatusDalamPemeriksaan,
	); err != nil {
		return nil, err
	}

	return inspection, nil
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

	// Update status inspection
	if err := u.inspectionRepo.UpdateStatus(
		inspection.ID,
		domain.InspectionStatusTerkirim,
		nil,
	); err != nil {
		return err
	}

	// ❗ Schedule JANGAN selesai di sini
	// Schedule menunggu keputusan supervisor

	// Audit
	_ = u.auditRepo.Save(&domain.AuditTrail{
		ID:        uuid.NewString(),
		UserID:    inspectorID,
		Action:    "submit_inspection",
		Entity:    "Inspection",
		EntityID:  inspection.ID,
		CreatedAt: time.Now(),
	})

	return nil
}

func (u *InspectionUsecase) GetByID(id string) (*domain.Inspection, error) {
	return u.inspectionRepo.GetByID(id)
}

func (u *InspectionUsecase) ListByStatus(status string) ([]domain.Inspection, error) {
	return u.inspectionRepo.ListByStatus(status)
}
