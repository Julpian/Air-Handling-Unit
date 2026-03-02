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

	ScanToken        *string
	ScanTokenExpires *time.Time

	ScannedNFCUID *string    // UID NFC yang discan
	InspectedAt   *time.Time // waktu scan NFC dilakukan

	InspectorSignature *string
	SPVSignature       *string
	SPVSignedAt        *time.Time
	ApprovedBy         *string

	ParentID  *string // jika re-inspection
	CreatedAt time.Time
}
