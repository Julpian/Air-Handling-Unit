package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type InspectionTaskUsecase struct {
	scheduleRepo repository.ScheduleRepository
}

func NewInspectionTaskUsecase(
	scheduleRepo repository.ScheduleRepository,
) *InspectionTaskUsecase {
	return &InspectionTaskUsecase{
		scheduleRepo: scheduleRepo,
	}
}

func (u *InspectionTaskUsecase) ListByInspector(
	inspectorID string,
) ([]domain.InspectionTask, error) {

	return u.scheduleRepo.ListTasksByInspector(inspectorID)
}
