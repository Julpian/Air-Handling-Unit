package usecase

import (
	"context"
	"errors"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type ScheduleBypassNFCUsecase struct {
	scheduleRepo repository.ScheduleRepository
	auditRepo    repository.AuditTrailRepository
}

func NewScheduleBypassNFCUsecase(
	scheduleRepo repository.ScheduleRepository,
	auditRepo repository.AuditTrailRepository,
) *ScheduleBypassNFCUsecase {
	return &ScheduleBypassNFCUsecase{
		scheduleRepo: scheduleRepo,
		auditRepo:    auditRepo,
	}
}

func (uc *ScheduleBypassNFCUsecase) Execute(
	ctx context.Context,
	scheduleID string,
	enabled bool,
	userID string,
	reason string,
) error {

	schedule, err := uc.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		return err
	}

	if schedule.Status == domain.ScheduleStatusSelesai {
		return errors.New("schedule sudah selesai, tidak bisa bypass NFC")
	}

	schedule.NFCBypass = enabled
	if err := uc.scheduleRepo.Update(schedule); err != nil {
		return err
	}

	// audit trail
	_ = uc.auditRepo.Save(&domain.AuditTrail{
		UserID:   userID,
		Action:   "BYPASS_NFC",
		Entity:   "Schedule",
		EntityID: scheduleID,
	})

	return nil
}
