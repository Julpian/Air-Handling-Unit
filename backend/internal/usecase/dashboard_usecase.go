package usecase

import (
	"context"
	"time"

	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type DashboardUsecase struct {
	scheduleRepo repository.ScheduleRepository
	inspectionRepo repository.InspectionRepository
}

func NewDashboardUsecase(
	scheduleRepo repository.ScheduleRepository,
	inspectionRepo repository.InspectionRepository,
) *DashboardUsecase {
	return &DashboardUsecase{
		scheduleRepo:   scheduleRepo,
		inspectionRepo: inspectionRepo,
	}
}

type DashboardStats struct {
	TodaySchedule int                        `json:"today_schedule"`
	Running       int                        `json:"running"`
	Overdue       int                        `json:"overdue"`
	RunningList   []domain.ScheduleWithDetail `json:"running_schedules"`
}

func (u *DashboardUsecase) GetStats(ctx context.Context) (*DashboardStats, error) {

	today := time.Now()

	schedules, err := u.scheduleRepo.ListWithDetail()
	if err != nil {
		return nil, err
	}

	var todayCount int
	var running int
	var overdue int
	var runningList []domain.ScheduleWithDetail

	for _, s := range schedules {

		if today.After(s.StartDate) && today.Before(s.EndDate) {
			todayCount++
		}

		if s.Status == "dalam_pemeriksaan" {
			running++
			runningList = append(runningList, *s)
		}

		if today.After(s.EndDate) && s.Status != "selesai" {
			overdue++
		}
	}

	return &DashboardStats{
		TodaySchedule: todayCount,
		Running:       running,
		Overdue:       overdue,
		RunningList:   runningList,
	}, nil
}

func (u *DashboardUsecase) GetFilterPressureChart() ([]domain.FilterPressureChartRow, error) {
	return u.inspectionRepo.GetFilterPressureChart()
}

