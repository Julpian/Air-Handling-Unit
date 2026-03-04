package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type SchedulePlanUsecase struct {
	repo repository.SchedulePlanRepository
}

func NewSchedulePlanUsecase(
	repo repository.SchedulePlanRepository,
) *SchedulePlanUsecase {
	return &SchedulePlanUsecase{
		repo: repo,
	}
}

// ✅ CREATE DENGAN STRUCT
func (u *SchedulePlanUsecase) Create(
	ahuID string,
	period string,
	week int,
	month *int,
) error {

	if ahuID == "" {
		return errors.New("ahu wajib diisi")
	}

	if week < 1 || week > 4 {
		return errors.New("minggu tidak valid")
	}

	switch period {

	case domain.PeriodMonthly: // "bulanan"
		month = nil

	case domain.PeriodSixMonth: // "enam_bulan"
		if month == nil {
			return errors.New("bulan wajib diisi untuk enam_bulan")
		}
		if *month < 1 || *month > 12 {
			return errors.New("bulan tidak valid")
		}

	case domain.PeriodYearly: // "tahunan"
		if month == nil {
			return errors.New("bulan wajib diisi untuk tahunan")
		}
		if *month < 1 || *month > 12 {
			return errors.New("bulan tidak valid")
		}

	default:
		return errors.New("periode tidak valid")
	}

	return u.repo.Create(&domain.SchedulePlan{
		ID:          uuid.NewString(),
		AHUId:       ahuID,
		Period:      period,
		WeekOfMonth: week,
		Month:       month,
	})
}

func (u *SchedulePlanUsecase) ListAllWithAHU() ([]domain.SchedulePlanWithAHU, error) {
	return u.repo.ListAllWithAHU()
}

func (u *SchedulePlanUsecase) Update(
	id string,
	period string,
	week int,
	month *int,
) error {

	if week < 1 || week > 4 {
		return errors.New("minggu tidak valid")
	}

	switch period {

	case domain.PeriodMonthly:
		month = nil

	case domain.PeriodSixMonth, domain.PeriodYearly:
		if month == nil {
			return errors.New("bulan wajib diisi")
		}

	default:
		return errors.New("periode tidak valid")
	}

	return u.repo.Update(&domain.SchedulePlan{
		ID:          id,
		Period:      period,
		WeekOfMonth: week,
		Month:       month,
	})
}

func (u *SchedulePlanUsecase) Delete(id string) error {
	if id == "" {
		return errors.New("id tidak valid")
	}
	return u.repo.Delete(id)
}
