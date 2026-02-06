package usecase

import (
	"errors"
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

type InspectionApprovalUsecase struct {
	inspectionRepo repository.InspectionRepository
	scheduleRepo   repository.ScheduleRepository
	auditRepo      repository.AuditTrailRepository
}

func NewInspectionApprovalUsecase(
	inspectionRepo repository.InspectionRepository,
	scheduleRepo repository.ScheduleRepository,
	auditRepo repository.AuditTrailRepository,
) *InspectionApprovalUsecase {
	return &InspectionApprovalUsecase{
		inspectionRepo: inspectionRepo,
		scheduleRepo:   scheduleRepo,
		auditRepo:      auditRepo,
	}
}

func (u *InspectionApprovalUsecase) ApproveInspection(
	inspectionID string,
	approverID string,
) error {

	inspection, err := u.inspectionRepo.GetByID(inspectionID)
	if err != nil || inspection == nil {
		return errors.New("inspection tidak ditemukan")
	}

	if inspection.Status != "terkirim" {
		return errors.New("inspection belum siap approve")
	}

	now := time.Now()

	// 1️⃣ Update inspection
	if err := u.inspectionRepo.Approve(
		inspectionID,
		approverID,
		now,
		map[string]any{
			"signed_at": now,
			"method":    "jwt",
		},
	); err != nil {
		return err
	}

	// 2️⃣ Lock schedule
	_ = u.scheduleRepo.UpdateStatus(inspection.ScheduleID, "selesai")

	// 3️⃣ Audit trail
	_ = u.auditRepo.Save(&domain.AuditTrail{
		ID:       uuid.NewString(),
		UserID:   approverID,
		Action:   "approve_inspection",
		Entity:   "inspection",
		EntityID: inspection.ID,
		Metadata: map[string]any{
			"status": "disetujui",
		},
		CreatedAt: now,
	})

	return nil
}

func (u *InspectionApprovalUsecase) RejectInspection(
	inspectionID string,
	approverID string,
	reason string,
) error {

	inspection, err := u.inspectionRepo.GetByID(inspectionID)
	if err != nil || inspection == nil {
		return errors.New("inspection tidak ditemukan")
	}

	if inspection.Status != "terkirim" {
		return errors.New("inspection belum siap direview")
	}

	// update
	_ = u.inspectionRepo.UpdateStatus(inspectionID, "revisi", &reason)
	_ = u.scheduleRepo.UpdateStatus(inspection.ScheduleID, "siap_diperiksa")

	// audit
	return u.auditRepo.Save(&domain.AuditTrail{
		ID:       uuid.NewString(),
		UserID:   approverID,
		Action:   "reject_inspection",
		Entity:   "inspection",
		EntityID: inspection.ID,
		Metadata: map[string]any{
			"reason": reason,
		},
		CreatedAt: time.Now(),
	})
}
