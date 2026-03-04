package domain

import "time"

type InspectionReport struct {
	InspectionID string
	ScheduleID   string
	Status       string
	UnitCode     string
	Inspector    string
	Signature    string // ✅ SATU SAJA

	SPVName      string
	SPVSignature string
	SPVSignedAt  *time.Time

	AreaName         string // Dari tabel areas
	CleanlinessClass string // Dari tabel ahus
	Vendor           string // Dari tabel ahus
	RoomName         string // Dari tabel ahus (Location)
	Period           string
	InspectedAt      *time.Time

	Sections []InspectionReportSection
}

type InspectionReportSection struct {
	Code  string
	Title string
	Items []InspectionReportItem
}

type InspectionReportItem struct {
	Label  string
	Value  string
	Result string
}
