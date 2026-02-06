package usecase

import (
	"errors"

	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

type ScheduleAssignUsecase struct {
	scheduleRepo repository.ScheduleRepository
	auditRepo    repository.AuditTrailRepository
}

func NewScheduleAssignUsecase(
	scheduleRepo repository.ScheduleRepository,
	auditRepo repository.AuditTrailRepository,
) *ScheduleAssignUsecase {
	return &ScheduleAssignUsecase{
		scheduleRepo: scheduleRepo,
		auditRepo:    auditRepo,
	}
}

func (u *ScheduleAssignUsecase) AssignInspector(
	scheduleID string,
	inspectorID string,
	adminID string,
) error {

	if _, err := uuid.Parse(inspectorID); err != nil {
		return errors.New("inspector_id tidak valid")
	}

	return u.scheduleRepo.AssignInspector(scheduleID, inspectorID)
}
