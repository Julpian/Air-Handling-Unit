package domain

import "time"

type SchedulePlan struct {
	ID          string    `json:"id"`
	AHUId       string    `json:"ahu_id"`
	UnitCode    string    `json:"unit_code"`
	Period      string    `json:"period"`
	WeekOfMonth int       `json:"week_of_month"`
	Month       *int      `json:"month"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type SchedulePlanWithAHU struct {
	SchedulePlan
	UnitCode string `json:"unit_code"` // ✅ GANTI
}
