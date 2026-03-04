package dto

type InspectionVerificationResponse struct {
	InspectionID  string `json:"inspection_id"`
	UnitCode      string `json:"unit_code"`
	Period        string `json:"period"`
	InspectorName string `json:"inspector_name"`
	SPVName       string `json:"spv_name"`
	InspectedAt   string `json:"inspected_at"`
}
