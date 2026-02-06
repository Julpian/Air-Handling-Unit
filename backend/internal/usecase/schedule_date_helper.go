package usecase

import (
	"time"

	"ahu-backend/internal/domain"
)

// weekRangeOfMonth → menghasilkan 1 minggu (Minggu–Sabtu)
func weekRangeOfMonth(
	year int,
	month time.Month,
	weekOfMonth int,
) (time.Time, time.Time) {

	t := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	// Geser ke Minggu pertama
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}

	// Geser ke minggu ke-N
	t = t.AddDate(0, 0, (weekOfMonth-1)*7)

	start := t
	end := t.AddDate(0, 0, 6)

	// Guard: jangan spill ke bulan lain
	if start.Month() != month {
		return time.Time{}, time.Time{}
	}

	return start, end
}

// buildWeekRanges → CORE generator schedule mingguan
func buildWeekRanges(
	plan *domain.SchedulePlan,
	year int,
) [][2]time.Time {

	var ranges [][2]time.Time
	var months []int

	switch plan.Period {

	case "monthly":
		months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	case "semester":
		if plan.Month == nil {
			return ranges
		}
		m := *plan.Month
		months = []int{m}
		if m+6 <= 12 {
			months = append(months, m+6)
		}

	case "yearly":
		if plan.Month == nil {
			return ranges
		}
		months = []int{*plan.Month}
	}

	for _, m := range months {
		start, end := weekRangeOfMonth(
			year,
			time.Month(m),
			plan.WeekOfMonth,
		)

		if start.IsZero() {
			continue
		}

		ranges = append(ranges, [2]time.Time{start, end})
	}

	return ranges
}
