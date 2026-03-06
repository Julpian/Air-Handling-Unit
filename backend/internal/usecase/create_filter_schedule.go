package usecase

import (
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

type CreateFilterScheduleUsecase struct {
	scheduleRepo repository.ScheduleRepository
	planRepo     repository.SchedulePlanRepository
}

func NewCreateFilterScheduleUsecase(
	scheduleRepo repository.ScheduleRepository,
	planRepo repository.SchedulePlanRepository,
) *CreateFilterScheduleUsecase {
	return &CreateFilterScheduleUsecase{
		scheduleRepo: scheduleRepo,
		planRepo:     planRepo,
	}
}

func (u *CreateFilterScheduleUsecase) Execute(
	ahuID string,
	startDate time.Time,
) error {

	plan, err := u.planRepo.GetByPeriod("ganti_filter", ahuID)
	if err != nil {
		return err
	}

	if plan == nil {
		plan = &domain.SchedulePlan{
			ID:          uuid.NewString(),
			AHUId:       ahuID,
			Period:      "ganti_filter",
			WeekOfMonth: 1,
			CreatedAt:   time.Now(),
		}

		err = u.planRepo.Create(plan)
		if err != nil {
			return err
		}
	}

	endDate := startDate.AddDate(0, 0, 7)

	return u.scheduleRepo.Create(&domain.Schedule{
		ID:        uuid.NewString(),
		PlanID:    plan.ID,   // ✅ pakai plan ID
		AHUId:     ahuID,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    domain.ScheduleStatusSiapDiperiksa,
		NFCBypass: false,
		CreatedAt: time.Now(),
	})
}