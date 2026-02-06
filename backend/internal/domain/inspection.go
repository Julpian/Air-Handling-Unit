package domain

import "time"

// Inspection merepresentasikan satu proses pemeriksaan AHU
type Inspection struct {
	ID          string
	ScheduleID  string
	InspectorID string

	FormTemplateID string // ⬅️ TAMBAHAN (penting)

	Status string  // pending | inspected | approved | rejected
	Note   *string // catatan jika reject / remark

	ScannedNFCUID *string    // UID NFC yang discan
	InspectedAt   *time.Time // waktu scan NFC dilakukan

	ParentID  *string // jika re-inspection
	CreatedAt time.Time
}
