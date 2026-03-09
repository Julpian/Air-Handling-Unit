package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type SchedulePlanUsecase struct {
	repo    repository.SchedulePlanRepository
	auditUC *AuditTrailUsecase
}

func NewSchedulePlanUsecase(
	repo repository.SchedulePlanRepository,
	auditUC *AuditTrailUsecase,
) *SchedulePlanUsecase {
	return &SchedulePlanUsecase{
		repo:    repo,
		auditUC: auditUC,
	}
}

// ✅ CREATE DENGAN STRUCT
func (u *SchedulePlanUsecase) Create(
	ahuID string,
	period string,
	week int,
	month *int,
	adminID string,
	adminName string,
) error {

	if ahuID == "" {
		return errors.New("ahu wajib diisi")
	}

	if week < 1 || week > 4 {
		return errors.New("minggu tidak valid")
	}

	switch period {

	case domain.PeriodMonthly:
		month = nil

	case domain.PeriodSixMonth:
		if month == nil {
			return errors.New("bulan wajib diisi untuk enam_bulan")
		}

	case domain.PeriodYearly:
		if month == nil {
			return errors.New("bulan wajib diisi untuk tahunan")
		}

	default:
		return errors.New("periode tidak valid")
	}

	plan := &domain.SchedulePlan{
		ID:          uuid.NewString(),
		AHUId:       ahuID,
		Period:      period,
		WeekOfMonth: week,
		Month:       month,
	}

	err := u.repo.Create(plan)

	if err == nil {
		u.auditUC.Log(&domain.AuditTrail{
			UserID:   adminID,
			Action:   "CREATE_SCHEDULE_PLAN",
			Entity:   "SchedulePlan",
			EntityID: plan.ID,
			Metadata: map[string]interface{}{
				"ahu_id":     ahuID,
				"period":     period,
				"week":       week,
				"admin_name": adminName,
			},
		})
	}

	return err
}

func (u *SchedulePlanUsecase) ListAllWithAHU() ([]domain.SchedulePlanWithAHU, error) {
	return u.repo.ListAllWithAHU()
}

func (u *SchedulePlanUsecase) Update(
	id string,
	period string,
	week int,
	month *int,
	adminID string,
	adminName string,
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

	err := u.repo.Update(&domain.SchedulePlan{
		ID:          id,
		Period:      period,
		WeekOfMonth: week,
		Month:       month,
	})

	if err == nil {
		u.auditUC.Log(&domain.AuditTrail{
			UserID:   adminID,
			Action:   "UPDATE_SCHEDULE_PLAN",
			Entity:   "SchedulePlan",
			EntityID: id,
			Metadata: map[string]interface{}{
				"period":     period,
				"week":       week,
				"admin_name": adminName,
			},
		})
	}

	return err
}

func (u *SchedulePlanUsecase) Delete(
	id string,
	adminID string,
	adminName string,
) error {

	if id == "" {
		return errors.New("id tidak valid")
	}

	err := u.repo.Delete(id)

	if err == nil {
		u.auditUC.Log(&domain.AuditTrail{
			UserID:   adminID,
			Action:   "DELETE_SCHEDULE_PLAN",
			Entity:   "SchedulePlan",
			EntityID: id,
			Metadata: map[string]interface{}{
				"admin_name": adminName,
			},
		})
	}

	return err
}
