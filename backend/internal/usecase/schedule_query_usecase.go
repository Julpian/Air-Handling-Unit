package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"ahu-backend/internal/usecase/dto"
)

type ScheduleQueryUsecase struct {
	repo repository.ScheduleRepository
}

func NewScheduleQueryUsecase(
	repo repository.ScheduleRepository,
) *ScheduleQueryUsecase {
	return &ScheduleQueryUsecase{repo: repo}
}

func (u *ScheduleQueryUsecase) List() ([]*dto.ScheduleListDTO, error) {
	data, err := u.repo.ListWithDetail()
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ScheduleListDTO, 0, len(data))

	for _, d := range data {
		result = append(result, &dto.ScheduleListDTO{
			ID:        d.ID,
			StartDate: d.StartDate,
			EndDate:   d.EndDate,
			Status:    d.Status,

			PlanID:      d.PlanID,
			Period:      d.Period,
			WeekOfMonth: d.WeekOfMonth,
			Month:       d.Month,

			InspectorID:   d.InspectorID,
			InspectorName: d.InspectorName,

			AHUID:    d.AHUID,
			UnitCode: d.UnitCode,
			RoomName: d.RoomName,
			NFCUID:   d.NFCUID,
		})
	}

	return result, nil
}

func (u *ScheduleQueryUsecase) ListByYear(year int) ([]*domain.ScheduleWithDetail, error) {
	return u.repo.ListByYear(year)
}
