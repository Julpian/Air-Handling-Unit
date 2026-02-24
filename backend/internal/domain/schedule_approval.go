package domain

import "time"

type ScheduleApproval struct {
	ID   string
	Year int

	SVPID        *string
	SVPSignature *string
	SVPSignedAt  *time.Time

	AsmenID        *string
	AsmenSignature *string
	AsmenSignedAt  *time.Time

	PDFPath *string
	PDFHash *string

	Status      string
	VerifyToken *string `json:"verify_token"`

	CreatedAt time.Time
}
