package domain

import "time"

type ScheduleWithDetail struct {
	ID        string    `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
	NFCBypass bool      `json:"nfc_bypass"` // ✅ TAMBAHKAN INI

	// 🔗 PLAN
	PlanID      string `json:"plan_id"`
	Period      string `json:"period"`
	WeekOfMonth int    `json:"week_of_month"`
	Month       *int   `json:"month"`

	// 👷 Inspector
	InspectorID   *string `json:"inspector_id"`
	InspectorName *string `json:"inspector_name"`

	// 🏭 AHU
	AHUID    string  `json:"ahu_id"`
	UnitCode string  `json:"unit_code"`
	RoomName *string `json:"room_name"`
	NFCUID   *string `json:"nfc_uid"`
}
