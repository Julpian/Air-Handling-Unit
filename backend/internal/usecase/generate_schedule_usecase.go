package usecase

import (
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"

	"github.com/google/uuid"
)

type GenerateScheduleUsecase struct {
	planRepo     repository.SchedulePlanRepository
	scheduleRepo repository.ScheduleRepository
}

func NewGenerateScheduleUsecase(
	planRepo repository.SchedulePlanRepository,
	scheduleRepo repository.ScheduleRepository,
) *GenerateScheduleUsecase {
	return &GenerateScheduleUsecase{
		planRepo:     planRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (u *GenerateScheduleUsecase) Generate(year int) error {
	plans, err := u.planRepo.GetActiveByYear(year)
	if err != nil {
		return err
	}

	for i := range plans {
		plan := &plans[i]

		// 1️⃣ hapus jadwal lama
		if err := u.scheduleRepo.DeleteByPlan(plan.ID); err != nil {
			return err
		}

		ranges := buildWeekRanges(plan, year)

		for _, r := range ranges {
			err := u.scheduleRepo.Create(&domain.Schedule{
				ID:        uuid.NewString(),
				PlanID:    plan.ID,
				AHUId:     plan.AHUId,
				StartDate: r[0],
				EndDate:   r[1],
				Status:    domain.ScheduleStatusSiapDiperiksa,
				NFCBypass: false,
				CreatedAt: time.Now(),
			})
			if err != nil {
				return err
			}
		}

		// 3️⃣ MARK PLAN AS GENERATED
		if err := u.planRepo.MarkGenerated(plan.ID); err != nil {
			return err
		}
	}

	return nil
}
