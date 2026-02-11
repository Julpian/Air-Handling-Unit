package domain

const (
	PeriodMonthly  = "bulanan"
	PeriodSixMonth = "enam_bulan"
	PeriodYearly   = "tahunan"
)

// =======================
// INSPECTION STATUS
// =======================

const (
	InspectionStatusSedangDiisi = "sedang_diisi"
	InspectionStatusTerkirim    = "terkirim"
	InspectionStatusDisetujui   = "disetujui"
	InspectionStatusRevisi      = "revisi"
)

// =======================
// SCHEDULE STATUS
// =======================

const (
	ScheduleStatusSiapDiperiksa    = "siap_diperiksa"
	ScheduleStatusDalamPemeriksaan = "dalam_pemeriksaan"
	ScheduleStatusSelesai          = "selesai"
	ScheduleStatusRevisi           = "revisi"
)

// =======================
// AUDIT ACTION
// =======================

const (
	AuditActionScanNFC           = "scan_nfc"
	AuditActionSubmitInspection  = "submit_inspection"
	AuditActionApproveInspection = "approve_inspection"
	AuditActionRejectInspection  = "reject_inspection"
)
